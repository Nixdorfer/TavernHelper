use crate::database;
use crate::types::{Config, SessionState};

#[tauri::command]
pub fn load_config() -> Result<Config, String> {
    let config_json = database::db_get_config("app_config")
        .map_err(|e| e.to_string())?
        .unwrap_or_default();
    if config_json.is_empty() {
        return Ok(Config::default());
    }
    serde_json::from_str(&config_json).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn save_config(config: Config) -> Result<(), String> {
    let config_json = serde_json::to_string(&config).map_err(|e| e.to_string())?;
    database::db_set_config("app_config", &config_json).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn update_config(key: String, value: serde_json::Value) -> Result<(), String> {
    let mut config = load_config()?;
    match key.as_str() {
        "theme" => config.theme = value.as_str().unwrap_or_default().to_string(),
        "language" => config.language = value.as_str().unwrap_or_default().to_string(),
        "lastOpenedProject" => config.last_opened_project = value.as_str().unwrap_or_default().to_string(),
        "colorScheme" => config.color_scheme = value.as_i64().unwrap_or(0) as i32,
        "colorMode" => config.color_mode = value.as_str().unwrap_or_default().to_string(),
        "systemPrompt" => config.system_prompt = value.as_str().unwrap_or_default().to_string(),
        "systemPromptType" => config.system_prompt_type = value.as_str().unwrap_or_default().to_string(),
        "debugMode" => config.debug_mode = value.as_bool().unwrap_or(false),
        "safeMode" => config.safe_mode = value.as_bool().unwrap_or(false),
        "safeModeAction" => config.safe_mode_action = value.as_str().unwrap_or_default().to_string(),
        "safeModeTemplate" => config.safe_mode_template = value.as_str().unwrap_or_default().to_string(),
        "debugTestReply" => config.debug_test_reply = value.as_str().unwrap_or_default().to_string(),
        "bytedanceApiKey" => config.bytedance_api_key = value.as_str().unwrap_or_default().to_string(),
        "autoGenerateImage" => config.auto_generate_image = value.as_bool().unwrap_or(false),
        "noImageMode" => config.no_image_mode = value.as_bool().unwrap_or(false),
        "strictMode" => config.strict_mode = value.as_bool().unwrap_or(false),
        "claudeApiKey" => config.claude_api_key = value.as_str().unwrap_or_default().to_string(),
        "geminiApiKey" => config.gemini_api_key = value.as_str().unwrap_or_default().to_string(),
        "grokApiKey" => config.grok_api_key = value.as_str().unwrap_or_default().to_string(),
        _ => return Err(format!("Unknown config key: {}", key)),
    }
    save_config(config)
}

#[tauri::command]
pub fn load_session_state() -> Result<SessionState, String> {
    let state_json = database::db_get_config("session_state")
        .map_err(|e| e.to_string())?
        .unwrap_or_default();
    if state_json.is_empty() {
        return Ok(SessionState::default());
    }
    serde_json::from_str(&state_json).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn save_session_state(state: SessionState) -> Result<(), String> {
    let state_json = serde_json::to_string(&state).map_err(|e| e.to_string())?;
    database::db_set_config("session_state", &state_json).map_err(|e| e.to_string())
}
