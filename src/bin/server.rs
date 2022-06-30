use anyhow::{Result, anyhow};
use futures::future::join;
use log::debug;
use serde::{Deserialize, Serialize};
use tokio::sync::mpsc::channel;
use vim_with_me::{quirk::run_forver_quirky, prime::ws::server};

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

    let (tx, rx) = channel(10);

    let res = join(
        tokio::spawn(server("0.0.0.0:42069", rx)),
        tokio::spawn(run_forver_quirky(tx))
    ).await;

    res.0??;
    res.1??;

    return Ok(());
}
