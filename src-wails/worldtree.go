package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

var serialChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateSerial(tx *sql.Tx, projectID int) string {
	var maxSN sql.NullString
	tx.QueryRow("SELECT sn FROM wt_line WHERE project_id = ? ORDER BY LENGTH(sn) DESC, sn DESC LIMIT 1", projectID).Scan(&maxSN)
	if !maxSN.Valid || maxSN.String == "" {
		return "aaaa"
	}
	return nextSerial(maxSN.String)
}

func nextSerial(sn string) string {
	base := len(serialChars)
	digits := []byte(sn)
	for i := len(digits) - 1; i >= 0; i-- {
		idx := strings.IndexByte(serialChars, digits[i])
		if idx < base-1 {
			digits[i] = serialChars[idx+1]
			return string(digits)
		}
		digits[i] = serialChars[0]
	}
	if len(digits) >= 8 {
		return string(digits)
	}
	return strings.Repeat("a", len(digits)+1)
}

func (a *App) GetNodeContent(nodeID int) (*NodeContent, error) {
	fmt.Println("[GetNodeContent] 开始, nodeID:", nodeID)
	path, err := a.GetNodePath(nodeID)
	if err != nil {
		fmt.Println("[GetNodeContent] GetNodePath失败:", err)
		return nil, err
	}
	fmt.Println("[GetNodeContent] path:", path)
	if len(path) == 0 {
		return nil, fmt.Errorf("节点不存在")
	}
	fmt.Println("[GetNodeContent] 开始replayFolders...")
	folders, err := a.replayFolders(path)
	if err != nil {
		fmt.Println("[GetNodeContent] replayFolders失败:", err)
		return nil, err
	}
	fmt.Println("[GetNodeContent] folders完成, 数量:", len(folders))
	pre, _ := a.replayBlocksByZone(path, "pre")
	fmt.Println("[GetNodeContent] pre完成")
	post, _ := a.replayBlocksByZone(path, "post")
	fmt.Println("[GetNodeContent] post完成")
	global, _ := a.replayBlocksByZone(path, "global")
	fmt.Println("[GetNodeContent] 完成")
	return &NodeContent{
		NodeID:  nodeID,
		Folders: folders,
		Pre:     pre,
		Post:    post,
		Global:  global,
	}, nil
}

func (a *App) replayFolders(path []int) ([]ReplayedFolder, error) {
	fmt.Println("[replayFolders] 开始, path:", path)
	type folderState struct {
		folder   WTFolder
		changeID int
		deleted  bool
	}
	type delRecord struct {
		changeID int
	}
	folderMap := make(map[int]*folderState)
	var delRecords []delRecord
	for _, nodeID := range path {
		fmt.Println("[replayFolders] 查询nodeID:", nodeID)
		rows, err := db.Query(`
			SELECT nc.id, nc.action, f.id, f.name
			FROM wt_node_change nc
			LEFT JOIN wt_folder f ON nc.detail_folder = f.id
			WHERE nc.node_id = ? AND nc.level = 'folder'
			ORDER BY nc.id
		`, nodeID)
		if err != nil {
			fmt.Println("[replayFolders] 查询失败:", err)
			continue
		}
		for rows.Next() {
			var changeID int
			var action string
			var folderID sql.NullInt64
			var folderName sql.NullString
			rows.Scan(&changeID, &action, &folderID, &folderName)
			if action == "add" && folderID.Valid {
				folderMap[int(folderID.Int64)] = &folderState{
					folder:   WTFolder{ID: int(folderID.Int64), Name: folderName.String},
					changeID: changeID,
					deleted:  false,
				}
			} else if action == "del" {
				delRecords = append(delRecords, delRecord{changeID: changeID})
			}
		}
		rows.Close()
	}
	for _, del := range delRecords {
		var targetChangeID int
		db.QueryRow("SELECT target FROM wt_node_change WHERE id = ?", del.changeID).Scan(&targetChangeID)
		var targetFolderID int
		db.QueryRow("SELECT detail_folder FROM wt_node_change WHERE id = ?", targetChangeID).Scan(&targetFolderID)
		if state, ok := folderMap[targetFolderID]; ok {
			state.deleted = true
		}
	}
	fmt.Println("[replayFolders] folderMap数量:", len(folderMap))
	var result []ReplayedFolder
	for _, state := range folderMap {
		if !state.deleted {
			cards, _ := a.replayCardsInFolder(path, state.changeID)
			result = append(result, ReplayedFolder{
				ID:       state.folder.ID,
				Name:     state.folder.Name,
				ChangeID: state.changeID,
				Cards:    cards,
			})
		}
	}
	return result, nil
}

func (a *App) replayCardsInFolder(path []int, folderChangeID int) ([]ReplayedCard, error) {
	type cardState struct {
		card     WTCard
		changeID int
		deleted  bool
	}
	type delRecord struct {
		changeID int
	}
	cardMap := make(map[int]*cardState)
	var delRecords []delRecord
	for _, nodeID := range path {
		rows, err := db.Query(`
			SELECT nc.id, nc.action, nc.target, c.id, c.name, c.desc, c.key_word, c.trigger_system, c.trigger_user, c.trigger_ai
			FROM wt_node_change nc
			LEFT JOIN wt_card c ON nc.detail_card = c.id
			WHERE nc.node_id = ? AND nc.level = 'card'
			ORDER BY nc.id
		`, nodeID)
		if err != nil {
			continue
		}
		for rows.Next() {
			var changeID int
			var action string
			var target sql.NullInt64
			var cardID sql.NullInt64
			var cardName, cardDesc, keyWord sql.NullString
			var triggerSystem, triggerUser, triggerAI sql.NullBool
			rows.Scan(&changeID, &action, &target, &cardID, &cardName, &cardDesc, &keyWord, &triggerSystem, &triggerUser, &triggerAI)
			if action == "add" && cardID.Valid {
				if target.Valid && int(target.Int64) == folderChangeID || !target.Valid && folderChangeID == 0 {
					cardMap[int(cardID.Int64)] = &cardState{
						card: WTCard{
							ID:            int(cardID.Int64),
							Name:          cardName.String,
							Desc:          cardDesc.String,
							KeyWord:       keyWord.String,
							TriggerSystem: triggerSystem.Bool,
							TriggerUser:   triggerUser.Bool,
							TriggerAI:     triggerAI.Bool,
						},
						changeID: changeID,
						deleted:  false,
					}
				}
			} else if action == "del" {
				delRecords = append(delRecords, delRecord{changeID: changeID})
			}
		}
		rows.Close()
	}
	for _, del := range delRecords {
		var targetChangeID int
		db.QueryRow("SELECT target FROM wt_node_change WHERE id = ?", del.changeID).Scan(&targetChangeID)
		var targetCardID int
		db.QueryRow("SELECT detail_card FROM wt_node_change WHERE id = ?", targetChangeID).Scan(&targetCardID)
		if state, ok := cardMap[targetCardID]; ok {
			state.deleted = true
		}
	}
	var result []ReplayedCard
	for _, state := range cardMap {
		if !state.deleted {
			blocks, _ := a.replayBlocksInCard(path, state.changeID)
			result = append(result, ReplayedCard{
				ID:            state.card.ID,
				Name:          state.card.Name,
				Desc:          state.card.Desc,
				KeyWord:       state.card.KeyWord,
				TriggerSystem: state.card.TriggerSystem,
				TriggerUser:   state.card.TriggerUser,
				TriggerAI:     state.card.TriggerAI,
				ChangeID:      state.changeID,
				Blocks:        blocks,
			})
		}
	}
	return result, nil
}

