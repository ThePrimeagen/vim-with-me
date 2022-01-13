use std::{time::Duration, sync::Arc};

use log::info;

use crate::{events::TwitchCommand, error::VWMError, opts::ClientOpts};

// "Xrandr": new SystemCommand("xrandr --output HDMI-0 --brightness 0.05", "xrandr --output HDMI-0
// --brightness 1", 5000),

pub enum Command {
    Xrandr(usize),
    Noop,
}

impl From<String> for Command {
    fn from(value: String) -> Self {
        if let Ok(TwitchCommand::TWITCH_EVENTSUB { data }) = serde_json::from_str(&value) {
            match data.reward.title.as_str() {
                "Xrandr" => return Command::Xrandr(5),
                _ => return Command::Noop,
            }
        }

        return Command::Noop;
    }
}

impl Command {
    pub fn duration(&self) -> Duration {
        return match self {
            // TODO: Macro
            Command::Xrandr(d) => Duration::from_secs(*d as u64),
            _ => Duration::from_secs(0),
        };
    }
    pub fn on(&self, opts: Arc<ClientOpts>) {
        match self {
            // TODO: Macro
            Command::Xrandr(_) => {
                info!("xrandr --output {} --brightness 0.05", opts.monitor);
                std::process::Command::new("xrandr")
                    .args(["--output", &opts.monitor, "--brightness", "0.05"])
                    .spawn()
                    .expect("xrandr command failed to start");
            },
            _ => {}
        }
    }

    pub fn off(&self, opts: Arc<ClientOpts>) {
        match self {
            Command::Xrandr(_) => {
                info!("xrandr --output {} --brightness 1", opts.monitor);
                std::process::Command::new("xrandr")
                    .args(["--output", &opts.monitor, "--brightness", "1"])
                    .spawn()
                    .expect("xrandr command failed to start");
            },
            _ => {}
        }
    }
}
