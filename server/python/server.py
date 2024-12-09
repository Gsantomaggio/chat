import socket
from threading import Thread, active_count

from source.handle_client_message import read_message, send_user_messages, send_message
from source.users import logout

##############################################
# send_user_messages to go client doesn't work
##############################################

def handle_client_connection(conn, addr):
    with conn:
        print(f"Connected by {addr[0]}:{addr[1]}")
        u = None
        data = conn.recv(2048)
        while data:
            try:
                result, command = read_message(data, conn, u)
                if command == "CommandLogin":
                    u = result
                    send_user_messages(u)
                elif result and command == "CommandMessage":
                    send_message(result)
                data = conn.recv(2048)
            except socket.error:
                break
            except (Exception, ValueError) as e:
                print(f"Generic error...\n{e}")
                break
        logout(u)


def accept_connections(sock: socket.socket) -> None:
    conn, addr = sock.accept()
    thread = Thread(target=handle_client_connection, args=(conn, addr))
    thread.start()
    # print(f"Active connections: {active_count()-1}")


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
    except KeyboardInterrupt:
        print()


if __name__ == "__main__":
    run_server(port=5555)
