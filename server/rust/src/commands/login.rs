use std::io::Write;
use crate::codec::{Decoder, Encoder};
use crate::commands::Command;
use crate::error::{DecodeError, EncodeError};
use crate::response::response::ResponseCode;
use crate::types::version::COMMAND_LOGIN;

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


impl Encoder for LoginRequest {
    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        self.correlation_id.encode(writer)?;
        self.user_name.as_str().encode(writer)?;
        Ok(())
    }

    fn encoded_size(&self) -> u32 {
        self.correlation_id.encoded_size()
            + self.user_name.as_str().encoded_size()
    }
}
impl Decoder for LoginRequest {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError> {
        let (input, correlation_id) = u32::decode(input)?;
        let (input, opt_user_name) =    Option::decode(input)?;

        Ok((
            input,
            LoginRequest {
                correlation_id,
                user_name: opt_user_name.unwrap_or("".to_string()),
            },
        ))
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
    pub(crate) response_code: ResponseCode,
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

