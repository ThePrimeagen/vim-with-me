use anyhow::Result;
use dotenv::dotenv;

use futures_util::StreamExt;

use log::error;
use tokio_tungstenite::connect_async;
use url::Url;
use vim_with_me::client::{handle_message, self};

#[tokio::main]
async fn main() -> Result<()> {
    dotenv().expect("dotenv to work");
    env_logger::init();

    let url = "ws://0.0.0.0:42069";

    let (socket, _) = connect_async(Url::parse(url).unwrap()).await?;
    let (_, mut incoming) = socket.split();
    let sender = client::vim::handle_tcp_to_vim("0.0.0.0:6969");

    // So far, we don't need async beyond simple async await
    while let Some(Ok(msg)) = incoming.next().await {
        match handle_message(msg, sender.clone()).await {
            Err(e) => {
                error!("error from handle_message {}", e);
            },
            _ => {}
        }
    };

    return Ok(());
}
