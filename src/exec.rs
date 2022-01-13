use std::sync::Arc;

use log::info;
use tokio::task::JoinHandle;

use crate::{error::VWMError, command, opts::ClientOpts};
pub fn execute_command(websocket_message: String, opts: Arc<ClientOpts>) -> Result<Option<JoinHandle<()>>, VWMError> {

    info!("execute_command: Processing websocket_message {:?}", websocket_message);
    let command: command::Command = websocket_message.into();

    if let command::Command::Noop = command {
        return Ok(None);
    }

    // TOKIO SPAWN?? I Still don't understand this very well
    let join_handle = tokio::spawn(async move {
        command.on(opts.clone());
        tokio::time::sleep(command.duration()).await;
        command.off(opts.clone());
    });

    return Ok(Some(join_handle));
}
