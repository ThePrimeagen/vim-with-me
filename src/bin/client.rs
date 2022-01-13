use std::sync::Arc;

use dotenv::dotenv;

use futures_util::StreamExt;

use log::{info, error};
use master::{opts::ClientOpts, error::VWMError, exec::execute_command};
use structopt::StructOpt;
use tokio_tungstenite::{connect_async, tungstenite::Message};
use url::Url;

#[tokio::main]
async fn main() -> Result<(), VWMError> {
    dotenv().expect("dotenv to work");
    env_logger::init();

    let opts = Arc::new(ClientOpts::from_args());
    let url = format!("ws://{}", opts.server);

    let (socket, _) = connect_async(Url::parse(url.as_str()).unwrap()).await?;

    let (_, mut incoming) = socket.split();

    // So far, we don't need async beyond simple async await
    while let Some(msg) = incoming.next().await {
        info!("Message Received {:?}", msg);

        if let Ok(Message::Text(m)) = msg {
            if let Ok(Some(join_handle)) = execute_command(m, opts.clone()) {
                join_handle.await?;
            }
        } else if let Err(e) = msg {
            error!("Error from websocket_message: {:?}", e);
        }
    };

    return Ok(());
}
