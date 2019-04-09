var bodyParser = require('body-parser');
var express = require('express');
var OAuthServer = require('express-oauth-server');
var InMemoryCache = require('./model.js');
const http = require("http");
const https = require("https");

// Create an Express application.
var app = express();

// Add body parser.
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

var model = new InMemoryCache();
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
app.listen(8080);

const server = http.createServer().listen(8090);

