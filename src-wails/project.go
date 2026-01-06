package main

import (
	"database/sql"
	"fmt"
	"time"
)

func (a *App) GetProjects() ([]ProjectInfo, error) {
	rows, err := db.Query(`SELECT id, name, type, update_time FROM wt_project ORDER BY update_time DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []ProjectInfo
	for rows.Next() {
		var p ProjectInfo
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.UpdatedAt); err != nil {
			continue
		}
		p.FileName = p.Name
		projects = append(projects, p)
	}
	return projects, nil
}

func (a *App) LoadProjectByName(name string) (*Project, error) {
	fmt.Println("[LoadProjectByName] 开始加载:", name)
	var projectID int
	err := db.QueryRow(`SELECT id FROM wt_project WHERE name = ?`, name).Scan(&projectID)
	if err != nil {
		fmt.Println("[LoadProjectByName] 查询项目ID失败:", err)
		return nil, err
	}
	fmt.Println("[LoadProjectByName] 找到项目ID:", projectID)
	return a.LoadProject(projectID)
}

func (a *App) LoadProject(projectID int) (*Project, error) {
	fmt.Println("[LoadProject] 开始加载项目:", projectID)
	var proj Project
	err := db.QueryRow(`
		SELECT id, name, type, current_node, create_time, update_time FROM wt_project WHERE id = ?
	`, projectID).Scan(&proj.ID, &proj.Name, &proj.Type, &proj.CurrentNode, &proj.CreateTime, &proj.UpdateTime)
	if err != nil {
		fmt.Println("[LoadProject] 查询项目失败:", err)
		return nil, err
	}
	fmt.Println("[LoadProject] 项目基本信息加载完成:", proj.Name)
	proj.FileName = proj.Name
	nodes, err := a.loadProjectNodes(projectID)
	if err != nil {
		fmt.Println("[LoadProject] 加载节点失败:", err)
		return nil, err
	}
	fmt.Println("[LoadProject] 节点加载完成, 数量:", len(nodes))
	proj.Timeline = nodes
	if proj.CurrentNode == nil && len(nodes) > 0 {
		rootNode := nodes[0].ID
		proj.CurrentNode = &rootNode
	}
	return &proj, nil
}

func (a *App) loadProjectNodes(projectID int) ([]NodeInfo, error) {
	rows, err := db.Query(`
		SELECT id, parent_id, name, COALESCE(desc, '')
		FROM wt_node WHERE project_id = ? ORDER BY id
	`, projectID)
	if err != nil {
		return nil, err
	}
	var nodes []NodeInfo
	for rows.Next() {
		var node NodeInfo
		var parentID sql.NullInt64
		var desc string
		if err := rows.Scan(&node.ID, &parentID, &node.Name, &desc); err != nil {
			continue
		}
		if parentID.Valid {
			pid := int(parentID.Int64)
			node.ParentID = &pid
		}
		node.Note = desc
		node.Tags = []string{}
		nodes = append(nodes, node)
	}
	rows.Close()
	for i := range nodes {
		content, err := a.GetNodeContent(nodes[i].ID)
		if err == nil && content != nil {
			nodes[i].PreText = a.convertBlocksToPromptEntries(content.Pre)
			nodes[i].PostText = a.convertBlocksToPromptEntries(content.Post)
			nodes[i].PrePrompt = a.convertBlocksToPromptEntries(content.Global)
			nodes[i].WorldBook = a.convertFoldersToWorldBook(content.Folders)
		} else {
			nodes[i].PreText = []PromptEntry{}
			nodes[i].PostText = []PromptEntry{}
			nodes[i].PrePrompt = []PromptEntry{}
			nodes[i].WorldBook = []WorldBookEntry{}
		}
	}
	return nodes, nil
}
func (a *App) convertBlocksToPromptEntries(blocks []ReplayedBlock) []PromptEntry {
	var entries []PromptEntry
	for _, block := range blocks {
		var content string
		for i, line := range block.Lines {
			if i > 0 {
				content += "\n"
			}
			if line.Content != nil {
				content += *line.Content
			}
		}
		entries = append(entries, PromptEntry{
			ID:      fmt.Sprintf("%d", block.ID),
			Name:    block.Title,
			Content: content,
			Enabled: true,
		})
	}
	return entries
}
func (a *App) convertFoldersToWorldBook(folders []ReplayedFolder) []WorldBookEntry {
	var entries []WorldBookEntry
	for _, folder := range folders {
		folderID := fmt.Sprintf("folder_%d", folder.ID)
		entries = append(entries, WorldBookEntry{
			ID:       folderID,
			Name:     folder.Name,
			IsFolder: true,
		})
		for _, card := range folder.Cards {
			cardID := fmt.Sprintf("card_%d", card.ID)
			var contentItems []ContentItem
			for _, block := range card.Blocks {
				var content string
				for i, line := range block.Lines {
					if i > 0 {
						content += "\n"
					}
					if line.Content != nil {
						content += *line.Content
					}
				}
				contentItems = append(contentItems, ContentItem{
					ID:      fmt.Sprintf("block_%d", block.ID),
					Title:   block.Title,
					Content: content,
				})
			}
			entries = append(entries, WorldBookEntry{
				ID:           cardID,
				Name:         card.Name,
				ParentID:     folderID,
				IsFolder:     false,
				Key:          card.KeyWord,
				ContentItems: contentItems,
			})
		}
	}
	return entries
}

func (a *App) CreateProject(name string, projectType string) (int, error) {
	now := time.Now().Format(time.RFC3339)
	if projectType != "play" && projectType != "create" {
		projectType = "create"
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.Exec(`
		INSERT INTO wt_project (name, type, create_time, update_time) VALUES (?, ?, ?, ?)
	`, name, projectType, now, now)
	if err != nil {
		return 0, fmt.Errorf("创建项目失败: %w", err)
	}
	projectID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	result, err = tx.Exec(`
		INSERT INTO wt_node (project_id, parent_id, name, desc) VALUES (?, NULL, ?, '')
	`, projectID, "根节点")
	if err != nil {
		return 0, fmt.Errorf("创建根节点失败: %w", err)
	}
	rootNodeID, _ := result.LastInsertId()
	sysFolders := []struct {
		folderName string
		cardName   string
	}{
		{"SYS_PRE", "SYS_PRE_BASE"},
		{"SYS_POST", "SYS_POST_BASE"},
		{"SYS_GLOBAL", "SYS_GLOBAL_BASE"},
	}
	for _, sf := range sysFolders {
		result, err = tx.Exec("INSERT INTO wt_folder (name) VALUES (?)", sf.folderName)
		if err != nil {
			return 0, fmt.Errorf("创建系统文件夹失败: %w", err)
		}
		folderID, _ := result.LastInsertId()
		result, err = tx.Exec(`
			INSERT INTO wt_node_change (action, level, node_id, target, detail_folder)
			VALUES ('add', 'folder', ?, NULL, ?)
		`, rootNodeID, folderID)
		if err != nil {
			return 0, fmt.Errorf("创建系统文件夹变更失败: %w", err)
		}
		folderChangeID, _ := result.LastInsertId()
		result, err = tx.Exec(`
			INSERT INTO wt_card (name, desc, key_word, trigger_system, trigger_user, trigger_ai)
			VALUES (?, '', '', 1, 1, 1)
		`, sf.cardName)
		if err != nil {
			return 0, fmt.Errorf("创建系统卡片失败: %w", err)
		}
		cardID, _ := result.LastInsertId()
		_, err = tx.Exec(`
			INSERT INTO wt_node_change (action, level, node_id, target, detail_card)
			VALUES ('add', 'card', ?, ?, ?)
		`, rootNodeID, folderChangeID, cardID)
		if err != nil {
			return 0, fmt.Errorf("创建系统卡片变更失败: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return int(projectID), nil
}

func (a *App) DeleteProject(projectName string) error {
	var projectID int
	err := db.QueryRow("SELECT id FROM wt_project WHERE name = ?", projectName).Scan(&projectID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE wt_project SET current_node = NULL WHERE id = ?", projectID)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		DELETE FROM wt_node_change WHERE node_id IN (
			SELECT id FROM wt_node WHERE project_id = ?
		)
	`, projectID)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM wt_line WHERE project_id = ?", projectID)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		DELETE FROM wt_block WHERE id NOT IN (
			SELECT DISTINCT detail_block FROM wt_node_change WHERE detail_block IS NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		DELETE FROM wt_card WHERE id NOT IN (
			SELECT DISTINCT detail_card FROM wt_node_change WHERE detail_card IS NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		DELETE FROM wt_folder WHERE id NOT IN (
			SELECT DISTINCT detail_folder FROM wt_node_change WHERE detail_folder IS NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM wt_project WHERE id = ?", projectID)
	return err
}

func (a *App) RenameProject(oldName string, newName string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec("UPDATE wt_project SET name = ?, update_time = ? WHERE name = ?", newName, now, oldName)
	return err
}

func (a *App) SaveProject(projectName string, projectData any) error {
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec("UPDATE wt_project SET update_time = ? WHERE name = ?", now, projectName)
	return err
}

func (a *App) UpdateProjectTime(projectID int) error {
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec("UPDATE wt_project SET update_time = ? WHERE id = ?", now, projectID)
	return err
}
func (a *App) UpdateProjectCurrentNode(projectName string, nodeID *int) error {
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec("UPDATE wt_project SET current_node = ?, update_time = ? WHERE name = ?", nodeID, now, projectName)
	return err
}

func (a *App) GenerateEncryptionKeyPair() (map[string]string, error) {
	pubKey, privKey, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"publicKey":  pubKey,
		"privateKey": privKey,
	}, nil
}
