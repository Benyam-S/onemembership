package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/deleted"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/jinzhu/gorm"
)

// DeletedServiceProviderRepository is a type that defines a repository for deleted service provider
type DeletedServiceProviderRepository struct {
	conn *gorm.DB
}

// NewDeletedServiceProviderRepository is a function that returns a new deleted service provider repository
func NewDeletedServiceProviderRepository(connection *gorm.DB) deleted.IDeletedServiceProviderRepository {
	return &DeletedServiceProviderRepository{conn: connection}
}

// Create is a method that adds a deleted service provider to the database
func (repo *DeletedServiceProviderRepository) Create(deletedServiceProvider *entity.DeletedServiceProvider) error {

	totalNumOfDeletedSPs := tools.CountMembers("deleted_service_providers", repo.conn)
	deletedServiceProvider.ID = fmt.Sprintf("DSP_%s%d", tools.RandomStringGN(7), totalNumOfDeletedSPs+1)

	for !tools.IsUnique("id", deletedServiceProvider.ID, "deleted_service_providers", repo.conn) {
		totalNumOfDeletedSPs++
		deletedServiceProvider.ID = fmt.Sprintf("DSP_%s%d", tools.RandomStringGN(7), totalNumOfDeletedSPs+1)
	}

	// Adding prefix to the ProviderID so we can uniquely identify it from other deleted service providers
	deletedServiceProvider.ProviderID = deletedServiceProvider.ID + "_" + deletedServiceProvider.ProviderID

	err := repo.conn.Create(deletedServiceProvider).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain deleted service provider from the database using an id,
// also Find() uses id as a key for selection
func (repo *DeletedServiceProviderRepository) Find(id string) (*entity.DeletedServiceProvider, error) {
	deletedServiceProvider := new(entity.DeletedServiceProvider)
	err := repo.conn.Model(deletedServiceProvider).Where("id = ? ", id).
		First(deletedServiceProvider).Error

	if err != nil {
		return nil, err
	}

	return deletedServiceProvider, nil
}

// Search is a method that search and returns a set of deleted service providers from the database using an identifier.
func (repo *DeletedServiceProviderRepository) Search(key string, pageNum int64, columns ...string) []*entity.DeletedServiceProvider {

	var deletedServiceProviders []*entity.DeletedServiceProvider
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		// modifying the key so that it can match the database phone number values
		if column == "phone_number" {
			splitKey := strings.Split(key, "")
			if splitKey[0] == "0" {
				modifiedKey := "+251" + strings.Join(splitKey[1:], "")
				whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
				sqlValues = append(sqlValues, modifiedKey)
				continue
			}
		}
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_service_providers WHERE ("+strings.Join(whereStmt, "||")+") "+
		"ORDER BY first_name ASC LIMIT ?, 30", sqlValues...).Scan(&deletedServiceProviders)

	return deletedServiceProviders
}

// SearchWRegx is a method that searchs and returns set of deleted service providers limited to the key identifier and page number using regular expersions
func (repo *DeletedServiceProviderRepository) SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedServiceProvider {
	var deletedServiceProviders []*entity.DeletedServiceProvider
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_service_providers WHERE "+strings.Join(whereStmt, "||")+
		" ORDER BY first_name ASC LIMIT ?, 30", sqlValues...).Scan(&deletedServiceProviders)

	return deletedServiceProviders
}

// All is a method that returns all the deleted service providers from the database limited with the pageNum
func (repo *DeletedServiceProviderRepository) All(pageNum int64) []*entity.DeletedServiceProvider {

	var deletedServiceProviders []*entity.DeletedServiceProvider
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM deleted_service_providers ORDER BY first_name ASC LIMIT ?, 30",
		limit).Scan(&deletedServiceProviders)
	return deletedServiceProviders
}

// Update is a method that updates a certain deleted service provider value in the database
func (repo *DeletedServiceProviderRepository) Update(deletedServiceProvider *entity.DeletedServiceProvider) error {

	prevServiceProvider := new(entity.DeletedServiceProvider)
	err := repo.conn.Model(prevServiceProvider).Where("id = ?", deletedServiceProvider.ID).
		First(prevServiceProvider).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(deletedServiceProvider).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that removes a certain deleted service provider from the database using an id.
// In Delete() id is only used as a key
func (repo *DeletedServiceProviderRepository) Delete(id string) (*entity.DeletedServiceProvider, error) {
	deletedServiceProvider := new(entity.DeletedServiceProvider)
	err := repo.conn.Model(deletedServiceProvider).Where("id = ?", id).
		First(deletedServiceProvider).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(deletedServiceProvider)
	return deletedServiceProvider, nil
}
