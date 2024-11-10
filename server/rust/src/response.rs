use std::io::Write;
use crate::login::LoginResponse;
use crate::response::responses::RESPONSE_CODE_OK;
use crate::types::{read_u16, DecodeError, Decoder, EncodeError, Encoder, Header};

#[derive(Debug, PartialEq, Eq)]
pub struct Response {
    header: Header,
    pub(crate) kind: ResponseKind,
}

pub mod responses {
    pub const RESPONSE_CODE_OK: u16 = 1;
}


#[derive(Debug, PartialEq, Eq, Clone)]
#[repr(u16)]
pub enum ResponseCode {
    Ok = 0x01,
    UserNotFound = 0x03,
    UserAlreadyLogged = 0x04,
}


#[derive(Debug, PartialEq, Eq)]
pub enum ResponseKind {
    Login(LoginResponse),
}

impl Decoder for ResponseCode {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError> {
        let (input, code) = read_u16(input)?;
        Ok((input, code.try_into()?))
    }
}

impl ResponseCode {
    fn to_u16(self) -> u16 {
        self as u16
    }
}

impl TryFrom<u16> for ResponseCode {
    type Error = DecodeError;

    fn try_from(value: u16) -> Result<Self, Self::Error> {
        match value {
            RESPONSE_CODE_OK => Ok(ResponseCode::Ok),
            _ => Ok(ResponseCode::Ok)
        }
    }
}

impl From<&ResponseCode> for u16 {
    fn from(code: &ResponseCode) -> Self {
        match code {
            ResponseCode::Ok => RESPONSE_CODE_OK,
            _ => { RESPONSE_CODE_OK }
        }
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
