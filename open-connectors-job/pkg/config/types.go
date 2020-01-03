package config

type JobConfig struct {
	LogLevel       string               `json:"logLevel"` //Log level
	Compass        CompassConfig        `json:"compass"`
	OpenConnectors OpenConnectorsConfig `json:"openConnectors"`
}

type CompassConfig struct {
	DirectorURL       string `json:"directorURL"`       //Url pointing to OAuth2 secured Director URL of Compass
	TenantId          string `json:"tenantId"`          //Url pointing to OAuth2 secured Director URL of Compass
	TokenUrl          string `json:"tokenUrl"`          //Url pointing towards client credentials token provider
	ClientID          string `json:"clientID"`          //Client id of the Compass Integration System
	ClientSecret      string `json:"clientSecret"`      //Client Secret of the Compass Integration System
	ApplicationPrefix string `json:"applicationPrefix"` //Prefix assigned to applications created
	TimeoutMills      int    `json:"timeoutMills"`      //Timeout in milliseconds for Compass API Calls
}

type OpenConnectorsConfig struct {
	Hostname           string   `json:"hostname"`           //Hostname of the Open Connectors instance
	OrganizationSecret string   `json:"organizationSecret"` //Organization Secret of the Open Connectors instance
	UserSecret         string   `json:"userSecret"`         //User Secret of the Open Connectors instance
	TimeoutMills       int      `json:"timeoutMills"`       //Timeout in milliseconds for Open Connectors API Calls
	Tags               []string `json:"tags"`               //Tags to select instances in open connectors, empty means
	//all
}
