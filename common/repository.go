package common

import "github.com/Benyam-S/onemembership/entity"

// ICommonRepository is an interface that defines all the common repository methods
type ICommonRepository interface {
	IsUnique(columnName string, columnValue interface{}, tableName string) bool
}

// ILanguageRepository is an interface that defines all the repository methods of a language struct
type ILanguageRepository interface {
	Create(newLanguageEntry *entity.Language) error
	Find(identifier string) (*entity.Language, error)
	All() []*entity.Language
	Update(languageEntry *entity.Language) error
	Delete(identifier string) (*entity.Language, error)
}

// ILanguageEntryRepository is an interface that defines all the repository methods of a language entry struct
type ILanguageEntryRepository interface {
	Create(newLanguageEntry *entity.LanguageEntry) error
	Find(id int64) (*entity.LanguageEntry, error)
	FindWCode(identifier, code string) (*entity.LanguageEntry, error)
	Update(languageEntry *entity.LanguageEntry) error
	UpdateWCode(languageEntry *entity.LanguageEntry) error
	Delete(id int64) (*entity.LanguageEntry, error)
	DeleteWCode(identifier, code string) (*entity.LanguageEntry, error)
}
