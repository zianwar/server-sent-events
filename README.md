# Server-Sent Events (SSE) Demo

This is a simple demo application that uses [Server-Sent Events (SSE)](https://en.wikipedia.org/wiki/Server-sent_events) to push updates from a Go server to a web clients. SSE is a standard mechanism that allows servers to send real-time updates to clients over a single, long-lived HTTP connection.

SSE is part of HTML5 and is standardized in the EventSource API.

- It's simple to Implement: SSE is simpler to implement compared to WebSockets.
- Most modern browsers support it
- SSE is a text-based protocol which is suitable for applications where the updates are primarily text.

Comparison with WebSockets:

- SSE is one-way (server to client), while WebSockets support full-duplex communication (both server to client and client to server).
- SSE automatically handles reconnections, whereas WebSockets require manual handling of reconnections.
- SSE is ideal for applications like live news updates, notifications, and real-time feeds. WebSockets are more suitable for real-time gaming, chat applications, and other use cases that require bi-directional communication.

# Usage

Clone repo and run server:

```sh
git clone https://github.com/zianwar/server-sent-events
cd server-sent-events
go mod tidy
go run .
```

Navigate to http://localhost:9000

- The server will send periodic messages to the client, which will be displayed on the web page.
- The web page will automatically connect to the server and display any messages sent by the server.

<p align="center">
  <img src="https://pub.anw.sh/sse.jpg" style="width:600px" />
</p>
