use log::{info, error, warn};
use systemstat::Duration;
use tokio::{net::TcpListener, sync::mpsc::{Sender, channel}, io::AsyncWriteExt};
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

pub async fn handle_message(msg: Message, sender: VimSender) -> Result<()> {
    let msg = msg.to_text()?;
    let msg = serde_json::from_str::<PrimeMessage>(msg)?;

    match msg.content {
        PrimeMessageContent::SystemCommand(s) => {
            handle_system_command(s).await?;
        },

        PrimeMessageContent::VimCommand(cmd) => {
        }

        _ => {
            info!("please handle this arm immediately {:?}", msg);
        }
    }

    return Ok(());
}

pub enum VimMessage {
    VimMotion(usize, String)
}

impl VimMessage {
    pub fn motion(s: String) -> VimMessage {
        return VimMessage::VimMotion(0, s);
    }
}

pub type VimSender = Sender<VimMessage>;

fn encode_vim_message(msg: VimMessage) -> Vec<u8> {
    let mut out = vec![];

    match msg {
        VimMessage::VimMotion(r#type, motion) => {
            out.push((1 + motion.len()) as u8);
            out.push(r#type as u8);
            motion.chars().for_each(|c| out.push(c as u8));
        }
    }

    return out;
}

pub fn handle_tcp_to_vim(addr: &'static str) -> VimSender {
    let (tx, mut rx) = channel(100);

    tokio::spawn(async move {
        let mut list_of_listeners = vec![];
        'outer_loop: loop {
            let listener = match TcpListener::bind(addr).await {
                Err(e) => {
                    error!("couldn't create the tcp listener: {}", e);
                    tokio::time::sleep(Duration::from_millis(5000)).await;
                    continue;
                },
                Ok(v) => v,
            };

            loop {
                tokio::select! {
                    connection = listener.accept() => {
                        if let Ok((tcp_connection, _)) = connection {
                            list_of_listeners.push(tcp_connection);
                        } else {
                            break 'outer_loop;
                        }
                    },
                    vim_msg = rx.recv() => {
                        if list_of_listeners.is_empty() {
                            warn!("no listeners for vim commands.  this means you have screwed up.");
                            continue;
                        }

                        if let Some(msg) = vim_msg {
                            let mut error_idx = vec![];
                            let msg = encode_vim_message(msg);
                            for (idx, listener) in list_of_listeners.iter_mut().enumerate() {
                                match listener.write(&msg).await {
                                    Err(_) => {
                                        error_idx.push(idx);
                                    },
                                    _ => {}
                                }
                            }

                            error_idx
                                .into_iter()
                                .rev()
                                .for_each(|idx| {
                                    list_of_listeners.remove(idx);
                                });
                        }
                    }
                };
            }
        }
    });

    return tx;
}

