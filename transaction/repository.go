package transaction

import (
	"github.com/Benyam-S/onemembership/entity"
)

// IPaymentGatewayRepository is an interface that defines all the repository methods of a payment gateway struct
type IPaymentGatewayRepository interface {
	Create(newPaymentGateway *entity.PaymentGateway) error
	Find(id int64) (*entity.PaymentGateway, error)
	All() []*entity.PaymentGateway
	Update(paymentGateway *entity.PaymentGateway) error
	Delete(id int64) (*entity.PaymentGateway, error)
}

// ISubscriptionTransactionRepository is an interface that defines all the repository methods of a subscription transaction struct
type ISubscriptionTransactionRepository interface {
	Create(newTransaction *entity.SubscriptionTransaction) error
	Find(identifier string) (*entity.SubscriptionTransaction, error)
	FindMultiple(identifier string) []*entity.SubscriptionTransaction
	Update(transaction *entity.SubscriptionTransaction) error
	Delete(id string) (*entity.SubscriptionTransaction, error)
	DeleteMultiple(identifier string) []*entity.SubscriptionTransaction
}

// ISPSubscriptionTransactionRepository is an interface that defines all the repository methods of a service provider subscription transaction struct
type ISPSubscriptionTransactionRepository interface {
	Create(newTransaction *entity.SPSubscriptionTransaction) error
	Find(identifier string) (*entity.SPSubscriptionTransaction, error)
	FindMultiple(identifier string) []*entity.SPSubscriptionTransaction
	Update(transaction *entity.SPSubscriptionTransaction) error
	Delete(id string) (*entity.SPSubscriptionTransaction, error)
	DeleteMultiple(identifier string) []*entity.SPSubscriptionTransaction
}

// ISPPayrollTransactionRepository is an interface that defines all the repository methods of a service provider payroll transaction struct
type ISPPayrollTransactionRepository interface {
	Create(newTransaction *entity.SPPayrollTransaction) error
	Find(id string) (*entity.SPPayrollTransaction, error)
	FindMultiple(providerID string) []*entity.SPPayrollTransaction
	Update(transaction *entity.SPPayrollTransaction) error
	Delete(id string) (*entity.SPPayrollTransaction, error)
	DeleteMultiple(providerID string) []*entity.SPPayrollTransaction
}
