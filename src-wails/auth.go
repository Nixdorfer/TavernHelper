package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (a *App) AuthLogin(email, password string, rememberMe bool) (map[string]any, error) {
	url := "https://aipornhub.ltd/console/api/login"
	body := map[string]any{
		"email":              email,
		"password":           password,
		"remember_me":        rememberMe,
		"interface_language": "zh-Hans",
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) AuthGetProfile(token string) (map[string]any, error) {
	url := "https://aipornhub.ltd/go/api/account/profile"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) AuthLogout(token string) error {
	client := &http.Client{}

	req1, err := http.NewRequest("GET", "https://aipornhub.ltd/console/api/app_site/passport_callback_url?action=logout", nil)
	if err != nil {
		return err
	}
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	req1.Header.Set("Origin", "https://aipornhub.ltd")
	req1.Header.Set("Referer", "https://aipornhub.ltd/")
	req1.Header.Set("x-language", "zh-Hans")

	resp1, err := client.Do(req1)
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	respBody1, _ := io.ReadAll(resp1.Body)
	var result1 map[string]any
	if err := json.Unmarshal(respBody1, &result1); err == nil {
		if data, ok := result1["data"].(map[string]any); ok {
			if urls, ok := data["callback_urls"].([]any); ok && len(urls) > 0 {
				if callbackURL, ok := urls[0].(string); ok {
					reqCb, _ := http.NewRequest("GET", callbackURL, nil)
					reqCb.Header.Set("Accept", "application/json")
					reqCb.Header.Set("Referer", "https://aipornhub.ltd/")
					client.Do(reqCb)
				}
			}
		}
	}

	req2, err := http.NewRequest("GET", "https://aipornhub.ltd/console/api/logout", nil)
	if err != nil {
		return err
	}
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Origin", "https://aipornhub.ltd")
	req2.Header.Set("Referer", "https://aipornhub.ltd/")
	req2.Header.Set("x-language", "zh-Hans")

	resp2, err := client.Do(req2)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	return nil
}

func (a *App) AuthGetPoints(token, userId string) (map[string]any, error) {
	url := fmt.Sprintf("https://aipornhub.ltd/go/api/account/point?target=%s", userId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "https://aipornhub.ltd")
	req.Header.Set("Referer", "https://aipornhub.ltd/zh")
	req.Header.Set("x-language", "zh-Hans")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}
