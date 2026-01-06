use crate::database;
use crate::services::llm_api;
use crate::types::{ChatRequest, ChatResponse};
use chrono::Utc;
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;
use tauri::Manager;

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct LocalCreation {
    pub folder_name: String,
    pub folder_path: String,
    pub name: String,
    pub description: String,
    #[serde(default)]
    pub remote_id: Option<String>,
    pub created_at: String,
    pub updated_at: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CreationConfig {
    #[serde(default)]
    pub name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default)]
    pub language: String,
    #[serde(default)]
    pub summary: String,
    #[serde(default)]
    pub world_book: Vec<serde_json::Value>,
    #[serde(default)]
    pub cg_book: Vec<serde_json::Value>,
    #[serde(default)]
    pub pre_text: String,
    #[serde(default)]
    pub post_text: String,
    #[serde(default)]
    pub is_anonymous: bool,
    #[serde(default)]
    pub banned_words: Vec<String>,
    #[serde(default)]
    pub pre_prompt: Option<String>,
    #[serde(default)]
    pub cover: Option<String>,
    #[serde(default)]
    pub cover_tiny: Option<String>,
    #[serde(default)]
    pub opening_statement: Option<String>,
    #[serde(default)]
    pub suggested_questions: Vec<String>,
    #[serde(default)]
    pub bg_image: Option<String>,
    #[serde(default)]
    pub category: i32,
    #[serde(default)]
    pub remote_id: Option<String>,
    #[serde(default)]
    pub created_at: String,
    #[serde(default)]
    pub updated_at: String,
    #[serde(flatten)]
    pub extra: std::collections::HashMap<String, serde_json::Value>,
}

fn get_creations_dir(app: &tauri::AppHandle) -> Result<PathBuf, String> {
    let exe_path = std::env::current_exe().map_err(|e| e.to_string())?;
    let exe_dir = exe_path.parent().ok_or("Cannot get exe directory")?;
    let creations_dir = exe_dir.join("creations");
    fs::create_dir_all(&creations_dir).map_err(|e| e.to_string())?;
    Ok(creations_dir)
}

#[tauri::command]
pub fn get_local_creations(app: tauri::AppHandle) -> Result<Vec<LocalCreation>, String> {
    let creations_dir = get_creations_dir(&app)?;
    let mut creations = Vec::new();
    if let Ok(entries) = fs::read_dir(&creations_dir) {
        for entry in entries.flatten() {
            if entry.path().is_dir() {
                let config_path = entry.path().join("config.json");
                if config_path.exists() {
                    if let Ok(content) = fs::read_to_string(&config_path) {
                        if let Ok(config) = serde_json::from_str::<CreationConfig>(&content) {
                            let metadata = fs::metadata(&config_path).ok();
                            let updated_at = metadata
                                .as_ref()
                                .and_then(|m| m.modified().ok())
                                .map(|t| chrono::DateTime::<chrono::Utc>::from(t).format("%Y-%m-%dT%H:%M:%SZ").to_string())
                                .unwrap_or_default();
                            creations.push(LocalCreation {
                                folder_name: entry.file_name().to_string_lossy().to_string(),
                                folder_path: entry.path().to_string_lossy().to_string(),
                                name: config.name,
                                description: config.description,
                                remote_id: config.remote_id,
                                created_at: config.created_at,
                                updated_at,
                            });
                        }
                    }
                }
            }
        }
    }
    Ok(creations)
}

