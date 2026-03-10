package golang_orm

import "time"

type Wallet struct {
	ID        string    `gorm:"primaryKey;column:id"`
	UserID    string    `gorm:"column:user_id"`
	Balance   int       `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdateAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
	User      *User     `gorm:"foreign_key:UserID;references:id"`
}

func (w Wallet) TableName() string {
	return "wallets"
}
