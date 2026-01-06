use crate::database;
use crate::types::{WTApp, WTConversation, WTDialogue, WTDialogueImage, DialoguesResult};
use chrono::Utc;

#[tauri::command]
pub fn get_apps() -> Result<Vec<WTApp>, String> {
    database::with_db_log("get_apps", |conn| {
        let mut stmt = conn.prepare("SELECT id, name FROM wt_app ORDER BY id")?;
        let apps = stmt.query_map([], |row| {
            Ok(WTApp {
                id: row.get(0)?,
                name: row.get(1)?,
            })
        })?;
        apps.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn create_app(name: String) -> Result<i64, String> {
    database::with_db_log(&format!("create_app: {}", name), |conn| {
        conn.execute("INSERT INTO wt_app (name) VALUES (?)", [&name])?;
        Ok(conn.last_insert_rowid())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_app(app_id: i64) -> Result<(), String> {
    database::with_db_log(&format!("delete_app: {}", app_id), |conn| {
        conn.execute("DELETE FROM wt_app WHERE id = ?", [app_id])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn rename_app(app_id: i64, new_name: String) -> Result<(), String> {
    database::with_db_log(&format!("rename_app: {}", app_id), |conn| {
        conn.execute("UPDATE wt_app SET name = ? WHERE id = ?", rusqlite::params![new_name, app_id])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_local_conversations(app_id: i64) -> Result<Vec<WTConversation>, String> {
    database::with_db_log(&format!("get_local_conversations: {}", app_id), |conn| {
        let mut stmt = conn.prepare("SELECT id, app_id, name, current_node FROM wt_conversation WHERE app_id = ? ORDER BY id DESC")?;
        let conversations = stmt.query_map([app_id], |row| {
            Ok(WTConversation {
                id: row.get(0)?,
                app_id: row.get(1)?,
                name: row.get(2)?,
                current_node: row.get(3)?,
            })
        })?;
        conversations.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn create_local_conversation(app_id: i64, name: String) -> Result<i64, String> {
    database::with_db_log(&format!("create_local_conversation: app={}", app_id), |conn| {
        conn.execute("INSERT INTO wt_conversation (app_id, name) VALUES (?, ?)", rusqlite::params![app_id, name])?;
        Ok(conn.last_insert_rowid())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_local_conversation(conversation_id: i64) -> Result<(), String> {
    database::with_db_log(&format!("delete_local_conversation: {}", conversation_id), |conn| {
        conn.execute("DELETE FROM wt_conversation WHERE id = ?", [conversation_id])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn rename_local_conversation(conversation_id: i64, new_name: String) -> Result<(), String> {
    database::with_db_log(&format!("rename_local_conversation: {}", conversation_id), |conn| {
        conn.execute("UPDATE wt_conversation SET name = ? WHERE id = ?", rusqlite::params![new_name, conversation_id])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn set_conversation_node(conversation_id: i64, node_id: Option<i64>) -> Result<(), String> {
    database::with_db_log(&format!("set_conversation_node: {}", conversation_id), |conn| {
        match node_id {
            Some(nid) => conn.execute("UPDATE wt_conversation SET current_node = ? WHERE id = ?", rusqlite::params![nid, conversation_id])?,
            None => conn.execute("UPDATE wt_conversation SET current_node = NULL WHERE id = ?", [conversation_id])?,
        };
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_conversation_node(conversation_id: i64) -> Result<Option<i64>, String> {
    database::with_db_log(&format!("get_conversation_node: {}", conversation_id), |conn| {
        conn.query_row(
            "SELECT current_node FROM wt_conversation WHERE id = ?",
            [conversation_id],
            |row| row.get(0),
        )
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_dialogues(conversation_id: i64, page: i64, limit: i64) -> Result<DialoguesResult, String> {
    database::with_db_log(&format!("get_dialogues: conv={}", conversation_id), |conn| {
        let total: i64 = conn.query_row(
            "SELECT COUNT(*) FROM wt_dialogue WHERE conversation_id = ?",
            [conversation_id],
            |row| row.get(0),
        )?;
        let offset = (page - 1) * limit;
        let mut stmt = conn.prepare(
            "SELECT id, conversation_id, create_time, request_content, response_content, COALESCE(request_system_prompt, ''), COALESCE(response_system_prompt, ''), node_id, request_point, response_point, request_token, response_token FROM wt_dialogue WHERE conversation_id = ? ORDER BY create_time DESC LIMIT ? OFFSET ?"
        )?;
        let dialogues = stmt.query_map(rusqlite::params![conversation_id, limit, offset], |row| {
            Ok(WTDialogue {
                id: row.get(0)?,
                conversation_id: row.get(1)?,
                create_time: row.get(2)?,
                request_content: row.get(3)?,
                response_content: row.get(4)?,
                request_system_prompt: row.get(5)?,
                response_system_prompt: row.get(6)?,
                node_id: row.get(7)?,
                request_point: row.get(8)?,
                response_point: row.get(9)?,
                request_token: row.get(10)?,
                response_token: row.get(11)?,
            })
        })?;
        Ok(DialoguesResult {
            dialogues: dialogues.collect::<Result<Vec<_>, _>>()?,
            total,
        })
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn create_dialogue(conversation_id: i64, request_content: String, response_content: String) -> Result<i64, String> {
    let now = Utc::now().format("%Y-%m-%dT%H:%M:%SZ").to_string();
    database::with_db_log(&format!("create_dialogue: conv={}", conversation_id), |conn| {
        conn.execute(
            "INSERT INTO wt_dialogue (conversation_id, create_time, request_content, response_content) VALUES (?, ?, ?, ?)",
            rusqlite::params![conversation_id, now, request_content, response_content],
        )?;
        Ok(conn.last_insert_rowid())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_dialogue(dialogue_id: i64, request_content: String, response_content: String) -> Result<(), String> {
    database::with_db_log(&format!("update_dialogue: {}", dialogue_id), |conn| {
        conn.execute(
            "UPDATE wt_dialogue SET request_content = ?, response_content = ? WHERE id = ?",
            rusqlite::params![request_content, response_content, dialogue_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_dialogue(dialogue_id: i64) -> Result<(), String> {
    database::with_db_log(&format!("delete_dialogue: {}", dialogue_id), |conn| {
        conn.execute("DELETE FROM wt_dialogue WHERE id = ?", [dialogue_id])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn add_dialogue_image(dialogue_id: i64, image_url: String, image_path: String, prompt: String) -> Result<i64, String> {
    database::with_db_log(&format!("add_dialogue_image: {}", dialogue_id), |conn| {
        conn.execute(
            "INSERT INTO wt_dialogue_image (dialogue_id, image_url, image_path, prompt) VALUES (?, ?, ?, ?)",
            rusqlite::params![dialogue_id, image_url, image_path, prompt],
        )?;
        Ok(conn.last_insert_rowid())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_dialogue_images(dialogue_id: i64) -> Result<Vec<WTDialogueImage>, String> {
    database::with_db_log(&format!("get_dialogue_images: {}", dialogue_id), |conn| {
        let mut stmt = conn.prepare(
            "SELECT id, dialogue_id, COALESCE(image_url, ''), COALESCE(image_path, ''), prompt FROM wt_dialogue_image WHERE dialogue_id = ?"
        )?;
        let images = stmt.query_map([dialogue_id], |row| {
            Ok(WTDialogueImage {
                id: row.get(0)?,
                dialogue_id: row.get(1)?,
                image_url: row.get(2)?,
                image_path: row.get(3)?,
                prompt: row.get(4)?,
            })
        })?;
        images.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}
