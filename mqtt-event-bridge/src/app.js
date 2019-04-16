"use strict"
const websocket = require('websocket-stream')
const WebSocketServer = require('ws').Server
const Connection = require('mqtt-connection')
const http = require('http')
const request = require('request-promise-native');
const server = http.createServer().listen(8080);

let envVariables = {
  sourceId: process.env.SOURCE_ID,
  applicationName: process.env.APPLICATION_NAME,
  auth_username: process.env.auth_username,
  auth_password: process.env.auth_password,
  oauth2_client_id: process.env.oauth2_client_id,
  oauth2_client_secret: process.env.oauth2_client_secret,
  oauth_server: process.env.oauth_server
}

const wss = new WebSocketServer({
  server: server
});
console.log("Application started")
let eventEndpoint = "http://event-bus-publish.kyma-system.svc.cluster.local:8080/v1/events";
console.log(process.env.DEBUG)
if (process.env.DEBUG) {
  console.log(envVariables)
  eventEndpoint = "http://localhost:4000/v1/events"
}
wss.on('connection', function (ws) {
  let stream = websocket(ws)
  // create MQTT Connection
  let connection = new Connection(stream)

  // client published
  connection.on('publish', function (packet) {

    let event = createEvent(JSON.parse(packet.payload.toString()));
    sendEvent(event).then((response) => {
      if (packet.messageId && response['event-id']) {
        console.log(`puback for message ${packet.messageId}`);
        connection.puback({ messageId: packet.messageId });
        // send a puback with messageId (for QoS > 0) if we got a valid response
      }
    }).catch((err) => console.error(err));
  })

  connection.on('connect', function (packet) {
    console.log('Client connecting');
    connection.connack({ returnCode: 0 });
  });

  // client pinged
  connection.on('pingreq', function () {
    client.pingresp()
  });

  // client subscribed
  connection.on('subscribe', function (packet) {
    // send a suback with messageId and granted QoS level
    connection.suback({ granted: [packet.qos], messageId: packet.messageId })
  })

});

function createEvent(msg) {
  return {
    "source-id": envVariables.sourceId.replace(/(^\w+:|^)\/\//, ''), //remove protocol : https://asd.com -> asd.com, else kyma dont accept
    "event-type": msg.eventType,
    "event-type-version": msg.eventVersion,
    "event-time": msg.eventTime,
    "data": msg.data
  }

}

function sendEvent(eventData) {
  console.log("Publishing event to event bus: " + JSON.stringify(eventData));
  return new Promise((resolve, reject) => {
    request.post({ url: eventEndpoint, json: eventData }).then((response) => {
      resolve(response)
    }
    ).catch((err) => reject(err))
  })
}
