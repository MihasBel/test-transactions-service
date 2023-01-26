package pg

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	model "github.com/MihasBel/test-transactions-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	tableName      = "transactions"
	statusSuccess  = 1
	statusRejected = 2
)

type nullAmount struct{}

func (n nullAmount) Error() string {
	return "Transaction with null amount"
}

// PlaceTransaction create transaction in pg DB
func (pg *PG) PlaceTransaction(_ context.Context, tran model.Transaction) error {
	if tran.Amount == 0 {
		err := nullAmount{}
		pg.log.Error().Err(err)
		return err
	}

	err := pg.gorm.Transaction(func(tx *gorm.DB) error {
		if tran.Amount > 0 {
			tran.Status = statusSuccess
			tran.Description = "success"
			if result := tx.Create(&tran); result.Error != nil {
				return result.Error
			}
			return nil
		}
		balance := 0
		row := tx.Table("transactions").Select("COALESCE(SUM(amount), 0)").Where("user_id = ? AND status = 1", tran.UserID.String()).Row()
		if err := row.Scan(&balance); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}

		}
		if balance+tran.Amount >= 0 {
			tran.Status = statusSuccess
			tran.Description = "success"
			if result := tx.Create(&tran); result.Error != nil {
				return result.Error
			}
			return nil
		}
		tran.Status = statusRejected
		tran.Description = "rejected"
		if result := tx.Table(tableName).Create(&tran); result.Error != nil {
			return result.Error
		}
		return nil

	})
	if err != nil {
		pg.log.Error().Err(err)
	}
	return nil
}

// GetTransactionByID get transaction from pg DB
func (pg *PG) GetTransactionByID(_ context.Context, id uuid.UUID) (model.Transaction, error) {
	tran := model.Transaction{}
	if err := pg.gorm.Table(tableName).Where("id=?", id.String()).First(&tran).Error; err != nil {
		pg.log.Error().Err(err)
	}

	if tran.ID == uuid.Nil {
		err := errors.New("transaction not found")
		pg.log.Error().Err(err)
		return tran, err
	}
	return tran, nil
}
