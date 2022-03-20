package preference

import (
	"github.com/Benyam-S/onemembership/entity"
)

// IService is an interface that defines all the service methods of a client preference struct
type IService interface {
	AddClientPreference(newClientPreference *entity.ClientPreference) error
	ValidateClientPreference(clientPreference *entity.ClientPreference) entity.ErrMap
	FindClientPreference(clientID string) (*entity.ClientPreference, error)
	UpdateClientPreference(clientPreference *entity.ClientPreference) error
	UpdateClientPreferenceSingleValue(userID, columnName string, columnValue interface{}) error
	DeleteClientPreference(clientID string) (*entity.ClientPreference, error)
}
