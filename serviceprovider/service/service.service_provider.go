package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/deleted"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/feedback"
	"github.com/Benyam-S/onemembership/log"
	"github.com/Benyam-S/onemembership/preference"
	"github.com/Benyam-S/onemembership/serviceprovider"
	"github.com/nyaruka/phonenumbers"
)

// Service is a type that defines a service provider service
type Service struct {
	serviceProviderRepo serviceprovider.IServiceProviderRepository
	spPasswordRepo      serviceprovider.ISPPasswordRepository
	spWalletRepo        serviceprovider.ISPWalletRepository
	preferenceService   preference.IService
	feedbackService     feedback.IService
	deletedService      deleted.IService
	cmService           common.IService
	logger              *log.Logger
}

// NewServiceProviderService is a function that returns a new service provider service
func NewServiceProviderService(serviceProviderRepository serviceprovider.IServiceProviderRepository,
	spPasswordRepository serviceprovider.ISPPasswordRepository, spWalletRepository serviceprovider.ISPWalletRepository,
	preferenceService preference.IService, feedbackService feedback.IService, deletedService deleted.IService,
	commonService common.IService, serviceProviderLogger *log.Logger) serviceprovider.IService {
	return &Service{serviceProviderRepo: serviceProviderRepository, spPasswordRepo: spPasswordRepository,
		spWalletRepo: spWalletRepository, preferenceService: preferenceService,
		feedbackService: feedbackService, deletedService: deletedService, cmService: commonService,
		logger: serviceProviderLogger}
}

// AddServiceProvider is a method that adds a new service provider to the system
func (service *Service) AddServiceProvider(newServiceProvider *entity.ServiceProvider) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider adding process, Service Provider => %s",
		newServiceProvider.ToString()), service.logger.Logs.ServiceProviderLogFile)

	err := service.serviceProviderRepo.Create(newServiceProvider)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Service Provider => %s, %s",
			newServiceProvider.ToString(), err.Error()))

		return errors.New("unable to add new service provider")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider adding process, Service Provider => %s",
		newServiceProvider.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return nil
}