func (a *App) replayBlocksByZone(path []int, zone string) ([]ReplayedBlock, error) {
	type blockState struct {
		block    WTBlock
		changeID int
		deleted  bool
	}
	type delRecord struct {
		changeID int
	}
	blockMap := make(map[int]*blockState)
	var delRecords []delRecord
	for _, nodeID := range path {
		rows, err := db.Query(`
			SELECT nc.id, nc.action, b.id, b.title, b.zone
			FROM wt_node_change nc
			LEFT JOIN wt_block b ON nc.detail_block = b.id
			WHERE nc.node_id = ? AND nc.level = 'block' AND nc.target IS NULL
			ORDER BY nc.id
		`, nodeID)
		if err != nil {
			continue
		}
		for rows.Next() {
			var changeID int
			var action string
			var blockID sql.NullInt64
			var blockTitle sql.NullString
			var blockZone sql.NullString
			rows.Scan(&changeID, &action, &blockID, &blockTitle, &blockZone)
			if action == "add" && blockID.Valid && blockZone.String == zone {
				blockMap[int(blockID.Int64)] = &blockState{
					block:    WTBlock{ID: int(blockID.Int64), Title: blockTitle.String, Zone: blockZone.String},
					changeID: changeID,
					deleted:  false,
				}
			} else if action == "del" {
				delRecords = append(delRecords, delRecord{changeID: changeID})
			}
		}
		rows.Close()
	}
	for _, del := range delRecords {
		var targetChangeID int
		db.QueryRow("SELECT target FROM wt_node_change WHERE id = ?", del.changeID).Scan(&targetChangeID)
		var targetBlockID int
		db.QueryRow("SELECT detail_block FROM wt_node_change WHERE id = ?", targetChangeID).Scan(&targetBlockID)
		if state, ok := blockMap[targetBlockID]; ok {
			state.deleted = true
		}
	}
	var result []ReplayedBlock
	for _, state := range blockMap {
		if !state.deleted {
			lines, _ := a.replayLinesInBlock(path, state.changeID)
			result = append(result, ReplayedBlock{
				ID:       state.block.ID,
				Title:    state.block.Title,
				Zone:     state.block.Zone,
				ChangeID: state.changeID,
				Lines:    lines,
			})
		}
	}
	return result, nil
}

func (a *App) replayBlocksInCard(path []int, cardChangeID int) ([]ReplayedBlock, error) {
	type blockState struct {
		block    WTBlock
		changeID int
		deleted  bool
	}
	type delRecord struct {
		changeID int
	}
	blockMap := make(map[int]*blockState)
	var delRecords []delRecord
	for _, nodeID := range path {
		rows, err := db.Query(`
			SELECT nc.id, nc.action, nc.target, b.id, b.title, b.zone
			FROM wt_node_change nc
			LEFT JOIN wt_block b ON nc.detail_block = b.id
			WHERE nc.node_id = ? AND nc.level = 'block'
			ORDER BY nc.id
		`, nodeID)
		if err != nil {
			continue
		}
		for rows.Next() {
			var changeID int
			var action string
			var target sql.NullInt64
			var blockID sql.NullInt64
			var blockTitle, blockZone sql.NullString
			rows.Scan(&changeID, &action, &target, &blockID, &blockTitle, &blockZone)
			if action == "add" && blockID.Valid && target.Valid && int(target.Int64) == cardChangeID {
				blockMap[int(blockID.Int64)] = &blockState{
					block:    WTBlock{ID: int(blockID.Int64), Title: blockTitle.String, Zone: blockZone.String},
					changeID: changeID,
					deleted:  false,
				}
			} else if action == "del" {
				delRecords = append(delRecords, delRecord{changeID: changeID})
			}
		}
		rows.Close()
	}
	for _, del := range delRecords {
		var targetChangeID int
		db.QueryRow("SELECT target FROM wt_node_change WHERE id = ?", del.changeID).Scan(&targetChangeID)
		var targetBlockID int
		db.QueryRow("SELECT detail_block FROM wt_node_change WHERE id = ?", targetChangeID).Scan(&targetBlockID)
		if state, ok := blockMap[targetBlockID]; ok {
			state.deleted = true
		}
	}
	var result []ReplayedBlock
	for _, state := range blockMap {
		if !state.deleted {
			lines, _ := a.replayLinesInBlock(path, state.changeID)
			result = append(result, ReplayedBlock{
				ID:       state.block.ID,
				Title:    state.block.Title,
				Zone:     state.block.Zone,
				ChangeID: state.changeID,
				Lines:    lines,
			})
		}
	}
	return result, nil
}

type replayLineState struct {
	line      WTLine
	changeID  int
	position  *int
	deleted   bool
	nodeID    int
	nodeIndex int
}

