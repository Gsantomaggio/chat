import socket
import struct
import time


class TcpClient:
    def __init__(self, address):
        self.address = address
        self.sock = socket.socket()
        self.correlation_id = 0

    def connect(self):
        self.sock.connect(self.address)

    def close(self):
        self.sock.close()

    def _send(self, data):
        length = struct.pack("!I", len(data))
        self.sock.sendall(length + data)

    def _receive(self):
        length_data = self.sock.recv(4)
        length = struct.unpack("!I", length_data)[0]
        return self.sock.recv(length)

    def login(self, username):
        self.correlation_id += 1
        header = struct.pack("!B H I", 0x01, 0x01, self.correlation_id)
        username_data = username.encode("utf-8")
        username_length = struct.pack("!H", len(username_data))
        message = header + username_length + username_data
        self._send(message)
        response = self._receive()
        return self._parse_response(response)

    def send_message(self, message, to_user, from_user):
        self.correlation_id += 1
        header = struct.pack("!B H I", 0x01, 0x02, self.correlation_id)
        message_data = message.encode("utf-8")
        message_length = struct.pack("!H", len(message_data))
        to_user_data = to_user.encode("utf-8")
        to_user_length = struct.pack("!H", len(to_user_data))
        from_user_data = from_user.encode("utf-8")
        from_user_length = struct.pack("!H", len(from_user_data))
        timestamp = struct.pack("!Q", int(time.time()))
        message = (
            header
            + message_length
            + message_data
            + from_user_length
            + from_user_data
            + to_user_length
            + to_user_data
            + timestamp
        )
        self._send(message)
        response = self._receive()
        return self._parse_response(response)

    def _parse_response(self, data):
        version, key, correlation_id, code = struct.unpack("!B H I H", data[:9])
        return {
            "version": version,
            "key": key,
            "correlation_id": correlation_id,
            "code": code,
        }
