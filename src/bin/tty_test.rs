use core::slice::SlicePattern;
use std::{thread, io::Write};

use anyhow::Result;
use log::{warn, info};
use robust_arduino_serial::{Order, write_i8, read_i8};
use serial::prelude::*;
use vim_with_me::arduino::Arduino;


fn main() -> Result<()> {
    env_logger::init();
    let serial_port = "/dev/ttyACM1";

    let mut ard = Arduino::create(serial_port)?;

    ard.write_str("JS Really Sucks")?;

    let mut vec = vec![];
    'outer: loop {
        std::thread::sleep(std::time::Duration::from_secs(1));

        while let Ok(b) = ard.read_i8() {
            vec.push(b as u8);

            if let Ok(s) = std::str::from_utf8(vec.as_slice()) {
                if s == "JS Really Sucks" {
                    break 'outer;
                }
            }
        }

    }

    println!("WE DID IT! {:?}", vec);

    return Ok(());
}
