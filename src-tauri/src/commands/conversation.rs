fn build_request(client: &reqwest::Client, method: reqwest::Method, url: &str, token: &str) -> reqwest::RequestBuilder {
    client.request(method, url)
        .header("Content-Type", "application/json")
        .header("Accept", "*/*")
        .header("Authorization", format!("Bearer {}", token))
        .header("Origin", "https://aipornhub.ltd")
        .header("Referer", "https://aipornhub.ltd/")
        .header("x-language", "zh-Hans")
}

#[tauri::command]
pub async fn get_conversations(token: String, app_id: String, page: i64, limit: i64) -> Result<serde_json::Value, String> {
    let url = format!("https://aipornhub.ltd/console/api/installed-apps/{}/conversations?page={}&limit={}", app_id, page, limit);
    let client = reqwest::Client::new();
    let resp = build_request(&client, reqwest::Method::GET, &url, &token)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn get_conversation_detail(token: String, app_id: String, conversation_id: String) -> Result<serde_json::Value, String> {
    let url = format!("https://aipornhub.ltd/console/api/installed-apps/{}/conversations/{}", app_id, conversation_id);
    let client = reqwest::Client::new();
    let resp = build_request(&client, reqwest::Method::GET, &url, &token)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    if !resp.status().is_success() {
        let body = resp.text().await.unwrap_or_default();
        return Err(format!("获取会话详情失败: {}", body));
    }
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn delete_conversation(token: String, app_id: String, conversation_id: String) -> Result<(), String> {
    let url = format!("https://aipornhub.ltd/console/api/installed-apps/{}/conversations/{}", app_id, conversation_id);
    let client = reqwest::Client::new();
    let resp = build_request(&client, reqwest::Method::DELETE, &url, &token)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let status = resp.status();
    if status.as_u16() != 200 && status.as_u16() != 204 {
        let body = resp.text().await.unwrap_or_default();
        return Err(format!("删除对话失败: {}", body));
    }
    Ok(())
}

#[tauri::command]
pub async fn rename_conversation(token: String, app_id: String, conversation_id: String, new_name: String) -> Result<(), String> {
    let url = format!("https://aipornhub.ltd/console/api/installed-apps/{}/conversations/{}/name", app_id, conversation_id);
    let client = reqwest::Client::new();
    let resp = build_request(&client, reqwest::Method::POST, &url, &token)
        .json(&serde_json::json!({ "name": new_name }))
        .send()
        .await
        .map_err(|e| e.to_string())?;
    if !resp.status().is_success() {
        let body = resp.text().await.unwrap_or_default();
        return Err(format!("重命名对话失败: {}", body));
    }
    Ok(())
}

#[tauri::command]
pub async fn create_new_conversation(token: String, app_id: String, query: String, conversation_name: String) -> Result<serde_json::Value, String> {
    let url = format!("https://aipornhub.ltd/console/api/installed-apps/{}/chat-messages", app_id);
    let body = serde_json::json!({
        "response_mode": "blocking",
        "conversation_id": "",
        "query": query,
        "inputs": {},
        "conversation_name": conversation_name,
        "history_start_at": null
    });
    let client = reqwest::Client::builder()
        .timeout(std::time::Duration::from_secs(120))
        .build()
        .map_err(|e| e.to_string())?;
    let resp = build_request(&client, reqwest::Method::POST, &url, &token)
        .json(&body)
        .send()
        .await
        .map_err(|e| e.to_string())?;
    if !resp.status().is_success() {
        let body = resp.text().await.unwrap_or_default();
        return Err(format!("创建对话失败: {}", body));
    }
    let text = resp.text().await.map_err(|e| e.to_string())?;
    let result: serde_json::Value = serde_json::from_str(&text).map_err(|e| e.to_string())?;
    if let Some(conv_id) = result.get("conversation_id").and_then(|v| v.as_str()) {
        if !conv_id.is_empty() {
            return Ok(serde_json::json!({ "conversation_id": conv_id }));
        }
    }
    Err("未能从响应中获取 conversation_id".to_string())
}
