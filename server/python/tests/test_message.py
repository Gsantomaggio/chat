import unittest
from datetime import datetime
from source.message import (
    Message,
)  # Assicurati di importare correttamente la tua classe


class TestMessage(unittest.TestCase):
    def setUp(self):
        self.correlationId = 1
        self.message = "Hello, World!"
        self.from_field = "user1"
        self.to_field = "user2"
        self.timestamp = int(datetime.now().timestamp())
        self.msg = Message(
            self.correlationId,
            self.message,
            self.from_field,
            self.to_field,
            self.timestamp,
        )

    def test_initialization(self):
        self.assertEqual(self.msg.correlationId, self.correlationId)
        self.assertEqual(self.msg.message, self.message)
        self.assertEqual(self.msg.from_field, self.from_field)
        self.assertEqual(self.msg.to_field, self.to_field)
        self.assertEqual(self.msg.timestamp, self.timestamp)

    def test_str(self):
        expected_str = (
            f"\nmessage: {self.message}\nfrom: {self.from_field}\nto: {self.to_field}"
        )
        self.assertEqual(str(self.msg), expected_str)


if __name__ == "__main__":
    unittest.main()
