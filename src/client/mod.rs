use log::info;
use std::process::Command;
use tokio_tungstenite::tungstenite::Message;
use anyhow::Result;

use crate::prime::message::{PrimeMessage, PrimeMessageContent};

async fn handle_system_command(command: String) -> Result<()> {

    match command.as_str() {
            // "Xrandr": new SystemCommand("xrandr --output HDMI-0 --brightness 0.05", "xrandr --output HDMI-0
            // --brightness 1", 5000),
        "xrandr" => {
            Command::new("xrandr")
                .arg("--output")
                .arg("DP-0")
                .arg("--brightness")
                .arg("0.05")
                .spawn()?;

            tokio::time::sleep(tokio::time::Duration::from_millis(5000)).await;

            Command::new("xrandr")
                .arg("--output")
                .arg("DP-0")
                .arg("--brightness")
                .arg("1")
                .spawn()?;
        },

        _ => {
            info!("not handling system command {}", command);
        }
    }

    return Ok(());
}

pub async fn handle_message(msg: Message) -> Result<()> {
    let msg = msg.to_text()?;
    let msg = serde_json::from_str::<PrimeMessage>(msg)?;

    match msg.content {
        PrimeMessageContent::SystemCommand(s) => {
            handle_system_command(s).await?;
        },

        _ => {
            info!("please handle this arm immediately {:?}", msg);
        }
    }

    return Ok(());
}