func (a *App) replayLinesInBlock(path []int, blockChangeID int) ([]ReplayedLine, error) {
	if len(path) == 0 {
		return nil, nil
	}
	currentNodeID := path[len(path)-1]
	nodeIndexMap := make(map[int]int)
	for idx, nid := range path {
		nodeIndexMap[nid] = idx
	}
	lineMap := make(map[int]*replayLineState)
	type delRecord struct {
		changeID int
	}
	var delRecords []delRecord
	for _, nodeID := range path {
		rows, err := db.Query(`
			SELECT nc.id, nc.action, nc.target, l.id, l.sn, l.content, l.position
			FROM wt_node_change nc
			LEFT JOIN wt_line l ON nc.detail_line = l.id
			WHERE nc.node_id = ? AND nc.level = 'line'
			ORDER BY nc.id
		`, nodeID)
		if err != nil {
			continue
		}
		for rows.Next() {
			var changeID int
			var action string
			var target sql.NullInt64
			var lineID sql.NullInt64
			var lineSN, lineContent sql.NullString
			var position sql.NullInt64
			rows.Scan(&changeID, &action, &target, &lineID, &lineSN, &lineContent, &position)
			if action == "add" && lineID.Valid && target.Valid && int(target.Int64) == blockChangeID {
				var pos *int
				if position.Valid {
					p := int(position.Int64)
					pos = &p
				}
				lineMap[int(lineID.Int64)] = &replayLineState{
					line:      WTLine{ID: int(lineID.Int64), SN: lineSN.String, Content: lineContent.String},
					changeID:  changeID,
					position:  pos,
					deleted:   false,
					nodeID:    nodeID,
					nodeIndex: nodeIndexMap[nodeID],
				}
			} else if action == "del" {
				delRecords = append(delRecords, delRecord{changeID: changeID})
			}
		}
		rows.Close()
	}
	for _, del := range delRecords {
		var targetChangeID int
		db.QueryRow("SELECT target FROM wt_node_change WHERE id = ?", del.changeID).Scan(&targetChangeID)
		var targetLineID int
		db.QueryRow("SELECT detail_line FROM wt_node_change WHERE id = ?", targetChangeID).Scan(&targetLineID)
		if state, ok := lineMap[targetLineID]; ok {
			state.deleted = true
		}
	}
	var childNodes []int
	rows, err := db.Query("SELECT id FROM wt_node WHERE parent_id = ?", currentNodeID)
	if err == nil {
		for rows.Next() {
			var childID int
			rows.Scan(&childID)
			childNodes = append(childNodes, childID)
		}
		rows.Close()
	}
	childDeleteMap := make(map[int]map[int]bool)
	for _, childID := range childNodes {
		childDeleteMap[childID] = make(map[int]bool)
		var childDelTargets []int
		delRows, err := db.Query(`
			SELECT nc.target FROM wt_node_change nc
			WHERE nc.node_id = ? AND nc.level = 'line' AND nc.action = 'del'
		`, childID)
		if err != nil {
			continue
		}
		for delRows.Next() {
			var targetChangeID int
			delRows.Scan(&targetChangeID)
			childDelTargets = append(childDelTargets, targetChangeID)
		}
		delRows.Close()
		for _, targetChangeID := range childDelTargets {
			var targetLineID int
			db.QueryRow("SELECT detail_line FROM wt_node_change WHERE id = ?", targetChangeID).Scan(&targetLineID)
			childDeleteMap[childID][targetLineID] = true
		}
	}
	var unsorted []ReplayedLine
	for lineID, state := range lineMap {
		if !state.deleted {
			shape := "dot"
			if state.nodeID == currentNodeID {
				shape = "square"
			}
			color := "green"
			if len(childNodes) > 0 {
				deleteCount := 0
				for _, childID := range childNodes {
					if childDeleteMap[childID][lineID] {
						deleteCount++
					}
				}
				if deleteCount == len(childNodes) {
					color = "red"
				} else if deleteCount > 0 {
					color = "yellow"
				}
			}
			var content *string
			if strings.TrimSpace(state.line.Content) == "" {
				content = nil
			} else {
				content = &state.line.Content
			}
			unsorted = append(unsorted, ReplayedLine{
				ID:       state.line.ID,
				SN:       state.line.SN,
				Content:  content,
				SyncDot:  color + "-" + shape,
				ChangeID: state.changeID,
			})
		}
	}
	return sortLinesByPosition(unsorted, lineMap), nil
}

func sortLinesByPosition(lines []ReplayedLine, lineMap map[int]*replayLineState) []ReplayedLine {
	if len(lines) == 0 {
		return lines
	}
	idToLine := make(map[int]ReplayedLine)
	for _, line := range lines {
		idToLine[line.ID] = line
	}
	afterMap := make(map[int][]int)
	var firstLineIDs []int
	for lineID, state := range lineMap {
		if state.deleted {
			continue
		}
		if state.position == nil {
			firstLineIDs = append(firstLineIDs, lineID)
		} else {
			afterMap[*state.position] = append(afterMap[*state.position], lineID)
		}
	}
	for k := range afterMap {
		sort.Slice(afterMap[k], func(i, j int) bool {
			ni := lineMap[afterMap[k][i]].nodeIndex
			nj := lineMap[afterMap[k][j]].nodeIndex
			if ni != nj {
				return ni < nj
			}
			return afterMap[k][i] < afterMap[k][j]
		})
	}
	sort.Slice(firstLineIDs, func(i, j int) bool {
		ni := lineMap[firstLineIDs[i]].nodeIndex
		nj := lineMap[firstLineIDs[j]].nodeIndex
		if ni != nj {
			return ni < nj
		}
		return firstLineIDs[i] < firstLineIDs[j]
	})
	if len(firstLineIDs) == 0 {
		return lines
	}
	var result []ReplayedLine
	visited := make(map[int]bool)
	var queue []int
	queue = append(queue, firstLineIDs...)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true
		if line, ok := idToLine[current]; ok {
			result = append(result, line)
		}
		if nexts, ok := afterMap[current]; ok {
			queue = append(queue, nexts...)
		}
	}
	if len(result) < len(lines) {
		for _, line := range lines {
			if !visited[line.ID] {
				result = append(result, line)
			}
		}
	}
	return result
}

