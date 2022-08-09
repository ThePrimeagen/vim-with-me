use anyhow::Result;
use std::process::Command;

use log::info;

pub async fn handle_system_command(command: String) -> Result<()> {

    match command.as_str() {
            // "Xrandr": new SystemCommand("xrandr --output HDMI-0 --brightness 0.05", "xrandr --output HDMI-0
            // --brightness 1", 5000),
        "xrandr" => {
            Command::new("xrandr")
                .arg("--output")
                .arg("DP-0")
                .arg("--brightness")
                .arg("0.02")
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


