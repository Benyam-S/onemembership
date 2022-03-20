package project

import "github.com/Benyam-S/onemembership/entity"

// IProjectRepository is an interface that defines all the repository methods of a project struct
type IProjectRepository interface {
	Create(newProject *entity.Project) error
	Find(identifier string) (*entity.Project, error)
	FindMultiple(providerID string) []*entity.Project
	Update(project *entity.Project) error
	Delete(id string) (*entity.Project, error)
	DeleteMultiple(providerID string) []*entity.Project
}

// IProjectChatLinkRepository is an interface that defines all the repository methods of a project to chat link (ProjectChatLink) struct
type IProjectChatLinkRepository interface {
	Create(newProjectChatLink *entity.ProjectChatLink) error
	Find(projectID string, chatID int64) (*entity.ProjectChatLink, error)
	FindMultiple(identifier interface{}) []*entity.ProjectChatLink
	Delete(projectID string, chatID int64) (*entity.ProjectChatLink, error)
	DeleteMultiple(identifier interface{}) []*entity.ProjectChatLink
}
