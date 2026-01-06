use once_cell::sync::OnceCell;
use std::fs::{self, File, OpenOptions};
use std::io::Write;
use std::path::PathBuf;
use std::sync::Mutex;
use tracing_subscriber::fmt::format::FmtSpan;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::util::SubscriberInitExt;

static LOG_ENABLED: OnceCell<bool> = OnceCell::new();
static LOG_DIR: OnceCell<PathBuf> = OnceCell::new();
static VUE_LOG_FILE: OnceCell<Mutex<File>> = OnceCell::new();

fn get_timestamp() -> String {
    chrono::Local::now().format("%Y-%m-%d-%H-%M-%S").to_string()
}

pub fn init_with_dir(_log_dir: &PathBuf) {
    let args: Vec<String> = std::env::args().collect();
    println!("[Logger] args: {:?}", args);
    let has_d = args.iter().any(|a| a == "-d");
    let has_b = args.iter().any(|a| a == "-b");
    let enabled = has_d && !has_b;
    println!("[Logger] has_d: {}, has_b: {}, enabled: {}", has_d, has_b, enabled);
    LOG_ENABLED.set(enabled).ok();
    if !enabled {
        return;
    }
    let exe_dir = std::env::current_exe()
        .and_then(|p| p.canonicalize())
        .map(|p| p.parent().unwrap_or(&p).to_path_buf())
        .unwrap_or_default();
    let exe_dir_str = exe_dir.to_string_lossy();
    let logs_dir = if exe_dir_str.contains("target\\debug") || exe_dir_str.contains("target/debug")
        || exe_dir_str.contains("target\\release") || exe_dir_str.contains("target/release") {
        exe_dir.join("..").join("..").join("..").join("runtime").join("logs")
    } else {
        exe_dir.join("..").join("logs")
    };
    println!("[Logger] exe_dir: {:?}", exe_dir);
    println!("[Logger] logs_dir: {:?}", logs_dir);
    if let Err(e) = fs::create_dir_all(&logs_dir) {
        println!("[Logger] Failed to create logs dir: {:?}", e);
    }
    LOG_DIR.set(logs_dir.clone()).ok();
    let tauri_log_path = logs_dir.join(format!("tauri-{}.log", get_timestamp()));
    let vue_log_path = logs_dir.join(format!("vue-{}.log", get_timestamp()));
    println!("[Logger] tauri_log_path: {:?}", tauri_log_path);
    println!("[Logger] vue_log_path: {:?}", vue_log_path);
    match File::create(&tauri_log_path) {
        Ok(file) => {
            let file_layer = tracing_subscriber::fmt::layer()
                .with_writer(Mutex::new(file))
                .with_ansi(false);
            let stdout_layer = tracing_subscriber::fmt::layer()
                .with_span_events(FmtSpan::CLOSE)
                .with_target(true)
                .with_level(true);
            tracing_subscriber::registry()
                .with(file_layer)
                .with(stdout_layer)
                .init();
            tracing::info!("Logger initialized, log file: {:?}", tauri_log_path);
        }
        Err(e) => {
            println!("[Logger] Failed to create tauri log file: {:?}", e);
        }
    }
    match OpenOptions::new().create(true).append(true).open(&vue_log_path) {
        Ok(file) => {
            VUE_LOG_FILE.set(Mutex::new(file)).ok();
            println!("[Logger] Vue log file created successfully");
        }
        Err(e) => {
            println!("[Logger] Failed to create vue log file: {:?}", e);
        }
    }
}

pub fn is_enabled() -> bool {
    *LOG_ENABLED.get().unwrap_or(&false)
}

pub fn write_vue_log(level: &str, message: &str) {
    if !is_enabled() {
        return;
    }
    if let Some(file) = VUE_LOG_FILE.get() {
        if let Ok(mut f) = file.lock() {
            let timestamp = chrono::Local::now().format("%Y-%m-%d %H:%M:%S%.3f");
            writeln!(f, "{} [{}] {}", timestamp, level, message).ok();
            f.flush().ok();
        }
    }
}
