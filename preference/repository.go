package preference

import (
	"github.com/Benyam-S/onemembership/entity"
)

// IPreferenceRepository is an interface that defines all the repository methods of a client preference struct
type IPreferenceRepository interface {
	Create(newClientPreference *entity.ClientPreference) error
	Find(clientID string) (*entity.ClientPreference, error)
	Update(clientPreference *entity.ClientPreference) error
	UpdateValue(clientPreference *entity.ClientPreference, columnName string, columnValue interface{}) error
	Delete(clientID string) (*entity.ClientPreference, error)
}
