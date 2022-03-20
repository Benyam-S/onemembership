package service

import (
	"errors"
	"fmt"

	"regexp"

	"github.com/Benyam-S/onemembership/common"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
	"github.com/Benyam-S/onemembership/preference"
)

// Service is a type that defines a client preference service
type Service struct {
	preferenceRepo preference.IPreferenceRepository
	languageRepo   common.ILanguageRepository
	logger         *log.Logger
}

// NewClientPreferenceService is a function that returns a new client preference service
func NewClientPreferenceService(preferenceRepository preference.IPreferenceRepository,
	languageRepository common.ILanguageRepository, preferenceLogger *log.Logger) preference.IService {
	return &Service{preferenceRepo: preferenceRepository, languageRepo: languageRepository,
		logger: preferenceLogger}
}

// AddClientPreference is a method that adds a new client preference to the system
func (service *Service) AddClientPreference(newClientPreference *entity.ClientPreference) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started client preference adding process, Client Preference => %s",
		newClientPreference.ToString()), service.logger.Logs.ServerLogFile)

	err := service.preferenceRepo.Create(newClientPreference)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Client Preference => %s, %s",
			newClientPreference.ToString(), err.Error()))

		return errors.New("unable to add new client preference")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished client preference adding process, Client Preference => %s",
		newClientPreference.ToString()), service.logger.Logs.ServerLogFile)

	return nil
}

// ValidateClientPreference is a method that validates a client preference values.
// It checks if the client preference has a valid entries or not and return map of errors if any.
func (service *Service) ValidateClientPreference(clientPreference *entity.ClientPreference) entity.ErrMap {

	errMap := make(map[string]error)
	emptyLanguage, _ := regexp.MatchString(`^\s*$`, clientPreference.Language)
	if !emptyLanguage {
		language, err := service.languageRepo.Find(clientPreference.Language)
		if err != nil {
			errMap["language"] = errors.New("preferred language not found")
		} else {
			// Setting the correct language code if language name was provided
			clientPreference.Language = language.Code
		}
	} else {

		// Setting the default language code if non provided
		clientPreference.Language = entity.DefaultLanguage
	}

	if len(errMap) > 0 {
		return errMap
	}
	return nil
}

// FindClientPreference is a method that find and return a client preference that matchs the clientID value
func (service *Service) FindClientPreference(clientID string) (*entity.ClientPreference, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Client Preference finding process { Client ID : %s }", clientID),
		service.logger.Logs.ServerLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, clientID)
	if empty {
		return nil, errors.New("client preference not found")
	}

	clientPreference, err := service.preferenceRepo.Find(clientID)
	if err != nil {
		return nil, errors.New("client preference not found")
	}

	return clientPreference, nil
}

// UpdateClientPreference is a method that updates a certain client's preference
func (service *Service) UpdateClientPreference(clientPreference *entity.ClientPreference) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started client preference updating process, Client Preference => %s",
		clientPreference.ToString()), service.logger.Logs.ServerLogFile)

	err := service.preferenceRepo.Update(clientPreference)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating Client Preference => %s, %s",
			clientPreference.ToString(), err.Error()))

		return errors.New("unable to update client preference")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished client preference updating process, Client Preference => %s",
		clientPreference.ToString()), service.logger.Logs.ServerLogFile)

	return nil
}

// UpdateClientPreferenceSingleValue is a method that updates a single column entry of a client preference
func (service *Service) UpdateClientPreferenceSingleValue(clientID, columnName string, columnValue interface{}) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started single client preference value updating process "+
		"{ ClientID : %s, ColumnName : %s, ColumnValue : %s }", clientID, columnName, fmt.Sprint(columnValue)),
		service.logger.Logs.ServerLogFile)

	clientPreference := entity.ClientPreference{ClientID: clientID}
	err := service.preferenceRepo.UpdateValue(&clientPreference, columnName, columnValue)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating single client preference value "+
			"{ ClientID : %s, ColumnName : %s, ColumnValue : %s }, %s", clientID, columnName,
			fmt.Sprint(columnValue), err.Error()))

		return errors.New("unable to update client preference")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished single client preference value updating process, Client Preference => %s",
		clientPreference.ToString()), service.logger.Logs.ServerLogFile)

	return nil
}

// DeleteClientPreference is a method that deletes a certain client's preference
func (service *Service) DeleteClientPreference(clientID string) (*entity.ClientPreference, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started client preference deleting process { Client ID : %s }", clientID),
		service.logger.Logs.ServerLogFile)

	clientPreference, err := service.preferenceRepo.Delete(clientID)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting client preference { Client ID : %s }, %s",
			clientID, err.Error()))

		return nil, errors.New("unable to delete client preference")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished client preference deleting process, Deleted Client Preference => %s",
		clientPreference.ToString()), service.logger.Logs.ServerLogFile)

	return clientPreference, nil
}
