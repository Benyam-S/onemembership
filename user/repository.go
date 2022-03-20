package user

import (
	"time"

	"github.com/Benyam-S/onemembership/entity"
)

// IUserRepository is an interface that defines all the repository methods of a user struct
type IUserRepository interface {
	Create(newUser *entity.User) error
	Find(identifier string) (*entity.User, error)
	FindAll(pageNum int64) ([]*entity.User, int64)
	SearchWRegx(key string, pageNum int64, columns ...string) ([]*entity.User, int64)
	Search(key string, pageNum int64, columns ...string) ([]*entity.User, int64)
	All() []*entity.User
	Total() int64
	FromTo(start, end time.Time) int64
	Update(user *entity.User) error
	UpdateValue(user *entity.User, columnName string, columnValue interface{}) error
	Delete(id string) (*entity.User, error)
}

// IUserPasswordRepository is an interface that defines all the repository methods of a user's password struct
type IUserPasswordRepository interface {
	Create(newUserPassword *entity.UserPassword) error
	Find(userID string) (*entity.UserPassword, error)
	Update(userPassword *entity.UserPassword) error
	Delete(userID string) (*entity.UserPassword, error)
}
