Golang Chat Server
================

This is a simple chat server written in Golang.

## Running the server
- `go run run/server/main.go`

### Features

- [x] Login (without password. It is enough to send the username)
- [x] Send message and dispatch to the correct user
- [x] Store in memory the users with the status (online/offline)
- [x] Store in memory the off-line messages when the user is not online
- [x] Send the off-line messages when the user logs in
- [x] Check if the user is already logged in
- [x] Check if the destination user exists

