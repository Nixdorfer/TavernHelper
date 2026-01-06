package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type DirtyItem struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	OriginalHash string    `json:"originalHash"`
	CurrentHash  string    `json:"currentHash"`
	IsDirty      bool      `json:"isDirty"`
	LastModified time.Time `json:"lastModified"`
	LastSaved    time.Time `json:"lastSaved"`
}

type DirtyManager struct {
	items map[string]*DirtyItem
	mu    sync.RWMutex
}

func NewDirtyManager() *DirtyManager {
	return &DirtyManager{
		items: make(map[string]*DirtyItem),
	}
}

func computeHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

func (dm *DirtyManager) SetOriginal(id, itemType string, content []byte) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	hash := computeHash(content)
	dm.items[id] = &DirtyItem{
		ID:           id,
		Type:         itemType,
		OriginalHash: hash,
		CurrentHash:  hash,
		IsDirty:      false,
		LastModified: time.Now(),
		LastSaved:    time.Now(),
	}
}

func (dm *DirtyManager) UpdateCurrent(id string, content []byte) bool {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	item, exists := dm.items[id]
	if !exists {
		return false
	}
	newHash := computeHash(content)
	if newHash != item.CurrentHash {
		item.CurrentHash = newHash
		item.IsDirty = newHash != item.OriginalHash
		item.LastModified = time.Now()
		return true
	}
	return false
}

func (dm *DirtyManager) MarkSaved(id string, content []byte) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	item, exists := dm.items[id]
	if !exists {
		return
	}
	hash := computeHash(content)
	item.OriginalHash = hash
	item.CurrentHash = hash
	item.IsDirty = false
	item.LastSaved = time.Now()
}

func (dm *DirtyManager) IsDirty(id string) bool {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	item, exists := dm.items[id]
	if !exists {
		return false
	}
	return item.IsDirty
}

func (dm *DirtyManager) GetItem(id string) *DirtyItem {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	item, exists := dm.items[id]
	if !exists {
		return nil
	}
	itemCopy := *item
	return &itemCopy
}

func (dm *DirtyManager) GetDirtyItems() []DirtyItem {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	var result []DirtyItem
	for _, item := range dm.items {
		if item.IsDirty {
			result = append(result, *item)
		}
	}
	return result
}

func (dm *DirtyManager) GetAllItems() []DirtyItem {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	var result []DirtyItem
	for _, item := range dm.items {
		result = append(result, *item)
	}
	return result
}

func (dm *DirtyManager) ClearItem(id string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	delete(dm.items, id)
}

func (dm *DirtyManager) HasAnyDirty() bool {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	for _, item := range dm.items {
		if item.IsDirty {
			return true
		}
	}
	return false
}

func (dm *DirtyManager) GetMinutesSinceLastSave(id string) int {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	item, exists := dm.items[id]
	if !exists {
		return 0
	}
	return int(time.Since(item.LastSaved).Minutes())
}

type DirtyState struct {
	ID                  string `json:"id"`
	IsDirty             bool   `json:"isDirty"`
	MinutesSinceLastSave int   `json:"minutesSinceLastSave"`
}

func (a *App) GetDirtyState(id string) *DirtyState {
	if a.dirtyManager == nil {
		return nil
	}
	item := a.dirtyManager.GetItem(id)
	if item == nil {
		return nil
	}
	return &DirtyState{
		ID:                  id,
		IsDirty:             item.IsDirty,
		MinutesSinceLastSave: a.dirtyManager.GetMinutesSinceLastSave(id),
	}
}

func (a *App) SetOriginalContent(id, itemType, content string) {
	if a.dirtyManager == nil {
		a.dirtyManager = NewDirtyManager()
	}
	a.dirtyManager.SetOriginal(id, itemType, []byte(content))
}

func (a *App) UpdateCurrentContent(id, content string) bool {
	if a.dirtyManager == nil {
		return false
	}
	return a.dirtyManager.UpdateCurrent(id, []byte(content))
}

func (a *App) MarkContentSaved(id, content string) {
	if a.dirtyManager == nil {
		return
	}
	a.dirtyManager.MarkSaved(id, []byte(content))
}

func (a *App) IsContentDirty(id string) bool {
	if a.dirtyManager == nil {
		return false
	}
	return a.dirtyManager.IsDirty(id)
}

func (a *App) HasAnyDirtyContent() bool {
	if a.dirtyManager == nil {
		return false
	}
	return a.dirtyManager.HasAnyDirty()
}

func (a *App) GetAllDirtyStates() []DirtyState {
	if a.dirtyManager == nil {
		return nil
	}
	items := a.dirtyManager.GetAllItems()
	var result []DirtyState
	for _, item := range items {
		result = append(result, DirtyState{
			ID:                  item.ID,
			IsDirty:             item.IsDirty,
			MinutesSinceLastSave: a.dirtyManager.GetMinutesSinceLastSave(item.ID),
		})
	}
	return result
}

func (a *App) ClearDirtyItem(id string) {
	if a.dirtyManager == nil {
		return
	}
	a.dirtyManager.ClearItem(id)
}
