use rust::tcp_server::TcpServer;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("Hello, world!");
    let tcp = TcpServer::new();
    tcp.start().await?;
    Ok(())
}
