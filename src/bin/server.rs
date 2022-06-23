use anyhow::{Result, anyhow};
use futures::StreamExt;
use log::debug;
use tokio_tungstenite::connect_async;
use serde::{Deserialize, Serialize};
use vim_with_me::quirk::{Quirk, QuirkMessage};

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

#[tokio::main]
async fn main() -> Result<()> {
    env_logger::init();

    debug!("starting program");

    if dotenv::dotenv().is_err() {
        return Err(anyhow!("Dotenv really didn't make it"));
    }

    loop {
        let mut quirk = Quirk::new();
        quirk.connect("wss://websocket.quirk.tools/").await?;

        let mut rx = quirk.get_receiver().unwrap();

        while let Some(msg) = rx.recv().await {
            match msg {
                QuirkMessage::Close => {
                    break;
                },

                QuirkMessage::Message(msg) => {
                    println!("got message");
                    println!("message {:?}", msg);
                    println!("done printing message");
                }
            }
        }
    }
}
