package repository

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/serviceprovider"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/jinzhu/gorm"
)

// ServiceProviderRepository is a type that defines a service provider repository type
type ServiceProviderRepository struct {
	conn *gorm.DB
}

// NewServiceProviderRepository is a function that creates a new service provider repository type
func NewServiceProviderRepository(connection *gorm.DB) serviceprovider.IServiceProviderRepository {
	return &ServiceProviderRepository{conn: connection}
}

// Create is a method that adds a new service provider to the database
func (repo *ServiceProviderRepository) Create(newServiceProvider *entity.ServiceProvider) error {
	totalNumOfMembers := tools.CountMembers("service_providers", repo.conn)
	newServiceProvider.ID = fmt.Sprintf("SP-%s%d", tools.RandomStringGN(7), totalNumOfMembers+1)

	for !tools.IsUnique("id", newServiceProvider.ID, "service_providers", repo.conn) {
		totalNumOfMembers++
		newServiceProvider.ID = fmt.Sprintf("SP-%s%d", tools.RandomStringGN(7), totalNumOfMembers+1)
	}

	err := repo.conn.Create(newServiceProvider).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain service provider from the database using an identifier,
// also Find() uses id, user_name, email and phone_number as a key for selection
func (repo *ServiceProviderRepository) Find(identifier string) (*entity.ServiceProvider, error) {

	modifiedIdentifier := identifier
	splitIdentifier := strings.Split(identifier, "")
	if splitIdentifier[0] == "0" {
		modifiedIdentifier = "+251" + strings.Join(splitIdentifier[1:], "")
	}

	serviceProvider := new(entity.ServiceProvider)
	err := repo.conn.Model(serviceProvider).
		Where("id = ? || user_name = ? || email = ?  || phone_number = ?",
			identifier, identifier, identifier, modifiedIdentifier).First(serviceProvider).Error

	if err != nil {
		return nil, err
	}
	return serviceProvider, nil
}

// FindAll is a method that returns set of service providers limited to the page number and category
func (repo *ServiceProviderRepository) FindAll(pageNum int64) ([]*entity.ServiceProvider, int64) {

	var serviceProviders []*entity.ServiceProvider
	var count float64

	repo.conn.Raw("SELECT * FROM service_providers ORDER BY first_name ASC LIMIT ?, 20", pageNum*20).Scan(&serviceProviders)
	repo.conn.Raw("SELECT COUNT(*) FROM service_providers").Count(&count)

	var pageCount int64 = int64(math.Ceil(count / 20.0))
	return serviceProviders, pageCount
}

// SearchWRegx is a method that searchs and returns set of service providers limited to the key identifier and page number using regular expersions
func (repo *ServiceProviderRepository) SearchWRegx(key string, pageNum int64, columns ...string) ([]*entity.ServiceProvider, int64) {
	var serviceProviders []*entity.ServiceProvider
	var whereStmt []string
	var sqlValues []interface{}
	var count float64

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	repo.conn.Raw("SELECT COUNT(*) FROM service_providers WHERE ("+strings.Join(whereStmt, "||")+") ",
		sqlValues...).Count(&count)

	sqlValues = append(sqlValues, pageNum*20)
	repo.conn.Raw("SELECT * FROM service_providers WHERE ("+strings.Join(whereStmt, "||")+") "+
		"ORDER BY first_name ASC LIMIT ?, 20", sqlValues...).Scan(&serviceProviders)

	var pageCount int64 = int64(math.Ceil(count / 20.0))
	return serviceProviders, pageCount
}

// Search is a method that searchs and returns set of service_providers limited to the key identifier and page number
func (repo *ServiceProviderRepository) Search(key string, pageNum int64, columns ...string) ([]*entity.ServiceProvider, int64) {
	var serviceProviders []*entity.ServiceProvider
	var whereStmt []string
	var sqlValues []interface{}
	var count float64

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

	repo.conn.Raw("SELECT COUNT(*) FROM service_providers WHERE ("+strings.Join(whereStmt, "||")+") ",
		sqlValues...).Count(&count)

	sqlValues = append(sqlValues, pageNum*20)
	repo.conn.Raw("SELECT * FROM service_providers WHERE ("+strings.Join(whereStmt, "||")+") "+
		"ORDER BY first_name ASC LIMIT ?, 20", sqlValues...).Scan(&serviceProviders)

	var pageCount int64 = int64(math.Ceil(count / 20.0))
	return serviceProviders, pageCount
}

// All is a method that returns all the service providers found in the database
func (repo *ServiceProviderRepository) All() []*entity.ServiceProvider {

	var serviceProviders []*entity.ServiceProvider

	repo.conn.Model(entity.ServiceProvider{}).Find(&serviceProviders).Order("created_at ASC")

	return serviceProviders
}

// Total is a method that retruns the total number of service providers
func (repo *ServiceProviderRepository) Total() int64 {

	var count int64
	repo.conn.Raw("SELECT COUNT(*) FROM service_providers").Count(&count)
	return count
}

// FromTo is a method that returns total number of service providers between start and end time
func (repo *ServiceProviderRepository) FromTo(start, end time.Time) int64 {

	var count int64
	repo.conn.Raw("SELECT COUNT(*) FROM service_providers WHERE created_at >= ? && created_at <= ?",
		start, end).Count(&count)
	return count
}

// Update is a method that updates a certain service provider entries in the database
func (repo *ServiceProviderRepository) Update(serviceProvider *entity.ServiceProvider) error {

	prevServiceProvider := new(entity.ServiceProvider)
	err := repo.conn.Model(prevServiceProvider).Where("id = ?", serviceProvider.ID).
		First(prevServiceProvider).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	serviceProvider.CreatedAt = prevServiceProvider.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(serviceProvider).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateValue is a method that updates a certain service provider single column value in the database
func (repo *ServiceProviderRepository) UpdateValue(serviceProvider *entity.ServiceProvider, columnName string, columnValue interface{}) error {

	prevServiceProvider := new(entity.ServiceProvider)
	err := repo.conn.Model(prevServiceProvider).Where("id = ?", serviceProvider.ID).First(prevServiceProvider).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.ServiceProvider{}).Where("id = ?", serviceProvider.ID).
		Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain service provider from the database using an service provider id.
// In Delete() id is only used as an key
func (repo *ServiceProviderRepository) Delete(id string) (*entity.ServiceProvider, error) {
	serviceProvider := new(entity.ServiceProvider)
	err := repo.conn.Model(serviceProvider).Where("id = ?", id).First(serviceProvider).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(serviceProvider)
	return serviceProvider, nil
}
