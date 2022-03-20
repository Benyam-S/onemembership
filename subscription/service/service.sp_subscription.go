package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
)

// AddSPSubscription is a method that adds a new service provider subscription to the system
func (service *Service) AddSPSubscription(newSubscription *entity.SPSubscription) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider subscription adding process, SP Subscription => %s",
		newSubscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	err := service.spSubscriptionRepo.Create(newSubscription)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding SP Subscription => %s, %s",
			newSubscription.ToString(), err.Error()))

		return errors.New("unable to add new subscription")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider subscription adding process, SP Subscription => %s",
		newSubscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	return nil
}

// FindSPSubscription is a method that find and return a service provider subscription that matches the providerID value
func (service *Service) FindSPSubscription(providerID string) (*entity.SPSubscription, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single service provider subscription finding process { Provider ID : %s }", providerID),
		service.logger.Logs.SubscriptionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, providerID)
	if empty {
		return nil, errors.New("no subscription found")
	}

	subscription, err := service.spSubscriptionRepo.Find(providerID)
	if err != nil {
		return nil, errors.New("no subscription found")
	}
	return subscription, nil
}

// FindMultipleSPSubscriptions is a method that find and return multiple service provider subscriptions that matchs the planID value
func (service *Service) FindMultipleSPSubscriptions(planID string) []*entity.SPSubscription {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple service provider subscriptions finding process { Plan ID : %s }", planID),
		service.logger.Logs.SubscriptionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, planID)
	if empty {
		return []*entity.SPSubscription{}
	}

	return service.spSubscriptionRepo.FindMultiple(planID)
}

// UpdateSPSubscription is a method that updates a service provider subscription in the system
func (service *Service) UpdateSPSubscription(subscription *entity.SPSubscription) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider subscription updating process, SP Subscription => %s",
		subscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	err := service.spSubscriptionRepo.Update(subscription)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating SP Subscription => %s, %s",
			subscription.ToString(), err.Error()))

		return errors.New("unable to update subscription")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider subscription updating process, SP Subscription => %s",
		subscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	return nil
}

// DeleteSPSubscription is a method that deletes a service provider subscription from the system using an providerID
func (service *Service) DeleteSPSubscription(providerID string) (*entity.SPSubscription, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider subscription deleting process { Provider ID : %s }",
		providerID), service.logger.Logs.SubscriptionLogFile)

	subscription, err := service.spSubscriptionRepo.Delete(providerID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting service provider subscription "+
			"{ Provider ID : %s }, %s", providerID, err.Error()))

		return nil, errors.New("unable to delete subscription")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider subscription deleting process, "+
		"Deleted SP Subscription => %s", subscription.ToString()), service.logger.Logs.SubscriptionLogFile)

	return subscription, nil
}

// DeleteMultipleSPSubscriptions is a method that deletes multiple service provider subscriptions from the system that match the given planID
func (service *Service) DeleteMultipleSPSubscriptions(planID string) []*entity.SPSubscription {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple service provider subscriptions deleting process { Plan ID : %s }",
		planID), service.logger.Logs.SubscriptionLogFile)

	return service.spSubscriptionRepo.DeleteMultiple(planID)
}