// ValidateProviderProfile is a method that validates a service provider profile.
// It checks if the service provider has a valid entries or not and return map of errors if any.
// Also it will add country code to the phone number value if not included: default country code +251
func (service *Service) ValidateProviderProfile(serviceProvider *entity.ServiceProvider) entity.ErrMap {

	errMap := make(map[string]error)

	// Removing all whitespaces
	phoneNumber := strings.Join(strings.Fields(serviceProvider.PhoneNumber), "")

	// Checking for local phone number
	isLocalPhoneNumber, _ := regexp.MatchString(`^0\d{9}$`, phoneNumber)

	if isLocalPhoneNumber {
		phoneNumberSlice := strings.Split(phoneNumber, "")
		if phoneNumberSlice[0] == "0" {
			phoneNumberSlice = phoneNumberSlice[1:]
			internationalPhoneNumber := "+251" + strings.Join(phoneNumberSlice, "")
			phoneNumber = internationalPhoneNumber
		}
	} else {
		// Making the phone number international if it is not local
		if len(phoneNumber) != 0 && string(phoneNumber[0]) != "+" {
			phoneNumber = "+" + phoneNumber
		}
	}

	parsedPhoneNumber, _ := phonenumbers.Parse(phoneNumber, "")
	isValidUsername, _ := regexp.MatchString(`^\w+$`, serviceProvider.UserName)
	isValidPhoneNumber := phonenumbers.IsValidNumber(parsedPhoneNumber)
	isValidEmail, _ := regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|`+
		`(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`,
		serviceProvider.Email)

	if len(serviceProvider.FirstName) > 255 {
		errMap["first_name"] = errors.New(`first name should not be longer than 255 characters`)
	}

	if len(serviceProvider.LastName) > 255 {
		errMap["last_name"] = errors.New(`last name should not be longer than 255 characters`)
	}

	serviceProvider.UserName = strings.ToLower(serviceProvider.UserName)
	emptyUsername, _ := regexp.MatchString(`^\s*$`, serviceProvider.UserName) // Username can be empty
	if !emptyUsername && !isValidUsername {
		errMap["user_name"] = errors.New(`username shouldn't contain space or any special characters`)
	} else if len(serviceProvider.UserName) > 20 {
		errMap["user_name"] = errors.New(`username should not be longer than 20 characters`)
	}

	emptyEmail, _ := regexp.MatchString(`^\s*$`, serviceProvider.Email) // Email can be empty
	if !emptyEmail && !isValidEmail {
		errMap["email"] = errors.New("invalid email address used")
	}

	emptyPhoneNumber, _ := regexp.MatchString(`^\s*$`, serviceProvider.PhoneNumber) // Phonenumber can be empty
	if !emptyPhoneNumber && !isValidPhoneNumber {
		errMap["phone_number"] = errors.New("invalid phonenumber used")
	} else if isValidPhoneNumber {
		// If a valid phone number is provided, adjust the phone number to fit the database
		// Stored in +251900010197 format
		phoneNumber = fmt.Sprintf("+%d%d", parsedPhoneNumber.GetCountryCode(),
			parsedPhoneNumber.GetNationalNumber())

		serviceProvider.PhoneNumber = phoneNumber
	}

	// Meaning a new serviceProvider is being add
	if serviceProvider.ID == "" {
		if errMap["user_name"] == nil && !emptyUsername && !service.cmService.IsUnique("user_name", serviceProvider.UserName, "service_providers") {
			errMap["user_name"] = errors.New(`username is taken, username should be unique`)
		}

		if errMap["email"] == nil && !emptyEmail && !service.cmService.IsUnique("email", serviceProvider.Email, "service_providers") {
			errMap["email"] = errors.New("email address already exists")
		}

		if isValidPhoneNumber && !service.cmService.IsUnique("phone_number", serviceProvider.PhoneNumber, "service_providers") {
			errMap["phone_number"] = errors.New("phone number already exists")
		}
	} else {
		// Meaning trying to update serviceProvider
		prevUser, err := service.serviceProviderRepo.Find(serviceProvider.ID)

		// Checking for err isn't relevant but to make it robust check for nil pointer
		if err == nil && errMap["user_name"] == nil && prevUser.UserName != serviceProvider.UserName && !emptyUsername {
			if !service.cmService.IsUnique("user_name", serviceProvider.UserName, "service_providers") {
				errMap["user_name"] = errors.New(`username is taken, username should be unique`)
			}
		}

		// checking uniqueness only for email that isn't identical to the provider's previous email
		if err == nil && errMap["email"] == nil && prevUser.Email != serviceProvider.Email && !emptyEmail {
			if !service.cmService.IsUnique("email", serviceProvider.Email, "service_providers") {
				errMap["email"] = errors.New("email address already exists")
			}
		}

		// Checking for err isn't relevant but to make it robust check for nil pointer
		if err == nil && isValidPhoneNumber && prevUser.PhoneNumber != serviceProvider.PhoneNumber {
			if !service.cmService.IsUnique("phone_number", serviceProvider.PhoneNumber, "service_providers") {
				errMap["phone_number"] = errors.New("phone number already exists")
			}
		}
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindServiceProvider is a method that find and return a service provider that matchs the identifier value
func (service *Service) FindServiceProvider(identifier string) (*entity.ServiceProvider, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Servie Provier finding process { Service Provider Identifier : %s }", identifier),
		service.logger.Logs.ServiceProviderLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("no service provider found")
	}

	serviceProvider, err := service.serviceProviderRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no service provider found")
	}
	return serviceProvider, nil
}

// AllServiceProviders is a method that returns all the service providers in the system
func (service *Service) AllServiceProviders() []*entity.ServiceProvider {
	return service.serviceProviderRepo.All()
}

// AllServiceProvidersWithPagination is a method that returns all the service providers with pagination
func (service *Service) AllServiceProvidersWithPagination(pageNum int64) ([]*entity.ServiceProvider, int64) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Getting all service providers process { Page Number : %d }", pageNum),
		service.logger.Logs.ServiceProviderLogFile)

	return service.serviceProviderRepo.FindAll(pageNum)
}

