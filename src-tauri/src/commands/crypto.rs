use base64::{engine::general_purpose::URL_SAFE_NO_PAD, Engine};
use ed25519_dalek::SigningKey;
use rand::rngs::OsRng;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct KeyPair {
    pub public_key: String,
    pub private_key: String,
}

#[tauri::command]
pub fn generate_encryption_key_pair() -> Result<KeyPair, String> {
    let signing_key = SigningKey::generate(&mut OsRng);
    let verifying_key = signing_key.verifying_key();
    let public_key = URL_SAFE_NO_PAD.encode(verifying_key.as_bytes());
    let private_key = URL_SAFE_NO_PAD.encode(signing_key.to_bytes());
    Ok(KeyPair {
        public_key,
        private_key,
    })
}
