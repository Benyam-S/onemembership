package subscription

import "github.com/Benyam-S/onemembership/entity"

// ISubscriptionRepository is an interface that defines all the repository methods of a subscription struct
type ISubscriptionRepository interface {
	Construct(subscriberID, subscriptionPlanID string) (*entity.Subscription, error)
	Create(newSubscription *entity.Subscription) error
	Find(id string) (*entity.Subscription, error)
	FindMultiple(identifier string) []*entity.Subscription
	Update(subscription *entity.Subscription) error
	Delete(id string) (*entity.Subscription, error)
	DeleteMultiple(identifier string) []*entity.Subscription
}

// ISPSubscriptionRepository is an interface that defines all the repository methods of a service provider subscription struct
type ISPSubscriptionRepository interface {
	Create(newSubscription *entity.SPSubscription) error
	Find(providerID string) (*entity.SPSubscription, error)
	FindMultiple(planID string) []*entity.SPSubscription
	Update(subscription *entity.SPSubscription) error
	Delete(providerID string) (*entity.SPSubscription, error)
	DeleteMultiple(planID string) []*entity.SPSubscription
}
