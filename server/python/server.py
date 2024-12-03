import socket
from threading import Thread, active_count

from source.handle_client_message import read_message
from source.users import users, User


def handle_client_connection(conn, addr):
    with conn:
        print(f"Connected by {addr[0]}:{addr[1]}")
        is_logged = False
        while True:
            data = conn.recv(1024)
            try:
                result, command = read_message(data, conn, is_logged)
                if command == "CommandLogin":
                    is_logged = True
                    conn.send(f"user: {result} logged in correcly")
                elif command == "CommandMessage":
                    message = result
                    user: User = users[message.to_field]
                    user.messages.append()
                else:
                    print(f"Generic Error occurred. Command: {command}")
            except ValueError as verr:
                conn.send(verr)

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
