package transaction

import (
	"context"
	"github.com/MihasBel/test-transactions-service/delivery/grpc/gen/v1/transaction"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) ByID(ctx context.Context, r *transaction.ByIDRequest) (*transaction.Transaction, error) {
	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return nil, err
	}
	tran, err := s.s.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &transaction.Transaction{
		Id:          tran.ID.String(),
		UserId:      tran.UserID.String(),
		Amount:      int64(tran.Amount),
		CreatedAt:   timestamppb.New(tran.CreatedAt),
		Status:      int32(tran.Status),
		Description: tran.Description,
	}, nil
}
