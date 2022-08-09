use log::{info, error, warn};
use tokio_tungstenite::tungstenite::Message;
use anyhow::Result;

pub mod system;
pub mod vim;

use crate::prime::message::{PrimeMessage, PrimeMessageContent};

use self::{system::handle_system_command, vim::{VimSender, VimMessage}};

pub async fn handle_message(msg: Message, sender: VimSender) -> Result<()> {
    let msg = msg.to_text()?;
    let msg = serde_json::from_str::<PrimeMessage>(msg)?;

    warn!("where am i? {:?}", msg);
    match msg.content {
        PrimeMessageContent::SystemCommand(s) => handle_system_command(s).await?,
        PrimeMessageContent::VimRTL => {
            warn!("rtl");
            let cmd = VimMessage::rtl();
            if !cmd.is_valid() {
                error!("invalid rtl {:?}", cmd);
                return Ok(());
            }

            match sender.send(cmd).await {
                Err(e) => {
                    error!("got an error from sending vim command {}", e);
                },
                _ => {}
            }
        },
        PrimeMessageContent::VimMotion(cmd) => {
            let cmd = VimMessage::motion(cmd);
            if !cmd.is_valid() {
                error!("invalid vim motion {:?}", cmd);
                return Ok(());
            }

            match sender.send(cmd).await {
                Err(e) => {
                    error!("got an error from sending vim command {}", e);
                },
                _ => {}
            }
        }

        _ => {
            info!("please handle this arm immediately {:?}", msg);
        }
    }

    return Ok(());
}

