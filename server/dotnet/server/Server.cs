using Microsoft.Extensions.Logging;
using System.Net.Sockets;
using System.Net;
using System.Text;
using server.src;

namespace server
{


    namespace server
    {
        internal class ServerTCP
        {
            private static readonly Encoding encoding = Encoding.UTF8;
            private static readonly CancellationTokenSource _cancellationTokenSource = new();
            private const int PORT = 5555;
            private static readonly TcpListener _server = new(IPAddress.Any, PORT);
            private static readonly ILogger _logger = LoggerFactory.Create(builder => builder.AddConsole()).CreateLogger(nameof(ServerTCP));

            public static async Task Main()
            {
                _logger.LogInformation("Server started on port {Port}", PORT);

                _server.Start();

                var handleClientTask = HandleClientsAsync();
                var waitStopServerTask = WaitStopServer();

                await waitStopServerTask;

                _cancellationTokenSource.Cancel();
                _server.Stop();

                await handleClientTask;
            }

            private static async Task WaitStopServer()
            {
                _logger.LogInformation("Press any key to stop the server");
                await Task.Run(() => Console.ReadKey());
            }

            private static async Task HandleClientsAsync()
            {
                try
                {
                    while (!_cancellationTokenSource.IsCancellationRequested)
                    {
                        TcpClient client = await _server.AcceptTcpClientAsync();
                        _logger.LogInformation("Client connected from {Endpoint}", client.Client.RemoteEndPoint);
                        _ = Task.Run(() => HandleClientAsync(client));
                    }
                }
                catch (Exception ex)
                {
                    _logger.LogInformation("Server stopped: {Message}", ex.Message);
                }
            }

            private static async Task HandleClientAsync(TcpClient client)
            {
                NetworkStream stream = client.GetStream();
                SingleUser? usr = null;

                try
                {
                    while (true)
                    {
                        byte[] buffer = new byte[1024];
                        int bytesRead = await stream.ReadAsync(buffer);
                        if (bytesRead == 0)
                        {
                            Users.Logout(usr);
                            _logger.LogInformation("Connection closed by {Endpoint}", client.Client.RemoteEndPoint);
                            break;
                        }

                        await using var memoryStream = new MemoryStream(buffer, 0, bytesRead);
                        using var reader = new BinaryReader(memoryStream, encoding, leaveOpen: true);

                        MessageHandler messageHandler = new(stream, reader, _logger);
                        usr = await messageHandler.HandleMessage(usr);
                    }
                }
                finally
                {
                    client.Close();
                }
            }
        }
    }
}
