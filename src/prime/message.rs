use serde::{Serialize, Deserialize};

use crate::quirk::QuirkMessage;

#[derive(Debug, Serialize, Deserialize)]
pub enum PrimeMessageContent {
    SystemCommand(String),
    StatusLineUpdate(String, String),

    // TODO: we don't even know yet
    VimCommand(String),
    VimColorScheme(),
}

// TODO: Why why why why
#[derive(Debug, Serialize, Deserialize)]
pub struct PrimeMessage {
    pub content: PrimeMessageContent,
}

