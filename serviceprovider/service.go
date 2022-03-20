package serviceprovider

import "github.com/Benyam-S/onemembership/entity"

// IService is an interface that defines all the service methods of a service provider struct
type IService interface {
	AddServiceProvider(newServiceProvider *entity.ServiceProvider) error
	ValidateProviderProfile(serviceProvider *entity.ServiceProvider) entity.ErrMap
	FindServiceProvider(identifier string) (*entity.ServiceProvider, error)
	AllServiceProviders() []*entity.ServiceProvider
	AllServiceProvidersWithPagination(pageNum int64) ([]*entity.ServiceProvider, int64)
	SearchServiceProviders(key string, pageNum int64, extra ...string) ([]*entity.ServiceProvider, int64)
	TotalServiceProviders() int64
	UpdateServiceProvider(serviceProvider *entity.ServiceProvider) error
	UpdateProviderSingleValue(providerID, columnName string, columnValue interface{}) error
	DeleteServiceProvider(providerID string) (*entity.ServiceProvider, error)

	AddSPPassword(newSPPassword *entity.SPPassword) error
	VerifySPPassword(spPassword *entity.SPPassword, verifyPassword string) error
	FindSPPassword(providerID string) (*entity.SPPassword, error)
	UpdateSPPassword(spPassword *entity.SPPassword) error
	DeleteSPPassword(providerID string) (*entity.SPPassword, error)

	AddSPWallet(newServiceProviderWallet *entity.SPWallet) error
	ValidateSPWallet(serviceProviderWallet *entity.SPWallet) entity.ErrMap
	FindSPWallet(identifier string) (*entity.SPWallet, error)
	UpdateSPWallet(serviceProviderWallet *entity.SPWallet) error
	UpdateSPWalletSingleValue(providerID, columnName string, columnValue interface{}) error
	DeleteSPWallet(providerID string) (*entity.SPWallet, error)
}
