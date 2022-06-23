
use futures::{StreamExt, TryStreamExt, future};
use futures_channel::mpsc::unbounded;

use serde::{Deserialize, Serialize};

use tokio::task::JoinHandle;
use tokio_tungstenite::{connect_async};

use url::Url;

use crate::error::DoxMeDaddyError;
use crate::forwarder::{ForwarderEvent, ReceiverGiverAsync};
use crate::opts::ServerOpts;
use crate::{async_receiver_giver};

pub struct Quirk {
    pub join_handle: JoinHandle<Result<(), tokio_tungstenite::tungstenite::Error>>,
    rx: Option<futures_channel::mpsc::UnboundedReceiver<ForwarderEvent>>
}

#[derive(Deserialize, Serialize, Debug)]
struct RequestBody {
    access_token: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct ResponseBody {
    access_token: String,
}


async_receiver_giver!(Quirk);

impl Quirk {
    pub async fn new(_opts: &ServerOpts) -> Result<Quirk, DoxMeDaddyError> {
        let quirk_token = get_quirk_token().await?;
        let url = format!("wss://websocket.quirk.tools?access_token={}", quirk_token);

        let (socket, _) = connect_async(Url::parse(url.as_str()).unwrap()).await.expect("Can't connect");
        let (_, incoming) = socket.split();
        let (tx, rx) = unbounded();

        let join_handle = tokio::spawn(incoming.try_for_each(move |msg| {
            if !msg.is_text() {
                return future::ok(());
            }
            if let Ok(text) = msg.into_text() {
                tx.unbounded_send(ForwarderEvent::QuirkMessageRaw(text)).expect("test");
            }
            return future::ok(());
        }));

        return Ok(Quirk {
            join_handle,
            rx: Some(rx),
        });
    }
}
