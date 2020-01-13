package open_connectors

type connectorInstanceElement struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type connectorInstanceResponse struct {
	ID      int64                    `json:"id"`
	Name    string                   `json:"name"`
	Token   string                   `json:"token"`
	Element connectorInstanceElement `json:"element"`
}
