package config

type ServerOptions struct {
	Port string `json:"port" mapstructure:"port"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		Port: "18080",
	}
}
