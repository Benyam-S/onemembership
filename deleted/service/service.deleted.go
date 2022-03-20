package service

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Benyam-S/onemembership/deleted"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
)

// Service is a struct that defines a service that manages deleted struct
type Service struct {
	deletedUserRepo                  deleted.IDeletedUserRepository
	deletedServiceProvider           deleted.IDeletedServiceProviderRepository
	deletedSubscriptionTransaction   deleted.IDeletedSubscriptionTransactionRepository
	deletedSPSubscriptionTransaction deleted.IDeletedSPSubscriptionTransactionRepository
	deletedSPPayrollTransaction      deleted.IDeletedSPPayrollTransactionRepository
	logger                           *log.Logger
}

// NewDeletedService is a function that returns a new deleted service
func NewDeletedService(deletedUserRepository deleted.IDeletedUserRepository,
	deletedServiceProviderRepository deleted.IDeletedServiceProviderRepository,
	deletedSubscriptionTransactionRepository deleted.IDeletedSubscriptionTransactionRepository,
	deletedSPSubscriptionTransactionRepository deleted.IDeletedSPSubscriptionTransactionRepository,
	deletedSPPayrollTransactionRepository deleted.IDeletedSPPayrollTransactionRepository,
	deletedLogger *log.Logger) deleted.IService {

	return &Service{deletedUserRepo: deletedUserRepository, deletedServiceProvider: deletedServiceProviderRepository,
		deletedSubscriptionTransaction:   deletedSubscriptionTransactionRepository,
		deletedSPSubscriptionTransaction: deletedSPSubscriptionTransactionRepository,
		deletedSPPayrollTransaction:      deletedSPPayrollTransactionRepository, logger: deletedLogger}
}

// AddUserToTrash is a method that adds a user to deleted table
func (service *Service) AddUserToTrash(user *entity.User) (*entity.DeletedUser, error) {

	deletedUser := new(entity.DeletedUser)
	deletedUser.ID = user.ID
	deletedUser.FirstName = user.FirstName
	deletedUser.LastName = user.LastName
	deletedUser.Username = user.UserName
	deletedUser.PhoneNumber = user.PhoneNumber
	deletedUser.Email = user.Email

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started trashed user adding process, Deleted User => %s",
		deletedUser.ToString()), service.logger.Logs.DeletedLogFile)

	err := service.deletedUserRepo.Create(deletedUser)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Deleted User => %s, %s",
			deletedUser.ToString(), err.Error()))

		return nil, errors.New("unable to add user to trash")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished trashed user adding process, Deleted User => %s",
		deletedUser.ToString()), service.logger.Logs.DeletedLogFile)
	return deletedUser, nil
}

// AddServiceProviderToTrash is a method that adds a service provider to deleted table
func (service *Service) AddServiceProviderToTrash(serviceProvider *entity.ServiceProvider) (*entity.DeletedServiceProvider, error) {

	deletedServiceProvider := new(entity.DeletedServiceProvider)
	deletedServiceProvider.FirstName = serviceProvider.FirstName
	deletedServiceProvider.LastName = serviceProvider.LastName
	deletedServiceProvider.UserName = serviceProvider.UserName
	deletedServiceProvider.Email = serviceProvider.Email
	deletedServiceProvider.PhoneNumber = serviceProvider.PhoneNumber

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started trashed service provider adding process, Deleted Service Provider => %s",
		deletedServiceProvider.ToString()), service.logger.Logs.DeletedLogFile)

	err := service.deletedServiceProvider.Create(deletedServiceProvider)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Deleted Service Provider => %s, %s",
			deletedServiceProvider.ToString(), err.Error()))

		return nil, errors.New("unable to add service provider to trash")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished trashed service provider adding process, Deleted Service Provider => %s",
		deletedServiceProvider.ToString()), service.logger.Logs.DeletedLogFile)

	return deletedServiceProvider, nil
}

// AddSubscriptionTranstactionsToTrash is a method that adds subscritpion transactions to deleted table using userID
func (service *Service) AddSubscriptionTranstactionsToTrash(userID, prefixedUserID string) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription transactions trashing process, "+
		"{ User ID : %s, Prefixed User ID : %s }", userID, prefixedUserID), service.logger.Logs.DeletedLogFile)

	service.deletedSubscriptionTransaction.Create(userID, prefixedUserID)

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished trashing subscription transactions, { User ID : %s, Prefixed User ID : %s }",
		userID, prefixedUserID), service.logger.Logs.DeletedLogFile)
}

// AddSPSubscriptionTranstactionsToTrash is a method that adds service providers subscritpion transactions to deleted table using providerID
func (service *Service) AddSPSubscriptionTranstactionsToTrash(providerID, prefixedProviderID string) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider subscription transactions trashing process, "+
		"{ Provider ID : %s, Prefixed Provider ID : %s }", providerID, prefixedProviderID), service.logger.Logs.DeletedLogFile)

	service.deletedSPSubscriptionTransaction.Create(providerID, prefixedProviderID)

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished trashing service provider subscription transactions, "+
		"{ Provider ID : %s, Prefixed Provider ID : %s }", providerID, prefixedProviderID), service.logger.Logs.DeletedLogFile)
}

