package repository

import (
	"fmt"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/Benyam-S/onemembership/transaction"
	"github.com/jinzhu/gorm"
)

// SPSubscriptionTransactionRepository is a type that defines a service provider subscription transaction repository type
type SPSubscriptionTransactionRepository struct {
	conn *gorm.DB
}

// NewSPSubscriptionTransactionRepository is a function that creates a new service provider subscription transaction repository type
func NewSPSubscriptionTransactionRepository(connection *gorm.DB) transaction.ISPSubscriptionTransactionRepository {
	return &SPSubscriptionTransactionRepository{conn: connection}
}

// Create is a method that adds a new service provider subscription transaction to the database
func (repo *SPSubscriptionTransactionRepository) Create(newSubscriptionTransaction *entity.SPSubscriptionTransaction) error {
	totalNumOfSubscriptionTransactions := tools.CountMembers("sp_subscription_transactions", repo.conn)
	newSubscriptionTransaction.ID = fmt.Sprintf("SBT-%s%d", tools.RandomStringGN(20), totalNumOfSubscriptionTransactions+1)

	for !tools.IsUnique("id", newSubscriptionTransaction.ID, "sp_subscription_transactions", repo.conn) {
		totalNumOfSubscriptionTransactions++
		newSubscriptionTransaction.ID = fmt.Sprintf("SBT-%s%d", tools.RandomStringGN(20), totalNumOfSubscriptionTransactions+1)
	}

	err := repo.conn.Create(newSubscriptionTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain service provider subscription transaction from the database using an identifier,
// also Find() uses id and out_trade_no as a key for selection
func (repo *SPSubscriptionTransactionRepository) Find(identifier string) (*entity.SPSubscriptionTransaction, error) {

	subscriptionTransaction := new(entity.SPSubscriptionTransaction)
	err := repo.conn.Model(subscriptionTransaction).Where("id = ? || out_trade_no = ?", identifier, identifier).
		First(subscriptionTransaction).Error

	if err != nil {
		return nil, err
	}
	return subscriptionTransaction, nil
}

// FindMultiple is a method that finds multiple service provider subscription transactions from the database the matches the given identifier
// In FindMultiple() provider_id, plan_id and app_id is used as a key
func (repo *SPSubscriptionTransactionRepository) FindMultiple(identifier string) []*entity.SPSubscriptionTransaction {

	var subscriptionTransactions []*entity.SPSubscriptionTransaction
	err := repo.conn.Model(entity.SPSubscriptionTransaction{}).Where("provider_id = ? || plan_id = ? || app_id",
		identifier, identifier, identifier).Find(&subscriptionTransactions).Error

	if err != nil {
		return []*entity.SPSubscriptionTransaction{}
	}
	return subscriptionTransactions
}

// Update is a method that updates a certain service provider subscription transaction entries in the database
func (repo *SPSubscriptionTransactionRepository) Update(subscriptionTransaction *entity.SPSubscriptionTransaction) error {

	prevSubscriptionTransaction := new(entity.SPSubscriptionTransaction)
	err := repo.conn.Model(prevSubscriptionTransaction).Where("id = ?", subscriptionTransaction.ID).
		First(prevSubscriptionTransaction).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	subscriptionTransaction.CreatedAt = prevSubscriptionTransaction.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(subscriptionTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain service provider subscription transaction from the database using an transaction id.
// In Delete() id is only used as an key
func (repo *SPSubscriptionTransactionRepository) Delete(id string) (*entity.SPSubscriptionTransaction, error) {
	subscriptionTransaction := new(entity.SPSubscriptionTransaction)
	err := repo.conn.Model(subscriptionTransaction).Where("id = ?", id).
		First(subscriptionTransaction).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(subscriptionTransaction)
	return subscriptionTransaction, nil
}

// DeleteMultiple is a method that deletes a set of subscription transactions from the database using an identifier.
// In DeleteMultiple() provider_id, plan_id and app_id is used as a key
func (repo *SPSubscriptionTransactionRepository) DeleteMultiple(identifier string) []*entity.SPSubscriptionTransaction {
	var subscriptionTransactions []*entity.SPSubscriptionTransaction
	repo.conn.Model(subscriptionTransactions).Where("provider_id = ? || plan_id = ? || app_id",
		identifier, identifier, identifier).Find(&subscriptionTransactions)

	for _, subscriptionTransaction := range subscriptionTransactions {
		repo.conn.Delete(subscriptionTransaction)
	}

	return subscriptionTransactions
}
