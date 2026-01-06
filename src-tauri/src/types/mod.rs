use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTProject {
    pub id: i64,
    pub name: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub current_node: Option<i64>,
    pub create_time: String,
    pub update_time: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTNode {
    pub id: i64,
    pub project_id: i64,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub parent_id: Option<i64>,
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub desc: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTBranchTag {
    pub id: i64,
    pub parent_id: i64,
    pub child_id: i64,
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub desc: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WTFolder {
    pub id: i64,
    pub name: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTCard {
    pub id: i64,
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub desc: String,
    pub key_word: String,
    pub image_word: String,
    pub trigger_system: bool,
    pub trigger_user: bool,
    pub trigger_ai: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WTBlock {
    pub id: i64,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub title: String,
    pub zone: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTLine {
    pub id: i64,
    pub sn: String,
    pub project_id: i64,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub content: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub position: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NodeDetailBlock {
    pub title: String,
    pub lines: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NodeDetailTrigger {
    pub mode: String,
    pub words: Vec<String>,
    pub system: bool,
    pub user: bool,
    pub ai: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NodeDetailCard {
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub desc: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub trigger: Option<NodeDetailTrigger>,
    pub content: HashMap<String, NodeDetailBlock>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub image: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NodeDetailFolder {
    pub name: String,
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub cards: HashMap<String, NodeDetailCard>,
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub folders: HashMap<String, NodeDetailFolder>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct NodeDetail {
    pub node_id: i64,
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub desc: String,
    #[serde(rename = "struct")]
    pub structure: HashMap<String, NodeDetailFolder>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ImmediateChange {
    pub name: String,
    pub action: String,
    pub level: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub target: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SaveFolderChange {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub name: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub desc: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SaveCardChange {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub name: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub desc: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub trigger: Option<NodeDetailTrigger>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub image: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SaveBlockChange {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub name: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub content: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SaveChanges {
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub folder: HashMap<String, SaveFolderChange>,
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub card: HashMap<String, SaveCardChange>,
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub block: HashMap<String, SaveBlockChange>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTNodeChange {
    pub id: i64,
    pub action: String,
    pub level: String,
    pub node_id: i64,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub target: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub detail_folder: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub detail_card: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub detail_block: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub detail_line: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ProjectInfo {
    pub id: i64,
    pub name: String,
    pub file_name: String,
    #[serde(rename = "type")]
    pub project_type: String,
    pub updated_at: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Project {
    pub id: i64,
    pub name: String,
    pub file_name: String,
    #[serde(rename = "type")]
    pub project_type: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub current_node: Option<i64>,
    pub create_time: String,
    pub update_time: String,
    pub timeline: Vec<NodeInfo>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct NodeInfo {
    pub id: i64,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub parent_id: Option<i64>,
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub note: String,
    #[serde(default)]
    pub tags: Vec<String>,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub created_at: String,
    #[serde(default)]
    pub pre_text: Vec<PromptEntry>,
    #[serde(default)]
    pub post_text: Vec<PromptEntry>,
    #[serde(default)]
    pub pre_prompt: Vec<PromptEntry>,
    #[serde(default)]
    pub world_book: Vec<WorldBookEntry>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ReplayedFolder {
    pub id: i64,
    pub name: String,
    pub change_id: i64,
    pub cards: Vec<ReplayedCard>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ReplayedCard {
    pub id: i64,
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub desc: String,
    pub key_word: String,
    pub trigger_system: bool,
    pub trigger_user: bool,
    pub trigger_ai: bool,
    pub change_id: i64,
    pub blocks: Vec<ReplayedBlock>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ReplayedBlock {
    pub id: i64,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub title: String,
    pub zone: String,
    pub change_id: i64,
    pub lines: Vec<ReplayedLine>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ReplayedLine {
    pub id: i64,
    pub sn: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub content: Option<String>,
    pub sync_dot: String,
    pub change_id: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct NodeContent {
    pub node_id: i64,
    pub folders: Vec<ReplayedFolder>,
    pub pre: Vec<ReplayedBlock>,
    pub post: Vec<ReplayedBlock>,
    pub global: Vec<ReplayedBlock>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct LineMarker {
    pub serial: String,
    #[serde(rename = "type")]
    pub marker_type: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub content: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub insert_before: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ContentItem {
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub id: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub serial: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub title: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub content: String,
    #[serde(default)]
    pub key_system: bool,
    #[serde(default)]
    pub key_user: bool,
    #[serde(default)]
    pub key_ai: bool,
    #[serde(default)]
    pub key_region: i32,
    #[serde(default)]
    pub value_region: i32,
    #[serde(default)]
    pub collapsed: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WorldBookEntry {
    pub id: String,
    pub key: String,
    pub value: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub parent_id: String,
    #[serde(default)]
    pub is_folder: bool,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub match_mode: String,
    #[serde(default)]
    pub key_region: i32,
    #[serde(default)]
    pub value_region: i32,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub group: String,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub content_items: Vec<ContentItem>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub keywords: Vec<String>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub t2i_keywords: Vec<String>,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub source_entry_id: String,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub line_markers: Vec<LineMarker>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub deleted_item_ids: Vec<String>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub local_items: Vec<ContentItem>,
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub item_overrides: HashMap<String, HashMap<String, serde_json::Value>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PromptEntry {
    pub id: String,
    pub name: String,
    pub content: String,
    pub enabled: bool,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub source_entry_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct NodeBinding {
    pub app_id: String,
    pub app_name: String,
    pub conversation_id: String,
    pub config_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct TimelineNode {
    pub id: String,
    pub name: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub note: String,
    #[serde(default)]
    pub tags: Vec<String>,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub parent_id: String,
    pub created_at: String,
    #[serde(default)]
    pub pre_text: Vec<PromptEntry>,
    #[serde(default)]
    pub post_text: Vec<PromptEntry>,
    #[serde(default)]
    pub pre_prompt: Vec<PromptEntry>,
    #[serde(default)]
    pub world_book: Vec<WorldBookEntry>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub binding: Option<NodeBinding>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct Config {
    #[serde(default)]
    pub theme: String,
    #[serde(default)]
    pub language: String,
    #[serde(default)]
    pub last_opened_project: String,
    #[serde(default)]
    pub color_scheme: i32,
    #[serde(default)]
    pub color_mode: String,
    #[serde(default)]
    pub system_prompt: String,
    #[serde(default)]
    pub system_prompt_type: String,
    #[serde(default)]
    pub debug_mode: bool,
    #[serde(default)]
    pub safe_mode: bool,
    #[serde(default)]
    pub safe_mode_action: String,
    #[serde(default)]
    pub safe_mode_template: String,
    #[serde(default)]
    pub debug_test_reply: String,
    #[serde(default)]
    pub bytedance_api_key: String,
    #[serde(default)]
    pub auto_generate_image: bool,
    #[serde(default)]
    pub no_image_mode: bool,
    #[serde(default)]
    pub strict_mode: bool,
    #[serde(default)]
    pub claude_api_key: String,
    #[serde(default)]
    pub gemini_api_key: String,
    #[serde(default)]
    pub grok_api_key: String,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
#[serde(rename_all = "camelCase")]
pub struct SessionState {
    #[serde(default)]
    pub open_tabs: Vec<OpenTabState>,
    #[serde(default)]
    pub active_tab_index: i32,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub active_conversation: Option<ConversationState>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct OpenTabState {
    pub project_file_name: String,
    pub current_node_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ConversationState {
    pub app_id: String,
    pub conversation_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Draft {
    pub id: i64,
    pub name: String,
    pub content: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub parent_id: Option<i64>,
    pub created_at: String,
    pub updated_at: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ClipboardCapture {
    pub id: i64,
    pub content: String,
    pub captured_at: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct GalleryImage {
    pub id: String,
    pub hash: String,
    pub local_path: String,
    #[serde(default)]
    pub remote_url: String,
    #[serde(default)]
    pub file_name: String,
    #[serde(default)]
    pub file_size: i64,
    pub created_at: String,
    #[serde(default)]
    pub folder_path: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct GalleryFolder {
    pub name: String,
    pub path: String,
    pub image_count: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ChatRequest {
    pub messages: Vec<ChatMessage>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub system_prompt: Option<String>,
    pub model: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub max_tokens: Option<i32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChatMessage {
    pub role: String,
    pub content: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChatResponse {
    pub content: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub error: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WorldBookPatch {
    #[serde(rename = "ln")]
    pub line_no: i32,
    #[serde(rename = "id")]
    pub entry_id: String,
    pub op: String,
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub changes: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PromptPatch {
    #[serde(rename = "ln")]
    pub line_no: i32,
    pub id: String,
    pub op: String,
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub changes: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct NodeDelta {
    pub node_id: String,
    pub parent_id: String,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub wb_added: Vec<WorldBookEntry>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub wb_removed: Vec<String>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub wb_modified: Vec<WorldBookPatch>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub pre_prompt_delta: Vec<PromptPatch>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub pre_text_delta: Vec<PromptPatch>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub post_text_delta: Vec<PromptPatch>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub name: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub note: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tags: Option<Vec<String>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DeltaResult {
    pub is_root: bool,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub delta: Option<NodeDelta>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub full_node: Option<TimelineNode>,
    pub compression_ratio: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WorldBookWithLine {
    #[serde(flatten)]
    pub entry: WorldBookEntry,
    #[serde(rename = "ln")]
    pub line_no: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTApp {
    pub id: i64,
    pub name: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTConversation {
    pub id: i64,
    pub app_id: i64,
    pub name: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub current_node: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTDialogue {
    pub id: i64,
    pub conversation_id: i64,
    pub create_time: String,
    pub request_content: String,
    pub response_content: String,
    #[serde(default)]
    pub request_system_prompt: String,
    #[serde(default)]
    pub response_system_prompt: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub node_id: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub request_point: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub response_point: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub request_token: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub response_token: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct WTDialogueImage {
    pub id: i64,
    pub dialogue_id: i64,
    #[serde(default)]
    pub image_url: String,
    #[serde(default)]
    pub image_path: String,
    pub prompt: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DialoguesResult {
    pub dialogues: Vec<WTDialogue>,
    pub total: i64,
}
