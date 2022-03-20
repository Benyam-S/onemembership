package subscriptionplan

import "github.com/Benyam-S/onemembership/entity"

// IService is an interface that defines all the service methods of a subscription plan struct
type IService interface {
	AddSubscriptionPlan(newSubscriptionPlan *entity.SubscriptionPlan) error
	ValidateSubscriptionPlan(subscriptionPlan *entity.SubscriptionPlan) entity.ErrMap
	FindSubscriptionPlan(id string) (*entity.SubscriptionPlan, error)
	FindMultipleSubscriptionPlans(projectID string) []*entity.SubscriptionPlan
	UpdateSubscriptionPlan(subscriptionPlan *entity.SubscriptionPlan) error
	DeleteSubscriptionPlan(id string) (*entity.SubscriptionPlan, error)
	DeleteMultipleSubscriptionPlans(projectID string) []*entity.SubscriptionPlan

	AddSPSubscriptionPlan(newSubscriptionPlan *entity.SPSubscriptionPlan) error
	ValidateSPSubscriptionPlan(subscriptionPlan *entity.SPSubscriptionPlan) entity.ErrMap
	FindSPSubscriptionPlan(id string) (*entity.SPSubscriptionPlan, error)
	UpdateSPSubscriptionPlan(subscriptionPlan *entity.SPSubscriptionPlan) error
	DeleteSPSubscriptionPlan(id string) (*entity.SPSubscriptionPlan, error)

	AddPlanChatLink(newPlanChatLink *entity.PlanChatLink) error
	FindPlanChatLink(planID string, chatID int64) (*entity.PlanChatLink, error)
	FindMultiplePlanChatLinks(identifier interface{}) []*entity.PlanChatLink
	DeletePlanChatLink(planID string, chatID int64) (*entity.PlanChatLink, error)
	DeleteMultiplePlanChatLinks(identifier interface{}) []*entity.PlanChatLink

	AddUserChatLink(newUserChatLink *entity.UserChatLink) error
	FindUserChatLink(userID, planID string, chatID int64) (*entity.UserChatLink, error)
	FindMultipleUserChatLinks(identifier interface{}) []*entity.UserChatLink
	DeleteUserChatLink(userID, planID string, chatID int64) (*entity.UserChatLink, error)
	DeleteMultipleUserChatLinks(identifier interface{}) []*entity.UserChatLink
}
