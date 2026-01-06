use crate::database;
use crate::types::WTNode;
const SERIAL_CHARS: &[u8] = b"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

#[tauri::command]
pub fn create_child_node(project_name: String, parent_id: i64, name: String) -> Result<WTNode, String> {
    database::with_db(|conn| {
        let project_id: i64 = conn.query_row(
            "SELECT id FROM wt_project WHERE name = ?",
            [&project_name],
            |row| row.get(0),
        )?;
        conn.execute(
            "INSERT INTO wt_node (project_id, parent_id, name) VALUES (?, ?, ?)",
            rusqlite::params![project_id, parent_id, name],
        )?;
        let id = conn.last_insert_rowid();
        Ok(WTNode {
            id,
            project_id,
            parent_id: Some(parent_id),
            name,
            desc: String::new(),
        })
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn create_brother_node(project_name: String, sibling_id: i64, name: String) -> Result<WTNode, String> {
    database::with_db(|conn| {
        let (project_id, parent_id): (i64, Option<i64>) = conn.query_row(
            "SELECT project_id, parent_id FROM wt_node WHERE id = ?",
            [sibling_id],
            |row| Ok((row.get(0)?, row.get(1)?)),
        )?;
        let _project_check: i64 = conn.query_row(
            "SELECT id FROM wt_project WHERE name = ? AND id = ?",
            rusqlite::params![project_name, project_id],
            |row| row.get(0),
        )?;
        conn.execute(
            "INSERT INTO wt_node (project_id, parent_id, name) VALUES (?, ?, ?)",
            rusqlite::params![project_id, parent_id, name],
        )?;
        let id = conn.last_insert_rowid();
        Ok(WTNode {
            id,
            project_id,
            parent_id,
            name,
            desc: String::new(),
        })
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn rename_node(node_id: i64, new_name: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_node SET name = ? WHERE id = ?",
            rusqlite::params![new_name, node_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_node_note(_project_name: String, node_id: i64, note: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_node SET desc = ? WHERE id = ?",
            rusqlite::params![note, node_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_node_desc(node_id: i64, desc: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_node SET desc = ? WHERE id = ?",
            rusqlite::params![desc, node_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_node(_project_name: String, node_id: i64) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute("DELETE FROM wt_node WHERE id = ?", [node_id])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_node_path(node_id: i64) -> Result<Vec<i64>, String> {
    database::with_db(|conn| {
        let mut path = Vec::new();
        let mut current = node_id;
        loop {
            path.insert(0, current);
            let parent_id: Option<i64> = conn.query_row(
                "SELECT parent_id FROM wt_node WHERE id = ?",
                [current],
                |row| row.get(0),
            ).ok().flatten();
            match parent_id {
                Some(pid) => current = pid,
                None => break,
            }
        }
        Ok(path)
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn rebase_node(node_id: i64, new_parent_id: Option<i64>) -> Result<(), String> {
    database::with_db(|conn| {
        if let Some(new_pid) = new_parent_id {
            let mut current = new_pid;
            loop {
                if current == node_id {
                    return Err(rusqlite::Error::InvalidQuery);
                }
                let parent_id: Option<i64> = conn.query_row(
                    "SELECT parent_id FROM wt_node WHERE id = ?",
                    [current],
                    |row| row.get(0),
                ).ok().flatten();
                match parent_id {
                    Some(pid) => current = pid,
                    None => break,
                }
            }
            conn.execute(
                "UPDATE wt_node SET parent_id = ? WHERE id = ?",
                rusqlite::params![new_pid, node_id],
            )?;
        } else {
            conn.execute(
                "UPDATE wt_node SET parent_id = NULL WHERE id = ?",
                [node_id],
            )?;
        }
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_branch_tag(parent_id: i64, child_id: i64) -> Result<Option<crate::types::WTBranchTag>, String> {
    database::with_db(|conn| {
        let result = conn.query_row(
            "SELECT id, parent_id, child_id, name, COALESCE(desc, '') FROM wt_branch_tag WHERE parent_id = ? AND child_id = ?",
            rusqlite::params![parent_id, child_id],
            |row| Ok(crate::types::WTBranchTag {
                id: row.get(0)?,
                parent_id: row.get(1)?,
                child_id: row.get(2)?,
                name: row.get(3)?,
                desc: row.get(4)?,
            }),
        );
        match result {
            Ok(tag) => Ok(Some(tag)),
            Err(rusqlite::Error::QueryReturnedNoRows) => Ok(None),
            Err(e) => Err(e),
        }
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn set_branch_tag(parent_id: i64, child_id: i64, name: String, desc: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_branch_tag (parent_id, child_id, name, desc) VALUES (?, ?, ?, ?) ON CONFLICT(parent_id, child_id) DO UPDATE SET name = excluded.name, desc = excluded.desc",
            rusqlite::params![parent_id, child_id, name, desc],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_branch_tag(parent_id: i64, child_id: i64) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "DELETE FROM wt_branch_tag WHERE parent_id = ? AND child_id = ?",
            rusqlite::params![parent_id, child_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_node_by_id(node_id: i64, name: String, desc: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_node SET name = ?, desc = ? WHERE id = ?",
            rusqlite::params![name, desc, node_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_node(_project_name: String, node_data: serde_json::Value) -> Result<(), String> {
    let node_id = node_data.get("id")
        .and_then(|v| v.as_i64())
        .ok_or("节点数据缺少id")?;
    let name = node_data.get("name")
        .and_then(|v| v.as_str())
        .unwrap_or("");
    let desc = node_data.get("note")
        .and_then(|v| v.as_str())
        .unwrap_or("");
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_node SET name = ?, desc = ? WHERE id = ?",
            rusqlite::params![name, desc, node_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

fn next_serial_str(sn: &str) -> String {
    let base = SERIAL_CHARS.len();
    let mut digits: Vec<u8> = sn.as_bytes().to_vec();
    for i in (0..digits.len()).rev() {
        if let Some(idx) = SERIAL_CHARS.iter().position(|&c| c == digits[i]) {
            if idx < base - 1 {
                digits[i] = SERIAL_CHARS[idx + 1];
                return String::from_utf8(digits).unwrap_or_else(|_| "aaaa".to_string());
            }
            digits[i] = SERIAL_CHARS[0];
        }
    }
    if digits.len() >= 8 {
        return String::from_utf8(digits).unwrap_or_else(|_| "aaaa".to_string());
    }
    "a".repeat(digits.len() + 1)
}

#[tauri::command]
pub fn generate_line_serial(project_name: String) -> Result<String, String> {
    database::with_db(|conn| {
        let project_id: Option<i64> = conn.query_row(
            "SELECT id FROM wt_project WHERE name = ?",
            [&project_name],
            |row| row.get(0),
        ).ok();
        let pid = project_id.unwrap_or(0);
        let max_sn: Option<String> = conn.query_row(
            "SELECT sn FROM wt_line WHERE project_id = ? ORDER BY LENGTH(sn) DESC, sn DESC LIMIT 1",
            [pid],
            |row| row.get(0),
        ).ok();
        Ok(match max_sn {
            Some(sn) if !sn.is_empty() => next_serial_str(&sn),
            _ => "aaaa".to_string(),
        })
    }).map_err(|e| e.to_string())
}
