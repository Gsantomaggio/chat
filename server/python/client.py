import socket
import sys

from source import protocol


def send_messages(s, user):
    if user == 2:
        message = protocol.login_message
    # while True:
    for i in range(2):
        if i == 0 and user == 1:
            message = protocol.login_message
        elif i == 0 and user == 2:
            message = protocol.login_message2
        elif i == 1 and user == 1:
            message = protocol.message_message
        else:
            message = None
        if message:
            s.send(message)
            ##################################
            # message = input("msg: ->\t ")
            # s.send(message.encode())
            ##################################
            if message.upper() == "ESC":
                print("Closing connection to the server.")
                s.close()
                break
        data = s.recv(4096)
        print(str(data, "utf-8"))


def conn_server(serverAddress: tuple, user: int) -> None:
    try:
        with socket.socket() as s:
            s.connect(serverAddress)
            print(f"Connected to the server:\t{serverAddress=}")
            send_messages(s, user)
    except socket.error as err:
        print(f"Socket Error: {err} | exiting...")
    except Exception as e:
        print(f"Generic exception, exiting...\n{e}")


if __name__ == "__main__":
    HOST = sys.argv[1]
    PORT = int(sys.argv[2])
    user = int(sys.argv[3])
    conn_server((HOST, PORT), user)
