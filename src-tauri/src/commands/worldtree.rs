use crate::database;
use crate::types::{
    ImmediateChange, NodeContent, NodeDetail, ReplayedBlock, ReplayedCard,
    ReplayedFolder, ReplayedLine, SaveChanges,
};
use std::collections::HashMap;

struct FolderState {
    id: i64,
    name: String,
    change_id: i64,
    deleted: bool,
}

struct CardState {
    id: i64,
    name: String,
    desc: String,
    key_word: String,
    trigger_system: bool,
    trigger_user: bool,
    trigger_ai: bool,
    change_id: i64,
    deleted: bool,
}

struct BlockState {
    id: i64,
    title: String,
    zone: String,
    change_id: i64,
    deleted: bool,
}

struct LineState {
    id: i64,
    sn: String,
    content: String,
    change_id: i64,
    position: Option<i64>,
    node_id: i64,
    node_index: usize,
    deleted: bool,
}

fn get_node_path_internal(conn: &rusqlite::Connection, node_id: i64) -> Vec<i64> {
    let mut path = Vec::new();
    let mut current = node_id;
    loop {
        path.insert(0, current);
        let parent_id: Option<i64> = conn
            .query_row(
                "SELECT parent_id FROM wt_node WHERE id = ?",
                [current],
                |row| row.get(0),
            )
            .ok()
            .flatten();
        match parent_id {
            Some(pid) => current = pid,
            None => break,
        }
    }
    path
}

fn replay_folders(conn: &rusqlite::Connection, path: &[i64]) -> Vec<ReplayedFolder> {
    let mut folder_map: HashMap<i64, FolderState> = HashMap::new();
    let mut del_records: Vec<i64> = Vec::new();
    for &node_id in path {
        let mut stmt = match conn.prepare(
            "SELECT nc.id, nc.action, f.id, f.name
             FROM wt_node_change nc
             LEFT JOIN wt_folder f ON nc.detail_folder = f.id
             WHERE nc.node_id = ? AND nc.level = 'folder'
             ORDER BY nc.id",
        ) {
            Ok(s) => s,
            Err(_) => continue,
        };
        let rows = stmt.query_map([node_id], |row| {
            Ok((
                row.get::<_, i64>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, Option<i64>>(2)?,
                row.get::<_, Option<String>>(3)?,
            ))
        });
        if let Ok(rows) = rows {
            for row in rows.flatten() {
                let (change_id, action, folder_id, folder_name) = row;
                if action == "add" {
                    if let (Some(fid), Some(fname)) = (folder_id, folder_name) {
                        folder_map.insert(
                            fid,
                            FolderState {
                                id: fid,
                                name: fname,
                                change_id,
                                deleted: false,
                            },
                        );
                    }
                } else if action == "del" {
                    del_records.push(change_id);
                }
            }
        }
    }
    for del_change_id in del_records {
        let target_change_id: Option<i64> = conn
            .query_row(
                "SELECT target FROM wt_node_change WHERE id = ?",
                [del_change_id],
                |row| row.get(0),
            )
            .ok()
            .flatten();
        if let Some(target_cid) = target_change_id {
            let target_folder_id: Option<i64> = conn
                .query_row(
                    "SELECT detail_folder FROM wt_node_change WHERE id = ?",
                    [target_cid],
                    |row| row.get(0),
                )
                .ok()
                .flatten();
            if let Some(fid) = target_folder_id {
                if let Some(state) = folder_map.get_mut(&fid) {
                    state.deleted = true;
                }
            }
        }
    }
    let mut result = Vec::new();
    for state in folder_map.values() {
        if !state.deleted {
            let cards = replay_cards_in_folder(conn, path, state.change_id);
            result.push(ReplayedFolder {
                id: state.id,
                name: state.name.clone(),
                change_id: state.change_id,
                cards,
            });
        }
    }
    result
}

