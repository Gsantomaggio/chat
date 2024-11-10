Rust Chat Server
================

This is a simple chat server written in Rust. It uses the `tokio`.

## Running the server
- `cargo run`

### Features

- [x] Login (without password. It is enough to send the username)
- [x] Send message and dispatch to the correct user
- [ ] Store in memory the users with the status (online/offline)
- [ ] Store in memory the off-line messages when the user is not online
- [ ] Send the off-line messages when the user logs in
- [ ] Check if the user is already logged in
- [x] Check if the destination user exists

