use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use std::collections::HashMap;
use std::sync::RwLock;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DirtyItem {
    pub id: String,
    #[serde(rename = "type")]
    pub item_type: String,
    pub original_hash: String,
    pub current_hash: String,
    pub is_dirty: bool,
    pub last_modified: DateTime<Utc>,
    pub last_saved: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DirtyState {
    pub id: String,
    pub is_dirty: bool,
    pub minutes_since_last_save: i64,
}

pub struct DirtyManager {
    items: RwLock<HashMap<String, DirtyItem>>,
}

impl DirtyManager {
    pub fn new() -> Self {
        Self {
            items: RwLock::new(HashMap::new()),
        }
    }
    fn compute_hash(content: &[u8]) -> String {
        let mut hasher = Sha256::new();
        hasher.update(content);
        format!("{:x}", hasher.finalize())
    }
    pub fn set_original(&self, id: &str, item_type: &str, content: &[u8]) {
        let hash = Self::compute_hash(content);
        let now = Utc::now();
        let mut items = self.items.write().unwrap();
        items.insert(
            id.to_string(),
            DirtyItem {
                id: id.to_string(),
                item_type: item_type.to_string(),
                original_hash: hash.clone(),
                current_hash: hash,
                is_dirty: false,
                last_modified: now,
                last_saved: now,
            },
        );
    }
    pub fn update_current(&self, id: &str, content: &[u8]) -> bool {
        let mut items = self.items.write().unwrap();
        if let Some(item) = items.get_mut(id) {
            let new_hash = Self::compute_hash(content);
            if new_hash != item.current_hash {
                item.current_hash = new_hash.clone();
                item.is_dirty = new_hash != item.original_hash;
                item.last_modified = Utc::now();
                return true;
            }
        }
        false
    }
    pub fn mark_saved(&self, id: &str, content: &[u8]) {
        let mut items = self.items.write().unwrap();
        if let Some(item) = items.get_mut(id) {
            let hash = Self::compute_hash(content);
            item.original_hash = hash.clone();
            item.current_hash = hash;
            item.is_dirty = false;
            item.last_saved = Utc::now();
        }
    }
    pub fn is_dirty(&self, id: &str) -> bool {
        let items = self.items.read().unwrap();
        items.get(id).map(|item| item.is_dirty).unwrap_or(false)
    }
    pub fn get_item(&self, id: &str) -> Option<DirtyItem> {
        let items = self.items.read().unwrap();
        items.get(id).cloned()
    }
    #[allow(dead_code)]
    pub fn get_dirty_items(&self) -> Vec<DirtyItem> {
        let items = self.items.read().unwrap();
        items.values().filter(|item| item.is_dirty).cloned().collect()
    }
    pub fn get_all_items(&self) -> Vec<DirtyItem> {
        let items = self.items.read().unwrap();
        items.values().cloned().collect()
    }
    pub fn clear_item(&self, id: &str) {
        let mut items = self.items.write().unwrap();
        items.remove(id);
    }
    pub fn has_any_dirty(&self) -> bool {
        let items = self.items.read().unwrap();
        items.values().any(|item| item.is_dirty)
    }
    pub fn get_minutes_since_last_save(&self, id: &str) -> i64 {
        let items = self.items.read().unwrap();
        items
            .get(id)
            .map(|item| Utc::now().signed_duration_since(item.last_saved).num_minutes())
            .unwrap_or(0)
    }
}

impl Default for DirtyManager {
    fn default() -> Self {
        Self::new()
    }
}

use std::sync::OnceLock;
static DIRTY_MANAGER: OnceLock<DirtyManager> = OnceLock::new();

fn get_manager() -> &'static DirtyManager {
    DIRTY_MANAGER.get_or_init(DirtyManager::new)
}

#[tauri::command]
pub fn get_dirty_state(id: String) -> Option<DirtyState> {
    let manager = get_manager();
    manager.get_item(&id).map(|item| DirtyState {
        id: item.id,
        is_dirty: item.is_dirty,
        minutes_since_last_save: manager.get_minutes_since_last_save(&id),
    })
}

#[tauri::command]
pub fn set_original_content(id: String, item_type: String, content: String) {
    get_manager().set_original(&id, &item_type, content.as_bytes());
}

#[tauri::command]
pub fn update_current_content(id: String, content: String) -> bool {
    get_manager().update_current(&id, content.as_bytes())
}

#[tauri::command]
pub fn mark_content_saved(id: String, content: String) {
    get_manager().mark_saved(&id, content.as_bytes());
}

#[tauri::command]
pub fn is_content_dirty(id: String) -> bool {
    get_manager().is_dirty(&id)
}

#[tauri::command]
pub fn has_any_dirty_content() -> bool {
    get_manager().has_any_dirty()
}

#[tauri::command]
pub fn get_all_dirty_states() -> Vec<DirtyState> {
    let manager = get_manager();
    manager
        .get_all_items()
        .into_iter()
        .map(|item| DirtyState {
            id: item.id.clone(),
            is_dirty: item.is_dirty,
            minutes_since_last_save: manager.get_minutes_since_last_save(&item.id),
        })
        .collect()
}

#[tauri::command]
pub fn clear_dirty_item(id: String) {
    get_manager().clear_item(&id);
}
