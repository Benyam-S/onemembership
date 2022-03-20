package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/user"
	"github.com/jinzhu/gorm"
)

// UserPasswordRepository is a type that defines a user password repository
type UserPasswordRepository struct {
	conn *gorm.DB
}

// NewUserPasswordRepository is a function that returns a new user password repository
func NewUserPasswordRepository(connection *gorm.DB) user.IUserPasswordRepository {
	return &UserPasswordRepository{conn: connection}
}

// Create is a method that adds a new user password to the database
func (repo *UserPasswordRepository) Create(newUserPassword *entity.UserPassword) error {

	err := repo.conn.Create(newUserPassword).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain user's password from the database using an identifier.
// In Find() user_id is only used as a key
func (repo *UserPasswordRepository) Find(userID string) (*entity.UserPassword, error) {
	userPassword := new(entity.UserPassword)
	err := repo.conn.Model(userPassword).Where("user_id = ?", userID).
		First(userPassword).Error

	if err != nil {
		return nil, err
	}
	return userPassword, nil
}

// Update is a method that updates a certain user's password value in the database
func (repo *UserPasswordRepository) Update(userPassword *entity.UserPassword) error {

	prevUserPassword := new(entity.UserPassword)
	err := repo.conn.Model(prevUserPassword).Where("user_id = ?", userPassword.UserID).
		First(prevUserPassword).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	userPassword.CreatedAt = prevUserPassword.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(userPassword).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user's password from the database using an identifier.
// In Delete() user_id is only used as a key
func (repo *UserPasswordRepository) Delete(userID string) (*entity.UserPassword, error) {
	userPassword := new(entity.UserPassword)
	err := repo.conn.Model(userPassword).Where("user_id = ?", userID).First(userPassword).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(userPassword)
	return userPassword, nil
}
