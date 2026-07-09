use std::collections::HashMap;
use std::error::Error;
use std::fs;
use std::path::{Path, PathBuf};
use std::process::{Command, exit};

use clap::{CommandFactory, Parser, Subcommand};
use serde::Deserialize;

const APP_VERSION: &str = "v3.1.2";
const CONFIG_ENV_VAR: &str = "ZIT_CONFIG";
const XDG_CONFIG_HOME_VAR: &str = "XDG_CONFIG_HOME";

type Result<T> = std::result::Result<T, Box<dyn Error>>;

#[derive(Parser)]
#[command(
    name = "zit",
    about = "git identity manager",
    disable_help_subcommand = false,
)]
struct Cli {
    #[command(subcommand)]
    command: Option<Cmd>,
}

#[derive(Subcommand)]
enum Cmd {
    /// Print version
    Version,
    /// Check git setup for potential problems
    Doctor,
    /// Set git identity
    Set {
        #[arg(long, help = "run without applying configurations")]
        dry_run: bool,
    },
    /// Configuration management
    Config {
        #[command(subcommand)]
        sub: ConfigCmd,
    },
}

#[derive(Subcommand)]
enum ConfigCmd {
    /// Initialize configuration file
    Init,
    /// Show path to configuration file
    Path,
    /// Show configuration file contents
    Show,
}

#[derive(Debug, Deserialize)]
struct ConfigRoot {
    hosts: HashMap<String, HostConfig>,
}

#[derive(Debug, Deserialize)]
struct HostConfig {
    default: Option<User>,
    #[serde(default)]
    overrides: Vec<Override>,
}

#[derive(Debug, Deserialize, Clone)]
struct User {
    name: String,
    email: String,
    #[serde(rename = "sign")]
    signing: Option<Signing>,
}

#[derive(Debug, Deserialize, Clone)]
struct Signing {
    key: String,
    format: String,
}

#[derive(Debug, Deserialize)]
struct Override {
    owner: String,
    #[serde(default)]
    repo: String,
    user: User,
}

fn home_dir() -> Result<PathBuf> {
    std::env::var("HOME")
        .map(PathBuf::from)
        .map_err(|_| "HOME is not set".into())
}

fn xdg_config_path(home_dir: &Path) -> PathBuf {
    match std::env::var(XDG_CONFIG_HOME_VAR) {
        Ok(v) if !v.is_empty() => PathBuf::from(v).join("zit").join("config.yaml"),
        _ => home_dir.join(".config").join("zit").join("config.yaml"),
    }
}

fn locate_config(home_dir: &Path) -> Result<PathBuf> {
    if let Ok(v) = std::env::var(CONFIG_ENV_VAR) {
        if !v.is_empty() {
            let p = PathBuf::from(&v);
            if p.is_file() {
                return Ok(p);
            }
            return Err(format!(
                "config file defined in {CONFIG_ENV_VAR} variable is not found at '{v}'"
            )
            .into());
        }
    }

    let xdg = xdg_config_path(home_dir);
    if xdg.is_file() {
        return Ok(xdg);
    }

    let dot = home_dir.join(".zit").join("config.yaml");
    if dot.is_file() {
        return Ok(dot);
    }

    Err(format!(
        "config file is not found at neither {} nor {}",
        xdg.display(),
        dot.display()
    )
    .into())
}

fn load_config(path: &Path) -> Result<ConfigRoot> {
    let contents = fs::read_to_string(path)
        .map_err(|e| format!("cannot read config: {e}"))?;
    yaml_serde::from_str(&contents).map_err(|e| format!("cannot parse config: {e}").into())
}

struct GitOutput {
    stdout: String,
    exit_code: i32,
}

fn git(args: &[&str]) -> Result<GitOutput> {
    let o = Command::new("git")
        .args(args)
        .output()
        .map_err(|e| format!("failed to run git: {e}"))?;
    Ok(GitOutput {
        stdout: String::from_utf8_lossy(&o.stdout).trim().to_string(),
        exit_code: o.status.code().unwrap_or(-1),
    })
}

