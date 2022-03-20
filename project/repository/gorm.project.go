package repository

import (
	"fmt"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/project"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/jinzhu/gorm"
)

// ProjectRepository is a type that defines a project repository type
type ProjectRepository struct {
	conn *gorm.DB
}

// NewProjectRepository is a function that creates a new project repository type
func NewProjectRepository(connection *gorm.DB) project.IProjectRepository {
	return &ProjectRepository{conn: connection}
}

// Create is a method that adds a new project to the database
func (repo *ProjectRepository) Create(newProject *entity.Project) error {
	totalNumOfProjects := tools.CountMembers("projects", repo.conn)
	newProject.ID = fmt.Sprintf("Pr-%s%d", tools.RandomStringGN(7), totalNumOfProjects+1)

	for !tools.IsUnique("id", newProject.ID, "projects", repo.conn) {
		totalNumOfProjects++
		newProject.ID = fmt.Sprintf("Pr-%s%d", tools.RandomStringGN(7), totalNumOfProjects+1)
	}

	err := repo.conn.Create(newProject).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain project from the database using an identifier,
// also Find() uses id or project_link as a key for selection
func (repo *ProjectRepository) Find(identifier string) (*entity.Project, error) {

	project := new(entity.Project)
	err := repo.conn.Model(project).Where("id = ? || project_link = ?", identifier, identifier).
		First(project).Error

	if err != nil {
		return nil, err
	}
	return project, nil
}

// FindMultiple is a method that finds multiple projects from the database the matches the given providerID
// In FindMultiple() provider_id is used as a key
func (repo *ProjectRepository) FindMultiple(providerID string) []*entity.Project {

	var projects []*entity.Project
	err := repo.conn.Model(entity.Project{}).Where("provider_id = ?", providerID).
		Find(&projects).Error

	if err != nil {
		return []*entity.Project{}
	}
	return projects
}

// Update is a method that updates a certain project entries in the database
func (repo *ProjectRepository) Update(project *entity.Project) error {

	prevProject := new(entity.Project)
	err := repo.conn.Model(prevProject).Where("id = ?", project.ID).
		First(prevProject).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	project.CreatedAt = prevProject.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(project).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain project from the database using an project id.
// In Delete() id is only used as an key
func (repo *ProjectRepository) Delete(id string) (*entity.Project, error) {
	project := new(entity.Project)
	err := repo.conn.Model(project).Where("id = ?", id).First(project).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(project)
	return project, nil
}

// DeleteMultiple is a method that deletes a set of projects from the database using an providerID.
// In DeleteMultiple() provider_id is used as an key
func (repo *ProjectRepository) DeleteMultiple(providerID string) []*entity.Project {
	var projects []*entity.Project
	repo.conn.Model(projects).Where("provider_id = ?", providerID).
		Find(&projects)

	for _, project := range projects {
		repo.conn.Delete(project)
	}

	return projects
}
