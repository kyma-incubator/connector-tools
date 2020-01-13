package compass

type compassCreateAPIResponseId struct {
	ID string `json:"id"`
}

type compassCreateAPIResponse struct {
	AddAPI compassCreateAPIResponseId `json:"addAPI"`
}