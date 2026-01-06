use crate::database;
use crate::types::{Project, ProjectInfo};
use chrono::Utc;

#[tauri::command]
pub fn get_projects() -> Result<Vec<ProjectInfo>, String> {
    database::with_db(|conn| {
        let mut stmt = conn.prepare(
            "SELECT id, name, type, update_time FROM wt_project ORDER BY update_time DESC"
        )?;
        let projects = stmt.query_map([], |row| {
            Ok(ProjectInfo {
                id: row.get(0)?,
                name: row.get::<_, String>(1)?,
                file_name: row.get::<_, String>(1)?,
                project_type: row.get(2)?,
                updated_at: row.get(3)?,
            })
        })?;
        projects.collect::<Result<Vec<_>, _>>()
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn load_project_by_name(name: String) -> Result<Project, String> {
    database::with_db(|conn| {
        let mut stmt = conn.prepare(
            "SELECT id, name, type, current_node, create_time, update_time FROM wt_project WHERE name = ?"
        )?;
        let project = stmt.query_row([&name], |row| {
            Ok(Project {
                id: row.get(0)?,
                name: row.get::<_, String>(1)?,
                file_name: row.get::<_, String>(1)?,
                project_type: row.get(2)?,
                current_node: row.get(3)?,
                create_time: row.get(4)?,
                update_time: row.get(5)?,
                timeline: vec![],
            })
        })?;
        Ok(project)
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn save_project(file_name: String, _data: serde_json::Value) -> Result<(), String> {
    let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_project SET update_time = ? WHERE name = ?",
            [&now, &file_name],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_project(file_name: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute("DELETE FROM wt_project WHERE name = ?", [&file_name])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn rename_project(old_name: String, new_name: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_project SET name = ? WHERE name = ?",
            [&new_name, &old_name],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn create_project(name: String, project_type: String) -> Result<Project, String> {
    let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_project (name, type, create_time, update_time) VALUES (?, ?, ?, ?)",
            [&name, &project_type, &now, &now],
        )?;
        let id = conn.last_insert_rowid();
        conn.execute(
            "INSERT INTO wt_node (project_id, name) VALUES (?, ?)",
            rusqlite::params![id, "root"],
        )?;
        let node_id = conn.last_insert_rowid();
        conn.execute(
            "UPDATE wt_project SET current_node = ? WHERE id = ?",
            rusqlite::params![node_id, id],
        )?;
        Ok(Project {
            id,
            name: name.clone(),
            file_name: name,
            project_type,
            current_node: Some(node_id),
            create_time: now.clone(),
            update_time: now,
            timeline: vec![],
        })
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_project_current_node(project_name: String, node_id: i64) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_project SET current_node = ? WHERE name = ?",
            rusqlite::params![node_id, project_name],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn load_project(project_id: i64) -> Result<Project, String> {
    database::with_db(|conn| {
        let mut stmt = conn.prepare(
            "SELECT id, name, type, current_node, create_time, update_time FROM wt_project WHERE id = ?"
        )?;
        let project = stmt.query_row([project_id], |row| {
            Ok(Project {
                id: row.get(0)?,
                name: row.get::<_, String>(1)?,
                file_name: row.get::<_, String>(1)?,
                project_type: row.get(2)?,
                current_node: row.get(3)?,
                create_time: row.get(4)?,
                update_time: row.get(5)?,
                timeline: vec![],
            })
        })?;
        Ok(project)
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_project_time(project_id: i64) -> Result<(), String> {
    let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_project SET update_time = ? WHERE id = ?",
            rusqlite::params![now, project_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}
