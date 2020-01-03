package compass

type compassCreateApplicationResponseId struct {
	ID string `json:"id"`
}

type compassCreateApplicationResponse struct {
	CreateApplication compassCreateApplicationResponseId `json:"createApplication"`
}