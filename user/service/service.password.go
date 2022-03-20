package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/tools"
	"golang.org/x/crypto/bcrypt"
)

// AddUserPassword is a method that adds new user password to the system
func (service *Service) AddUserPassword(newUserPassword *entity.UserPassword) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started user password adding process, User Password => %s",
		newUserPassword.ToString()), service.logger.Logs.UserLogFile)

	err := service.passwordRepo.Create(newUserPassword)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding User Password => %s, %s",
			newUserPassword.ToString(), err.Error()))

		return errors.New("unable to add new password")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished user password adding process, User Password => %s",
		newUserPassword.ToString()), service.logger.Logs.UserLogFile)

	return nil
}

// VerifyUserPassword is a method that verify a user has provided a valid password with a matching verifyPassword entry
func (service *Service) VerifyUserPassword(userPassword *entity.UserPassword, verifyPassword string) error {
	matchPassword, _ := regexp.MatchString(`^[a-zA-Z0-9\._\-&!?=#]{8}[a-zA-Z0-9\._\-&!?=#]*$`, userPassword.Password)

	if len(userPassword.Password) < 8 {
		return errors.New("password should contain at least 8 characters")
	}

	if !matchPassword {
		return errors.New("invalid characters used in password")
	}

	if userPassword.Password != verifyPassword {
		return errors.New("passwords do not match")
	}

	userPassword.Salt = tools.RandomStringGN(30)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword.Password+userPassword.Salt), 12)
	userPassword.Password = base64.StdEncoding.EncodeToString(hashedPassword)

	return nil
}

// FindUserPassword is a method that find and return a user's password that matchs the identifier value
func (service *Service) FindUserPassword(userID string) (*entity.UserPassword, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("User Password finding process { User ID : %s }", userID),
		service.logger.Logs.UserLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, userID)
	if empty {
		return nil, errors.New("password not found")
	}

	userPassword, err := service.passwordRepo.Find(userID)
	if err != nil {
		return nil, errors.New("password not found")
	}

	return userPassword, nil
}

// UpdateUserPassword is a method that updates a certain user's password
func (service *Service) UpdateUserPassword(userPassword *entity.UserPassword) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started user password updating process, User Password => %s",
		userPassword.ToString()), service.logger.Logs.UserLogFile)

	err := service.passwordRepo.Update(userPassword)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating User Password => %s, %s",
			userPassword.ToString(), err.Error()))

		return errors.New("unable to update password")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished user password updating process, User Password => %s",
		userPassword.ToString()), service.logger.Logs.UserLogFile)

	return nil
}

// DeleteUserPassword is a method that deletes a certain user's password
func (service *Service) DeleteUserPassword(userID string) (*entity.UserPassword, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started user password deleting process { User ID : %s }", userID),
		service.logger.Logs.UserLogFile)

	userPassword, err := service.passwordRepo.Delete(userID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting user password { User ID : %s }, %s",
			userID, err.Error()))

		return nil, errors.New("unable to delete password")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished user password deleting process, Deleted User Password => %s",
		userPassword.ToString()), service.logger.Logs.UserLogFile)

	return userPassword, nil
}
