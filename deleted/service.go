package deleted

import "github.com/Benyam-S/onemembership/entity"

// IService is a method that defines all the service methods for managing deleted struct
type IService interface {
	AddUserToTrash(user *entity.User) (*entity.DeletedUser, error)
	AddServiceProviderToTrash(serviceProvider *entity.ServiceProvider) (*entity.DeletedServiceProvider, error)
	AddSubscriptionTranstactionsToTrash(userID, prefixedUserID string)
	AddSPSubscriptionTranstactionsToTrash(providerID, prefixedProviderID string)
	AddPayrollTranstactionsToTrash(providerID, prefixedProviderID string)

	FindDeletedUser(id string) (*entity.DeletedUser, error)
	SearchDeletedUsers(key, pagination string, extra ...string) []*entity.DeletedUser

	FindDeletedServiceProvider(id string) (*entity.DeletedServiceProvider, error)
	SearchDeletedServiceProviders(key, pagination string, extra ...string) []*entity.DeletedServiceProvider

	// FindDeletedSPSubTranstaction(id string) (*entity.DeletedSPSubscriptionTransaction, error)
	// SearchDeletedSPSubTranstactions(key, pagination string, extra ...string) []*entity.DeletedSPSubscriptionTransaction

	// FindDeletedSubTranstaction(id string) (*entity.DeletedUser, error)
	// SearchDeletedSubTranstactions(key, pagination string, extra ...string) []*entity.DeletedUser

	// FindDeletedPayrollTranstaction(id string) (*entity.DeletedUser, error)
	// SearchDeletedPayrollTranstactions(key, pagination string, extra ...string) []*entity.DeletedUser
}
