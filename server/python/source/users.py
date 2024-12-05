from socket import socket
import time

from source.wire_formatting import read_string


users = {}


class User:
    def __init__(self, username: str) -> None:
        self.username = username
        self.lastlogin = None
        self.isonline = False
        self.conn = None
        self.messages = []

    def __str__(self) -> str:
        if self.lastlogin:
            from datetime import datetime

            format = "%d-%m-%Y %H:%M:%S"
            lastlogin = datetime.fromtimestamp(self.lastlogin).strftime(format)
        else:
            lastlogin = ""
        return f"\nuser: {self.username}\nlast login: {lastlogin}\nmessages to recv: {len(self.messages)}"

    def update_lastlogin(self):
        self.lastlogin = time.time()


def check_user(username: str) -> User:
    return users.setdefault(username, User(username))


def login(buffer: bytes, offset: int, conn: socket) -> tuple:
    username, _ = read_string(buffer, offset)
    user = check_user(username)
    if user.isonline:
        raise ValueError(f"user: {username} already logged!")
    else:
        user.isonline = True
        user.conn = conn
        user.update_lastlogin()

        return username


def logout(username: str) -> None:
    users[username].isonline = False
