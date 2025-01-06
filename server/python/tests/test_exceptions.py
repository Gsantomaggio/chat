import unittest
from source.exceptions import AlreadyLoggedException


class TestAlreadyLoggedException(unittest.TestCase):
    def test_exception_message(self):
        message = "User is already logged in"
        exception = AlreadyLoggedException(message)
        self.assertEqual(str(exception), message)
        self.assertIsNone(exception.errors)

    def test_exception_with_errors(self):
        message = "User is already logged in"
        errors = {"error_code": 123}
        exception = AlreadyLoggedException(message, errors)
        self.assertEqual(str(exception), message)
        self.assertEqual(exception.errors, errors)


if __name__ == "__main__":
    unittest.main()
