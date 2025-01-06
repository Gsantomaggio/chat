import unittest
from unittest.mock import Mock
from socket import socket
import time
from datetime import datetime
from source.users import User, logout


class TestUser(unittest.TestCase):
    def setUp(self):
        self.mock_socket = Mock(spec=socket)
        self.user = User("testuser")

    def test_initial_state(self):
        self.assertEqual(self.user.username, "testuser")
        self.assertIsNone(self.user.lastlogin)
        self.assertFalse(self.user.isonline)
        self.assertEqual(self.user.status, "offline")
        self.assertIsNone(self.user.conn)
        self.assertEqual(self.user.messages, [])

    def test_login_success(self):
        response_code = self.user.login(self.mock_socket)
        self.assertEqual(response_code, 1)
        self.assertTrue(self.user.isonline)
        self.assertEqual(self.user.status, "online")
        self.assertEqual(self.user.conn, self.mock_socket)
        self.assertIsNotNone(self.user.lastlogin)

    def test_login_already_logged(self):
        self.user.isonline = True
        self.user.status = "online"
        response_code = self.user.login(self.mock_socket)
        self.assertEqual(response_code, 4)
        self.assertTrue(self.user.isonline)
        self.assertEqual(self.user.status, "online")
        self.assertIsNone(self.user.conn)  # Connessione non dovrebbe cambiare

    def test_logout(self):
        self.user.login(self.mock_socket)
        logout(self.user)
        self.assertFalse(self.user.isonline)
        self.assertEqual(self.user.status, "offline")
        self.assertIsNone(self.user.conn)

    def test_printlastlogin(self):
        self.assertEqual(self.user.printlastlogin(), "Never")
        self.user.lastlogin = time.time()
        expected_time = datetime.fromtimestamp(self.user.lastlogin).strftime(
            "%a, %d %b %Y %H:%M:%S"
        )
        self.assertEqual(self.user.printlastlogin(), expected_time)

    def test_str(self):
        self.assertEqual(
            str(self.user),
            "testuser - status: offline - last login: Never - messages in queue: 0",
        )
        self.user.login(self.mock_socket)
        expected_time = datetime.fromtimestamp(self.user.lastlogin).strftime(
            "%a, %d %b %Y %H:%M:%S"
        )
        self.assertEqual(
            str(self.user),
            f"testuser - status: online - last login: {expected_time} - messages in queue: 0",
        )


if __name__ == "__main__":
    unittest.main()
