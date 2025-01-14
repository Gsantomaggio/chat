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
                    int messageLength = await ReadStreamLengthAsync(4, stream);
                    byte[] buffer = await ReadMessageStreamAsync(messageLength, stream);

                    await using var memoryStream = new MemoryStream(buffer);
                    using var reader = new BinaryReader(memoryStream, encoding, leaveOpen: true);

                    MessageHandler messageHandler = new(stream, reader);
                    usr = await messageHandler.HandleMessage(usr);
                }
            }
            catch (SocketException)
            {
                if (usr != null)
                {
                    Users.Logout(usr);
                    Logger.LogInformation("User {username} logged out", usr.Username);
                }
            }
            catch (Exception ex)
            {
                Logger.LogError("Error message: {errorMessage}", ex.Message);
            }
            finally
            {
                client.Close();
            }
        }

        /// <summary>
        /// Asynchronously reads a specified number of bytes from a NetworkStream and converts it to an integer.
        /// </summary>
        /// <param name="length">The number of bytes to read from the stream.</param>
        /// <param name="stream">The NetworkStream to read from.</param>
        /// <returns>The integer value represented by the bytes read from the stream, or 0 if the length is 0.</returns>
        private static async Task<int> ReadStreamLengthAsync(uint length, NetworkStream stream)
        {
            if (length == 0)
                return 0;
            byte[] buffer = new byte[length];
            await stream.ReadAsync(buffer);
            if (BitConverter.IsLittleEndian)
                Array.Reverse(buffer);
            int messageLength = BitConverter.ToInt32(buffer);

            return messageLength;
        }

        /// <summary>
        /// Asynchronously reads a specified number of bytes from a NetworkStream and returns them as a byte array.
        /// Throws a SocketException if the length is zero or if the number of bytes read does not match the specified length.
        /// </summary>
        /// <param name="length">The number of bytes to read from the stream.</param>
        /// <param name="stream">The NetworkStream to read from.</param>
        /// <returns>A byte array containing the bytes read from the stream.</returns>
        /// <exception cref="SocketException">Thrown if the length is zero or if the number of bytes read does not match the specified length.</exception>
        private static async Task<byte[]> ReadMessageStreamAsync(int length, NetworkStream stream)
        {
            if (length == 0)
                throw new SocketException((int)SocketError.ConnectionReset);

            byte[] buffer = new byte[length];
            int bytesRead = await stream.ReadAsync(buffer);

            if (bytesRead != length)
                throw new SocketException((int)SocketError.ConnectionReset);

            return buffer;
        }
    }
}
