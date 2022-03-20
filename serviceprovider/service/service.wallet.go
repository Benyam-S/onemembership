package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/entity"
)

// AddSPWallet is a method that adds a new service provider wallet to the system
func (service *Service) AddSPWallet(newSPWallet *entity.SPWallet) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider wallet adding process, SP Wallet => %s",
		newSPWallet.ToString()), service.logger.Logs.ServiceProviderLogFile)

	err := service.spWalletRepo.Create(newSPWallet)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding SP Wallet => %s, %s",
			newSPWallet.ToString(), err.Error()))

		return errors.New("unable to add new service provider wallet")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider wallet adding process, SP Wallet => %s",
		newSPWallet.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return nil
}

// ValidateSPWallet is a method that validates a service provider wallet.
// It checks if the service provider wallet has a valid entries or not and return map of errors if any.
func (service *Service) ValidateSPWallet(spWallet *entity.SPWallet) entity.ErrMap {

	errMap := make(map[string]error)

	if len(spWallet.LinkedAccount) > 255 {
		errMap["linked_account"] = errors.New(`service provider account should not be longer than 255 characters`)
	}

	var isValidLinkedAccountProvider bool
	linkedAccountProviders := service.cmService.GetAllValidLinkedAccountProviders()
	for _, linkedAccountProvider := range linkedAccountProviders {
		if strings.ToUpper(linkedAccountProvider) == strings.ToUpper(spWallet.LinkedAccountProvider) {
			isValidLinkedAccountProvider = true
			break
		}
	}

	if !isValidLinkedAccountProvider {
		errMap["linked_account_provider"] = errors.New(`invalid account provider selected`)
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindSPWallet is a method that find and return a service provider wallet that matches the service provider id value
func (service *Service) FindSPWallet(providerID string) (*entity.SPWallet, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("SP Wallet finding process { Provider ID : %s }", providerID),
		service.logger.Logs.ServiceProviderLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, providerID)
	if empty {
		return nil, errors.New("no service provider wallet found")
	}

	spWallet, err := service.spWalletRepo.Find(providerID)
	if err != nil {
		return nil, errors.New("no service provider wallet found")
	}
	return spWallet, nil
}

// UpdateSPWallet is a method that updates a service provider wallet in the system
func (service *Service) UpdateSPWallet(spWallet *entity.SPWallet) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider wallet updating process, SP Wallet => %s",
		spWallet.ToString()), service.logger.Logs.ServiceProviderLogFile)

	err := service.spWalletRepo.Update(spWallet)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating SP Wallet => %s, %s",
			spWallet.ToString(), err.Error()))

		return errors.New("unable to update service provider wallet")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider wallet updating process, SP => %s",
		spWallet.ToString()), service.logger.Logs.ServiceProviderLogFile)
	return nil
}

// UpdateSPWalletSingleValue is a method that updates a single column entry of a service provider wallet
func (service *Service) UpdateSPWalletSingleValue(providerID, columnName string, columnValue interface{}) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started single service provider wallet value updating process "+
		"{ Provider ID : %s, Column Name : %s, Column Value : %s }", providerID, columnName, columnValue),
		service.logger.Logs.ServiceProviderLogFile)

	spWallet := entity.SPWallet{ProviderID: providerID}
	err := service.spWalletRepo.UpdateValue(&spWallet, columnName, columnValue)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For single service provider wallet value updating "+
			"{ Provider ID : %s, Column Name : %s, Column Value : %s }, %s", providerID,
			columnName, columnValue, err.Error()))

		return errors.New("unable to update service provider wallet")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished single service provider wallet value updating process, "+
		"SP Wallet => %s", spWallet.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return nil
}

// DeleteSPWallet is a method that deletes a service provider wallet from the system
func (service *Service) DeleteSPWallet(providerID string) (*entity.SPWallet, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started service provider wallet deleting process { Provider ID : %s }", providerID),
		service.logger.Logs.ServiceProviderLogFile)

	spWallet, err := service.spWalletRepo.Delete(providerID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting service provider wallet "+
			"{ Provider ID : %s }, %s", providerID, err.Error()))

		return nil, errors.New("unable to delete service provider")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished service provider wallet deleting process, Deleted SP Wallet => %s",
		spWallet.ToString()), service.logger.Logs.ServiceProviderLogFile)

	return spWallet, nil
}
