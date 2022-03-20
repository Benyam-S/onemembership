package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
)

// AddSubscriptionTransaction is a method that adds a new subscription transaction to the system
func (service *Service) AddSubscriptionTransaction(newSubscriptionTransaction *entity.SubscriptionTransaction) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription transaction adding process, Subscription Transaction => %s",
		newSubscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	err := service.subTransactionRepo.Create(newSubscriptionTransaction)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding Subscription Transaction => %s, %s",
			newSubscriptionTransaction.ToString(), err.Error()))

		return errors.New("unable to add new subscription transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription transaction adding process, Subscription Transaction => %s",
		newSubscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	return nil
}

// ValidateSubscriptionTransaction is a method that validates a subscription transaction entries.
// It checks if the subscription transaction has a valid entries or not and return map of errors if any.
func (service *Service) ValidateSubscriptionTransaction(subscriptionTransaction *entity.SubscriptionTransaction) entity.ErrMap {

	errMap := make(map[string]error)

	// Checking uniqueness in subscription_transactions
	isUniqueNonce := service.cmService.IsUnique("nonce", subscriptionTransaction.Nonce, "subscription_transactions")
	if !isUniqueNonce {
		errMap["nonce"] = errors.New(`subscription transaction nonce should be unique`)
	}

	// Checking uniqueness in subscription_transactions
	isUniqueOutTradeNo := service.cmService.IsUnique("out_trade_no",
		subscriptionTransaction.OutTradeNo, "subscription_transactions")
	if !isUniqueOutTradeNo {
		errMap["out_trade_no"] = errors.New(`subscription transaction out trade number should be unique`)
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindSubscriptionTransaction is a method that find and return a subscription transaction that matches the transaction id value
func (service *Service) FindSubscriptionTransaction(id string) (*entity.SubscriptionTransaction, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single subscription transaction finding process { Subscription Transaction ID : %s }", id),
		service.logger.Logs.TransactionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, id)
	if empty {
		return nil, errors.New("no subscription transaction found")
	}

	subscriptionTransaction, err := service.subTransactionRepo.Find(id)
	if err != nil {
		return nil, errors.New("no subscription transaction found")
	}
	return subscriptionTransaction, nil
}

// FindMultipleSubscriptionTransactions is a method that find and return multiple subscription transactions that matchs the identifier value
func (service *Service) FindMultipleSubscriptionTransactions(identifier string) []*entity.SubscriptionTransaction {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple subscription transaction finding process { Subscription Transaction Identifier : %s }",
		identifier), service.logger.Logs.TransactionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, identifier)
	if empty {
		return []*entity.SubscriptionTransaction{}
	}

	return service.subTransactionRepo.FindMultiple(identifier)
}

// UpdateSubscriptionTransaction is a method that updates a subscription transaction in the system
func (service *Service) UpdateSubscriptionTransaction(subscriptionTransaction *entity.SubscriptionTransaction) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription transaction updating process, Subscription Transaction => %s",
		subscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	err := service.subTransactionRepo.Update(subscriptionTransaction)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating Subscription Transaction => %s, %s",
			subscriptionTransaction.ToString(), err.Error()))

		return errors.New("unable to update subscription transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription transaction updating process, Subscription Transaction => %s",
		subscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	return nil
}

// DeleteSubscriptionTransaction is a method that deletes a subscription transaction from the system using an id
func (service *Service) DeleteSubscriptionTransaction(id string) (*entity.SubscriptionTransaction, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started subscription transaction deleting process { Subscription Transaction ID : %s }",
		id), service.logger.Logs.TransactionLogFile)

	subscriptionTransaction, err := service.subTransactionRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For deleting subscription transaction "+
			"{ Subscription Transaction ID : %s }, %s", id, err.Error()))

		return nil, errors.New("unable to delete subscription transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished subscription transaction deleting process, Deleted Subscription Transaction => %s",
		subscriptionTransaction.ToString()), service.logger.Logs.TransactionLogFile)
	return subscriptionTransaction, nil
}

// DeleteMultipleSubscriptionTransactions is a method that deletes multiple subscription transactions from the system that match the given identifier
func (service *Service) DeleteMultipleSubscriptionTransactions(identifier string) []*entity.SubscriptionTransaction {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple subscription transaction deleting { Subscription Transaction Identifier : %s }",
		identifier), service.logger.Logs.TransactionLogFile)

	return service.subTransactionRepo.DeleteMultiple(identifier)
}
