package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/deleted"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/jinzhu/gorm"
)

// DeletedSPSubscriptionTransactionRepository is a type that defines a repository for deleted service provider subscription transaction
type DeletedSPSubscriptionTransactionRepository struct {
	conn *gorm.DB
}

// NewDeletedSPSubscriptionTransactionRepository is a function that returns a new deleted service provider subscription transaction repository
func NewDeletedSPSubscriptionTransactionRepository(connection *gorm.DB) deleted.IDeletedSPSubscriptionTransactionRepository {
	return &DeletedSPSubscriptionTransactionRepository{conn: connection}
}

// Create is a method that adds a deleted service provider subscription transactions to the database in bulk using the providerID
func (repo *DeletedSPSubscriptionTransactionRepository) Create(providerID, prefixedProviderID string) {

	repo.conn.Exec(`INSERT INTO deleted_sp_subscription_transactions SELECT id, provider_id, plan_id, app_id, receiver_name, `+
		`subject, received_amount, transaction_fee, currency_type, timeout_express, nonce, out_trade_no, status, `+
		`FROM sp_subscription_transactions WHERE provider_id = ?`, providerID)

	repo.conn.Exec("UPDATE deleted_sp_subscription_transactions SET provider_id = ? WHERE provider_id = ?",
		prefixedProviderID, providerID)
}

// Find is a method that finds a certain deleted service provider subscription transaction from the database using an id,
// also Find() uses id as a key for selection
func (repo *DeletedSPSubscriptionTransactionRepository) Find(id string) (*entity.DeletedSPSubscriptionTransaction, error) {
	deletedSPST := new(entity.DeletedSPSubscriptionTransaction)
	err := repo.conn.Model(deletedSPST).Where("id = ? ", id).
		First(deletedSPST).Error

	if err != nil {
		return nil, err
	}
	return deletedSPST, nil
}

// Search is a method that search and returns a set of deleted service provider subscription transactions from the database using an id.
func (repo *DeletedSPSubscriptionTransactionRepository) Search(key string, pageNum int64,
	columns ...string) []*entity.DeletedSPSubscriptionTransaction {

	var deletedSPSTs []*entity.DeletedSPSubscriptionTransaction
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_sp_subscription_transactions WHERE ("+strings.Join(whereStmt, "||")+") "+
		"ORDER BY receiver_name ASC LIMIT ?, 30", sqlValues...).Scan(&deletedSPSTs)

	return deletedSPSTs
}

// SearchWRegx is a method that searchs and returns set of deleted service provider subscription transactions limited to the key id and page number using regular expersions
func (repo *DeletedSPSubscriptionTransactionRepository) SearchWRegx(key string,
	pageNum int64, columns ...string) []*entity.DeletedSPSubscriptionTransaction {
	var deletedSPSTs []*entity.DeletedSPSubscriptionTransaction
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_sp_subscription_transactions WHERE "+strings.Join(whereStmt, "||")+
		" ORDER BY receiver_name ASC LIMIT ?, 30", sqlValues...).Scan(&deletedSPSTs)

	return deletedSPSTs
}

// All is a method that returns all the deleted service provider subscription transactions from the database limited with the pageNum
func (repo *DeletedSPSubscriptionTransactionRepository) All(pageNum int64) []*entity.DeletedSPSubscriptionTransaction {

	var deletedSPSTs []*entity.DeletedSPSubscriptionTransaction
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM deleted_sp_subscription_transactions ORDER BY receiver_name ASC LIMIT ?, 30",
		limit).Scan(&deletedSPSTs)
	return deletedSPSTs
}

// Update is a method that updates a certain deleted service provider subscription transaction value in the database
func (repo *DeletedSPSubscriptionTransactionRepository) Update(
	deletedSPST *entity.DeletedSPSubscriptionTransaction) error {

	prevDeleteSPST := new(entity.DeletedSPSubscriptionTransaction)
	err := repo.conn.Model(prevDeleteSPST).Where("id = ?", deletedSPST.ID).First(prevDeleteSPST).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(deletedSPST).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete is a method that removes a certain deleted service provider subscription transaction from the database using an id.
// In Delete() id is only used as a key
func (repo *DeletedSPSubscriptionTransactionRepository) Delete(id string) (*entity.DeletedSPSubscriptionTransaction, error) {
	deletedSPST := new(entity.DeletedSPSubscriptionTransaction)
	err := repo.conn.Model(deletedSPST).Where("id = ?", id).First(deletedSPST).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(deletedSPST)
	return deletedSPST, nil
}
