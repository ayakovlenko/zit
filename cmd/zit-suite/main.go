package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Spec struct {
	Name    string     `yaml:"name"`
	Version string     `yaml:"version"`
	Tests   []TestCase `yaml:"tests"`
}

type TestCase struct {
	ID     string   `yaml:"id"`
	Argv   []string `yaml:"argv"`
	Setup  Setup    `yaml:"setup"`
	Expect Expect   `yaml:"expect"`
}

type Setup struct {
	GitInit         bool              `yaml:"git_init"`
	GitRemote       string            `yaml:"git_remote"`
	ConfigFile      string            `yaml:"config_file"`
	GitGlobalConfig string            `yaml:"git_global_config"`
	GitSystemConfig string            `yaml:"git_system_config"`
	ExistingFile    *ExistingFile     `yaml:"existing_file"`
	Env             map[string]string `yaml:"env"`
}

type ExistingFile struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}

type Expect struct {
	Stdout             string        `yaml:"stdout"`
	StdoutContainsAll  []string      `yaml:"stdout_contains_all"`
	Stderr             string        `yaml:"stderr"`
	StderrContains     string        `yaml:"stderr_contains"`
	Exit               int           `yaml:"exit"`
	GitConfig          []GitKV       `yaml:"git_config"`
	GitConfigUnchanged bool          `yaml:"git_config_unchanged"`
	FsExists           []string      `yaml:"fs_exists"`
	FsContent          []FsAssertion `yaml:"fs_content"`
}

type GitKV struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type FsAssertion struct {
	Path        string `yaml:"path"`
	NotContains string `yaml:"not_contains"`
}

// Fixed paths inside every test container.
const (
	containerHome      = "/home/zit"
	containerRepoDir   = "/home/zit/repo"
	containerGitGlobal = "/home/zit/git-global-config"
	containerGitSystem = "/home/zit/git-system-config"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: zit-suite <image>")
		os.Exit(2)
	}
	image := os.Args[1]

	data, err := os.ReadFile("spec.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading spec: %v\n", err)
		os.Exit(2)
	}

	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing spec: %v\n", err)
		os.Exit(2)
	}

	passed, failed := 0, 0
	for _, tc := range spec.Tests {
		failures := runTest(tc, image)
		if len(failures) == 0 {
			passed++
		} else {
			failed++
			fmt.Printf("--- FAIL: %s\n", tc.ID)
			for _, f := range failures {
				fmt.Printf("    %s\n", strings.ReplaceAll(f, "\n", "\n    "))
			}
		}
	}

	status := "ok "
	if failed > 0 {
		status = "FAIL"
	}
	fmt.Printf("%s  %d passed, %d failed\n", status, passed, failed)

	if failed > 0 {
		os.Exit(1)
	}
}

