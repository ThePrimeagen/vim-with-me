use anyhow::{Result, anyhow};
use futures::{StreamExt, SinkExt};
use log::{debug, error};
use tokio::{net::TcpListener, sync::mpsc};
use tokio_tungstenite::tungstenite::Message;

use crate::quirk::{Receiver, QuirkMessage};

use super::message::{PrimeMessage, PrimeMessageContent};

pub type Sender = mpsc::Sender<PrimeMessage>;

impl TryFrom<QuirkMessage> for PrimeMessage {
    type Error = anyhow::Error;

    fn try_from(value: QuirkMessage) -> Result<Self> {
        if let QuirkMessage::Message(msg) = value {
            match msg.data.data.reward.title.as_str() {
                "Xrandr" => {
                    return Ok(PrimeMessage {
                        content: PrimeMessageContent::SystemCommand("xrandr".to_string()),
                    });
                },

                "VimMotion" => {
                    println!("HELP ME {:?}", msg.data);

                    return Ok(PrimeMessage {
                        content: PrimeMessageContent::VimMotion(msg.data.data.user_input),
                    });
                },

                "Giveaway FEMboi" => {
                    return Err(anyhow!("unsupported"));
                },

                _ => {
                    return Err(anyhow!("unrecognized reward"));
                }
            }
        }
        return Err(anyhow!("help me"));
    }
}

pub async fn server(addr: &str, mut rx: Receiver) -> Result<()> {

    // Create the event loop and TCP listener we'll accept connections on.
    let try_socket = TcpListener::bind(&addr).await;
    let listener = try_socket.expect("Failed to bind");
    debug!("Listening on: {}", addr);

    let mut outgoing = vec![];
    // Let's spawn the handling of each connection in a separate task.
    loop {
        tokio::select! {
            Ok((stream, _)) = listener.accept() => {
                let ws_stream = match tokio_tungstenite::accept_async(stream)
                    .await {
                        Err(e) => {
                            error!("unable to accept the connection {}", e);
                            continue;
                        },
                        Ok(v) => v
                    };

                let (o, _) = ws_stream.split();

                outgoing.push(o);
            }

            Some(msg) = rx.recv() => {
                let mut errs = vec![];
                let msg: Result<PrimeMessage> = msg.try_into();
                if msg.is_err() {
                    continue;
                }
                let msg = msg.expect("this should never fail");
                let out_msg = serde_json::to_string(&msg).expect("there should never be an error");

                for (idx, out) in outgoing.iter_mut().enumerate() {
                    match out.send(Message::Text(out_msg.clone())).await {
                        Err(e) => {
                            error!("sending to client, but error'd {:?}", e);
                            errs.push(idx);
                        },
                        _ => {}
                    }
                }

                errs.into_iter().rev().for_each(|x| {
                    let _ = outgoing.remove(x);
                });
            }
        }
    }
}


