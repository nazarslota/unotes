package inputports

import (
	"github.com/udholdenhed/unotes/auth/internal/application"
	"github.com/udholdenhed/unotes/auth/internal/inputports/http"
)

type Services struct {
	HTTPServer *http.Server
}

func NewServices(services *application.Services, ho http.ServerOptions) *Services {
	s := &Services{HTTPServer: http.NewServer(services, ho)}
	return s
}