fn is_git_dir() -> Result<bool> {
    let o = git(&["rev-parse", "--is-inside-work-tree"])?;
    Ok(o.exit_code == 0 && o.stdout == "true")
}

fn git_config_get(scope: &str, key: &str) -> Result<Option<String>> {
    let o = git(&["config", scope, key])?;
    Ok(if o.exit_code == 0 && !o.stdout.is_empty() {
        Some(o.stdout)
    } else {
        None
    })
}

fn git_set_local(key: &str, value: &str) -> Result<()> {
    let o = git(&["config", "--local", key, value])?;
    if o.exit_code != 0 {
        return Err(format!("set local config {key}={value}: git exited {}", o.exit_code).into());
    }
    Ok(())
}

fn git_remote_url(name: &str) -> Result<String> {
    let o = git(&["remote", "get-url", name])?;
    match o.exit_code {
        0 => Ok(o.stdout),
        2 => Err(format!("remote \"{name}\" is not set").into()),
        128 => Err(o.stdout.into()),
        c => Err(format!("git remote get-url {name}: exited {c}").into()),
    }
}

struct RepoInfo {
    host: String,
    owner: String,
    name: String,
}

fn parse_remote_url(url: &str) -> Result<RepoInfo> {
    use git_url_parse::GitUrl;
    use git_url_parse::types::provider::GenericProvider;

    let parsed = GitUrl::parse(url).map_err(|e| format!("cannot parse remote url: {e}"))?;
    let host = parsed.host().unwrap_or("").to_string();
    let provider: GenericProvider = parsed
        .provider_info()
        .map_err(|e| format!("cannot extract repo info: {e}"))?;
    Ok(RepoInfo {
        host,
        owner: provider.owner().to_string(),
        name: provider.repo().to_string(),
    })
}

fn find_best_match<'a>(host_conf: &'a HostConfig, repo: &RepoInfo) -> Option<&'a User> {
    let result = host_conf.default.as_ref();

    for ov in &host_conf.overrides {
        if !ov.repo.is_empty() {
            if ov.owner == repo.owner && ov.repo == repo.name {
                return Some(&ov.user);
            }
        }
        if ov.owner == repo.owner {
            return Some(&ov.user);
        }
    }

    result
}

fn cmd_version() {
    println!("{APP_VERSION}");
}

fn cmd_doctor() -> Result<()> {
    let checks: &[(&str, fn() -> Result<bool>)] = &[
        (
            "git config --global user.useConfigOnly true",
            || Ok(git_config_get("--global", "user.useConfigOnly")?.as_deref() == Some("true")),
        ),
        (
            "git config --unset-all --global user.name",
            || Ok(git_config_get("--global", "user.name")?.is_none()),
        ),
        (
            "git config --unset-all --global user.email",
            || Ok(git_config_get("--global", "user.email")?.is_none()),
        ),
        (
            "git config --unset-all --system user.name",
            || Ok(git_config_get("--system", "user.name")?.is_none()),
        ),
        (
            "git config --unset-all --system user.email",
            || Ok(git_config_get("--system", "user.email")?.is_none()),
        ),
    ];

    let lines = checks
        .iter()
        .map(|(name, check)| {
            let tick = if check()? { "x" } else { " " };
            Ok(format!("- [{tick}] {name}"))
        })
        .collect::<Result<Vec<_>>>()?;

    println!("{}", lines.join("\n"));
    Ok(())
}

