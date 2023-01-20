package domain

// HTTPServerConfig, config for server to serve http connections
type HTTPServerConfig struct {
	Host string
	Port int
}

type RequestModel[T any] struct {
	Data *T `json:"data,omitempty"`
}

type ResponseModel[T any] struct {
	Data  *T     `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}
