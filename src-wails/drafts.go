package main

import (
	"fmt"
	"time"
)

func (a *App) GetAllDrafts() ([]Draft, error) {
	return DBGetAllDrafts()
}

func (a *App) CreateDraft(draft Draft) error {
	return DBCreateDraft(&draft)
}

func (a *App) UpdateDraft(draft Draft) error {
	return DBUpdateDraft(&draft)
}

func (a *App) DeleteDraft(id string) error {
	return DBDeleteDraft(id)
}

func (a *App) DeleteDraftRecursive(id string) error {
	return DBDeleteDraftRecursive(id)
}

func (a *App) GetClipboardCaptures() ([]ClipboardCapture, error) {
	return DBGetClipboardCaptures(50)
}

func (a *App) DeleteClipboardCapture(id string) error {
	return DBDeleteClipboardCapture(id)
}

func (a *App) ClearAllClipboardCaptures() error {
	return DBClearAllClipboardCaptures()
}

func (a *App) MoveClipboardToDraft(captureID, name, parentID string) (*Draft, error) {
	captures, err := DBGetClipboardCaptures(50)
	if err != nil {
		return nil, err
	}

	var capture *ClipboardCapture
	for _, c := range captures {
		if c.ID == captureID {
			capture = &c
			break
		}
	}
	if capture == nil {
		return nil, fmt.Errorf("未找到剪贴板条目")
	}

	now := time.Now().Format(time.RFC3339)
	draft := &Draft{
		ID:        fmt.Sprintf("draft_%d", time.Now().UnixNano()),
		Name:      name,
		Content:   capture.Content,
		ParentID:  parentID,
		IsFolder:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := DBCreateDraft(draft); err != nil {
		return nil, err
	}

	DBDeleteClipboardCapture(captureID)

	return draft, nil
}
