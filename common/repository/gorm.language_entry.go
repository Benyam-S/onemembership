package repository

import (
	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/jinzhu/gorm"
)

// LanguageEntryRepository is a type that defines a language entry repository type
type LanguageEntryRepository struct {
	conn *gorm.DB
}

// NewLanguageEntryRepository is a function that creates a new language entry repository type
func NewLanguageEntryRepository(connection *gorm.DB) common.ILanguageEntryRepository {
	return &LanguageEntryRepository{conn: connection}
}

// Create is a method that adds a new language entry to the database
func (repo *LanguageEntryRepository) Create(newLanguageEntry *entity.LanguageEntry) error {
	err := repo.conn.Create(newLanguageEntry).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain language entry from the database using its ID,
// also Find() uses only ID as a key for selection
func (repo *LanguageEntryRepository) Find(id int64) (*entity.LanguageEntry, error) {

	languageEntry := new(entity.LanguageEntry)
	err := repo.conn.Model(languageEntry).Where("id = ?", id).First(languageEntry).Error

	if err != nil {
		return nil, err
	}
	return languageEntry, nil
}

// FindWCode is a method that finds a certain language entry from the database using its language code,
// FindWCode() uses identifier and code as a key for selection
func (repo *LanguageEntryRepository) FindWCode(identifier, code string) (*entity.LanguageEntry, error) {

	languageEntry := new(entity.LanguageEntry)
	err := repo.conn.Model(languageEntry).Where("identifier = ? && code = ?", identifier, code).
		First(languageEntry).Error

	if err != nil {
		return nil, err
	}
	return languageEntry, nil
}

// Update is a method that updates a certain language entry in the database
// Update() uses only ID as a key for selection
func (repo *LanguageEntryRepository) Update(languageEntry *entity.LanguageEntry) error {

	prevLanguageEntry := new(entity.LanguageEntry)
	err := repo.conn.Model(prevLanguageEntry).Where("id = ?", languageEntry.ID).
		First(prevLanguageEntry).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(languageEntry).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateWCode is a method that updates a certain language entry in the database using its langauge code,
// UpdateWCode() uses identifier and code as a key for selection
func (repo *LanguageEntryRepository) UpdateWCode(languageEntry *entity.LanguageEntry) error {

	prevLanguageEntry := new(entity.LanguageEntry)
	err := repo.conn.Model(prevLanguageEntry).Where("identifier = ? && code = ?", languageEntry.Identifier,
		languageEntry.Code).First(prevLanguageEntry).Error

	if err != nil {
		return err
	}

	// Setting the primary key of the langague entry
	languageEntry.ID = prevLanguageEntry.ID

	err = repo.conn.Save(languageEntry).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain language entry from the database using an identifier.
// In Delete() ID is only used as an key
func (repo *LanguageEntryRepository) Delete(id int64) (*entity.LanguageEntry, error) {
	languageEntry := new(entity.LanguageEntry)
	err := repo.conn.Model(languageEntry).Where("id = ?", id).First(languageEntry).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(languageEntry)
	return languageEntry, nil
}

// DeleteWCode is a method that deletes a certain language entry from the database using its language code.
// DeleteWCode() uses identifier and code as a key for selection
func (repo *LanguageEntryRepository) DeleteWCode(identifier, code string) (*entity.LanguageEntry, error) {
	languageEntry := new(entity.LanguageEntry)
	err := repo.conn.Model(languageEntry).Where("identifier = ? && code = ?", identifier, code).
		First(languageEntry).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(languageEntry)
	return languageEntry, nil
}
