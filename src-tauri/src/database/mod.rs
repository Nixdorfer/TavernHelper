pub mod schema;

use rusqlite::{Connection, Result};
use std::path::Path;
use std::sync::{Mutex, OnceLock};

static DB: OnceLock<Mutex<Connection>> = OnceLock::new();

pub fn init_db(app_data_dir: &Path) -> Result<()> {
    let db_path = app_data_dir.join("tavern.db");
    tracing::info!("Initializing database at: {:?}", db_path);
    let conn = Connection::open(&db_path)?;
    conn.execute_batch("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;")?;
    schema::init_schema(&conn)?;
    DB.set(Mutex::new(conn)).map_err(|_| rusqlite::Error::InvalidQuery)?;
    tracing::info!("Database initialized successfully");
    Ok(())
}

pub fn with_db<F, T>(f: F) -> Result<T>
where
    F: FnOnce(&Connection) -> Result<T>,
{
    let guard = DB.get().ok_or(rusqlite::Error::InvalidQuery)?.lock().unwrap();
    f(&guard)
}

pub fn with_db_log<F, T>(op: &str, f: F) -> Result<T>
where
    F: FnOnce(&Connection) -> Result<T>,
{
    tracing::debug!("[DB] {}", op);
    let guard = DB.get().ok_or(rusqlite::Error::InvalidQuery)?.lock().unwrap();
    let result = f(&guard);
    match &result {
        Ok(_) => tracing::debug!("[DB] {} - OK", op),
        Err(e) => tracing::error!("[DB] {} - Error: {}", op, e),
    }
    result
}

pub fn with_db_tx<F, T>(op: &str, f: F) -> Result<T>
where
    F: FnOnce(&Connection) -> Result<T>,
{
    tracing::debug!("[DB TX] {} - BEGIN", op);
    let guard = DB.get().ok_or(rusqlite::Error::InvalidQuery)?.lock().unwrap();
    guard.execute("BEGIN", [])?;
    let result = f(&guard);
    match &result {
        Ok(_) => {
            guard.execute("COMMIT", [])?;
            tracing::debug!("[DB TX] {} - COMMIT", op);
        }
        Err(e) => {
            let _ = guard.execute("ROLLBACK", []);
            tracing::error!("[DB TX] {} - ROLLBACK: {}", op, e);
        }
    }
    result
}

#[allow(dead_code)]
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
