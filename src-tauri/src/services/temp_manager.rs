use std::fs;
use std::path::PathBuf;
use chrono::Utc;

pub fn get_temp_dir(app_data_dir: &PathBuf) -> Result<PathBuf, String> {
    let temp_dir = app_data_dir.join("temp");
    fs::create_dir_all(&temp_dir).map_err(|e| e.to_string())?;
    Ok(temp_dir)
}

pub fn save_temp_file(temp_dir: &PathBuf, name: &str, content: &str) -> Result<PathBuf, String> {
    let file_path = temp_dir.join(name);
    fs::write(&file_path, content).map_err(|e| e.to_string())?;
    Ok(file_path)
}

pub fn load_temp_file(temp_dir: &PathBuf, name: &str) -> Result<String, String> {
    let file_path = temp_dir.join(name);
    fs::read_to_string(&file_path).map_err(|e| e.to_string())
}

pub fn delete_temp_file(temp_dir: &PathBuf, name: &str) -> Result<(), String> {
    let file_path = temp_dir.join(name);
    if file_path.exists() {
        fs::remove_file(&file_path).map_err(|e| e.to_string())?;
    }
    Ok(())
}

pub fn list_temp_files(temp_dir: &PathBuf) -> Result<Vec<String>, String> {
    let mut files = Vec::new();
    if let Ok(entries) = fs::read_dir(temp_dir) {
        for entry in entries.flatten() {
            if entry.path().is_file() {
                if let Some(name) = entry.file_name().to_str() {
                    files.push(name.to_string());
                }
            }
        }
    }
    Ok(files)
}

pub fn cleanup_old_temp_files(temp_dir: &PathBuf, max_age_hours: i64) -> Result<i32, String> {
    let now = Utc::now();
    let mut cleaned = 0;
    if let Ok(entries) = fs::read_dir(temp_dir) {
        for entry in entries.flatten() {
            let path = entry.path();
            if path.is_file() {
                if let Ok(metadata) = fs::metadata(&path) {
                    if let Ok(modified) = metadata.modified() {
                        let modified_time = chrono::DateTime::<Utc>::from(modified);
                        let age = now.signed_duration_since(modified_time);
                        if age.num_hours() > max_age_hours {
                            if fs::remove_file(&path).is_ok() {
                                cleaned += 1;
                            }
                        }
                    }
                }
            }
        }
    }
    Ok(cleaned)
}
