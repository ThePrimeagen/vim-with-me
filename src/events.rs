use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct Reward {
    pub title: String,
}

#[derive(Deserialize, Debug)]
pub struct TwitchData {
    pub reward: Reward,
    pub user_name: String,
}

#[derive(Deserialize, Debug)]
#[serde(tag = "source")]
#[allow(non_camel_case_types)]
pub enum TwitchCommand {
    TWITCH_CHAT {
        text: String,
        name: String,
        sub: bool,
    },

    TWITCH_EVENTSUB {
        data: TwitchData
    }
}


