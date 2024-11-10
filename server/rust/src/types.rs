#[allow(unused)]
pub mod version {
    pub const PROTOCOL_VERSION: u8 = 0x01;

    pub const COMMAND_LOGIN: u16 = 0x01;
    pub const COMMAND_MESSAGE: u16 = 0x02;


    pub const GENERIC_RESPONSE: u16 = 0x03;
}

impl Decoder for Option<String> {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError> {
        let (input, len) = read_i16(input)?;

        if len == 0 {
            return Ok((input, None));
        }
        let (bytes, input) = input.split_at(len as usize);
        let string = String::from_utf8(bytes.to_vec())?;
        Ok((input, Some(string)))
    }
}

impl Decoder for u32 {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError> {
        read_u32(input).map(Ok)?
    }
}

impl Decoder for u64 {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError> {
        read_u64(input).map(Ok)?
    }
}



impl Encoder for u32 {
    fn encoded_size(&self) -> u32 {
        4
    }

    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        writer.write_u32::<BigEndian>(*self)?;
        Ok(())
    }
}


impl Encoder for u64 {
    fn encoded_size(&self) -> u32 {
        8
    }

    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        writer.write_u64::<BigEndian>(*self)?;
        Ok(())
    }
}



impl Encoder for u16 {
    fn encoded_size(&self) -> u32 {
        2
    }

    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        writer.write_u16::<BigEndian>(*self)?;
        Ok(())
    }
}

impl Encoder for &str {
    fn encoded_size(&self) -> u32 {
        2 + self.len() as u32
    }

    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        writer.write_i16::<BigEndian>(self.len() as i16)?;
        writer.write_all(self.as_bytes())?;
        Ok(())
    }
}

macro_rules! reader {
    ( $fn:ident, $size:expr, $ret:ty) => {
        #[allow(unused)]
        pub fn $fn(input: &[u8]) -> Result<(&[u8], $ret), IncompleteError> {
            check_len(input, $size)?;
            let x = byteorder::BigEndian::$fn(input);
            Ok((&input[$size..], x))
        }
    };
}

reader!(read_u32, 4, u32);
reader!(read_u64, 8, u64);
reader!(read_i16, 2, i16);
reader!(read_u16, 2, u16);

pub fn check_len(input: &[u8], size: usize) -> Result<(), IncompleteError> {
    if input.len() < size {
        return Err(IncompleteError(size));
    }
    Ok(())
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

impl Decoder for Header {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError> {
        if input.len() < 3 {
            return Err(DecodeError::Incomplete(IncompleteError(3)));
        }
        let version = input[0];
        let key = u16::from_be_bytes([input[1], input[2]]);
        Ok((&input[3..], Self::new(version, key)))
    }
}

impl Encoder for Header {
    fn encoded_size(&self) -> u32 {
        1 + 2
    }

    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        writer.write_u8(self.version())?;
        writer.write_u16::<BigEndian>(self.key())?;
        Ok(())
    }
}

use std::io::Write;
use std::string::FromUtf8Error;
use byteorder::{BigEndian, ByteOrder, WriteBytesExt};

impl From<IncompleteError> for DecodeError {
    fn from(err: IncompleteError) -> Self {
        DecodeError::Incomplete(err)
    }
}

#[derive(PartialEq, Eq, Debug)]
pub struct LoginRequest {
    pub correlation_id: u32,
    pub user_name: String,
}


impl LoginRequest {
    pub fn new(correlation_id: u32, user_name: String) -> Self {
        Self {
            correlation_id,
            user_name,
        }
    }
}

impl Decoder for LoginRequest {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError> {
        let (input, correlation_id) = u32::decode(input)?;
        let (input, user_name) = Option::decode(input)?;
        Ok((input, Self::new(correlation_id, user_name.unwrap())))
    }
}

#[derive(Debug)]
pub struct IncompleteError(pub usize);

#[derive(Debug)]
pub enum DecodeError {
    Incomplete(IncompleteError),
    Utf8Error(FromUtf8Error),
    UnknownResponseCode(u16),
    UnsupportedResponseType(u16),
    MismatchSize(usize),
    MessageParse(String),
    InvalidFormatCode(u8),
    Empty,
}

#[derive(Debug)]
pub enum EncodeError {
    Io(std::io::Error),
    MaxSizeError(usize),
}

impl From<std::io::Error> for EncodeError {
    fn from(err: std::io::Error) -> Self {
        EncodeError::Io(err)
    }
}

impl From<FromUtf8Error> for DecodeError {
    fn from(err: FromUtf8Error) -> Self {
        DecodeError::Utf8Error(err)
    }
}


pub trait Decoder
where
    Self: Sized,
{
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError>;
}


pub trait Encoder {
    fn encoded_size(&self) -> u32;
    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError>;
}


// pub trait Command {
//     fn key(&self) -> u16;
//     fn version(&self) -> u8 {
//         PROTOCOL_VERSION
//     }
// }