package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/deleted"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/jinzhu/gorm"
)

// DeletedSPPayrollTransactionRepository is a type that defines a repository for deleted service provider payroll transaction
type DeletedSPPayrollTransactionRepository struct {
	conn *gorm.DB
}

// NewDeletedSPPayrollTransactionRepository is a function that returns a new deleted service provider payroll transaction repository
func NewDeletedSPPayrollTransactionRepository(connection *gorm.DB) deleted.IDeletedSPPayrollTransactionRepository {
	return &DeletedSPPayrollTransactionRepository{conn: connection}
}

// Create is a method that adds a deleted service provider payroll transactions to the database in bulk using the providerID
func (repo *DeletedSPPayrollTransactionRepository) Create(providerID, prefixedProviderID string) {

	repo.conn.Exec(`INSERT INTO deleted_sp_payroll_transactions SELECT id, provider_id, payed_amount, linked_account, `+
		`linked_account_provider, status, FROM sp_payroll_transactions WHERE provider_id = ?`, providerID)

	repo.conn.Exec("UPDATE deleted_sp_payroll_transactions SET provider_id = ? WHERE provider_id = ?",
		prefixedProviderID, providerID)
}

// Find is a method that finds a certain deleted service provider payroll transaction from the database using an id,
// also Find() uses id as a key for selection
func (repo *DeletedSPPayrollTransactionRepository) Find(id string) (*entity.DeletedSPPayrollTransaction, error) {
	deletedSPPT := new(entity.DeletedSPPayrollTransaction)
	err := repo.conn.Model(deletedSPPT).Where("id = ? ", id).
		First(deletedSPPT).Error

	if err != nil {
		return nil, err
	}
	return deletedSPPT, nil
}

// Search is a method that search and returns a set of deleted service provider payroll transactions from the database using an id.
func (repo *DeletedSPPayrollTransactionRepository) Search(key string, pageNum int64,
	columns ...string) []*entity.DeletedSPPayrollTransaction {

	var deletedSPPTs []*entity.DeletedSPPayrollTransaction
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_sp_payroll_transactions WHERE ("+strings.Join(whereStmt, "||")+") "+
		"ORDER BY linked_account ASC LIMIT ?, 30", sqlValues...).Scan(&deletedSPPTs)

	return deletedSPPTs
}

// SearchWRegx is a method that searchs and returns set of deleted service provider payroll transactions limited to the key id and page number using regular expersions
func (repo *DeletedSPPayrollTransactionRepository) SearchWRegx(key string,
	pageNum int64, columns ...string) []*entity.DeletedSPPayrollTransaction {
	var deletedSPPTs []*entity.DeletedSPPayrollTransaction
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_sp_payroll_transactions WHERE "+strings.Join(whereStmt, "||")+
		" ORDER BY linked_account ASC LIMIT ?, 30", sqlValues...).Scan(&deletedSPPTs)

	return deletedSPPTs
}

// All is a method that returns all the deleted service provider payroll transactions from the database limited with the pageNum
func (repo *DeletedSPPayrollTransactionRepository) All(pageNum int64) []*entity.DeletedSPPayrollTransaction {

	var deletedSPPTs []*entity.DeletedSPPayrollTransaction
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM deleted_sp_payroll_transactions ORDER BY linked_account ASC LIMIT ?, 30",
		limit).Scan(&deletedSPPTs)
	return deletedSPPTs
}

// Update is a method that updates a certain deleted service provider payroll transaction value in the database
func (repo *DeletedSPPayrollTransactionRepository) Update(
	deletedSPPT *entity.DeletedSPPayrollTransaction) error {

	prevDeleteSPST := new(entity.DeletedSPPayrollTransaction)
	err := repo.conn.Model(prevDeleteSPST).Where("id = ?", deletedSPPT.ID).First(prevDeleteSPST).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(deletedSPPT).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that removes a certain deleted service provider payroll transaction from the database using an id.
// In Delete() id is only used as a key
func (repo *DeletedSPPayrollTransactionRepository) Delete(id string) (*entity.DeletedSPPayrollTransaction, error) {
	deletedSPPT := new(entity.DeletedSPPayrollTransaction)
	err := repo.conn.Model(deletedSPPT).Where("id = ?", id).First(deletedSPPT).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(deletedSPPT)
	return deletedSPPT, nil
}
