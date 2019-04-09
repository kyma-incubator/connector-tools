const MAX_TOKENS = 100;

/**
 * Constructor.
 */
function InMemoryCache() {
  this.clients = [
      { 
        clientId : process.env.oauth2_client_id,
        clientSecret : process.env.oauth2_client_secret,
        redirectUris : [''],
        grants: ['client_credentials'],
      }];
  this.tokens = [];
  this.users = [{ id : '1', username: 'kyma', password: 'kyma' }];
}

/*
 * Get access token.
 */
InMemoryCache.prototype.getAccessToken = function(bearerToken) {
  console.log('called getAccessToken, bearerToken=', bearerToken);
  var tokens = this.tokens.filter(function(token) {
    return token.accessToken === bearerToken;
  });

  return tokens.length ? tokens[0] : false;
};

/**
 * Get refresh token.
 */
InMemoryCache.prototype.getRefreshToken = function(bearerToken) {
  console.log('called getRefreshToken, bearerToken=', bearerToken);
  var tokens = this.tokens.filter(function(token) {
    return token.refreshToken === bearerToken;
  });

  return tokens.length ? tokens[0] : false;
};

/**
 * Get client.
 */
InMemoryCache.prototype.getClient = function(clientId, clientSecret) {
  console.log(`called InMemoryCache.getClient - clientId=${clientId}, clientSecret=${clientSecret}`);
  var clients = this.clients.filter(function(client) {
    return client.clientId === clientId &&
           client.clientSecret === clientSecret;
  });
  return clients.length ? clients[0] : false;
};

/**
 * Save token.
 */
InMemoryCache.prototype.saveToken = function(token, client, user) {
  //console.log('called saveToken', arguments);
  var newToken = {
    accessToken: token.accessToken,
    accessTokenExpiresAt: token.accessTokenExpiresAt,
    clientId: client.clientId,
    refreshToken: token.refreshToken,
    refreshTokenExpiresAt: token.refreshTokenExpiresAt,
    userId: user.id,
    client: client,
    user:user,
    scope: null, //where are we taking scope from? maybe client?
  };

  // make sure our array of tokens never get's bigger than MAX_TOKENS
  if (this.tokens.length >= MAX_TOKENS)
  {
    this.tokens.shift();
  }
  this.tokens.push(newToken);
  return newToken;
};

/*
 * Get user.
 */
InMemoryCache.prototype.getUser = function(username, password) {
  var users = this.users.filter(function(user) {
    return user.username === username && user.password === password;
  });

  return users.length ? users[0] : false;
};

InMemoryCache.prototype.getUserFromClient = function(){
  console.log('called prototype.getUserFromClient', arguments);
  //todo: find correct user.
  return this.users[0];
}

InMemoryCache.prototype.saveAuthorizationCode = function(){
    console.log('how is this implemented!?', arguments);
}

/**
 * Export constructor.
 */
module.exports = InMemoryCache;