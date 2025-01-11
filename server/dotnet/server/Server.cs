using Microsoft.Extensions.Logging;
using System.Net.Sockets;
using System.Net;
using System.Text;
using server.src;

namespace server
{
    namespace server
    {
        internal static class ServerTCP
        {
            private static readonly Encoding encoding = Encoding.UTF8;
            private static readonly CancellationTokenSource _cancellationTokenSource = new();
            private const int PORT = 5555;
            private static readonly TcpListener _server = new(IPAddress.Any, PORT);

            public static async Task Main()
            {
                Logger.LogInformation("Server started on port {Port}", PORT);

                _server.Start();

                var handleClientTask = HandleClientsAsync();
                var waitStopServerTask = WaitStopServer();
                var logUsersStatus = LogUsersStatus();

                await waitStopServerTask;

                _cancellationTokenSource.Cancel();
                _server.Stop();

                await logUsersStatus;
                await handleClientTask;
            }

            private static async Task LogUsersStatus()
            {
                await Task.Run(() =>
                {
                    while (!_cancellationTokenSource.IsCancellationRequested)
                    {
                        string usersToPrint = "Users:\n\t";
                        if (Users.Instance.Length() < 1)
                        {
                            usersToPrint = "Users: []\n";
                        }
                        else
                        {
                            foreach (var u in Users.Instance.GetUsers())
                            {
                                usersToPrint += $"{u.Username} is {u.Status}, last login: {u.LastLogin} UTC\n\t";
                            }
                        }
                        Logger.LogDebug("{message}", usersToPrint);
                        Thread.Sleep(3000);
                    }
                });
            }

            private static async Task WaitStopServer()
            {
                Logger.LogInformation("Press any key to stop the server");
                await Task.Run(() => Console.ReadKey());
            }

            private static async Task HandleClientsAsync()
            {
                try
                {
                    while (!_cancellationTokenSource.IsCancellationRequested)
                    {
                        TcpClient client = await _server.AcceptTcpClientAsync();
                        Logger.LogInformation("Client connected from {Endpoint}", client.Client.RemoteEndPoint);
                        _ = Task.Run(() => HandleClientAsync(client));
                    }
                }
                catch
                {
                    Logger.LogInformation("Server stopped...");
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
                            if (usr != null)
                            {
                                Users.Logout(usr);
                                Logger.LogInformation("User {username} logged out", usr.Username);
                            }
                            break;
                        }

                        await using var memoryStream = new MemoryStream(buffer, 0, bytesRead);
                        using var reader = new BinaryReader(memoryStream, encoding, leaveOpen: true);

                        MessageHandler messageHandler = new(stream, reader);
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
