package entity

import (
	"time"
)

// DeletedUser is a type that defines a user that has been deleted
// This struct is used to store and identify a pervious user
type DeletedUser struct {
	ID          string `gorm:"primary_key; unique;"`
	UserID      string
	FirstName   string
	LastName    string
	Username    string
	PhoneNumber string
	Email       string
	CreatedAt   time.Time
}

// DeletedServiceProvider is a type that defines a service provider that has been deleted
// This struct is used to store and identify a pervious service provider
type DeletedServiceProvider struct {
	ID          string `gorm:"primary_key; unique;"`
	ProviderID  string
	FirstName   string
	LastName    string
	UserName    string
	PhoneNumber string
	Email       string
	CreatedAt   time.Time
}

// DeletedSPSubscriptionTransaction (DeletedServiceProviderSubscriptionTransaction) is a type that defines
// a deleted transaction performed during subscription by service provider
type DeletedSPSubscriptionTransaction struct {
	ID             string
	ProviderID     string
	PlanID         string
	AppID          string
	ReceiverName   string
	Subject        string
	ReceivedAmount float64
	TransactionFee float64
	CurrencyType   string
	TimeoutExpress int64
	Nonce          string
	OutTradeNo     string
	TradeNo        string
	Status         string
	CreatedAt      time.Time
}

// DeletedSPPayrollTransaction (DeletedServiceProviderSubscriptionTransaction) is a type that defines
// a deleted transaction performed during payroll to service provider
type DeletedSPPayrollTransaction struct {
	ID                    string
	ProviderID            string
	PayedAmount           float64
	LinkedAccount         string
	LinkedAccountProvider string
	Status                string
	CreatedAt             time.Time
}

// DeletedSubscriptionTransaction is a type that defines a deleted transaction performed during subscription
type DeletedSubscriptionTransaction struct {
	ID             string
	UserID         string
	PlanID         string
	AppID          string
	ReceiverName   string
	Subject        string
	ReceivedAmount float64
	TransactionFee float64
	CurrencyType   string
	TimeoutExpress int64
	Nonce          string
	OutTradeNo     string
	TradeNo        string
	Status         string
	InitiatedFrom  string
	CreatedAt      time.Time
}
