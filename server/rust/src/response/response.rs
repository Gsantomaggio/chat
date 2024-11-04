use std::io::Write;
use crate::codec::{Decoder, Encoder};
use crate::codec::decoder::read_u16;
use crate::commands::login::LoginResponse;
use crate::error::{DecodeError, EncodeError};
use crate::types::Header;

#[derive(Debug, PartialEq, Eq)]
pub struct Response {
    header: Header,
    pub(crate) kind: ResponseKind,
}

#[derive(Debug, PartialEq, Eq, Clone)]
#[repr(u16)]
pub enum ResponseCode {
    Ok = 1,
    UserNotFound = 2,
    UserAlreadyLogged = 3,
}



impl ResponseCode {
    fn to_u16(self) -> u16 {
        self as u16
    }
}


#[derive(Debug, PartialEq, Eq)]
pub enum ResponseKind {
    Login(LoginResponse),
}


impl Decoder for ResponseCode {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), crate::error::DecodeError> {
        let (input, code) = read_u16(input)?;
        Ok((input, code.try_into()?))
    }
}

impl Response {
    pub fn new(header: Header, kind: ResponseKind) -> Self {
        Self { header, kind }
    }

    pub fn correlation_id(&self) -> Option<u32> {
        match &self.kind {
            ResponseKind::Login(login) => Some(login.correlation_id)
        }
    }

    pub fn response_code(&self) -> Option<u16> {
        match &self.kind {
            ResponseKind::Login(login) => Some(login.response_code.clone().to_u16())
        }
    }
}

impl Encoder for Response {
    fn encoded_size(&self) -> u32 {
        self.header.encoded_size() +
            4 + // correlation id
            2 // response code
    }

    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        self.header.encode(writer)?;
        self.correlation_id().unwrap().encode(writer)?;
        self.response_code().unwrap().encode(writer)?;
        Ok(())
    }
}



