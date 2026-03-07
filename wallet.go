package golang_orm

import "time"

type Wallet struct {
	ID        string    `gorm:"primary_key;column:id"`
	UserID    string    `gorm:"column:user_id"`
	Balance   int       `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdateAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (w Wallet) TableName() string {
	return "wallets"
}
