using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Net.Sockets;

try
{
    TcpClient client = new("127.0.0.1", 5555);
    Console.WriteLine("Connected to server");

    NetworkStream stream = client.GetStream();

    while (true)
    {
        Console.Write("Enter a message: ");
        string message = Console.ReadLine() ?? string.Empty;
        if (message == "exit") break;

        byte[] buffer = Encoding.UTF8.GetBytes(message);
        stream.Write(buffer, 0, buffer.Length);

        byte[] response = new byte[2048];
        int bytesRead = stream.Read(response, 0, response.Length);
        string data = Encoding.UTF8.GetString(response, 0, bytesRead);
        Console.WriteLine($"Received: {data}");
    }
}
catch (Exception exec)
{
    Console.WriteLine($"Error: {exec}");
}
