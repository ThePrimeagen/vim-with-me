use ipipe::Pipe;
use std::io::{BufRead, BufWriter, BufReader};

use std::{sync::mpsc::channel, fs};
use std::thread;
use std::time::Duration;

use std::os::unix::net::UnixStream;
use std::io::prelude::*;

use systemstat::{System, Platform};

use anyhow::Result;

#[tokio::main]
async fn main() -> Result<()> {
    let sys = System::new();
    let mut stream: UnixStream;

    let path = "/tmp/vim_me_daddy.sock";
    if fs::metadata(path).is_ok() {
        fs::remove_file(path)?;
    }

    loop {
        match UnixStream::connect(path) {
            Err(e) => {
                println!("errored {}.  Waiting 5 seconds", e);
                tokio::time::sleep(tokio::time::Duration::from_secs(5)).await;
            },

            Ok(s) => {
                stream = s;
                break;
            }
        }

    }

    let (tx, rx) = channel();

    thread::spawn(move || {
        loop {
            thread::sleep(Duration::from_secs(3));

            match sys.cpu_temp() {
                Ok(cpu_temp) => {
                    tx.send(format!("Temp is: {}", cpu_temp));
                },
                Err(x) => {
                    tx.send("error".to_string());
                }
            };
        }

    });

    thread::spawn(move || async {
        let path = "/tmp/vim_me_daddy.sock";
        loop {
            if let Ok(data) = fs::read(path) {
                println!("got the data {:?}", data);
            } else {
                println!("couuldn't rnnead the data");
            }
            tokio::time::sleep(Duration::from_secs(5)).await;
        }
    });

    loop {
        let _ = rx
            .try_recv()
            .map(|reply| stream.write_all(reply.as_bytes()));
    }

    /*
    let path = "/tmp/vim_me_daddy.sock";
    if fs::metadata(path).is_ok() {
        fs::remove_file(path)?;
    }

    tokio::time::sleep(Duration::from_secs(1)).await;

    unix_named_pipe::create(path, None)?;
    let mut writer = unix_named_pipe::open_write(path)?;

    writer.write(b"hello world")?;

    return Ok(());
    */
}

    /*
    let mut pipe = Pipe::create().unwrap();
    println!("Name: {}", pipe.path().display());

    let writer = pipe.clone();
    thread::spawn(move || print_nums(writer));
    for line in BufReader::new(pipe).lines() {
        println!("{}", line.unwrap());
    }

    return Ok(());
}

fn print_nums(mut pipe: Pipe) {
    for i in 1..=10 {
        writeln!(&mut pipe, "{}", i).unwrap();
    }
    write!(&mut pipe, "{}", CANCEL as char).unwrap();
}
*/

