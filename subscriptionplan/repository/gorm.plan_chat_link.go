package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/subscriptionplan"
	"github.com/jinzhu/gorm"
)

// PlanChatLinkRepository is a type that defines a subscription plan to chat link repository type
type PlanChatLinkRepository struct {
	conn *gorm.DB
}

// NewPlanChatLinkRepository is a function that creates a new subscription plan to chat link repository type
func NewPlanChatLinkRepository(connection *gorm.DB) subscriptionplan.IPlanChatLinkRepository {
	return &PlanChatLinkRepository{conn: connection}
}

// Create is a method that adds a new subscription plan to chat link to the database
func (repo *PlanChatLinkRepository) Create(newPlanChatLink *entity.PlanChatLink) error {
	err := repo.conn.Create(newPlanChatLink).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain subscription plan to chat link from the database using planID and chatID
func (repo *PlanChatLinkRepository) Find(planID string, chatID int64) (*entity.PlanChatLink, error) {

	planChatLink := new(entity.PlanChatLink)
	err := repo.conn.Model(planChatLink).Where("plan_id = ? && chat_id = ?", planID, chatID).
		First(planChatLink).Error

	if err != nil {
		return nil, err
	}

	return planChatLink, nil
}

// FindMultiple is a method that finds multiple subscription plan to chat link from the database the matches the given identifier
// In FindMultiple() plan_id and chat_id are used as a key
func (repo *PlanChatLinkRepository) FindMultiple(identifier interface{}) []*entity.PlanChatLink {

	var planChatLinks []*entity.PlanChatLink
	err := repo.conn.Model(entity.PlanChatLink{}).Where("plan_id = ? || chat_id = ?", identifier, identifier).
		Find(&planChatLinks).Error

	if err != nil {
		return []*entity.PlanChatLink{}
	}

	return planChatLinks
}

// Delete is a method that deletes a certain subscription plan to chat link from the database using planID and chatID.
func (repo *PlanChatLinkRepository) Delete(planID string, chatID int64) (*entity.PlanChatLink, error) {
	planChatLink := new(entity.PlanChatLink)
	err := repo.conn.Model(planChatLink).Where("plan_id = ? && chat_id = ?", planID, chatID).
		First(planChatLink).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Exec("DELETE FROM plan_chat_links WHERE plan_id = ? && chat_id = ?", planID, chatID)

	return planChatLink, nil
}

// DeleteMultiple is a method that deletes a set of subscription plan to chat link from the database using an identifier.
// In DeleteMultiple() plan_id and chat_id are used as a key
func (repo *PlanChatLinkRepository) DeleteMultiple(identifier interface{}) []*entity.PlanChatLink {
	var planChatLinks []*entity.PlanChatLink
	repo.conn.Model(planChatLinks).Where("plan_id = ? || chat_id = ?", identifier, identifier).
		Find(&planChatLinks)

	planID, ok1 := identifier.(string)
	if ok1 {
		repo.conn.Exec("DELETE FROM plan_chat_links WHERE plan_id = ?", planID)
		return planChatLinks
	}

	chatID, ok2 := identifier.(int64)
	if ok2 {
		repo.conn.Exec("DELETE FROM plan_chat_links WHERE chat_id = ?", chatID)
		return planChatLinks
	}

	return planChatLinks
}
