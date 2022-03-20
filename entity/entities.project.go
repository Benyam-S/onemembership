package entity

import (
	"time"
)

// Project is a type that defines a project that a service provider can own
type Project struct {
	ID          string `gorm:"primary_key; unique;"`
	ProviderID  string
	Name        string `gorm:"type:blob;"` // The name may contain special characters or emojis
	Description string `gorm:"type:blob;"` // The description may contain special characters or emojis
	ProjectLink string
	Status      string // To identify the project has been active, disabled or compeleted or not
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ProjectChatLink is a type that defines a link between project and telegram chat
type ProjectChatLink struct {
	ProjectID string `gorm:"unique_index:unique_project_to_chat_link_relation;"` // Defining composite unique key
	ChatID    int64  `gorm:"unique_index:unique_project_to_chat_link_relation;"` // Defining composite unique key
	Type      string // Can be used to identify whether it is a channel or group
}

// SubscriptionPlan is a type that defines a subscription plan for a given project
type SubscriptionPlan struct {
	ID          string `gorm:"primary_key; unique;"`
	ProjectID   string
	Name        string
	Benfits     string `gorm:"type:blob;"` // The benfits may contain special characters or emojis
	Duration    int64  // Represents the number of days the subscription plan lengthen
	Price       float64
	Currency    string
	IsRecurring bool
	Status      string // Can be used to identify the status of the plan in order to take action
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PlanChatLink is a type that defines a link between subscription plan and telegram chat
type PlanChatLink struct {
	PlanID string `gorm:"unique_index:unique_plan_to_chat_link_relation;"` // Defining composite unique key
	ChatID int64  `gorm:"unique_index:unique_plan_to_chat_link_relation;"` // Defining composite unique key
}

// UserChatLink is a type that defines the chat link generated for the user
type UserChatLink struct {
	UserID     string
	PlanID     string
	ChatID     int64
	InviteLink string
	CreatedAt  time.Time
}

// Subscription is a type that defines user subscription and subscription history
// Contains all the information needed for reserving history
type Subscription struct {
	ID string `gorm:"primary_key; unique;"`

	// For storing the current transaction subscriber detail
	SubscriberID          string
	SubscriberFirstName   string
	SubscriberLastName    string
	SubscriberUserName    string
	SubscriberPhoneNumber string
	SubscriberEmail       string

	// For storing the current transaction subscription provider detail
	ProviderID          string
	ProviderFirstName   string
	ProviderLastName    string
	ProviderUserName    string
	ProviderPhoneNumber string
	ProviderEmail       string

	// For storing the current project detail
	ProjectID          string
	ProjectName        string `gorm:"type:blob;"`
	ProjectDescription string `gorm:"type:blob;"`
	ProjectLink        string

	// For storing the current transaction subscription plan detail
	SubscriptionPlanID          string
	SubscriptionPlanName        string
	SubscriptionPlanBenfits     string `gorm:"type:blob;"`
	SubscriptionPlanDuration    int64
	SubscriptionPlanPrice       float64
	SubscriptionPlanIsRecurring bool
	SubscriptionPlanCurrency    string

	// TimeStamp for the created history or subscription
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
