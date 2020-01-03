package compass

type compassDeleteApplicationResponseId struct {
	ID string `json:"id"`
}

type compassDeleteApplicationResponse struct {
	DeleteApplication compassDeleteApplicationResponseId `json:"deleteApplication"`
}