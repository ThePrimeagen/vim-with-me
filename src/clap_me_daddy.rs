use clap::Parser;

#[derive(Debug, Parser)]
#[clap(name = "me daddy")]
pub struct Opts {
    #[clap(short = 'a', long = "address", default_value="0.0.0.0")]
    pub address: String,

    #[clap(short = 'p', long = "port", default_value="42069")]
    pub port: u16,
}

