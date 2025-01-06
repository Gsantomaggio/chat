using System.Collections.Concurrent;
using System.Net.Sockets;

namespace server.src
{
    internal class SingleUser(string username)
    {
        public string Username { get; set; } = username;

        public NetworkStream? Stream { get; set; } = null;
        public string? Status { get; set; } = "offline";
        public DateTime? LastLogin { get; set; } = null;
        public Queue<Message> Messages { get; set; } = new Queue<Message>();
    }

    internal sealed class Users
    {
        private static readonly Lazy<Users> lazy = new(() => new Users());
        public static Users Instance => lazy.Value;
        private readonly ConcurrentDictionary<string, SingleUser> users = new();

        private Users() { }

        public bool TryGetValue(string key, out SingleUser user) => users.TryGetValue(key, out user);
        public SingleUser GetOrAdd(string key, SingleUser user) => users.GetOrAdd(key, user);


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
