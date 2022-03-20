package entity

import (
	"encoding/json"
	"fmt"
)

// ToString is a method that converts a User struct to readable JSON string format
func (user *User) ToString() string {
	output, err := json.Marshal(user)
	if err != nil {
		return fmt.Sprint(user)
	}

	return string(output)
}

// ToString is a method that converts a User Password struct to readable JSON string format
func (userPassword *UserPassword) ToString() string {
	output, err := json.Marshal(userPassword)
	if err != nil {
		return fmt.Sprint(userPassword)
	}

	return string(output)
}

// ToString is a method that converts a Service Provider struct to readable JSON string format
func (serviceProvider *ServiceProvider) ToString() string {
	output, err := json.Marshal(serviceProvider)
	if err != nil {
		return fmt.Sprint(serviceProvider)
	}

	return string(output)
}

// ToString is a method that converts a Service Provider Password struct to readable JSON string format
func (serviceProviderPassword *SPPassword) ToString() string {
	output, err := json.Marshal(serviceProviderPassword)
	if err != nil {
		return fmt.Sprint(serviceProviderPassword)
	}

	return string(output)
}

// ToString is a method that converts a Service Provider Wallet struct to readable JSON string format
func (spWallet *SPWallet) ToString() string {
	output, err := json.Marshal(spWallet)
	if err != nil {
		return fmt.Sprint(spWallet)
	}

	return string(output)
}

// ToString is a method that converts a Service Provider Subscription struct to readable JSON string format
func (spSubscription *SPSubscription) ToString() string {
	output, err := json.Marshal(spSubscription)
	if err != nil {
		return fmt.Sprint(spSubscription)
	}

	return string(output)
}

// ToString is a method that converts a Service Provider Subscription Transaction struct to readable JSON string format
func (spSubscriptionTransaction *SPSubscriptionTransaction) ToString() string {
	output, err := json.Marshal(spSubscriptionTransaction)
	if err != nil {
		return fmt.Sprint(spSubscriptionTransaction)
	}

	return string(output)
}

// ToString is a method that converts a Service Provider Payroll Transaction struct to readable JSON string format
func (spPayrollTransaction *SPPayrollTransaction) ToString() string {
	output, err := json.Marshal(spPayrollTransaction)
	if err != nil {
		return fmt.Sprint(spPayrollTransaction)
	}

	return string(output)
}

// ToString is a method that converts a Service Provider Subscription Plan struct to readable JSON string format
func (spSubscriptionPlan *SPSubscriptionPlan) ToString() string {
	output, err := json.Marshal(spSubscriptionPlan)
	if err != nil {
		return fmt.Sprint(spSubscriptionPlan)
	}

	return string(output)
}

// ToString is a method that converts a Subscription struct to readable JSON string format
func (subscription *Subscription) ToString() string {
	output, err := json.Marshal(subscription)
	if err != nil {
		return fmt.Sprint(subscription)
	}

	return string(output)
}

// ToString is a method that converts a Subscription Transaction struct to readable JSON string format
func (subscriptionTransaction *SubscriptionTransaction) ToString() string {
	output, err := json.Marshal(subscriptionTransaction)
	if err != nil {
		return fmt.Sprint(subscriptionTransaction)
	}

	return string(output)
}

// ToString is a method that converts a Subscription Plan struct to readable JSON string format
func (subscriptionPlan *SubscriptionPlan) ToString() string {
	output, err := json.Marshal(subscriptionPlan)
	if err != nil {
		return fmt.Sprint(subscriptionPlan)
	}

	return string(output)
}

// ToString is a method that converts a PlanChatLink struct to readable JSON string format
func (planChatLink *PlanChatLink) ToString() string {
	output, err := json.Marshal(planChatLink)
	if err != nil {
		return fmt.Sprint(planChatLink)
	}

	return string(output)
}

// ToString is a method that converts a UserChatLink struct to readable JSON string format
func (userChatLink *UserChatLink) ToString() string {
	output, err := json.Marshal(userChatLink)
	if err != nil {
		return fmt.Sprint(userChatLink)
	}

	return string(output)
}

// ToString is a method that converts a Project struct to readable JSON string format
func (project *Project) ToString() string {
	output, err := json.Marshal(project)
	if err != nil {
		return fmt.Sprint(project)
	}

	return string(output)
}

// ToString is a method that converts a ProjectChatLink struct to readable JSON string format
func (projectChatLink *ProjectChatLink) ToString() string {
	output, err := json.Marshal(projectChatLink)
	if err != nil {
		return fmt.Sprint(projectChatLink)
	}

	return string(output)
}

// ToString is a method that converts a PaymentGateway struct to readable JSON string format
func (paymentGateway *PaymentGateway) ToString() string {
	output, err := json.Marshal(paymentGateway)
	if err != nil {
		return fmt.Sprint(paymentGateway)
	}

	return string(output)
}

// ToString is a method that converts a Feedback struct to readable JSON string format
func (feedback *Feedback) ToString() string {
	output, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Sprint(feedback)
	}

	return string(output)
}

// ToString is a method that converts a Deleted User struct to readable JSON string format
func (deletedUser *DeletedUser) ToString() string {
	output, err := json.Marshal(deletedUser)
	if err != nil {
		return fmt.Sprint(deletedUser)
	}

	return string(output)
}

// ToString is a method that converts a Deleted Service Provider struct to readable JSON string format
func (deletedServiceProvider *DeletedServiceProvider) ToString() string {
	output, err := json.Marshal(deletedServiceProvider)
	if err != nil {
		return fmt.Sprint(deletedServiceProvider)
	}

	return string(output)
}

// ToString is a method that converts a Deleted Service Provider Subscription Transaction struct to readable JSON string format
func (deletedSPSubscriptionTransaction *DeletedSPSubscriptionTransaction) ToString() string {
	output, err := json.Marshal(deletedSPSubscriptionTransaction)
	if err != nil {
		return fmt.Sprint(deletedSPSubscriptionTransaction)
	}

	return string(output)
}

// ToString is a method that converts a Deleted Service Provider Payroll Transaction struct to readable JSON string format
func (deletedSPPayrollTransaction *DeletedSPPayrollTransaction) ToString() string {
	output, err := json.Marshal(deletedSPPayrollTransaction)
	if err != nil {
		return fmt.Sprint(deletedSPPayrollTransaction)
	}

	return string(output)
}

// ToString is a method that converts a Deleted User Subscription Transaction struct to readable JSON string format
func (deletedSubscriptionTransaction *DeletedSubscriptionTransaction) ToString() string {
	output, err := json.Marshal(deletedSubscriptionTransaction)
	if err != nil {
		return fmt.Sprint(deletedSubscriptionTransaction)
	}

	return string(output)
}

// ToString is a method that converts a Client Preference struct to readable JSON string format
func (preference *ClientPreference) ToString() string {
	output, err := json.Marshal(preference)
	if err != nil {
		return fmt.Sprint(preference)
	}

	return string(output)
}

// ToString is a method that converts a Language Entry struct to readable JSON string format
func (language *LanguageEntry) ToString() string {
	output, err := json.Marshal(language)
	if err != nil {
		return fmt.Sprint(language)
	}

	return string(output)
}
