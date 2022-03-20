package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/serviceprovider"
	"github.com/jinzhu/gorm"
)

// SPPasswordRepository is a type that defines a service provider's password repository
type SPPasswordRepository struct {
	conn *gorm.DB
}

// NewSPPasswordRepository is a function that returns a new service provider's password repository
func NewSPPasswordRepository(connection *gorm.DB) serviceprovider.ISPPasswordRepository {
	return &SPPasswordRepository{conn: connection}
}

// Create is a method that adds a new service provider password to the database
func (repo *SPPasswordRepository) Create(newSPPassword *entity.SPPassword) error {

	err := repo.conn.Create(newSPPassword).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain service provider's password from the database using an identifier.
// In Find() provider_id is only used as a key
func (repo *SPPasswordRepository) Find(providerID string) (*entity.SPPassword, error) {
	spPassword := new(entity.SPPassword)
	err := repo.conn.Model(spPassword).Where("provider_id = ?", providerID).
		First(spPassword).Error

	if err != nil {
		return nil, err
	}

	return spPassword, nil
}

// Update is a method that updates a certain service provider's password value in the database
func (repo *SPPasswordRepository) Update(spPassword *entity.SPPassword) error {

	prevSPPassword := new(entity.SPPassword)
	err := repo.conn.Model(prevSPPassword).Where("provider_id = ?", spPassword.ProviderID).
		First(prevSPPassword).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	spPassword.CreatedAt = prevSPPassword.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(spPassword).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete is a method that deletes a certain service provider's password from the database using an identifier.
// In Delete() provider_id is only used as a key
func (repo *SPPasswordRepository) Delete(providerID string) (*entity.SPPassword, error) {
	spPassword := new(entity.SPPassword)
	err := repo.conn.Model(spPassword).Where("provider_id = ?", providerID).First(spPassword).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(spPassword)
	return spPassword, nil
}
