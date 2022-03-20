package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

// Service is a type that defines a common service
type Service struct {
	commonRepo        common.ICommonRepository
	languageRepo      common.ILanguageRepository
	languageEntryRepo common.ILanguageEntryRepository
	logger            *log.Logger
}

// NewCommonService is a function that returns a new common service
func NewCommonService(commonRepository common.ICommonRepository, languageRepository common.ILanguageRepository,
	langagueEntryRepository common.ILanguageEntryRepository, languageLogger *log.Logger) common.IService {
	return &Service{commonRepo: commonRepository, languageRepo: languageRepository,
		languageEntryRepo: langagueEntryRepository, logger: languageLogger}
}

// IsUnique is a method that checks if a given column value is unique in a certain table
func (service *Service) IsUnique(columnName string, columnValue interface{}, tableName string) bool {
	return service.commonRepo.IsUnique(columnName, columnValue, tableName)
}

// FindLanguage is a method that find and return a certain language that matches the given identifier
func (service *Service) FindLanguage(identifier string) (*entity.Language, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Language finding process { Identifier : %s }", identifier),
		service.logger.Logs.ServerLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("no language found")
	}

	language, err := service.languageRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no language found")
	}

	return language, nil
}

// FindLanguageEntry is a method that find and return a language entry that matches the given identifier and code
func (service *Service) FindLanguageEntry(identifier, code string) string {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Language entry finding process { Identifier : %s, Code : %s }", identifier, code),
		service.logger.Logs.ServerLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return identifier
	}

	// Removing any special characters from identifier since identifier doesn't support special characters
	languageEntry, err := service.languageEntryRepo.FindWCode(emoji.RemoveAll(identifier), code)
	if err != nil {
		return identifier
	}

	return languageEntry.Value
}

// AllLanguages is a method that returns all the languages registered by the system
func (service *Service) AllLanguages() []*entity.Language {
	return service.languageRepo.All()
}

// GetAllValidChatTypes is a method that returns all the valid resource types that are supported by the system
func (service *Service) GetAllValidChatTypes() []string {
	return []string{"CHANNEL", "GROUP"}
}

// GetAllValidCurrencyTypes is a method that returns all the valid currency types that are supported by the system
func (service *Service) GetAllValidCurrencyTypes() []string {
	return []string{"ETB"}
}

// GetAllValidLinkedAccountProviders is a method that returns all the valid linked account providers that are supported by the system
func (service *Service) GetAllValidLinkedAccountProviders() []string {
	return []string{"CBE", "Telebirr"}
}
