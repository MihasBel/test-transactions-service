package transaction

import (
	"github.com/MihasBel/test-transactions-servise/internal/rep"
)

type Server struct {
	s rep.Storage
}

func New(s rep.Storage) *Server {
	return &Server{
		s: s,
	}
}
