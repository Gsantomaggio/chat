/// <summary>
/// The ServerTCP class is responsible for managing a TCP server that handles client connections,
/// processes incoming messages, and logs user statuses.
/// </summary>

using System.Net.Sockets;
using System.Net;
using System.Text;
using server.src;

namespace server
{
    internal static class ServerTCP
    {
        private static readonly Encoding encoding = Encoding.UTF8;
        private static readonly CancellationTokenSource _cancellationTokenSource = new();
        private const int PORT = 5555;
        private static readonly TcpListener _server = new(IPAddress.Any, PORT);

        /// <summary>
        /// The main entry point for the server application. Starts the server, handles client connections,
        /// waits for a stop signal, and logs user statuses.
        /// </summary>
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

        /// <summary>
        /// Logs the status of connected users at regular intervals until the server is stopped.
        /// </summary>
        private static async Task LogUsersStatus()
        {
            await Task.Run(() =>
            {
                while (!_cancellationTokenSource.IsCancellationRequested)
                {
                    string usersToPrint = string.Join("\n\t", Users.Instance.GetUsers()
                        .Select(u => $"{u.Username} is {u.Status}, last login: {u.LastLogin} UTC"));
                    
                    usersToPrint = usersToPrint != string.Empty ? $"Users:\n\t{usersToPrint}" : "Users: []\n";
                    
                    Logger.LogDebug("{message}", usersToPrint);
                    Thread.Sleep(3000);
                }
            });
        }

        /// <summary>
        /// Waits for a key press to stop the server.
        /// </summary>
        private static async Task WaitStopServer()
        {
            Logger.LogInformation("Press any key to stop the server");
            await Task.Run(() => Console.ReadKey());
        }

        /// <summary>
        /// Handles incoming client connections asynchronously.
        /// </summary>
        private static async Task HandleClientsAsync()
        {
            try
            {
                while (!_cancellationTokenSource.IsCancellationRequested)
                {
                    TcpClient client = await _server.AcceptTcpClientAsync();
                    if (client.Client.RemoteEndPoint is null)
                    {
                        client.Close();
                        continue;
                    }
                    Logger.LogInformation("Client connected from {Endpoint}", client.Client.RemoteEndPoint);
                    _ = Task.Run(() => HandleClientAsync(client));
                }
            }
            catch
            {
                Logger.LogInformation("Server stopped...");
            }
        }

        /// <summary>
        /// Processes messages from a connected client and handles user login/logout.
        /// </summary>
        /// <param name="client">The connected client.</param>
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
