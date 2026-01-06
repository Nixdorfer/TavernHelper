package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

type WorldBookPatch struct {
	LineNo  int            `json:"ln"`
	EntryID string         `json:"id"`
	Op      string         `json:"op"`
	Changes map[string]any `json:"changes,omitempty"`
}

type PromptPatch struct {
	LineNo  int            `json:"ln"`
	ID      string         `json:"id"`
	Op      string         `json:"op"`
	Changes map[string]any `json:"changes,omitempty"`
}

type NodeDelta struct {
	NodeID            string           `json:"nodeId"`
	ParentID          string           `json:"parentId"`
	WorldBookAdded    []WorldBookEntry `json:"wb_added,omitempty"`
	WorldBookRemoved  []string         `json:"wb_removed,omitempty"`
	WorldBookModified []WorldBookPatch `json:"wb_modified,omitempty"`
	PrePromptDelta    []PromptPatch    `json:"pre_prompt_delta,omitempty"`
	PreTextDelta      []PromptPatch    `json:"pre_text_delta,omitempty"`
	PostTextDelta     []PromptPatch    `json:"post_text_delta,omitempty"`
	NameChanged       *string          `json:"name,omitempty"`
	NoteChanged       *string          `json:"note,omitempty"`
	TagsChanged       []string         `json:"tags,omitempty"`
}

type DeltaResult struct {
	IsRoot           bool          `json:"isRoot"`
	Delta            *NodeDelta    `json:"delta,omitempty"`
	FullNode         *TimelineNode `json:"fullNode,omitempty"`
	CompressionRatio float64       `json:"compressionRatio"`
}

func computeContentHash(data any) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:])
}

func (a *App) ComputeNodeDelta(parentNode, childNode *TimelineNode) (*NodeDelta, error) {
	if parentNode == nil || childNode == nil {
		return nil, fmt.Errorf("nodes cannot be nil")
	}
	delta := &NodeDelta{
		NodeID:   childNode.ID,
		ParentID: parentNode.ID,
	}
	if childNode.Name != parentNode.Name {
		name := childNode.Name
		delta.NameChanged = &name
	}
	if childNode.Note != parentNode.Note {
		note := childNode.Note
		delta.NoteChanged = &note
	}
	if !stringSliceEqual(childNode.Tags, parentNode.Tags) {
		delta.TagsChanged = childNode.Tags
	}
	parentWBMap := make(map[string]WorldBookEntry)
	for _, wb := range parentNode.WorldBook {
		parentWBMap[wb.ID] = wb
	}
	childWBMap := make(map[string]WorldBookEntry)
	for _, wb := range childNode.WorldBook {
		childWBMap[wb.ID] = wb
	}
	for id := range parentWBMap {
		if _, exists := childWBMap[id]; !exists {
			delta.WorldBookRemoved = append(delta.WorldBookRemoved, id)
		}
	}
	for id, childWB := range childWBMap {
		parentWB, exists := parentWBMap[id]
		if !exists {
			delta.WorldBookAdded = append(delta.WorldBookAdded, childWB)
		} else {
			if computeContentHash(parentWB) != computeContentHash(childWB) {
				changes := make(map[string]any)
				if parentWB.Key != childWB.Key {
					changes["key"] = childWB.Key
				}
				if parentWB.Value != childWB.Value {
					changes["value"] = childWB.Value
				}
				if parentWB.Name != childWB.Name {
					changes["name"] = childWB.Name
				}
				if computeContentHash(parentWB.ContentItems) != computeContentHash(childWB.ContentItems) {
					changes["contentItems"] = childWB.ContentItems
				}
				if computeContentHash(parentWB.Keywords) != computeContentHash(childWB.Keywords) {
					changes["keywords"] = childWB.Keywords
				}
				if len(changes) > 0 {
					patch := WorldBookPatch{
						EntryID: id,
						Op:      "modify",
						Changes: changes,
					}
					delta.WorldBookModified = append(delta.WorldBookModified, patch)
				}
			}
		}
	}
	delta.PrePromptDelta = computePromptDeltaTyped(parentNode.PrePrompt, childNode.PrePrompt)
	delta.PreTextDelta = computePromptDeltaTyped(parentNode.PreText, childNode.PreText)
	delta.PostTextDelta = computePromptDeltaTyped(parentNode.PostText, childNode.PostText)
	return delta, nil
}

