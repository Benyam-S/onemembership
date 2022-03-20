package repository

import (
	"fmt"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/subscription"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/jinzhu/gorm"
)

// SubscriptionRepository is a type that defines a subscription repository type
type SubscriptionRepository struct {
	conn *gorm.DB
}

// NewSubscriptionRepository is a function that creates a new subscription repository type
func NewSubscriptionRepository(connection *gorm.DB) subscription.ISubscriptionRepository {
	return &SubscriptionRepository{conn: connection}
}

// Construct is a method that construct a subscription from given collection of identifiers
func (repo *SubscriptionRepository) Construct(subscriberID, subscriptionPlanID string) (*entity.Subscription, error) {
	sqlStatement := `SELECT A.id AS subscriber_id, A.first_name AS subscriber_first_name, ` +
		`A.last_name AS subscriber_last_name, A.user_name AS subscriber_user_name, ` +
		`A.phone_number AS subscriber_phone_number, A.email AS subscriber_email, ` +
		`B.id AS provider_id, B.first_name AS provider_first_name, B.last_name AS provider_last_name, ` +
		`B.user_name AS provider_user_name, B.phone_number AS provider_phone_number, B.email AS provider_email, ` +
		`C.id AS project_id, C.name AS project_name, C.description AS project_description, C.project_link As project_link, ` +
		`D.id AS subscription_plan_id, D.name AS subscription_plan_name, D.benfits AS subscription_plan_benfits, ` +
		`D.duration AS subscription_plan_duration, D.price AS subscription_plan_price, ` +
		`D.is_recurring AS subscription_plan_is_recurring, D.currency AS subscription_plan_currency ` +
		`FROM users A, ((subscription_plans D INNER JOIN projects C ON D.project_id = C.id) ` +
		`INNER JOIN service_providers B ON c.provider_id = B.id) WHERE A.id = ? AND D.id = ?`

	type SubscriptionData struct {
		SubscriberID          string
		SubscriberFirstName   string
		SubscriberLastName    string
		SubscriberUserName    string
		SubscriberPhoneNumber string
		SubscriberEmail       string

		ProviderID          string
		ProviderFirstName   string
		ProviderLastName    string
		ProviderUserName    string
		ProviderPhoneNumber string
		ProviderEmail       string

		ProjectID          string
		ProjectName        string
		ProjectDescription string
		ProjectLink        string

		SubscriptionPlanID          string
		SubscriptionPlanName        string
		SubscriptionPlanBenfits     string
		SubscriptionPlanDuration    int64
		SubscriptionPlanPrice       float64
		SubscriptionPlanIsRecurring bool
		SubscriptionPlanCurrency    string
	}

	subscriptionData := new(SubscriptionData)
	err := repo.conn.Raw(sqlStatement, subscriberID, subscriptionPlanID).Scan(subscriptionData).Error
	if err != nil {
		return nil, err
	}

	subscription := &entity.Subscription{
		SubscriberID:          subscriptionData.SubscriberID,
		SubscriberFirstName:   subscriptionData.SubscriberFirstName,
		SubscriberLastName:    subscriptionData.SubscriberLastName,
		SubscriberUserName:    subscriptionData.SubscriberUserName,
		SubscriberPhoneNumber: subscriptionData.SubscriberPhoneNumber,
		SubscriberEmail:       subscriptionData.SubscriberEmail,

		ProviderID:          subscriptionData.ProviderID,
		ProviderFirstName:   subscriptionData.ProviderFirstName,
		ProviderLastName:    subscriptionData.ProviderLastName,
		ProviderUserName:    subscriptionData.ProviderUserName,
		ProviderPhoneNumber: subscriptionData.ProviderPhoneNumber,
		ProviderEmail:       subscriptionData.ProviderEmail,

		ProjectID:          subscriptionData.ProjectID,
		ProjectName:        subscriptionData.ProjectName,
		ProjectDescription: subscriptionData.ProjectDescription,
		ProjectLink:        subscriptionData.ProjectLink,

		SubscriptionPlanID:          subscriptionData.SubscriptionPlanID,
		SubscriptionPlanName:        subscriptionData.SubscriptionPlanName,
		SubscriptionPlanBenfits:     subscriptionData.SubscriptionPlanBenfits,
		SubscriptionPlanDuration:    subscriptionData.SubscriptionPlanDuration,
		SubscriptionPlanPrice:       subscriptionData.SubscriptionPlanPrice,
		SubscriptionPlanIsRecurring: subscriptionData.SubscriptionPlanIsRecurring,
		SubscriptionPlanCurrency:    subscriptionData.SubscriptionPlanCurrency,
	}

	return subscription, nil
}

// Create is a method that adds a new subscription to the database
func (repo *SubscriptionRepository) Create(newSubscription *entity.Subscription) error {
	totalNumOfSubscriptions := tools.CountMembers("subscriptions", repo.conn)
	newSubscription.ID = fmt.Sprintf("SUB-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptions+1)

	for !tools.IsUnique("id", newSubscription.ID, "subscriptions", repo.conn) {
		totalNumOfSubscriptions++
		newSubscription.ID = fmt.Sprintf("SUB-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptions+1)
	}

	err := repo.conn.Create(newSubscription).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain subscription from the database using an subscription id,
// also Find() uses only id as a key for selection
func (repo *SubscriptionRepository) Find(id string) (*entity.Subscription, error) {

	subscription := new(entity.Subscription)
	err := repo.conn.Model(subscription).Where("id = ?", id).First(subscription).Error

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// FindMultiple is a method that finds multiple subscriptions from the database the matches the given identifier
// In FindMultiple() subscriber_id, project_id, subscription_plan_id or provider_id can used as a key
func (repo *SubscriptionRepository) FindMultiple(identifier string) []*entity.Subscription {

	var subscriptions []*entity.Subscription
	err := repo.conn.Model(entity.Subscription{}).Where(
		"subscriber_id = ?|| provider_id = ? || project_id = ? || subscription_plan_id = ?",
		identifier, identifier, identifier, identifier).Find(&subscriptions).Error

	if err != nil {
		return []*entity.Subscription{}
	}
	return subscriptions
}

// Update is a method that updates a certain subscription entries in the database
func (repo *SubscriptionRepository) Update(subscription *entity.Subscription) error {

	prevSubscription := new(entity.Subscription)
	err := repo.conn.Model(prevSubscription).Where("id = ?", subscription.ID).First(prevSubscription).Error

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

// Delete is a method that deletes a certain subscription from the database using an subscription id.
// In Delete() id is only used as an key
func (repo *SubscriptionRepository) Delete(id string) (*entity.Subscription, error) {
	subscription := new(entity.Subscription)
	err := repo.conn.Model(subscription).Where("id = ?", id).First(subscription).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(subscription)
	return subscription, nil
}

// DeleteMultiple is a method that deletes a set of subscriptions from the database using an identifier.
// In DeleteMultiple() subscriber_id, project_id, subscription_plan_id or provider_id can be used as an key
func (repo *SubscriptionRepository) DeleteMultiple(identifier string) []*entity.Subscription {
	var subscriptions []*entity.Subscription
	repo.conn.Model(subscriptions).Where(
		"subscriber_id = ?|| provider_id = ? || project_id = ? || subscription_plan_id = ?",
		identifier, identifier, identifier, identifier).
		Find(&subscriptions)

	for _, subscription := range subscriptions {
		repo.conn.Delete(subscription)
	}

	return subscriptions
}
