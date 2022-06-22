use std::net::TcpListener;

use anyhow::{Result, anyhow};
use clap::Parser;
use tokio_tungstenite::connect_async;
use vim_with_me::clap_me_daddy::Opts;

#[tokio::main]
async fn main() -> Result<()> {
    env_logger::init();

    if dotenv::dotenv().is_err() {
        return Err(anyhow!("Dotenv really didn't make it"));
    }

    let opts = Opts::parse();
    let addr = Url::new(format!("ws://{}:{}", opts.address, opts.port));
    let (ws_stream, _) = connect_async(url).await.expect("Failed to connect");

    let (write, read) = ws_stream.split();

    return Ok(());
}
