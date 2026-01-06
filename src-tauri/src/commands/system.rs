use std::fs;
use std::path::PathBuf;
use tauri::Manager;
use tauri_plugin_dialog::DialogExt;

#[tauri::command]
pub fn hide_window(app: tauri::AppHandle) -> Result<(), String> {
    if let Some(window) = app.get_webview_window("main") {
        window.hide().map_err(|e| e.to_string())?;
    }
    Ok(())
}

#[tauri::command]
pub fn show_window(app: tauri::AppHandle) -> Result<(), String> {
    if let Some(window) = app.get_webview_window("main") {
        window.show().map_err(|e| e.to_string())?;
        window.set_focus().map_err(|e| e.to_string())?;
    }
    Ok(())
}

#[tauri::command]
pub fn is_debug_mode() -> bool {
    cfg!(debug_assertions)
}

#[tauri::command]
pub async fn open_file_dialog(app: tauri::AppHandle, title: String, filters: Vec<(String, Vec<String>)>) -> Result<Option<String>, String> {
    let mut dialog = app.dialog().file().set_title(&title);
    for (name, extensions) in filters {
        let ext_refs: Vec<&str> = extensions.iter().map(|s| s.as_str()).collect();
        dialog = dialog.add_filter(&name, &ext_refs);
    }
    let result = dialog.blocking_pick_file();
    Ok(result.and_then(|f| f.as_path().map(|p| p.to_string_lossy().to_string())))
}

#[tauri::command]
pub async fn save_file_dialog(app: tauri::AppHandle, title: String, default_name: String, filters: Vec<(String, Vec<String>)>) -> Result<Option<String>, String> {
    let mut dialog = app.dialog().file().set_title(&title).set_file_name(&default_name);
    for (name, extensions) in filters {
        let ext_refs: Vec<&str> = extensions.iter().map(|s| s.as_str()).collect();
        dialog = dialog.add_filter(&name, &ext_refs);
    }
    let result = dialog.blocking_save_file();
    Ok(result.and_then(|f| f.as_path().map(|p| p.to_string_lossy().to_string())))
}

fn get_exe_dir() -> Result<PathBuf, String> {
    std::env::current_exe()
        .map(|p| p.parent().unwrap_or(&p).to_path_buf())
        .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_safe_mode_txt_files() -> Result<Vec<String>, String> {
    let exe_dir = get_exe_dir()?;
    let text_dir = exe_dir.join("text");
    if !text_dir.exists() {
        return Ok(vec![]);
    }
    let mut txt_files = Vec::new();
    let entries = fs::read_dir(&text_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_file() {
            if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
                if name.to_lowercase().ends_with(".txt") {
                    txt_files.push(name.to_string());
                }
            }
        }
    }
    Ok(txt_files)
}

#[tauri::command]
pub fn read_safe_mode_template(filename: String) -> Result<String, String> {
    let exe_dir = get_exe_dir()?;
    let file_path = exe_dir.join("text").join(&filename);
    fs::read_to_string(&file_path).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn read_prompt_file(filename: String) -> Result<String, String> {
    let exe_dir = get_exe_dir()?;
    let file_path = exe_dir.join("..").join("..").join("prompts").join(&filename);
    fs::read_to_string(&file_path).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn write_prompt_file(filename: String, content: String) -> Result<(), String> {
    let exe_dir = get_exe_dir()?;
    let prompts_dir = exe_dir.join("..").join("..").join("prompts");
    fs::create_dir_all(&prompts_dir).map_err(|e| e.to_string())?;
    let file_path = prompts_dir.join(&filename);
    fs::write(&file_path, &content).map_err(|e| e.to_string())
}

#[cfg(target_os = "windows")]
#[tauri::command]
pub fn is_caps_lock_on() -> bool {
    use windows_sys::Win32::UI::Input::KeyboardAndMouse::GetKeyState;
    const VK_CAPITAL: i32 = 0x14;
    unsafe { (GetKeyState(VK_CAPITAL) & 1) != 0 }
}

#[cfg(not(target_os = "windows"))]
#[tauri::command]
pub fn is_caps_lock_on() -> bool {
    false
}
