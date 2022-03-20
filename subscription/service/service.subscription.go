package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
	"github.com/Benyam-S/onemembership/subscription"
)

// Service is a type that defines a subscription service
type Service struct {
	subscriptionRepo   subscription.ISubscriptionRepository
	spSubscriptionRepo subscription.ISPSubscriptionRepository
	logger             *log.Logger
}

// NewSubscriptionService is a function that returns a new subscription service
func NewSubscriptionService(subscriptionRepository subscription.ISubscriptionRepository,
	spSubscriptionRepository subscription.ISPSubscriptionRepository, subscriptionLogger *log.Logger) subscription.IService {
	return &Service{subscriptionRepo: subscriptionRepository, spSubscriptionRepo: spSubscriptionRepository,
		logger: subscriptionLogger}
}

// ConstructSubscription is a method that constructs a new subscription using subscribers id and plan id
func (service *Service) ConstructSubscription(subscriberID, subscriptionPlanID string) (*entity.Subscription, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription construction process { Subscriber ID : %s, Subscription Plan ID : %s }",
		subscriberID, subscriptionPlanID), service.logger.Logs.SubscriptionLogFile)

	newSubscription, err := service.subscriptionRepo.Construct(subscriberID, subscriptionPlanID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf(
			"Error: For constructing subscription { Subscriber ID : %s, Subscription Plan ID : %s }, %s",
			subscriberID, subscriptionPlanID, err.Error()))

		return nil, errors.New("unable to construct new subscription")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription construction process, Subscription => %s",
		newSubscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	return newSubscription, nil
}

// AddSubscription is a method that adds a new subscription to the system
func (service *Service) AddSubscription(newSubscription *entity.Subscription) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription adding process, Subscription => %s",
		newSubscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	err := service.subscriptionRepo.Create(newSubscription)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Subscription  => %s, %s",
			newSubscription.ToString(), err.Error()))

		return errors.New("unable to add new subscription")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription adding process, Subscription  => %s",
		newSubscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	return nil
}

// FindSubscription is a method that find and return a subscription that matches the id value
func (service *Service) FindSubscription(id string) (*entity.Subscription, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single subscription finding process { Subscription ID : %s }", id),
		service.logger.Logs.SubscriptionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, id)
	if empty {
		return nil, errors.New("no subscription found")
	}

	subscription, err := service.subscriptionRepo.Find(id)
	if err != nil {
		return nil, errors.New("no subscription found")
	}
	return subscription, nil
}

// FindMultipleSubscriptions is a method that find and return multiple subscriptions that matchs the identifier value
func (service *Service) FindMultipleSubscriptions(identifier string) []*entity.Subscription {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple subscriptions finding process { Subscription Identifier : %s }", identifier),
		service.logger.Logs.SubscriptionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return []*entity.Subscription{}
	}

	return service.subscriptionRepo.FindMultiple(identifier)
}

// UpdateSubscription is a method that updates a subscription in the system
func (service *Service) UpdateSubscription(subscription *entity.Subscription) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription updating process, Subscription => %s",
		subscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	err := service.subscriptionRepo.Update(subscription)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating Subscription => %s, %s",
			subscription.ToString(), err.Error()))

		return errors.New("unable to update subscription")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription updating process, Subscription => %s",
		subscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	return nil
}

// DeleteSubscription is a method that deletes a subscription from the system using an id
func (service *Service) DeleteSubscription(id string) (*entity.Subscription, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription deleting process { Subscription ID : %s }",
		id), service.logger.Logs.SubscriptionLogFile)

	subscription, err := service.subscriptionRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting subscription { Subscription ID : %s }, %s",
			id, err.Error()))

		return nil, errors.New("unable to delete subscription")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription deleting process, Deleted Subscription => %s",
		subscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	return subscription, nil
}

// DeleteMultipleSubscriptions is a method that deletes multiple subscriptions from the system that match the given identifier
func (service *Service) DeleteMultipleSubscriptions(identifier string) []*entity.Subscription {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple subscriptions deleting process { Subscription Identifier : %s }",
		identifier), service.logger.Logs.SubscriptionLogFile)

	return service.subscriptionRepo.DeleteMultiple(identifier)
}
