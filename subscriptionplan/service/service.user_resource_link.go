package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
)

// AddUserChatLink is a method that adds a new user to chat link to the system
func (service *Service) AddUserChatLink(newUserChatLink *entity.UserChatLink) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf(
		"Started user to chat link adding process, User Chat Link => %s", newUserChatLink.ToString()),
		service.logger.Logs.SubscriptionPlanLogFile)

	err := service.userChatLinkRepo.Create(newUserChatLink)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding User Chat Link => %s, %s",
			newUserChatLink.ToString(), err.Error()))

		return errors.New("unable to add new user to chat link")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished user to chat link adding process, User Chat Link => %s",
		newUserChatLink.ToString()), service.logger.Logs.SubscriptionPlanLogFile)

	return nil
}

// FindUserChatLink is a method that find and return a user to chat link that matches the user id, plan id and chat id
func (service *Service) FindUserChatLink(userID, planID string, chatID int64) (*entity.UserChatLink, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single user to chat link finding process { User ID : %s, Plan ID : %s, Chat ID : %d }",
		userID, planID, chatID), service.logger.Logs.SubscriptionPlanLogFile)

	emptyPlanID, _ := regexp.MatchString(`^\s*$`, planID)
	if emptyPlanID {
		return nil, errors.New("no user to chat link found")
	}

	spChatLink, err := service.userChatLinkRepo.Find(userID, planID, chatID)
	if err != nil {
		return nil, errors.New("no user to chat link found")
	}
	return spChatLink, nil
}

// FindMultipleUserChatLinks is a method that find and return multiple user to chat links that matchs the identifier value
func (service *Service) FindMultipleUserChatLinks(identifier interface{}) []*entity.UserChatLink {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf(
		"Multiple user to chat link finding process { User Chat Link Identifier : %s }",
		identifier), service.logger.Logs.SubscriptionPlanLogFile)

	identifierS, ok := identifier.(string)
	empty, _ := regexp.MatchString(`^\s*$`, identifierS)
	if ok && empty {
		return []*entity.UserChatLink{}
	}

	return service.userChatLinkRepo.FindMultiple(identifier)
}

// DeleteUserChatLink is a method that deletes a user to chat link from the system using user id, plan id and chat id
func (service *Service) DeleteUserChatLink(userID, planID string, chatID int64) (*entity.UserChatLink, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started user to chat link deleting process { User ID : %s, Plan ID : %s, Chat ID : %d }",
		userID, planID, chatID), service.logger.Logs.SubscriptionPlanLogFile)

	spChatLink, err := service.userChatLinkRepo.Delete(userID, planID, chatID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf(
			"Error: For deleting user to chat link { User ID : %s, Plan ID : %s, Chat ID : %d }, %s",
			userID, planID, chatID, err.Error()))

		return nil, errors.New("unable to delete user to chat link")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf(
		"Finished user to chat link deleting process, Deleted User Chat Link => %s",
		spChatLink.ToString()), service.logger.Logs.SubscriptionPlanLogFile)
	return spChatLink, nil
}

// DeleteMultipleUserChatLinks is a method that deletes multiple user to chat links from the system that match the given identifier
func (service *Service) DeleteMultipleUserChatLinks(identifier interface{}) []*entity.UserChatLink {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf(
		"Multiple user to chat link deleting { User Chat Link Identifier : %s }",
		identifier), service.logger.Logs.SubscriptionPlanLogFile)

	return service.userChatLinkRepo.DeleteMultiple(identifier)
}
