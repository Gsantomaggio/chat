import time
import socket
from threading import Thread, Event

from source.message_handler import MessageHandler
from source.wire_formatting import read_uint32
from source.users import logout

from source import Logger
from source.exceptions import AlreadyLoggedException

logger = Logger(__name__)

"""
TcpServer is a class that sets up a TCP server to handle multiple client connections.

Attributes:
    host (str): The host address for the server.
    port (int): The port number for the server.
    backlog (int): The maximum number of queued connections.
    stop_event (Event): An event to signal the server to stop.
    users (dict): A dictionary to store user information.
    server_thread (Thread): A thread to run the server.

Methods:
    log_users_status(): Logs the status of connected users periodically.
    handle_client_connection(conn, addr): Handles the connection with a client.
    accept_connections(sock): Accepts incoming client connections.
    run_server(): Runs the server, accepting connections and handling them.
    stop_server(): Stops the server when the user presses ENTER.
"""


class TcpServer:
    def __init__(self, host="0.0.0.0", port=0, backlog=5):
        self.host = host
        self.port = port
        self.backlog = backlog
        self.stop_event = Event()
        self.users = {}
        self.server_thread = Thread(target=self.run_server)
        self.server_thread.start()
        self.stop_server()

    def log_users_status(self):
        def print_users_status(self):
            while not self.stop_event.is_set():
                if not self.users:
                    users_to_print = "Users: []\n"
                else:
                    users_to_print = "Users:\n\t"
                    for user in self.users.values():
                        single_user_to_print = f"{user.username} is {user.status}, last login: {user.printlastlogin()} UTC\n\t"
                        users_to_print += single_user_to_print
                logger.debug(users_to_print)
                time.sleep(3)

        thread = Thread(target=print_users_status, args=(self,))
        thread.start()

    def handle_client_connection(self, conn, addr):
        conn.settimeout(1)
        with conn:
            client_refs = f"{addr[0]}:{addr[1]}"
            logger.info(f"Connected by {client_refs}")
            usr = None
            while not self.stop_event.is_set():
                try:
                    data_length = conn.recv(4)
                    data_message = None
                    if data_length:
                        message_length, _ = read_uint32(data_length, 0)
                        data_message = conn.recv(message_length)
                    if data_message:
                        usr = MessageHandler(conn, self.users).read_message(
                            data_message, usr
                        )
                    else:
                        logout(usr)
                        logger.info(f"User {usr.username} logged out")
                        break
                except socket.timeout:
                    continue
                except AlreadyLoggedException as e:
                    logger.warning(e)
                    break
                except ValueError as e:
                    logger.warning(f"{e}\nClosing connection with {client_refs}")

    def accept_connections(self, sock: socket.socket):
        while not self.stop_event.is_set():
            try:
                conn, addr = sock.accept()
                thread = Thread(target=self.handle_client_connection, args=(conn, addr))
                thread.start()
            except socket.timeout:
                continue

    def run_server(self):
        try:
            with socket.socket() as s:
                s.bind((self.host, self.port))
                s.listen(self.backlog)
                s.settimeout(1)
                host, port = s.getsockname()
                logger.info(f"Server listening on address: {host}:{port}")
                self.log_users_status()
                self.accept_connections(s)
        except socket.error as err:
            logger.error(f"Socket Error: {err} | exiting...")
        except Exception as e:
            logger.critical(f"Generic exception, exiting...\n{e}")

    def stop_server(self):
        input("Press ENTER to stop\n")
        self.stop_event.set()
        self.server_thread.join()