func computePromptDeltaTyped(parentPrompts, childPrompts []PromptEntry) []PromptPatch {
	var patches []PromptPatch
	parentMap := make(map[string]PromptEntry)
	for _, p := range parentPrompts {
		parentMap[p.ID] = p
	}
	childMap := make(map[string]PromptEntry)
	for _, p := range childPrompts {
		childMap[p.ID] = p
	}
	for id := range parentMap {
		if _, exists := childMap[id]; !exists {
			patches = append(patches, PromptPatch{
				ID: id,
				Op: "remove",
			})
		}
	}
	for id, childP := range childMap {
		parentP, exists := parentMap[id]
		if !exists {
			patches = append(patches, PromptPatch{
				ID: id,
				Op: "add",
				Changes: map[string]any{
					"id":      childP.ID,
					"name":    childP.Name,
					"content": childP.Content,
					"enabled": childP.Enabled,
				},
			})
		} else {
			if computeContentHash(parentP) != computeContentHash(childP) {
				changes := make(map[string]any)
				if parentP.Name != childP.Name {
					changes["name"] = childP.Name
				}
				if parentP.Content != childP.Content {
					changes["content"] = childP.Content
				}
				if parentP.Enabled != childP.Enabled {
					changes["enabled"] = childP.Enabled
				}
				if len(changes) > 0 {
					patches = append(patches, PromptPatch{
						ID:      id,
						Op:      "modify",
						Changes: changes,
					})
				}
			}
		}
	}
	return patches
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

func (a *App) ApplyDelta(baseNode *TimelineNode, delta *NodeDelta) (*TimelineNode, error) {
	if baseNode == nil || delta == nil {
		return nil, fmt.Errorf("base node and delta cannot be nil")
	}
	result := &TimelineNode{
		ID:        delta.NodeID,
		Name:      baseNode.Name,
		Note:      baseNode.Note,
		Tags:      append([]string{}, baseNode.Tags...),
		WorldBook: make([]WorldBookEntry, 0),
		PrePrompt: make([]PromptEntry, 0),
		PreText:   make([]PromptEntry, 0),
		PostText:  make([]PromptEntry, 0),
		ParentID:  delta.ParentID,
	}
	if delta.NameChanged != nil {
		result.Name = *delta.NameChanged
	}
	if delta.NoteChanged != nil {
		result.Note = *delta.NoteChanged
	}
	if delta.TagsChanged != nil {
		result.Tags = delta.TagsChanged
	}
	wbMap := make(map[string]WorldBookEntry)
	for _, wb := range baseNode.WorldBook {
		wbMap[wb.ID] = wb
	}
	for _, id := range delta.WorldBookRemoved {
		delete(wbMap, id)
	}
	for _, added := range delta.WorldBookAdded {
		wbMap[added.ID] = added
	}
	for _, patch := range delta.WorldBookModified {
		if existing, ok := wbMap[patch.EntryID]; ok {
			if v, ok := patch.Changes["key"].(string); ok {
				existing.Key = v
			}
			if v, ok := patch.Changes["value"].(string); ok {
				existing.Value = v
			}
			if v, ok := patch.Changes["name"].(string); ok {
				existing.Name = v
			}
			if v, ok := patch.Changes["contentItems"].([]ContentItem); ok {
				existing.ContentItems = v
			}
			if v, ok := patch.Changes["keywords"].([]string); ok {
				existing.Keywords = v
			}
			wbMap[patch.EntryID] = existing
		}
	}
	for _, wb := range wbMap {
		result.WorldBook = append(result.WorldBook, wb)
	}
	result.PrePrompt = applyPromptDeltaTyped(baseNode.PrePrompt, delta.PrePromptDelta)
	result.PreText = applyPromptDeltaTyped(baseNode.PreText, delta.PreTextDelta)
	result.PostText = applyPromptDeltaTyped(baseNode.PostText, delta.PostTextDelta)
	return result, nil
}

func applyPromptDeltaTyped(basePrompts []PromptEntry, patches []PromptPatch) []PromptEntry {
	promptMap := make(map[string]PromptEntry)
	for _, p := range basePrompts {
		promptMap[p.ID] = p
	}
	for _, patch := range patches {
		switch patch.Op {
		case "remove":
			delete(promptMap, patch.ID)
		case "add":
			entry := PromptEntry{ID: patch.ID}
			if v, ok := patch.Changes["name"].(string); ok {
				entry.Name = v
			}
			if v, ok := patch.Changes["content"].(string); ok {
				entry.Content = v
			}
			if v, ok := patch.Changes["enabled"].(bool); ok {
				entry.Enabled = v
			}
			promptMap[patch.ID] = entry
		case "modify":
			if existing, ok := promptMap[patch.ID]; ok {
				if v, ok := patch.Changes["name"].(string); ok {
					existing.Name = v
				}
				if v, ok := patch.Changes["content"].(string); ok {
					existing.Content = v
				}
				if v, ok := patch.Changes["enabled"].(bool); ok {
					existing.Enabled = v
				}
				promptMap[patch.ID] = existing
			}
		}
	}
	result := make([]PromptEntry, 0, len(promptMap))
	for _, p := range promptMap {
		result = append(result, p)
	}
	return result
}

func (a *App) CalculateDeltaCompressionRatio(fullNode *TimelineNode, delta *NodeDelta) float64 {
	fullBytes, _ := json.Marshal(fullNode)
	deltaBytes, _ := json.Marshal(delta)
	if len(fullBytes) == 0 {
		return 1.0
	}
	return float64(len(deltaBytes)) / float64(len(fullBytes))
}

func (a *App) GetNodeWithDelta(projectId, nodeId string, timeline []TimelineNode) (*DeltaResult, error) {
	var targetNode *TimelineNode
	var parentNode *TimelineNode
	for i := range timeline {
		if timeline[i].ID == nodeId {
			targetNode = &timeline[i]
			break
		}
	}
	if targetNode == nil {
		return nil, fmt.Errorf("node not found: %s", nodeId)
	}
	if targetNode.ParentID == "" {
		return &DeltaResult{
			IsRoot:           true,
			FullNode:         targetNode,
			CompressionRatio: 1.0,
		}, nil
	}
	for i := range timeline {
		if timeline[i].ID == targetNode.ParentID {
			parentNode = &timeline[i]
			break
		}
	}
	if parentNode == nil {
		return &DeltaResult{
			IsRoot:           true,
			FullNode:         targetNode,
			CompressionRatio: 1.0,
		}, nil
	}
	delta, err := a.ComputeNodeDelta(parentNode, targetNode)
	if err != nil {
		return nil, err
	}
	ratio := a.CalculateDeltaCompressionRatio(targetNode, delta)
	return &DeltaResult{
		IsRoot:           false,
		Delta:            delta,
		CompressionRatio: ratio,
	}, nil
}

func (a *App) ReconstructNodeFromDelta(projectId, nodeId string, timeline []TimelineNode) (*TimelineNode, error) {
	nodeMap := make(map[string]*TimelineNode)
	for i := range timeline {
		nodeMap[timeline[i].ID] = &timeline[i]
	}
	target, exists := nodeMap[nodeId]
	if !exists {
		return nil, fmt.Errorf("node not found: %s", nodeId)
	}
	var path []*TimelineNode
	current := target
	for current != nil {
		path = append([]*TimelineNode{current}, path...)
		if current.ParentID == "" {
			break
		}
		parent, exists := nodeMap[current.ParentID]
		if !exists {
			break
		}
		current = parent
	}
	if len(path) == 0 {
		return nil, fmt.Errorf("could not build path to node")
	}
	result := path[0]
	for i := 1; i < len(path); i++ {
		delta, err := a.ComputeNodeDelta(result, path[i])
		if err != nil {
			return nil, fmt.Errorf("failed to compute delta: %w", err)
		}
		result, err = a.ApplyDelta(result, delta)
		if err != nil {
			return nil, fmt.Errorf("failed to apply delta: %w", err)
		}
	}
	return result, nil
}

type WorldBookWithLine struct {
	WorldBookEntry
	LineNo int `json:"ln"`
}

func (a *App) AddLineNumbersToWorldBook(worldBook []WorldBookEntry) []WorldBookWithLine {
	result := make([]WorldBookWithLine, len(worldBook))
	for i, wb := range worldBook {
		result[i] = WorldBookWithLine{
			WorldBookEntry: wb,
			LineNo:         i + 1,
		}
	}
	return result
}
