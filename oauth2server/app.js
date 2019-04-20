"use strict"
const bodyParser = require('body-parser');
const express = require('express');
const OAuthServer = require('express-oauth-server');
const InMemoryCache = require('./model.js');

const port = 8080

// Create an Express application.
let app = express();

// Add body parser.
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

let model = new InMemoryCache();
// Add OAuth server.
app.oauth = new OAuthServer({
  debug: true,
  model: model
});

// Retrieve a token
app.post('/oauth/token', app.oauth.token());

// Validate a token
app.get('/validate', app.oauth.authenticate(), function(req, res) {
  res.send('OK');
});

// Start listening for requests.
app.listen(port);

console.log(`Started server on port ${port} configured with clientId ${process.env.OAUTH_CLIENT_ID}`)

