package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
)

// AddPlanChatLink is a method that adds a new subscription plan to chat link to the system
func (service *Service) AddPlanChatLink(newPlanChatLink *entity.PlanChatLink) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf(
		"Started subscription plan to chat link adding process, Plan Chat Link => %s", newPlanChatLink.ToString()),
		service.logger.Logs.SubscriptionPlanLogFile)

	err := service.planChatLinkRepo.Create(newPlanChatLink)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Plan Chat Link => %s, %s",
			newPlanChatLink.ToString(), err.Error()))

		return errors.New("unable to add new subscription plan to chat link")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription plan to chat link adding process, Plan Chat Link => %s",
		newPlanChatLink.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	return nil
}

// FindPlanChatLink is a method that find and return a subscription plan to chat link that matches the plan id and chat id
func (service *Service) FindPlanChatLink(planID string, chatID int64) (*entity.PlanChatLink, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single subscription plan to chat link finding process { Plan ID : %s, Chat ID : %d }",
		planID, chatID), service.logger.Logs.SubscriptionPlanLogFile)

	emptyPlanID, _ := regexp.MatchString(`^\s*$`, planID)
	if emptyPlanID {
		return nil, errors.New("no subscription plan to chat link found")
	}

	planChatLink, err := service.planChatLinkRepo.Find(planID, chatID)
	if err != nil {
		return nil, errors.New("no subscription plan to chat link found")
	}
	return planChatLink, nil
}

// FindMultiplePlanChatLinks is a method that find and return multiple subscription plan to chat links that matchs the identifier value
func (service *Service) FindMultiplePlanChatLinks(identifier interface{}) []*entity.PlanChatLink {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf(
		"Multiple subscription plan to chat link finding process { Plan Chat Link Identifier : %s }",
		identifier), service.logger.Logs.SubscriptionPlanLogFile)

	identifierS, ok := identifier.(string)
	empty, _ := regexp.MatchString(`^\s*$`, identifierS)
	if ok && empty {
		return []*entity.PlanChatLink{}
	}

	return service.planChatLinkRepo.FindMultiple(identifier)
}

// DeletePlanChatLink is a method that deletes a subscription plan to chat link from the system using plan id and chat id
func (service *Service) DeletePlanChatLink(planID string, chatID int64) (*entity.PlanChatLink, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription plan to chat link deleting process { Plan ID : %s, Chat ID : %d }",
		planID, chatID), service.logger.Logs.SubscriptionPlanLogFile)

	planChatLink, err := service.planChatLinkRepo.Delete(planID, chatID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf(
			"Error: For deleting subscription plan to chat link { Plan ID : %s, Chat ID : %d }, %s",
			planID, chatID, err.Error()))

		return nil, errors.New("unable to delete subscription plan to chat link")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf(
		"Finished subscription plan to chat link deleting process, Deleted Plan Chat Link => %s",
		planChatLink.ToString()), service.logger.Logs.SubscriptionPlanLogFile)
	return planChatLink, nil
}

// DeleteMultiplePlanChatLinks is a method that deletes multiple subscription plan to chat links from the system that match the given identifier
func (service *Service) DeleteMultiplePlanChatLinks(identifier interface{}) []*entity.PlanChatLink {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf(
		"Multiple subscription plan to chat link deleting { Plan Chat Link Identifier : %s }",
		identifier), service.logger.Logs.SubscriptionPlanLogFile)

	return service.planChatLinkRepo.DeleteMultiple(identifier)
}
