import socket
import sys

from source import protocol


def send_messages(s):
    # while True:
    for i in range(2):
        if i == 0:
            message = protocol.login_message
        else:
            message = protocol.message_message
        s.send(message)
        ##################################
        # message = input("msg: ->\t ")
        # s.send(message.encode())
        ##################################
        if message.upper() == "ESC":
            print("Closing connection to the server.")
            s.close()
            break
        else:
            data = s.recv(4096)
            print(str(data, "utf-8"))


def conn_server(serverAddress: tuple) -> None:
    try:
        with socket.socket() as s:
            s.connect(serverAddress)
            print(f"Connected to the server:\t{serverAddress=}")
            send_messages(s)
    except socket.error as err:
        print(f"Socket Error: {err} | exiting...")
    except Exception:
        print("Generic exception, exiting...")


if __name__ == "__main__":
    HOST = sys.argv[1]
    PORT = int(sys.argv[2])
    conn_server((HOST, PORT))
