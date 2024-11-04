use std::error::Error;
use std::ptr::write;
use byteorder::BigEndian;
use bytes::BytesMut;
// use tokio::net::{TcpListener, TcpStream};
use tokio::io::{AsyncBufReadExt, AsyncReadExt, AsyncWriteExt, BufReader, BufWriter, WriteHalf};
use tokio::net::{TcpListener, TcpStream};
use crate::codec::{Decoder, Encoder};
use crate::commands::login::{LoginRequest, LoginResponse};
use crate::response::response::{Response, ResponseCode};
use crate::response::response::ResponseKind::Login;
use crate::types::Header;
use crate::types::version::{COMMAND_LOGIN, GENERIC_RESPONSE, PROTOCOL_VERSION};

#[derive(PartialEq, Eq, Debug)]
pub struct User {
    pub(crate) user_name: String,
}

#[derive(PartialEq, Eq, Debug)]
pub struct TcpServer {
    pub(crate) users: Vec<User>,
}

impl TcpServer {
    pub fn new() -> TcpServer {
        TcpServer {
            users: Vec::new()
        }
    }

    pub async fn start(&self) -> std::io::Result<()> {
        let listener = TcpListener::bind("127.0.0.1:5555").await?;
        println!("Server listening on 127.0.0.1:5555");
        loop {
            let (mut socket, _) = listener.accept().await?;
            tokio::spawn(async move {
                if let Err(e) = handle_client(socket).await {
                    println!("Error handling client: {}", e);
                }
            });
        }
    }
}


async fn handle_client(mut socket: TcpStream) -> Result<(), Box<dyn Error>> {
    let (reader, mut writer) = socket.split();
    let mut reader = BufReader::new(reader);
    // let mut writer = BufWriter::new(writer);


    loop {
        let mut buffer = BytesMut::new();
        let mut len_bytes = [0u8; 4];
        reader.read_exact(&mut len_bytes).await?;

        let len = u32::from_be_bytes(len_bytes) as usize;

        // Read the data.  This will wait until the specified number of bytes is available.
        buffer.resize(len, 0);
        reader.read_exact(&mut buffer).await?;

        let (h, hq) = Header::decode(&buffer).unwrap();
        match hq.key()
        {
            COMMAND_LOGIN => {
                let (_, login) = LoginRequest::decode(h).unwrap();
                println!("Logged{}", login.user_name);
                users
                let response = Response::new(
                    Header::new(PROTOCOL_VERSION, GENERIC_RESPONSE),
                    Login(LoginResponse::new(login.correlation_id, ResponseCode::Ok)),
                );
                writer.write_all(response_buffer(&response).await.unwrap().as_slice()).await?;
            }
            _ => {
                println!("error");
            }
        }
        writer.flush().await?;
    }
}

async fn response_buffer(response: &Response) -> Result<Vec<u8>, Box<dyn Error>> {
    let mut writer_tmp = Vec::new();
    writer_tmp.write_u32(response.encoded_size()).await?;
    response.encode(&mut writer_tmp).unwrap();
    writer_tmp.flush().await?;
    Ok(writer_tmp)
}