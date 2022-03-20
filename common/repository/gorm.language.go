package repository

import (
	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/jinzhu/gorm"
)

// LanguageRepository is a type that defines a language repository type
type LanguageRepository struct {
	conn *gorm.DB
}

// NewLanguageRepository is a function that creates a new language repository type
func NewLanguageRepository(connection *gorm.DB) common.ILanguageRepository {
	return &LanguageRepository{conn: connection}
}

// Create is a method that adds a new language to the database
func (repo *LanguageRepository) Create(newLanguage *entity.Language) error {
	err := repo.conn.Create(newLanguage).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain language from the database using an identifier,
// also Find() uses both code and name as a key for selection
func (repo *LanguageRepository) Find(identifier string) (*entity.Language, error) {

	language := new(entity.Language)
	err := repo.conn.Model(language).Where("code = ? || name = ?", identifier, identifier).
		First(language).Error

	if err != nil {
		return nil, err
	}
	return language, nil
}

// All is a method that returns all the languages found in the database
func (repo *LanguageRepository) All() []*entity.Language {

	var languages []*entity.Language

	repo.conn.Model(entity.Language{}).Order("display_order ASC").Find(&languages)

	return languages
}

// Update is a method that updates a certain language in the database
func (repo *LanguageRepository) Update(language *entity.Language) error {

	// If you want to change the code but keep the name or vise versa
	err := repo.conn.Exec("UPDATE languages SET code = ?, name = ? WHERE code = ? || name = ?",
		language.Code, language.Name, language.Code, language.Name).Error

	return err
}

// Delete is a method that deletes a certain language from the database using an identifier.
// In Delete() both code and name can be used as an key
func (repo *LanguageRepository) Delete(identifier string) (*entity.Language, error) {
	language := new(entity.Language)
	err := repo.conn.Model(language).Where("code = ? || name = ?", identifier, identifier).First(language).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(language)
	return language, nil
}
