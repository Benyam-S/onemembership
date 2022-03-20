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
	"github.com/Benyam-S/onemembership/user"
	"github.com/nyaruka/phonenumbers"
)

// Service is a type that defines a user service
type Service struct {
	userRepo          user.IUserRepository
	passwordRepo      user.IUserPasswordRepository
	preferenceService preference.IService
	feedbackService   feedback.IService
	deletedService    deleted.IService
	cmService         common.IService
	logger            *log.Logger
}

// NewUserService is a function that returns a new user service
func NewUserService(userRepository user.IUserRepository, passwordRepository user.IUserPasswordRepository,
	preferenceService preference.IService, feedbackService feedback.IService, deletedService deleted.IService,
	commonService common.IService, userLogger *log.Logger) user.IService {
	return &Service{userRepo: userRepository, passwordRepo: passwordRepository, preferenceService: preferenceService,
		feedbackService: feedbackService, deletedService: deletedService, cmService: commonService, logger: userLogger}
}

// AddUser is a method that adds a new user to the system
func (service *Service) AddUser(newUser *entity.User) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started user adding process, User => %s", newUser.ToString()),
		service.logger.Logs.UserLogFile)

	err := service.userRepo.Create(newUser)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding User => %s, %s", newUser.ToString(), err.Error()))

		return errors.New("unable to add new user")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished user adding process, User => %s", newUser.ToString()),
		service.logger.Logs.UserLogFile)

	return nil
}

// ValidateUserProfile is a method that validates a user profile.
// It checks if the user has a valid entries or not and return map of errors if any.
// Also it will add country code to the phone number value if not included: default country code +251
func (service *Service) ValidateUserProfile(user *entity.User) entity.ErrMap {

	errMap := make(map[string]error)

	// Removing all whitespaces
	phoneNumber := strings.Join(strings.Fields(user.PhoneNumber), "")

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
	isValidUsername, _ := regexp.MatchString(`^\w+$`, user.UserName)
	isValidPhoneNumber := phonenumbers.IsValidNumber(parsedPhoneNumber)
	isValidEmail, _ := regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|`+
		`(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`,
		user.Email)

	if len(user.FirstName) > 255 {
		errMap["first_name"] = errors.New(`first name should not be longer than 255 characters`)
	}

	if len(user.LastName) > 255 {
		errMap["last_name"] = errors.New(`last name should not be longer than 255 characters`)
	}

	// Casting username to lower case so case sensitive problem doesn't occur
	user.UserName = strings.ToLower(user.UserName)
	emptyUsername, _ := regexp.MatchString(`^\s*$`, user.UserName) // Username can be empty
	if !emptyUsername && !isValidUsername {
		errMap["user_name"] = errors.New(`username shouldn't contain space or any special characters`)
	} else if len(user.UserName) > 20 {
		errMap["user_name"] = errors.New(`username should not be longer than 20 characters`)
	}

	emptyEmail, _ := regexp.MatchString(`^\s*$`, user.Email) // Email can be empty
	if !emptyEmail && !isValidEmail {
		errMap["email"] = errors.New("invalid email address used")
	}

	emptyPhoneNumber, _ := regexp.MatchString(`^\s*$`, user.PhoneNumber) // Phonenumber can be empty
	if !emptyPhoneNumber && !isValidPhoneNumber {
		errMap["phone_number"] = errors.New("invalid phonenumber used")
	} else if isValidPhoneNumber {
		// If a valid phone number is provided, adjust the phone number to fit the database
		// Stored in +251900010197 format
		phoneNumber = fmt.Sprintf("+%d%d", parsedPhoneNumber.GetCountryCode(),
			parsedPhoneNumber.GetNationalNumber())

		user.PhoneNumber = phoneNumber
	}

	// Meaning a new user is being add
	if user.ID == "" {
		if errMap["user_name"] == nil && !emptyUsername && !service.cmService.IsUnique("user_name", user.UserName, "users") {
			errMap["user_name"] = errors.New(`username is taken, username should be unique`)
		}

		if errMap["email"] == nil && !emptyEmail && !service.cmService.IsUnique("email", user.Email, "users") {
			errMap["email"] = errors.New("email address already exists")
		}

		if isValidPhoneNumber && !service.cmService.IsUnique("phone_number", user.PhoneNumber, "users") {
			errMap["phone_number"] = errors.New("phone number already exists")
		}
	} else {
		// Meaning trying to update user
		prevUser, err := service.userRepo.Find(user.ID)

		// Checking for err isn't relevant but to make it robust check for nil pointer
		if err == nil && errMap["user_name"] == nil && prevUser.UserName != user.UserName && !emptyUsername {
			if !service.cmService.IsUnique("user_name", user.UserName, "users") {
				errMap["user_name"] = errors.New(`username is taken, username should be unique`)
			}
		}

		// checking uniqueness only for email that isn't identical to the provider's previous email
		if err == nil && errMap["email"] == nil && prevUser.Email != user.Email && !emptyEmail {
			if !service.cmService.IsUnique("email", user.Email, "users") {
				errMap["email"] = errors.New("email address already exists")
			}
		}

		// Checking for err isn't relevant but to make it robust check for nil pointer
		if err == nil && isValidPhoneNumber && prevUser.PhoneNumber != user.PhoneNumber {
			if !service.cmService.IsUnique("phone_number", user.PhoneNumber, "users") {
				errMap["phone_number"] = errors.New("phone number already exists")
			}
		}
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindUser is a method that find and return a user that matchs the identifier value
func (service *Service) FindUser(identifier string) (*entity.User, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("User finding process { Identifier : %s }", identifier),
		service.logger.Logs.UserLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return nil, errors.New("no user found")
	}

	user, err := service.userRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("no user found")
	}
	return user, nil
}

