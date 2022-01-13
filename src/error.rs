use tokio::task::JoinError;


#[derive(Debug)]
pub enum VWMError {
    Unknown,
    WebsocketError(tokio_tungstenite::tungstenite::Error),
    JsonParseError(serde_json::Error),
    JoinHandleError(JoinError),
    ToCommandError(String),
}

impl From<tokio_tungstenite::tungstenite::Error> for VWMError {
    fn from(e: tokio_tungstenite::tungstenite::Error) -> Self {
        return Self::WebsocketError(e);
    }
}

impl From<JoinError> for VWMError {
    fn from(e: JoinError) ->  Self {
        return Self::JoinHandleError(e);
    }
}

impl From<serde_json::Error> for VWMError {
    fn from(e: serde_json::Error) ->  Self {
        return Self::JsonParseError(e);
    }
}


