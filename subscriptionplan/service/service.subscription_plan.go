package service

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
	"github.com/Benyam-S/onemembership/subscriptionplan"
)

// Service is a type that defines a subscription plan service
type Service struct {
	subscriptionPlanRepo   subscriptionplan.ISubscriptionPlanRepository
	spSubscriptionPlanRepo subscriptionplan.ISPSubscriptionPlanRepository
	planChatLinkRepo       subscriptionplan.IPlanChatLinkRepository
	userChatLinkRepo       subscriptionplan.IUserChatLinkRepository
	cmService              common.IService
	logger                 *log.Logger
}

// NewSubscriptionPlanService is a function that returns a new subscription plan service
func NewSubscriptionPlanService(subscriptionPlanRepository subscriptionplan.ISubscriptionPlanRepository,
	spSubscriptionPlanRepository subscriptionplan.ISPSubscriptionPlanRepository,
	planChatLinkRepository subscriptionplan.IPlanChatLinkRepository,
	userChatLinkRepository subscriptionplan.IUserChatLinkRepository, commonService common.IService,
	subscriptionPlanLogger *log.Logger) subscriptionplan.IService {
	return &Service{subscriptionPlanRepo: subscriptionPlanRepository, spSubscriptionPlanRepo: spSubscriptionPlanRepository,
		userChatLinkRepo: userChatLinkRepository, planChatLinkRepo: planChatLinkRepository,
		cmService: commonService, logger: subscriptionPlanLogger}
}

// AddSubscriptionPlan is a method that adds a new subscription plan to the system
func (service *Service) AddSubscriptionPlan(newSubscriptionPlan *entity.SubscriptionPlan) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription plan adding process, Subscription Plan => %s",
		newSubscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	err := service.subscriptionPlanRepo.Create(newSubscriptionPlan)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Subscription Plan => %s, %s",
			newSubscriptionPlan.ToString(), err.Error()))

		return errors.New("unable to add new subscription plan")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription plan adding process, Subscription Plan => %s",
		newSubscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	return nil
}

// ValidateSubscriptionPlan is a method that validates a subscription plan entries.
// It checks if the subscription plan has a valid entries or not and return map of errors if any.
func (service *Service) ValidateSubscriptionPlan(subscriptionPlan *entity.SubscriptionPlan) entity.ErrMap {

	errMap := make(map[string]error)

	empty, _ := regexp.MatchString(`^\s*$`, subscriptionPlan.Name)
	if empty {
		errMap["name"] = errors.New(`subscription plan name can not be empty`)
	} else if len(subscriptionPlan.Name) > 255 {
		errMap["name"] = errors.New(`subscription plan name should not be longer than 255 characters`)
	}

	if errMap["name"] == nil {
		subscriptionPlans := service.FindMultipleSubscriptionPlans(subscriptionPlan.ProjectID)
		for _, prevSubscriptionPlan := range subscriptionPlans {
			if subscriptionPlan.ID != prevSubscriptionPlan.ID && prevSubscriptionPlan.Name == subscriptionPlan.Name {
				errMap["name"] = errors.New(`subscription plan name already exists in the current project`)
				break
			}
		}
	}

	if len(subscriptionPlan.Benfits) > 2000 {
		errMap["benfits"] = errors.New(`subscription plan benfits should not be longer than 2000 characters`)
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

	var validCurrencyType bool
	currencyTypes := service.cmService.GetAllValidCurrencyTypes()
	for _, currencyType := range currencyTypes {
		if strings.ToUpper(currencyType) == strings.ToUpper(subscriptionPlan.Currency) {
			validCurrencyType = true
			break
		}
	}

	if !validCurrencyType {
		errMap["currency"] = errors.New(`invalid currency type selected`)
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindSubscriptionPlan is a method that find and return a subscription plan that matches the id value
func (service *Service) FindSubscriptionPlan(id string) (*entity.SubscriptionPlan, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single subscription plan finding process { Subscription Plan ID : %s }", id),
		service.logger.Logs.SubscriptionPlanLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, id)
	if empty {
		return nil, errors.New("no subscription plan found")
	}

	subscriptionPlan, err := service.subscriptionPlanRepo.Find(id)
	if err != nil {
		return nil, errors.New("no subscription plan found")
	}
	return subscriptionPlan, nil
}

// FindMultipleSubscriptionPlans is a method that find and return multiple subscription plans that matchs the projectID value
func (service *Service) FindMultipleSubscriptionPlans(projectID string) []*entity.SubscriptionPlan {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple subscription plan finding process { Project ID : %s }",
		projectID), service.logger.Logs.SubscriptionPlanLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, projectID)
	if empty {
		return []*entity.SubscriptionPlan{}
	}

	return service.subscriptionPlanRepo.FindMultiple(projectID)
}

// UpdateSubscriptionPlan is a method that updates a subscription plan in the system
func (service *Service) UpdateSubscriptionPlan(subscriptionPlan *entity.SubscriptionPlan) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription plan updating process, Subscription Plan => %s",
		subscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	err := service.subscriptionPlanRepo.Update(subscriptionPlan)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating Subscription Plan => %s, %s",
			subscriptionPlan.ToString(), err.Error()))

		return errors.New("unable to update subscription plan")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription plan updating process, Subscription Plan => %s",
		subscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	return nil
}

// DeleteSubscriptionPlan is a method that deletes a subscription plan from the system using an id
func (service *Service) DeleteSubscriptionPlan(id string) (*entity.SubscriptionPlan, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription plan deleting process { Subscription Plan ID : %s }",
		id), service.logger.Logs.SubscriptionPlanLogFile)

	subscriptionPlan, err := service.subscriptionPlanRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting subscription plan { Subscription Plan ID : %s }, %s",
			id, err.Error()))

		return nil, errors.New("unable to delete subscription plan")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription plan deleting process, Deleted Subscription Plan => %s",
		subscriptionPlan.ToString()), service.logger.Logs.SubscriptionPlanLogFile)
	return subscriptionPlan, nil
}

// DeleteMultipleSubscriptionPlans is a method that deletes multiple subscription plans from the system that match the given projectID
func (service *Service) DeleteMultipleSubscriptionPlans(projectID string) []*entity.SubscriptionPlan {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple subscription plan deleting { Project ID : %s }",
		projectID), service.logger.Logs.SubscriptionPlanLogFile)

	return service.subscriptionPlanRepo.DeleteMultiple(projectID)
}
