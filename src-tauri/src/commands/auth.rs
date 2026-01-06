fn build_request(client: &reqwest::Client, method: reqwest::Method, url: &str) -> reqwest::RequestBuilder {
    client.request(method, url)
        .header("Content-Type", "application/json")
        .header("Accept", "application/json")
        .header("Origin", "https://aipornhub.ltd")
        .header("Referer", "https://aipornhub.ltd/")
        .header("x-language", "zh-Hans")
}

#[tauri::command]
pub async fn auth_login(email: String, password: String, remember_me: bool) -> Result<serde_json::Value, String> {
    let client = reqwest::Client::new();
    let resp = build_request(&client, reqwest::Method::POST, "https://aipornhub.ltd/console/api/login")
        .json(&serde_json::json!({
            "email": email,
            "password": password,
            "remember_me": remember_me,
            "interface_language": "zh-Hans"
        }))
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn auth_get_profile(token: String) -> Result<serde_json::Value, String> {
    let client = reqwest::Client::new();
    let resp = build_request(&client, reqwest::Method::GET, "https://aipornhub.ltd/go/api/account/profile")
        .header("Authorization", format!("Bearer {}", token))
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}

#[tauri::command]
pub async fn auth_logout(token: String) -> Result<(), String> {
    let client = reqwest::Client::new();
    let resp1 = build_request(&client, reqwest::Method::GET, "https://aipornhub.ltd/console/api/app_site/passport_callback_url?action=logout")
        .header("Authorization", format!("Bearer {}", token))
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text1 = resp1.text().await.unwrap_or_default();
    if let Ok(result1) = serde_json::from_str::<serde_json::Value>(&text1) {
        if let Some(urls) = result1.get("data")
            .and_then(|d| d.get("callback_urls"))
            .and_then(|u| u.as_array()) {
            if let Some(callback_url) = urls.first().and_then(|u| u.as_str()) {
                let _ = client.get(callback_url)
                    .header("Accept", "application/json")
                    .header("Referer", "https://aipornhub.ltd/")
                    .send()
                    .await;
            }
        }
    }
    let _ = build_request(&client, reqwest::Method::GET, "https://aipornhub.ltd/console/api/logout")
        .header("Authorization", format!("Bearer {}", token))
        .send()
        .await
        .map_err(|e| e.to_string())?;
    Ok(())
}

#[tauri::command]
pub async fn auth_get_points(token: String, user_id: String) -> Result<serde_json::Value, String> {
    let url = format!("https://aipornhub.ltd/go/api/account/point?target={}", user_id);
    let client = reqwest::Client::new();
    let resp = build_request(&client, reqwest::Method::GET, &url)
        .header("Authorization", format!("Bearer {}", token))
        .send()
        .await
        .map_err(|e| e.to_string())?;
    let text = resp.text().await.map_err(|e| e.to_string())?;
    serde_json::from_str(&text).map_err(|e| e.to_string())
}
