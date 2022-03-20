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

// AddSPPassword is a method that adds new service provider password to the system
func (service *Service) AddSPPassword(newSPPassword *entity.SPPassword) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider password adding process, SP Password => %s",
		newSPPassword.ToString()), service.logger.Logs.ServiceProviderLogFile)

	err := service.spPasswordRepo.Create(newSPPassword)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding SP Password => %s, %s",
			newSPPassword.ToString(), err.Error()))

		return errors.New("unable to add new password")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider password adding process, SP Password => %s",
		newSPPassword.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return nil
}

// VerifySPPassword is a method that verify a service provider has provided a valid password with a matching verifyPassword entry
func (service *Service) VerifySPPassword(spPassword *entity.SPPassword, verifyPassword string) error {
	matchPassword, _ := regexp.MatchString(`^[a-zA-Z0-9\._\-&!?=#]{8}[a-zA-Z0-9\._\-&!?=#]*$`, spPassword.Password)

	if len(spPassword.Password) < 8 {
		return errors.New("password should contain at least 8 characters")
	}

	if !matchPassword {
		return errors.New("invalid characters used in password")
	}

	if spPassword.Password != verifyPassword {
		return errors.New("passwords do not match")
	}

	spPassword.Salt = tools.RandomStringGN(30)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(spPassword.Password+spPassword.Salt), 12)
	spPassword.Password = base64.StdEncoding.EncodeToString(hashedPassword)

	return nil
}

// FindSPPassword is a method that find and return a service provider's password that matchs the identifier value
func (service *Service) FindSPPassword(providerID string) (*entity.SPPassword, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("SP Password finding process { Provider ID : %s }", providerID),
		service.logger.Logs.ServiceProviderLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, providerID)
	if empty {
		return nil, errors.New("password not found")
	}

	spPassword, err := service.spPasswordRepo.Find(providerID)
	if err != nil {
		return nil, errors.New("password not found")
	}

	return spPassword, nil
}

// UpdateSPPassword is a method that updates a certain service provider's password
func (service *Service) UpdateSPPassword(spPassword *entity.SPPassword) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider password updating process, SP Password => %s",
		spPassword.ToString()), service.logger.Logs.ServiceProviderLogFile)

	err := service.spPasswordRepo.Update(spPassword)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating SP Password => %s, %s",
			spPassword.ToString(), err.Error()))

		return errors.New("unable to update password")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider password updating process, SP Password => %s",
		spPassword.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return nil
}

// DeleteSPPassword is a method that deletes a certain service provider's password
func (service *Service) DeleteSPPassword(providerID string) (*entity.SPPassword, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider password deleting process { Provider ID : %s }", providerID),
		service.logger.Logs.ServiceProviderLogFile)

	spPassword, err := service.spPasswordRepo.Delete(providerID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting service provider password { Provider ID : %s }, %s",
			providerID, err.Error()))

		return nil, errors.New("unable to delete password")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider password deleting process, Deleted SP Password => %s",
		spPassword.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return spPassword, nil
}
