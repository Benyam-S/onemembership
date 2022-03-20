package repository

import (
	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/transaction"
	"github.com/jinzhu/gorm"
)

// PaymentGatewayRepository is a type that defines a payment gateway repository type
type PaymentGatewayRepository struct {
	conn *gorm.DB
}

// NewPaymentGatewayRepository is a function that creates a new payment gateway repository type
func NewPaymentGatewayRepository(connection *gorm.DB) transaction.IPaymentGatewayRepository {
	return &PaymentGatewayRepository{conn: connection}
}

// Create is a method that adds a new payment gateway to the database
func (repo *PaymentGatewayRepository) Create(newPaymentGateway *entity.PaymentGateway) error {
	err := repo.conn.Create(newPaymentGateway).Error
	if err != nil {
		return err
	}

	return nil
}

// Find is a method that finds a certain payment gateway from the database using an gateway id,
// also Find() uses only id as a key for selection
func (repo *PaymentGatewayRepository) Find(id int64) (*entity.PaymentGateway, error) {

	paymentGateway := new(entity.PaymentGateway)
	err := repo.conn.Model(paymentGateway).Where("id = ?", id).First(paymentGateway).Error

	if err != nil {
		return nil, err
	}
	return paymentGateway, nil
}

// All is a method that returns all the payment gateway found in the database
func (repo *PaymentGatewayRepository) All() []*entity.PaymentGateway {

	var paymentGateways []*entity.PaymentGateway

	repo.conn.Model(entity.PaymentGateway{}).Find(&paymentGateways).Order("created_at ASC")

	return paymentGateways
}

// Update is a method that updates a certain payment gateway entries in the database
func (repo *PaymentGatewayRepository) Update(paymentGateway *entity.PaymentGateway) error {

	prevPaymentGateway := new(entity.PaymentGateway)
	err := repo.conn.Model(prevPaymentGateway).Where("id = ?", paymentGateway.ID).
		First(prevPaymentGateway).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	paymentGateway.CreatedAt = prevPaymentGateway.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(paymentGateway).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain payment gateway from the database using an gateway id.
// In Delete() id is only used as an key
func (repo *PaymentGatewayRepository) Delete(id int64) (*entity.PaymentGateway, error) {
	paymentGateway := new(entity.PaymentGateway)
	err := repo.conn.Model(paymentGateway).Where("id = ?", id).First(paymentGateway).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(paymentGateway)
	return paymentGateway, nil
}
