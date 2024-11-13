import { Buffer } from "node:buffer";

const SIZE_BYTES_COUNT = 4;
const HEADER_BYTES_COUNT = 3;
const COMMAND_OFFSET = SIZE_BYTES_COUNT + HEADER_BYTES_COUNT;

const readMessangeLength = (buffer) => {
  return buffer.readUInt32BE();
};

const readHeader = (buffer) => {
  const version = buffer.readUInt8(SIZE_BYTES_COUNT);
  const command = buffer.readUInt16BE(SIZE_BYTES_COUNT + 1);

  return {
    version,
    command,
  };
};

const readCommandLoginBody = (buffer) => {
  const correlationId = buffer.readUInt32BE(COMMAND_OFFSET);
  const stringLength = buffer.readUInt16BE(COMMAND_OFFSET + 4);
  const username = buffer.toString(
    "utf8",
    COMMAND_OFFSET + 6,
    COMMAND_OFFSET + 6 + stringLength
  );

  return {
    correlationId,
    username,
  };
};

const readCommandMessageBody = (buffer) => {
  let offset = COMMAND_OFFSET;
  const correlationId = buffer.readUInt32BE(offset);
  offset += 4;

  let stringLength = buffer.readUInt16BE(offset);
  offset += 2;
  const message = buffer.toString("utf8", offset, offset + stringLength);
  offset += stringLength;

  stringLength = buffer.readUInt16BE(offset);
  offset += 2;
  const from = buffer.toString("utf8", offset, offset + stringLength);
  offset += stringLength;

  stringLength = buffer.readUInt16BE(offset);
  offset += 2;
  const to = buffer.toString("utf8", offset, offset + stringLength);
  offset += stringLength;

  const time = buffer.readBigUInt64BE(offset);
  offset += 8;

  return {
    correlationId,
    message,
    from,
    to,
    time,
  };
};

const createResponse = (correlationId, code) => {
  const buffer = Buffer.alloc(13);
  let offset = 0;

  buffer.writeUInt32BE(0x09, offset);
  offset += 4;

  buffer.writeUInt16BE(0x03, offset);
  offset += 2;

  buffer.writeUInt8(0x01, offset);
  offset += 1;

  buffer.writeUInt32BE(correlationId, offset);
  offset += 4;

  buffer.writeUInt16BE(code, offset);

  return buffer;
};

export {
  SIZE_BYTES_COUNT,
  HEADER_BYTES_COUNT,
  COMMAND_OFFSET,
  readMessangeLength,
  readHeader,
  readCommandLoginBody,
  readCommandMessageBody,
  createResponse,
};