fn replay_cards_in_folder(
    conn: &rusqlite::Connection,
    path: &[i64],
    folder_change_id: i64,
) -> Vec<ReplayedCard> {
    let mut card_map: HashMap<i64, CardState> = HashMap::new();
    let mut del_records: Vec<i64> = Vec::new();
    for &node_id in path {
        let mut stmt = match conn.prepare(
            "SELECT nc.id, nc.action, nc.target, c.id, c.name, c.desc, c.key_word, c.trigger_system, c.trigger_user, c.trigger_ai
             FROM wt_node_change nc
             LEFT JOIN wt_card c ON nc.detail_card = c.id
             WHERE nc.node_id = ? AND nc.level = 'card'
             ORDER BY nc.id",
        ) {
            Ok(s) => s,
            Err(_) => continue,
        };
        let rows = stmt.query_map([node_id], |row| {
            Ok((
                row.get::<_, i64>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, Option<i64>>(2)?,
                row.get::<_, Option<i64>>(3)?,
                row.get::<_, Option<String>>(4)?,
                row.get::<_, Option<String>>(5)?,
                row.get::<_, Option<String>>(6)?,
                row.get::<_, Option<i32>>(7)?,
                row.get::<_, Option<i32>>(8)?,
                row.get::<_, Option<i32>>(9)?,
            ))
        });
        if let Ok(rows) = rows {
            for row in rows.flatten() {
                let (change_id, action, target, card_id, name, desc, key_word, ts, tu, ta) = row;
                if action == "add" {
                    if let Some(cid) = card_id {
                        let matches = match target {
                            Some(t) => t == folder_change_id,
                            None => folder_change_id == 0,
                        };
                        if matches {
                            card_map.insert(
                                cid,
                                CardState {
                                    id: cid,
                                    name: name.unwrap_or_default(),
                                    desc: desc.unwrap_or_default(),
                                    key_word: key_word.unwrap_or_default(),
                                    trigger_system: ts.unwrap_or(0) == 1,
                                    trigger_user: tu.unwrap_or(0) == 1,
                                    trigger_ai: ta.unwrap_or(0) == 1,
                                    change_id,
                                    deleted: false,
                                },
                            );
                        }
                    }
                } else if action == "del" {
                    del_records.push(change_id);
                }
            }
        }
    }
    for del_change_id in del_records {
        let target_change_id: Option<i64> = conn
            .query_row(
                "SELECT target FROM wt_node_change WHERE id = ?",
                [del_change_id],
                |row| row.get(0),
            )
            .ok()
            .flatten();
        if let Some(target_cid) = target_change_id {
            let target_card_id: Option<i64> = conn
                .query_row(
                    "SELECT detail_card FROM wt_node_change WHERE id = ?",
                    [target_cid],
                    |row| row.get(0),
                )
                .ok()
                .flatten();
            if let Some(cid) = target_card_id {
                if let Some(state) = card_map.get_mut(&cid) {
                    state.deleted = true;
                }
            }
        }
    }
    let mut result = Vec::new();
    for state in card_map.values() {
        if !state.deleted {
            let blocks = replay_blocks_in_card(conn, path, state.change_id);
            result.push(ReplayedCard {
                id: state.id,
                name: state.name.clone(),
                desc: state.desc.clone(),
                key_word: state.key_word.clone(),
                trigger_system: state.trigger_system,
                trigger_user: state.trigger_user,
                trigger_ai: state.trigger_ai,
                change_id: state.change_id,
                blocks,
            });
        }
    }
    result
}

fn replay_blocks_by_zone(
    conn: &rusqlite::Connection,
    path: &[i64],
    zone: &str,
) -> Vec<ReplayedBlock> {
    let mut block_map: HashMap<i64, BlockState> = HashMap::new();
    let mut del_records: Vec<i64> = Vec::new();
    for &node_id in path {
        let mut stmt = match conn.prepare(
            "SELECT nc.id, nc.action, b.id, b.title, b.zone
             FROM wt_node_change nc
             LEFT JOIN wt_block b ON nc.detail_block = b.id
             WHERE nc.node_id = ? AND nc.level = 'block' AND nc.target IS NULL
             ORDER BY nc.id",
        ) {
            Ok(s) => s,
            Err(_) => continue,
        };
        let rows = stmt.query_map([node_id], |row| {
            Ok((
                row.get::<_, i64>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, Option<i64>>(2)?,
                row.get::<_, Option<String>>(3)?,
                row.get::<_, Option<String>>(4)?,
            ))
        });
        if let Ok(rows) = rows {
            for row in rows.flatten() {
                let (change_id, action, block_id, title, block_zone) = row;
                if action == "add" {
                    if let (Some(bid), Some(bzone)) = (block_id, block_zone.clone()) {
                        if bzone == zone {
                            block_map.insert(
                                bid,
                                BlockState {
                                    id: bid,
                                    title: title.unwrap_or_default(),
                                    zone: bzone,
                                    change_id,
                                    deleted: false,
                                },
                            );
                        }
                    }
                } else if action == "del" {
                    del_records.push(change_id);
                }
            }
        }
    }
    for del_change_id in del_records {
        let target_change_id: Option<i64> = conn
            .query_row(
                "SELECT target FROM wt_node_change WHERE id = ?",
                [del_change_id],
                |row| row.get(0),
            )
            .ok()
            .flatten();
        if let Some(target_cid) = target_change_id {
            let target_block_id: Option<i64> = conn
                .query_row(
                    "SELECT detail_block FROM wt_node_change WHERE id = ?",
                    [target_cid],
                    |row| row.get(0),
                )
                .ok()
                .flatten();
            if let Some(bid) = target_block_id {
                if let Some(state) = block_map.get_mut(&bid) {
                    state.deleted = true;
                }
            }
        }
    }
    let mut result = Vec::new();
    for state in block_map.values() {
        if !state.deleted {
            let lines = replay_lines_in_block(conn, path, state.change_id);
            result.push(ReplayedBlock {
                id: state.id,
                title: state.title.clone(),
                zone: state.zone.clone(),
                change_id: state.change_id,
                lines,
            });
        }
    }
    result
}

