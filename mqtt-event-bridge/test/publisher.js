"use strict"
const mqtt = require('mqtt');
const request = require('request-promise-native');

const mqttServer = process.env.MQTT_SERVER;
const oauthServer = process.env.OAUTH_SERVER;
const client_id = process.env.OAUTH_CLIENT_ID;
const client_secret = process.env.OAUTH_CLIENT_SECRET;

//const mqttServrer= "wss://hb-marketing-default-4c5417f4-6040-11e9-82a1-0a580a40-mqtt.sjanota.kyma.pro"

(async () => {

    const oauthToken = await request({
        url: oauthServer, method: 'POST', json: true,
        form: {
            'grant_type': 'client_credentials',
            'client_id': client_id,
            'client_secret': client_secret
        }
    });

    console.log(`token = ${oauthToken.access_token}`);

    const client = mqtt.connect(mqttServer, { wsOptions: { headers: { 'authorization': `Bearer ${oauthToken.access_token}` } } });
    const delay = 1000

    const sampleEvent = {
        "eventType": "User.registered",
        "cloudEventsVersion": "0.1",
        "source": "https://example.com",
        "eventTime": "2019-03-14T02:30:16Z",
        "schemaURL": "https://example.com/ODATA_SPEC/",
        "contentType": "application/json",
        "data": { "myKey": "myValue" }
    };

    client.on('connect', function () {
        console.log("connected");
        setInterval(sendMessage, delay);
    });

    client.on('close', function (err) {
        console.log("error", err.message);
    });

    function sendMessage() {
        client.publish('EXTFACTORY', JSON.stringify(sampleEvent), { qos: 1 });
        console.log('Message Sent');
    }

})();
