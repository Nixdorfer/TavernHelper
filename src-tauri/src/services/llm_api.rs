use crate::types::{ChatRequest, ChatResponse};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize)]
struct ClaudeRequest {
    model: String,
    max_tokens: i32,
    messages: Vec<ClaudeMessage>,
    #[serde(skip_serializing_if = "Option::is_none")]
    system: Option<String>,
}

#[derive(Debug, Serialize)]
struct ClaudeMessage {
    role: String,
    content: String,
}

#[derive(Debug, Deserialize)]
struct ClaudeResponse {
    content: Vec<ClaudeContent>,
}

#[derive(Debug, Deserialize)]
struct ClaudeContent {
    text: String,
}

#[derive(Debug, Serialize)]
struct GeminiRequest {
    contents: Vec<GeminiContent>,
    #[serde(skip_serializing_if = "Option::is_none")]
    system_instruction: Option<GeminiSystemInstruction>,
}

#[derive(Debug, Serialize)]
struct GeminiContent {
    role: String,
    parts: Vec<GeminiPart>,
}

#[derive(Debug, Serialize, Deserialize)]
struct GeminiPart {
    text: String,
}

#[derive(Debug, Serialize)]
struct GeminiSystemInstruction {
    parts: Vec<GeminiPart>,
}

#[derive(Debug, Deserialize)]
struct GeminiResponse {
    candidates: Vec<GeminiCandidate>,
}

#[derive(Debug, Deserialize)]
struct GeminiCandidate {
    content: GeminiCandidateContent,
}

#[derive(Debug, Deserialize)]
struct GeminiCandidateContent {
    parts: Vec<GeminiPart>,
}

pub async fn send_chat(request: ChatRequest) -> Result<ChatResponse, String> {
    let model = request.model.to_lowercase();
    if model.contains("claude") {
        call_claude_api(request).await
    } else if model.contains("gemini") {
        call_gemini_api(request).await
    } else if model.contains("grok") {
        call_grok_api(request).await
    } else {
        Err(format!("Unsupported model: {}", request.model))
    }
}

async fn call_claude_api(request: ChatRequest) -> Result<ChatResponse, String> {
    let api_key = crate::database::db_get_config("claude_api_key")
        .map_err(|e| e.to_string())?
        .unwrap_or_default();
    if api_key.is_empty() {
        return Err("Claude API key not configured".to_string());
    }
    let messages: Vec<ClaudeMessage> = request
        .messages
        .iter()
        .map(|m| ClaudeMessage {
            role: m.role.clone(),
            content: m.content.clone(),
        })
        .collect();
    let claude_request = ClaudeRequest {
        model: request.model,
        max_tokens: request.max_tokens.unwrap_or(4096),
        messages,
        system: request.system_prompt,
    };
    let client = reqwest::Client::new();
    let resp = client
        .post("https://api.anthropic.com/v1/messages")
        .header("x-api-key", &api_key)
        .header("anthropic-version", "2023-06-01")
        .header("Content-Type", "application/json")
        .json(&claude_request)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    if !resp.status().is_success() {
        let error_text = resp.text().await.unwrap_or_default();
        return Ok(ChatResponse {
            content: String::new(),
            error: Some(error_text),
        });
    }
    let claude_response: ClaudeResponse = resp.json().await.map_err(|e| e.to_string())?;
    let content = claude_response
        .content
        .first()
        .map(|c| c.text.clone())
        .unwrap_or_default();
    Ok(ChatResponse {
        content,
        error: None,
    })
}

async fn call_gemini_api(request: ChatRequest) -> Result<ChatResponse, String> {
    let api_key = crate::database::db_get_config("gemini_api_key")
        .map_err(|e| e.to_string())?
        .unwrap_or_default();
    if api_key.is_empty() {
        return Err("Gemini API key not configured".to_string());
    }
    let contents: Vec<GeminiContent> = request
        .messages
        .iter()
        .map(|m| GeminiContent {
            role: if m.role == "assistant" { "model".to_string() } else { "user".to_string() },
            parts: vec![GeminiPart { text: m.content.clone() }],
        })
        .collect();
    let system_instruction = request.system_prompt.map(|s| GeminiSystemInstruction {
        parts: vec![GeminiPart { text: s }],
    });
    let gemini_request = GeminiRequest {
        contents,
        system_instruction,
    };
    let url = format!(
        "https://generativelanguage.googleapis.com/v1beta/models/{}:generateContent?key={}",
        request.model, api_key
    );
    let client = reqwest::Client::new();
    let resp = client
        .post(&url)
        .header("Content-Type", "application/json")
        .json(&gemini_request)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    if !resp.status().is_success() {
        let error_text = resp.text().await.unwrap_or_default();
        return Ok(ChatResponse {
            content: String::new(),
            error: Some(error_text),
        });
    }
    let gemini_response: GeminiResponse = resp.json().await.map_err(|e| e.to_string())?;
    let content = gemini_response
        .candidates
        .first()
        .and_then(|c| c.content.parts.first())
        .map(|p| p.text.clone())
        .unwrap_or_default();
    Ok(ChatResponse {
        content,
        error: None,
    })
}

async fn call_grok_api(request: ChatRequest) -> Result<ChatResponse, String> {
    let api_key = crate::database::db_get_config("grok_api_key")
        .map_err(|e| e.to_string())?
        .unwrap_or_default();
    if api_key.is_empty() {
        return Err("Grok API key not configured".to_string());
    }
    let mut messages: Vec<serde_json::Value> = vec![];
    if let Some(system) = &request.system_prompt {
        messages.push(serde_json::json!({
            "role": "system",
            "content": system
        }));
    }
    for m in &request.messages {
        messages.push(serde_json::json!({
            "role": m.role,
            "content": m.content
        }));
    }
    let grok_request = serde_json::json!({
        "model": request.model,
        "messages": messages,
        "max_tokens": request.max_tokens.unwrap_or(4096)
    });
    let client = reqwest::Client::new();
    let resp = client
        .post("https://api.x.ai/v1/chat/completions")
        .header("Authorization", format!("Bearer {}", api_key))
        .header("Content-Type", "application/json")
        .json(&grok_request)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    if !resp.status().is_success() {
        let error_text = resp.text().await.unwrap_or_default();
        return Ok(ChatResponse {
            content: String::new(),
            error: Some(error_text),
        });
    }
    let grok_response: serde_json::Value = resp.json().await.map_err(|e| e.to_string())?;
    let content = grok_response["choices"][0]["message"]["content"]
        .as_str()
        .unwrap_or_default()
        .to_string();
    Ok(ChatResponse {
        content,
        error: None,
    })
}
