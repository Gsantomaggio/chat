import net from "node:net";
import createDebug from "debug";
const debug = createDebug("server");
import {
  COMMAND_CODES,
  createResponse,
  readCommandLoginBody,
  readCommandMessageBody,
  readHeader,
  readMessangeLength,
  RESPONSE_CODES,
} from "./protocol.mjs";
import { socketToId, getSocketsExcluding } from "./utils.mjs";

const users = {};
const mailboxes = {};
let sockets = [];

const setUserOnline = (username, socketId) => {
  users[username] = { socketId, status: "online" };
  if (!mailboxes[username]) {
    mailboxes[username] = [];
  }
  console.log(JSON.stringify(users));
};
const setUserOffline = (username) => {
  users[username].status = "offline";
};
const isUserOnline = (username) => {
  if (!users[username]) {
    return false;
  }
  debug(`${username} found`);
  return users[username].status === "online";
};

const queueMessage = (to, cmdMsgBuf) => {
  mailboxes[to].push(cmdMsgBuf);
};
const getUsernameFromSocketId = (socketId) =>
  Object.keys(users).find((key) => users[key].socketId === socketId);

const server = net.createServer((socket) => {
  sockets.push(socket);
  const socketId = socketToId(socket);
  console.log(`Client connected: ${socketId}`);

  //TODO: check for bytes read to match the expected length (at least)
  socket.on("data", (data) => {
    const bytesToRead = readMessangeLength(data);
    debug(`Bytes to read are ${bytesToRead} and data is ${data.length} bytes`);

    const { command, version } = readHeader(data);
    debug(`Received header with command ${command} and version ${version}`);

    if (command === COMMAND_CODES.LOGIN) {
      const { username, correlationId } = readCommandLoginBody(data);
      setUserOnline(username, socketId);
      console.log(`User ${username} is now online`);

      const response = createResponse(correlationId, RESPONSE_CODES.OK);
      socket.write(response);
      debug(`Login response sent: ${response.toString("hex")}`);

      //TODO: improve
      mailboxes[username].forEach((cmdMsgBuf) => {
        socket.write(cmdMsgBuf);
      });
      mailboxes[username] = [];
      // ---
    } else if (command === COMMAND_CODES.MESSAGE) {
      const { to, from, message, time, correlationId } =
        readCommandMessageBody(data);

      /*
      // NOT REQUESTED FROM SPEC
      const isFromCurrentUser =
        users[from] && users[from].socketId === socketId;
      if (!isFromCurrentUser) {
        console.error(`Invalid From field: ${from}`);
        const response = createResponse(correlationId, 0x02);
        socket.write(response);
        return;
      }
      */

      const recipientExists = users[to];
      if (!recipientExists) {
        console.error(`User ${to} not found`);
        const response = createResponse(
          correlationId,
          RESPONSE_CODES.USER_NOT_FOUND
        );
        socket.write(response);
        debug(`Message response sent: ${response.toString("hex")}`);
        return;
      }

      debug(JSON.stringify(users));
      if (isUserOnline(to)) {
        debug(`${to} is online`);
        const recipientSocketId = users[to].socketId;
        const recipientSocket = sockets.find(
          (s) => socketToId(s) === recipientSocketId
        );
        recipientSocket.write(data); //hack OR perf, choose one
      } else {
        debug(`${to} is offline`);
        queueMessage(to, data);
      }
      const response = createResponse(correlationId, RESPONSE_CODES.OK);
      socket.write(response);
      debug(`Message response sent: ${response.toString("hex")}`);
    } else {
      console.error(`Command not supported: ${command}`);
    }
  });

  socket.on("end", () => {
    sockets = getSocketsExcluding(sockets, socket);
    const username = getUsernameFromSocketId(socketId);
    setUserOffline(username);
    console.log(`User ${username} is now offline`);
    debug(`Client disconnected: ${socketId}`);
  });

  socket.on("error", (err) => {
    console.error(`Error: ${err}`);
  });
});

const port = process.env.PORT || 5555;

server.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
