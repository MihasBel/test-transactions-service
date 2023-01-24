package transaction

import (
	"github.com/MihasBel/test-transactions-service/internal/rep"
)

// Server storage
type Server struct {
	s rep.Storage
}

// New constructor
func New(s rep.Storage) *Server {
	return &Server{
		s: s,
	}
}
