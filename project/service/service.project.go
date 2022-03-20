package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
	"github.com/Benyam-S/onemembership/project"
)

// Service is a type that defines a project and ptspLink service
type Service struct {
	projectRepo         project.IProjectRepository
	projectChatLinkRepo project.IProjectChatLinkRepository
	cmService           common.IService
	logger              *log.Logger
}

// NewProjectService is a function that returns a new project and ptspLink service
func NewProjectService(projectRepository project.IProjectRepository,
	projectChatLinkRepository project.IProjectChatLinkRepository, commonService common.IService,
	projectLogger *log.Logger) project.IService {
	return &Service{projectRepo: projectRepository, projectChatLinkRepo: projectChatLinkRepository,
		cmService: commonService, logger: projectLogger}
}

// AddProject is a method that adds a new project to the system
func (service *Service) AddProject(newProject *entity.Project) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started project adding process, Project => %s",
		newProject.ToString()), service.logger.Logs.ProjectLogFile)

	err := service.projectRepo.Create(newProject)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Project => %s, %s",
			newProject.ToString(), err.Error()))

		return errors.New("unable to add new project")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished project adding process, Project => %s",
		newProject.ToString()), service.logger.Logs.ProjectLogFile)

	return nil
}

// ValidateProject is a method that validates a project entries.
// It checks if the project has a valid entries or not and return map of errors if any.
func (service *Service) ValidateProject(project *entity.Project) entity.ErrMap {

	errMap := make(map[string]error)

	emptyName, _ := regexp.MatchString(`^\s*$`, project.Name)
	if emptyName {
		errMap["name"] = errors.New(`project name can not be empty`)
	} else if len(project.Name) > 1000 { // Since it may contain special character it is better to use '1000'
		errMap["name"] = errors.New(`project name should not be longer than 1000 characters`)
	}

	if errMap["name"] == nil {
		projects := service.FindMultipleProjects(project.ProviderID)
		for _, prevProject := range projects {
			if prevProject.ID != project.ID && prevProject.Name == project.Name {
				errMap["name"] = errors.New(`project name already exists in your project list`)
				break
			}
		}
	}

	if len(project.Description) > 2000 {
		errMap["description"] = errors.New(`project description should not be longer than 2000 characters`)
	}

	// Changing the project link to lower case
	project.ProjectLink = strings.ToLower(project.ProjectLink)

	// Project Link shouldn't contain space
	isValidProjectLink, _ := regexp.MatchString(`^\w+$`, project.ProjectLink)
	if !isValidProjectLink {
		errMap["project_link"] = errors.New(`project link shouldn't contain space or any special characters`)
	} else if len(project.ProjectLink) > 20 {
		errMap["project_link"] = errors.New(`project link should not be longer than 20 characters`)
	}

	// Meaning a new project is being add
	if project.ID == "" {
		if errMap["project_link"] == nil && !service.cmService.IsUnique("project_link", project.ProjectLink, "projects") {
			errMap["project_link"] = errors.New(`project link is taken, link should be unique`)
		}
	} else {
		// Meaning trying to update user
		prevProject, err := service.projectRepo.Find(project.ID)

		// Checking for err isn't relevant but to make it robust check for nil pointer
		if err == nil && errMap["project_link"] == nil && prevProject.ProjectLink != project.ProjectLink {
			if !service.cmService.IsUnique("project_link", project.ProjectLink, "projects") {
				errMap["project_link"] = errors.New(`project link is taken, link should be unique`)
			}
		}
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindProject is a method that find and return a project that matches the identifier value
func (service *Service) FindProject(identifier string) (*entity.Project, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single project finding process { Project Identifier : %s }", identifier),
		service.logger.Logs.ProjectLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("no project found")
	}

	project, err := service.projectRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no project found")
	}
	return project, nil
}

// FindMultipleProjects is a method that find and return multiple projects that matchs the providerID value
func (service *Service) FindMultipleProjects(providerID string) []*entity.Project {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple project finding process { Provider ID : %s }", providerID),
		service.logger.Logs.ProjectLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, providerID)
	if empty {
		return []*entity.Project{}
	}

	return service.projectRepo.FindMultiple(providerID)
}

// UpdateProject is a method that updates a project in the system
func (service *Service) UpdateProject(project *entity.Project) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started project updating process, Project => %s",
		project.ToString()), service.logger.Logs.ProjectLogFile)

	err := service.projectRepo.Update(project)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating Project => %s, %s",
			project.ToString(), err.Error()))

		return errors.New("unable to update project")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished project updating process, Project => %s",
		project.ToString()), service.logger.Logs.ProjectLogFile)

	return nil
}

// DeleteProject is a method that deletes a project from the system using an id
func (service *Service) DeleteProject(id string) (*entity.Project, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started project deleting process { Project ID : %s }",
		id), service.logger.Logs.ProjectLogFile)

	project, err := service.projectRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting project { Project ID : %s }, %s", id, err.Error()))

		return nil, errors.New("unable to delete project")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished project deleting process, Deleted Project => %s",
		project.ToString()), service.logger.Logs.ProjectLogFile)
	return project, nil
}

// DeleteMultipleProjects is a method that deletes multiple projects from the system that match the given identifier
func (service *Service) DeleteMultipleProjects(providerID string) []*entity.Project {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple project deleting { Provider ID : %s }",
		providerID), service.logger.Logs.ProjectLogFile)

	return service.projectRepo.DeleteMultiple(providerID)
}