func sortIntSlice(s []int) {
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[j] < s[i] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

func (a *App) AddFolder(nodeID int, name string) (int, error) {
	return a.AddFolderWithParent(nodeID, nil, name)
}

func (a *App) AddFolderWithParent(nodeID int, parentChangeID *int, name string) (int, error) {
	fmt.Println("[AddFolderWithParent] 开始, nodeID:", nodeID, "parentChangeID:", parentChangeID, "name:", name)
	if strings.HasPrefix(name, "SYS_") {
		fmt.Println("[AddFolderWithParent] 拒绝: SYS_开头")
		return 0, fmt.Errorf("不允许创建SYS_开头的文件夹")
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("[AddFolderWithParent] 事务开始失败:", err)
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.Exec("INSERT INTO wt_folder (name) VALUES (?)", name)
	if err != nil {
		fmt.Println("[AddFolderWithParent] 插入wt_folder失败:", err)
		return 0, err
	}
	folderID, _ := result.LastInsertId()
	fmt.Println("[AddFolderWithParent] folderID:", folderID)
	var targetVal any
	if parentChangeID != nil {
		targetVal = *parentChangeID
	}
	result, err = tx.Exec(`
		INSERT INTO wt_node_change (action, level, node_id, target, detail_folder)
		VALUES ('add', 'folder', ?, ?, ?)
	`, nodeID, targetVal, folderID)
	if err != nil {
		fmt.Println("[AddFolderWithParent] 插入wt_node_change失败:", err)
		return 0, err
	}
	changeID, _ := result.LastInsertId()
	fmt.Println("[AddFolderWithParent] changeID:", changeID)
	if err := tx.Commit(); err != nil {
		fmt.Println("[AddFolderWithParent] 提交失败:", err)
		return 0, err
	}
	fmt.Println("[AddFolderWithParent] 成功")
	return int(changeID), nil
}

func (a *App) DeleteFolder(nodeID int, folderChangeID int) error {
	_, err := db.Exec(`
		INSERT INTO wt_node_change (action, level, node_id, target)
		VALUES ('del', 'folder', ?, ?)
	`, nodeID, folderChangeID)
	return err
}

func (a *App) AddCard(nodeID int, folderChangeID *int, name string, keyWord string) (int, error) {
	fmt.Println("[AddCard] 开始, nodeID:", nodeID, "folderChangeID:", folderChangeID, "name:", name, "keyWord:", keyWord)
	if strings.HasPrefix(name, "SYS_") {
		fmt.Println("[AddCard] 拒绝: SYS_开头")
		return 0, fmt.Errorf("不允许创建SYS_开头的卡片")
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("[AddCard] 事务开始失败:", err)
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.Exec(`
		INSERT INTO wt_card (name, desc, key_word, trigger_system, trigger_user, trigger_ai)
		VALUES (?, '', ?, 0, 1, 1)
	`, name, keyWord)
	if err != nil {
		fmt.Println("[AddCard] 插入wt_card失败:", err)
		return 0, err
	}
	cardID, _ := result.LastInsertId()
	fmt.Println("[AddCard] cardID:", cardID)
	var targetVal any
	if folderChangeID != nil {
		targetVal = *folderChangeID
	}
	result, err = tx.Exec(`
		INSERT INTO wt_node_change (action, level, node_id, target, detail_card)
		VALUES ('add', 'card', ?, ?, ?)
	`, nodeID, targetVal, cardID)
	if err != nil {
		fmt.Println("[AddCard] 插入wt_node_change失败:", err)
		return 0, err
	}
	changeID, _ := result.LastInsertId()
	fmt.Println("[AddCard] changeID:", changeID)
	if err := tx.Commit(); err != nil {
		fmt.Println("[AddCard] 提交失败:", err)
		return 0, err
	}
	fmt.Println("[AddCard] 成功")
	return int(changeID), nil
}

func (a *App) DeleteCard(nodeID int, cardChangeID int) error {
	_, err := db.Exec(`
		INSERT INTO wt_node_change (action, level, node_id, target)
		VALUES ('del', 'card', ?, ?)
	`, nodeID, cardChangeID)
	return err
}

func (a *App) AddBlock(nodeID int, cardChangeID *int, title string, zone string) (int, error) {
	fmt.Println("[AddBlock] 开始, nodeID:", nodeID, "cardChangeID:", cardChangeID, "title:", title, "zone:", zone)
	if strings.HasPrefix(title, "SYS_") {
		fmt.Println("[AddBlock] 拒绝: SYS_开头")
		return 0, fmt.Errorf("不允许创建SYS_开头的区块")
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("[AddBlock] 事务开始失败:", err)
		return 0, err
	}
	defer tx.Rollback()
	result, err := tx.Exec("INSERT INTO wt_block (title, zone) VALUES (?, ?)", title, zone)
	if err != nil {
		fmt.Println("[AddBlock] 插入wt_block失败:", err)
		return 0, err
	}
	blockID, _ := result.LastInsertId()
	fmt.Println("[AddBlock] blockID:", blockID)
	var targetVal any
	if cardChangeID != nil {
		targetVal = *cardChangeID
	}
	result, err = tx.Exec(`
		INSERT INTO wt_node_change (action, level, node_id, target, detail_block)
		VALUES ('add', 'block', ?, ?, ?)
	`, nodeID, targetVal, blockID)
	if err != nil {
		fmt.Println("[AddBlock] 插入wt_node_change失败:", err)
		return 0, err
	}
	changeID, _ := result.LastInsertId()
	fmt.Println("[AddBlock] changeID:", changeID)
	if err := tx.Commit(); err != nil {
		fmt.Println("[AddBlock] 提交失败:", err)
		return 0, err
	}
	fmt.Println("[AddBlock] 成功")
	return int(changeID), nil
}

func (a *App) DeleteBlock(nodeID int, blockChangeID int) error {
	_, err := db.Exec(`
		INSERT INTO wt_node_change (action, level, node_id, target)
		VALUES ('del', 'block', ?, ?)
	`, nodeID, blockChangeID)
	return err
}

func (a *App) AddLine(nodeID int, projectID int, blockChangeID int, content string) (int, error) {
	fmt.Println("[AddLine] 开始, nodeID:", nodeID, "projectID:", projectID, "blockChangeID:", blockChangeID, "content:", content)
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("[AddLine] 事务开始失败:", err)
		return 0, err
	}
	defer tx.Rollback()
	sn := generateSerial(tx, projectID)
	fmt.Println("[AddLine] sn:", sn)
	result, err := tx.Exec("INSERT INTO wt_line (sn, project_id, content, node_id) VALUES (?, ?, ?, ?)", sn, projectID, content, nodeID)
	if err != nil {
		fmt.Println("[AddLine] 插入wt_line失败:", err)
		return 0, err
	}
	lineID, _ := result.LastInsertId()
	fmt.Println("[AddLine] lineID:", lineID)
	result, err = tx.Exec(`
		INSERT INTO wt_node_change (action, level, node_id, target, detail_line)
		VALUES ('add', 'line', ?, ?, ?)
	`, nodeID, blockChangeID, lineID)
	if err != nil {
		fmt.Println("[AddLine] 插入wt_node_change失败:", err)
		return 0, err
	}
	changeID, _ := result.LastInsertId()
	fmt.Println("[AddLine] changeID:", changeID)
	if err := tx.Commit(); err != nil {
		fmt.Println("[AddLine] 提交失败:", err)
		return 0, err
	}
	fmt.Println("[AddLine] 成功")
	return int(changeID), nil
}

func (a *App) DeleteLine(nodeID int, lineChangeID int) error {
	_, err := db.Exec(`
		INSERT INTO wt_node_change (action, level, node_id, target)
		VALUES ('del', 'line', ?, ?)
	`, nodeID, lineChangeID)
	return err
}

func (a *App) UpdateCardKeyWord(cardID int, keyWord string) error {
	_, err := db.Exec("UPDATE wt_card SET key_word = ? WHERE id = ?", keyWord, cardID)
	return err
}

func (a *App) UpdateCardTriggers(cardID int, triggerSystem, triggerUser, triggerAI bool) error {
	ts, tu, ta := 0, 0, 0
	if triggerSystem {
		ts = 1
	}
	if triggerUser {
		tu = 1
	}
	if triggerAI {
		ta = 1
	}
	_, err := db.Exec("UPDATE wt_card SET trigger_system = ?, trigger_user = ?, trigger_ai = ? WHERE id = ?", ts, tu, ta, cardID)
	return err
}

func (a *App) UpdateBlockTitle(blockID int, title string) error {
	_, err := db.Exec("UPDATE wt_block SET title = ? WHERE id = ?", title, blockID)
	return err
}

func (a *App) UpdateLineContent(lineID int, content string) error {
	_, err := db.Exec("UPDATE wt_line SET content = ? WHERE id = ?", content, lineID)
	return err
}

func (a *App) GetSystemFolderChangeIds(projectID int) (map[string]int, error) {
	var rootNodeID int
	err := db.QueryRow(`SELECT id FROM wt_node WHERE project_id = ? AND parent_id IS NULL`, projectID).Scan(&rootNodeID)
	if err != nil {
		return nil, fmt.Errorf("获取根节点失败: %w", err)
	}
	result := make(map[string]int)
	rows, err := db.Query(`
		SELECT f.name, nc.id
		FROM wt_node_change nc
		JOIN wt_folder f ON nc.detail_folder = f.id
		WHERE nc.node_id = ? AND nc.level = 'folder' AND nc.action = 'add'
		AND f.name IN ('SYS_PRE', 'SYS_POST', 'SYS_GLOBAL')
	`, rootNodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var folderName string
		var changeID int
		if err := rows.Scan(&folderName, &changeID); err != nil {
			continue
		}
		result[folderName] = changeID
	}
	return result, nil
}

type NodeTemplate struct {
	Name    string                    `json:"name"`
	Desc    string                    `json:"desc,omitempty"`
	Folders map[string]FolderTemplate `json:"folders"`
}

type FolderTemplate struct {
	Name    string                    `json:"name"`
	Desc    string                    `json:"desc,omitempty"`
	Folders map[string]FolderTemplate `json:"folders"`
	Cards   map[string]CardTemplate   `json:"cards"`
}

type CardTemplate struct {
	Name       string                   `json:"name"`
	Desc       string                   `json:"desc,omitempty"`
	Trigger    *TriggerTemplate         `json:"trigger"`
	Blocks     map[string]BlockTemplate `json:"blocks"`
	ImageWords []string                 `json:"image"`
}

type TriggerTemplate struct {
	Mode   string   `json:"mode"`
	Words  []string `json:"words"`
	System bool     `json:"system"`
	User   bool     `json:"user"`
	AI     bool     `json:"ai"`
}

type LineData struct {
	Content *string `json:"content"`
	SyncDot string  `json:"syncDot"`
}
type BlockTemplate struct {
	Title string              `json:"title"`
	Lines map[string]LineData `json:"lines"`
}

type folderState struct {
	Name           string
	Desc           string
	ParentChangeID int
	Deleted        bool
}

type cardState struct {
	Name           string
	Desc           string
	KeyWord        string
	ImageWord      string
	TriggerMode    int
	TriggerSystem  bool
	TriggerUser    bool
	TriggerAI      bool
	ParentChangeID int
	Deleted        bool
}

type blockState struct {
	Title          string
	Zone           string
	ParentChangeID int
	Deleted        bool
}

type lineState struct {
	SN             string
	Content        string
	Position       int
	ParentChangeID int
	Deleted        bool
	NodeID         int
	LineID         int
}

type replayState struct {
	Folders        map[int]*folderState
	Cards          map[int]*cardState
	Blocks         map[int]*blockState
	Lines          map[int]*lineState
	CurrentNodeID  int
	ChildNodes     []int
	ChildDeleteMap map[int]map[int]bool
}

func (a *App) GetNodeDetail(nodeID int) (*NodeTemplate, error) {
	fmt.Println("[GetNodeDetail] 请求nodeID:", nodeID)
	var nodeName, nodeDesc string
	err := db.QueryRow(`SELECT name, COALESCE(desc, '') FROM wt_node WHERE id = ?`, nodeID).Scan(&nodeName, &nodeDesc)
	if err != nil {
		fmt.Println("[GetNodeDetail] 获取节点信息失败:", err)
		return nil, fmt.Errorf("获取节点信息失败: %w", err)
	}
	path, err := a.GetNodePath(nodeID)
	if err != nil {
		fmt.Println("[GetNodeDetail] GetNodePath失败:", err)
		return nil, err
	}
	fmt.Println("[GetNodeDetail] path:", path)
	state := a.replayAllChanges(path)
	state.CurrentNodeID = nodeID
	var childNodes []int
	rows, err := db.Query("SELECT id FROM wt_node WHERE parent_id = ?", nodeID)
	if err == nil {
		for rows.Next() {
			var childID int
			rows.Scan(&childID)
			childNodes = append(childNodes, childID)
		}
		rows.Close()
	}
	state.ChildNodes = childNodes
	state.ChildDeleteMap = make(map[int]map[int]bool)
	for _, childID := range childNodes {
		state.ChildDeleteMap[childID] = make(map[int]bool)
		delRows, err := db.Query(`
			SELECT nc.target FROM wt_node_change nc
			WHERE nc.node_id = ? AND nc.level = 'line' AND nc.action = 'del'
		`, childID)
		if err != nil {
			continue
		}
		for delRows.Next() {
			var targetChangeID int
			delRows.Scan(&targetChangeID)
			if line, ok := state.Lines[targetChangeID]; ok {
				state.ChildDeleteMap[childID][line.LineID] = true
			}
		}
		delRows.Close()
	}
	result := a.buildNodeTemplate(nodeName, nodeDesc, state)
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("[GetNodeDetail] 返回结果:\n", string(resultJSON))
	return result, nil
}

func (a *App) replayAllChanges(path []int) *replayState {
	state := &replayState{
		Folders: make(map[int]*folderState),
		Cards:   make(map[int]*cardState),
		Blocks:  make(map[int]*blockState),
		Lines:   make(map[int]*lineState),
	}
	if len(path) == 0 {
		return state
	}
	placeholders := make([]string, len(path))
	args := make([]any, len(path))
	for i, id := range path {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(`
		SELECT nc.id, nc.action, nc.level, nc.target, nc.node_id,
		       f.name AS folder_name,
		       c.name AS card_name, c.desc AS card_desc, c.key_word, c.image_word,
		       c.trigger_mode, c.trigger_system, c.trigger_user, c.trigger_ai,
		       b.title AS block_title, b.zone AS block_zone,
		       l.id AS line_id, l.sn, l.content, l.position
		FROM wt_node_change nc
		LEFT JOIN wt_folder f ON nc.detail_folder = f.id
		LEFT JOIN wt_card c ON nc.detail_card = c.id
		LEFT JOIN wt_block b ON nc.detail_block = b.id
		LEFT JOIN wt_line l ON nc.detail_line = l.id
		WHERE nc.node_id IN (%s)
		ORDER BY nc.id
	`, strings.Join(placeholders, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println("[replayAllChanges] 查询失败:", err)
		return state
	}
	defer rows.Close()
	for rows.Next() {
		var changeID int
		var action, level string
		var target sql.NullInt64
		var nodeID int
		var folderName sql.NullString
		var cardName, cardDesc, keyWord, imageWord sql.NullString
		var triggerMode, triggerSystem, triggerUser, triggerAI sql.NullInt64
		var blockTitle, blockZone sql.NullString
		var lineID sql.NullInt64
		var lineSN, lineContent sql.NullString
		var linePosition sql.NullInt64
		rows.Scan(&changeID, &action, &level, &target, &nodeID,
			&folderName,
			&cardName, &cardDesc, &keyWord, &imageWord,
			&triggerMode, &triggerSystem, &triggerUser, &triggerAI,
			&blockTitle, &blockZone,
			&lineID, &lineSN, &lineContent, &linePosition)
		parentID := 0
		if target.Valid {
			parentID = int(target.Int64)
		}
		switch level {
		case "folder":
			if action == "add" {
				state.Folders[changeID] = &folderState{
					Name:           folderName.String,
					ParentChangeID: parentID,
				}
			} else if action == "del" && parentID > 0 {
				if f, ok := state.Folders[parentID]; ok {
					f.Deleted = true
				}
			}
		case "card":
			if action == "add" {
				state.Cards[changeID] = &cardState{
					Name:           cardName.String,
					Desc:           cardDesc.String,
					KeyWord:        keyWord.String,
					ImageWord:      imageWord.String,
					TriggerMode:    int(triggerMode.Int64),
					TriggerSystem:  triggerSystem.Int64 == 1,
					TriggerUser:    triggerUser.Int64 == 1,
					TriggerAI:      triggerAI.Int64 == 1,
					ParentChangeID: parentID,
				}
			} else if action == "del" && parentID > 0 {
				if c, ok := state.Cards[parentID]; ok {
					c.Deleted = true
				}
			}
		case "block":
			if action == "add" {
				state.Blocks[changeID] = &blockState{
					Title:          blockTitle.String,
					Zone:           blockZone.String,
					ParentChangeID: parentID,
				}
			} else if action == "del" && parentID > 0 {
				if b, ok := state.Blocks[parentID]; ok {
					b.Deleted = true
				}
			}
		case "line":
			if action == "add" {
				pos := 0
				if linePosition.Valid {
					pos = int(linePosition.Int64)
				}
				lid := 0
				if lineID.Valid {
					lid = int(lineID.Int64)
				}
				state.Lines[changeID] = &lineState{
					SN:             lineSN.String,
					Content:        lineContent.String,
					Position:       pos,
					ParentChangeID: parentID,
					NodeID:         nodeID,
					LineID:         lid,
				}
			} else if action == "del" && parentID > 0 {
				if l, ok := state.Lines[parentID]; ok {
					l.Deleted = true
				}
			}
		}
	}
	fmt.Println("[replayAllChanges] folders:", len(state.Folders), "cards:", len(state.Cards), "blocks:", len(state.Blocks), "lines:", len(state.Lines))
	return state
}

func (a *App) buildNodeTemplate(name, desc string, state *replayState) *NodeTemplate {
	result := &NodeTemplate{
		Name:    name,
		Desc:    desc,
		Folders: make(map[string]FolderTemplate),
	}
	for changeID, folder := range state.Folders {
		if folder.Deleted {
			continue
		}
		if folder.ParentChangeID == 0 {
			result.Folders[fmt.Sprintf("%d", changeID)] = a.buildFolderTemplate(changeID, state)
		}
	}
	return result
}

func (a *App) buildFolderTemplate(folderChangeID int, state *replayState) FolderTemplate {
	folder := state.Folders[folderChangeID]
	ft := FolderTemplate{
		Name:    folder.Name,
		Desc:    folder.Desc,
		Folders: make(map[string]FolderTemplate),
		Cards:   make(map[string]CardTemplate),
	}
	for changeID, subFolder := range state.Folders {
		if subFolder.Deleted || subFolder.ParentChangeID != folderChangeID {
			continue
		}
		ft.Folders[fmt.Sprintf("%d", changeID)] = a.buildFolderTemplate(changeID, state)
	}
	for changeID, card := range state.Cards {
		if card.Deleted || card.ParentChangeID != folderChangeID {
			continue
		}
		ft.Cards[fmt.Sprintf("%d", changeID)] = a.buildCardTemplate(changeID, state)
	}
	return ft
}

func (a *App) buildCardTemplate(cardChangeID int, state *replayState) CardTemplate {
	card := state.Cards[cardChangeID]
	ct := CardTemplate{
		Name:   card.Name,
		Desc:   card.Desc,
		Blocks: make(map[string]BlockTemplate),
	}
	mode := "or"
	if card.TriggerMode == 1 {
		mode = "and"
	}
	var words []string
	if card.KeyWord != "" {
		words = strings.Split(card.KeyWord, "@")
	}
	ct.Trigger = &TriggerTemplate{
		Mode:   mode,
		Words:  words,
		System: card.TriggerSystem,
		User:   card.TriggerUser,
		AI:     card.TriggerAI,
	}
	if card.ImageWord != "" {
		ct.ImageWords = strings.Split(card.ImageWord, "@")
	}
	for changeID, block := range state.Blocks {
		if block.Deleted || block.ParentChangeID != cardChangeID {
			continue
		}
		ct.Blocks[fmt.Sprintf("%d", changeID)] = a.buildBlockTemplate(changeID, state)
	}
	return ct
}

func (a *App) buildBlockTemplate(blockChangeID int, state *replayState) BlockTemplate {
	block := state.Blocks[blockChangeID]
	bt := BlockTemplate{
		Title: block.Title,
		Lines: make(map[string]LineData),
	}
	for _, line := range state.Lines {
		if line.Deleted || line.ParentChangeID != blockChangeID {
			continue
		}
		shape := "dot"
		if line.NodeID == state.CurrentNodeID {
			shape = "square"
		}
		color := "green"
		if len(state.ChildNodes) > 0 {
			deleteCount := 0
			for _, childID := range state.ChildNodes {
				if state.ChildDeleteMap[childID][line.LineID] {
					deleteCount++
				}
			}
			if deleteCount == len(state.ChildNodes) {
				color = "red"
			} else if deleteCount > 0 {
				color = "yellow"
			}
		}
		var content *string
		if strings.TrimSpace(line.Content) == "" {
			content = nil
		} else {
			content = &line.Content
		}
		bt.Lines[line.SN] = LineData{
			Content: content,
			SyncDot: color + "-" + shape,
		}
	}
	return bt
}

func parseKeyWords(kw string) []string {
	kw = strings.TrimPrefix(kw, "(")
	kw = strings.TrimSuffix(kw, ")")
	if kw == "" {
		return []string{}
	}
	return strings.Split(kw, ",")
}

func (a *App) ImmediateChange(nodeID int, change ImmediateChange) (int, error) {
	changeJSON, _ := json.Marshal(change)
	fmt.Println("[ImmediateChange] nodeID:", nodeID, "change:", string(changeJSON))
	switch change.Level {
	case "folder":
		switch change.Action {
		case "add":
			result, err := a.AddFolder(nodeID, change.Name)
			fmt.Println("[ImmediateChange] AddFolder结果:", result, err)
			return result, err
		case "del":
			if change.Target == nil {
				return 0, fmt.Errorf("删除folder需要target")
			}
			err := a.DeleteFolder(nodeID, *change.Target)
			fmt.Println("[ImmediateChange] DeleteFolder结果:", err)
			return 0, err
		}
	case "card":
		switch change.Action {
		case "add":
			result, err := a.AddCard(nodeID, change.Target, change.Name, "")
			fmt.Println("[ImmediateChange] AddCard结果:", result, err)
			return result, err
		case "del":
			if change.Target == nil {
				return 0, fmt.Errorf("删除card需要target")
			}
			err := a.DeleteCard(nodeID, *change.Target)
			fmt.Println("[ImmediateChange] DeleteCard结果:", err)
			return 0, err
		}
	case "block":
		switch change.Action {
		case "add":
			var blockInfo struct {
				Title string `json:"title"`
				Zone  string `json:"zone"`
			}
			if err := json.Unmarshal([]byte(change.Name), &blockInfo); err != nil {
				blockInfo.Title = change.Name
				blockInfo.Zone = "card"
			}
			result, err := a.AddBlock(nodeID, change.Target, blockInfo.Title, blockInfo.Zone)
			fmt.Println("[ImmediateChange] AddBlock结果:", result, err)
			return result, err
		case "del":
			if change.Target == nil {
				return 0, fmt.Errorf("删除block需要target")
			}
			err := a.DeleteBlock(nodeID, *change.Target)
			fmt.Println("[ImmediateChange] DeleteBlock结果:", err)
			return 0, err
		}
	}
	return 0, fmt.Errorf("未知的change类型: %s %s", change.Action, change.Level)
}

func parseID(s string) (int, error) {
	var id int
	_, err := fmt.Sscanf(s, "%d", &id)
	return id, err
}

type SaveNodeChangesResult struct {
	Blocks map[string]map[string]LineData `json:"blocks"`
}

func (a *App) SaveNodeChanges(nodeID int, projectID int, changes SaveChanges) (*SaveNodeChangesResult, error) {
	changesJSON, _ := json.MarshalIndent(changes, "", "  ")
	fmt.Println("[SaveNodeChanges] 收到请求 nodeID:", nodeID, "projectID:", projectID)
	fmt.Println("[SaveNodeChanges] changes:\n", string(changesJSON))
	var nodePath []int
	needPath := false
	savedBlockChangeIDs := []int{}
	for changeIDStr, blockChange := range changes.Block {
		if blockChange.Content != nil {
			needPath = true
			changeID, _ := parseID(changeIDStr)
			savedBlockChangeIDs = append(savedBlockChangeIDs, changeID)
		}
	}
	if needPath {
		var err error
		nodePath, err = a.GetNodePath(nodeID)
		if err != nil {
			fmt.Println("[SaveNodeChanges] GetNodePath失败:", err)
			return nil, err
		}
		fmt.Println("[SaveNodeChanges] 节点路径:", nodePath)
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("[SaveNodeChanges] 开启事务失败:", err)
		return nil, err
	}
	defer tx.Rollback()
	for changeIDStr, folderChange := range changes.Folder {
		changeID, _ := parseID(changeIDStr)
		var folderID int
		tx.QueryRow("SELECT detail_folder FROM wt_node_change WHERE id = ?", changeID).Scan(&folderID)
		fmt.Println("[SaveNodeChanges] 更新folder changeID:", changeID, "folderID:", folderID, folderChange)
		if folderChange.Name != nil {
			tx.Exec("UPDATE wt_folder SET name = ? WHERE id = ?", *folderChange.Name, folderID)
		}
	}
	for changeIDStr, cardChange := range changes.Card {
		changeID, _ := parseID(changeIDStr)
		var cardID int
		tx.QueryRow("SELECT detail_card FROM wt_node_change WHERE id = ?", changeID).Scan(&cardID)
		fmt.Println("[SaveNodeChanges] 更新card changeID:", changeID, "cardID:", cardID, cardChange)
		if cardChange.Name != nil {
			tx.Exec("UPDATE wt_card SET name = ? WHERE id = ?", *cardChange.Name, cardID)
		}
		if cardChange.Desc != nil {
			tx.Exec("UPDATE wt_card SET desc = ? WHERE id = ?", *cardChange.Desc, cardID)
		}
		if cardChange.Trigger != nil {
			kw := strings.Join(cardChange.Trigger.Words, "@")
			triggerMode := 0
			if cardChange.Trigger.Mode == "and" {
				triggerMode = 1
			}
			ts, tu, ta := 0, 0, 0
			if cardChange.Trigger.System {
				ts = 1
			}
			if cardChange.Trigger.User {
				tu = 1
			}
			if cardChange.Trigger.AI {
				ta = 1
			}
			tx.Exec("UPDATE wt_card SET key_word = ?, trigger_mode = ?, trigger_system = ?, trigger_user = ?, trigger_ai = ? WHERE id = ?",
				kw, triggerMode, ts, tu, ta, cardID)
		}
		if cardChange.Image != nil {
			imageWord := strings.Join(cardChange.Image, "@")
			tx.Exec("UPDATE wt_card SET image_word = ? WHERE id = ?", imageWord, cardID)
		}
	}
	for changeIDStr, blockChange := range changes.Block {
		changeID, _ := parseID(changeIDStr)
		var blockID int
		tx.QueryRow("SELECT detail_block FROM wt_node_change WHERE id = ?", changeID).Scan(&blockID)
		fmt.Println("[SaveNodeChanges] 更新block changeID:", changeID, "blockID:", blockID, "Name:", blockChange.Name, "Content是否nil:", blockChange.Content == nil)
		if blockChange.Name != nil {
			tx.Exec("UPDATE wt_block SET title = ? WHERE id = ?", *blockChange.Name, blockID)
		}
		if blockChange.Content != nil {
			fmt.Println("[SaveNodeChanges] 处理block内容 changeID:", changeID, "content长度:", len(*blockChange.Content))
			if err := a.diffAndSaveBlockContent(tx, nodePath, nodeID, projectID, changeID, *blockChange.Content); err != nil {
				fmt.Println("[SaveNodeChanges] diffAndSaveBlockContent失败:", err)
				return nil, err
			}
		} else {
			fmt.Println("[SaveNodeChanges] block.Content为nil, 跳过")
		}
	}
	fmt.Println("[SaveNodeChanges] 提交事务")
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	result := &SaveNodeChangesResult{
		Blocks: make(map[string]map[string]LineData),
	}
	for _, blockChangeID := range savedBlockChangeIDs {
		lines, err := a.replayLinesInBlock(nodePath, blockChangeID)
		if err != nil {
			continue
		}
		blockLines := make(map[string]LineData)
		for _, line := range lines {
			blockLines[line.SN] = LineData{
				Content: line.Content,
				SyncDot: line.SyncDot,
			}
		}
		result.Blocks[fmt.Sprintf("%d", blockChangeID)] = blockLines
	}
	return result, nil
}

func calcMD5(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

type diffLine struct {
	SN       string
	Content  string
	LineID   int
	ChangeID int
	NodeID   int
	Used     bool
}

func (a *App) diffAndSaveBlockContent(tx *sql.Tx, path []int, nodeID, projectID, blockChangeID int, newContent string) error {
	fmt.Println("[diffAndSaveBlockContent] nodeID:", nodeID, "projectID:", projectID, "blockChangeID:", blockChangeID, "newContent长度:", len(newContent), "内容:", newContent)
	fmt.Println("[diffAndSaveBlockContent] 节点路径:", path)
	var oldLines []diffLine
	deletedChangeIDs := make(map[int]bool)
	for _, pathNodeID := range path {
		addRows, err := tx.Query(`
			SELECT nc.id, l.id, l.sn, l.content, nc.node_id
			FROM wt_node_change nc
			JOIN wt_line l ON nc.detail_line = l.id
			WHERE nc.node_id = ? AND nc.level = 'line' AND nc.action = 'add' AND nc.target = ?
			ORDER BY nc.id
		`, pathNodeID, blockChangeID)
		if err != nil {
			continue
		}
		for addRows.Next() {
			var chgID int
			var lineID int
			var sn, content string
			var chgNodeID int
			addRows.Scan(&chgID, &lineID, &sn, &content, &chgNodeID)
			oldLines = append(oldLines, diffLine{
				SN:       sn,
				Content:  content,
				LineID:   lineID,
				ChangeID: chgID,
				NodeID:   chgNodeID,
				Used:     false,
			})
		}
		addRows.Close()
		delRows, err := tx.Query(`
			SELECT nc.target FROM wt_node_change nc
			WHERE nc.node_id = ? AND nc.level = 'line' AND nc.action = 'del'
		`, pathNodeID)
		if err != nil {
			continue
		}
		for delRows.Next() {
			var targetChangeID int
			delRows.Scan(&targetChangeID)
			deletedChangeIDs[targetChangeID] = true
		}
		delRows.Close()
	}
	var filtered []diffLine
	for _, line := range oldLines {
		if !deletedChangeIDs[line.ChangeID] {
			filtered = append(filtered, line)
		}
	}
	oldLines = filtered
	fmt.Println("[diffAndSaveBlockContent] 现有行数:", len(oldLines))
	newLines := strings.Split(newContent, "\n")
	if len(newLines) == 1 && newLines[0] == "" {
		newLines = []string{}
	}
	fmt.Println("[diffAndSaveBlockContent] 新行数:", len(newLines))
	oldMD5Map := make(map[string][]*diffLine)
	for i := range oldLines {
		h := calcMD5(oldLines[i].Content)
		oldMD5Map[h] = append(oldMD5Map[h], &oldLines[i])
	}
	var lastLineID *int
	for _, newLine := range newLines {
		h := calcMD5(newLine)
		if candidates, ok := oldMD5Map[h]; ok && len(candidates) > 0 {
			var found *diffLine
			for _, c := range candidates {
				if !c.Used {
					found = c
					break
				}
			}
			if found != nil {
				found.Used = true
				tx.Exec("UPDATE wt_line SET position = ? WHERE id = ?", lastLineID, found.LineID)
				lastLineID = &found.LineID
				fmt.Println("[diffAndSaveBlockContent] 匹配已有行 sn:", found.SN, "lineID:", found.LineID, "更新position")
				continue
			}
		}
		sn := generateSerial(tx, projectID)
		result, err := tx.Exec("INSERT INTO wt_line (sn, project_id, content, position, node_id) VALUES (?, ?, ?, ?, ?)",
			sn, projectID, newLine, lastLineID, nodeID)
		if err != nil {
			fmt.Println("[diffAndSaveBlockContent] 插入line失败:", err)
			return err
		}
		lineID, _ := result.LastInsertId()
		lineIDInt := int(lineID)
		lastLineID = &lineIDInt
		_, err = tx.Exec(`
			INSERT INTO wt_node_change (action, level, node_id, target, detail_line)
			VALUES ('add', 'line', ?, ?, ?)
		`, nodeID, blockChangeID, lineID)
		if err != nil {
			fmt.Println("[diffAndSaveBlockContent] 插入node_change失败:", err)
			return err
		}
		fmt.Println("[diffAndSaveBlockContent] 新建行 sn:", sn, "lineID:", lineID)
	}
	for _, line := range oldLines {
		if !line.Used {
			fmt.Println("[diffAndSaveBlockContent] 删除行 sn:", line.SN, "lineID:", line.LineID, "originNodeID:", line.NodeID)
			var delLinePosition sql.NullInt64
			tx.QueryRow("SELECT position FROM wt_line WHERE id = ?", line.LineID).Scan(&delLinePosition)
			tx.Exec("UPDATE wt_line SET position = ? WHERE position = ?", delLinePosition, line.LineID)
			_, err := tx.Exec(`
				INSERT INTO wt_node_change (action, level, node_id, target)
				VALUES ('del', 'line', ?, ?)
			`, nodeID, line.ChangeID)
			if err != nil {
				fmt.Println("[diffAndSaveBlockContent] 删除行node_change失败:", err)
				return err
			}
		}
	}
	return nil
}
