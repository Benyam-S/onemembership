package subscriptionplan

import "github.com/Benyam-S/onemembership/entity"

// ISubscriptionPlanRepository is an interface that defines all the repository methods of a subscription plan struct
type ISubscriptionPlanRepository interface {
	Create(newSubscriptionPlan *entity.SubscriptionPlan) error
	Find(id string) (*entity.SubscriptionPlan, error)
	FindMultiple(projectID string) []*entity.SubscriptionPlan
	Update(subscriptionPlan *entity.SubscriptionPlan) error
	Delete(id string) (*entity.SubscriptionPlan, error)
	DeleteMultiple(projectID string) []*entity.SubscriptionPlan
}

// IPlanChatLinkRepository is an interface that defines all the repository methods of a subscription plan to chat link (PlanChatLink) struct
type IPlanChatLinkRepository interface {
	Create(newPlanChatLink *entity.PlanChatLink) error
	Find(planID string, chatID int64) (*entity.PlanChatLink, error)
	FindMultiple(identifier interface{}) []*entity.PlanChatLink
	Delete(planID string, chatID int64) (*entity.PlanChatLink, error)
	DeleteMultiple(identifier interface{}) []*entity.PlanChatLink
}

// ISPSubscriptionPlanRepository is an interface that defines all the repository methods of a service provider subscription plan struct
type ISPSubscriptionPlanRepository interface {
	Create(newSubscriptionPlan *entity.SPSubscriptionPlan) error
	Find(id string) (*entity.SPSubscriptionPlan, error)
	Update(subscriptionPlan *entity.SPSubscriptionPlan) error
	Delete(id string) (*entity.SPSubscriptionPlan, error)
}

// IUserChatLinkRepository is an interface that defines all the repository methods of a user to chat link (UserChatLink) struct
type IUserChatLinkRepository interface {
	Create(newUserChatLink *entity.UserChatLink) error
	Find(userID, planID string, chatID int64) (*entity.UserChatLink, error)
	FindMultiple(identifier interface{}) []*entity.UserChatLink
	Delete(userID, planID string, chatID int64) (*entity.UserChatLink, error)
	DeleteMultiple(identifier interface{}) []*entity.UserChatLink
}
