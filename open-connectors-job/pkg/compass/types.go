package compass

type Application struct {
	ID                       string
	Name                     string
	ConnectorInstanceID      string
	ConnectorInstanceContext string
	APIs                     *[]API
}

type API struct {
	ID      string
	Name    string
	Version string
}


