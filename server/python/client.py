import socket

def send_messages(s):
    while True:
        message = input("msg: ->\t ")
        s.send(message.encode())
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


if __name__ == '__main__':
    HOST = "127.0.0.1"
    PORT = 65432
    conn_server((HOST, PORT))