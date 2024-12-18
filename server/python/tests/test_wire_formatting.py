import unittest
import time
from source.wire_formatting import (
    read_header,
    read_string,
    read_uint64,
    read_uint32,
    read_uint16,
    read_uint8,
    read_timestamp,
    write_uint8,
    write_uint16,
    write_uint32,
)


class TestWireFormatting(unittest.TestCase):
    def test_read_uint8(self):
        buffer = b"\x01"
        value, offset = read_uint8(buffer, 0)
        self.assertEqual(value, 1)
        self.assertEqual(offset, 1)

    def test_read_uint16(self):
        buffer = b"\x00\x01"
        value, offset = read_uint16(buffer, 0)
        self.assertEqual(value, 1)
        self.assertEqual(offset, 2)

    def test_read_uint32(self):
        buffer = b"\x00\x00\x00\x01"
        value, offset = read_uint32(buffer, 0)
        self.assertEqual(value, 1)
        self.assertEqual(offset, 4)

    def test_read_uint64(self):
        buffer = b"\x00\x00\x00\x00\x00\x00\x00\x01"
        value, offset = read_uint64(buffer, 0)
        self.assertEqual(value, 1)
        self.assertEqual(offset, 8)

    def test_read_string(self):
        buffer = b"\x00\x05hello"
        value, offset = read_string(buffer, 0)
        self.assertEqual(value, "hello")
        self.assertEqual(offset, 7)

    def test_read_header(self):
        buffer = b"\x01\x00\x02"
        version, command, offset = read_header(buffer, 0)
        self.assertEqual(version, 1)
        self.assertEqual(command, 2)
        self.assertEqual(offset, 3)

    def test_read_timestamp(self):
        timestamp = int(time.time())
        buffer = timestamp.to_bytes(8, "big")
        value, offset = read_timestamp(buffer, 0)
        self.assertEqual(value, timestamp)
        self.assertEqual(offset, 8)

    def test_write_uint8(self):
        value = 1
        result = write_uint8(value)
        self.assertEqual(result, b"\x01")

    def test_write_uint16(self):
        value = 1
        result = write_uint16(value)
        self.assertEqual(result, b"\x00\x01")

    def test_write_uint32(self):
        value = 1
        result = write_uint32(value)
        self.assertEqual(result, b"\x00\x00\x00\x01")


if __name__ == "__main__":
    unittest.main()
