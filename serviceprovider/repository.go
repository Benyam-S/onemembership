package serviceprovider

import (
	"time"

	"github.com/Benyam-S/onemembership/entity"
)

// IServiceProviderRepository is an interface that defines all the repository methods of a service provider struct
type IServiceProviderRepository interface {
	Create(newServiceProvider *entity.ServiceProvider) error
	Find(identifier string) (*entity.ServiceProvider, error)
	FindAll(pageNum int64) ([]*entity.ServiceProvider, int64)
	SearchWRegx(key string, pageNum int64, columns ...string) ([]*entity.ServiceProvider, int64)
	Search(key string, pageNum int64, columns ...string) ([]*entity.ServiceProvider, int64)
	All() []*entity.ServiceProvider
	Total() int64
	FromTo(start, end time.Time) int64
	Update(serviceprovider *entity.ServiceProvider) error
	UpdateValue(serviceprovider *entity.ServiceProvider, columnName string, columnValue interface{}) error
	Delete(id string) (*entity.ServiceProvider, error)
}

// ISPPasswordRepository is an interface that defines all the repository methods of a service provider's password struct
type ISPPasswordRepository interface {
	Create(newSPPassword *entity.SPPassword) error
	Find(providerID string) (*entity.SPPassword, error)
	Update(spPassword *entity.SPPassword) error
	Delete(providerID string) (*entity.SPPassword, error)
}

// ISPWalletRepository is an interface that defines all the repository methods of a service provider wallet struct
type ISPWalletRepository interface {
	Create(newSPWallet *entity.SPWallet) error
	Find(providerID string) (*entity.SPWallet, error)
	Update(spWallet *entity.SPWallet) error
	UpdateValue(spWallet *entity.SPWallet, columnName string, columnValue interface{}) error
	Delete(providerID string) (*entity.SPWallet, error)
}
