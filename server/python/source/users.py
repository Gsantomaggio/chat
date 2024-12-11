from socket import socket
import time
from datetime import datetime


class User:
    def __init__(self, username: str):
        self.username = username
        self.lastlogin = None
        self.isonline = False
        self.conn = None
        self.messages = []

    def __str__(self) -> str:
        if self.lastlogin:
            format = "%d-%m-%Y %H:%M:%S"
            lastlogin = datetime.fromtimestamp(self.lastlogin).strftime(format)
        else:
            lastlogin = ""
        return f"\nuser: {self.username}\nlast login: {lastlogin}\nmessages to recv: {len(self.messages)}"

    def login(self, conn: socket) -> int:
        if self.isonline:
            return 4
        else:
            self.isonline = True
            self.conn = conn
            self.lastlogin = time.time()
            return 1


def logout(user: User | None):
    if user:
        user.isonline = False
        user.conn = None
