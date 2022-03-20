package entity

import (
	"time"
)

// PaymentGateway is a type that defines all the available payment gateways or gateways
type PaymentGateway struct {
	ID        int64 `gorm:"primary_key; unique; auto_increment;"`
	Name      string
	CreatedAt time.Time
}

// SubscriptionTransaction is a type that defines a transaction performed during subscription
type SubscriptionTransaction struct {
	ID             string `gorm:"primary_key; unique;"`
	UserID         string
	PlanID         string
	AppID          string // Used to identify which gateway was used
	ReceiverName   string
	Subject        string
	ReceivedAmount float64 // Indicates the amount that will be provided to service provider
	TransactionFee float64 // The amount that will be deducted from the total amount where TotalAmount = ReceivedAmount + Fee
	CurrencyType   string
	TimeoutExpress int64
	Nonce          string
	OutTradeNo     string
	TradeNo        string
	Status         string // Can be used to identify the status of the transaction
	InitiatedFrom  string // Indicates from which interface the request was initiated such as from bot or web
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SPSubscriptionTransaction (ServiceProviderSubscriptionTransaction) is a type that defines
// a transaction performed during subscription by service provider
type SPSubscriptionTransaction struct {
	ID             string `gorm:"primary_key; unique;"`
	ProviderID     string
	PlanID         string
	AppID          string // Used to identify which gateway was used
	ReceiverName   string
	Subject        string
	ReceivedAmount float64 // Indicates the amount that will be provided to service provider
	TransactionFee float64 // The amount that will be deducted from the total amount where TotalAmount = ReceivedAmount + Fee
	CurrencyType   string
	TimeoutExpress int64
	Nonce          string
	OutTradeNo     string
	TradeNo        string
	Status         string // Can be used to identify the status of the transaction
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SPSubscriptionTransaction (ServiceProviderSubscriptionTransaction) is a type that defines
// a transaction performed during payroll to service provider
type SPPayrollTransaction struct {
	ID                    string `gorm:"primary_key; unique;"`
	ProviderID            string
	PayedAmount           float64
	LinkedAccount         string // Indicates to which account payment was made
	LinkedAccountProvider string // Indicates to which account provider payment was made
	Status                string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
