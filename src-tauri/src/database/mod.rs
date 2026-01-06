pub mod schema;

use rusqlite::{Connection, Result};
use std::path::Path;
use std::sync::{Mutex, OnceLock};

static DB: OnceLock<Mutex<Connection>> = OnceLock::new();

pub fn init_db(app_data_dir: &Path) -> Result<()> {
    let db_path = app_data_dir.join("tavern.db");
    let conn = Connection::open(&db_path)?;
    conn.execute_batch("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;")?;
    schema::init_schema(&conn)?;
    DB.set(Mutex::new(conn)).map_err(|_| rusqlite::Error::InvalidQuery)?;
    Ok(())
}

pub fn with_db<F, T>(f: F) -> Result<T>
where
    F: FnOnce(&Connection) -> Result<T>,
{
    let guard = DB.get().ok_or(rusqlite::Error::InvalidQuery)?.lock().unwrap();
    f(&guard)
}

pub fn with_db_mut<F, T>(f: F) -> Result<T>
where
    F: FnOnce(&mut Connection) -> Result<T>,
{
    let mut guard = DB.get().ok_or(rusqlite::Error::InvalidQuery)?.lock().unwrap();
    f(&mut guard)
}

pub fn db_get_config(key: &str) -> Result<Option<String>> {
    with_db(|conn| {
        let mut stmt = conn.prepare("SELECT value FROM wt_config WHERE key = ?")?;
        let result: Result<String> = stmt.query_row([key], |row| row.get(0));
        match result {
            Ok(value) => Ok(Some(value)),
            Err(rusqlite::Error::QueryReturnedNoRows) => Ok(None),
            Err(e) => Err(e),
        }
    })
}

pub fn db_set_config(key: &str, value: &str) -> Result<()> {
    with_db(|conn| {
        conn.execute(
            "INSERT OR REPLACE INTO wt_config (key, value) VALUES (?, ?)",
            [key, value],
        )?;
        Ok(())
    })
}
