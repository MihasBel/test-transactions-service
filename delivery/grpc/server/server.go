package server

import (
	"context"
	"encoding/json"
	v1transaction "github.com/MihasBel/test-transactions-service/delivery/grpc/gen/v1/transaction"
	"github.com/MihasBel/test-transactions-service/delivery/grpc/server/transaction"
	"github.com/MihasBel/test-transactions-service/internal/rep"
	"github.com/rs/zerolog"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/pkg/errors"
)

const (
	TCP = "tcp"

	MsgStartListening = "start listening grpc at %v"
	MsgStopListening  = "stop listening grpc at %v"
	MsgServerFailed   = "server failed"

	MsgErrFailedListen = "failed to listen GRPC server: %w"

	KeyLoggerDirection  = "direction"
	KeyLoggerGRPCStatus = "grpc_status"
	KeyLoggerDuration   = "duration"
	KeyLoggerAnswer     = "answer"
	KeyLoggerError      = "error"

	ValLoggerDirection = "in"

	KeyMetrics = "metrics"
)

type Server struct {
	cfg Config
	srv *grpc.Server
	s   rep.Storage
	l   zerolog.Logger
}

func New(
	cfg Config,
	s rep.Storage,
	l zerolog.Logger,
) *Server {
	return &Server{
		cfg: cfg,
		s:   s,
		l:   l,
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen(TCP, s.cfg.Address)
	if err != nil {
		return errors.Wrap(err, MsgErrFailedListen)
	}

	s.srv = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			s.ZerologInterceptor,
			recovery.UnaryServerInterceptor(),
		)),
	)
	v1transaction.RegisterTransactionAPIServer(s.srv, transaction.New(s.s))

	reflection.Register(s.srv)

	s.l.Info().Msgf(MsgStartListening, s.cfg.Address)
	errCh := make(chan error)
	go func() {
		if err := s.srv.Serve(lis); err != nil {
			s.l.Error().Err(err).Msg(MsgServerFailed)
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(time.Duration(s.cfg.StartTimeout) * time.Second):
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.l.Info().Msgf(MsgStopListening, s.cfg.Address)
	stopCh := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		stopCh <- struct{}{}
	}()
	select {
	case <-time.After(time.Duration(s.cfg.StopTimeout) * time.Second):
		return nil
	case <-stopCh:
		return nil
	}
}

func (s *Server) ZerologInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	body, _ := json.Marshal(resp)
	st, _ := status.FromError(err)

	lgr := s.l.Info().
		Str(KeyLoggerDirection, ValLoggerDirection).
		Str(KeyLoggerGRPCStatus, st.Code().String()).
		Dur(KeyLoggerDuration, time.Since(start)).
		Str(KeyLoggerAnswer, string(body))

	if err != nil {
		lgr.Str(KeyLoggerError, err.Error())
	}

	lgr.Msg(info.FullMethod)

	return resp, err
}
