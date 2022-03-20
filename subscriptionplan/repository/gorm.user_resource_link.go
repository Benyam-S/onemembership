package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/subscriptionplan"
	"github.com/jinzhu/gorm"
)

// UserChatLinkRepository is a type that defines a user to chat link repository type
type UserChatLinkRepository struct {
	conn *gorm.DB
}

// NewUserChatLinkRepository is a function that creates a new user to chat link repository type
func NewUserChatLinkRepository(connection *gorm.DB) subscriptionplan.IUserChatLinkRepository {
	return &UserChatLinkRepository{conn: connection}
}

// Create is a method that adds a new user to chat link to the database
func (repo *UserChatLinkRepository) Create(newUserChatLink *entity.UserChatLink) error {
	err := repo.conn.Create(newUserChatLink).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain user to chat link from the database using userID, planID and chatID
func (repo *UserChatLinkRepository) Find(userID, planID string, chatID int64) (*entity.UserChatLink, error) {

	userChatLink := new(entity.UserChatLink)
	err := repo.conn.Model(userChatLink).Where("user_id = ? && plan_id = ? && chat_id = ?", userID, planID, chatID).
		First(userChatLink).Error

	if err != nil {
		return nil, err
	}

	return userChatLink, nil
}

// FindMultiple is a method that finds multiple user to chat link from the database the matches the given identifier
// In FindMultiple() user_id, plan_id and chat_id are used as a key
func (repo *UserChatLinkRepository) FindMultiple(identifier interface{}) []*entity.UserChatLink {

	var userChatLinks []*entity.UserChatLink
	err := repo.conn.Model(entity.UserChatLink{}).Where("user_id = ? || plan_id = ? || chat_id = ?",
		identifier, identifier, identifier).Find(&userChatLinks).Error

	if err != nil {
		return []*entity.UserChatLink{}
	}

	return userChatLinks
}

// Delete is a method that deletes a certain user to chat link from the database using userID, planID and chatID.
func (repo *UserChatLinkRepository) Delete(userID, planID string, chatID int64) (*entity.UserChatLink, error) {
	userChatLink := new(entity.UserChatLink)
	err := repo.conn.Model(userChatLink).Where("user_id = ? && plan_id = ? && chat_id = ?", userID, planID, chatID).
		First(userChatLink).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Exec("DELETE FROM user_chat_links WHERE user_id = ? && plan_id = ? && chat_id = ?",
		userID, planID, chatID)

	return userChatLink, nil
}

// DeleteMultiple is a method that deletes a set of user to chat link from the database using an identifier.
// In DeleteMultiple() user_id, plan_id and chat_id are used as a key
func (repo *UserChatLinkRepository) DeleteMultiple(identifier interface{}) []*entity.UserChatLink {
	var userChatLinks []*entity.UserChatLink
	repo.conn.Model(userChatLinks).Where("user_id = ? || plan_id = ? || chat_id = ?", identifier, identifier, identifier).
		Find(&userChatLinks)

	userID, ok1 := identifier.(string)
	if ok1 {
		repo.conn.Exec("DELETE FROM user_chat_links WHERE user_id = ?", userID)
		return userChatLinks
	}

	planID, ok2 := identifier.(string)
	if ok2 {
		repo.conn.Exec("DELETE FROM user_chat_links WHERE plan_id = ?", planID)
		return userChatLinks
	}

	chatID, ok3 := identifier.(int64)
	if ok3 {
		repo.conn.Exec("DELETE FROM user_chat_links WHERE chat_id = ?", chatID)
		return userChatLinks
	}

	return userChatLinks
}
