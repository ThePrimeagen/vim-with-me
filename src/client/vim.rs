use log::{error, warn};
use num_enum::IntoPrimitive;
use systemstat::Duration;
use tokio::{net::TcpListener, sync::mpsc::{Sender, channel}, io::AsyncWriteExt};

#[derive(Debug, IntoPrimitive)]
#[repr(u8)]
pub enum VimEncodeType {
    Motion = 0,
    RTL = 1,
    Chat = 2,
}

#[derive(Debug)]
pub enum VimMessage {
    RTL(VimEncodeType),
    Motion(VimEncodeType, String),
    Chat(VimEncodeType, String),
}

impl VimMessage {
    pub fn motion(s: String) -> VimMessage {
        return VimMessage::Motion(VimEncodeType::Motion, s);
    }

    pub fn rtl() -> VimMessage {
        return VimMessage::RTL(VimEncodeType::RTL);
    }

    pub fn chat(yes_or_no: String) -> VimMessage {
        return VimMessage::Chat(VimEncodeType::Chat, yes_or_no);
    }

    pub fn is_valid(&self) -> bool {
        match self {
            VimMessage::Motion(_, cmd) => {
                return self.is_valid_motion(cmd);
            },

            VimMessage::RTL(_) => return true,

            VimMessage::Chat(_, yes_or_no) if yes_or_no == "yes" || yes_or_no == "no" => {
                return true;
            },

            _ => return false,
        }
    }

    fn is_valid_motion(&self, cmd: &String) -> bool {
        let index = cmd.chars().position(|c| !c.is_ascii_digit()).unwrap();
        let (num, cmd) = cmd.split_at(index);
        let num_res = str::parse::<usize>(num);

        if num.len() > 0 && !num_res.is_ok() {
            return false;
        }

        return match cmd {
            "j" | "h" | "k" | "l" => true,
            _ => false,
        }
    }
}

pub type VimSender = Sender<VimMessage>;

fn encode_message_with_string(r#type: VimEncodeType, str: String, mut out: Vec<u8>) -> Vec<u8> {
    out.push((1 + str.len()) as u8);
    out.push(r#type.into());
    str.chars().for_each(|c| out.push(c as u8));

    return out;
}

fn encode_vim_message(msg: VimMessage) -> Vec<u8> {
    let mut out = vec![];

    match msg {
        VimMessage::Motion(r#type, motion) => {
            out = encode_message_with_string(r#type, motion, out);
        },

        VimMessage::RTL(r#type) => {
            out.push(1 as u8);
            out.push(r#type.into());
        },

        VimMessage::Chat(r#type, yes_or_no) => {
            out = encode_message_with_string(r#type, yes_or_no, out);
        }
    }

    return out;
}

async fn handle_vim_message(msg: VimMessage, list_of_listeners: &mut Vec<tokio::net::TcpStream>) -> Vec<usize> {
    let mut out = vec![];

    warn!("sending vim message {:?}", msg);
    let msg = encode_vim_message(msg);
    for (idx, listener) in list_of_listeners.iter_mut().enumerate() {
        match listener.write(&msg).await {
            Err(_) => {
                out.push(idx);
            },
            _ => {}
        }
    }

    return out;
}

pub fn handle_tcp_to_vim(addr: &'static str) -> VimSender {
    let (tx, mut rx) = channel(100);

    tokio::spawn(async move {
        let mut list_of_listeners = vec![];
        'outer_loop: loop {
            let listener = match TcpListener::bind(addr).await {
                Err(e) => {
                    error!("couldn't create the tcp listener: {}", e);
                    tokio::time::sleep(Duration::from_millis(5000)).await;
                    continue;
                },
                Ok(v) => v,
            };

            loop {
                tokio::select! {
                    connection = listener.accept() => {
                        if let Ok((tcp_connection, _)) = connection {
                            list_of_listeners.push(tcp_connection);
                        } else {
                            break 'outer_loop;
                        }
                    },
                    vim_msg = rx.recv() => {
                        if list_of_listeners.is_empty() {
                            warn!("no listeners for vim commands.  this means you have screwed up.");
                            continue;
                        }

                        let error_idx = match vim_msg {
                            None => continue,
                            Some(v) => {
                                handle_vim_message(v, &mut list_of_listeners).await
                            }
                        };

                        error_idx
                            .into_iter()
                            .rev()
                            .for_each(|idx| {
                                list_of_listeners.remove(idx);
                            });
                    }
                };
            }
        }
    });

    return tx;
}


