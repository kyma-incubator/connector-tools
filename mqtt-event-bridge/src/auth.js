
const request = require("request")

const oauthEndpoint = process.env.OAUTH_ENDPOINT

/**
 *
 * @param {string} oauthClientId - client_id for the Oauth2 server that we will deploy.
 * @param {string} oauthClientSecret - client_secret for the Oauth2 server
 * @returns {Promise} Promise object represent the token returned by server
 *  */
function getToken(oauthClientId, oauthClientSecret) {
    return new Promise((resolve, reject) => {
        let registrationURL = `http://application-registry-external-api.kyma-integration.svc.cluster.local:8081/${applicationName}/v1/metadata/services` //TODO: change this to a oauth server

        let info = { "client_id": oauthClientId, "client_secret": oauthClientSecret }
        request.post({
            url: registrationURL,
            headers: {
                "Content-Type": "application/json"
            },
            json: info
        }, (error, httpResponse, body) =>
                error ? reject(error) : resolve(JSON.parse(body).token)
        )
    })
}