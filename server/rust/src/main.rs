/*
 *   Copyright (c) 2024 Nazmul Idris
 *   All rights reserved.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */
mod login;
mod message;
mod response;
mod types;

use crate::login::LoginResponse;
use crate::message::CommandMessage;
use crate::response::ResponseKind::Login;
use crate::response::{Response, ResponseCode};
use crate::types::version::{COMMAND_LOGIN, COMMAND_MESSAGE, GENERIC_RESPONSE, PROTOCOL_VERSION};
use crate::types::{Decoder, Encoder, Header, LoginRequest};
use bytes::BytesMut;
use dashmap::DashMap;
use log::{error, info};
use std::error::Error;
use std::sync::Arc;
use tokio::io::AsyncReadExt;
use tokio::sync::Mutex;
use tokio::{
    io::AsyncWriteExt,
    net::{TcpListener, TcpStream},
    // sync::broadcast::{self, Sender},
};

type IOResult<T> = std::io::Result<T>;

#[derive(Debug, Clone)]
pub struct Message {
    pub user_name: String,
    pub payload: String,
}

struct TokioConnection {
    reader: Mutex<tokio::io::ReadHalf<TcpStream>>,
    writer: Mutex<tokio::io::WriteHalf<TcpStream>>,
}
#[tokio::main]
async fn main() -> IOResult<()> {
    let addr = "0.0.0.0:5555";
    // Start logging.
    femme::start();

    let listener = TcpListener::bind(addr).await?;

    let connections: Arc<DashMap<String, Arc<TokioConnection>>> = Arc::new(DashMap::new());

    // prints connected users every 5 seconds
    let c = Arc::clone(&connections);
    tokio::spawn(async move {
        loop {
            tokio::time::sleep(tokio::time::Duration::from_secs(5)).await;
            info!(
                "Connected users: {:?}",
                &c.iter().map(|x| x.key().clone()).collect::<Vec<_>>()
            );
        }
    });

    // Server infinite loop.
    loop {
        info!("Listening for new connections, address: {}", addr);

        // Accept incoming connections.
        let (client_tcp_stream, _) = listener.accept().await?;
        let connections: Arc<DashMap<String, Arc<TokioConnection>>> = Arc::clone(&connections);

        tokio::spawn(async move {
            let c = Arc::clone(&connections);
            let user_name = match handle_client_task(client_tcp_stream, connections).await {
                Ok(user_name) => {
                    info!(
                        "({}) - Successfully ended client task",
                        user_name.as_ref().unwrap_or(&"Anonymous".to_string())
                    );
                    user_name
                }
                Err(error) => {
                    info!("Problem handling client task: {:?}", error);
                    error.user_name
                }
            };
            if let Some(user_name) = user_name {
                c.remove(&user_name);
                info!("{} Disconnected", user_name);
            }
        });

        info!("Released a connection");
    }
}

#[derive(Debug)]
struct UserConnectionError {
    user_name: Option<String>,
    inner: std::io::Error,
}

impl UserConnectionError {
    fn new(user_name: Option<String>, inner: std::io::Error) -> Self {
        UserConnectionError { user_name, inner }
    }
}

impl std::fmt::Display for UserConnectionError {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(
            f,
            "({}) - ConnectionError: {:?}",
            self.user_name.as_ref().unwrap_or(&"Anonymous".to_string()),
            self.inner
        )
    }
}

