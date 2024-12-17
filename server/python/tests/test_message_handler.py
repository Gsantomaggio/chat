import unittest
from unittest.mock import Mock, patch
import time
from socket import socket
from source.exceptions import AlreadyLoggedException
from source.message import Message
from source.users import User
from source.message_handler import MessageHandler


class TestMessageHandler(unittest.TestCase):
    def setUp(self):
        self.mock_socket = Mock(spec=socket)
        self.mock_user = Mock(spec=User)
        self.mock_users = {}
        self.handler = MessageHandler(self.mock_socket, self.mock_users)

    def test_read_message_login_success(self):
        buffer = b"\x00\x00\x00\x0e\x01\x00\x01\x00\x00\x00\x01\x00\x05user1"
        self.mock_user.login.return_value = 1
        self.mock_user.username = 'user1'
        self.mock_user.conn = self.mock_socket
        self.mock_user.messages = []
        self.mock_users["user1"] = self.mock_user

        with patch("source.message_handler.read_string", return_value=("user1", 14)):
            user = self.handler.read_message(buffer, None)

        self.assertEqual(user.username, "user1")
        self.mock_user.login.assert_called_once_with(self.mock_socket)
        self.mock_socket.send.assert_called()

    def test_read_message_login_already_logged(self):
        buffer = b'\x00\x00\x00\x0e\x01\x00\x01\x00\x00\x00\x01\x00\x05user1'
        self.mock_user.login.return_value = 4
        self.mock_user.username = 'user1'
        self.mock_user.conn = self.mock_socket
        self.mock_user.messages = []
        self.mock_users['user1'] = self.mock_user

        with patch('source.message_handler.read_string', return_value=('user1', 14)):
            with self.assertRaises(AlreadyLoggedException):
                self.handler.read_message(buffer, None)

        self.mock_user.login.assert_called_once_with(self.mock_socket)
        self.mock_socket.send.assert_called()

    def test_read_message_send(self):
        buffer = b'\x00\x00\x00\x24\x01\x00\x02\x00\x00\x00\x01\x00\x05Hello\x00\x05user2\x00\x05user1' + int(time.time()).to_bytes(8)
        user = User('user1')
        user.conn = self.mock_socket
        self.mock_users['user2'] = User('user2')

        with patch('source.message_handler.read_string', side_effect=[('Hello', 18), ('user2', 25), ('user1', 32)]):
            with patch('source.message_handler.read_timestamp', return_value=(int(time.time()), 40)):
                self.handler.read_message(buffer, user)

        self.mock_socket.send.assert_called()

    def test_read_message_length(self):
        buffer = b'\x00\x00\x00\x01\x01'
        length, offset = self.handler._read_message_length(buffer)
        self.assertEqual(length, 1)
        self.assertEqual(offset, 4)

    def test_read_correlationId(self):
        buffer = b'\x00\x00\x00\x01'
        correlationId, offset = self.handler._read_correlationId(buffer, 0)
        self.assertEqual(correlationId, 1)
        self.assertEqual(offset, 4)

    def test_read_command_message(self):
        buffer = b'\x00\x05Hello' \
                 b'\x00\x05user2' \
                 b'\x00\x05user1' + int(time.time()).to_bytes(8)

        message = self.handler._read_command_message(buffer, 0, 1)

        self.assertEqual(message.correlationId, 1)
        self.assertEqual(message.message, 'Hello')
        self.assertEqual(message.from_field, 'user2')
        self.assertEqual(message.to_field, 'user1')

    def test_create_command_message(self):
        timestamp = int(time.time())
        message = Message(1, 'Hello', 'user2', 'user1', timestamp)

        mex = self.handler._create_command_message(message)

        expected_message = b'\x01' + b'\x00\x02' + (1).to_bytes(4) + (5).to_bytes(2) + b'Hello' + (5).to_bytes(2) + b'user2' + (5).to_bytes(2) + b'user1' + (timestamp).to_bytes(8)
        expected_length = (len(expected_message)).to_bytes(4)
        expected_mex = expected_length + expected_message

        self.assertEqual(mex, expected_mex)

    # def test_send_message(self):
    #     message = Message(1, 'Hello', 'user2', 'user1', int(time.time()))

    #     user2 = User('user2')
    #     user2.conn = Mock(spec=socket)

    #     users = {'user2': user2}

    #     handler = MessageHandler(self.mock_socket, users)

    #     response_code = handler._send_message(message)

    #     user2.conn.send.assert_called()

    #     self.assertEqual(response_code, 1)

    # def test_send_user_messages(self):
    #     user = User('user1')

    #     message = Message(1, 'Hello', 'user2', 'user1', int(time.time()))

    #     user.messages.append(message)

    #     user.conn = Mock(spec=socket)

    #     remaining_messages = self.handler._send_user_messages(user)

    #     user.conn.send.assert_called()

    #     self.assertEqual(len(remaining_messages), 0)

    # def test_send_response(self):
    #     response = self.handler._send_response(1, 1, self.mock_socket)

    #     expected_response_prefix = (9).to_bytes(4) + (1).to_bytes(1) + (3).to_bytes(2) + (1).to_bytes(4) + (1).to_bytes(2)

    #     self.mock_socket.send.assert_called()

    #     self.assertTrue(response.startswith(expected_response_prefix))


if __name__ == "__main__":
    unittest.main()
