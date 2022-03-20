package deleted

import "github.com/Benyam-S/onemembership/entity"

// IDeletedUserRepository is an interface that defines all the repository methods for managing deleted users
type IDeletedUserRepository interface {
	Create(deletedUser *entity.DeletedUser) error
	Find(id string) (*entity.DeletedUser, error)
	Search(key string, pageNum int64, columns ...string) []*entity.DeletedUser
	SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedUser
	All(pageNum int64) []*entity.DeletedUser
	Update(deletedUser *entity.DeletedUser) error
	Delete(id string) (*entity.DeletedUser, error)
}

// IDeletedServiceProviderRepository is an interface that defines all the repository methods for managing deleted service providers
type IDeletedServiceProviderRepository interface {
	Create(deletedServiceProvider *entity.DeletedServiceProvider) error
	Find(id string) (*entity.DeletedServiceProvider, error)
	Search(key string, pageNum int64, columns ...string) []*entity.DeletedServiceProvider
	SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedServiceProvider
	All(pageNum int64) []*entity.DeletedServiceProvider
	Update(deletedServiceProvider *entity.DeletedServiceProvider) error
	Delete(id string) (*entity.DeletedServiceProvider, error)
}

// IDeletedSubscriptionTransactionRepository is an interface that defines all the repository methods for managing deleted subscription transaction
type IDeletedSubscriptionTransactionRepository interface {
	Create(userID, prefixedUserID string)
	Find(id string) (*entity.DeletedSubscriptionTransaction, error)
	Search(key string, pageNum int64, columns ...string) []*entity.DeletedSubscriptionTransaction
	SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedSubscriptionTransaction
	All(pageNum int64) []*entity.DeletedSubscriptionTransaction
	Update(deletedSubscriptonTransaction *entity.DeletedSubscriptionTransaction) error
	Delete(id string) (*entity.DeletedSubscriptionTransaction, error)
}

// IDeletedSPSubscriptionTransactionRepository is an interface that defines all the repository methods for managing deleted service provider subscription transaction
type IDeletedSPSubscriptionTransactionRepository interface {
	Create(providerID, prefixedProviderID string)
	Find(id string) (*entity.DeletedSPSubscriptionTransaction, error)
	Search(key string, pageNum int64, columns ...string) []*entity.DeletedSPSubscriptionTransaction
	SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedSPSubscriptionTransaction
	All(pageNum int64) []*entity.DeletedSPSubscriptionTransaction
	Update(deletedSubscriptonTransaction *entity.DeletedSPSubscriptionTransaction) error
	Delete(id string) (*entity.DeletedSPSubscriptionTransaction, error)
}

// IDeletedSPPayrollTransactionRepository is an interface that defines all the repository methods for managing deleted service provider payroll transaction
type IDeletedSPPayrollTransactionRepository interface {
	Create(providerID, prefixedProviderID string)
	Find(id string) (*entity.DeletedSPPayrollTransaction, error)
	Search(key string, pageNum int64, columns ...string) []*entity.DeletedSPPayrollTransaction
	SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedSPPayrollTransaction
	All(pageNum int64) []*entity.DeletedSPPayrollTransaction
	Update(deletedPayrollTransaction *entity.DeletedSPPayrollTransaction) error
	Delete(id string) (*entity.DeletedSPPayrollTransaction, error)
}
