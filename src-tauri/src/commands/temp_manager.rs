use std::fs;
use std::path::PathBuf;
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use tauri::Manager;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct TempFileInfo {
    pub filename: String,
    #[serde(rename = "type")]
    pub file_type: String,
    pub id: String,
    pub size: i64,
    pub created_at: String,
    pub last_modified: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct SessionRecoveryInfo {
    pub id: String,
    #[serde(rename = "type")]
    pub session_type: String,
    pub name: String,
    pub last_modified: String,
    pub preview: String,
    pub file_path: String,
}

fn get_temp_dir(app: &tauri::AppHandle) -> Result<PathBuf, String> {
    let app_data_dir = app.path().app_data_dir().map_err(|e| e.to_string())?;
    let temp_dir = app_data_dir.join("temp");
    fs::create_dir_all(&temp_dir).map_err(|e| e.to_string())?;
    Ok(temp_dir)
}

#[tauri::command]
pub fn get_temp_dir_path(app: tauri::AppHandle) -> Result<String, String> {
    let temp_dir = get_temp_dir(&app)?;
    Ok(temp_dir.to_string_lossy().to_string())
}

#[tauri::command]
pub fn save_temp_file(app: tauri::AppHandle, filename: String, content: String) -> Result<(), String> {
    let temp_dir = get_temp_dir(&app)?;
    let file_path = temp_dir.join(&filename);
    fs::write(&file_path, &content).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn load_temp_file(app: tauri::AppHandle, filename: String) -> Result<String, String> {
    let temp_dir = get_temp_dir(&app)?;
    let file_path = temp_dir.join(&filename);
    fs::read_to_string(&file_path).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_temp_file(app: tauri::AppHandle, filename: String) -> Result<(), String> {
    let temp_dir = get_temp_dir(&app)?;
    let file_path = temp_dir.join(&filename);
    if file_path.exists() {
        fs::remove_file(&file_path).map_err(|e| e.to_string())?;
    }
    Ok(())
}

#[tauri::command]
pub fn list_temp_files(app: tauri::AppHandle) -> Result<Vec<TempFileInfo>, String> {
    let temp_dir = get_temp_dir(&app)?;
    let mut files = Vec::new();
    let entries = fs::read_dir(&temp_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if !path.is_file() {
            continue;
        }
        let filename = path.file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("")
            .to_string();
        let metadata = fs::metadata(&path).ok();
        let size = metadata.as_ref().map(|m| m.len() as i64).unwrap_or(0);
        let modified = metadata.as_ref()
            .and_then(|m| m.modified().ok())
            .map(|t| DateTime::<Utc>::from(t).format("%Y-%m-%dT%H:%M:%SZ").to_string())
            .unwrap_or_default();
        let (file_type, id) = if filename.starts_with("creation_") {
            let parts: Vec<&str> = filename.trim_end_matches(".json").split('_').collect();
            let file_id = parts.get(1).unwrap_or(&"").to_string();
            ("creation".to_string(), file_id)
        } else if filename.starts_with("worldtree_") {
            let parts: Vec<&str> = filename.trim_end_matches(".json").split('_').collect();
            let file_id = parts.get(1).unwrap_or(&"").to_string();
            ("worldtree".to_string(), file_id)
        } else if filename.starts_with("session_") {
            ("session".to_string(), String::new())
        } else {
            (String::new(), String::new())
        };
        files.push(TempFileInfo {
            filename,
            file_type,
            id,
            size,
            created_at: modified.clone(),
            last_modified: modified,
        });
    }
    Ok(files)
}

#[tauri::command]
pub fn cleanup_temp_files(app: tauri::AppHandle, older_than_minutes: i64) -> Result<i32, String> {
    let temp_dir = get_temp_dir(&app)?;
    let now = Utc::now();
    let mut count = 0;
    let entries = fs::read_dir(&temp_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if !path.is_file() {
            continue;
        }
        if let Ok(metadata) = fs::metadata(&path) {
            if let Ok(modified) = metadata.modified() {
                let modified_time = DateTime::<Utc>::from(modified);
                let age = now.signed_duration_since(modified_time);
                if age.num_minutes() > older_than_minutes {
                    if fs::remove_file(&path).is_ok() {
                        count += 1;
                    }
                }
            }
        }
    }
    Ok(count)
}

#[tauri::command]
pub fn get_unsaved_sessions(app: tauri::AppHandle) -> Result<Vec<SessionRecoveryInfo>, String> {
    let temp_dir = get_temp_dir(&app)?;
    let mut sessions = Vec::new();
    let entries = fs::read_dir(&temp_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if !path.is_file() {
            continue;
        }
        let filename = path.file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("")
            .to_string();
        if !filename.ends_with(".json") {
            continue;
        }
        let metadata = fs::metadata(&path).ok();
        let modified = metadata.as_ref()
            .and_then(|m| m.modified().ok())
            .map(|t| DateTime::<Utc>::from(t).format("%Y-%m-%dT%H:%M:%SZ").to_string())
            .unwrap_or_default();
        let content = fs::read_to_string(&path).ok();
        let temp_data: Option<serde_json::Value> = content
            .as_ref()
            .and_then(|c| serde_json::from_str(c).ok());
        let session = if filename.starts_with("creation_") {
            let parts: Vec<&str> = filename.trim_end_matches(".json").split('_').collect();
            let id = parts.get(1).unwrap_or(&"").to_string();
            let name = temp_data.as_ref()
                .and_then(|d| d.get("name"))
                .and_then(|v| v.as_str())
                .unwrap_or("未命名创作")
                .to_string();
            let preview = temp_data.as_ref()
                .and_then(|d| d.get("description"))
                .and_then(|v| v.as_str())
                .map(|s| if s.len() > 100 { format!("{}...", &s[..100]) } else { s.to_string() })
                .unwrap_or_default();
            SessionRecoveryInfo {
                id,
                session_type: "creation".to_string(),
                name,
                last_modified: modified,
                preview,
                file_path: path.to_string_lossy().to_string(),
            }
        } else if filename.starts_with("worldtree_") {
            let parts: Vec<&str> = filename.trim_end_matches(".json").split('_').collect();
            let id = parts.get(1).unwrap_or(&"").to_string();
            let name = temp_data.as_ref()
                .and_then(|d| d.get("name"))
                .and_then(|v| v.as_str())
                .unwrap_or("未命名世界树")
                .to_string();
            SessionRecoveryInfo {
                id,
                session_type: "worldtree".to_string(),
                name,
                last_modified: modified,
                preview: String::new(),
                file_path: path.to_string_lossy().to_string(),
            }
        } else {
            continue;
        };
        sessions.push(session);
    }
    sessions.sort_by(|a, b| b.last_modified.cmp(&a.last_modified));
    Ok(sessions)
}

#[tauri::command]
pub fn recover_temp_session(file_path: String) -> Result<String, String> {
    fs::read_to_string(&file_path).map_err(|e| format!("读取临时文件失败: {}", e))
}

#[tauri::command]
pub fn discard_temp_session(file_path: String) -> Result<(), String> {
    fs::remove_file(&file_path).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn discard_all_temp_sessions(app: tauri::AppHandle) -> Result<(), String> {
    let temp_dir = get_temp_dir(&app)?;
    let entries = fs::read_dir(&temp_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if !path.is_file() {
            continue;
        }
        if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
            if name.starts_with("creation_") || name.starts_with("worldtree_") {
                let _ = fs::remove_file(&path);
            }
        }
    }
    Ok(())
}

#[tauri::command]
pub fn save_creation_temp(app: tauri::AppHandle, creation_id: String, content: String) -> Result<(), String> {
    let temp_dir = get_temp_dir(&app)?;
    let entries = fs::read_dir(&temp_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
            if name.starts_with(&format!("creation_{}_", creation_id)) {
                let _ = fs::remove_file(&path);
            }
        }
    }
    let filename = format!("creation_{}_{}.json", creation_id, chrono::Utc::now().timestamp_nanos_opt().unwrap_or(0));
    let file_path = temp_dir.join(&filename);
    fs::write(&file_path, &content).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn save_world_tree_temp(app: tauri::AppHandle, project_id: String, content: String) -> Result<(), String> {
    let temp_dir = get_temp_dir(&app)?;
    let entries = fs::read_dir(&temp_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
            if name.starts_with(&format!("worldtree_{}_", project_id)) {
                let _ = fs::remove_file(&path);
            }
        }
    }
    let filename = format!("worldtree_{}_{}.json", project_id, chrono::Utc::now().timestamp_nanos_opt().unwrap_or(0));
    let file_path = temp_dir.join(&filename);
    fs::write(&file_path, &content).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn clear_creation_temp(app: tauri::AppHandle, creation_id: String) -> Result<(), String> {
    let temp_dir = get_temp_dir(&app)?;
    let entries = fs::read_dir(&temp_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
            if name.starts_with(&format!("creation_{}_", creation_id)) {
                let _ = fs::remove_file(&path);
            }
        }
    }
    Ok(())
}

#[tauri::command]
pub fn clear_world_tree_temp(app: tauri::AppHandle, project_id: String) -> Result<(), String> {
    let temp_dir = get_temp_dir(&app)?;
    let entries = fs::read_dir(&temp_dir).map_err(|e| e.to_string())?;
    for entry in entries.flatten() {
        let path = entry.path();
        if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
            if name.starts_with(&format!("worldtree_{}_", project_id)) {
                let _ = fs::remove_file(&path);
            }
        }
    }
    Ok(())
}