fn cmd_set(dry_run: bool) -> Result<()> {
    if !is_git_dir()? {
        let cwd = std::env::current_dir().unwrap_or_default();
        // Printed to stderr in the Go impl but the spec expects it on stderr too --
        // however the set_not_a_git_dir test checks stderr, so we must use eprintln
        // and then exit rather than return Err (which would double-print via main).
        eprintln!(
            "Error: {:?} is not a git directory\n\nMake sure you are executing zit inside a git directory.\n\nIf you are, perhaps you have forgotten to initialize a new repository? In this\ncase, run:\n\n    git init\n\nOr, if you have an existing repository but haven't set up the remote URL:\n\n    git remote add origin <url>",
            cwd.display()
        );
        exit(1);
    }

    let home = home_dir()?;
    let config_path = locate_config(&home)?;
    let conf = load_config(&config_path)?;

    let remote_url = match git_remote_url("origin") {
        Ok(u) => u,
        Err(e) if e.to_string().contains("is not set") => {
            println!(
                "Error: {e}\n\nAdd remote URL so that zit could use it for choosing the correct git identity as\ndefined in the configuration file:\n\ngit remote add origin <url>"
            );
            exit(1);
        }
        Err(e) => return Err(e),
    };

    let repo_info = parse_remote_url(&remote_url)?;

    let host_conf = conf
        .hosts
        .get(&repo_info.host)
        .ok_or_else(|| format!("cannot find config for host {:?}", repo_info.host))?;

    let user = find_best_match(host_conf, &repo_info)
        .ok_or_else(|| format!("cannot find a match for host {:?}", repo_info.host))?;

    if !dry_run {
        git_set_local("user.name", &user.name)?;
        git_set_local("user.email", &user.email)?;
    }
    println!("set user: {} <{}>", user.name, user.email);

    if let Some(sign) = &user.signing {
        if !dry_run {
            git_set_local("commit.gpgsign", "true")?;
            git_set_local("tag.gpgsign", "true")?;
            git_set_local("user.signingKey", &sign.key)?;
            git_set_local("gpg.format", &sign.format)?;
        }
        println!("set signing key: {} key at {}", sign.format, sign.key);
    }

    Ok(())
}

const SAMPLE_CONFIG: &str = r#"---
users:
  work: &work_user
    name: "John Doe"
    email: "john.doe@corp.com"

  personal:
    github_user: &personal_github_user
      name: "JD42"
      email: "JD42@users.noreply.github.com"
      sign:
        key: "~/.ssh/id_ed25519_github.pub"
        format: "ssh"

    gitlab_user: &personal_gitlab_user
      name: "JD42"
      email: "786972-JD42@users.noreply.gitlab.com"

hosts:

  github.com:
    default: *personal_github_user
    overrides:
      - owner: "corp"
        user: *work_user

  gitlab.com:
    default: *personal_gitlab_user
"#;

fn cmd_config_init() -> Result<()> {
    let home = home_dir()?;
    let path = xdg_config_path(&home);
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent)
            .map_err(|e| format!("cannot create config dir: {e}"))?;
    }
    fs::write(&path, SAMPLE_CONFIG).map_err(|e| format!("cannot write config: {e}").into())
}

fn cmd_config_path() -> Result<()> {
    let home = home_dir()?;
    match locate_config(&home) {
        Ok(p) => {
            println!("{}", p.display());
            Ok(())
        }
        Err(_) => {
            eprint!("error: config not found; run:\n\n    zit config init\n");
            exit(1);
        }
    }
}

fn cmd_config_show() -> Result<()> {
    let home = home_dir()?;
    match locate_config(&home) {
        Ok(p) => {
            let contents = fs::read_to_string(&p)
                .map_err(|e| format!("cannot read config: {e}"))?;
            print!("{contents}");
            Ok(())
        }
        Err(_) => {
            eprint!("error: config not found; run:\n\n    zit config init\n");
            exit(1);
        }
    }
}

fn run() -> Result<()> {
    let cli = Cli::parse();

    match cli.command {
        None => { let _ = Cli::command().print_help(); },
        Some(Cmd::Version) => cmd_version(),
        Some(Cmd::Doctor) => cmd_doctor()?,
        Some(Cmd::Set { dry_run }) => cmd_set(dry_run)?,
        Some(Cmd::Config { sub }) => match sub {
            ConfigCmd::Init => cmd_config_init()?,
            ConfigCmd::Path => cmd_config_path()?,
            ConfigCmd::Show => cmd_config_show()?,
        },
    }

    Ok(())
}

fn main() {
    if let Err(e) = run() {
        eprintln!("{e}");
        exit(1);
    }
}
