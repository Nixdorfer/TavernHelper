package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	execPath, _ := os.Executable()
	dbPath := filepath.Join(filepath.Dir(execPath), "worldbook.db")
	var err error
	db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := createTables(); err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}
	if err := initSchema(db); err != nil {
		return fmt.Errorf("创建表结构失败: %w", err)
	}
	return nil
}
func GetDB() *sql.DB {
	return db
}
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func createTables() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS wt_config (
			key TEXT PRIMARY KEY,
			value TEXT
		);
		CREATE TABLE IF NOT EXISTS wt_draft (
			id TEXT PRIMARY KEY,
			name TEXT DEFAULT '',
			content TEXT DEFAULT '',
			parent_id TEXT DEFAULT '',
			is_folder INTEGER DEFAULT 0,
			sort_order INTEGER DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_wt_draft_parent ON wt_draft(parent_id);
		CREATE TABLE IF NOT EXISTS wt_clipboard (
			id TEXT PRIMARY KEY,
			content TEXT DEFAULT '',
			created_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_wt_clipboard_created_at ON wt_clipboard(created_at DESC);
		CREATE TABLE IF NOT EXISTS wt_creation (
			id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			app_name TEXT DEFAULT '',
			app_icon TEXT DEFAULT '',
			model_config TEXT DEFAULT '{}',
			updated_at TEXT NOT NULL,
			PRIMARY KEY (id, user_id)
		);
		CREATE INDEX IF NOT EXISTS idx_wt_creation_user ON wt_creation(user_id);
	`)
	return err
}
func DBGetConfig(key string) (string, error) {
	var value string
	err := db.QueryRow("SELECT value FROM wt_config WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}
func DBSetConfig(key, value string) error {
	_, err := db.Exec(`
		INSERT INTO wt_config (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)
	return err
}
func (a *App) GetApps() ([]WTApp, error) {
	rows, err := db.Query(`SELECT id, name FROM wt_app ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []WTApp
	for rows.Next() {
		var app WTApp
		if err := rows.Scan(&app.ID, &app.Name); err != nil {
			continue
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (a *App) CreateApp(name string) (int, error) {
	result, err := db.Exec(`INSERT INTO wt_app (name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (a *App) DeleteApp(appID int) error {
	_, err := db.Exec(`DELETE FROM wt_app WHERE id = ?`, appID)
	return err
}

func (a *App) RenameApp(appID int, newName string) error {
	_, err := db.Exec(`UPDATE wt_app SET name = ? WHERE id = ?`, newName, appID)
	return err
}

func (a *App) GetLocalConversations(appID int) ([]WTConversation, error) {
	rows, err := db.Query(`SELECT id, app_id, name, current_node FROM wt_conversation WHERE app_id = ? ORDER BY id DESC`, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var conversations []WTConversation
	for rows.Next() {
		var conv WTConversation
		var currentNode sql.NullInt64
		if err := rows.Scan(&conv.ID, &conv.AppID, &conv.Name, &currentNode); err != nil {
			continue
		}
		if currentNode.Valid {
			node := int(currentNode.Int64)
			conv.CurrentNode = &node
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}

func (a *App) CreateConversation(appID int, name string) (int, error) {
	result, err := db.Exec(`INSERT INTO wt_conversation (app_id, name) VALUES (?, ?)`, appID, name)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (a *App) DeleteLocalConversation(conversationID int) error {
	_, err := db.Exec(`DELETE FROM wt_conversation WHERE id = ?`, conversationID)
	return err
}

func (a *App) RenameLocalConversation(conversationID int, newName string) error {
	_, err := db.Exec(`UPDATE wt_conversation SET name = ? WHERE id = ?`, newName, conversationID)
	return err
}

func (a *App) SetConversationNode(conversationID int, nodeID *int) error {
	if nodeID == nil {
		_, err := db.Exec(`UPDATE wt_conversation SET current_node = NULL WHERE id = ?`, conversationID)
		return err
	}
	_, err := db.Exec(`UPDATE wt_conversation SET current_node = ? WHERE id = ?`, *nodeID, conversationID)
	return err
}

func (a *App) GetConversationNode(conversationID int) (*int, error) {
	var currentNode sql.NullInt64
	err := db.QueryRow(`SELECT current_node FROM wt_conversation WHERE id = ?`, conversationID).Scan(&currentNode)
	if err != nil {
		return nil, err
	}
	if currentNode.Valid {
		node := int(currentNode.Int64)
		return &node, nil
	}
	return nil, nil
}

func (a *App) GetDialogues(conversationID int, page, limit int) ([]WTDialogue, int, error) {
	var total int
	if err := db.QueryRow(`SELECT COUNT(*) FROM wt_dialogue WHERE conversation_id = ?`, conversationID).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	rows, err := db.Query(`
		SELECT id, conversation_id, create_time, request_content, response_content,
			   COALESCE(request_system_prompt, ''), COALESCE(response_system_prompt, ''),
			   node_id, request_point, response_point, request_token, response_token
		FROM wt_dialogue WHERE conversation_id = ? ORDER BY create_time DESC LIMIT ? OFFSET ?
	`, conversationID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var dialogues []WTDialogue
	for rows.Next() {
		var d WTDialogue
		var nodeID, reqPoint, respPoint, reqToken, respToken sql.NullInt64
		if err := rows.Scan(&d.ID, &d.ConversationID, &d.CreateTime, &d.RequestContent, &d.ResponseContent,
			&d.RequestSystemPrompt, &d.ResponseSystemPrompt, &nodeID, &reqPoint, &respPoint, &reqToken, &respToken); err != nil {
			continue
		}
		if nodeID.Valid {
			n := int(nodeID.Int64)
			d.NodeID = &n
		}
		if reqPoint.Valid {
			n := int(reqPoint.Int64)
			d.RequestPoint = &n
		}
		if respPoint.Valid {
			n := int(respPoint.Int64)
			d.ResponsePoint = &n
		}
		if reqToken.Valid {
			n := int(reqToken.Int64)
			d.RequestToken = &n
		}
		if respToken.Valid {
			n := int(respToken.Int64)
			d.ResponseToken = &n
		}
		dialogues = append(dialogues, d)
	}
	return dialogues, total, nil
}

func (a *App) CreateDialogue(conversationID int, requestContent, responseContent string) (int, error) {
	now := time.Now().Format(time.RFC3339)
	result, err := db.Exec(`
		INSERT INTO wt_dialogue (conversation_id, create_time, request_content, response_content)
		VALUES (?, ?, ?, ?)
	`, conversationID, now, requestContent, responseContent)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (a *App) UpdateDialogue(dialogueID int, requestContent, responseContent string) error {
	_, err := db.Exec(`UPDATE wt_dialogue SET request_content = ?, response_content = ? WHERE id = ?`,
		requestContent, responseContent, dialogueID)
	return err
}

func (a *App) DeleteDialogue(dialogueID int) error {
	_, err := db.Exec(`DELETE FROM wt_dialogue WHERE id = ?`, dialogueID)
	return err
}

func (a *App) AddDialogueImage(dialogueID int, imageURL, imagePath, prompt string) (int, error) {
	result, err := db.Exec(`INSERT INTO wt_dialogue_image (dialogue_id, image_url, image_path, prompt) VALUES (?, ?, ?, ?)`,
		dialogueID, imageURL, imagePath, prompt)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (a *App) GetDialogueImages(dialogueID int) ([]WTDialogueImage, error) {
	rows, err := db.Query(`SELECT id, dialogue_id, COALESCE(image_url, ''), COALESCE(image_path, ''), prompt FROM wt_dialogue_image WHERE dialogue_id = ?`, dialogueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var images []WTDialogueImage
	for rows.Next() {
		var img WTDialogueImage
		if err := rows.Scan(&img.ID, &img.DialogueID, &img.ImageURL, &img.ImagePath, &img.Prompt); err != nil {
			continue
		}
		images = append(images, img)
	}
	return images, nil
}

type Draft struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	ParentID  string `json:"parentId"`
	IsFolder  bool   `json:"isFolder"`
	SortOrder int    `json:"sortOrder"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func DBGetAllDrafts() ([]Draft, error) {
	rows, err := db.Query(`
		SELECT id, name, content, parent_id, is_folder, sort_order, created_at, updated_at
		FROM wt_draft ORDER BY sort_order, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var drafts []Draft
	for rows.Next() {
		var d Draft
		var isFolder int
		if err := rows.Scan(&d.ID, &d.Name, &d.Content, &d.ParentID, &isFolder, &d.SortOrder, &d.CreatedAt, &d.UpdatedAt); err != nil {
			continue
		}
		d.IsFolder = isFolder == 1
		drafts = append(drafts, d)
	}
	return drafts, nil
}

func DBCreateDraft(draft *Draft) error {
	isFolder := 0
	if draft.IsFolder {
		isFolder = 1
	}
	_, err := db.Exec(`
		INSERT INTO wt_draft (id, name, content, parent_id, is_folder, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, draft.ID, draft.Name, draft.Content, draft.ParentID, isFolder, draft.SortOrder, draft.CreatedAt, draft.UpdatedAt)
	return err
}

func DBUpdateDraft(draft *Draft) error {
	isFolder := 0
	if draft.IsFolder {
		isFolder = 1
	}
	_, err := db.Exec(`
		UPDATE wt_draft SET name = ?, content = ?, parent_id = ?, is_folder = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`, draft.Name, draft.Content, draft.ParentID, isFolder, draft.SortOrder, draft.UpdatedAt, draft.ID)
	return err
}

func DBDeleteDraft(id string) error {
	_, err := db.Exec("DELETE FROM wt_draft WHERE id = ?", id)
	return err
}

func DBDeleteDraftRecursive(id string) error {
	toDelete := []string{id}
	collectDraftChildren(id, &toDelete)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, draftID := range toDelete {
		if _, err := tx.Exec("DELETE FROM wt_draft WHERE id = ?", draftID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func collectDraftChildren(parentID string, result *[]string) {
	rows, err := db.Query("SELECT id FROM wt_draft WHERE parent_id = ?", parentID)
	if err != nil {
		return
	}
	var childIDs []string
	for rows.Next() {
		var childID string
		if err := rows.Scan(&childID); err != nil {
			continue
		}
		childIDs = append(childIDs, childID)
	}
	rows.Close()
	for _, childID := range childIDs {
		*result = append(*result, childID)
		collectDraftChildren(childID, result)
	}
}

type ClipboardCapture struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

func DBGetClipboardCaptures(limit int) ([]ClipboardCapture, error) {
	rows, err := db.Query(`
		SELECT id, content, created_at
		FROM wt_clipboard ORDER BY created_at DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var captures []ClipboardCapture
	for rows.Next() {
		var c ClipboardCapture
		if err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt); err != nil {
			continue
		}
		captures = append(captures, c)
	}
	return captures, nil
}

func DBAddClipboardCapture(capture *ClipboardCapture) error {
	db.Exec("DELETE FROM wt_clipboard WHERE content = ?", capture.Content)
	_, err := db.Exec(`
		INSERT INTO wt_clipboard (id, content, created_at)
		VALUES (?, ?, ?)
	`, capture.ID, capture.Content, capture.CreatedAt)
	return err
}

func DBDeleteClipboardCapture(id string) error {
	_, err := db.Exec("DELETE FROM wt_clipboard WHERE id = ?", id)
	return err
}

func DBClearAllClipboardCaptures() error {
	_, err := db.Exec("DELETE FROM wt_clipboard")
	return err
}

func DBCleanupOldClipboardCaptures(keepCount int) error {
	_, err := db.Exec(`
		DELETE FROM wt_clipboard WHERE id NOT IN (
			SELECT id FROM wt_clipboard ORDER BY created_at DESC LIMIT ?
		)
	`, keepCount)
	return err
}


type GalleryImage struct {
	ID         string `json:"id"`
	Hash       string `json:"hash"`
	LocalPath  string `json:"localPath"`
	RemoteURL  string `json:"remoteUrl"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	CreatedAt  string `json:"createdAt"`
	FolderPath string `json:"folderPath"`
	IsValid    bool   `json:"isValid"`
}

type GalleryFolder struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	CreatedAt string `json:"createdAt"`
}

func DBGetAllGalleryImages() ([]GalleryImage, error) {
	rows, err := db.Query(`SELECT id, hash, local_path, remote_url, file_name, file_size, created_at, COALESCE(folder_path, '') FROM wt_image ORDER BY folder_path, created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var images []GalleryImage
	for rows.Next() {
		var img GalleryImage
		if err := rows.Scan(&img.ID, &img.Hash, &img.LocalPath, &img.RemoteURL, &img.FileName, &img.FileSize, &img.CreatedAt, &img.FolderPath); err != nil {
			continue
		}
		images = append(images, img)
	}
	return images, nil
}

func DBGetGalleryImageByHash(hash string) (*GalleryImage, error) {
	var img GalleryImage
	err := db.QueryRow(`SELECT id, hash, local_path, remote_url, file_name, file_size, created_at, COALESCE(folder_path, '') FROM wt_image WHERE hash = ?`, hash).Scan(&img.ID, &img.Hash, &img.LocalPath, &img.RemoteURL, &img.FileName, &img.FileSize, &img.CreatedAt, &img.FolderPath)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func DBAddGalleryImage(img *GalleryImage) error {
	_, err := db.Exec(`INSERT INTO wt_image (id, hash, local_path, remote_url, file_name, file_size, created_at, folder_path) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		img.ID, img.Hash, img.LocalPath, img.RemoteURL, img.FileName, img.FileSize, img.CreatedAt, img.FolderPath)
	return err
}

func DBUpdateGalleryImageFolder(id, folderPath string) error {
	_, err := db.Exec(`UPDATE wt_image SET folder_path = ? WHERE id = ?`, folderPath, id)
	return err
}

func DBUpdateGalleryImageURL(id, remoteURL string) error {
	_, err := db.Exec(`UPDATE wt_image SET remote_url = ? WHERE id = ?`, remoteURL, id)
	return err
}

func DBUpdateGalleryImageName(id, fileName string) error {
	_, err := db.Exec(`UPDATE wt_image SET file_name = ? WHERE id = ?`, fileName, id)
	return err
}

func DBRenameFolderImages(oldPath, newPath string) error {
	_, err := db.Exec(`UPDATE wt_image SET folder_path = ? WHERE folder_path = ?`, newPath, oldPath)
	return err
}

func DBDeleteGalleryImage(id string) error {
	_, err := db.Exec(`DELETE FROM wt_image WHERE id = ?`, id)
	return err
}
