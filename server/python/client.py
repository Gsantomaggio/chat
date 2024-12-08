import socket
import sys

from source import protocol


def check_close_conn(s: socket, text: str) -> str:
    if text.upper() == "ESC":
        print("Closing connection to the server.")
        s.close()
        sys.exit()
    return text


def send_messages(s: socket) -> None:
    usr = None
    while s:
        if not usr:
            usr = check_close_conn(
                s,
                input(
                    "\nTo send messages you have to login.\nPlease enter your username (ESC to quit):\n-> "
                ),
            )
            message = protocol.create_login_message(usr)
        else:
            msg = check_close_conn(s, input("Enter the message (ESC to quit): "))
            to = check_close_conn(s, input("Enter the receiver (ESC to quit): "))
            message = protocol.create_msg_message(msg, to, usr)

        s.send(message)
        data = s.recv(4096)
        print(f"{str(data, "utf-8")}")


def conn_server(serverAddress: tuple) -> None:
    try:
        with socket.socket() as s:
            s.connect(serverAddress)
            print(f"Connected to the server:\t{serverAddress=}")
            send_messages(s)
    except socket.error as err:
        print(f"Socket Error: {err} | exiting...")
    except Exception as e:
        print(f"Generic exception, exiting...\n{e}")


if __name__ == "__main__":
    HOST = "0.0.0.0"
    # PORT = int(sys.argv[1])
    PORT = 5555
    conn_server((HOST, PORT))
