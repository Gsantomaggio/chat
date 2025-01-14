using System.Net.Sockets;
using server.src;

/// <summary>
/// Provides extension methods for reading big-endian values using BinaryReader.
/// </summary>
internal static class EndianHelpers
{
    public static ushort ReadUInt16BE(this BinaryReader reader)
    {
        byte[] bytes = reader.ReadBytes(Constants.protocolUint16SizeBytes);
        if (BitConverter.IsLittleEndian)
            Array.Reverse(bytes);
        return BitConverter.ToUInt16(bytes, 0);
    }

    public static uint ReadUInt32BE(this BinaryReader reader)
    {
        byte[] bytes = reader.ReadBytes(Constants.protocolUint32SizeBytes);
        if (BitConverter.IsLittleEndian)
            Array.Reverse(bytes);
        return BitConverter.ToUInt32(bytes, 0);
    }

    public static ulong ReadUInt64BE(this BinaryReader reader)
    {
        byte[] bytes = reader.ReadBytes(Constants.protocolUint64SizeBytes);
        if (BitConverter.IsLittleEndian)
            Array.Reverse(bytes);
        return BitConverter.ToUInt64(bytes, 0);
    }

    public static byte[] GetBytesUInt16BE(ushort value)
    {
        byte[] bytes = BitConverter.GetBytes(value);
        if (BitConverter.IsLittleEndian)
            Array.Reverse(bytes);
        return bytes;
    }

    public static byte[] GetBytesUInt32BE(uint value)
    {
        byte[] bytes = BitConverter.GetBytes(value);
        if (BitConverter.IsLittleEndian)
            Array.Reverse(bytes);
        return bytes;
    }

    public static byte[] GetBytesUInt64BE(ulong value)
    {
        byte[] bytes = BitConverter.GetBytes(value);
        if (BitConverter.IsLittleEndian)
            Array.Reverse(bytes);
        return bytes;
    }
}