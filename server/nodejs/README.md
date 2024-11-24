# Node.js Chat Server

This is a simple chat server written in Node.js (JavaScript).

## Running the server

- `npm install` (dependencies)
- `node server.mjs`
- `DEBUG=server,protocol node server.mjs` (to see more logs)

### Features

- [x] Login (without password. It is enough to send the username)
- [x] Send message and dispatch to the correct user
- [x] Store in memory the users with the status (online/offline)
- [x] Store in memory the off-line messages when the user is not online
- [ ] Send the off-line messages when the user logs in
- [ ] Check if the user is already logged in
- [x] Check if the destination user exists
