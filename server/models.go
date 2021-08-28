package server

type Pusher struct {
	Name    string `json:"name" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type Puller struct {
	Name string `json:"name" validate:"required"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type Client struct{}
