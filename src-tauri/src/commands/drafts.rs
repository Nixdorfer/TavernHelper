use crate::database;
use crate::types::{ClipboardCapture, Draft};
use chrono::Utc;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use tauri::Emitter;

static CLIPBOARD_MONITOR_RUNNING: AtomicBool = AtomicBool::new(false);

#[tauri::command]
pub fn get_all_drafts() -> Result<Vec<Draft>, String> {
    database::with_db_log("get_all_drafts", |conn| {
        let mut stmt = conn.prepare(
            "SELECT id, name, content, parent_id, created_at, updated_at FROM wt_draft ORDER BY updated_at DESC"
        )?;
        let drafts = stmt.query_map([], |row| {
            Ok(Draft {
                id: row.get(0)?,
                name: row.get(1)?,
                content: row.get(2)?,
                parent_id: row.get(3)?,
                created_at: row.get(4)?,
                updated_at: row.get(5)?,
            })
        })?;
        drafts.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn create_draft(draft: Draft) -> Result<Draft, String> {
    let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
    database::with_db_log("create_draft", |conn| {
        conn.execute(
            "INSERT INTO wt_draft (name, content, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
            rusqlite::params![draft.name, draft.content, draft.parent_id, now, now],
        )?;
        let id = conn.last_insert_rowid();
        Ok(Draft {
            id,
            name: draft.name,
            content: draft.content,
            parent_id: draft.parent_id,
            created_at: now.clone(),
            updated_at: now,
        })
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_draft(draft: Draft) -> Result<(), String> {
    let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
    database::with_db_log(&format!("update_draft: {}", draft.id), |conn| {
        conn.execute(
            "UPDATE wt_draft SET name = ?, content = ?, parent_id = ?, updated_at = ? WHERE id = ?",
            rusqlite::params![draft.name, draft.content, draft.parent_id, now, draft.id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_draft(id: i64) -> Result<(), String> {
    database::with_db_log(&format!("delete_draft: {}", id), |conn| {
        conn.execute("DELETE FROM wt_draft WHERE id = ?", [id])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_clipboard_captures() -> Result<Vec<ClipboardCapture>, String> {
    database::with_db_log("get_clipboard_captures", |conn| {
        let mut stmt = conn.prepare(
            "SELECT id, content, captured_at FROM wt_clipboard ORDER BY captured_at DESC"
        )?;
        let captures = stmt.query_map([], |row| {
            Ok(ClipboardCapture {
                id: row.get(0)?,
                content: row.get(1)?,
                captured_at: row.get(2)?,
            })
        })?;
        captures.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn move_clipboard_to_draft(capture_id: i64, name: String, parent_id: Option<i64>) -> Result<Draft, String> {
    let content: String = database::with_db_log(&format!("move_clipboard_to_draft: {}", capture_id), |conn| {
        conn.query_row(
            "SELECT content FROM wt_clipboard WHERE id = ?",
            [capture_id],
            |row| row.get(0),
        )
    }).map_err(|e| e.to_string())?;
    let draft = create_draft(Draft {
        id: 0,
        name,
        content,
        parent_id,
        created_at: String::new(),
        updated_at: String::new(),
    })?;
    database::with_db_log(&format!("delete_clipboard: {}", capture_id), |conn| {
        conn.execute("DELETE FROM wt_clipboard WHERE id = ?", [capture_id])?;
        Ok(())
    }).map_err(|e| e.to_string())?;
    Ok(draft)
}

#[tauri::command]
pub fn start_clipboard_monitor(app: tauri::AppHandle) -> Result<(), String> {
    if CLIPBOARD_MONITOR_RUNNING.swap(true, Ordering::SeqCst) {
        return Ok(());
    }
    let running = Arc::new(AtomicBool::new(true));
    let running_clone = running.clone();
    std::thread::spawn(move || {
        let mut clipboard = arboard::Clipboard::new().ok();
        let mut last_content = String::new();
        while running_clone.load(Ordering::SeqCst) && CLIPBOARD_MONITOR_RUNNING.load(Ordering::SeqCst) {
            if let Some(ref mut cb) = clipboard {
                if let Ok(content) = cb.get_text() {
                    if !content.is_empty() && content != last_content {
                        last_content = content.clone();
                        let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
                        let _ = database::with_db_log("clipboard_capture", |conn| {
                            conn.execute(
                                "INSERT INTO wt_clipboard (content, captured_at) VALUES (?, ?)",
                                [&content, &now],
                            )
                        });
                        let _ = app.emit("clipboard-captured", &content);
                    }
                }
            }
            std::thread::sleep(std::time::Duration::from_millis(500));
        }
    });
    Ok(())
}

#[tauri::command]
pub fn stop_clipboard_monitor() -> Result<(), String> {
    CLIPBOARD_MONITOR_RUNNING.store(false, Ordering::SeqCst);
    Ok(())
}

#[tauri::command]
pub fn clear_all_clipboard_captures() -> Result<(), String> {
    database::with_db_log("clear_all_clipboard_captures", |conn| {
        conn.execute("DELETE FROM wt_clipboard", [])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn copy_to_clipboard(content: String) -> Result<(), String> {
    let mut clipboard = arboard::Clipboard::new().map_err(|e| e.to_string())?;
    clipboard.set_text(content).map_err(|e| e.to_string())
}
