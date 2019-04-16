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
 }

/*
 * Get access token.
 */
InMemoryCache.prototype.getAccessToken = function(bearerToken) {
  var tokens = this.tokens.filter(function(token) {
    return token.accessToken === bearerToken;
  });

  return tokens.length ? tokens[0] : false;
};

/**
 * Get refresh token.
 */
InMemoryCache.prototype.getRefreshToken = function(bearerToken) {
  var tokens = this.tokens.filter(function(token) {
    return token.refreshToken === bearerToken;
  });

  return tokens.length ? tokens[0] : false;
};

/**
 * Get client.
 */
InMemoryCache.prototype.getClient = function(clientId, clientSecret) {
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
  var newToken = {
    accessToken: token.accessToken,
    accessTokenExpiresAt: token.accessTokenExpiresAt,
    clientId: client.clientId,
    refreshToken: token.refreshToken,
    refreshTokenExpiresAt: token.refreshTokenExpiresAt,
    userId: user.id,
    client: client,
    user:user,
    scope: null, 
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
  return {}
};

InMemoryCache.prototype.getUserFromClient = function(){
  return {};
}

InMemoryCache.prototype.saveAuthorizationCode = function(){
    // do nothing
}

/**
 * Export constructor.
 */
module.exports = InMemoryCache;