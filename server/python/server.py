import socket
from threading import Thread, Event

from source.handle_client_message import read_message
from source.users import logout


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

    def handle_client_connection(self, conn, addr):
        conn.settimeout(1)
        with conn:
            client_refs = f"{addr[0]}:{addr[1]}"
            print(f"Connected by {client_refs}")
            usr = None
            while not self.stop_event.is_set():
                try:
                    data = conn.recv(2048)
                    if data:
                        usr = read_message(data, conn, usr, self.users)
                    else:
                        logout(usr)
                        break
                except socket.timeout:
                    continue
                except ValueError as e:
                    print(f"{e}\nClosing connection with {client_refs}")

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
                print(f"Server listening on address: {host}:{port}")
                self.accept_connections(s)
        except socket.error as err:
            print(f"Socket Error: {err} | exiting...")
        except Exception as e:
            print(f"Generic exception, exiting...\n{e}")

    def stop_server(self):
        input("Press ENTER to stop\n")
        self.stop_event.set()
        self.server_thread.join()


if __name__ == "__main__":
    server = TcpServer(port=5555)
