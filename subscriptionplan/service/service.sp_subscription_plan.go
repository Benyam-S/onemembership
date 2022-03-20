package service

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/entity"
)

// AddSPSubscriptionPlan is a method that adds a new service provider subscription plan to the system
func (service *Service) AddSPSubscriptionPlan(newSPSubscriptionPlan *entity.SPSubscriptionPlan) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider subscription plan adding process, SP Subscription Plan => %s",
		newSPSubscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	err := service.spSubscriptionPlanRepo.Create(newSPSubscriptionPlan)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding SP Subscription Plan => %s, %s",
			newSPSubscriptionPlan.ToString(), err.Error()))

		return errors.New("unable to add new subscription plan")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider subscription plan adding process, SP Subscription Plan => %s",
		newSPSubscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	return nil
}

// ValidateSPSubscriptionPlan is a method that validates a service provider subscription plan entries.
// It checks if the service provider subscription plan has a valid entries or not and return map of errors if any.
func (service *Service) ValidateSPSubscriptionPlan(subscriptionPlan *entity.SPSubscriptionPlan) entity.ErrMap {

	errMap := make(map[string]error)

	empty, _ := regexp.MatchString(`^\s*$`, subscriptionPlan.Name)
	if empty {
		errMap["name"] = errors.New(`subscription plan name can not be empty`)
	} else if len(subscriptionPlan.Name) > 255 {
		errMap["name"] = errors.New(`subscription plan name should not be longer than 255 characters`)
	}

	if subscriptionPlan.Duration < 1 || subscriptionPlan.Duration > 1000000 {
		errMap["duration"] = errors.New(`invalid subscription plan duration used`)
	}

	if subscriptionPlan.Price < 0 {
		errMap["price"] = errors.New(`invalid subscription plan price used`)
	} else {
		// Formatting the price
		subscriptionPlan.Price = math.Round(subscriptionPlan.Price*100) / 100
	}

	var isValidCurrencyType bool
	currencyTypes := service.cmService.GetAllValidCurrencyTypes()
	for _, currencyType := range currencyTypes {
		if strings.ToUpper(currencyType) == strings.ToUpper(subscriptionPlan.Currency) {
			isValidCurrencyType = true
			break
		}
	}

	if !isValidCurrencyType {
		errMap["currency"] = errors.New(`invalid currency type selected`)
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindSPSubscriptionPlan is a method that find and return a service provider subscription plan that matches the id value
func (service *Service) FindSPSubscriptionPlan(id string) (*entity.SPSubscriptionPlan, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single service provider subscription plan finding process "+
		"{ SP Subscription Plan ID : %s }", id), service.logger.Logs.SubscriptionPlanLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, id)
	if empty {
		return nil, errors.New("no subscription plan found")
	}

	subscriptionPlan, err := service.spSubscriptionPlanRepo.Find(id)
	if err != nil {
		return nil, errors.New("no subscription plan found")
	}
	return subscriptionPlan, nil
}

// UpdateSPSubscriptionPlan is a method that updates a service provider subscription plan in the system
func (service *Service) UpdateSPSubscriptionPlan(subscriptionPlan *entity.SPSubscriptionPlan) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider subscription plan updating process, SP Subscription Plan => %s",
		subscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	err := service.spSubscriptionPlanRepo.Update(subscriptionPlan)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating SP Subscription Plan => %s, %s",
			subscriptionPlan.ToString(), err.Error()))

		return errors.New("unable to update subscription plan")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider subscription plan updating process, SP Subscription Plan => %s",
		subscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	return nil
}

// DeleteSPSubscriptionPlan is a method that deletes a service provider subscription plan from the system using an id
func (service *Service) DeleteSPSubscriptionPlan(id string) (*entity.SPSubscriptionPlan, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider subscription plan deleting process "+
		"{ SP Subscription Plan ID : %s }", id), service.logger.Logs.SubscriptionPlanLogFile)

	subscriptionPlan, err := service.spSubscriptionPlanRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting subscription plan "+
			"{ SP Subscription Plan ID : %s }, %s", id, err.Error()))

		return nil, errors.New("unable to delete subscription plan")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider subscription plan deleting process, "+
		"Deleted SP Subscription Plan => %s", subscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)
	return subscriptionPlan, nil
}
