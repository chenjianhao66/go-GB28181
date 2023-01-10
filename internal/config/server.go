package config

type ServerOptions struct {
	Port string `json:"port"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		Port: "18080",
	}
}
