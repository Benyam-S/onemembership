package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onemembership/deleted"
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/tools"
	"github.com/jinzhu/gorm"
)

// DeletedUserRepository is a type that defines a repository for deleted user
type DeletedUserRepository struct {
	conn *gorm.DB
}

// NewDeletedUserRepository is a function that returns a new deleted user repository
func NewDeletedUserRepository(connection *gorm.DB) deleted.IDeletedUserRepository {
	return &DeletedUserRepository{conn: connection}
}

// Create is a method that adds a deleted user to the database
func (repo *DeletedUserRepository) Create(deletedUser *entity.DeletedUser) error {

	totalNumOfDeletedUsers := tools.CountMembers("deleted_users", repo.conn)
	deletedUser.ID = fmt.Sprintf("DUR_%s%d", tools.RandomStringGN(7), totalNumOfDeletedUsers+1)

	for !tools.IsUnique("id", deletedUser.ID, "deleted_users", repo.conn) {
		totalNumOfDeletedUsers++
		deletedUser.ID = fmt.Sprintf("DUR_%s%d", tools.RandomStringGN(7), totalNumOfDeletedUsers+1)
	}

	// Adding prefix to UserID so to uniquely identify it from other deleted users
	deletedUser.UserID = deletedUser.ID + "_" + deletedUser.UserID

	err := repo.conn.Create(deletedUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain deleted user from the database using an id,
// also Find() uses id as a key for selection
func (repo *DeletedUserRepository) Find(id string) (*entity.DeletedUser, error) {
	deletedUser := new(entity.DeletedUser)
	err := repo.conn.Model(deletedUser).Where("id = ? ", id).
		First(deletedUser).Error

	if err != nil {
		return nil, err
	}
	return deletedUser, nil
}

// Search is a method that search and returns a set of deleted users from the database using an id.
func (repo *DeletedUserRepository) Search(key string, pageNum int64, columns ...string) []*entity.DeletedUser {

	var deletedUsers []*entity.DeletedUser
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		// modifying the key so that it can match the database phone number values
		if column == "phone_number" {
			splitKey := strings.Split(key, "")
			if splitKey[0] == "0" {
				modifiedKey := "+251" + strings.Join(splitKey[1:], "")
				whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
				sqlValues = append(sqlValues, modifiedKey)
				continue
			}
		}
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_users WHERE ("+strings.Join(whereStmt, "||")+") ORDER BY first_name ASC LIMIT ?, 30",
		sqlValues...).Scan(&deletedUsers)

	return deletedUsers
}

// SearchWRegx is a method that searchs and returns set of deleted users limited to the key id and page number using regular expersions
func (repo *DeletedUserRepository) SearchWRegx(key string, pageNum int64, columns ...string) []*entity.DeletedUser {
	var deletedUsers []*entity.DeletedUser
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s regexp ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM deleted_users WHERE "+strings.Join(whereStmt, "||")+" ORDER BY first_name ASC LIMIT ?, 30",
		sqlValues...).Scan(&deletedUsers)

	return deletedUsers
}

// All is a method that returns all the deleted users from the database limited with the pageNum
func (repo *DeletedUserRepository) All(pageNum int64) []*entity.DeletedUser {

	var deletedUsers []*entity.DeletedUser
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM deleted_users ORDER BY first_name ASC LIMIT ?, 30", limit).Scan(&deletedUsers)
	return deletedUsers
}

// Update is a method that updates a certain deleted user value in the database
func (repo *DeletedUserRepository) Update(deletedUser *entity.DeletedUser) error {

	prevUser := new(entity.DeletedUser)
	err := repo.conn.Model(prevUser).Where("id = ?", deletedUser.ID).First(prevUser).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(deletedUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that removes a certain deleted user from the database using an id.
// In Delete() id is only used as a key
func (repo *DeletedUserRepository) Delete(id string) (*entity.DeletedUser, error) {
	deletedUser := new(entity.DeletedUser)
	err := repo.conn.Model(deletedUser).Where("id = ?", id).First(deletedUser).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(deletedUser)
	return deletedUser, nil
}