// AllUsers is a method that returns all the users in the system
func (service *Service) AllUsers() []*entity.User {
	return service.userRepo.All()
}

// AllUsersWithPagination is a method that returns all the users with pagination
func (service *Service) AllUsersWithPagination(pageNum int64) ([]*entity.User, int64) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Getting all users process { Page Number : %d }", pageNum),
		service.logger.Logs.UserLogFile)

	return service.userRepo.FindAll(pageNum)
}

// SearchUsers is a method that searchs and returns a set of users related to the key identifier
func (service *Service) SearchUsers(key string, pageNum int64, extra ...string) ([]*entity.User, int64) {

	/* ---------------------------- Logging ---------------------------- */
	extraLog := ""
	for index, extraValue := range extra {
		extraLog += fmt.Sprintf(", Extra%d : %s", index, extraValue)
	}
	service.logger.Log(fmt.Sprintf("Searching users process { Key : %s, Page Number : %d%s }", key, pageNum, extraLog),
		service.logger.Logs.UserLogFile)

	defaultSearchColumnsRegx := []string{"first_name"}
	defaultSearchColumnsRegx = append(defaultSearchColumnsRegx, extra...)
	defaultSearchColumns := []string{"id", "user_name", "phone_number", "email"}

	result2 := make([]*entity.User, 0)
	results := make([]*entity.User, 0)
	resultsMap := make(map[string]*entity.User)
	var pageCount1 int64 = 0
	var pageCount2 int64 = 0
	var pageCount int64 = 0

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return results, 0
	}

	result1, pageCount1 := service.userRepo.Search(key, pageNum, defaultSearchColumns...)
	if len(defaultSearchColumnsRegx) > 0 {
		result2, pageCount2 = service.userRepo.SearchWRegx(key, pageNum, defaultSearchColumnsRegx...)
	}

	for _, user := range result1 {
		resultsMap[user.ID] = user
	}

	for _, user := range result2 {
		resultsMap[user.ID] = user
	}

	for _, uniqueUser := range resultsMap {
		results = append(results, uniqueUser)
	}

	pageCount = pageCount1
	if pageCount < pageCount2 {
		pageCount = pageCount2
	}

	return results, pageCount
}

// TotalUsers is a method that returns the total number of users
func (service *Service) TotalUsers() int64 {
	return service.userRepo.Total()
}

// UpdateUser is a method that updates a user in the system
func (service *Service) UpdateUser(user *entity.User) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started user updating process, User => %s", user.ToString()),
		service.logger.Logs.UserLogFile)

	err := service.userRepo.Update(user)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating User => %s, %s", user.ToString(), err.Error()))

		return errors.New("unable to update user")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished user updating process, User => %s", user.ToString()),
		service.logger.Logs.UserLogFile)

	return nil
}

// UpdateUserSingleValue is a method that updates a single column entry of a user
func (service *Service) UpdateUserSingleValue(userID, columnName string, columnValue interface{}) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started single user value updating process "+
		"{ UserID : %s, ColumnName : %s, ColumnValue : %s }", userID, columnName, fmt.Sprint(columnValue)),
		service.logger.Logs.UserLogFile)

	user := entity.User{ID: userID}
	err := service.userRepo.UpdateValue(&user, columnName, columnValue)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating single user value "+
			"{ UserID : %s, ColumnName : %s, ColumnValue : %s }, %s", userID, columnName,
			fmt.Sprint(columnValue), err.Error()))

		return errors.New("unable to update user")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished single user value updating process, User => %s",
		user.ToString()), service.logger.Logs.UserLogFile)

	return nil
}

// DeleteUser is a method that deletes a user from the system
func (service *Service) DeleteUser(userID string) (*entity.User, error) {

	// Trashing user and user related data
	user, err := service.FindUser(userID)
	if err != nil {
		return nil, err
	}

	deletedUser, err := service.deletedService.AddUserToTrash(user)
	if err == nil {
		// Trashing the subscription transactions
		service.deletedService.AddSubscriptionTranstactionsToTrash(userID, deletedUser.UserID)
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started user deleting process { User ID : %s }", userID),
		service.logger.Logs.UserLogFile)

	user, err = service.userRepo.Delete(userID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting user { User ID : %s }, %s", userID, err.Error()))

		return nil, errors.New("unable to delete user")
	}

	// Deleting client preference
	service.preferenceService.DeleteClientPreference(userID)

	// Setting client id to null for feedbacks
	service.feedbackService.SetFeedbackClientIDNull(userID)

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished user deleting process, Deleted User => %s",
		user.ToString()), service.logger.Logs.UserLogFile)

	return user, nil
}