func runTest(tc TestCase, image string) []string {
	cid := "zit-suite-" + tc.ID

	configPath := resolveConfigPath(containerHome, tc.Setup.Env)

	// CWD inside the container: repo subdir when git_init, home otherwise.
	cwd := containerHome
	if tc.Setup.GitInit {
		cwd = containerRepoDir
	}

	var failures []string
	fail := func(format string, args ...any) {
		failures = append(failures, fmt.Sprintf(format, args...))
	}

	// Build env for docker create.
	env := buildContainerEnv(configPath, tc.Setup.Env)

	// Create the container.
	createArgs := []string{"create", "--name", cid}
	for _, kv := range env {
		createArgs = append(createArgs, "-e", kv)
	}
	createArgs = append(createArgs, image, "sleep", "infinity")
	if out, err := dockerCmd(createArgs...); err != nil {
		return []string{fmt.Sprintf("setup: docker create: %v\n%s", err, out)}
	}
	defer dockerCmd("rm", "-f", cid) //nolint:errcheck

	// Start it.
	if out, err := dockerCmd("start", cid); err != nil {
		return []string{fmt.Sprintf("setup: docker start: %v\n%s", err, out)}
	}

	// Write git global and system config files.
	if err := containerWriteFile(cid, containerGitGlobal, tc.Setup.GitGlobalConfig); err != nil {
		return []string{fmt.Sprintf("setup: write git global config: %v", err)}
	}
	if err := containerWriteFile(cid, containerGitSystem, tc.Setup.GitSystemConfig); err != nil {
		return []string{fmt.Sprintf("setup: write git system config: %v", err)}
	}

	// Git init.
	if tc.Setup.GitInit {
		if out, err := containerExec(cid, "git", "init", "--initial-branch=main", containerRepoDir); err != nil {
			// Older git versions don't support --initial-branch; retry without.
			if out2, err2 := containerExec(cid, "git", "init", containerRepoDir); err2 != nil {
				return []string{fmt.Sprintf("setup: git init: %v\n%s", err2, out2)}
			}
			_ = out
		}
		_, _ = containerExec(cid, "git", "-C", containerRepoDir, "config", "user.email", "test@example.com")
		_, _ = containerExec(cid, "git", "-C", containerRepoDir, "config", "user.name", "Test")
		if tc.Setup.GitRemote != "" {
			if out, err := containerExec(cid, "git", "-C", containerRepoDir, "remote", "add", "origin", tc.Setup.GitRemote); err != nil {
				return []string{fmt.Sprintf("setup: git remote add: %v\n%s", err, out)}
			}
		}
	}

	// Write zit config file.
	if tc.Setup.ConfigFile != "" {
		if err := containerWriteFile(cid, configPath, tc.Setup.ConfigFile); err != nil {
			return []string{fmt.Sprintf("setup: write config file: %v", err)}
		}
	}

	// Write existing_file.
	if tc.Setup.ExistingFile != nil {
		p := substituteVars(tc.Setup.ExistingFile.Path, containerHome, cwd, configPath)
		if err := containerWriteFile(cid, p, tc.Setup.ExistingFile.Content); err != nil {
			return []string{fmt.Sprintf("setup: write existing_file: %v", err)}
		}
	}

	// Snapshot .git/config before running zit, for git_config_unchanged check.
	var gitConfigSnapshot string
	if tc.Expect.GitConfigUnchanged && tc.Setup.GitInit {
		gitConfigSnapshot, _, _ = dockerExec(cid, "cat", containerRepoDir+"/.git/config")
	}

	// Run zit. cd into cwd first; redirect stdout/stderr to temp files so we
	// can capture them separately, then read them back.
	zitCmd := shellQuote(append([]string{"zit"}, tc.Argv...)...)
	script := fmt.Sprintf("cd %s && %s > /tmp/zit-stdout 2> /tmp/zit-stderr; echo $? > /tmp/zit-exit", shellEscape(cwd), zitCmd)
	containerExec(cid, "sh", "-c", script) //nolint:errcheck -- exit code comes from the file

	gotStdout, _, _ := dockerExec(cid, "cat", "/tmp/zit-stdout")
	gotStderr, _, _ := dockerExec(cid, "cat", "/tmp/zit-stderr")
	exitStr, _, _ := dockerExec(cid, "cat", "/tmp/zit-exit")
	exitCode := 0
	fmt.Sscanf(strings.TrimSpace(exitStr), "%d", &exitCode)

	wantStdout := substituteVars(tc.Expect.Stdout, containerHome, cwd, configPath)
	wantStderr := substituteVars(tc.Expect.Stderr, containerHome, cwd, configPath)
	wantStderrContains := substituteVars(tc.Expect.StderrContains, containerHome, cwd, configPath)

	if exitCode != tc.Expect.Exit {
		fail("exit code mismatch:\n  want: %d\n  got:  %d\n  stderr: %q", tc.Expect.Exit, exitCode, gotStderr)
	}
	if wantStdout != "" && gotStdout != wantStdout {
		fail("stdout mismatch:\n  want: %q\n  got:  %q", wantStdout, gotStdout)
	} else if wantStdout == "" && tc.Expect.Stdout == "" && len(tc.Expect.StdoutContainsAll) == 0 && gotStdout != "" {
		fail("stdout mismatch:\n  want: (empty)\n  got:  %q", gotStdout)
	}
	for _, s := range tc.Expect.StdoutContainsAll {
		if !strings.Contains(gotStdout, s) {
			fail("stdout does not contain %q:\n  got: %q", s, gotStdout)
		}
	}
	if wantStderrContains != "" {
		if !strings.Contains(gotStderr, wantStderrContains) {
			fail("stderr does not contain %q:\n  got: %q", wantStderrContains, gotStderr)
		}
	} else if wantStderr != "" && gotStderr != wantStderr {
		fail("stderr mismatch:\n  want: %q\n  got:  %q", wantStderr, gotStderr)
	} else if wantStderr == "" && tc.Expect.Stderr == "" && gotStderr != "" {
		fail("stderr mismatch:\n  want: (empty)\n  got:  %q", gotStderr)
	}

	// git_config assertions.
	for _, kv := range tc.Expect.GitConfig {
		out, _, code := dockerExec(cid, "git", "-C", containerRepoDir, "config", "--local", kv.Key)
		out = strings.TrimSpace(out)
		if code != 0 {
			fail("git config --local %s: exit %d", kv.Key, code)
			continue
		}
		if out != kv.Value {
			fail("git config %s:\n  want: %q\n  got:  %q", kv.Key, kv.Value, out)
		}
	}

	// git_config_unchanged assertion.
	if tc.Expect.GitConfigUnchanged && tc.Setup.GitInit {
		after, _, _ := dockerExec(cid, "cat", containerRepoDir+"/.git/config")
		if gitConfigSnapshot != after {
			fail("git config changed when it should not have")
		}
	}

	// fs_exists assertions.
	for _, p := range tc.Expect.FsExists {
		p = substituteVars(p, containerHome, cwd, configPath)
		_, _, code := dockerExec(cid, "test", "-e", p)
		if code != 0 {
			fail("fs_exists: %s does not exist", p)
		}
	}

	// fs_content assertions.
	for _, a := range tc.Expect.FsContent {
		p := substituteVars(a.Path, containerHome, cwd, configPath)
		content, _, code := dockerExec(cid, "cat", p)
		if code != 0 {
			fail("fs_content: cannot read %s", p)
			continue
		}
		if a.NotContains != "" && strings.Contains(content, a.NotContains) {
			fail("fs_content: %s contains %q but should not", p, a.NotContains)
		}
	}

	return failures
}

