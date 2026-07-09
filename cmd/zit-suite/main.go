package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

func main() {
	specPath := flag.String("spec", "spec.yaml", "path to spec.yaml")
	runFilter := flag.String("run", "", "regex filter on test ID")
	verbose := flag.Bool("v", false, "print output of passing tests")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: zit-suite [--spec <path>] [--run <regex>] [-v] <binary> [binary-args...]")
		os.Exit(2)
	}
	binaryTokens := flag.Args()

	// Resolve the binary to an absolute path so it can be found when
	// cmd.Dir is set to the test's temp directory.
	if resolved, err := exec.LookPath(binaryTokens[0]); err == nil {
		if abs, err := filepath.Abs(resolved); err == nil {
			binaryTokens[0] = abs
		}
	}

	var filter *regexp.Regexp
	if *runFilter != "" {
		var err error
		filter, err = regexp.Compile(*runFilter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid --run regex: %v\n", err)
			os.Exit(2)
		}
	}

	data, err := os.ReadFile(*specPath)
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
		if filter != nil && !filter.MatchString(tc.ID) {
			continue
		}
		failures := runTest(tc, binaryTokens)
		if len(failures) == 0 {
			passed++
			if *verbose {
				fmt.Printf("--- PASS: %s\n", tc.ID)
			}
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

func runTest(tc TestCase, binaryTokens []string) []string {
	tmpHome, err := os.MkdirTemp("", "zit-suite-*")
	if err != nil {
		return []string{fmt.Sprintf("setup: MkdirTemp: %v", err)}
	}
	defer os.RemoveAll(tmpHome)

	// Resolve symlinks so {{HOME}} matches what the OS reports to the binary.
	if real, err := filepath.EvalSymlinks(tmpHome); err == nil {
		tmpHome = real
	}

	var failures []string
	fail := func(format string, args ...any) {
		failures = append(failures, fmt.Sprintf(format, args...))
	}

	// CWD defaults to tmpHome; overridden when git_init is true.
	cwd := tmpHome

	if tc.Setup.GitInit {
		repoDir := filepath.Join(tmpHome, "repo")
		if err := os.MkdirAll(repoDir, 0755); err != nil {
			return []string{fmt.Sprintf("setup: mkdir repo: %v", err)}
		}
		if _, err := gitSetup(repoDir, "init", "--initial-branch=main"); err != nil {
			// older git versions don't support --initial-branch; retry without
			if _, err2 := gitSetup(repoDir, "init"); err2 != nil {
				return []string{fmt.Sprintf("setup: git init: %v", err2)}
			}
		}
		_, _ = gitSetup(repoDir, "config", "user.email", "test@example.com")
		_, _ = gitSetup(repoDir, "config", "user.name", "Test")
		if tc.Setup.GitRemote != "" {
			if _, err := gitSetup(repoDir, "remote", "add", "origin", tc.Setup.GitRemote); err != nil {
				return []string{fmt.Sprintf("setup: git remote add: %v", err)}
			}
		}
		cwd = repoDir
	}

	globalCfgPath := filepath.Join(tmpHome, "git-global-config")
	systemCfgPath := filepath.Join(tmpHome, "git-system-config")
	if err := os.WriteFile(globalCfgPath, []byte(tc.Setup.GitGlobalConfig), 0644); err != nil {
		return []string{fmt.Sprintf("setup: write git global config: %v", err)}
	}
	if err := os.WriteFile(systemCfgPath, []byte(tc.Setup.GitSystemConfig), 0644); err != nil {
		return []string{fmt.Sprintf("setup: write git system config: %v", err)}
	}

	configPath := resolveConfigPath(tmpHome, tc.Setup.Env)

	if tc.Setup.ConfigFile != "" {
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return []string{fmt.Sprintf("setup: mkdir config dir: %v", err)}
		}
		if err := os.WriteFile(configPath, []byte(tc.Setup.ConfigFile), 0644); err != nil {
			return []string{fmt.Sprintf("setup: write config file: %v", err)}
		}
	}

	if tc.Setup.ExistingFile != nil {
		p := substituteVars(tc.Setup.ExistingFile.Path, tmpHome, cwd, configPath)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			return []string{fmt.Sprintf("setup: mkdir existing_file dir: %v", err)}
		}
		if err := os.WriteFile(p, []byte(tc.Setup.ExistingFile.Content), 0644); err != nil {
			return []string{fmt.Sprintf("setup: write existing_file: %v", err)}
		}
	}

	var gitConfigSnapshot []byte
	if tc.Expect.GitConfigUnchanged && tc.Setup.GitInit {
		gitConfigSnapshot, _ = os.ReadFile(filepath.Join(cwd, ".git", "config"))
	}

	env := buildEnv(tmpHome, globalCfgPath, systemCfgPath, configPath, tc.Setup.Env)

	args := append(append([]string{}, binaryTokens[1:]...), tc.Argv...)
	cmd := exec.Command(binaryTokens[0], args...)
	cmd.Dir = cwd
	cmd.Env = env

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	_ = cmd.Run()

	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	wantStdout := substituteVars(tc.Expect.Stdout, tmpHome, cwd, configPath)
	wantStderr := substituteVars(tc.Expect.Stderr, tmpHome, cwd, configPath)
	wantStderrContains := substituteVars(tc.Expect.StderrContains, tmpHome, cwd, configPath)

	gotStdout := stdoutBuf.String()
	gotStderr := stderrBuf.String()

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

	for _, kv := range tc.Expect.GitConfig {
		out, err := gitLocal(cwd, "config", "--local", kv.Key)
		if err != nil {
			fail("git config --local %s: %v", kv.Key, err)
			continue
		}
		if out != kv.Value {
			fail("git config %s:\n  want: %q\n  got:  %q", kv.Key, kv.Value, out)
		}
	}

	if tc.Expect.GitConfigUnchanged && tc.Setup.GitInit {
		after, _ := os.ReadFile(filepath.Join(cwd, ".git", "config"))
		if !bytes.Equal(gitConfigSnapshot, after) {
			fail("git config changed when it should not have")
		}
	}

	for _, p := range tc.Expect.FsExists {
		p = substituteVars(p, tmpHome, cwd, configPath)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			fail("fs_exists: %s does not exist", p)
		}
	}

	for _, a := range tc.Expect.FsContent {
		p := substituteVars(a.Path, tmpHome, cwd, configPath)
		content, err := os.ReadFile(p)
		if err != nil {
			fail("fs_content: cannot read %s: %v", p, err)
			continue
		}
		if a.NotContains != "" && strings.Contains(string(content), a.NotContains) {
			fail("fs_content: %s contains %q but should not", p, a.NotContains)
		}
	}

	return failures
}

