from socket import socket
import time

from source.wire_formatting import read_string


users = {}


class User:
    def __init__(self, username: str, conn: socket) -> None:
        self.username = username
        self.lastlogin = time.time()
        self.isonline = False
        self.conn = conn
        self.messages = []

    def update_lastlogin(self):
        self.lastlogin = time.time()


def login(buffer: bytes, offset: int, conn: socket) -> tuple:
    username, _ = read_string(buffer, offset)
    user: User = users.setdefault(username, User(username, conn))
    if user.isonline:
        raise ValueError(f"user: {username} already logged!")
    else:
        user.isonline = True
        user.update_lastlogin()

        return username


def logout(username: str) -> None:
    users[username].isonline = False