// dockerCmd runs a top-level docker command (not exec) and returns combined output.
func dockerCmd(args ...string) (string, error) {
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// containerExec runs a command inside cid and returns combined output + error.
func containerExec(cid string, args ...string) (string, error) {
	out, _, err2 := dockerExec(cid, args...)
	var err error
	if err2 != 0 {
		err = fmt.Errorf("exit %d", err2)
	}
	return out, err
}

// dockerExec runs docker exec <cid> <args...> and returns stdout, stderr, exit code.
func dockerExec(cid string, args ...string) (stdout, stderr string, exitCode int) {
	dockerArgs := append([]string{"exec", cid}, args...)
	cmd := exec.Command("docker", dockerArgs...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	_ = cmd.Run()
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	return outBuf.String(), errBuf.String(), exitCode
}

// containerWriteFile writes content to path inside the container using sh -c with stdin.
func containerWriteFile(cid, path, content string) error {
	dir := filepath.Dir(path)
	script := fmt.Sprintf("mkdir -p %s && cat > %s", shellEscape(dir), shellEscape(path))
	cmd := exec.Command("docker", "exec", "-i", cid, "sh", "-c", script)
	cmd.Stdin = strings.NewReader(content)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}

// shellEscape wraps a string in single quotes, escaping any single quotes within.
func shellEscape(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// shellQuote joins args into a shell command string with each arg single-quoted.
func shellQuote(args ...string) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = shellEscape(a)
	}
	return strings.Join(parts, " ")
}

// buildContainerEnv builds the env slice for docker create.
func buildContainerEnv(configPath string, setupEnv map[string]string) []string {
	env := []string{
		"HOME=" + containerHome,
		"GIT_CONFIG_GLOBAL=" + containerGitGlobal,
		"GIT_CONFIG_SYSTEM=" + containerGitSystem,
		"GIT_CONFIG_NOSYSTEM=1",
	}
	for k, v := range setupEnv {
		if v == "" {
			continue
		}
		v = strings.ReplaceAll(v, "{{HOME}}", containerHome)
		v = strings.ReplaceAll(v, "{{CONFIG_PATH}}", configPath)
		env = append(env, k+"="+v)
	}
	return env
}

// resolveConfigPath returns the path the config file should be written to and
// read from, following the spec search order. Does not require the file to exist.
func resolveConfigPath(home string, setupEnv map[string]string) string {
	sub := func(s string) string { return strings.ReplaceAll(s, "{{HOME}}", home) }

	defaultPath := filepath.Join(home, ".config", "zit", "config.yaml")
	if v := setupEnv["XDG_CONFIG_HOME"]; v != "" {
		defaultPath = filepath.Join(sub(v), "zit", "config.yaml")
	}

	if v := setupEnv["ZIT_CONFIG"]; v != "" {
		v = sub(v)
		v = strings.ReplaceAll(v, "{{CONFIG_PATH}}", defaultPath)
		return v
	}
	return defaultPath
}

func substituteVars(s, home, cwd, configPath string) string {
	return substituteVarsRaw(s, home, cwd, configPath)
}

func substituteVarsRaw(s, home, cwd, configPath string) string {
	s = strings.ReplaceAll(s, "{{HOME}}", home)
	if cwd != "" {
		s = strings.ReplaceAll(s, "{{CWD}}", cwd)
	}
	if configPath != "" {
		s = strings.ReplaceAll(s, "{{CONFIG_PATH}}", configPath)
	}
	return s
}
