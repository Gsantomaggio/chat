import socket
from threading import Thread, active_count

from source.message import Message


def handle_client_message(buffer: bytes, user: str) -> Message:
    return Message(buffer, user)


def handle_client_connection(conn, addr):
    with conn:
        print(f"Connected by {addr[0]}:{addr[1]}")
        user = ""
        while True:
            data = conn.recv(1024)
            m = handle_client_message(data, user)
            print(m.username)
            conn.sendall(m)
    print(f"Active connections: {active_count()-1}")


def accept_connections(sock: socket.socket) -> None:
    conn, addr = sock.accept()
    thread = Thread(target=handle_client_connection, args=(conn, addr))
    thread.start()
    print(f"Active connections: {active_count()-1}")


def run_server(host: str = "0.0.0.0", port: int = 0, backlog: int = 5) -> None:
    try:
        with socket.socket() as s:
            s.bind((host, port))
            s.listen(backlog)
            host, port = s.getsockname()
            print(f"Server listening on address: {host}:{port}")
            while True:
                accept_connections(s)
    except socket.error as err:
        print(f"Socket Error: {err} | exiting...")
    except Exception:
        print("Generic exception, exiting...")


if __name__ == "__main__":
    run_server()
