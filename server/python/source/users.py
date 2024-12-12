from socket import socket
import time
from datetime import datetime


class User:
    def __init__(self, username: str):
        self.username = username
        self.lastlogin = None
        self.isonline = False
        self.status = "offline"
        self.conn = None
        self.messages = []

    def __str__(self) -> str:
        return f"{self.username} - status: {self.status} - last login: {self.printlastlogin()} - messages in queue: {len(self.messages)}"

    def login(self, conn: socket) -> int:
        if self.isonline:
            return 4
        else:
            self.isonline = True
            self.status = "online"
            self.conn = conn
            self.lastlogin = time.time()
            return 1
    
    def printlastlogin(self):
        if self.lastlogin:
            format = "%a, %d %b %Y %H:%M:%S"
            return datetime.fromtimestamp(self.lastlogin).strftime(format)
        return ""
        

def logout(user: User | None):
    if user:
        user.isonline = False
        user.conn = None
