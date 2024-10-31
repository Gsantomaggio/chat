use crate::types::version::PROTOCOL_VERSION;

pub mod login;



pub trait Command {
    fn key(&self) -> u16;
    fn version(&self) -> u8 {
        PROTOCOL_VERSION
    }
}