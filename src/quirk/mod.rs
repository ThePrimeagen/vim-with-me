use anyhow::{Result};
use futures::{StreamExt, stream::SplitStream};
use log::{error, warn, info};
use tokio::{net::TcpStream, sync::mpsc, time};
use tokio_tungstenite::{connect_async, MaybeTlsStream, WebSocketStream};
use serde::{Deserialize, Serialize};

pub type Receiver = mpsc::Receiver<QuirkMessage>;
pub type Sender = mpsc::Sender<QuirkMessage>;
type Read = SplitStream<WebSocketStream<MaybeTlsStream<TcpStream>>>;

#[derive(Debug, Deserialize, Serialize)]
pub struct TwitchReward {
    pub id: String,
    pub title: String,
    pub prompt: String,
    pub cost: usize,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct TwitchData {
    pub broadcaster_user_id: String,
    pub broadcaster_user_login: String,
    pub broadcaster_user_name: String,
    pub id: String,
    pub user_id: String,
    pub user_login: String,
    pub user_name: String,
    pub user_input: String,
    pub status: String,
    pub redeemed_at: String,
    pub reward: TwitchReward,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct TwitchRedeem {
    pub data: TwitchData,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct TwitchSubEvent {
    pub data: TwitchRedeem,
    pub source: String,
    pub r#type: String
}

#[derive(Debug)]
pub enum QuirkMessage {
    Close,
    Message(TwitchSubEvent),
}

pub struct Quirk {
    pub rx: Option<Receiver>,
    tx: Sender,
}

#[derive(Deserialize, Serialize, Debug)]
struct RequestBody {
    access_token: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct ResponseBody {
    access_token: String,
}

pub async fn get_quirk_token() -> Result<String> {
    let token = std::env::var("QUIRK_TOKEN").expect("QUIRK_TOKEN should be an env variable");
    let request = RequestBody {
        access_token: token,
    };

    let client = reqwest::Client::new();
    let res: ResponseBody = client
        .post("https://websocket.quirk.tools/token")
        .json(&request)
        .header("Content-Type", "application/json")
        .send()
        .await?
        .json()
        .await?;

    return Ok(res.access_token);
}

async fn run_quirk(tx: Sender, mut read: Read) {
    while let Some(Ok(msg)) = read.next().await {
        if let Ok(msg) = msg.to_text() {
            if let Ok(msg) = serde_json::from_str(msg) {
                info!("received message from quirk {:?}", msg);
                match tx.send(QuirkMessage::Message(msg)).await {
                    Err(_) => break,
                    _ => {}
                }
            }
        }
    }

    match tx.send(QuirkMessage::Close).await {
        Err(e) => {
            error!("unable to send out close message, probably need to restart program {}.", e);
        },
        _ => {}
    }
}

impl Quirk {
    pub fn new() -> Quirk {
        let (tx, rx) = mpsc::channel::<QuirkMessage>(10);

        return Quirk {
            rx: Some(rx),
            tx,
        }
    }

    pub async fn connect(&self, url: &str) -> Result<()> {
        let quirk_token = get_quirk_token().await?;
        let url = format!("{}?access_token={}", url, quirk_token);
        let (ws_stream, _) = connect_async(url).await.expect("Failed to connect");
        let (_, read) = ws_stream.split();

        tokio::spawn(run_quirk(self.tx.clone(), read));

        return Ok(());
    }

    pub fn get_receiver(&mut self) -> Option<Receiver> {
        return self.rx.take();
    }
}

pub async fn run_forver_quirky(tx: Sender) -> Result<()> {
    loop {
        let mut quirk = Quirk::new();
        warn!("about to (re)connect to quirk");
        match quirk.connect("wss://websocket.quirk.tools/").await {
            Ok(_) => {
                let mut rx = quirk.get_receiver().expect("quirk should always have an rx to give");
                loop {
                    tokio::select! {
                        msg = rx.recv() => {
                            match msg {
                                Some(QuirkMessage::Close) => {
                                    warn!("quirk has closed");
                                    break;
                                },

                                Some(msg) => {
                                    if let Err(e) = tx.send(msg).await {
                                        error!("error'd or emitting quirk message {}", e);
                                        break;
                                    }
                                },

                                None => break,
                            }
                        },

                        _ = tokio::time::sleep(tokio::time::Duration::from_secs(60 * 30)) => {
                            error!("reconnecting to quirk (forceful disconnect)");
                            break;
                        }
                    }
                }

            },
            _ => {},
        }

        warn!("Disconnected from quirk, reconnecting in 5");
        tokio::time::sleep(time::Duration::from_secs(5)).await;
    }
}
