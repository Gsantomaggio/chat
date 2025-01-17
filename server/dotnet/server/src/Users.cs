using System.Collections.Concurrent;
using System.Net.Sockets;

namespace server.src
{
    /// <summary>
    /// Represents a single user with properties for username, network stream, status, last login time, and a queue of messages.
    /// </summary>
    /// <param name="username">The username of the user.</param>
    internal class SingleUser(string username)
    {
        public string Username { get; set; } = username;

        public NetworkStream? Stream { get; set; } = null;
        public string? Status { get; set; } = "offline";
        public DateTime? LastLogin { get; set; } = null;
        public Queue<Message> Messages { get; set; } = new Queue<Message>();
    }

    /// <summary>
    /// Singleton class that manages a collection of users, providing methods for user login, logout, and retrieval.
    /// </summary>
    internal sealed class Users
    {
        /// <summary>
        /// Gets the singleton instance of the Users class.
        /// </summary>
        private static readonly Lazy<Users> lazy = new(() => new Users());
        public static Users Instance => lazy.Value;
        private readonly ConcurrentDictionary<string, SingleUser> users = new();

        private Users() { }

        /// <summary>
        /// Tries to get a user from the collection by their username.
        /// </summary>
        /// <param name="key">The username of the user to retrieve.</param>
        /// <param name="user">The retrieved user, if found.</param>
        /// <returns>True if the user was found; otherwise, false.</returns>
        public bool TryGetValue(string key, out SingleUser user) => users.TryGetValue(key, out user);

        /// <summary>
        /// Gets an existing user or adds a new user to the collection.
        /// </summary>
        /// <param name="key">The username of the user.</param>
        /// <param name="user">The user to add if not already present.</param>
        /// <returns>The existing or newly added user.</returns>
        public SingleUser GetOrAdd(string key, SingleUser user) => users.GetOrAdd(key, user);

        /// <summary>
        /// Gets an array of all users in the collection, ordered by username in ascending order.
        /// </summary>
        /// <returns>An array of users.</returns>
        public SingleUser[] GetUsers() => [.. users.Values.OrderBy(user => user.Username)];


        /// <summary>
        /// Logs in a user by setting their status to online and updating their last login time and stream.
        /// </summary>
        /// <param name="username">The username of the user to log in.</param>
        /// <param name="stream">The network stream associated with the user.</param>
        /// <param name="user">The logged-in user.</param>
        /// <returns>A response code indicating the result of the login attempt.</returns>
        public static ushort Login(string username, NetworkStream stream, out SingleUser user)
        {
            user = Instance.GetOrAdd(username, new SingleUser(username));

            if (user.Status == "online")
                return Constants.ResponseCodeErrorUserAlreadyLogged;

            user.Status = "online";
            user.LastLogin = DateTime.Now;
            user.Stream = stream;
            
            return Constants.ResponseCodeOk;
        }

        /// <summary>
        /// Logs out a user by setting its status to offline and clearing its stream.
        /// </summary>
        /// <param name="user">The user to log out.</param>
        public static void Logout(SingleUser? user)
        {
            if (user is not null)
            {
                user.Status = "offline";
                user.Stream = null;
            }

        }
    }
}
