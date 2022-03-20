package project

import "github.com/Benyam-S/onemembership/entity"

// IService is an interface that defines all the service methods of a project struct
type IService interface {
	AddProject(newProject *entity.Project) error
	ValidateProject(project *entity.Project) entity.ErrMap
	FindProject(identifier string) (*entity.Project, error)
	FindMultipleProjects(providerID string) []*entity.Project
	UpdateProject(project *entity.Project) error
	DeleteProject(id string) (*entity.Project, error)
	DeleteMultipleProjects(providerID string) []*entity.Project

	AddProjectChatLink(newProjectChatLink *entity.ProjectChatLink) error
	ValidateProjectChatLink(projectChatLink *entity.ProjectChatLink) entity.ErrMap
	FindProjectChatLink(projectID string, chatID int64) (*entity.ProjectChatLink, error)
	FindMultipleProjectChatLinks(identifier interface{}) []*entity.ProjectChatLink
	DeleteProjectChatLink(projectID string, chatID int64) (*entity.ProjectChatLink, error)
	DeleteMultipleProjectChatLinks(identifier interface{}) []*entity.ProjectChatLink
}
