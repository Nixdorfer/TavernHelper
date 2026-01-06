package main

type LineInfo struct {
	Serial  string
	Content string
}

type SyncDotStatus struct {
	LineIndex int    `json:"lineIndex"`
	Serial    string `json:"serial"`
	Shape     string `json:"shape"`
	Color     string `json:"color"`
}

func (a *App) ComputeEffectiveContent(projectID int, nodeId int, entryId string, contentItemSerial string) (string, error) {
	return "", nil
}

func (a *App) ComputeSyncDotStatus(projectID int, nodeId int, entryId string, contentItemSerial string) ([]SyncDotStatus, error) {
	return nil, nil
}

func (a *App) AddLineMarker(projectID int, nodeId int, entryId string, marker LineMarker) error {
	return nil
}

func (a *App) RemoveLineMarker(projectID int, nodeId int, entryId string, serial string, markerType string) error {
	return nil
}

func (a *App) DeleteLineInNode(projectID int, nodeId int, entryId string, lineSerial string, contentItemSerial string) error {
	return nil
}

func (a *App) AddLineInNode(projectID int, nodeId int, entryId string, content string, insertBeforeSerial string, newSerial string, contentItemSerial string) error {
	return nil
}

func (a *App) GetNodeLineMarkers(projectID int, nodeId int, entryId string) ([]LineMarker, error) {
	return nil, nil
}

func (a *App) ComputeNodeInheritedContentItems(projectID int, nodeId int, entryId string) ([]ContentItem, error) {
	return []ContentItem{}, nil
}

func (a *App) ComputeNodeEffectiveContentItems(projectID int, nodeId int, entryId string) ([]ContentItem, error) {
	return []ContentItem{}, nil
}
