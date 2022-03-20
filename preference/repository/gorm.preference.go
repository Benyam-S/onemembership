package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/preference"
	"github.com/jinzhu/gorm"
)

// PreferenceRepository is a type that defines a client preference repository
type PreferenceRepository struct {
	conn *gorm.DB
}

// NewPreferenceRepository is a function that returns a new client preference repository
func NewPreferenceRepository(connection *gorm.DB) preference.IPreferenceRepository {
	return &PreferenceRepository{conn: connection}
}

// Create is a method that adds a new client preference to the database
func (repo *PreferenceRepository) Create(newClientPreference *entity.ClientPreference) error {
	err := repo.conn.Create(newClientPreference).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain client preference from the database using an clientID,
// also Find() uses only client_id as a key for selection
func (repo *PreferenceRepository) Find(clientID string) (*entity.ClientPreference, error) {
	clientPreference := new(entity.ClientPreference)
	err := repo.conn.Model(clientPreference).Where("client_id = ?", clientID).
		First(clientPreference).Error

	if err != nil {
		return nil, err
	}
	return clientPreference, nil
}

// Update is a method that updates a certain client preference value in the database
func (repo *PreferenceRepository) Update(clientPreference *entity.ClientPreference) error {

	prevClientPreference := new(entity.ClientPreference)
	err := repo.conn.Model(prevClientPreference).Where("client_id = ?", clientPreference.ClientID).
		First(prevClientPreference).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(clientPreference).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain client preference single column value in the database
func (repo *PreferenceRepository) UpdateValue(clientPreference *entity.ClientPreference, columnName string, columnValue interface{}) error {

	prevClientPreference := new(entity.ClientPreference)
	err := repo.conn.Model(prevClientPreference).Where("client_id = ?", clientPreference.ClientID).
		First(prevClientPreference).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.ClientPreference{}).Where("client_id = ?", clientPreference.ClientID).
		Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain client preference from the database using an clientID.
// In Delete() client_id is only used as a key
func (repo *PreferenceRepository) Delete(clientID string) (*entity.ClientPreference, error) {
	clientPreference := new(entity.ClientPreference)
	err := repo.conn.Model(clientPreference).Where("client_id = ?", clientID).
		First(clientPreference).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(clientPreference)
	return clientPreference, nil
}
