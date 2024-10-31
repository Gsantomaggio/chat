
#[allow(unused)]
pub mod version {
    pub const PROTOCOL_VERSION: u8 = 0x01;

    pub const COMMAND_LOGIN: u16 = 0x01;
}


#[derive(Debug, PartialEq, Eq)]
pub struct Header {
    key: u16,
    version: u8,
}

impl Header {
    pub fn new(key: u16, version: u8) -> Self {
        Self { key, version }
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