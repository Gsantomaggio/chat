import socket


def run_server(host: str, port: str, backlog: int = 1) -> None:
    try:
        with socket.socket() as s:
            s.bind((host, port))
            s.listen(backlog)
            conn, addr = s.accept()
            with conn:
                print(f"Connected by {addr}")
                while True:
                    data = conn.recv(1024)
                    data_string = str(data, "utf-8")
                    is_exit = data_string.upper() == "ESC"
                    if is_exit:
                        print("Closing connection...")
                        break
                    print(data_string)
                    new_data = f"Received message: {data_string}"
                    conn.sendall(new_data.encode())
    except socket.error as err:
        print(f"Socket Error: {err} | exiting...")
    except Exception:
        print("Generic exception, exiting...")


if __name__ == "__main__":
    HOST = "127.0.0.1"
    PORT = 65432
    run_server(HOST, PORT)