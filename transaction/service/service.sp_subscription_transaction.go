package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
)

// AddSPSubscriptionTransaction is a method that adds a new service provider subscription transaction to the system
func (service *Service) AddSPSubscriptionTransaction(newSubscriptionTransaction *entity.SPSubscriptionTransaction) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription transaction adding process, "+
		"SP Subscription Transaction => %s", newSubscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	err := service.spSubscriptionTransactionRepo.Create(newSubscriptionTransaction)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding SP Subscription Transaction => %s, %s",
			newSubscriptionTransaction.ToString(), err.Error()))

		return errors.New("unable to add new subscription transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription transaction adding process, "+
		"SP Subscription Transaction => %s", newSubscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	return nil
}

// ValidateSPSubscriptionTransaction is a method that validates a service provider subscription transaction entries.
// It checks if the service provider subscription transaction has a valid entries or not and return map of errors if any.
func (service *Service) ValidateSPSubscriptionTransaction(subscriptionTransaction *entity.SPSubscriptionTransaction) entity.ErrMap {

	errMap := make(map[string]error)

	isValidSubject, _ := regexp.MatchString(`^\w+[\s\w]*$`, subscriptionTransaction.Subject)
	if !isValidSubject {
		errMap["subject"] = errors.New(`subscription transaction subject should not contain any special characters`)
	}

	// Checking uniqueness in both subscription_transactions and sp_subscription_transactions
	isUniqueNonce1 := service.cmService.IsUnique("nonce", subscriptionTransaction.Nonce, "subscription_transactions")
	isUniqueNonce2 := service.cmService.IsUnique("nonce", subscriptionTransaction.Nonce, "sp_subscription_transactions")
	if !isUniqueNonce1 || !isUniqueNonce2 {
		errMap["nonce"] = errors.New(`subscription transaction nonce should be unique`)
	}

	// Checking uniqueness in both subscription_transactions and sp_subscription_transactions
	isUniqueOutTradeNo1 := service.cmService.IsUnique("out_trade_no",
		subscriptionTransaction.OutTradeNo, "subscription_transactions")
	isUniqueOutTradeNo2 := service.cmService.IsUnique("out_trade_no",
		subscriptionTransaction.OutTradeNo, "sp_subscription_transactions")
	if !isUniqueOutTradeNo1 || !isUniqueOutTradeNo2 {
		errMap["out_trade_no"] = errors.New(`subscription transaction out trade number should be unique`)
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindSPSubscriptionTransaction is a method that find and return a service provider subscription transaction that matches the id value
func (service *Service) FindSPSubscriptionTransaction(id string) (*entity.SPSubscriptionTransaction, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single subscription transaction finding process "+
		"{ SP Subscription Transaction ID : %s }", id), service.logger.Logs.TransactionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, id)
	if empty {
		return nil, errors.New("no subscription transaction found")
	}

	subscriptionTransaction, err := service.spSubscriptionTransactionRepo.Find(id)
	if err != nil {
		return nil, errors.New("no subscription transaction found")
	}
	return subscriptionTransaction, nil
}

// FindMultipleSPSubscriptionTransactions is a method that find and return multiple service provider subscription transactions that matchs the identifier value
func (service *Service) FindMultipleSPSubscriptionTransactions(identifier string) []*entity.SPSubscriptionTransaction {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple subscription transaction finding process "+
		"{ SP Subscription Transaction Identifier : %s }", identifier), service.logger.Logs.TransactionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return []*entity.SPSubscriptionTransaction{}
	}

	return service.spSubscriptionTransactionRepo.FindMultiple(identifier)
}

// UpdateSPSubscriptionTransaction is a method that updates a service provider subscription transaction in the system
func (service *Service) UpdateSPSubscriptionTransaction(subscriptionTransaction *entity.SPSubscriptionTransaction) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription transaction updating process, SP Subscription Transaction => %s",
		subscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	err := service.spSubscriptionTransactionRepo.Update(subscriptionTransaction)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating SP Subscription Transaction => %s, %s",
			subscriptionTransaction.ToString(), err.Error()))

		return errors.New("unable to update subscription transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription transaction updating process, SP Subscription Transaction => %s",
		subscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	return nil
}

// DeleteSPSubscriptionTransaction is a method that deletes a service provider subscription transaction from the system using an id
func (service *Service) DeleteSPSubscriptionTransaction(id string) (*entity.SPSubscriptionTransaction, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription transaction deleting process { SP Subscription Transaction ID : %s }",
		id), service.logger.Logs.TransactionLogFile)

	subscriptionTransaction, err := service.spSubscriptionTransactionRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(
			fmt.Sprintf("Error: For deleting subscription transaction { SP Subscription Transaction ID : %s }, %s",
				id, err.Error()))

		return nil, errors.New("unable to delete subscription transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription transaction deleting process, "+
		"Deleted SP Subscription Transaction => %s", subscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)
	return subscriptionTransaction, nil
}

// DeleteMultipleSPSubscriptionTransactions is a method that deletes multiple service provider subscription transactions from the system that match the given identifier
func (service *Service) DeleteMultipleSPSubscriptionTransactions(identifier string) []*entity.SPSubscriptionTransaction {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple subscription transaction deleting { SP Subscription Transaction Identifier : %s }",
		identifier), service.logger.Logs.TransactionLogFile)

	return service.spSubscriptionTransactionRepo.DeleteMultiple(identifier)
}
