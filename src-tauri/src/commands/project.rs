use crate::database;
use crate::types::{NodeInfo, Project, ProjectInfo};
use chrono::Utc;

#[tauri::command]
pub fn get_projects() -> Result<Vec<ProjectInfo>, String> {
    database::with_db_log("get_projects", |conn| {
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

fn load_nodes_for_project(conn: &rusqlite::Connection, project_id: i64) -> rusqlite::Result<Vec<NodeInfo>> {
    let mut stmt = conn.prepare(
        "SELECT id, parent_id, name, COALESCE(desc, '') FROM wt_node WHERE project_id = ? ORDER BY id"
    )?;
    let nodes = stmt.query_map([project_id], |row| {
        Ok(NodeInfo {
            id: row.get(0)?,
            parent_id: row.get(1)?,
            name: row.get(2)?,
            note: row.get(3)?,
            tags: vec![],
            created_at: String::new(),
            pre_text: vec![],
            post_text: vec![],
            pre_prompt: vec![],
            world_book: vec![],
        })
    })?;
    nodes.collect()
}

#[tauri::command]
pub fn load_project_by_name(name: String) -> Result<Project, String> {
    database::with_db_log(&format!("load_project_by_name: {}", name), |conn| {
        let mut stmt = conn.prepare(
            "SELECT id, name, type, current_node, create_time, update_time FROM wt_project WHERE name = ?"
        )?;
        let mut project = stmt.query_row([&name], |row| {
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
        project.timeline = load_nodes_for_project(conn, project.id)?;
        if project.current_node.is_none() && !project.timeline.is_empty() {
            project.current_node = Some(project.timeline[0].id);
        }
        Ok(project)
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn save_project(file_name: String, _data: serde_json::Value) -> Result<(), String> {
    let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
    database::with_db_log(&format!("save_project: {}", file_name), |conn| {
        conn.execute(
            "UPDATE wt_project SET update_time = ? WHERE name = ?",
            [&now, &file_name],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_project(file_name: String) -> Result<(), String> {
    database::with_db_log(&format!("delete_project: {}", file_name), |conn| {
        conn.execute("DELETE FROM wt_project WHERE name = ?", [&file_name])?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn rename_project(old_name: String, new_name: String) -> Result<(), String> {
    database::with_db_log(&format!("rename_project: {} -> {}", old_name, new_name), |conn| {
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
    database::with_db_log(&format!("create_project: {}", name), |conn| {
        conn.execute("BEGIN", [])?;
        let result = (|| -> rusqlite::Result<(i64, i64)> {
            conn.execute(
                "INSERT INTO wt_project (name, type, create_time, update_time) VALUES (?, ?, ?, ?)",
                [&name, &project_type, &now, &now],
            )?;
            let id = conn.last_insert_rowid();
            conn.execute(
                "INSERT INTO wt_node (project_id, name) VALUES (?, ?)",
                rusqlite::params![id, "根节点"],
            )?;
            let node_id = conn.last_insert_rowid();
            let sys_folders = [
                ("SYS_PRE", "SYS_PRE_BASE"),
                ("SYS_POST", "SYS_POST_BASE"),
                ("SYS_GLOBAL", "SYS_GLOBAL_BASE"),
            ];
            for (folder_name, card_name) in sys_folders {
                conn.execute("INSERT INTO wt_folder (name) VALUES (?)", [folder_name])?;
                let folder_id = conn.last_insert_rowid();
                conn.execute(
                    "INSERT INTO wt_node_change (action, level, node_id, target, detail_folder) VALUES ('add', 'folder', ?, NULL, ?)",
                    rusqlite::params![node_id, folder_id],
                )?;
                let folder_change_id = conn.last_insert_rowid();
                conn.execute(
                    "INSERT INTO wt_card (name, desc, key_word, trigger_system, trigger_user, trigger_ai) VALUES (?, '', '', 1, 1, 1)",
                    [card_name],
                )?;
                let card_id = conn.last_insert_rowid();
                conn.execute(
                    "INSERT INTO wt_node_change (action, level, node_id, target, detail_card) VALUES ('add', 'card', ?, ?, ?)",
                    rusqlite::params![node_id, folder_change_id, card_id],
                )?;
            }
            conn.execute(
                "UPDATE wt_project SET current_node = ? WHERE id = ?",
                rusqlite::params![node_id, id],
            )?;
            Ok((id, node_id))
        })();
        match result {
            Ok((id, node_id)) => {
                conn.execute("COMMIT", [])?;
                let root_node = NodeInfo {
                    id: node_id,
                    parent_id: None,
                    name: "根节点".to_string(),
                    note: String::new(),
                    tags: vec![],
                    created_at: String::new(),
                    pre_text: vec![],
                    post_text: vec![],
                    pre_prompt: vec![],
                    world_book: vec![],
                };
                Ok(Project {
                    id,
                    name: name.clone(),
                    file_name: name,
                    project_type,
                    current_node: Some(node_id),
                    create_time: now.clone(),
                    update_time: now,
                    timeline: vec![root_node],
                })
            }
            Err(e) => {
                conn.execute("ROLLBACK", []).ok();
                Err(e)
            }
        }
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_project_current_node(project_name: String, node_id: i64) -> Result<(), String> {
    database::with_db_log(&format!("update_project_current_node: {} node={}", project_name, node_id), |conn| {
        conn.execute(
            "UPDATE wt_project SET current_node = ? WHERE name = ?",
            rusqlite::params![node_id, project_name],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn load_project(project_id: i64) -> Result<Project, String> {
    database::with_db_log(&format!("load_project: {}", project_id), |conn| {
        let mut stmt = conn.prepare(
            "SELECT id, name, type, current_node, create_time, update_time FROM wt_project WHERE id = ?"
        )?;
        let mut project = stmt.query_row([project_id], |row| {
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
        project.timeline = load_nodes_for_project(conn, project.id)?;
        if project.current_node.is_none() && !project.timeline.is_empty() {
            project.current_node = Some(project.timeline[0].id);
        }
        Ok(project)
    }).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_project_time(project_id: i64) -> Result<(), String> {
    let now = Utc::now().format("%Y-%m-%d %H:%M:%S").to_string();
    database::with_db_log(&format!("update_project_time: {}", project_id), |conn| {
        conn.execute(
            "UPDATE wt_project SET update_time = ? WHERE id = ?",
            rusqlite::params![now, project_id],
        )?;
        Ok(())
    }).map_err(|e| e.to_string())
}
