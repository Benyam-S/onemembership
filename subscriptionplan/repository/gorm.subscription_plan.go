package repository

import (
	"fmt"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/subscriptionplan"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/jinzhu/gorm"
)

// SubscriptionPlanRepository is a type that defines a subscription plan repository type
type SubscriptionPlanRepository struct {
	conn *gorm.DB
}

// NewSubscriptionPlanRepository is a function that creates a new subscription plan repository type
func NewSubscriptionPlanRepository(connection *gorm.DB) subscriptionplan.ISubscriptionPlanRepository {
	return &SubscriptionPlanRepository{conn: connection}
}

// Create is a method that adds a new subscription plan to the database
func (repo *SubscriptionPlanRepository) Create(newSubscriptionPlan *entity.SubscriptionPlan) error {
	totalNumOfSubscriptionPlans := tools.CountMembers("subscription_plans", repo.conn)
	newSubscriptionPlan.ID = fmt.Sprintf("SBP-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptionPlans+1)

	for !tools.IsUnique("id", newSubscriptionPlan.ID, "subscription_plans", repo.conn) {
		totalNumOfSubscriptionPlans++
		newSubscriptionPlan.ID = fmt.Sprintf("SBP-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptionPlans+1)
	}

	err := repo.conn.Create(newSubscriptionPlan).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain subscription plan from the database using an subscription plan id,
// also Find() uses only id as a key for selection
func (repo *SubscriptionPlanRepository) Find(id string) (*entity.SubscriptionPlan, error) {

	subscriptionPlan := new(entity.SubscriptionPlan)
	err := repo.conn.Model(subscriptionPlan).Where("id = ?", id).First(subscriptionPlan).Error

	if err != nil {
		return nil, err
	}
	return subscriptionPlan, nil
}

// FindMultiple is a method that finds multiple subscription plans from the database the matches the given projectID
// In FindMultiple() project_id is used as a key
func (repo *SubscriptionPlanRepository) FindMultiple(projectID string) []*entity.SubscriptionPlan {

	var subscriptionPlans []*entity.SubscriptionPlan
	err := repo.conn.Model(entity.SubscriptionPlan{}).Where("project_id = ?", projectID).
		Find(&subscriptionPlans).Error

	if err != nil {
		return []*entity.SubscriptionPlan{}
	}
	return subscriptionPlans
}

// Update is a method that updates a certain subscription plan entries in the database
func (repo *SubscriptionPlanRepository) Update(subscriptionPlan *entity.SubscriptionPlan) error {

	prevSubscriptionPlan := new(entity.SubscriptionPlan)
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

// Delete is a method that deletes a certain subscription plan from the database using an subscription plan id.
// In Delete() id is only used as an key
func (repo *SubscriptionPlanRepository) Delete(id string) (*entity.SubscriptionPlan, error) {
	subscriptionPlan := new(entity.SubscriptionPlan)
	err := repo.conn.Model(subscriptionPlan).Where("id = ?", id).First(subscriptionPlan).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(subscriptionPlan)
	return subscriptionPlan, nil
}

// DeleteMultiple is a method that deletes a set of subscription plans from the database using an projectID.
// In DeleteMultiple() project_id is used as an key
func (repo *SubscriptionPlanRepository) DeleteMultiple(projectID string) []*entity.SubscriptionPlan {
	var subscriptionPlans []*entity.SubscriptionPlan
	repo.conn.Model(subscriptionPlans).Where("project_id = ?", projectID).
		Find(&subscriptionPlans)

	for _, subscriptionPlan := range subscriptionPlans {
		repo.conn.Delete(subscriptionPlan)
	}

	return subscriptionPlans
}
