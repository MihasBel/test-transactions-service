package pg

import (
	"context"

	model "github.com/MihasBel/test-transactions-servise/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	tableName = "transactions"
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
	anyTran := model.Transaction{}
	emptyUUID := uuid.UUID{}
	pg.gorm.Table("transactions").Where("user_id = ? AND status = 1", tran.UserID.String()).First(&anyTran)
	if anyTran.ID == emptyUUID {
		return transactionWithNullBalance(pg.gorm, tran)
	}
	err := pg.gorm.Transaction(func(tx *gorm.DB) error {
		if tran.Amount > 0 {
			tran.Status = 1
			tran.Description = "success"
			if result := tx.Create(&tran); result.Error != nil {
				return result.Error
			}
			return nil
		}
		balance := 0
		row := tx.Table("transactions").Select("SUM(amount)").Where("user_id = ? AND status = 1", tran.UserID.String()).Row()
		if err := row.Scan(&balance); err != nil {
			return err
		}
		if balance+tran.Amount > 0 {
			tran.Status = 1
			tran.Description = "success"
			if result := tx.Create(&tran); result.Error != nil {
				return result.Error
			}
			return nil
		}
		tran.Status = 2
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
func transactionWithNullBalance(db *gorm.DB, tran model.Transaction) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if tran.Amount > 0 {
			tran.Status = 1
			tran.Description = "success"
			if result := tx.Create(&tran); result.Error != nil {
				return result.Error
			}
			return nil
		}
		tran.Status = 2
		tran.Description = "rejected"
		if result := tx.Table(tableName).Create(&tran); result.Error != nil {
			return result.Error
		}
		return nil

	})
	return err
}

// GetTransactionByID get transaction from pg DB
func (pg *PG) GetTransactionByID(_ context.Context, id uuid.UUID) (model.Transaction, error) {
	tran := model.Transaction{}
	if err := pg.gorm.Table(tableName).Where("id=?", id.String()).First(&tran).Error; err != nil {
		pg.log.Error().Err(err)
	}
	return tran, nil
}
