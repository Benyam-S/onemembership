package subscription

import "github.com/Benyam-S/onemembership/entity"

// IService is an interface that defines all the service methods of a subscription struct
type IService interface {
	ConstructSubscription(subscriberID, subscriptionPlanID string) (*entity.Subscription, error)
	AddSubscription(newSubscription *entity.Subscription) error
	FindSubscription(id string) (*entity.Subscription, error)
	FindMultipleSubscriptions(identifier string) []*entity.Subscription
	UpdateSubscription(subscription *entity.Subscription) error
	DeleteSubscription(id string) (*entity.Subscription, error)
	DeleteMultipleSubscriptions(identifier string) []*entity.Subscription

	AddSPSubscription(newSubscription *entity.SPSubscription) error
	FindSPSubscription(providerID string) (*entity.SPSubscription, error)
	FindMultipleSPSubscriptions(planID string) []*entity.SPSubscription
	UpdateSPSubscription(subscription *entity.SPSubscription) error
	DeleteSPSubscription(providerID string) (*entity.SPSubscription, error)
	DeleteMultipleSPSubscriptions(planID string) []*entity.SPSubscription
}
