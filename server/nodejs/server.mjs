import net from "node:net";
import {
  createResponse,
  readCommandLoginBody,
  readCommandMessageBody,
  readHeader,
  readMessangeLength,
} from "./protocol.mjs";

const users = {};

const getSocketId = (socket) => {
  return `${socket.remoteAddress}:${socket.remotePort}`;
};

const server = net.createServer((socket) => {
  const socketId = getSocketId(socket);
  console.log(`Client connected: ${socketId}`);

  //TODO: check for bytes read to match the expected length (at least)
  socket.on("data", (data) => {
    const bytesToRead = readMessangeLength(data);

    console.log(
      `Bytes to read are ${bytesToRead} and data is ${data.length} bytes`
    );

    const header = readHeader(data);

    console.log(`received header is ${JSON.stringify(header)}`);

    let body;
    if (header.command === 0x01) {
      body = readCommandLoginBody(data);
      users[body.username] = { socketId, status: "online" };
      console.log(`User ${body.username} is now online`);
      console.log(`CorrelationId is ${body.correlationId}`);
      const response = createResponse(body.correlationId, 0x01);
      socket.write(response);
      console.log("Response sent:", response);
    } else if (header.command === 0x02) {
      body = readCommandMessageBody(data);

      // Verifica che nel campo From ci sia lo username che corrisponde all'utente attualmente connesso sul socket
      if (users[body.from] && users[body.from].socketId === socketId) {
        // Verifica che nel campo To ci sia lo username di un utente presente nell'array users
        if (users[body.to]) {
          console.log(
            `Message from ${body.from} to ${body.to} at ${body.time}: ${body.message}`
          );

          // Invia il messaggio al destinatario
          const recipientSocketId = users[body.to].socketId;
          const recipientSocket = server.connections.find(
            (conn) => getSocketId(conn) === recipientSocketId
          );
          if (recipientSocket) {
            recipientSocket.write(body.message);
          }

          // Costruisce la risposta
          const response = createResponse(body.correlationId, 0); // Utilizziamo il codice 0 per indicare successo

          // Invia la risposta al client
          socket.write(response);
        } else {
          console.error(`User ${body.to} not found`);
          // Costruisce la risposta con codice di errore
          const response = createResponse(body.correlationId, 0x03); // Utilizziamo il codice 0x03 per indicare utente non trovato
          socket.write(response);
        }
      } else {
        console.error(`Invalid From field: ${body.from}`);
        // Costruisce la risposta con codice di errore
        const response = createResponse(body.correlationId, 0x02); // Utilizziamo il codice 0x02 per indicare errore di autenticazione
        socket.write(response);
      }
    } else {
      // TODO: throw Unsupported Error
    }
    console.log("Received message:", { header, body });
  });

  socket.on("end", () => {
    const username = Object.keys(users).find(
      (key) => users[key].socketId === socketId
    );
    if (username) {
      users[username].status = "offline";
      console.log(`User ${username} is now offline`);
    }
    console.log(`Client disconnected: ${socketId}`);
  });

  socket.on("error", (err) => {
    console.error(`Error: ${err}`);
  });
});

const port = process.env.PORT || 5555;

server.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