// SearchServiceProviders is a method that searchs and returns a set of service providers related to the key identifier
func (service *Service) SearchServiceProviders(key string, pageNum int64, extra ...string) ([]*entity.ServiceProvider, int64) {

	/* ---------------------------- Logging ---------------------------- */
	extraLog := ""
	for index, extraValue := range extra {
		extraLog += fmt.Sprintf(", Extra%d : %s", index, extraValue)
	}

	service.logger.Log(fmt.Sprintf("Searching service providers process { Key : %s, Page Number : %d%s }",
		key, pageNum, extraLog), service.logger.Logs.ServiceProviderLogFile)

	defaultSearchColumnsRegx := []string{"first_name"}
	defaultSearchColumnsRegx = append(defaultSearchColumnsRegx, extra...)
	defaultSearchColumns := []string{"id", "user_name", "phone_number", "email"}

	result2 := make([]*entity.ServiceProvider, 0)
	results := make([]*entity.ServiceProvider, 0)
	resultsMap := make(map[string]*entity.ServiceProvider)
	var pageCount1 int64 = 0
	var pageCount2 int64 = 0
	var pageCount int64 = 0

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return results, 0
	}

	result1, pageCount1 := service.serviceProviderRepo.Search(key, pageNum, defaultSearchColumns...)
	if len(defaultSearchColumnsRegx) > 0 {
		result2, pageCount2 = service.serviceProviderRepo.SearchWRegx(key, pageNum, defaultSearchColumnsRegx...)
	}

	for _, serviceProvider := range result1 {
		resultsMap[serviceProvider.ID] = serviceProvider
	}

	for _, serviceProvider := range result2 {
		resultsMap[serviceProvider.ID] = serviceProvider
	}

	for _, uniqueServiceProvider := range resultsMap {
		results = append(results, uniqueServiceProvider)
	}

	pageCount = pageCount1
	if pageCount < pageCount2 {
		pageCount = pageCount2
	}

	return results, pageCount
}

// TotalServiceProviders is a method that returns the total number of service providers
func (service *Service) TotalServiceProviders() int64 {
	return service.serviceProviderRepo.Total()
}

// UpdateServiceProvider is a method that updates a service provider in the system
func (service *Service) UpdateServiceProvider(serviceProvider *entity.ServiceProvider) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider updating process, Service Provider => %s",
		serviceProvider.ToString()), service.logger.Logs.ServiceProviderLogFile)

	err := service.serviceProviderRepo.Update(serviceProvider)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating Service Provider => %s, %s",
			serviceProvider.ToString(), err.Error()))

		return errors.New("unable to update service provider")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider updating process, Service Provider => %s",
		serviceProvider.ToString()), service.logger.Logs.ServiceProviderLogFile)
	return nil
}

// UpdateProviderSingleValue is a method that updates a single column entry of a service provider
func (service *Service) UpdateProviderSingleValue(providerID, columnName string, columnValue interface{}) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started single service provider value updating process "+
		"{ Service Provider ID : %s, Column Name : %s, Column Value : %s }", providerID, columnName, columnValue),
		service.logger.Logs.ServiceProviderLogFile)

	serviceProvider := entity.ServiceProvider{ID: providerID}
	err := service.serviceProviderRepo.UpdateValue(&serviceProvider, columnName, columnValue)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For single service provider value updating "+
			"{ Service Provider ID : %s, Column Name : %s, Column Value : %s }, %s", providerID, columnName,
			columnValue, err.Error()))

		return errors.New("unable to update service provider")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished single service provider value updating process, Service Provider => %s",
		serviceProvider.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return nil
}

// DeleteServiceProvider is a method that deletes a service provider from the system
func (service *Service) DeleteServiceProvider(providerID string) (*entity.ServiceProvider, error) {

	// Trashing service provider and service provider related data
	serviceProvider, err := service.FindServiceProvider(providerID)
	if err != nil {
		return nil, err
	}

	deletedServiceProvider, err := service.deletedService.AddServiceProviderToTrash(serviceProvider)
	if err == nil {
		// Trashing the service provider subscription transactions
		service.deletedService.AddSPSubscriptionTranstactionsToTrash(providerID, deletedServiceProvider.ProviderID)

		// Trashing the service provider payroll transactions
		service.deletedService.AddPayrollTranstactionsToTrash(providerID, deletedServiceProvider.ProviderID)
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider deleting process { Service Provider ID : %s }", providerID),
		service.logger.Logs.ServiceProviderLogFile)

	serviceProvider, err = service.serviceProviderRepo.Delete(providerID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting service provider { Service Provider ID : %s }, %s",
			providerID, err.Error()))

		return nil, errors.New("unable to delete service provider")
	}

	// Deleting client preference
	service.preferenceService.DeleteClientPreference(providerID)

	// Setting client id to null for feedbacks
	service.feedbackService.SetFeedbackClientIDNull(providerID)

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider deleting process, Deleted Service Provider => %s",
		serviceProvider.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return serviceProvider, nil
}
