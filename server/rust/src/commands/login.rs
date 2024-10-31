use std::io::Write;
use crate::codec::Encoder;
use crate::commands::Command;
use crate::error::EncodeError;
use crate::response::response::ResponseCode;
use crate::types::version::COMMAND_LOGIN;

#[derive(PartialEq, Eq, Debug)]
pub struct LoginRequest {
    correlation_id: u32,
    closing_code: ResponseCode,
}


impl LoginRequest {
    pub fn new(correlation_id: u32, closing_code: ResponseCode) -> Self {
        Self {
            correlation_id,
            closing_code,
        }
    }
}


impl Encoder for LoginRequest {
    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        self.correlation_id.encode(writer)?;
        self.closing_code.encode(writer)?;
        Ok(())
    }

    fn encoded_size(&self) -> u32 {
        self.correlation_id.encoded_size()
            + self.closing_code.encoded_size()
    }
}

impl Command for LoginRequest {
    fn key(&self) -> u16 {
        return COMMAND_LOGIN;
    }
}


#[derive(PartialEq, Eq, Debug)]
pub struct LoginResponse {
    pub(crate) correlation_id: u32,
    response_code: ResponseCode,
}

impl LoginResponse {
    pub fn new(correlation_id: u32, response_code: ResponseCode) -> Self {
        Self {
            correlation_id,
            response_code,
        }
    }
    pub fn is_ok(&self) -> bool {
        self.response_code == ResponseCode::Ok
    }
}

