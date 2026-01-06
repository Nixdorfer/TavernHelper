use crate::database;
use crate::types::{GalleryFolder, GalleryImage};
use chrono::Utc;
use sha2::{Sha256, Digest};
use std::fs;
use std::path::PathBuf;
use tauri::Manager;

fn get_gallery_dir(app: &tauri::AppHandle) -> Result<PathBuf, String> {
    let app_data = app.path().app_data_dir().map_err(|e| e.to_string())?;
    let gallery_dir = app_data.join("gallery");
    fs::create_dir_all(&gallery_dir).map_err(|e| e.to_string())?;
    Ok(gallery_dir)
}

#[tauri::command]
pub fn get_gallery_images() -> Result<Vec<GalleryImage>, String> {
    database::with_db_log("get_gallery_images", |conn| {
        let mut stmt = conn.prepare(
            "SELECT id, hash, local_path, remote_url, file_name, file_size, created_at, folder_path FROM wt_image ORDER BY created_at DESC"
        )?;
        let images = stmt.query_map([], |row| {
            Ok(GalleryImage {
                id: row.get(0)?,
                hash: row.get(1)?,
                local_path: row.get(2)?,
                remote_url: row.get::<_, Option<String>>(3)?.unwrap_or_default(),
                file_name: row.get::<_, Option<String>>(4)?.unwrap_or_default(),
                file_size: row.get::<_, Option<i64>>(5)?.unwrap_or(0),
                created_at: row.get(6)?,
                folder_path: row.get::<_, Option<String>>(7)?.unwrap_or_default(),
            })
        })?;
        images.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_gallery_folders() -> Result<Vec<GalleryFolder>, String> {
    database::with_db_log("get_gallery_folders", |conn| {
        let mut stmt = conn.prepare(
            "SELECT folder_path, COUNT(*) as count FROM wt_image WHERE folder_path != '' GROUP BY folder_path"
        )?;
        let folders = stmt.query_map([], |row| {
            let path: String = row.get(0)?;
            Ok(GalleryFolder {
                name: path.clone(),
                path,
                image_count: row.get(1)?,
            })
        })?;
        folders.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn create_gallery_folder(name: String) -> Result<GalleryFolder, String> {
    Ok(GalleryFolder {
        name: name.clone(),
        path: name,
        image_count: 0,
    })
}

#[tauri::command]
pub fn delete_gallery_folder(name: String) -> Result<(), String> {
    database::with_db_log(&format!("delete_gallery_folder: {}", name), |conn| {
        conn.execute(
            "UPDATE wt_image SET folder_path = '' WHERE folder_path = ?",
            [&name],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_gallery_image(app: tauri::AppHandle, id: String) -> Result<(), String> {
    let local_path: String = database::with_db_log(&format!("get_image_path: {}", id), |conn| {
        conn.query_row(
            "SELECT local_path FROM wt_image WHERE id = ?",
            [&id],
            |row| row.get(0),
        )
    }).map_err(|e| e.to_string())?;
    let gallery_dir = get_gallery_dir(&app)?;
    let file_path = gallery_dir.join(&local_path);
    if file_path.exists() {
        fs::remove_file(file_path).map_err(|e| e.to_string())?;
    }
    database::with_db_log(&format!("delete_gallery_image: {}", id), |conn| {
        conn.execute("DELETE FROM wt_image WHERE id = ?", [&id])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_gallery_images(app: tauri::AppHandle, ids: Vec<String>) -> Result<(), String> {
    for id in ids {
        delete_gallery_image(app.clone(), id)?;
    }
    Ok(())
}

#[tauri::command]
pub fn move_gallery_image_to_folder(image_id: String, folder_path: String) -> Result<(), String> {
    database::with_db_log(&format!("move_gallery_image: {}", image_id), |conn| {
        conn.execute(
            "UPDATE wt_image SET folder_path = ? WHERE id = ?",
            [&folder_path, &image_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn read_gallery_image_as_base64(app: tauri::AppHandle, id: String) -> Result<String, String> {
    let local_path: String = database::with_db_log(&format!("get_image_path: {}", id), |conn| {
        conn.query_row(
            "SELECT local_path FROM wt_image WHERE id = ?",
            [&id],
            |row| row.get(0),
        )
    }).map_err(|e| e.to_string())?;
    let gallery_dir = get_gallery_dir(&app)?;
    let file_path = gallery_dir.join(&local_path);
    let data = fs::read(&file_path).map_err(|e| e.to_string())?;
    Ok(base64::Engine::encode(&base64::engine::general_purpose::STANDARD, &data))
}

#[tauri::command]
pub async fn select_and_add_gallery_image(app: tauri::AppHandle) -> Result<GalleryImage, String> {
    use tauri_plugin_dialog::DialogExt;
    let file_path = app.dialog()
        .file()
        .add_filter("Images", &["png", "jpg", "jpeg", "gif", "webp"])
        .blocking_pick_file()
        .ok_or("No file selected")?;
    let path_str = file_path.as_path().map(|p| p.to_string_lossy().to_string()).ok_or("Invalid file path")?;
    add_gallery_image_from_path(app, path_str, String::new()).await
}

#[tauri::command]
pub async fn select_and_add_to_folder(app: tauri::AppHandle, folder_path: String) -> Result<GalleryImage, String> {
    use tauri_plugin_dialog::DialogExt;
    let file_path = app.dialog()
        .file()
        .add_filter("Images", &["png", "jpg", "jpeg", "gif", "webp"])
        .blocking_pick_file()
        .ok_or("No file selected")?;
    let path_str = file_path.as_path().map(|p| p.to_string_lossy().to_string()).ok_or("Invalid file path")?;
    add_gallery_image_from_path(app, path_str, folder_path).await
}

async fn add_gallery_image_from_path(app: tauri::AppHandle, path: String, folder_path: String) -> Result<GalleryImage, String> {
    let data = fs::read(&path).map_err(|e| e.to_string())?;
    let file_name = PathBuf::from(&path)
        .file_name()
        .map(|s| s.to_string_lossy().to_string())
        .unwrap_or_default();
    add_image_internal(app, data, file_name, folder_path)
}

#[tauri::command]
pub fn add_gallery_image_from_base64(app: tauri::AppHandle, base64_data: String, file_name: String) -> Result<GalleryImage, String> {
    let data = base64::Engine::decode(&base64::engine::general_purpose::STANDARD, &base64_data)
        .map_err(|e| e.to_string())?;
    add_image_internal(app, data, file_name, String::new())
}

#[tauri::command]
pub fn add_gallery_image_from_base64_to_folder(app: tauri::AppHandle, base64_data: String, file_name: String, folder_path: String) -> Result<GalleryImage, String> {
    let data = base64::Engine::decode(&base64::engine::general_purpose::STANDARD, &base64_data)
        .map_err(|e| e.to_string())?;
    add_image_internal(app, data, file_name, folder_path)
}

fn add_image_internal(app: tauri::AppHandle, data: Vec<u8>, file_name: String, folder_path: String) -> Result<GalleryImage, String> {
    let mut hasher = Sha256::new();
    hasher.update(&data);
    let hash = format!("{:x}", hasher.finalize());
    let existing: Option<String> = database::with_db_log("check_image_exists", |conn| {
        Ok(conn.query_row(
            "SELECT id FROM wt_image WHERE hash = ?",
            [&hash],
            |row| row.get(0),
        ).ok())
    }).map_err(|e: rusqlite::Error| e.to_string())?;
    if let Some(id) = existing {
        return Err(format!("Image already exists with id: {}", id));
    }
    let id = uuid::Uuid::new_v4().to_string();
    let ext = PathBuf::from(&file_name)
        .extension()
        .map(|s| s.to_string_lossy().to_string())
        .unwrap_or_else(|| "png".to_string());
    let local_path = format!("{}.{}", id, ext);
    let gallery_dir = get_gallery_dir(&app)?;
    let full_path = gallery_dir.join(&local_path);
    fs::write(&full_path, &data).map_err(|e| e.to_string())?;
    let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
    let file_size = data.len() as i64;
    database::with_db_log(&format!("add_gallery_image: {}", id), |conn| {
        conn.execute(
            "INSERT INTO wt_image (id, hash, local_path, file_name, file_size, created_at, folder_path) VALUES (?, ?, ?, ?, ?, ?, ?)",
            rusqlite::params![id, hash, local_path, file_name, file_size, now, folder_path],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())?;
    Ok(GalleryImage {
        id,
        hash,
        local_path,
        remote_url: String::new(),
        file_name,
        file_size,
        created_at: now,
        folder_path,
    })
}

#[tauri::command]
pub fn rename_gallery_image(id: String, new_name: String) -> Result<(), String> {
    database::with_db_log(&format!("rename_gallery_image: {}", id), |conn| {
        conn.execute(
            "UPDATE wt_image SET file_name = ? WHERE id = ?",
            [&new_name, &id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn rename_gallery_folder(old_name: String, new_name: String) -> Result<(), String> {
    database::with_db_log(&format!("rename_gallery_folder: {}", old_name), |conn| {
        conn.execute(
            "UPDATE wt_image SET folder_path = ? WHERE folder_path = ?",
            [&new_name, &old_name],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_gallery_image_url(id: String, url: String) -> Result<(), String> {
    database::with_db_log(&format!("update_gallery_image_url: {}", id), |conn| {
        conn.execute(
            "UPDATE wt_image SET remote_url = ? WHERE id = ?",
            [&url, &id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}
