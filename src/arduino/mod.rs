use anyhow::{Result, Context};

use std::time::Duration;
use std::thread;
use log::{warn, info};
use serial::{prelude::*, unix::TTYPort};
use robust_arduino_serial::*;

// Default settings of Arduino
// see: https://www.arduino.cc/en/Serial/Begin
const SETTINGS: serial::PortSettings = serial::PortSettings {
    baud_rate:    serial::Baud115200,
    char_size:    serial::Bits8,
    parity:       serial::ParityNone,
    stop_bits:    serial::Stop1,
    flow_control: serial::FlowNone,
};

struct Arduino {
    port: TTYPort
}

impl Arduino {
    pub fn create(serial_port: &str) -> Result<Arduino> {

        warn!("Opening port: {:?}", serial_port);

        let mut port = serial::open(&serial_port)?;
        port.configure(&SETTINGS)?;
        port.set_timeout(Duration::from_secs(30))?;

        loop
        {
            info!("Waiting for Arduino...");
            let order = Order::HELLO as i8;
            write_i8(&mut port, order).unwrap();
            let received_order = Order::from_i8(read_i8(&mut port)?);
            if let Some(received_order) = received_order {
                if received_order == Order::ALREADY_CONNECTED {
                    break;
                }
            }
            thread::sleep(Duration::from_secs(1));
        }

        warn!("Connected to Arduino");

        return Ok(Arduino {
            port,
        });
    }

    pub fn write_i8(&mut self, val: i8) -> Result<usize> {
        return write_i8(&mut self.port, val).context("unable to write_i8");
    }

    pub fn write_i16(&mut self, val: i16) -> Result<usize> {
        return write_i16(&mut self.port, val).context("unable to write_i8");
    }
}