async fn handle_client_task(
    client_tcp_stream: TcpStream,
    connections: Arc<DashMap<String, Arc<TokioConnection>>>,
) -> Result<Option<String>, UserConnectionError> {
    // Get reader and writer from tcp stream.
    // let (reader, writer) = tokio::io::split(client_tcp_stream);
    let mut user_name = None;

    let (reader, writer) = tokio::io::split(client_tcp_stream);
    let connection = Arc::new(TokioConnection {
        reader: Mutex::new(reader),
        writer: Mutex::new(writer),
    });
    loop {
        let mut locked_reader = connection.reader.lock().await;
        let mut buffer = BytesMut::new();
        let mut len_bytes = [0u8; 4];
        locked_reader
            .read_exact(&mut len_bytes)
            .await
            .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
        let len = u32::from_be_bytes(len_bytes) as usize;

        // Read the data.  This will wait until the specified number of bytes is available.
        buffer.resize(len, 0);
        locked_reader
            .read_exact(&mut buffer)
            .await
            .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
        let (h, hq) = Header::decode(&buffer).unwrap();
        match hq.key() {
            COMMAND_LOGIN => {
                let (_, login) = LoginRequest::decode(h).unwrap();
                info!("Logged {}", login.user_name);
                let response = Response::new(
                    Header::new(PROTOCOL_VERSION, GENERIC_RESPONSE),
                    Login(LoginResponse::new(login.correlation_id, ResponseCode::Ok)),
                );
                let mut locked_writer = connection.writer.lock().await;
                locked_writer
                    .write_all(response_buffer(&response).await.unwrap().as_slice())
                    .await
                    .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
                locked_writer
                    .flush()
                    .await
                    .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
                info!("[{}] handle_client: start", login.user_name);
                info!("[END] Logged {}", login.user_name);
                let user = login.user_name.clone();
                user_name = Some(user.clone());
                connections.insert(user, Arc::clone(&connection));
            }

            COMMAND_MESSAGE => {
                let (_, command_message_request) = CommandMessage::decode(h).unwrap();
                let cloned_connections = Arc::clone(&connections);

                let connection_map = cloned_connections.get(&command_message_request.to);

                match connection_map {
                    None => {
                        let response = Response::new(
                            Header::new(PROTOCOL_VERSION, GENERIC_RESPONSE),
                            Login(LoginResponse::new(
                                command_message_request.correlation_id,
                                ResponseCode::UserNotFound,
                            )),
                        );
                        let mut locked_writer = connection.writer.lock().await;
                        locked_writer
                            .write_all(response_buffer(&response).await.unwrap().as_slice())
                            .await
                            .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
                        locked_writer
                            .flush()
                            .await
                            .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
                        drop(locked_writer);
                        return Ok(user_name.clone());
                    }
                    Some(ok_connection) => {
                        let response = Response::new(
                            Header::new(PROTOCOL_VERSION, GENERIC_RESPONSE),
                            Login(LoginResponse::new(
                                command_message_request.correlation_id,
                                ResponseCode::Ok,
                            )),
                        );
                        let mut locked_writer = connection.writer.lock().await;
                        locked_writer
                            .write_all(response_buffer(&response).await.unwrap().as_slice())
                            .await
                            .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
                        locked_writer
                            .flush()
                            .await
                            .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
                        drop(locked_writer);

                        let destination_header = Header::new(PROTOCOL_VERSION, COMMAND_MESSAGE);
                        let destination_command_message = CommandMessage::new(
                            0,
                            command_message_request.message,
                            command_message_request.from,
                            command_message_request.to,
                            command_message_request.time,
                        );

                        let mut locked_writer_destination =
                            ok_connection.value().writer.lock().await;

                        let mut writer_tmp = Vec::new();
                        writer_tmp
                            .write_u32(
                                destination_header.encoded_size()
                                    + destination_command_message.encoded_size(),
                            )
                            .await
                            .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
                        destination_header.encode(&mut writer_tmp).unwrap();
                        destination_command_message.encode(&mut writer_tmp).unwrap();
                        match locked_writer_destination
                            .write_all(writer_tmp.as_slice())
                            .await
                        {
                            Ok(_) => info!(
                                "Sent message to the user {}",
                                destination_command_message.to
                            ),
                            Err(e) => error!("Error writing to connection {}: {:?}", "id", e),
                        }

                        locked_writer_destination
                            .flush()
                            .await
                            .map_err(|e| UserConnectionError::new(user_name.clone(), e))?;
                    }
                }
            }

            _ => {}
        }
    }
}

async fn response_buffer(response: &Response) -> Result<Vec<u8>, Box<dyn Error>> {
    let mut writer_tmp = Vec::new();
    writer_tmp.write_u32(response.encoded_size()).await?;
    response.encode(&mut writer_tmp).unwrap();
    writer_tmp.flush().await?;
    Ok(writer_tmp)
}
