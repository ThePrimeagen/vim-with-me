use serde::{Serialize, Deserialize};
use tokio::sync::mpsc;

use crate::quirk::{Receiver, QuirkMessage};

#[derive(Debug, Serialize, Deserialize)]
pub enum PrimeMessageType {
    StatusLineUpdate,
    SystemCommand,
}

#[derive(Debug, Serialize, Deserialize)]
pub enum PrimeMessageContent {
    SystemCommand(String),
    StatusLineUpdate(String, String),

    // TODO: we don't even know yet
    VimCommand(),
    VimColorScheme(),
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PrimeMessage {
    pub r#type: PrimeMessageType,
    pub content: PrimeMessageContent,
}

pub type Sender = mpsc::Sender<PrimeMessage>;

pub async fn receive(mut rx: Receiver, tx: Sender) {
    while let Some(QuirkMessage::Message(msg)) = rx.recv().await {
    }
}

