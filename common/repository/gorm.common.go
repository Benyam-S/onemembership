package repository

import (
	"regexp"

	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/jinzhu/gorm"
)

// CommonRepository is a type that defines a repository for common use
type CommonRepository struct {
	conn *gorm.DB
}

// NewCommonRepository is a function that returns a new common repository type
func NewCommonRepository(connection *gorm.DB) common.ICommonRepository {
	return &CommonRepository{conn: connection}
}

// IsUnique is a methods that checks if a given column value is unique in a certain table
func (repo *CommonRepository) IsUnique(columnName string, columnValue interface{}, tableName string) bool {

	// Checking wether the value is empty or not, incase of string
	// If the value is empty then say it is unique
	if stringValue, ok := columnValue.(string); ok {
		emptyValue, _ := regexp.MatchString(`^\s*$`, stringValue)
		if emptyValue {
			return true
		}
	}

	return tools.IsUnique(columnName, columnValue, tableName, repo.conn)
}
