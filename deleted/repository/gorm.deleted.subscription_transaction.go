package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/deleted"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/jinzhu/gorm"
)

// DeletedSubscriptionTransactionRepository is a type that defines a repository for deleted subscription transaction
type DeletedSubscriptionTransactionRepository struct {
	conn *gorm.DB
}

// NewDeletedSubscriptionTransactionRepository is a function that returns a new deleted subscription transaction repository
func NewDeletedSubscriptionTransactionRepository(connection *gorm.DB) deleted.IDeletedSubscriptionTransactionRepository {
	return &DeletedSubscriptionTransactionRepository{conn: connection}
}

// Create is a method that adds a deleted subscription transactions to the database in bulk using the userID
func (repo *DeletedSubscriptionTransactionRepository) Create(userID, prefixedUserID string) {

	repo.conn.Exec(`INSERT INTO deleted_subscription_transactions SELECT id, user_id, plan_id, app_id, receiver_name, subject, `+
		`received_amount, transaction_fee, currency_type, timeout_express, nonce, out_trade_no, status, initiated_from `+
		`FROM subscription_transactions WHERE user_id = ?`, userID)

	repo.conn.Exec("UPDATE deleted_subscription_transactions SET user_id = ? WHERE user_id = ?", prefixedUserID, userID)
}

// Find is a method that finds a certain deleted subscription transaction from the database using an id,
// also Find() uses id as a key for selection
func (repo *DeletedSubscriptionTransactionRepository) Find(id string) (*entity.DeletedSubscriptionTransaction, error) {
	deletedSubscriptionTransaction := new(entity.DeletedSubscriptionTransaction)
	err := repo.conn.Model(deletedSubscriptionTransaction).Where("id = ? ", id).
		First(deletedSubscriptionTransaction).Error

	if err != nil {
		return nil, err
	}
	return deletedSubscriptionTransaction, nil
}

// Search is a method that search and returns a set of deleted subscription transactions from the database using an id.
func (repo *DeletedSubscriptionTransactionRepository) Search(key string, pageNum int64,
	columns ...string) []*entity.DeletedSubscriptionTransaction {

	var deletedSubscriptionTransactions []*entity.DeletedSubscriptionTransaction
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_subscription_transactions WHERE ("+strings.Join(whereStmt, "||")+") "+
		"ORDER BY receiver_name ASC LIMIT ?, 30", sqlValues...).Scan(&deletedSubscriptionTransactions)

	return deletedSubscriptionTransactions
}

// SearchWRegx is a method that searchs and returns set of deleted subscription transactions limited to the key id and page number using regular expersions
func (repo *DeletedSubscriptionTransactionRepository) SearchWRegx(key string,
	pageNum int64, columns ...string) []*entity.DeletedSubscriptionTransaction {
	var deletedSubscriptionTransactions []*entity.DeletedSubscriptionTransaction
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_subscription_transactions WHERE "+strings.Join(whereStmt, "||")+
		" ORDER BY receiver_name ASC LIMIT ?, 30", sqlValues...).Scan(&deletedSubscriptionTransactions)

	return deletedSubscriptionTransactions
}

// All is a method that returns all the deleted subscription transactions from the database limited with the pageNum
func (repo *DeletedSubscriptionTransactionRepository) All(pageNum int64) []*entity.DeletedSubscriptionTransaction {

	var deletedSubscriptionTransactions []*entity.DeletedSubscriptionTransaction
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM deleted_subscription_transactions ORDER BY receiver_name ASC LIMIT ?, 30",
		limit).Scan(&deletedSubscriptionTransactions)
	return deletedSubscriptionTransactions
}

// Update is a method that updates a certain deleted subscription transaction value in the database
func (repo *DeletedSubscriptionTransactionRepository) Update(
	deletedSubscriptionTransaction *entity.DeletedSubscriptionTransaction) error {

	prevDeleteSubscriptionTransaction := new(entity.DeletedSubscriptionTransaction)
	err := repo.conn.Model(prevDeleteSubscriptionTransaction).Where("id = ?", deletedSubscriptionTransaction.ID).
		First(prevDeleteSubscriptionTransaction).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(deletedSubscriptionTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that removes a certain deleted subscription transaction from the database using an id.
// In Delete() id is only used as a key
func (repo *DeletedSubscriptionTransactionRepository) Delete(id string) (*entity.DeletedSubscriptionTransaction, error) {
	deletedSubscriptionTransaction := new(entity.DeletedSubscriptionTransaction)
	err := repo.conn.Model(deletedSubscriptionTransaction).Where("id = ?", id).First(deletedSubscriptionTransaction).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(deletedSubscriptionTransaction)
	return deletedSubscriptionTransaction, nil
}