#[tauri::command]
pub fn get_local_creation_config(app: tauri::AppHandle, folder_name: String) -> Result<serde_json::Value, String> {
    let creations_dir = get_creations_dir(&app)?;
    let config_path = creations_dir.join(&folder_name).join("config.json");
    let content = fs::read_to_string(&config_path).map_err(|e| e.to_string())?;
    serde_json::from_str(&content).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_local_creation_by_remote_id(app: tauri::AppHandle, remote_id: String) -> Result<Option<serde_json::Value>, String> {
    let creations_dir = get_creations_dir(&app)?;
    if let Ok(entries) = fs::read_dir(&creations_dir) {
        for entry in entries.flatten() {
            if entry.path().is_dir() {
                let config_path = entry.path().join("config.json");
                if config_path.exists() {
                    if let Ok(content) = fs::read_to_string(&config_path) {
                        if let Ok(mut config) = serde_json::from_str::<serde_json::Value>(&content) {
                            if let Some(rid) = config.get("remote_id").and_then(|v| v.as_str()) {
                                if rid == remote_id {
                                    if let Some(obj) = config.as_object_mut() {
                                        obj.insert("folder_name".to_string(), serde_json::Value::String(entry.file_name().to_string_lossy().to_string()));
                                        obj.insert("folder_path".to_string(), serde_json::Value::String(entry.path().to_string_lossy().to_string()));
                                    }
                                    return Ok(Some(config));
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    Ok(None)
}

#[tauri::command]
pub fn save_local_creation(app: tauri::AppHandle, folder_name: String, config: serde_json::Value) -> Result<(), String> {
    let creations_dir = get_creations_dir(&app)?;
    let creation_dir = creations_dir.join(&folder_name);
    fs::create_dir_all(&creation_dir).map_err(|e| e.to_string())?;
    let pages_dir = creation_dir.join("pages");
    fs::create_dir_all(&pages_dir).map_err(|e| e.to_string())?;
    let config_path = creation_dir.join("config.json");
    let content = serde_json::to_string_pretty(&config).map_err(|e| e.to_string())?;
    fs::write(&config_path, content).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn rename_local_creation(app: tauri::AppHandle, old_name: String, new_name: String) -> Result<(), String> {
    let creations_dir = get_creations_dir(&app)?;
    let old_path = creations_dir.join(&old_name);
    let new_path = creations_dir.join(&new_name);
    if new_path.exists() {
        return Err("目标文件夹已存在".to_string());
    }
    fs::rename(&old_path, &new_path).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_local_creation(app: tauri::AppHandle, folder_name: String) -> Result<(), String> {
    let creations_dir = get_creations_dir(&app)?;
    let creation_dir = creations_dir.join(&folder_name);
    if creation_dir.exists() {
        fs::remove_dir_all(&creation_dir).map_err(|e| e.to_string())?;
    }
    Ok(())
}

#[tauri::command]
pub fn save_local_creation_page(app: tauri::AppHandle, folder_name: String, page_name: String, content: String) -> Result<(), String> {
    let creations_dir = get_creations_dir(&app)?;
    let pages_dir = creations_dir.join(&folder_name).join("pages");
    fs::create_dir_all(&pages_dir).map_err(|e| e.to_string())?;
    let page_path = pages_dir.join(format!("{}.html", page_name));
    fs::write(&page_path, content).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_local_creation_page(app: tauri::AppHandle, folder_name: String, page_name: String) -> Result<String, String> {
    let creations_dir = get_creations_dir(&app)?;
    let page_path = creations_dir.join(&folder_name).join("pages").join(format!("{}.html", page_name));
    fs::read_to_string(&page_path).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn create_new_local_creation(app: tauri::AppHandle, name: String) -> Result<String, String> {
    let creations_dir = get_creations_dir(&app)?;
    let folder_name = format!("{}_{}", name, Utc::now().timestamp_nanos_opt().unwrap_or(0));
    let project_dir = creations_dir.join(&folder_name);
    fs::create_dir_all(&project_dir).map_err(|e| e.to_string())?;
    let pages_dir = project_dir.join("pages");
    fs::create_dir_all(&pages_dir).map_err(|e| e.to_string())?;
    let now = Utc::now().format("%Y-%m-%dT%H:%M:%SZ").to_string();
    let config = serde_json::json!({
        "name": name,
        "description": "",
        "language": "",
        "summary": "",
        "world_book": [],
        "cg_book": [],
        "pre_text": "",
        "post_text": "",
        "is_anonymous": false,
        "banned_words": [],
        "pre_prompt": null,
        "cover": null,
        "cover_tiny": null,
        "opening_statement": null,
        "suggested_questions": [],
        "suggested_questions_after": {"enabled": false},
        "bg_image": null,
        "bg_mobile": null,
        "category": 0,
        "builtInCss": null,
        "mod_permission": 4,
        "shortcut_commands": [],
        "is_available_not_public": true,
        "preset_type": 1,
        "preset_chats": [],
        "tags": [],
        "created_at": now,
        "updated_at": now
    });
    let config_path = project_dir.join("config.json");
    let content = serde_json::to_string_pretty(&config).map_err(|e| e.to_string())?;
    fs::write(&config_path, content).map_err(|e| e.to_string())?;
    Ok(folder_name)
}

#[tauri::command]
pub async fn fetch_with_auth(token: String, url: String, method: String, body: Option<String>) -> Result<String, String> {
    let client = reqwest::Client::new();
    let mut request = match method.to_uppercase().as_str() {
        "GET" => client.get(&url),
        "POST" => client.post(&url),
        "PUT" => client.put(&url),
        "DELETE" => client.delete(&url),
        "PATCH" => client.patch(&url),
        _ => return Err(format!("Unsupported method: {}", method)),
    };
    request = request
        .header("Authorization", format!("Bearer {}", token))
        .header("Accept", "*/*")
        .header("Origin", "https://aipornhub.ltd")
        .header("Referer", "https://aipornhub.ltd/")
        .header("x-language", "zh-Hans");
    if let Some(body_content) = body {
        request = request
            .header("Content-Type", "application/json")
            .body(body_content);
    }
    let resp = request.send().await.map_err(|e| e.to_string())?;
    resp.text().await.map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn get_app_info(token: String, app_id: String) -> Result<serde_json::Value, String> {
    let url = format!("https://aipornhub.ltd/go/api/apps/{}", app_id);
    let client = reqwest::Client::new();
    let resp = client.get(&url)
        .header("Content-Type", "application/json")
        .header("Accept", "*/*")
        .header("Authorization", format!("Bearer {}", token))
        .header("Origin", "https://aipornhub.ltd")
        .header("Referer", "https://aipornhub.ltd/")
        .header("x-language", "zh-Hans")
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn get_user_app_model_config(token: String, app_id: String) -> Result<serde_json::Value, String> {
    let url = format!("https://aipornhub.ltd/console/api/apps/{}/user_app_model_config", app_id);
    let client = reqwest::Client::new();
    let resp = client.get(&url)
        .header("Accept", "*/*")
        .header("Authorization", format!("Bearer {}", token))
        .header("Content-Type", "application/json")
        .header("x-language", "zh-Hans")
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn export_user_app_model_config(token: String, app_id: String) -> Result<serde_json::Value, String> {
    let url = format!("https://aipornhub.ltd/console/api/apps/{}/user_app_model_config/export", app_id);
    let client = reqwest::Client::new();
    let resp = client.get(&url)
        .header("Accept", "*/*")
        .header("Authorization", format!("Bearer {}", token))
        .header("Content-Type", "application/json")
        .header("x-language", "zh-Hans")
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn update_user_app_model_config(token: String, app_id: String, payload: serde_json::Value) -> Result<serde_json::Value, String> {
    let url = format!("https://aipornhub.ltd/console/api/apps/{}/user_app_model_config", app_id);
    let client = reqwest::Client::new();
    let resp = client.post(&url)
        .header("Accept", "*/*")
        .header("Authorization", format!("Bearer {}", token))
        .header("Content-Type", "application/json")
        .header("Origin", "https://aipornhub.ltd")
        .header("Referer", "https://aipornhub.ltd/")
        .header("x-language", "zh-Hans")
        .json(&payload)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn upload_creation_image(app: tauri::AppHandle, token: String) -> Result<String, String> {
    use tauri_plugin_dialog::DialogExt;
    let file_path = app.dialog()
        .file()
        .add_filter("Images", &["png", "jpg", "jpeg", "gif", "webp"])
        .blocking_pick_file()
        .ok_or("No file selected")?;
    let path_str = file_path.as_path().map(|p| p.to_path_buf()).ok_or("Invalid file path")?;
    let file_data = fs::read(&path_str).map_err(|e| e.to_string())?;
    let file_name = path_str.file_name().unwrap_or_default().to_string_lossy().to_string();
    upload_image_to_remote(&token, &file_data, &file_name).await
}

#[tauri::command]
pub async fn upload_gallery_image_to_remote(app: tauri::AppHandle, image_id: String, token: String) -> Result<String, String> {
    let local_path: String = database::with_db(|conn| {
        conn.query_row(
            "SELECT local_path FROM wt_image WHERE id = ?",
            [&image_id],
            |row| row.get(0),
        )
    }).map_err(|e| e.to_string())?;
    let app_data = app.path().app_data_dir().map_err(|e| e.to_string())?;
    let full_path = app_data.join("gallery").join(&local_path);
    let file_data = fs::read(&full_path).map_err(|e| format!("读取图片失败: {}", e))?;
    let file_name = full_path.file_name().unwrap_or_default().to_string_lossy().to_string();
    upload_image_to_remote(&token, &file_data, &file_name).await
}

async fn upload_image_to_remote(token: &str, file_data: &[u8], file_name: &str) -> Result<String, String> {
    let boundary = format!("----WebKitFormBoundary{}", Utc::now().timestamp_nanos_opt().unwrap_or(0));
    let mut body = Vec::new();
    body.extend_from_slice(format!("--{}\r\n", boundary).as_bytes());
    body.extend_from_slice(format!("Content-Disposition: form-data; name=\"file\"; filename=\"{}\"\r\n", file_name).as_bytes());
    body.extend_from_slice(b"Content-Type: image/png\r\n\r\n");
    body.extend_from_slice(file_data);
    body.extend_from_slice(format!("\r\n--{}--\r\n", boundary).as_bytes());
    let client = reqwest::Client::new();
    let resp = client.post("https://aipornhub.ltd/go/api/app_manage/upload_cover")
        .header("Content-Type", format!("multipart/form-data; boundary={}", boundary))
        .header("Accept", "*/*")
        .header("Authorization", format!("Bearer {}", token))
        .header("Origin", "https://aipornhub.ltd")
        .header("Referer", "https://aipornhub.ltd/")
        .body(body)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    let result: serde_json::Value = serde_json::from_str(&text).map_err(|e| e.to_string())?;
    if result.get("code").and_then(|v| v.as_i64()) != Some(100000) {
        return Err(format!("上传失败: {}", result.get("msg").and_then(|v| v.as_str()).unwrap_or("未知错误")));
    }
    result.get("data")
        .and_then(|d| d.get("url"))
        .and_then(|u| u.as_str())
        .map(|s| s.to_string())
        .ok_or_else(|| "无法获取上传URL".to_string())
}

#[tauri::command]
pub fn save_creation_app_cache(user_id: String, app_id: String, app_name: String, app_icon: String, model_config: String) -> Result<(), String> {
    let now = Utc::now().format("%Y-%m-%dT%H:%M:%SZ").to_string();
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_creation (id, user_id, app_name, app_icon, model_config, updated_at) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT(id, user_id) DO UPDATE SET app_name = excluded.app_name, app_icon = excluded.app_icon, model_config = excluded.model_config, updated_at = excluded.updated_at",
            rusqlite::params![app_id, user_id, app_name, app_icon, model_config, now],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_creation_app_cache(user_id: String, app_id: String) -> Result<String, String> {
    database::with_db(|conn| {
        conn.query_row(
            "SELECT model_config FROM wt_creation WHERE id = ? AND user_id = ?",
            rusqlite::params![app_id, user_id],
            |row| row.get(0),
        )
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_all_creation_app_cache(user_id: String) -> Result<Vec<serde_json::Value>, String> {
    database::with_db(|conn| {
        let mut stmt = conn.prepare("SELECT id, app_name, app_icon, model_config, updated_at FROM wt_creation WHERE user_id = ? ORDER BY updated_at DESC")?;
        let rows = stmt.query_map([&user_id], |row| {
            Ok(serde_json::json!({
                "id": row.get::<_, String>(0)?,
                "app_name": row.get::<_, String>(1)?,
                "app_icon": row.get::<_, String>(2)?,
                "model_config": row.get::<_, String>(3)?,
                "updated_at": row.get::<_, String>(4)?
            }))
        })?;
        rows.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_creation_app_cache(user_id: String, app_id: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "DELETE FROM wt_creation WHERE id = ? AND user_id = ?",
            rusqlite::params![app_id, user_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ByteDanceImageResponse {
    pub data: Option<Vec<ByteDanceImageData>>,
    pub error: Option<ByteDanceError>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ByteDanceImageData {
    pub url: Option<String>,
    pub b64_json: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ByteDanceError {
    pub code: String,
    pub message: String,
}

#[tauri::command]
pub async fn generate_image_byte_dance(api_key: String, prompt: String) -> Result<ByteDanceImageResponse, String> {
    let client = reqwest::Client::builder()
        .timeout(std::time::Duration::from_secs(120))
        .build()
        .map_err(|e| e.to_string())?;
    let body = serde_json::json!({
        "model": "seedream-4-0-250828",
        "prompt": prompt,
        "sequential_image_generation": "disabled",
        "response_format": "url",
        "size": "2K",
        "stream": false,
        "watermark": true
    });
    let resp = client.post("https://ark.ap-southeast.bytepluses.com/api/v3/images/generations")
        .header("Content-Type", "application/json")
        .header("Authorization", format!("Bearer {}", api_key))
        .json(&body)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| format!("解析响应失败: {}", e))
}

#[tauri::command]
pub async fn send_test_chat(request: ChatRequest) -> Result<ChatResponse, String> {
    llm_api::send_chat(request).await
}
