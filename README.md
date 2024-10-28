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

 


      




