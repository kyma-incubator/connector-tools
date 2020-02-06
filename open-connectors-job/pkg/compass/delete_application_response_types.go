package compass

type compassUnregisterApplicationResponseId struct {
	ID string `json:"id"`
}

type compassUnregisterApplicationResponse struct {
	UnregisterApplication compassUnregisterApplicationResponseId `json:"unregisterApplication"`
}