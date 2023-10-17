// Run websocket server
const wsServer = Bun.serve({
  port: 8080,
  fetch(req, server) {
    // upgrade the request to a WebSocket
    if (server.upgrade(req)) {
      return; // do not return a Response
    }
    return new Response("Upgrade failed :(", { status: 500 });
  },
  websocket: {
    message(ws, message) { // a message is received
      ws.send(message); // echo back the message
    },
    open(ws) {}, // a socket is opened
    close(ws, code, message) {}, // a socket is closed
    drain(ws) {}, // the socket is ready to receive more data
  },
});

// Run web server
const webServer = Bun.serve({
  port: 8000,
  //unix: "/run/sanic.sock",
  fetch(req) {
    return new Response("Put frontend here!");
  },
});

console.log(`Listening on http://localhost:${webServer.port} ...`);
console.log(`Listening on ws://localhost:${wsServer.port} ...`);