fn replay_blocks_in_card(
    conn: &rusqlite::Connection,
    path: &[i64],
    card_change_id: i64,
) -> Vec<ReplayedBlock> {
    let mut block_map: HashMap<i64, BlockState> = HashMap::new();
    let mut del_records: Vec<i64> = Vec::new();
    for &node_id in path {
        let mut stmt = match conn.prepare(
            "SELECT nc.id, nc.action, nc.target, b.id, b.title, b.zone
             FROM wt_node_change nc
             LEFT JOIN wt_block b ON nc.detail_block = b.id
             WHERE nc.node_id = ? AND nc.level = 'block'
             ORDER BY nc.id",
        ) {
            Ok(s) => s,
            Err(_) => continue,
        };
        let rows = stmt.query_map([node_id], |row| {
            Ok((
                row.get::<_, i64>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, Option<i64>>(2)?,
                row.get::<_, Option<i64>>(3)?,
                row.get::<_, Option<String>>(4)?,
                row.get::<_, Option<String>>(5)?,
            ))
        });
        if let Ok(rows) = rows {
            for row in rows.flatten() {
                let (change_id, action, target, block_id, title, zone) = row;
                if action == "add" {
                    if let Some(bid) = block_id {
                        if target == Some(card_change_id) {
                            block_map.insert(
                                bid,
                                BlockState {
                                    id: bid,
                                    title: title.unwrap_or_default(),
                                    zone: zone.unwrap_or_default(),
                                    change_id,
                                    deleted: false,
                                },
                            );
                        }
                    }
                } else if action == "del" {
                    del_records.push(change_id);
                }
            }
        }
    }
    for del_change_id in del_records {
        let target_change_id: Option<i64> = conn
            .query_row(
                "SELECT target FROM wt_node_change WHERE id = ?",
                [del_change_id],
                |row| row.get(0),
            )
            .ok()
            .flatten();
        if let Some(target_cid) = target_change_id {
            let target_block_id: Option<i64> = conn
                .query_row(
                    "SELECT detail_block FROM wt_node_change WHERE id = ?",
                    [target_cid],
                    |row| row.get(0),
                )
                .ok()
                .flatten();
            if let Some(bid) = target_block_id {
                if let Some(state) = block_map.get_mut(&bid) {
                    state.deleted = true;
                }
            }
        }
    }
    let mut result = Vec::new();
    for state in block_map.values() {
        if !state.deleted {
            let lines = replay_lines_in_block(conn, path, state.change_id);
            result.push(ReplayedBlock {
                id: state.id,
                title: state.title.clone(),
                zone: state.zone.clone(),
                change_id: state.change_id,
                lines,
            });
        }
    }
    result
}

fn replay_lines_in_block(
    conn: &rusqlite::Connection,
    path: &[i64],
    block_change_id: i64,
) -> Vec<ReplayedLine> {
    if path.is_empty() {
        return Vec::new();
    }
    let current_node_id = *path.last().unwrap();
    let node_index_map: HashMap<i64, usize> = path.iter().enumerate().map(|(i, &n)| (n, i)).collect();
    let mut line_map: HashMap<i64, LineState> = HashMap::new();
    let mut del_records: Vec<i64> = Vec::new();
    for &node_id in path {
        let mut stmt = match conn.prepare(
            "SELECT nc.id, nc.action, nc.target, l.id, l.sn, l.content, l.position
             FROM wt_node_change nc
             LEFT JOIN wt_line l ON nc.detail_line = l.id
             WHERE nc.node_id = ? AND nc.level = 'line'
             ORDER BY nc.id",
        ) {
            Ok(s) => s,
            Err(_) => continue,
        };
        let rows = stmt.query_map([node_id], |row| {
            Ok((
                row.get::<_, i64>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, Option<i64>>(2)?,
                row.get::<_, Option<i64>>(3)?,
                row.get::<_, Option<String>>(4)?,
                row.get::<_, Option<String>>(5)?,
                row.get::<_, Option<i64>>(6)?,
            ))
        });
        if let Ok(rows) = rows {
            for row in rows.flatten() {
                let (change_id, action, target, line_id, sn, content, position) = row;
                if action == "add" {
                    if let Some(lid) = line_id {
                        if target == Some(block_change_id) {
                            line_map.insert(
                                lid,
                                LineState {
                                    id: lid,
                                    sn: sn.unwrap_or_default(),
                                    content: content.unwrap_or_default(),
                                    change_id,
                                    position,
                                    node_id,
                                    node_index: *node_index_map.get(&node_id).unwrap_or(&0),
                                    deleted: false,
                                },
                            );
                        }
                    }
                } else if action == "del" {
                    del_records.push(change_id);
                }
            }
        }
    }
    for del_change_id in del_records {
        let target_change_id: Option<i64> = conn
            .query_row(
                "SELECT target FROM wt_node_change WHERE id = ?",
                [del_change_id],
                |row| row.get(0),
            )
            .ok()
            .flatten();
        if let Some(target_cid) = target_change_id {
            let target_line_id: Option<i64> = conn
                .query_row(
                    "SELECT detail_line FROM wt_node_change WHERE id = ?",
                    [target_cid],
                    |row| row.get(0),
                )
                .ok()
                .flatten();
            if let Some(lid) = target_line_id {
                if let Some(state) = line_map.get_mut(&lid) {
                    state.deleted = true;
                }
            }
        }
    }
    let child_nodes: Vec<i64> = {
        let mut stmt = conn
            .prepare("SELECT id FROM wt_node WHERE parent_id = ?")
            .unwrap();
        stmt.query_map([current_node_id], |row| row.get(0))
            .unwrap()
            .flatten()
            .collect()
    };
    let mut child_delete_map: HashMap<i64, HashMap<i64, bool>> = HashMap::new();
    for &child_id in &child_nodes {
        let mut child_deletes: HashMap<i64, bool> = HashMap::new();
        let mut stmt = match conn.prepare(
            "SELECT nc.target FROM wt_node_change nc
             WHERE nc.node_id = ? AND nc.level = 'line' AND nc.action = 'del'",
        ) {
            Ok(s) => s,
            Err(_) => continue,
        };
        let targets: Vec<i64> = stmt
            .query_map([child_id], |row| row.get(0))
            .unwrap()
            .flatten()
            .collect();
        for target_change_id in targets {
            let target_line_id: Option<i64> = conn
                .query_row(
                    "SELECT detail_line FROM wt_node_change WHERE id = ?",
                    [target_change_id],
                    |row| row.get(0),
                )
                .ok()
                .flatten();
            if let Some(lid) = target_line_id {
                child_deletes.insert(lid, true);
            }
        }
        child_delete_map.insert(child_id, child_deletes);
    }
    let mut unsorted: Vec<(ReplayedLine, &LineState)> = Vec::new();
    for (line_id, state) in &line_map {
        if !state.deleted {
            let shape = if state.node_id == current_node_id {
                "square"
            } else {
                "dot"
            };
            let color = if !child_nodes.is_empty() {
                let delete_count = child_nodes
                    .iter()
                    .filter(|&&child_id| {
                        child_delete_map
                            .get(&child_id)
                            .map_or(false, |m| m.contains_key(line_id))
                    })
                    .count();
                if delete_count == child_nodes.len() {
                    "red"
                } else if delete_count > 0 {
                    "yellow"
                } else {
                    "green"
                }
            } else {
                "green"
            };
            let content = if state.content.trim().is_empty() {
                None
            } else {
                Some(state.content.clone())
            };
            unsorted.push((
                ReplayedLine {
                    id: state.id,
                    sn: state.sn.clone(),
                    content,
                    sync_dot: format!("{}-{}", color, shape),
                    change_id: state.change_id,
                },
                state,
            ));
        }
    }
    sort_lines_by_position(unsorted, &line_map)
}

