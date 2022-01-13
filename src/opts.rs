use structopt::StructOpt;

#[derive(Debug, StructOpt, Clone)]
pub struct ClientOpts {

    /// The port to use for the events to be served on
    #[structopt(short = "s", long = "server", default_value = "69420.theprimeagen.tv:42069")]
    pub server: String,

    #[structopt(short = "m", long = "monitor", default_value = "HDMI-0")]
    pub monitor: String,

}
