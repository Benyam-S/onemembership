package entity

import (
	"net/http"
	"time"
)

// User is a type that defines a user group the will use or subscribe to subscription plans
type User struct {
	ID          string `gorm:"primary_key; unique;"`
	FirstName   string
	LastName    string
	UserName    string
	PhoneNumber string
	Email       string
	Preference  *ClientPreference
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserPassword is a type that defines a user password
type UserPassword struct {
	UserID    string `gorm:"primary_key; unique;"`
	Password  string
	Salt      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ServiceProvider is a type that defines the service provider
type ServiceProvider struct {
	ID          string `gorm:"primary_key; unique;"`
	FirstName   string
	LastName    string
	UserName    string
	PhoneNumber string
	Email       string
	Preference  *ClientPreference
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SPPassword (ServiceProviderPassword) is a type that defines a user password
type SPPassword struct {
	ProviderID string `gorm:"primary_key; unique;"`
	Password   string
	Salt       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// SPWallet (ServiceProviderWallet) is a type that defines the service provider wallet account
type SPWallet struct {
	ProviderID    string  `gorm:"primary_key; unique;"`
	RunningAmount float64 // Amount that is actively counting the received value,
	//it is estimated value it mightn't be incorrect so we should recalculate it from the transactions on withdraw

	PendingAmount         float64 // Amount requested to be withdrawn
	LinkedAccount         string  // Third party banking account saved by service provider
	LinkedAccountProvider string  // Third party banking account provider
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// SPSubscription (ServiceProviderSubscription) is a type that defines service provider subscription
type SPSubscription struct {
	ProviderID string `gorm:"primary_key; unique;"`

	// For storing the current subscription plan detail
	SubscriptionPlanID       string
	SubscriptionPlanName     string
	SubscriptionPlanDuration int64
	SubscriptionPlanPrice    float64
	SubscriptionPlanCurrency string

	// TimeStamp for the created subscription
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

// SPSubscriptionPlan (ServiceProviderSubscriptionPlan) is a type that defines a subscription plan for service providers
// It will be provided by the administrators
type SPSubscriptionPlan struct {
	ID        string `gorm:"primary_key; unique;"`
	Name      string
	Duration  int64 // Represents the number of days the subscription plan lengthen
	Price     float64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Feedback is a type that defines client feedback
type Feedback struct {
	ID        string `gorm:"primary_key; unique;"`
	ClientID  string
	Comment   string `gorm:"type:blob;"`
	Seen      bool
	CreatedAt time.Time
}

// ClientPreference is a type that defines a onemembership client preference
type ClientPreference struct {
	ClientID string `gorm:"primary_key; unique;"`
	Language string `gorm:"default: 'en'"` // The same as the DefaultLanguage constant
}

// Language is a type that defines the langauges available in the system
type Language struct {
	Code         string `gorm:"primary_key; unique;"`
	Name         string `gorm:"unique;"`
	Flag         string `gorm:"type:blob"`
	DisplayOrder int64  // This is only used for display purpose
}

// LanguageEntry is a struct that defines the different language entries that are supported by the system
type LanguageEntry struct {
	ID         int64  `gorm:"primary_key; auto_increment; unique;"`
	Identifier string `gorm:"unique_index:unique_language_entry;"` // Defining composite unique key
	Value      string `gorm:"type:blob;"`                          // Value may contain special characters
	Code       string `gorm:"unique_index:unique_language_entry;"` // Defining composite unique key
}

// SystemConfig is a type that defines a server system configuration file
type SystemConfig struct {
	RedisClient          map[string]string `json:"redis_client"`
	MysqlClient          map[string]string `json:"mysql_client"`
	CookieName           string            `json:"cookie_name"`
	SecretKey            string            `json:"secret_key"`
	SuperAdminEmail      string            `json:"super_admin_email"`
	HTTPDomainAddress    string            `json:"http_domain_address"`
	BotDomainAddress     string            `json:"bot_domain_address"`
	BotClientServerPort  string            `json:"bot_client_server_port"`
	HTTPClientServerPort string            `json:"http_client_server_port"`
	LogsPath             string            `json:"logs_path"`
	ArchivesPath         string            `json:"archives_path"`
	Logs                 map[string]string `json:"logs"`
}

// Key is a type that defines a key type that can be used a key value in context
type Key string

// ErrMap is a type that defines a map with string identifier and it's error
type ErrMap map[string]error

// StringMap is a method that returns string map corresponding to the ErrMap where the error type is converted to a string
func (errMap ErrMap) StringMap() map[string]string {
	stringMap := make(map[string]string)
	for key, value := range errMap {
		stringMap[key] = value.Error()
	}

	return stringMap
}

// Middleware is a type that defines a function that takes a handler func and return a new handler func type
type Middleware func(http.HandlerFunc) http.HandlerFunc
