package golang_orm

import "time"

type Address struct {
	ID        int64     `gorm:"primaryKey;column:id;autoIncrement"`
	UserId    string    `gorm:"column:user_id"`
	Address   string    `gorm:"column:address"`
	CreateAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User      User      `gorm:"foreign_key:user_id;references:id"`
}

func (a Address) TableName() string {
	return "addresses"
}
