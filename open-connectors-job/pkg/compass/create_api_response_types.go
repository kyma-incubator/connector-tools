package compass

type compassAddAPIResponseId struct {
	ID string `json:"id"`
}

type compassAddAPIResponse struct {
	AddAPIDefinition compassAddAPIResponseId `json:"addAPIDefinition"`
}