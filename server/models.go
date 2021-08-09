package server

type Setter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Pusher struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Puller struct {
	Name string `json:"name"`
}