fn sort_lines_by_position(
    lines: Vec<(ReplayedLine, &LineState)>,
    line_map: &HashMap<i64, LineState>,
) -> Vec<ReplayedLine> {
    if lines.is_empty() {
        return Vec::new();
    }
    let id_to_line: HashMap<i64, ReplayedLine> =
        lines.iter().map(|(l, _)| (l.id, l.clone())).collect();
    let mut after_map: HashMap<i64, Vec<i64>> = HashMap::new();
    let mut first_line_ids: Vec<i64> = Vec::new();
    for (_, state) in line_map {
        if state.deleted {
            continue;
        }
        match state.position {
            None => first_line_ids.push(state.id),
            Some(pos) => after_map.entry(pos).or_default().push(state.id),
        }
    }
    for nexts in after_map.values_mut() {
        nexts.sort_by(|a, b| {
            let ni_a = line_map.get(a).map(|s| s.node_index).unwrap_or(0);
            let ni_b = line_map.get(b).map(|s| s.node_index).unwrap_or(0);
            if ni_a != ni_b {
                ni_a.cmp(&ni_b)
            } else {
                a.cmp(b)
            }
        });
    }
    first_line_ids.sort_by(|a, b| {
        let ni_a = line_map.get(a).map(|s| s.node_index).unwrap_or(0);
        let ni_b = line_map.get(b).map(|s| s.node_index).unwrap_or(0);
        if ni_a != ni_b {
            ni_a.cmp(&ni_b)
        } else {
            a.cmp(b)
        }
    });
    if first_line_ids.is_empty() {
        return lines.into_iter().map(|(l, _)| l).collect();
    }
    let mut result = Vec::new();
    let mut visited: HashMap<i64, bool> = HashMap::new();
    let mut queue: Vec<i64> = first_line_ids;
    while let Some(current) = queue.first().cloned() {
        queue.remove(0);
        if visited.contains_key(&current) {
            continue;
        }
        visited.insert(current, true);
        if let Some(line) = id_to_line.get(&current) {
            result.push(line.clone());
        }
        if let Some(nexts) = after_map.get(&current) {
            for &next in nexts {
                queue.push(next);
            }
        }
    }
    if result.len() < lines.len() {
        for (line, _) in &lines {
            if !visited.contains_key(&line.id) {
                result.push(line.clone());
            }
        }
    }
    result
}

