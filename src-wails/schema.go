package main

import (
	"database/sql"
	"fmt"
)

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS wt_project (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL DEFAULT 'create' CHECK(type IN ('play', 'create')),
			current_node INTEGER REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE SET NULL,
			create_time TEXT NOT NULL,
			update_time TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS wt_node (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL REFERENCES wt_project(id) ON UPDATE CASCADE ON DELETE CASCADE,
			parent_id INTEGER REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE CASCADE,
			name TEXT NOT NULL,
			desc TEXT,
			UNIQUE(name, project_id)
		);
		CREATE INDEX IF NOT EXISTS idx_wt_node_project ON wt_node(project_id);
		CREATE INDEX IF NOT EXISTS idx_wt_node_parent ON wt_node(parent_id);
		CREATE TABLE IF NOT EXISTS wt_branch_tag (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			parent_id INTEGER NOT NULL REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE CASCADE,
			child_id INTEGER NOT NULL REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE CASCADE,
			name TEXT NOT NULL,
			desc TEXT,
			UNIQUE(parent_id, child_id)
		);
		CREATE TABLE IF NOT EXISTS wt_folder (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS wt_card (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			desc TEXT,
			key_word TEXT,
			image_word TEXT,
			trigger_mode INTEGER NOT NULL DEFAULT 0,
			trigger_system INTEGER NOT NULL DEFAULT 0,
			trigger_user INTEGER NOT NULL DEFAULT 1,
			trigger_ai INTEGER NOT NULL DEFAULT 1
		);
		CREATE TABLE IF NOT EXISTS wt_block (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			zone TEXT NOT NULL CHECK(zone IN ('pre', 'post', 'global', 'card'))
		);
		CREATE TABLE IF NOT EXISTS wt_line (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sn TEXT NOT NULL,
			project_id INTEGER NOT NULL REFERENCES wt_project(id) ON UPDATE CASCADE ON DELETE CASCADE,
			content TEXT,
			position INTEGER REFERENCES wt_line(id) ON UPDATE CASCADE ON DELETE SET NULL,
			node_id INTEGER REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE SET NULL,
			UNIQUE(sn, project_id)
		);
		CREATE INDEX IF NOT EXISTS idx_wt_line_project ON wt_line(project_id);
		CREATE INDEX IF NOT EXISTS idx_wt_line_position ON wt_line(position);
		CREATE INDEX IF NOT EXISTS idx_wt_line_node ON wt_line(node_id);
		CREATE TABLE IF NOT EXISTS wt_node_change (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT NOT NULL CHECK(action IN ('add', 'del')),
			level TEXT NOT NULL CHECK(level IN ('folder', 'card', 'block', 'line')),
			node_id INTEGER NOT NULL REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE CASCADE,
			target INTEGER REFERENCES wt_node_change(id) ON UPDATE CASCADE ON DELETE CASCADE,
			detail_folder INTEGER REFERENCES wt_folder(id) ON UPDATE CASCADE ON DELETE RESTRICT,
			detail_card INTEGER REFERENCES wt_card(id) ON UPDATE CASCADE ON DELETE RESTRICT,
			detail_block INTEGER REFERENCES wt_block(id) ON UPDATE CASCADE ON DELETE RESTRICT,
			detail_line INTEGER REFERENCES wt_line(id) ON UPDATE CASCADE ON DELETE RESTRICT
		);
		CREATE INDEX IF NOT EXISTS idx_wt_node_change_node ON wt_node_change(node_id);
		CREATE INDEX IF NOT EXISTS idx_wt_node_change_target ON wt_node_change(target);
		CREATE TABLE IF NOT EXISTS wt_app (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS wt_conversation (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			app_id INTEGER NOT NULL REFERENCES wt_app(id) ON UPDATE CASCADE ON DELETE CASCADE,
			name TEXT NOT NULL DEFAULT '新对话',
			current_node INTEGER REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE SET NULL
		);
		CREATE INDEX IF NOT EXISTS idx_wt_conversation_app ON wt_conversation(app_id);
		CREATE TABLE IF NOT EXISTS wt_dialogue (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			conversation_id INTEGER NOT NULL REFERENCES wt_conversation(id) ON UPDATE CASCADE ON DELETE CASCADE,
			create_time TEXT NOT NULL,
			request_content TEXT NOT NULL,
			response_content TEXT NOT NULL,
			request_system_prompt TEXT,
			response_system_prompt TEXT,
			node_id INTEGER REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE SET NULL,
			request_point INTEGER,
			response_point INTEGER,
			request_token INTEGER,
			response_token INTEGER
		);
		CREATE INDEX IF NOT EXISTS idx_wt_dialogue_conversation ON wt_dialogue(conversation_id);
		CREATE TABLE IF NOT EXISTS wt_dialogue_image (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dialogue_id INTEGER NOT NULL REFERENCES wt_dialogue(id) ON UPDATE CASCADE ON DELETE CASCADE,
			image_url TEXT,
			image_path TEXT,
			prompt TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_wt_dialogue_image_dialogue ON wt_dialogue_image(dialogue_id);
		CREATE TABLE IF NOT EXISTS wt_image (
			id TEXT PRIMARY KEY,
			hash TEXT NOT NULL UNIQUE,
			local_path TEXT NOT NULL,
			remote_url TEXT DEFAULT '',
			file_name TEXT DEFAULT '',
			file_size INTEGER DEFAULT 0,
			created_at TEXT NOT NULL,
			folder_path TEXT DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS idx_wt_image_hash ON wt_image(hash);
		CREATE INDEX IF NOT EXISTS idx_wt_image_folder ON wt_image(folder_path);
	`)
	if err != nil {
		return fmt.Errorf("创建表失败: %w", err)
	}
	if err := migrateSchema(db); err != nil {
		return fmt.Errorf("迁移表失败: %w", err)
	}
	return nil
}

func migrateSchema(db *sql.DB) error {
	var hasTypeColumn int
	err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('wt_project') WHERE name='type'`).Scan(&hasTypeColumn)
	if err != nil {
		return err
	}
	if hasTypeColumn == 0 {
		_, err = db.Exec(`ALTER TABLE wt_project ADD COLUMN type TEXT NOT NULL DEFAULT 'create' CHECK(type IN ('play', 'create'))`)
		if err != nil {
			return err
		}
	}
	var hasCurrentNodeColumn int
	err = db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('wt_project') WHERE name='current_node'`).Scan(&hasCurrentNodeColumn)
	if err != nil {
		return err
	}
	if hasCurrentNodeColumn == 0 {
		_, err = db.Exec(`ALTER TABLE wt_project ADD COLUMN current_node INTEGER REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE SET NULL`)
		if err != nil {
			return err
		}
	}
	var hasPositionColumn int
	err = db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('wt_line') WHERE name='position'`).Scan(&hasPositionColumn)
	if err != nil {
		return err
	}
	if hasPositionColumn == 0 {
		_, err = db.Exec(`ALTER TABLE wt_line ADD COLUMN position INTEGER REFERENCES wt_node_change(id) ON UPDATE CASCADE ON DELETE SET NULL`)
		if err != nil {
			return err
		}
		db.Exec(`CREATE INDEX IF NOT EXISTS idx_wt_line_position ON wt_line(position)`)
	}
	var hasTriggerModeColumn int
	err = db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('wt_card') WHERE name='trigger_mode'`).Scan(&hasTriggerModeColumn)
	if err != nil {
		return err
	}
	if hasTriggerModeColumn == 0 {
		_, err = db.Exec(`ALTER TABLE wt_card ADD COLUMN trigger_mode INTEGER NOT NULL DEFAULT 0`)
		if err != nil {
			return err
		}
	}
	var hasImageWordColumn int
	err = db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('wt_card') WHERE name='image_word'`).Scan(&hasImageWordColumn)
	if err != nil {
		return err
	}
	if hasImageWordColumn == 0 {
		_, err = db.Exec(`ALTER TABLE wt_card ADD COLUMN image_word TEXT`)
		if err != nil {
			return err
		}
	}
	db.Exec(`
		UPDATE wt_line SET position = (
			SELECT nc.detail_line FROM wt_node_change nc WHERE nc.id = wt_line.position
		) WHERE position IS NOT NULL AND EXISTS (
			SELECT 1 FROM wt_node_change nc WHERE nc.id = wt_line.position
		)
	`)
	var hasLineNodeIDColumn int
	db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('wt_line') WHERE name='node_id'`).Scan(&hasLineNodeIDColumn)
	if hasLineNodeIDColumn == 0 {
		db.Exec(`ALTER TABLE wt_line ADD COLUMN node_id INTEGER REFERENCES wt_node(id) ON UPDATE CASCADE ON DELETE SET NULL`)
		db.Exec(`CREATE INDEX IF NOT EXISTS idx_wt_line_node ON wt_line(node_id)`)
		db.Exec(`UPDATE wt_line SET node_id = (
			SELECT nc.node_id FROM wt_node_change nc
			WHERE nc.detail_line = wt_line.id AND nc.action = 'add'
			LIMIT 1
		)`)
	}
	return nil
}
