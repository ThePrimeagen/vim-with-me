use anyhow::{Result, Context};

use std::time::Duration;
use std::thread;
use log::{warn, info};
use serial::{prelude::*, unix::TTYPort};
use robust_arduino_serial::*;

// Default settings of Arduino
// see: https://www.arduino.cc/en/Serial/Begin
const SETTINGS: serial::PortSettings = serial::PortSettings {
    baud_rate:    serial::Baud57600,
    char_size:    serial::Bits8,
    parity:       serial::ParityNone,
    stop_bits:    serial::Stop1,
    flow_control: serial::FlowNone,
};

pub struct Arduino {
    port: TTYPort
}

impl Arduino {
    pub fn create(serial_port: &str) -> Result<Arduino> {
        warn!("Opening port: {:?}", serial_port);

        let mut port = serial::open(serial_port)?;
        port.configure(&SETTINGS)?;

        return Ok(Arduino {
            port,
        });
    }

    pub fn write_str(&mut self, s: &str) -> Result<usize> {
        for ele in s.bytes() {
            write_i8(&mut self.port, ele as i8).context("unable to write_i8")?;
        }
        write_i8(&mut self.port, 0).context("unable to write_i8")?;

        return Ok(s.len());
    }

    pub fn write_i8(&mut self, val: i8) -> Result<usize> {
        return write_i8(&mut self.port, val).context("unable to write_i8");
    }

    pub fn write_i16(&mut self, val: i16) -> Result<usize> {
        return write_i16(&mut self.port, val).context("unable to write_i8");
    }

    pub fn read_i8(&mut self) -> Result<i8> {
        return read_i8(&mut self.port).context("trying to read i8 from serial");
    }
}
