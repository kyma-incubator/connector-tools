"use strict"
const websocket = require('websocket-stream')
const WebSocketServer = require('ws').Server
const Connection = require('mqtt-connection')
const http = require('http')
const request = require('request-promise-native');

const envVariables = {
  appName: process.env.APPLICATION_NAME,
  oauthUrl: process.env.OAUTH_URL,
  port: process.env.PORT ? process.env.PORT : 8080,
  eventUrl: process.env.EVENT_URL ? process.env.EVENT_URL : "http://event-bus-publish.kyma-system.svc.cluster.local:8080/v1/events"
}

let wss = new WebSocketServer({
  server: http.createServer().listen(envVariables.port),
  verifyClient: verifyToken
});

console.log(`Application Started on port ${envVariables.port}`)

wss.on('connection', function (ws) {
  let stream = websocket(ws)
  // create MQTT Connection
  let connection = new Connection(stream)

  connection.on('publish', function (packet) {
    let event = createEvent(JSON.parse(packet.payload.toString()));

    console.log(`Publishing packet with ID ${packet.messageId} to event bus: ${JSON.stringify(event)}`);
    request.post({
      url: envVariables.eventUrl,
      json: event
    }).then((response) => {
      if (packet.messageId && response['event-id']) {
        console.log(`Sending back puback for message ${packet.messageId}`);
        connection.puback({ messageId: packet.messageId });
      }
    }).catch((err) => {
      console.error(`Error while forwarding packet with ID ${packet.messageId}, error is: ${JSON.stringify(err, null, 2)}`)
    });
  })

  connection.on('connect', function (packet) {
    console.log('Client connecting');
    connection.connack({ returnCode: 0 });
  });

  connection.on('pingreq', function () {
    client.pingresp()
  });

  connection.on('error', function (error) {
    console.log(`Error: ${JSON.stringify(error, null, 2)}`)
  })

  connection.on('subscribe', function (packet) {
    console.log('Client subscribing');
    connection.suback({
      granted: [packet.qos],
      messageId: packet.messageId
    })
  })
});

async function verifyToken(info, cb) {
  if (!envVariables.oauthUrl) {
    console.error("Skipping token validation as OAUTH2_URL is not configured")
    cb(true);
    return
  }

  let headerToken = info.req.headers['authorization'];
  if (!headerToken) {
    console.error("No Authorization header provided in request")
    cb(false, 401, 'Unauthorized request: no Authorization header provided');
    return;
  }

  if (!headerToken.startsWith("Bearer ")) {
    console.error(`Authorizatin header does not contain a Bearer token, header value is: ${headerToken}`)
    cb(false, 401, 'Unauthorized request: no Bearer token provided in Authorization header');
    return;
  }

  try {
    var response = await request.get({
      uri: envVariables.oauthUrl + "/validate",
      method: 'GET',
      headers: {
        'Authorization': headerToken
      }
    });
    if (response.statusCode >= 300) {
      console.log(`Validation of bearer token failed with status code ${response.statusCode}, token was: ${token}`);
      cb(false, 403, 'Validation of Bearer token failed');
      return
    }

    console.log('Token validation successful');
    cb(true);
  }
  catch (err) {
    console.log(`Error while calling validate method of OAuth server ${JSON.stringify(err, null, 2)}`);
    cb(false, err.statusCode, 'Problem while validating provided Bearer token');
  }
}

function createEvent(msg) {
  return {
    "source-id": envVariables.appName,
    "event-type": msg.eventType.toLowerCase(),
    "event-type-version": "v1",
    "event-time": msg.eventTime,
    "data": msg.data
  }
}
