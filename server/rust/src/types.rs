#[allow(unused)]
pub mod version {
    pub const PROTOCOL_VERSION: u8 = 0x01;

    pub const COMMAND_LOGIN: u16 = 0x01;

    pub const GENERIC_RESPONSE: u16 = 0x03;
}


#[derive(Debug, PartialEq, Eq)]
pub struct Header {
    version: u8,
    key: u16,
}

impl Header {
    pub fn new(version: u8, key: u16) -> Self {
        Self { version, key }
    }

    /// Get a reference to the request header's version.
    pub fn version(&self) -> u8 {
        self.version
    }

    /// Get a reference to the request header's key.
    pub fn key(&self) -> u16 {
        self.key
    }
}