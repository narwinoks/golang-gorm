package golang_orm

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"primaryKey;column:id;<-:create"`
	Name         Name      `gorm:"embedded"`
	Password     string    `gorm:"column:password"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Information  string    `gorm:"-"`
	Wallet       Wallet    `gorm:"foreignKey:user_id;references:id"`
	Address      []Address `gorm:"foreignKey:user_id;references:id"`
	LikeProducts []Product `gorm:"many2many:user_like_product;foreignKey:ID;joinForeignKey:user_id;references:ID;joinReferences:product_id"`
}

func (u User) TableName() string {
	return "users"
}

type Name struct {
	FirstName  string `gorm:"column:first_name"`
	LastName   string `gorm:"column:last_name"`
	MiddleName string `gorm:"column:middle_name"`
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.ID == "" {
		u.ID = "user-" + time.Now().Format("20060102150405")
	}
	return nil
}

//id is primary key
//default table name singular
//timestamp tracing
//field permission
//<- create only / update
//-> read only
// - ignore

type UserLog struct {
	ID        int    `gorm:"primaryKey,column:id,autoIncrement"`
	UserId    string `gorm:"column:user_id"`
	Action    string `gorm:"column:action"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli"`
}

func (receiver *UserLog) TableName() string {
	return "user_logs"
}
