package repository

import (
	"fmt"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/Benyam-S/onemembership/transaction"
	"github.com/jinzhu/gorm"
)

// SPPayrollTransactionRepository is a type that defines a service provider payroll transaction repository type
type SPPayrollTransactionRepository struct {
	conn *gorm.DB
}

// NewSPPayrollTransactionRepository is a function that creates a new service provider payroll transaction repository type
func NewSPPayrollTransactionRepository(connection *gorm.DB) transaction.ISPPayrollTransactionRepository {
	return &SPPayrollTransactionRepository{conn: connection}
}

// Create is a method that adds a new service provider payroll transaction to the database
func (repo *SPPayrollTransactionRepository) Create(newPayrollTransaction *entity.SPPayrollTransaction) error {
	totalNumOfPayrollTransactions := tools.CountMembers("sp_payroll_transactions", repo.conn)
	newPayrollTransaction.ID = fmt.Sprintf("PRT-%s%d", tools.RandomStringGN(20), totalNumOfPayrollTransactions+1)

	for !tools.IsUnique("id", newPayrollTransaction.ID, "sp_payroll_transactions", repo.conn) {
		totalNumOfPayrollTransactions++
		newPayrollTransaction.ID = fmt.Sprintf("PRT-%s%d", tools.RandomStringGN(20), totalNumOfPayrollTransactions+1)
	}

	err := repo.conn.Create(newPayrollTransaction).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain service provider payroll transaction from the database using an transaction id,
// also Find() uses only id as a key for selection
func (repo *SPPayrollTransactionRepository) Find(id string) (*entity.SPPayrollTransaction, error) {

	payrollTransaction := new(entity.SPPayrollTransaction)
	err := repo.conn.Model(payrollTransaction).Where("id = ?", id).
		First(payrollTransaction).Error

	if err != nil {
		return nil, err
	}

	return payrollTransaction, nil
}

// FindMultiple is a method that finds multiple service provider payroll transactions from the database the matches the given providerID
// In FindMultiple() provider_id is only used as a key
func (repo *SPPayrollTransactionRepository) FindMultiple(providerID string) []*entity.SPPayrollTransaction {

	var payrollTransactions []*entity.SPPayrollTransaction
	err := repo.conn.Model(entity.SPPayrollTransaction{}).Where("provider_id = ?", providerID).
		Find(&payrollTransactions).Error

	if err != nil {
		return []*entity.SPPayrollTransaction{}
	}
	return payrollTransactions
}

// Update is a method that updates a certain service provider payroll transaction entries in the database
func (repo *SPPayrollTransactionRepository) Update(payrollTransaction *entity.SPPayrollTransaction) error {

	prevPayrollTransaction := new(entity.SPPayrollTransaction)
	err := repo.conn.Model(prevPayrollTransaction).Where("id = ?", payrollTransaction.ID).
		First(prevPayrollTransaction).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	payrollTransaction.CreatedAt = prevPayrollTransaction.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(payrollTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain service provider payroll transaction from the database using an transaction id.
// In Delete() id is only used as an key
func (repo *SPPayrollTransactionRepository) Delete(id string) (*entity.SPPayrollTransaction, error) {
	payrollTransaction := new(entity.SPPayrollTransaction)
	err := repo.conn.Model(payrollTransaction).Where("id = ?", id).
		First(payrollTransaction).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(payrollTransaction)
	return payrollTransaction, nil
}

// DeleteMultiple is a method that deletes a set of service provider payroll transactions from the database using an providerID.
// In DeleteMultiple() provider_id is only used as a key
func (repo *SPPayrollTransactionRepository) DeleteMultiple(providerID string) []*entity.SPPayrollTransaction {
	var payrollTransactions []*entity.SPPayrollTransaction
	repo.conn.Model(payrollTransactions).Where("provider_id = ?", providerID).Find(&payrollTransactions)

	for _, payrollTransaction := range payrollTransactions {
		repo.conn.Delete(payrollTransaction)
	}

	return payrollTransactions
}
