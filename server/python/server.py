import socket
from threading import Thread

from source.handle_client_message import read_message, send_user_messages, send_message
from source.users import logout

running = True
def handle_client_connection(conn, addr):
    global running
    conn.settimeout(1)
    with conn:
        print(f"Connected by {addr[0]}:{addr[1]}")
        u = None
        while running:
            try:
                data = conn.recv(2048)
                if data:
                    result, command = read_message(data, conn, u)
                    if command == "CommandLogin":
                        u = result
                        send_user_messages(u)
                    elif result and command == "CommandMessage":
                        send_message(result)
                else:
                    logout(u)
                    break
            except socket.timeout:
                continue


def accept_connections(sock: socket.socket) -> None:
    global running
    while running:
        try:
            conn, addr = sock.accept()
            thread = Thread(target=handle_client_connection, args=(conn, addr))
            thread.start()
        except socket.timeout:
            continue


def run_server(host: str = "0.0.0.0", port: int = 0, backlog: int = 5) -> None:
    global running
    try:
        with socket.socket() as s:
            s.bind((host, port))
            s.listen(backlog)
            s.settimeout(1)
            host, port = s.getsockname()
            print(f"Server listening on address: {host}:{port}")
            accept_connections(s)
    except socket.error as err:
        print(f"Socket Error: {err} | exiting...")
    except Exception as e:
        print(f"Generic exception, exiting...\n{e}")


def stop():
    global running
    input("Press ENTER to stop\n")
    running = False


if __name__ == "__main__":
    server_thread = Thread(target=run_server, kwargs={'port': 5555})
    server_thread.start()
    stop()