// AddPayrollTranstactionToTrash is a method that adds service providers payroll transactions to deleted table using providerID
func (service *Service) AddPayrollTranstactionsToTrash(providerID, prefixedProviderID string) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider payroll transactions trashing process, "+
		"{ Provider ID : %s, Prefixed Provider ID : %s }", providerID, prefixedProviderID), service.logger.Logs.DeletedLogFile)

	service.deletedSPPayrollTransaction.Create(providerID, prefixedProviderID)

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished trashing service provider payroll transactions, "+
		"{ Provider ID : %s, Prefixed Provider ID : %s }", providerID, prefixedProviderID), service.logger.Logs.DeletedLogFile)
}

// FindDeletedUser is a method that find and return a deleted user that matchs the deleted user id value
func (service *Service) FindDeletedUser(id string) (*entity.DeletedUser, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Deleted User finding process { Deleted User ID : %s }", id),
		service.logger.Logs.DeletedLogFile)

	deletedUser, err := service.deletedUserRepo.Find(id)
	if err != nil {
		return nil, errors.New("no deleted user found")
	}
	return deletedUser, nil
}

// SearchDeletedUsers is a method that searchs and returns a set of deleted users related to the key identifier
func (service *Service) SearchDeletedUsers(key, pagination string, extra ...string) []*entity.DeletedUser {

	/* ---------------------------- Logging ---------------------------- */
	extraLog := ""
	for index, extraValue := range extra {
		extraLog += fmt.Sprintf(", Extra%d : %s", index, extraValue)
	}
	service.logger.Log(fmt.Sprintf("Searching deleted users process { Key : %s, Pagination : %s%s }",
		key, pagination, extraLog), service.logger.Logs.DeletedLogFile)

	defaultSearchColumnsRegx := []string{"first_name"}
	defaultSearchColumnsRegx = append(defaultSearchColumnsRegx, extra...)
	defaultSearchColumns := []string{"id", "user_name", "phone_number", "email"}
	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	result2 := make([]*entity.DeletedUser, 0)
	results := make([]*entity.DeletedUser, 0)
	resultsMap := make(map[string]*entity.DeletedUser)

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return results
	}

	result1 := service.deletedUserRepo.Search(key, pageNum, defaultSearchColumns...)
	if len(defaultSearchColumnsRegx) > 0 {
		result2 = service.deletedUserRepo.SearchWRegx(key, pageNum, defaultSearchColumnsRegx...)
	}

	for _, deletedUser := range result1 {
		resultsMap[deletedUser.ID] = deletedUser
	}

	for _, deletedUser := range result2 {
		resultsMap[deletedUser.ID] = deletedUser
	}

	for _, uniqueDeletedUser := range resultsMap {
		results = append(results, uniqueDeletedUser)
	}

	return results
}

// FindDeletedServiceProvider is a method that find and return a deleted service provider that matchs the deleted service provider id value
func (service *Service) FindDeletedServiceProvider(id string) (*entity.DeletedServiceProvider, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Deleted Servie Provider finding process { Deleted Service Provider ID : %s }",
		id), service.logger.Logs.DeletedLogFile)

	deletedServiceProvider, err := service.deletedServiceProvider.Find(id)
	if err != nil {
		return nil, errors.New("no deleted service provider found")
	}
	return deletedServiceProvider, nil
}

// SearchDeletedServiceProviders is a method that searchs and returns a set of deleted service providers related to the key identifier
func (service *Service) SearchDeletedServiceProviders(key, pagination string, extra ...string) []*entity.DeletedServiceProvider {
	/* ---------------------------- Logging ---------------------------- */
	extraLog := ""
	for index, extraValue := range extra {
		extraLog += fmt.Sprintf(", Extra%d : %s", index, extraValue)
	}
	service.logger.Log(fmt.Sprintf("Searching deleted service provider process { Key : %s, Pagination : %s%s }",
		key, pagination, extraLog), service.logger.Logs.DeletedLogFile)

	defaultSearchColumnsRegx := []string{"first_name"}
	defaultSearchColumnsRegx = append(defaultSearchColumnsRegx, extra...)
	defaultSearchColumns := []string{"id", "user_name", "phone_number", "email"}
	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	result2 := make([]*entity.DeletedServiceProvider, 0)
	results := make([]*entity.DeletedServiceProvider, 0)
	resultsMap := make(map[string]*entity.DeletedServiceProvider)

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return results
	}

	result1 := service.deletedServiceProvider.Search(key, pageNum, defaultSearchColumns...)
	if len(defaultSearchColumnsRegx) > 0 {
		result2 = service.deletedServiceProvider.SearchWRegx(key, pageNum, defaultSearchColumnsRegx...)
	}

	for _, deletedServiceProvider := range result1 {
		resultsMap[deletedServiceProvider.ID] = deletedServiceProvider
	}

	for _, deletedServiceProvider := range result2 {
		resultsMap[deletedServiceProvider.ID] = deletedServiceProvider
	}

	for _, uniqueDeletedServiceProvider := range resultsMap {
		results = append(results, uniqueDeletedServiceProvider)
	}

	return results
}
