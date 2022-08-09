use serde::{Serialize, Deserialize};

#[derive(Debug, Serialize, Deserialize)]
pub enum PrimeMessageContent {
    SystemCommand(String),
    StatusLineUpdate(String, String),

    // TODO: we don't even know yet
    VimMotion(String),
    VimRTL,
    VimColorScheme(),
}

// TODO: Why why why why
#[derive(Debug, Serialize, Deserialize)]
pub struct PrimeMessage {
    pub content: PrimeMessageContent,
}

