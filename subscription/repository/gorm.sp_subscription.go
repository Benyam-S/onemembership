package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/subscription"
	"github.com/jinzhu/gorm"
)

// SPSubscriptionRepository is a type that defines a service provider's subscription repository type
type SPSubscriptionRepository struct {
	conn *gorm.DB
}

// NewSPSubscriptionRepository is a function that creates a new service provider's subscription repository type
func NewSPSubscriptionRepository(connection *gorm.DB) subscription.ISPSubscriptionRepository {
	return &SPSubscriptionRepository{conn: connection}
}

// Create is a method that adds a new service provider subscription to the database
func (repo *SPSubscriptionRepository) Create(newSubscription *entity.SPSubscription) error {

	err := repo.conn.Create(newSubscription).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain service provider subscription from the database using an providerID,
// also Find() uses only provider_id as a key for selection
func (repo *SPSubscriptionRepository) Find(providerID string) (*entity.SPSubscription, error) {

	subscription := new(entity.SPSubscription)
	err := repo.conn.Model(subscription).Where("provider_id = ?", providerID).
		First(subscription).Error

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// FindMultiple is a method that finds multiple service provider subscriptions from the database the matches the given planID
// In FindMultiple() subscription_plan_id is only used as a key
func (repo *SPSubscriptionRepository) FindMultiple(planID string) []*entity.SPSubscription {

	var subscriptions []*entity.SPSubscription
	err := repo.conn.Model(entity.SPSubscription{}).Where("subscription_plan_id = ?", planID).
		Find(&subscriptions).Error

	if err != nil {
		return []*entity.SPSubscription{}
	}
	return subscriptions
}

// Update is a method that updates a certain service provider subscription entries in the database
func (repo *SPSubscriptionRepository) Update(subscription *entity.SPSubscription) error {

	prevSubscription := new(entity.SPSubscription)
	err := repo.conn.Model(prevSubscription).Where("provider_id = ?", subscription.ProviderID).
		First(prevSubscription).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	subscription.CreatedAt = prevSubscription.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(subscription).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete is a method that deletes a certain service provider subscription from the database using an providerID.
// In Delete() provider_id is only used as an key
func (repo *SPSubscriptionRepository) Delete(providerID string) (*entity.SPSubscription, error) {
	subscription := new(entity.SPSubscription)
	err := repo.conn.Model(subscription).Where("provider_id = ?", providerID).First(subscription).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(subscription)
	return subscription, nil
}

// DeleteMultiple is a method that deletes a set of service provider subscriptions from the database using an planID.
// In DeleteMultiple() subscription_plan_id is only used as an key
func (repo *SPSubscriptionRepository) DeleteMultiple(planID string) []*entity.SPSubscription {
	var subscriptions []*entity.SPSubscription
	repo.conn.Model(subscriptions).Where("subscription_plan_id = ?", planID).
		Find(&subscriptions)

	for _, subscription := range subscriptions {
		repo.conn.Delete(subscription)
	}

	return subscriptions
}
