from socket import socket
import time


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


def login(user: User, conn: socket) -> int:
    if user.isonline:
        return 4
    else:
        user.isonline = True
        user.conn = conn
        user.update_lastlogin()
        return 1


def logout(user: User | None) -> None:
    if user:
        user.isonline = False
        user.conn = None
