package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Drop is a users registered airdrop
type Drop struct {
	gorm.Model
	UserName   string `gorm:"type:varchar(255)" json:"user_name"`
	DropID     string `gorm:"type:varchar(255);unique" json:"drop"`
	EthAddress string `gorm:"Type:varchar(255);unique" json:"eth_address"`
}

// DropManager is used to manipualte airdrop models
type DropManager struct {
	DB *gorm.DB
}

// NewDropManager is used to create our airdrop model manager
func NewDropManager(db *gorm.DB) *DropManager {
	return &DropManager{db}
}

// RegisterAirDrop is used to register an airdrop
func (dm *DropManager) RegisterAirDrop(dropID, ethAddress, username string) (*Drop, error) {
	d, err := dm.FindByDropID(dropID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("unexpected error occured")
	}
	if err == nil {
		return nil, errors.New("airdrop with id already exists")
	}

	d, err = dm.FindByEthAddress(ethAddress)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("unexpected error occured")
	}
	if err == nil {
		return nil, errors.New("airdrop with eth address already exists")
	}

	d = &Drop{
		UserName:   username,
		DropID:     dropID,
		EthAddress: ethAddress,
	}
	if err = dm.DB.Create(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}

// FindByDropID is used to find an airdrop by its airdrop id
func (dm *DropManager) FindByDropID(dropID string) (*Drop, error) {
	d := Drop{}
	if err := dm.DB.Where("drop_id = ?", dropID).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

// FindByEthAddress is used to find an airdrop by an eth address
func (dm *DropManager) FindByEthAddress(ethAddress string) (*Drop, error) {
	d := Drop{}
	if err := dm.DB.Where("eth_address = ?", ethAddress).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
