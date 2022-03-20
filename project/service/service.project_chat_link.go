package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/entity"
)

// AddProjectChatLink is a method that adds a new projectChatLink to the system
func (service *Service) AddProjectChatLink(newProjectChatLink *entity.ProjectChatLink) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started project to chat link adding process, Project Chat Link => %s",
		newProjectChatLink.ToString()), service.logger.Logs.ProjectLogFile)

	err := service.projectChatLinkRepo.Create(newProjectChatLink)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Project Chat Link => %s, %s",
			newProjectChatLink.ToString(), err.Error()))

		return errors.New("unable to add new project to chat link")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished project to chat link adding process, Project Chat Link => %s",
		newProjectChatLink.ToString()), service.logger.Logs.ProjectLogFile)

	return nil
}

// ValidateProjectChatLink is a method that validates a projectChatLink entries.
// It checks if the projectChatLink has a valid entries or not and return map of errors if any.
func (service *Service) ValidateProjectChatLink(projectChatLink *entity.ProjectChatLink) entity.ErrMap {

	errMap := make(map[string]error)

	var validChatType bool
	chatTypes := service.cmService.GetAllValidChatTypes()
	for _, chatType := range chatTypes {
		if strings.ToUpper(chatType) == strings.ToUpper(projectChatLink.Type) {
			validChatType = true
			break
		}
	}

	if !validChatType {
		errMap["type"] = errors.New(`invalid chat type selected`)
	}

	// Chat ID == 0 means no username or chat id isn't provided
	if projectChatLink.ChatID == 0 {
		errMap["chat_id"] = errors.New(`invalid chat id used`)
	}

	// Checking if the chat has already been registered for the given project
	_, err := service.FindProjectChatLink(projectChatLink.ProjectID, projectChatLink.ChatID)
	if err == nil {
		errMap["chat_id"] = errors.New(`chat already linked to the project`)
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindProjectChatLink is a method that find and return a projectChatLink that matches the project id and chat id
func (service *Service) FindProjectChatLink(projectID string, chatID int64) (*entity.ProjectChatLink, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single project to chat link finding process { Project ID : %s, Chat ID : %d }",
		projectID, chatID), service.logger.Logs.ProjectLogFile)

	emptyProjectID, _ := regexp.MatchString(`^\s*$`, projectID)
	if emptyProjectID {
		return nil, errors.New("no project to chat link found")
	}

	projectChatLink, err := service.projectChatLinkRepo.Find(projectID, chatID)
	if err != nil {
		return nil, errors.New("no project to chat link found")
	}
	return projectChatLink, nil
}

// FindMultipleProjectChatLinks is a method that find and return multiple projectChatLinks that matchs the identifier value
func (service *Service) FindMultipleProjectChatLinks(identifier interface{}) []*entity.ProjectChatLink {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple project to chat link finding process { Project Chat Link Identifier : %s }",
		identifier), service.logger.Logs.ProjectLogFile)

	identifierS, ok := identifier.(string)
	empty, _ := regexp.MatchString(`^\s*$`, identifierS)
	if ok && empty {
		return []*entity.ProjectChatLink{}
	}

	return service.projectChatLinkRepo.FindMultiple(identifier)
}

// DeleteProjectChatLink is a method that deletes a projectChatLink from the system using an project id and chat id
func (service *Service) DeleteProjectChatLink(projectID string, chatID int64) (*entity.ProjectChatLink, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started project to chat link deleting process { Project ID : %s, Chat ID : %d }",
		projectID, chatID), service.logger.Logs.ProjectLogFile)

	projectChatLink, err := service.projectChatLinkRepo.Delete(projectID, chatID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf(
			"Error: For deleting project to chat link { Project ID : %s, Chat ID : %d }, %s",
			projectID, chatID, err.Error()))

		return nil, errors.New("unable to delete project to chat link")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished project to chat link deleting process, Deleted Project Chat Link => %s",
		projectChatLink.ToString()), service.logger.Logs.ProjectLogFile)
	return projectChatLink, nil
}

// DeleteMultipleProjectChatLinks is a method that deletes multiple projectChatLinks from the system that match the given identifier
func (service *Service) DeleteMultipleProjectChatLinks(identifier interface{}) []*entity.ProjectChatLink {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple project to chat link deleting { Project Chat Link Identifier : %s }",
		identifier), service.logger.Logs.ProjectLogFile)

	return service.projectChatLinkRepo.DeleteMultiple(identifier)
}
