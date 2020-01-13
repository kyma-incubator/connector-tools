package compass

type compassApplicationResponseAPIVersion struct {
	Value string `json:"value"`
}

type compassApplicationResponseAPI struct {
	ID      string                               `json:"id"`
	Name    string                               `json:"name"`
	Version compassApplicationResponseAPIVersion `json:"version"`
}

type compassApplicationResponseAPIData struct {
	Data []compassApplicationResponseAPI `json:"data"`
}

type compassApplicationResponseLabels struct {
	OpenConnectors compassApplicationResponseLabel `json:"open_connectors"`
}

type compassApplicationResponseLabel struct {
	ConnectorInstanceID      string `json:"connectorInstanceID"`
	ConnectorInstanceContext string `json:"connectorInstanceContext"`
}

type compassApplicationResponseApplicationData struct {
	ID     string                            `json:"id"`
	Name   string                            `json:"name"`
	Labels compassApplicationResponseLabels  `json:"labels"`
	APIs   compassApplicationResponseAPIData `json:"apis"`
}

type compassApplicationResponseApplication struct {
	Data []compassApplicationResponseApplicationData `json:"data"`
}

type compassApplicationResponseApplications struct {
	Applications compassApplicationResponseApplication `json:"applications"`
}
