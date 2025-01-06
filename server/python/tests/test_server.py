import unittest
import multiprocessing
import threading
import time
import os
import signal
from test_client import TcpClient
from source.server import TcpServer


class TcpServerTest(TcpServer):
    def stop_server(self):
        pass


def run_server():
    server = TcpServerTest(host="localhost", port=6666)
    server.run_server()


class TestTcpServer(unittest.TestCase):
    @classmethod
    def setUpClass(self):
        self.server_process = multiprocessing.Process(target=run_server)
        self.server_process.start()
        time.sleep(1)

    @classmethod
    def tearDownClass(self):
        self.stop_server()

    @classmethod
    def stop_server(self):
        os.kill(self.server_process.pid, signal.SIGINT)
        self.server_process.join()

    def run_client(self, func, *args):
        client_thread = threading.Thread(target=func, args=args)
        client_thread.start()
        client_thread.join()

    def test_login_success(self):
        self.run_client(self._test_login_success)

    def _test_login_success(self):
        client = TcpClient(("localhost", 6666))
        client.connect()
        response = client.login("user1")
        client.close()
        self.assertEqual(response["code"], 0x01)

    def test_login_duplicate(self):
        self.run_client(self._test_login_duplicate)

    def _test_login_duplicate(self):
        client1 = TcpClient(("localhost", 6666))
        client1.connect()
        response1 = client1.login("user1")

        client2 = TcpClient(("localhost", 6666))
        client2.connect()
        response2 = client2.login("user1")

        client1.close()
        client2.close()
        self.assertEqual(response1["code"], 0x01)
        self.assertEqual(response2["code"], 0x04)

    def test_send_message(self):
        self.run_client(self._test_send_message)

    def _test_send_message(self):
        client1 = TcpClient(("localhost", 6666))
        client1.connect()
        client1.login("user1")

        client2 = TcpClient(("localhost", 6666))
        client2.connect()
        client2.login("user2")

        response = client2.send_message("Hello", "user1", "user2")

        client1.close()
        client2.close()
        self.assertEqual(response["code"], 0x01)


if __name__ == "__main__":
    unittest.main()
