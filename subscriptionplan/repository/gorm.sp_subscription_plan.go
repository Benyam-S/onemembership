package repository

import (
	"fmt"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/subscriptionplan"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/jinzhu/gorm"
)

// SPSubscriptionPlanRepository is a type that defines a service provider's subscription plan repository type
type SPSubscriptionPlanRepository struct {
	conn *gorm.DB
}

// NewSPSubscriptionPlanRepository is a function that creates a new service provider's subscription plan repository type
func NewSPSubscriptionPlanRepository(connection *gorm.DB) subscriptionplan.ISPSubscriptionPlanRepository {
	return &SPSubscriptionPlanRepository{conn: connection}
}

// Create is a method that adds a new service provider subscription plan to the database
func (repo *SPSubscriptionPlanRepository) Create(newSubscriptionPlan *entity.SPSubscriptionPlan) error {
	totalNumOfSubscriptionPlans := tools.CountMembers("sp_subscription_plans", repo.conn)
	newSubscriptionPlan.ID = fmt.Sprintf("SBP-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptionPlans+1)

	for !tools.IsUnique("id", newSubscriptionPlan.ID, "sp_subscription_plans", repo.conn) {
		totalNumOfSubscriptionPlans++
		newSubscriptionPlan.ID = fmt.Sprintf("SBP-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptionPlans+1)
	}

	err := repo.conn.Create(newSubscriptionPlan).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain servie provider subscription plan from the database using an id,
// also Find() uses only id as a key for selection
func (repo *SPSubscriptionPlanRepository) Find(id string) (*entity.SPSubscriptionPlan, error) {

	subscriptionPlan := new(entity.SPSubscriptionPlan)
	err := repo.conn.Model(subscriptionPlan).Where("id = ?", id).First(subscriptionPlan).Error

	if err != nil {
		return nil, err
	}

	return subscriptionPlan, nil
}

// Update is a method that updates a certain service provider subscription plan entries in the database
func (repo *SPSubscriptionPlanRepository) Update(subscriptionPlan *entity.SPSubscriptionPlan) error {

	prevSubscriptionPlan := new(entity.SPSubscriptionPlan)
	err := repo.conn.Model(prevSubscriptionPlan).Where("id = ?", subscriptionPlan.ID).
		First(prevSubscriptionPlan).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	subscriptionPlan.CreatedAt = prevSubscriptionPlan.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(subscriptionPlan).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain service provider subscription plan from the database using an id.
// In Delete() id is only used as an key
func (repo *SPSubscriptionPlanRepository) Delete(id string) (*entity.SPSubscriptionPlan, error) {
	subscriptionPlan := new(entity.SPSubscriptionPlan)
	err := repo.conn.Model(subscriptionPlan).Where("id = ?", id).First(subscriptionPlan).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(subscriptionPlan)
	return subscriptionPlan, nil
}
