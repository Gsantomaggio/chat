use std::error::Error;
use bytes::BytesMut;
// use tokio::net::{TcpListener, TcpStream};
use tokio::io::{AsyncBufReadExt, AsyncReadExt, AsyncWriteExt, BufReader};
use tokio::net::{TcpListener, TcpStream};

#[derive(PartialEq, Eq, Debug)]
pub struct TcpServer {}

impl TcpServer {
    pub fn new() -> TcpServer {
        TcpServer {}
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

    loop {
        let mut buffer = BytesMut::new();
        let mut len_bytes = [0u8; 4];
        reader.read_exact(&mut len_bytes).await?;

        let len = u32::from_be_bytes(len_bytes) as usize;

        // Read the data.  This will wait until the specified number of bytes is available.
        buffer.resize(len, 0);
        reader.read_exact(&mut buffer).await?;

        let data = buffer.freeze();
        let message = String::from_utf8_lossy(&data);

        println!("Received: {}", message);

        // Respond to the client (optional)
        writer.write_all(b"Data received\r\n").await?;
        writer.flush().await?;

        // Check for connection closure.  If the read_exact fails, it indicates a closed connection.
        // You might handle this with a `break` or other error handling strategy.
    }
}