// gitSetup runs a git command during test setup with a clean env so ambient
// GIT_CONFIG_GLOBAL / GIT_CONFIG_SYSTEM don't interfere.
func gitSetup(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_CONFIG_NOSYSTEM=1",
	)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// gitLocal reads a git config value from a repo directory after the test has
// run, using the ambient env (which has the real git config).
func gitLocal(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_CONFIG_NOSYSTEM=1",
	)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// resolveConfigPath returns the path the config file should be written to and
// read from, following the spec search order. Does not require the file to exist.
func resolveConfigPath(home string, setupEnv map[string]string) string {
	// Substitute {{HOME}} in env values before using them.
	sub := func(s string) string { return strings.ReplaceAll(s, "{{HOME}}", home) }

	// ZIT_CONFIG takes priority, but it may contain {{CONFIG_PATH}} which
	// means "same as the default path" -- resolve the default first and
	// substitute back.
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

func buildEnv(home, globalCfg, systemCfg, configPath string, setupEnv map[string]string) []string {
	env := []string{
		"HOME=" + home,
		"PATH=" + os.Getenv("PATH"),
		"GIT_CONFIG_GLOBAL=" + globalCfg,
		"GIT_CONFIG_SYSTEM=" + systemCfg,
		"GIT_CONFIG_NOSYSTEM=1",
	}
	for k, v := range setupEnv {
		if v == "" {
			continue
		}
		// Substitute {{HOME}} and {{CONFIG_PATH}} in env values.
		v = strings.ReplaceAll(v, "{{HOME}}", home)
		v = strings.ReplaceAll(v, "{{CONFIG_PATH}}", configPath)
		env = append(env, k+"="+v)
	}
	return env
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
