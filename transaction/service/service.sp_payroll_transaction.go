package service

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Benyam-S/onemembership/entity"
)

// AddSPPayrollTransaction is a method that adds a new service provider payroll transaction to the system
func (service *Service) AddSPPayrollTransaction(newPayrollTransaction *entity.SPPayrollTransaction) error {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started payroll transaction adding process, "+
		"SP Payroll Transaction => %s", newPayrollTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	err := service.spPayrollTransactionRepo.Create(newPayrollTransaction)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For adding SP Payroll Transaction => %s, %s",
			newPayrollTransaction.ToString(), err.Error()))

		return errors.New("unable to add new payroll transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished payroll transaction adding process, "+
		"SP Payroll Transaction => %s", newPayrollTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	return nil
}

// FindSPPayrollTransaction is a method that find and return a service provider payroll transaction that matches the id value
func (service *Service) FindSPPayrollTransaction(id string) (*entity.SPPayrollTransaction, error) {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Single payroll transaction finding process { SP Payroll Transaction ID : %s }", id),
		service.logger.Logs.TransactionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, id)
	if empty {
		return nil, errors.New("no payroll transaction found")
	}

	payrollTransaction, err := service.spPayrollTransactionRepo.Find(id)
	if err != nil {
		return nil, errors.New("no payroll transaction found")
	}
	return payrollTransaction, nil
}

// FindMultipleSPPayrollTransactions is a method that find and return multiple service provider payroll transactions that matchs the identifier value
func (service *Service) FindMultipleSPPayrollTransactions(providerID string) []*entity.SPPayrollTransaction {

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple payroll transaction finding process { Provider ID : %s }",
		providerID), service.logger.Logs.TransactionLogFile)

	empty, _ := regexp.MatchString(`^\s*$`, providerID)
	if empty {
		return []*entity.SPPayrollTransaction{}
	}

	return service.spPayrollTransactionRepo.FindMultiple(providerID)
}

// UpdateSPPayrollTransaction is a method that updates a service provider payroll transaction in the system
func (service *Service) UpdateSPPayrollTransaction(payrollTransaction *entity.SPPayrollTransaction) error {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started payroll transaction updating process, SP Payroll Transaction => %s",
		payrollTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	err := service.spPayrollTransactionRepo.Update(payrollTransaction)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(fmt.Sprintf("Error: For updating SP Payroll Transaction => %s, %s",
			payrollTransaction.ToString(), err.Error()))

		return errors.New("unable to update payroll transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished payroll transaction updating process, SP Payroll Transaction => %s",
		payrollTransaction.ToString()), service.logger.Logs.TransactionLogFile)

	return nil
}

// DeleteSPPayrollTransaction is a method that deletes a service provider payroll transaction from the system using an id
func (service *Service) DeleteSPPayrollTransaction(id string) (*entity.SPPayrollTransaction, error) {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Started payroll transaction deleting process { SP Payroll Transaction ID : %s }",
		id), service.logger.Logs.TransactionLogFile)

	payrollTransaction, err := service.spPayrollTransactionRepo.Delete(id)
	if err != nil {
		/* ---------------------------- Logging ---------------------------- */
		service.logger.LogToErrorFile(
			fmt.Sprintf("Error: For deleting payroll transaction { SP Payroll Transaction ID : %s }, %s",
				id, err.Error()))

		return nil, errors.New("unable to delete payroll transaction")
	}

	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Finished payroll transaction deleting process, "+
		"Deleted SP Payroll Transaction => %s", payrollTransaction.ToString()), service.logger.Logs.TransactionLogFile)
	return payrollTransaction, nil
}

// DeleteMultipleSPPayrollTransactions is a method that deletes multiple service provider payroll transactions from the system that match the given identifier
func (service *Service) DeleteMultipleSPPayrollTransactions(providerID string) []*entity.SPPayrollTransaction {
	/* ---------------------------- Logging ---------------------------- */
	service.logger.Log(fmt.Sprintf("Multiple payroll transaction deleting { Provider ID : %s }",
		providerID), service.logger.Logs.TransactionLogFile)

	return service.spPayrollTransactionRepo.DeleteMultiple(providerID)
}
