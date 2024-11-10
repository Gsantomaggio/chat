use std::io::Write;
use crate::types::{ DecodeError, Decoder, EncodeError, Encoder};

// 	return readMany(reader, &m.correlationId, &m.Message, &m.From, &m.To, &m.Time)
pub struct CommandMessage {
    pub correlation_id: u32,
    pub message: String, // payload would be better as name
    pub from: String,
    pub to: String,
    pub time: u64,
}


// 	return readMany(reader, &m.correlationId, &m.Message, &m.From, &m.To, &m.Time)
impl Decoder for CommandMessage {
    fn decode(input: &[u8]) -> Result<(&[u8], Self), DecodeError> {
        let (input, correlation_id) = u32::decode(input)?;
        let (input, message1) = Option::decode(input)?;
        let (input, from1) =  Option::decode(input)?;
        let (input, to1) =  Option::decode(input)?;
        let (input, time) = u64::decode(input)?;
        Ok((input, CommandMessage {
            correlation_id,
            message: message1.unwrap_or("".to_string()),
            from:from1.unwrap(),
            to: to1.unwrap(),
            time,
        }))
    }
}

impl CommandMessage {
    pub fn new(correlation_id: u32, message: String, from: String, to: String, time: u64) -> Self {
        Self {
            correlation_id,
            message,
            from,
            to,
            time,
        }
    }
}

// impl Command for CommandMessage {
//     fn key(&self) -> u16 {
//         return COMMAND_MESSAGE;
//     }
// }

impl Encoder for CommandMessage {
    fn encoded_size(&self) -> u32 {
        self.correlation_id.encoded_size()
            + self.message.as_str().encoded_size()
            + self.from.as_str().encoded_size()
            + self.to.as_str().encoded_size()
            + self.time.encoded_size()
    }

    fn encode(&self, writer: &mut impl Write) -> Result<(), EncodeError> {
        self.correlation_id.encode(writer)?;
        self.message.as_str().encode(writer)?;
        self.from.as_str().encode(writer)?;
        self.to.as_str().encode(writer)?;
        self.time.encode(writer)?;
        Ok(())
    }
}