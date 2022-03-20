package common

import "github.com/Benyam-S/onemembership/entity"

// IService is an interface that defines all the common service methods
type IService interface {
	IsUnique(columnName string, columnValue interface{}, tableName string) bool
	FindLanguage(identifier string) (*entity.Language, error)
	FindLanguageEntry(identifier, code string) string
	AllLanguages() []*entity.Language
	GetAllValidChatTypes() []string
	GetAllValidCurrencyTypes() []string
	GetAllValidLinkedAccountProviders() []string
}
