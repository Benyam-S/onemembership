package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/serviceprovider"
	"github.com/jinzhu/gorm"
)

// SPWalletRepository is a type that defines a service provider's wallet repository type
type SPWalletRepository struct {
	conn *gorm.DB
}

// NewSPWalletRepository is a function that creates a new service provider's wallet repository type
func NewSPWalletRepository(connection *gorm.DB) serviceprovider.ISPWalletRepository {
	return &SPWalletRepository{conn: connection}
}

// Create is a method that adds a new service provider wallet to the database
func (repo *SPWalletRepository) Create(newSPWallet *entity.SPWallet) error {

	err := repo.conn.Create(newSPWallet).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain service provider wallet from the database using an identifier,
// also Find() uses provider_id only as a key for selection
func (repo *SPWalletRepository) Find(providerID string) (*entity.SPWallet, error) {

	spWallet := new(entity.SPWallet)
	err := repo.conn.Model(spWallet).Where("provider_id = ?", providerID).
		First(spWallet).Error

	if err != nil {
		return nil, err
	}

	return spWallet, nil
}

// Update is a method that updates a certain service provider wallet entries in the database
func (repo *SPWalletRepository) Update(spWallet *entity.SPWallet) error {

	prevSPWallet := new(entity.SPWallet)
	err := repo.conn.Model(prevSPWallet).Where("provider_id = ?", spWallet.ProviderID).
		First(prevSPWallet).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	spWallet.CreatedAt = prevSPWallet.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(spWallet).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateValue is a method that updates a certain service provider wallet single column value in the database
func (repo *SPWalletRepository) UpdateValue(spWallet *entity.SPWallet, columnName string, columnValue interface{}) error {

	prevSPWallet := new(entity.SPWallet)
	err := repo.conn.Model(prevSPWallet).Where("provider_id = ?", spWallet.ProviderID).First(prevSPWallet).Error
	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.SPWallet{}).Where("provider_id = ?", spWallet.ProviderID).
		Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete is a method that deletes a certain service provider wallet from the database using an identifier.
// In Delete() provider_id is only used as an key
func (repo *SPWalletRepository) Delete(providerID string) (*entity.SPWallet, error) {
	spWallet := new(entity.SPWallet)
	err := repo.conn.Model(spWallet).Where("provider_id = ?", providerID).First(spWallet).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(spWallet)
	return spWallet, nil
}
