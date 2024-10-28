# Console Chat

### Disclaimer:

This project's scope is only to learn how to write a client-server application from scratch.
The code is demonstrative and focuses more on understanding how the TCP works than production code.
The client applications are meant to test the server quickly.

We provide complete examples in the following languages:

- Server golang (server/go)
- Client golang (server/go/tcp_client)
  ..... ???

### How to run the golang stack:

- Server: /server/go/run/server and run `go run main.go`
- Client: /server/go/run/client and run `go run main.go localhost:5555`
- You can use two different terminals with two different users

### Protocol definition:

The protocol is a binary protocol with the following structure:

### Data types

- `byte` 1 byte
- `uint16` 2 bytes
- `uint32` 4 bytes
- `string` 2 bytes for the length + N bytes for the string
- `[]byte` 2 bytes for the length + N bytes for the string
- `uint64` 8 bytes

### Header

| Name      | Type     |
|-----------|----------|
| `version` | `byte`   | 
| `command` | `uint16` |

### CommandLogin

| Name            | Type     | value(s) | reference         |
|-----------------|----------|----------|-------------------|
| `key`           | `uint16` | 0x01     | `Header::command` |
| `version`       | `uint16` | 0x01     | `Header::version` |
| `correlationId` | `uint32` |          |                   |
| `username`      | `string` |          |                   |

### CommandMessage

| Name            | Type     | value(s) | reference         |
|-----------------|----------|----------|-------------------|
| `key`           | `uint16` | 0x02     | `Header::command` |
| `version`       | `uint16` | 0x01     | `Header::version` |
| `correlationId` | `uint32` |          |                   |
| `message`       | `string` |          |                   |
| `From`          | `string` |          |                   |
| `To`            | `string` |          |                   |
| `Time`          | `uint64` |          |                   |

## Responses

| Name            | Type     | value(s) | reference         |
|-----------------|----------|----------|-------------------|
| `key`           | `uint16` | 0x03     | `Header::command` |
| `version`       | `uint16` | 0x01     | `Header::version` |
| `correlationId` | `uint32` |          |                   |
| `code`          | `uint16` |          | `ResponseCodes`   |

### ResponseCodes

| Name                     | value(s) | 
|--------------------------|----------|
| `OK`                     | 0x01     |
| `ErrorUserNotFound`      | 0x02     |
| `ErrorUserAlreadyLogged` | 0x03     |


### Write data to the socket

- Write the length of whole message (header + command) as a `uint32`
- Write the header
- Write the command

For Example: `CommandLogin` with username `user1`

- `Header`:
  - `version` = 0x01 // 1 byte
  - `command` = 0x01 // 2 bytes
- `CommandLogin`:
  - `correlationId` = 0x01 // 4 bytes
  - `username` = "user1" // length = 2 + 5 = 7

- Write the length of the whole message: 3 + 11  = 14
- Write the `header` + `message`:
- 0x01 //  - version  (1 byte)
- 0x00 0x01  // command ( 2 bytes)
- 0x00 0x00 0x00 0x01  // correlationId (4 bytes) 
- 0x00 0x01 // username length  (2 bytes)
- 0x75 0x73 0x65 0x72 0x31 // username (5 bytes)
- Total bytes written: 14 + 4 = 20
- Send the message
- Read the response

### Read data from the socket

- Read the length of the whole message (header + command) as a `uint32`
- Read the header
- Read the command key 
- Read the command based on the key

For Example: `CommandLogin` with username `user1`

- Read the first 4 bytes for the length of the whole message: 14
- Ensure the socket buffer has at least 14 bytes
- Read the header:
  - Read the version: 0x01
  - Read the command: 0x01
  - Read the command based on the key: `CommandLogin`
    - Read the correlationId: 0x01
    - Read the username: "user1"
  - Process the command
  - Send the response



 


      




