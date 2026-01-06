package main

import (
	"database/sql"
	"fmt"
)

func (a *App) CreateChildNode(projectName string, parentNodeID any, name string) (*NodeInfo, error) {
	var projectID int
	err := db.QueryRow("SELECT id FROM wt_project WHERE name = ?", projectName).Scan(&projectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在")
	}
	var parentID int
	switch v := parentNodeID.(type) {
	case float64:
		parentID = int(v)
	case int:
		parentID = v
	}
	var parentIDPtr *int
	if parentID > 0 {
		var parentName string
		err = db.QueryRow("SELECT name FROM wt_node WHERE id = ?", parentID).Scan(&parentName)
		if err != nil {
			return nil, fmt.Errorf("父节点不存在")
		}
		parentIDPtr = &parentID
	}
	result, err := db.Exec(`
		INSERT INTO wt_node (project_id, parent_id, name, desc) VALUES (?, ?, ?, '')
	`, projectID, parentIDPtr, name)
	if err != nil {
		return nil, err
	}
	nodeID, _ := result.LastInsertId()
	return &NodeInfo{
		ID:       int(nodeID),
		ParentID: parentIDPtr,
		Name:     name,
	}, nil
}

func (a *App) CreateBrotherNode(projectName string, siblingNodeID any, name string) (*NodeInfo, error) {
	var projectID int
	err := db.QueryRow("SELECT id FROM wt_project WHERE name = ?", projectName).Scan(&projectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在")
	}
	var siblingID int
	switch v := siblingNodeID.(type) {
	case float64:
		siblingID = int(v)
	case int:
		siblingID = v
	}
	var parentID sql.NullInt64
	err = db.QueryRow("SELECT parent_id FROM wt_node WHERE id = ?", siblingID).Scan(&parentID)
	if err != nil {
		return nil, fmt.Errorf("节点不存在")
	}
	var parentIDVal any
	if parentID.Valid {
		parentIDVal = parentID.Int64
	}
	result, err := db.Exec(`
		INSERT INTO wt_node (project_id, parent_id, name, desc) VALUES (?, ?, ?, '')
	`, projectID, parentIDVal, name)
	if err != nil {
		return nil, err
	}
	nodeID, _ := result.LastInsertId()
	var parentIDPtr *int
	if parentID.Valid {
		pid := int(parentID.Int64)
		parentIDPtr = &pid
	}
	return &NodeInfo{
		ID:       int(nodeID),
		ParentID: parentIDPtr,
		Name:     name,
	}, nil
}

func (a *App) UpdateNodeByID(nodeID int, name string, desc string) error {
	_, err := db.Exec(`UPDATE wt_node SET name = ?, desc = ? WHERE id = ?`, name, desc, nodeID)
	return err
}

func (a *App) UpdateNode(projectName string, nodeData map[string]any) error {
	nodeID, ok := nodeData["id"]
	if !ok {
		return fmt.Errorf("节点数据缺少id")
	}
	var id int
	switch v := nodeID.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	}
	name, _ := nodeData["name"].(string)
	desc, _ := nodeData["note"].(string)
	_, err := db.Exec(`UPDATE wt_node SET name = ?, desc = ? WHERE id = ?`, name, desc, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteNode(projectName string, nodeID any) error {
	var id int
	switch v := nodeID.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	case string:
		return fmt.Errorf("invalid node ID type")
	}
	var projectID int
	err := db.QueryRow("SELECT project_id FROM wt_node WHERE id = ?", id).Scan(&projectID)
	if err != nil {
		return err
	}
	var count int
	db.QueryRow("SELECT COUNT(*) FROM wt_node WHERE project_id = ?", projectID).Scan(&count)
	if count <= 1 {
		return fmt.Errorf("不能删除所有节点")
	}
	_, err = db.Exec("DELETE FROM wt_node WHERE id = ?", id)
	return err
}

func (a *App) UpdateNodeNote(projectName string, nodeID any, note string) error {
	var id int
	switch v := nodeID.(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	}
	_, err := db.Exec("UPDATE wt_node SET desc = ? WHERE id = ?", note, id)
	return err
}

func (a *App) RenameNode(nodeID int, newName string) error {
	_, err := db.Exec("UPDATE wt_node SET name = ? WHERE id = ?", newName, nodeID)
	return err
}

func (a *App) UpdateNodeDesc(nodeID int, desc string) error {
	_, err := db.Exec("UPDATE wt_node SET desc = ? WHERE id = ?", desc, nodeID)
	return err
}

func (a *App) RebaseNode(nodeID int, newParentID *int) error {
	if newParentID == nil {
		_, err := db.Exec("UPDATE wt_node SET parent_id = NULL WHERE id = ?", nodeID)
		return err
	}
	current := *newParentID
	for {
		if current == nodeID {
			return fmt.Errorf("不能形成循环引用")
		}
		var parentID sql.NullInt64
		err := db.QueryRow("SELECT parent_id FROM wt_node WHERE id = ?", current).Scan(&parentID)
		if err != nil || !parentID.Valid {
			break
		}
		current = int(parentID.Int64)
	}
	_, err := db.Exec("UPDATE wt_node SET parent_id = ? WHERE id = ?", *newParentID, nodeID)
	return err
}

func (a *App) GetBranchTag(parentID int, childID int) (*WTBranchTag, error) {
	var tag WTBranchTag
	err := db.QueryRow(`
		SELECT id, parent_id, child_id, name, COALESCE(desc, '')
		FROM wt_branch_tag WHERE parent_id = ? AND child_id = ?
	`, parentID, childID).Scan(&tag.ID, &tag.ParentID, &tag.ChildID, &tag.Name, &tag.Desc)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (a *App) SetBranchTag(parentID int, childID int, name string, desc string) error {
	_, err := db.Exec(`
		INSERT INTO wt_branch_tag (parent_id, child_id, name, desc) VALUES (?, ?, ?, ?)
		ON CONFLICT(parent_id, child_id) DO UPDATE SET name = excluded.name, desc = excluded.desc
	`, parentID, childID, name, desc)
	return err
}

func (a *App) DeleteBranchTag(parentID int, childID int) error {
	_, err := db.Exec("DELETE FROM wt_branch_tag WHERE parent_id = ? AND child_id = ?", parentID, childID)
	return err
}

func (a *App) GetNodePath(nodeID int) ([]int, error) {
	fmt.Println("[GetNodePath] 开始, nodeID:", nodeID)
	var path []int
	current := nodeID
	for {
		fmt.Println("[GetNodePath] 处理current:", current)
		path = append([]int{current}, path...)
		var parentID sql.NullInt64
		fmt.Println("[GetNodePath] 查询parent_id...")
		err := db.QueryRow("SELECT parent_id FROM wt_node WHERE id = ?", current).Scan(&parentID)
		fmt.Println("[GetNodePath] 查询完成, err:", err, "parentID.Valid:", parentID.Valid)
		if err != nil || !parentID.Valid {
			break
		}
		current = int(parentID.Int64)
	}
	fmt.Println("[GetNodePath] 完成, path:", path)
	return path, nil
}

func (a *App) GenerateLineSerial(projectName string) (string, error) {
	var projectID int
	err := db.QueryRow("SELECT id FROM wt_project WHERE name = ?", projectName).Scan(&projectID)
	if err != nil {
		projectID = 0
	}
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	sn := generateSerial(tx, projectID)
	return sn, nil
}
