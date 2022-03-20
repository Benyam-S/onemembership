package transaction

import "github.com/Benyam-S/onemembership/entity"

// TelebirrAPIAccount is a struct that defines all the need entries for telebirr api
type TelebirrAPIAccount struct {
	AccessPoint    string  `json:"api_access_point"`
	AppID          string  `json:"app_id"`
	AppKey         string  `json:"app_key"`
	NotifyURL      string  `json:"notify_url"`
	ReturnURL      string  `json:"return_url"`
	ShortCode      string  `json:"short_code"`
	TransactionFee float64 `json:"transaction_fee"`
	PublicKey      []byte  `json:"-"`
}

// IService is an interface that defines all the service methods of a project struct
type IService interface {
	AddPaymentGateway(newPaymentGateway *entity.PaymentGateway) error
	ValidatePaymentGateway(paymentGateway *entity.PaymentGateway) entity.ErrMap
	FindPaymentGateway(id int64) (*entity.PaymentGateway, error)
	AllPaymentGateways() []*entity.PaymentGateway
	UpdatePaymentGateway(paymentGateway *entity.PaymentGateway) error
	DeletePaymentGateway(id int64) (*entity.PaymentGateway, error)

	AddSubscriptionTransaction(newSubscriptionTransaction *entity.SubscriptionTransaction) error
	ValidateSubscriptionTransaction(subscriptionTransaction *entity.SubscriptionTransaction) entity.ErrMap
	FindSubscriptionTransaction(id string) (*entity.SubscriptionTransaction, error)
	FindMultipleSubscriptionTransactions(identifier string) []*entity.SubscriptionTransaction
	UpdateSubscriptionTransaction(subscriptionTransaction *entity.SubscriptionTransaction) error
	DeleteSubscriptionTransaction(id string) (*entity.SubscriptionTransaction, error)
	DeleteMultipleSubscriptionTransactions(identifier string) []*entity.SubscriptionTransaction

	AddSPSubscriptionTransaction(newSubscriptionTransaction *entity.SPSubscriptionTransaction) error
	ValidateSPSubscriptionTransaction(subscriptionTransaction *entity.SPSubscriptionTransaction) entity.ErrMap
	FindSPSubscriptionTransaction(id string) (*entity.SPSubscriptionTransaction, error)
	FindMultipleSPSubscriptionTransactions(identifier string) []*entity.SPSubscriptionTransaction
	UpdateSPSubscriptionTransaction(subscriptionTransaction *entity.SPSubscriptionTransaction) error
	DeleteSPSubscriptionTransaction(id string) (*entity.SPSubscriptionTransaction, error)
	DeleteMultipleSPSubscriptionTransactions(identifier string) []*entity.SPSubscriptionTransaction

	AddSPPayrollTransaction(newPayrollTransaction *entity.SPPayrollTransaction) error
	FindSPPayrollTransaction(id string) (*entity.SPPayrollTransaction, error)
	FindMultipleSPPayrollTransactions(providerID string) []*entity.SPPayrollTransaction
	UpdateSPPayrollTransaction(payrollTransaction *entity.SPPayrollTransaction) error
	DeleteSPPayrollTransaction(id string) (*entity.SPPayrollTransaction, error)
	DeleteMultipleSPPayrollTransactions(providerID string) []*entity.SPPayrollTransaction

	GetTelebirrH5WebURL(userID, planID, receiverName, subject, currencyType, initiatedFrom string,
		receivedAmount float64) (string, error)
}
