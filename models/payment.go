package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Payments is our payment model
type Payments struct {
	gorm.Model
	Number         int64   `gorm:"type:integer"`
	DepositAddress string  `gorm:"type:varchar(255)"`
	TxHash         string  `gorm:"type:varchar(255)"`
	USDValue       float64 `gorm:"type:float"` // USDValue is also a "Credit" value, since 1 USD -> 1 Credit
	Blockchain     string  `gorm:"type:varchar(255)"`
	Type           string  `gorm:"type:varchar(255)"` // ETH, RTC, XMR, BTC, LTC
	UserName       string  `gorm:"type:varchar(255)"`
	Confirmed      bool    `gorm:"type:varchar(255)"`
}

// PaymentManager is used to interact with payment information in our database
type PaymentManager struct {
	DB *gorm.DB
}

// NewPaymentManager is used to generate our payment manager helper
func NewPaymentManager(db *gorm.DB) *PaymentManager {
	return &PaymentManager{DB: db}
}

// GetLatestPaymentNumber is used to get the latest payment number for a user
func (pm *PaymentManager) GetLatestPaymentNumber(username string) (int64, error) {
	p := Payments{}
	check := pm.DB.Where("user_name = ?", username).Order("number desc").First(&p)
	if check.Error != nil && check.Error != gorm.ErrRecordNotFound {
		return 0, check.Error
	}
	if check.Error == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return p.Number, nil
}

// NewPayment is used to create a payment in our database
func (pm *PaymentManager) NewPayment(number int64, depositAddress string, txHash string, usdValue float64, blockchain string, paymentType string, username string) (*Payments, error) {
	p := Payments{}
	check := pm.DB.Where("tx_hash = ?", txHash).First(&p)
	if check.Error == nil {
		return nil, errors.New("payment with tx hash already exists")
	} else if check.Error != nil && check.Error != gorm.ErrRecordNotFound {
		return nil, check.Error
	}
	check = pm.DB.Where("number = ?", number).First(&p)
	if check.Error == nil {
		return nil, errors.New("payment with number already exists")
	} else if check.Error != nil && check.Error != gorm.ErrRecordNotFound {
		return nil, check.Error
	}
	p = Payments{
		DepositAddress: depositAddress,
		TxHash:         txHash,
		USDValue:       usdValue,
		Blockchain:     blockchain,
		Type:           paymentType,
		UserName:       username,
		Confirmed:      false,
	}

	if check := pm.DB.Create(&p); check.Error != nil {
		return nil, check.Error
	}

	return &p, nil
}

// ConfirmPayment is used to mark a payment as confirmed
func (pm *PaymentManager) ConfirmPayment(txHash string) (*Payments, error) {
	p := Payments{}
	if check := pm.DB.Where("tx_hash = ?", txHash).First(&p); check.Error != nil {
		return nil, check.Error
	}
	p.Confirmed = true
	if check := pm.DB.Model(&p).Update("confirmed", p.Confirmed); check.Error != nil {
		return nil, check.Error
	}
	return &p, nil
}
