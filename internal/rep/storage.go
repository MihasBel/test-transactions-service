package rep

import (
	"context"

	model "github.com/MihasBel/test-transactions-servise/models"
	"github.com/google/uuid"
)

// Storage stores and manipulates with transactions data
//
//go:generate mockgen -source=storage.go -destination=../../mocks/storage.go -package=mocks
type Storage interface {
	PlaceTransaction(ctx context.Context, transaction model.Transaction) error
	GetTransactionByID(ctx context.Context, id uuid.UUID) (model.Transaction, error)
}