#[tauri::command]
pub fn get_node_content(node_id: i64) -> Result<NodeContent, String> {
    database::with_db(|conn| {
        let path = get_node_path_internal(conn, node_id);
        if path.is_empty() {
            return Err(rusqlite::Error::QueryReturnedNoRows);
        }
        let folders = replay_folders(conn, &path);
        let pre = replay_blocks_by_zone(conn, &path, "pre");
        let post = replay_blocks_by_zone(conn, &path, "post");
        let global = replay_blocks_by_zone(conn, &path, "global");
        Ok(NodeContent {
            node_id,
            folders,
            pre,
            post,
            global,
        })
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn add_folder(node_id: i64, name: String) -> Result<i64, String> {
    if name.starts_with("SYS_") {
        return Err("不允许创建SYS_开头的文件夹".to_string());
    }
    database::with_db(|conn| {
        conn.execute("INSERT INTO wt_folder (name) VALUES (?)", [&name])?;
        let folder_id = conn.last_insert_rowid();
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, detail_folder) VALUES ('add', 'folder', ?, ?)",
            rusqlite::params![node_id, folder_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn add_folder_with_parent(
    node_id: i64,
    parent_change_id: Option<i64>,
    name: String,
) -> Result<i64, String> {
    if name.starts_with("SYS_") {
        return Err("不允许创建SYS_开头的文件夹".to_string());
    }
    database::with_db(|conn| {
        conn.execute("INSERT INTO wt_folder (name) VALUES (?)", [&name])?;
        let folder_id = conn.last_insert_rowid();
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, target, detail_folder) VALUES ('add', 'folder', ?, ?, ?)",
            rusqlite::params![node_id, parent_change_id, folder_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn add_card(
    node_id: i64,
    folder_change_id: Option<i64>,
    name: String,
    key_word: String,
) -> Result<i64, String> {
    if name.starts_with("SYS_") {
        return Err("不允许创建SYS_开头的卡片".to_string());
    }
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_card (name, desc, key_word, trigger_system, trigger_user, trigger_ai) VALUES (?, '', ?, 0, 1, 1)",
            rusqlite::params![name, key_word],
        )?;
        let card_id = conn.last_insert_rowid();
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, target, detail_card) VALUES ('add', 'card', ?, ?, ?)",
            rusqlite::params![node_id, folder_change_id, card_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn add_block(
    node_id: i64,
    card_change_id: Option<i64>,
    title: String,
    zone: String,
) -> Result<i64, String> {
    if title.starts_with("SYS_") {
        return Err("不允许创建SYS_开头的区块".to_string());
    }
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_block (title, zone) VALUES (?, ?)",
            [&title, &zone],
        )?;
        let block_id = conn.last_insert_rowid();
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, target, detail_block) VALUES ('add', 'block', ?, ?, ?)",
            rusqlite::params![node_id, card_change_id, block_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn add_line(
    node_id: i64,
    project_id: i64,
    block_change_id: i64,
    content: String,
) -> Result<i64, String> {
    database::with_db(|conn| {
        let sn = generate_serial(conn, project_id);
        conn.execute(
            "INSERT INTO wt_line (sn, project_id, content, node_id) VALUES (?, ?, ?, ?)",
            rusqlite::params![sn, project_id, content, node_id],
        )?;
        let line_id = conn.last_insert_rowid();
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, target, detail_line) VALUES ('add', 'line', ?, ?, ?)",
            rusqlite::params![node_id, block_change_id, line_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

fn generate_serial(conn: &rusqlite::Connection, project_id: i64) -> String {
    let max_sn: Option<String> = conn
        .query_row(
            "SELECT sn FROM wt_line WHERE project_id = ? ORDER BY LENGTH(sn) DESC, sn DESC LIMIT 1",
            [project_id],
            |row| row.get(0),
        )
        .ok();
    match max_sn {
        Some(sn) if !sn.is_empty() => next_serial(&sn),
        _ => "aaaa".to_string(),
    }
}

fn next_serial(sn: &str) -> String {
    const CHARS: &[u8] = b"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    let base = CHARS.len();
    let mut digits: Vec<u8> = sn.as_bytes().to_vec();
    for i in (0..digits.len()).rev() {
        if let Some(idx) = CHARS.iter().position(|&c| c == digits[i]) {
            if idx < base - 1 {
                digits[i] = CHARS[idx + 1];
                return String::from_utf8(digits).unwrap_or_else(|_| "aaaa".to_string());
            }
            digits[i] = CHARS[0];
        }
    }
    if digits.len() >= 8 {
        return String::from_utf8(digits).unwrap_or_else(|_| "aaaa".to_string());
    }
    "a".repeat(digits.len() + 1)
}

#[tauri::command]
pub fn update_line_content(line_id: i64, content: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_line SET content = ? WHERE id = ?",
            rusqlite::params![content, line_id],
        )?;
        Ok(())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_card_key_word(card_id: i64, key_word: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_card SET key_word = ? WHERE id = ?",
            rusqlite::params![key_word, card_id],
        )?;
        Ok(())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_card_triggers(
    card_id: i64,
    trigger_system: bool,
    trigger_user: bool,
    trigger_ai: bool,
) -> Result<(), String> {
    let ts: i32 = if trigger_system { 1 } else { 0 };
    let tu: i32 = if trigger_user { 1 } else { 0 };
    let ta: i32 = if trigger_ai { 1 } else { 0 };
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_card SET trigger_system = ?, trigger_user = ?, trigger_ai = ? WHERE id = ?",
            rusqlite::params![ts, tu, ta, card_id],
        )?;
        Ok(())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_block_title(block_id: i64, title: String) -> Result<(), String> {
    database::with_db(|conn| {
        conn.execute(
            "UPDATE wt_block SET title = ? WHERE id = ?",
            rusqlite::params![title, block_id],
        )?;
        Ok(())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_folder(node_id: i64, folder_change_id: i64) -> Result<i64, String> {
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, target) VALUES ('del', 'folder', ?, ?)",
            rusqlite::params![node_id, folder_change_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_card(node_id: i64, card_change_id: i64) -> Result<i64, String> {
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, target) VALUES ('del', 'card', ?, ?)",
            rusqlite::params![node_id, card_change_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_block(node_id: i64, block_change_id: i64) -> Result<i64, String> {
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, target) VALUES ('del', 'block', ?, ?)",
            rusqlite::params![node_id, block_change_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn delete_line(node_id: i64, line_change_id: i64) -> Result<i64, String> {
    database::with_db(|conn| {
        conn.execute(
            "INSERT INTO wt_node_change (action, level, node_id, target) VALUES ('del', 'line', ?, ?)",
            rusqlite::params![node_id, line_change_id],
        )?;
        Ok(conn.last_insert_rowid())
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn get_node_detail(node_id: i64) -> Result<NodeDetail, String> {
    database::with_db(|conn| {
        let (name, desc): (String, Option<String>) = conn.query_row(
            "SELECT name, desc FROM wt_node WHERE id = ?",
            [node_id],
            |row| Ok((row.get(0)?, row.get(1)?)),
        )?;
        Ok(NodeDetail {
            node_id,
            name,
            desc: desc.unwrap_or_default(),
            structure: HashMap::new(),
        })
    })
    .map_err(|e| e.to_string())
}

#[tauri::command]
pub fn immediate_change(node_id: i64, change: ImmediateChange) -> Result<i64, String> {
    match change.level.as_str() {
        "folder" => match change.action.as_str() {
            "add" => add_folder(node_id, change.name),
            "del" => {
                let target = change.target.ok_or("删除folder需要target")?;
                delete_folder(node_id, target)
            }
            _ => Err("未知的action".to_string()),
        },
        "card" => match change.action.as_str() {
            "add" => add_card(node_id, change.target, change.name, String::new()),
            "del" => {
                let target = change.target.ok_or("删除card需要target")?;
                delete_card(node_id, target)
            }
            _ => Err("未知的action".to_string()),
        },
        "block" => match change.action.as_str() {
            "add" => {
                let block_info: Result<serde_json::Value, _> = serde_json::from_str(&change.name);
                let (title, zone) = match block_info {
                    Ok(v) => (
                        v.get("title")
                            .and_then(|t| t.as_str())
                            .unwrap_or(&change.name)
                            .to_string(),
                        v.get("zone")
                            .and_then(|z| z.as_str())
                            .unwrap_or("card")
                            .to_string(),
                    ),
                    Err(_) => (change.name.clone(), "card".to_string()),
                };
                add_block(node_id, change.target, title, zone)
            }
            "del" => {
                let target = change.target.ok_or("删除block需要target")?;
                delete_block(node_id, target)
            }
            _ => Err("未知的action".to_string()),
        },
        _ => Err(format!("未知的level: {}", change.level)),
    }
}

fn calc_md5(s: &str) -> String {
    format!("{:x}", md5::compute(s.as_bytes()))
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct LineData {
    pub content: Option<String>,
    pub sync_dot: String,
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct SaveNodeChangesResult {
    pub blocks: HashMap<String, HashMap<String, LineData>>,
}

#[tauri::command]
pub fn save_node_changes(
    node_id: i64,
    project_id: i64,
    changes: SaveChanges,
) -> Result<SaveNodeChangesResult, String> {
    database::with_db(|conn| {
        let path = get_node_path_internal(conn, node_id);
        let mut saved_block_change_ids: Vec<i64> = Vec::new();
        for (change_id_str, block_change) in &changes.block {
            if block_change.content.is_some() {
                if let Ok(cid) = change_id_str.parse::<i64>() {
                    saved_block_change_ids.push(cid);
                }
            }
        }
        for (change_id_str, folder_change) in &changes.folder {
            if let Ok(change_id) = change_id_str.parse::<i64>() {
                let folder_id: Option<i64> = conn
                    .query_row(
                        "SELECT detail_folder FROM wt_node_change WHERE id = ?",
                        [change_id],
                        |row| row.get(0),
                    )
                    .ok()
                    .flatten();
                if let Some(fid) = folder_id {
                    if let Some(name) = &folder_change.name {
                        let _ = conn.execute(
                            "UPDATE wt_folder SET name = ? WHERE id = ?",
                            rusqlite::params![name, fid],
                        );
                    }
                }
            }
        }
        for (change_id_str, card_change) in &changes.card {
            if let Ok(change_id) = change_id_str.parse::<i64>() {
                let card_id: Option<i64> = conn
                    .query_row(
                        "SELECT detail_card FROM wt_node_change WHERE id = ?",
                        [change_id],
                        |row| row.get(0),
                    )
                    .ok()
                    .flatten();
                if let Some(cid) = card_id {
                    if let Some(name) = &card_change.name {
                        let _ = conn.execute(
                            "UPDATE wt_card SET name = ? WHERE id = ?",
                            rusqlite::params![name, cid],
                        );
                    }
                    if let Some(desc) = &card_change.desc {
                        let _ = conn.execute(
                            "UPDATE wt_card SET desc = ? WHERE id = ?",
                            rusqlite::params![desc, cid],
                        );
                    }
                    if let Some(trigger) = &card_change.trigger {
                        let kw = trigger.words.join("@");
                        let trigger_mode = if trigger.mode == "and" { 1 } else { 0 };
                        let ts = if trigger.system { 1 } else { 0 };
                        let tu = if trigger.user { 1 } else { 0 };
                        let ta = if trigger.ai { 1 } else { 0 };
                        let _ = conn.execute(
                            "UPDATE wt_card SET key_word = ?, trigger_mode = ?, trigger_system = ?, trigger_user = ?, trigger_ai = ? WHERE id = ?",
                            rusqlite::params![kw, trigger_mode, ts, tu, ta, cid],
                        );
                    }
                    if !card_change.image.is_empty() {
                        let image_word = card_change.image.join("@");
                        let _ = conn.execute(
                            "UPDATE wt_card SET image_word = ? WHERE id = ?",
                            rusqlite::params![image_word, cid],
                        );
                    }
                }
            }
        }
        for (change_id_str, block_change) in &changes.block {
            if let Ok(change_id) = change_id_str.parse::<i64>() {
                let block_id: Option<i64> = conn
                    .query_row(
                        "SELECT detail_block FROM wt_node_change WHERE id = ?",
                        [change_id],
                        |row| row.get(0),
                    )
                    .ok()
                    .flatten();
                if let Some(bid) = block_id {
                    if let Some(name) = &block_change.name {
                        let _ = conn.execute(
                            "UPDATE wt_block SET title = ? WHERE id = ?",
                            rusqlite::params![name, bid],
                        );
                    }
                }
                if let Some(content) = &block_change.content {
                    diff_and_save_block_content(conn, &path, node_id, project_id, change_id, content);
                }
            }
        }
        let mut result_blocks: HashMap<String, HashMap<String, LineData>> = HashMap::new();
        for block_change_id in saved_block_change_ids {
            let lines = replay_lines_in_block(conn, &path, block_change_id);
            let mut block_lines: HashMap<String, LineData> = HashMap::new();
            for line in lines {
                block_lines.insert(
                    line.sn.clone(),
                    LineData {
                        content: line.content,
                        sync_dot: line.sync_dot,
                    },
                );
            }
            result_blocks.insert(block_change_id.to_string(), block_lines);
        }
        Ok(SaveNodeChangesResult {
            blocks: result_blocks,
        })
    })
    .map_err(|e| e.to_string())
}

struct DiffLine {
    sn: String,
    content: String,
    line_id: i64,
    change_id: i64,
    node_id: i64,
    used: bool,
}

fn diff_and_save_block_content(
    conn: &rusqlite::Connection,
    path: &[i64],
    node_id: i64,
    project_id: i64,
    block_change_id: i64,
    new_content: &str,
) {
    let mut old_lines: Vec<DiffLine> = Vec::new();
    let mut deleted_change_ids: HashMap<i64, bool> = HashMap::new();
    for &path_node_id in path {
        let mut stmt = match conn.prepare(
            "SELECT nc.id, l.id, l.sn, l.content, nc.node_id
             FROM wt_node_change nc
             JOIN wt_line l ON nc.detail_line = l.id
             WHERE nc.node_id = ? AND nc.level = 'line' AND nc.action = 'add' AND nc.target = ?
             ORDER BY nc.id",
        ) {
            Ok(s) => s,
            Err(_) => continue,
        };
        let rows = stmt.query_map(rusqlite::params![path_node_id, block_change_id], |row| {
            Ok((
                row.get::<_, i64>(0)?,
                row.get::<_, i64>(1)?,
                row.get::<_, String>(2)?,
                row.get::<_, String>(3)?,
                row.get::<_, i64>(4)?,
            ))
        });
        if let Ok(rows) = rows {
            for row in rows.flatten() {
                let (chg_id, line_id, sn, content, chg_node_id) = row;
                old_lines.push(DiffLine {
                    sn,
                    content,
                    line_id,
                    change_id: chg_id,
                    node_id: chg_node_id,
                    used: false,
                });
            }
        }
        let mut del_stmt = match conn.prepare(
            "SELECT nc.target FROM wt_node_change nc
             WHERE nc.node_id = ? AND nc.level = 'line' AND nc.action = 'del'",
        ) {
            Ok(s) => s,
            Err(_) => continue,
        };
        let del_targets: Vec<i64> = del_stmt
            .query_map([path_node_id], |row| row.get(0))
            .unwrap()
            .flatten()
            .collect();
        for target_change_id in del_targets {
            deleted_change_ids.insert(target_change_id, true);
        }
    }
    old_lines.retain(|line| !deleted_change_ids.contains_key(&line.change_id));
    let new_lines: Vec<&str> = if new_content.is_empty() {
        Vec::new()
    } else {
        new_content.split('\n').collect()
    };
    let mut old_md5_map: HashMap<String, Vec<usize>> = HashMap::new();
    for (i, line) in old_lines.iter().enumerate() {
        let h = calc_md5(&line.content);
        old_md5_map.entry(h).or_default().push(i);
    }
    let mut last_line_id: Option<i64> = None;
    for new_line in &new_lines {
        let h = calc_md5(new_line);
        let mut found_idx: Option<usize> = None;
        if let Some(candidates) = old_md5_map.get(&h) {
            for &idx in candidates {
                if !old_lines[idx].used {
                    found_idx = Some(idx);
                    break;
                }
            }
        }
        if let Some(idx) = found_idx {
            old_lines[idx].used = true;
            let _ = conn.execute(
                "UPDATE wt_line SET position = ? WHERE id = ?",
                rusqlite::params![last_line_id, old_lines[idx].line_id],
            );
            last_line_id = Some(old_lines[idx].line_id);
        } else {
            let sn = generate_serial(conn, project_id);
            let _ = conn.execute(
                "INSERT INTO wt_line (sn, project_id, content, position, node_id) VALUES (?, ?, ?, ?, ?)",
                rusqlite::params![sn, project_id, new_line, last_line_id, node_id],
            );
            let line_id = conn.last_insert_rowid();
            let _ = conn.execute(
                "INSERT INTO wt_node_change (action, level, node_id, target, detail_line) VALUES ('add', 'line', ?, ?, ?)",
                rusqlite::params![node_id, block_change_id, line_id],
            );
            last_line_id = Some(line_id);
        }
    }
    for line in &old_lines {
        if !line.used {
            let del_line_position: Option<i64> = conn
                .query_row(
                    "SELECT position FROM wt_line WHERE id = ?",
                    [line.line_id],
                    |row| row.get(0),
                )
                .ok()
                .flatten();
            let _ = conn.execute(
                "UPDATE wt_line SET position = ? WHERE position = ?",
                rusqlite::params![del_line_position, line.line_id],
            );
            let _ = conn.execute(
                "INSERT INTO wt_node_change (action, level, node_id, target) VALUES ('del', 'line', ?, ?)",
                rusqlite::params![node_id, line.change_id],
            );
        }
    }
}

#[tauri::command]
pub fn get_system_folder_change_ids(project_id: i64) -> Result<HashMap<String, i64>, String> {
    database::with_db(|conn| {
        let root_node_id: i64 = conn.query_row(
            "SELECT id FROM wt_node WHERE project_id = ? AND parent_id IS NULL",
            [project_id],
            |row| row.get(0),
        )?;
        let mut stmt = conn.prepare(
            "SELECT f.name, nc.id
             FROM wt_node_change nc
             JOIN wt_folder f ON nc.detail_folder = f.id
             WHERE nc.node_id = ? AND nc.level = 'folder' AND nc.action = 'add'
             AND f.name IN ('SYS_PRE', 'SYS_POST', 'SYS_GLOBAL')",
        )?;
        let mut result: HashMap<String, i64> = HashMap::new();
        let rows = stmt.query_map([root_node_id], |row| {
            Ok((row.get::<_, String>(0)?, row.get::<_, i64>(1)?))
        })?;
        for row in rows.flatten() {
            result.insert(row.0, row.1);
        }
        Ok(result)
    })
    .map_err(|e| e.to_string())
}

use crate::types::{ContentItem, LineMarker};

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct SyncDotStatus {
    pub line_index: i32,
    pub serial: String,
    pub shape: String,
    pub color: String,
}

#[tauri::command]
pub fn compute_effective_content(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
    _content_item_serial: String,
) -> Result<String, String> {
    Ok(String::new())
}

#[tauri::command]
pub fn compute_sync_dot_status(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
    _content_item_serial: String,
) -> Result<Vec<SyncDotStatus>, String> {
    Ok(vec![])
}

#[tauri::command]
pub fn add_line_marker(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
    _marker: LineMarker,
) -> Result<(), String> {
    Ok(())
}

#[tauri::command]
pub fn remove_line_marker(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
    _serial: String,
    _marker_type: String,
) -> Result<(), String> {
    Ok(())
}

#[tauri::command]
pub fn delete_line_in_node(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
    _line_serial: String,
    _content_item_serial: String,
) -> Result<(), String> {
    Ok(())
}

#[tauri::command]
pub fn add_line_in_node(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
    _content: String,
    _insert_before_serial: String,
    _new_serial: String,
    _content_item_serial: String,
) -> Result<(), String> {
    Ok(())
}

#[tauri::command]
pub fn get_node_line_markers(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
) -> Result<Vec<LineMarker>, String> {
    Ok(vec![])
}

#[tauri::command]
pub fn compute_node_inherited_content_items(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
) -> Result<Vec<ContentItem>, String> {
    Ok(vec![])
}

#[tauri::command]
pub fn compute_node_effective_content_items(
    _project_id: i64,
    _node_id: i64,
    _entry_id: String,
) -> Result<Vec<ContentItem>, String> {
    Ok(vec![])
}

use crate::types::{DeltaResult, NodeDelta, TimelineNode, WorldBookEntry, WorldBookWithLine};

#[tauri::command]
pub fn compute_node_delta(
    _parent_node: TimelineNode,
    _child_node: TimelineNode,
) -> Result<NodeDelta, String> {
    Ok(NodeDelta {
        node_id: String::new(),
        parent_id: String::new(),
        wb_added: vec![],
        wb_removed: vec![],
        wb_modified: vec![],
        pre_prompt_delta: vec![],
        pre_text_delta: vec![],
        post_text_delta: vec![],
        name: None,
        note: None,
        tags: None,
    })
}

#[tauri::command]
pub fn apply_delta(_base_node: TimelineNode, _delta: NodeDelta) -> Result<TimelineNode, String> {
    Err("Not implemented".to_string())
}

#[tauri::command]
pub fn get_node_with_delta(
    _project_id: String,
    _node_id: String,
    _timeline: Vec<TimelineNode>,
) -> Result<DeltaResult, String> {
    Ok(DeltaResult {
        is_root: true,
        delta: None,
        full_node: None,
        compression_ratio: 1.0,
    })
}

#[tauri::command]
pub fn reconstruct_node_from_delta(
    _project_id: String,
    _node_id: String,
    _timeline: Vec<TimelineNode>,
) -> Result<TimelineNode, String> {
    Err("Not implemented".to_string())
}

#[tauri::command]
pub fn add_line_numbers_to_world_book(world_book: Vec<WorldBookEntry>) -> Vec<WorldBookWithLine> {
    world_book
        .into_iter()
        .enumerate()
        .map(|(i, entry)| WorldBookWithLine {
            entry,
            line_no: (i + 1) as i32,
        })
        .collect()
}
