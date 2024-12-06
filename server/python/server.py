import socket
from threading import Thread, active_count

from source.handle_client_message import read_message, send_user_messages
from source.users import check_user


def handle_client_connection(conn, addr):
    with conn:
        print(f"Connected by {addr[0]}:{addr[1]}")
        u = None
        while True:
            data = conn.recv(2048)
            if not data:
                print("Connection closed")
                break
            else:
                try:
                    result, command = read_message(data, conn, u)
                    if command == "CommandLogin":
                        u = result
                    elif result and command == "CommandMessage":
                        user = check_user(result.to_field)
                        user.messages.append(result)
                        if user.isonline:
                            send_user_messages(user)
                except ValueError as verr:
                    conn.send(str(verr).encode())
                    break
                except Exception as e:
                    print(e)
                    break


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
                """
                import time
                time.sleep(1)
                if active_count() == 1:
                    print("Chiusura beccata correttamente...")
                    break
                """
    except socket.error as err:
        print(f"Socket Error: {err} | exiting...")
    except Exception:
        print("Generic exception, exiting...")


if __name__ == "__main__":
    run_server()
