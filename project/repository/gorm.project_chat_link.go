package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/project"
	"github.com/jinzhu/gorm"
)

// ProjectChatLinkRepository is a type that defines a projectChatLink repository type
type ProjectChatLinkRepository struct {
	conn *gorm.DB
}

// NewProjectChatLinkRepository is a function that creates a new projectChatLink repository type
func NewProjectChatLinkRepository(connection *gorm.DB) project.IProjectChatLinkRepository {
	return &ProjectChatLinkRepository{conn: connection}
}

// Create is a method that adds a new projectChatLink to the database
func (repo *ProjectChatLinkRepository) Create(newProjectChatLink *entity.ProjectChatLink) error {
	err := repo.conn.Create(newProjectChatLink).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain projectChatLink from the database using projectID and chatID
func (repo *ProjectChatLinkRepository) Find(projectID string, chatID int64) (*entity.ProjectChatLink, error) {

	projectChatLink := new(entity.ProjectChatLink)
	err := repo.conn.Model(projectChatLink).Where("project_id = ? && chat_id = ?", projectID, chatID).
		First(projectChatLink).Error

	if err != nil {
		return nil, err
	}

	return projectChatLink, nil
}

// FindMultiple is a method that finds multiple projectChatLink from the database the matches the given identifier
// In FindMultiple() project_id and chat_id are used as a key
func (repo *ProjectChatLinkRepository) FindMultiple(identifier interface{}) []*entity.ProjectChatLink {

	var projectChatLinks []*entity.ProjectChatLink
	err := repo.conn.Model(entity.ProjectChatLink{}).Where("project_id = ? || chat_id = ?", identifier, identifier).
		Find(&projectChatLinks).Error

	if err != nil {
		return []*entity.ProjectChatLink{}
	}

	return projectChatLinks
}

// Delete is a method that deletes a certain projectChatLink from the database using projectID and chatID.
func (repo *ProjectChatLinkRepository) Delete(projectID string, chatID int64) (*entity.ProjectChatLink, error) {
	projectChatLink := new(entity.ProjectChatLink)
	err := repo.conn.Model(projectChatLink).Where("project_id = ? && chat_id = ?", projectID, chatID).
		First(projectChatLink).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Exec("DELETE FROM project_chat_links WHERE project_id = ? && chat_id = ?", projectID, chatID)

	return projectChatLink, nil
}

// DeleteMultiple is a method that deletes a set of projectChatLinks from the database using an identifier.
// In DeleteMultiple() project_id and chat_id are used as a key
func (repo *ProjectChatLinkRepository) DeleteMultiple(identifier interface{}) []*entity.ProjectChatLink {
	var projectChatLinks []*entity.ProjectChatLink
	repo.conn.Model(projectChatLinks).Where("project_id = ? || chat_id = ?", identifier, identifier).
		Find(&projectChatLinks)

	projectID, ok1 := identifier.(string)
	if ok1 {
		repo.conn.Exec("DELETE FROM project_chat_links WHERE project_id = ?", projectID)
		return projectChatLinks
	}

	chatID, ok2 := identifier.(int64)
	if ok2 {
		repo.conn.Exec("DELETE FROM project_chat_links WHERE chat_id = ?", chatID)
		return projectChatLinks
	}

	return projectChatLinks
}
