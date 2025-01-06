from socket import socket
import time
from datetime import datetime


class User:
    """
    User class represents a user in the system with attributes and methods to manage their status and connection.

    Attributes:
        username (str): The username of the user.
        lastlogin (float): The timestamp of the user's last login.
        isonline (bool): The online status of the user.
        status (str): The current status of the user (online/offline).
        conn (socket): The socket connection associated with the user.
        messages (list): A list to store messages for the user.

    Methods:
        __str__(): Returns a string representation of the user.
        login(conn: socket) -> int: Logs the user in and updates their status and connection.
        printlastlogin(): Returns the last login time as a formatted string.
    """

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
        return "Never"


def logout(user: User | None):
    """
    Function:
    logout(user: User | None): Logs out the user by updating their status and connection.
    """

    if user:
        user.isonline = False
        user.status = "offline"
        user.conn = None
