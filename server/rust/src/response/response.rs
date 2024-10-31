use crate::codec::Decoder;
use crate::codec::decoder::read_u16;
use crate::commands::login::LoginResponse;
use crate::types::Header;

#[derive(Debug, PartialEq, Eq)]
pub struct Response {
    header: Header,
    pub(crate) kind: ResponseKind,
}

#[derive(Debug, PartialEq, Eq, Clone)]
pub enum ResponseCode {
    Ok,
    UserNotFound,
    UserAlreadyLogged,
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

}
