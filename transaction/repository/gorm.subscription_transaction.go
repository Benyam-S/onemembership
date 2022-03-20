package repository

import (
	"fmt"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/Benyam-S/onemembership/transaction"
	"github.com/jinzhu/gorm"
)

// SubscriptionTransactionRepository is a type that defines a subscription transaction repository type
type SubscriptionTransactionRepository struct {
	conn *gorm.DB
}

// NewSubscriptionTransactionRepository is a function that creates a new subscription transaction repository type
func NewSubscriptionTransactionRepository(connection *gorm.DB) transaction.ISubscriptionTransactionRepository {
	return &SubscriptionTransactionRepository{conn: connection}
}

// Create is a method that adds a new subscription transaction to the database
func (repo *SubscriptionTransactionRepository) Create(newSubscriptionTransaction *entity.SubscriptionTransaction) error {
	totalNumOfSubscriptionTransactions := tools.CountMembers("subscription_transactions", repo.conn)
	newSubscriptionTransaction.ID = fmt.Sprintf("SBT-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptionTransactions+1)

	for !tools.IsUnique("id", newSubscriptionTransaction.ID, "subscription_transactions", repo.conn) {
		totalNumOfSubscriptionTransactions++
		newSubscriptionTransaction.ID = fmt.Sprintf("SBT-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptionTransactions+1)
	}

	err := repo.conn.Create(newSubscriptionTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain subscription transaction from the database using an identifier,
// also Find() uses id and out_trade_no as a key for selection
func (repo *SubscriptionTransactionRepository) Find(identifier string) (*entity.SubscriptionTransaction, error) {

	subscriptionTransaction := new(entity.SubscriptionTransaction)
	err := repo.conn.Model(subscriptionTransaction).Where("id = ? || out_trade_no = ?", identifier, identifier).
		First(subscriptionTransaction).Error

	if err != nil {
		return nil, err
	}
	return subscriptionTransaction, nil
}

// FindMultiple is a method that finds multiple subscription transactions from the database the matches the given identifier
// In FindMultiple() user_id, plan_id and app_id is used as a key
func (repo *SubscriptionTransactionRepository) FindMultiple(identifier string) []*entity.SubscriptionTransaction {

	var subscriptionTransactions []*entity.SubscriptionTransaction
	err := repo.conn.Model(entity.SubscriptionTransaction{}).Where("user_id = ? || plan_id = ? || app_id",
		identifier, identifier, identifier).Find(&subscriptionTransactions).Error

	if err != nil {
		return []*entity.SubscriptionTransaction{}
	}
	return subscriptionTransactions
}

// Update is a method that updates a certain subscription transaction entries in the database
func (repo *SubscriptionTransactionRepository) Update(subscriptionTransaction *entity.SubscriptionTransaction) error {

	prevSubscriptionTransaction := new(entity.SubscriptionTransaction)
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

// Delete is a method that deletes a certain subscription transaction from the database using an transaction id.
// In Delete() id is only used as an key
func (repo *SubscriptionTransactionRepository) Delete(id string) (*entity.SubscriptionTransaction, error) {
	subscriptionTransaction := new(entity.SubscriptionTransaction)
	err := repo.conn.Model(subscriptionTransaction).Where("id = ?", id).
		First(subscriptionTransaction).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(subscriptionTransaction)
	return subscriptionTransaction, nil
}

// DeleteMultiple is a method that deletes a set of subscription transactions from the database using an identifier.
// In DeleteMultiple() user_id, plan_id and app_id is used as a key
func (repo *SubscriptionTransactionRepository) DeleteMultiple(identifier string) []*entity.SubscriptionTransaction {
	var subscriptionTransactions []*entity.SubscriptionTransaction
	repo.conn.Model(subscriptionTransactions).Where("user_id = ? || plan_id = ? || app_id",
		identifier, identifier, identifier).Find(&subscriptionTransactions)

	for _, subscriptionTransaction := range subscriptionTransactions {
		repo.conn.Delete(subscriptionTransaction)
	}

	return subscriptionTransactions
}
