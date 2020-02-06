# kyma-oauth-proxy
This provides support for APIs that require a request body form to be sent during an oauth token request.  It does not manage the tokens in any way, this should rely on the kyma application connector.  When registering an oauth api in kyma this should be used as the oauth endpoint.  An  example can be found in the addons folder. 


## Parmaters
| Name             | Description                                   | Type    | Default Value | Example                       |
| ---------------- | :-------------------------------------------- | :------ | :------------ | :---------------------------- |
| oauthURL         | The url of the oauth request                  | string  |               | https://myoauthurl.com        |
| dumpRequest      | Will output the request to the logs           | boolean | false         |                               |
| removeAuthHeader | Will remove the received Authorization header | boolean | true          |                               |
| headers          | comma seperated listing of headers            | map     |               | header1=Value1,Header2=Value2 |
| requestBodyForm  | comma seperated listing of form parameters    | map     |               | username=user,password=pw     |

## Example Usage
go run cmd/kyma-oauth-proxy/main.go --headers Content-Type=application/x-www-form-urlencoded --requestBodyForm grant_type=password,username=user,password=pw,domain=mydomain --oauthURL https://httpbin.org/anything

curl -X POST http://localhost:8080/