package compass

type compassRegisterApplicationResponseId struct {
	ID string `json:"id"`
}

type compassRegisterApplicationResponse struct {
	RegisterApplication compassRegisterApplicationResponseId `json:"registerApplication"`
}