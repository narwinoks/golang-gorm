package golang_orm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Sample struct {
	Id   string
	Name string
}

func OpenConnection() *gorm.DB {
	driver := mysql.Open("root:root@tcp(127.0.0.1:3306)/golang_orm?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(driver, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return db
}

var db = OpenConnection()

func TestOpenConnection(t *testing.T) {
	assert.NotNil(t, db)
}

func TestExecuteSQL(t *testing.T) {
	err := db.Exec("INSERT INTO sample (id,name) VALUES(?,?)", 1, "EKO").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample (id,name) VALUES(?,?)", 2, "RULI").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample (id,name) VALUES(?,?)", 3, "JOKO").Error
	assert.Nil(t, err)

}

func TestRawSQL(t *testing.T) {
	var sample Sample
	err := db.Raw("SELECT id,name FROM sample where id = ?", "1").Scan(&sample).Error
	assert.Nil(t, err)
	assert.Equal(t, "EKO", sample.Name)

	var samples []Sample
	err = db.Raw("SELECT id,name FROM sample order by id ASC").Scan(&samples).Error
	assert.Nil(t, err)
	assert.Equal(t, 3, len(samples))
	//
}

func TestSqlRow(t *testing.T) {
	rows, err := db.Raw("SELECT id,name FROM sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		var id string
		var name string
		err := rows.Scan(&id, &name)
		assert.Nil(t, err)

		samples = append(samples, Sample{Id: id, Name: name})

	}
	assert.Equal(t, 3, len(samples))

}
func TestScanRow(t *testing.T) {
	rows, err := db.Raw("SELECT id,name FROM sample").Rows()
	assert.Nil(t, err)
	defer rows.Close()

	var samples []Sample
	for rows.Next() {
		db.ScanRows(rows, &samples)

	}
	assert.Equal(t, 3, len(samples))

}

func TestCreateUser(t *testing.T) {
	user := User{
		ID: "1",
		Name: Name{
			FirstName:  "eko",
			LastName:   "eko",
			MiddleName: "eko",
		},
		Password:    "rahasia",
		Information: "information",
	}
	tx := db.Create(&user)
	assert.Nil(t, tx.Error)
	assert.Equal(t, int64(1), tx.RowsAffected)
}
func TestBatchInsert(t *testing.T) {
	var users []User
	for i := 2; i < 12; i++ {
		users = append(users, User{
			ID: strconv.Itoa(i),
			Name: Name{
				FirstName:  "first name" + strconv.Itoa(i),
				LastName:   "middle name" + strconv.Itoa(i),
				MiddleName: "last name" + strconv.Itoa(i),
			},
			Password: "Rahasia",
		})
	}
	result := db.Create(&users)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(10), result.RowsAffected)
}

func TestTransactionSuccess(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{ID: "13", Password: "rahasia", Name: Name{FirstName: "User 13", MiddleName: "User 13", LastName: "User 13"}}).Error
		if err != nil {
			return err
		}
		err = tx.Create(&User{ID: "14", Password: "rahasia", Name: Name{FirstName: "User 14", MiddleName: "User 14", LastName: "User 14"}}).Error
		if err != nil {
			return err
		}
		err = tx.Create(&User{ID: "15", Password: "rahasia", Name: Name{FirstName: "User 15", MiddleName: "User 15", LastName: "User 15"}}).Error
		if err != nil {
			return err
		}
		return nil
	})
	assert.Nil(t, err)
}

func TestTransactionFailed(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&User{ID: "16", Password: "rahasia", Name: Name{FirstName: "User 13", MiddleName: "User 13", LastName: "User 13"}}).Error
		if err != nil {
			return err
		}
		err = tx.Create(&User{ID: "14", Password: "rahasia", Name: Name{FirstName: "User 14", MiddleName: "User 14", LastName: "User 14"}}).Error
		if err != nil {
			return err
		}
		return nil
	})
	assert.NotNil(t, err)
}

func TestTransactionManualSuccess(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()
	err := tx.Create(&User{ID: "16", Password: "rahasia", Name: Name{FirstName: "User 16", MiddleName: "User 16", LastName: "User 16"}}).Error
	assert.Nil(t, err)
	err = tx.Create(&User{ID: "17", Password: "rahasia", Name: Name{FirstName: "User 17", MiddleName: "User 17", LastName: "User 17"}}).Error
	assert.Nil(t, err)
	if err == nil {
		tx.Commit()
	}
}
func TestTransactionManualFailed(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()
	err := tx.Create(&User{ID: "17", Password: "rahasia", Name: Name{FirstName: "User 17", MiddleName: "User 17", LastName: "User 17"}}).Error
	assert.NotNil(t, err)
	err = tx.Create(&User{ID: "17", Password: "rahasia", Name: Name{FirstName: "User 17", MiddleName: "User 17", LastName: "User 17"}}).Error
	assert.NotNil(t, err)
	if err == nil {
		tx.Commit()
	}
}

// single object orm
func TestQuerySingleObject(t *testing.T) {
	user := User{}
	tx := db.First(&user).Error
	assert.Nil(t, tx)
	assert.Equal(t, "10", user.ID)

	user = User{}
	tx = db.Last(&user).Error
	assert.Nil(t, tx)
	assert.Equal(t, "9", user.ID)
}

func TestQuerySingleObjectInlineCondition(t *testing.T) {
	user := User{}
	err := db.First(&user, "id = ?", "5").Error
	assert.Nil(t, err)
	assert.Equal(t, "5", user.ID)
}

//query all object

func TestQueryAllObject(t *testing.T) {
	var users []User
	err := db.Find(&users, "id IN ?", []string{"1", "2", "3", "4", "5"}).Error
	assert.Nil(t, err)
	assert.Equal(t, 4, len(users))
}

func TestQueryCondition(t *testing.T) {
	var users []User
	err := db.Where("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 6, len(users))
}

func TestOrOperator(t *testing.T) {
	var users []User
	err := db.Where("first_name like ?", "%User%").Or("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
}

func TestNotOperator(t *testing.T) {
	var users []User
	err := db.Not("first_name like ?", "%User%").Where("password = ?", "rahasia").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 11, len(users))
}

// select not all field
func TestSelectField(t *testing.T) {
	var users []User
	err := db.Select("id", "first_name").Find(&users).Error
	assert.Nil(t, err)
	for _, user := range users {
		assert.NotNil(t, user.ID)
		assert.NotEqual(t, "", user.Name.FirstName)
	}
	assert.Equal(t, 17, len(users))
}

//where dinamis condition

func TestStructCondition(t *testing.T) {
	userCondition := User{
		Name: Name{
			FirstName: "User 13",
			LastName:  "", // tidak bisa karna dianggap default value
		},
		Password: "rahasia",
	}
	var users []User
	err := db.Where(userCondition).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))

}

func TestMapCondition(t *testing.T) {
	mapCondtion := map[string]interface{}{
		"middle_name": "",
	}
	var users []User
	err := db.Where(mapCondtion).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(users))
}

// order limit offset
func TestOrderLimitOffset(t *testing.T) {
	var users []User
	err := db.Order("id asc,first_name asc").Limit(5).Offset(5).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))
}

type UserResponse struct {
	ID        string
	FirstName string
	LastName  string
}

func TestQueryNonModel(t *testing.T) {
	var users []UserResponse
	err := db.Model(&User{}).Select("id", "first_name", "last_name").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 17, len(users))
	fmt.Print(users)
}
func TestUpdate(t *testing.T) {
	user := User{}
	err := db.Take(&user, "id = ?", "2").Error
	assert.Nil(t, err)
	user.Name.FirstName = "Budi"
	user.Name.MiddleName = ""
	user.Name.LastName = "Nugraha"
	user.Password = "rahasia123"
	err = db.Save(&user).Error
	assert.Nil(t, err)
}

func TestSelectedColumn(t *testing.T) {
	err := db.Model(&User{}).Where("id = ?", "2").Updates(map[string]interface{}{
		"middle_name": "",
		"last_name":   "Morro",
	}).Error
	assert.Nil(t, err)
	err = db.Model(&User{}).Where("id = ?", "2").Update("password", "ubahlagi").Error
	assert.Nil(t, err)
	err = db.Model(&User{}).Where("id  = ? ", "2").Updates(User{
		Name: Name{
			FirstName:  "EKO UPDATE",
			MiddleName: "EKO UPDATE",
			LastName:   "TEST UPDATE",
		},
	}).Error
	assert.Nil(t, err)
}

func TestAutoIncrement(t *testing.T) {
	for i := 0; i <= 10; i++ {
		userLog := UserLog{
			UserId: "1",
			Action: "This ACTION",
		}
		err := db.Create(&userLog).Error
		assert.Nil(t, err)
		assert.NotEqual(t, 0, userLog.ID)
		fmt.Print(userLog.ID)
	}
}
func TestUpdateOrCreate(t *testing.T) {
	//auto increment
	userLog := UserLog{
		UserId: "2",
		Action: "This Is Action",
	}
	err := db.Save(&userLog).Error //user log
	assert.Nil(t, err)
	userLog.UserId = "2"
	err = db.Save(&userLog).Error // update
	assert.Nil(t, err)
}

func TestSaveOrUpdateNonAutoIncrement(t *testing.T) {
	user := User{
		ID: "99",
		Name: Name{
			FirstName: "user 99",
		},
	}
	err := db.Save(&user).Error //insert
	assert.Nil(t, err)

	user.Name.FirstName = "user 99 update"
	err = db.Save(&user).Error
	assert.Nil(t, err)
}

func TestConflict(t *testing.T) {
	user := User{
		ID: "88",
		Name: Name{
			FirstName: "user 88",
		},
	}
	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error
	assert.Nil(t, err)
}

//delete

func TestDeleteCondition(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "88").Error
	assert.Nil(t, err)
	err = db.Delete(&user).Error
	assert.Nil(t, err)

	err = db.Delete(&user, "id = ? ", "99").Delete(&User{}).Error
	assert.Nil(t, err)

	err = db.Where("id = ?", "77").Delete(&User{}).Error
	assert.Nil(t, err)
}

// soft delete

func TestSoftDelete(t *testing.T) {
	todo := Todo{
		UserId:      "1",
		Title:       "This Is Title",
		Description: "This Is description",
	}
	err := db.Create(&todo).Error
	assert.Nil(t, err)

	err = db.Delete(&todo).Error
	assert.Nil(t, err)

	assert.NotNil(t, todo.DeletedAt)

	var todos []Todo
	err = db.Find(&todos).Error
	assert.Nil(t, err)

	assert.Equal(t, 0, len(todos))
}

func TestUnscoped(t *testing.T) {
	var todo Todo
	err := db.Unscoped().First(&todo, "id = ?", "7").Error
	assert.Nil(t, err)

	err = db.Unscoped().Delete(&todo).Error
	assert.Nil(t, err)
	var todos []Todo
	err = db.Unscoped().Find(&todos).Error
	assert.Nil(t, err)

}

func TestLock(t *testing.T) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, "id = ?", "1").Error
		if err != nil {
			return err
		}
		user.Name.FirstName = "Joko"
		user.Name.LastName = "Moro"
		err = tx.Save(&user).Error
		return err
	})
	assert.Nil(t, err)
}

func TestCreateWallet(t *testing.T) {
	wallet := Wallet{
		UserID:  "1",
		ID:      "2",
		Balance: 100000,
	}
	err := db.Create(&wallet).Error
	assert.Nil(t, err)
}

// preload
func TestRetrieveRelation(t *testing.T) {
	var user User
	err := db.Model(&User{}).Preload("Wallet").Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "2", user.Wallet.ID)
}

func TestRetrieveJoins(t *testing.T) {
	var user User
	err := db.Model(&User{}).Joins("Wallet").Take(&user, "users.id = ?", "1").Error
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "2", user.Wallet.ID)
}

func TestAutoCreateOrUpdate(t *testing.T) {
	user := User{
		ID:       "20",
		Password: "Rahasia",
		Name: Name{
			FirstName: "First Name",
		},
		Wallet: Wallet{
			ID:      "20",
			UserID:  "20",
			Balance: 1000000,
		},
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)
}
func TestSkipAutoCreateOrUpdate(t *testing.T) {
	user := User{
		ID:       "21",
		Password: "RAHASIA",
		Name: Name{
			FirstName: "First Name",
		},
		Wallet: Wallet{
			ID:      "21",
			UserID:  "21",
			Balance: 100000,
		},
	}
	err := db.Omit(clause.Associations).Create(&user).Error
	assert.Nil(t, err)
}

func TestUserAndAddress(t *testing.T) {
	user := User{
		ID:       "51",
		Password: "Rahasia",
		Name: Name{
			FirstName: "user is 51",
		},
		Wallet: Wallet{
			UserID:  "51",
			ID:      "51",
			Balance: 1000000,
		},
		Address: []Address{
			{
				UserId:  "51",
				Address: "Bandung 1",
			},
			{
				UserId:  "51",
				Address: "Bandung 2",
			},
		},
	}
	err := db.Create(&user).Error
	assert.Nil(t, err)
}
func TestPreloadJoinOneToMany(t *testing.T) {
	var users []User
	err := db.Model(&User{}).Preload("Address").Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
}
func TestTakePreloadJoinOneToMany(t *testing.T) {
	var user User
	err := db.Model(&User{}).Preload("Address").Joins("Wallet").Take(&user, "users.id = ?", "50").Error
	assert.Nil(t, err)
}

func TestBelongTO(t *testing.T) {
	fmt.Print("Preload")
	var address []Address
	err := db.Model(&Address{}).Preload("User").Find(&address).Error
	assert.Nil(t, err)
	fmt.Print("Joins")
	address = []Address{}
	err = db.Model(&Address{}).Joins("User").Find(&address).Error
	assert.Nil(t, err)
}

func TestBelongsToWallet(t *testing.T) {
	fmt.Print("Preload")
	var wallet []Wallet
	err := db.Model(&Wallet{}).Preload("User").Find(&wallet).Error
	assert.Nil(t, err)
	fmt.Print("Preload")
	wallet = []Wallet{}
	err = db.Model(&Wallet{}).Joins("User").Find(&Wallet{}).Error
	assert.Nil(t, err)
}

func TestCreateManyToMany(t *testing.T) {
	product := Product{
		Name:  "PRODUCT 001",
		Price: 100000,
		ID:    "P001",
	}
	err := db.Create(&product).Error
	assert.Nil(t, err)
	err = db.Table("user_like_product").Create(map[string]interface{}{
		"user_id":    "1",
		"product_id": "P001",
	}).Error
	err = db.Table("user_like_product").Create(map[string]interface{}{
		"user_id":    "2",
		"product_id": "P001",
	}).Error
	assert.Nil(t, err)
}

func TestPreloadManyToMany(t *testing.T) {
	var product Product
	err := db.Preload("LikedByUsers").Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(product.LikedByUsers))
}
func TestPreloadManyToManyUser(t *testing.T) {
	var user User
	err := db.Preload("LikeProducts").Take(&user, "id  = ?", 2).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, len(user.LikeProducts))
}

// association mode
func TestAssociationMode(t *testing.T) {
	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	var users []User
	err = db.Model(&product).Where("users.first_name = ?", "Budi").Association("LikedByUsers").Find(&users)
	assert.Equal(t, 1, len(users))
}

func TestAssociationAdd(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "3").Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	err = db.Model(&product).Association("LikedByUsers").Append(&user)
	assert.Nil(t, err)
}

func TestAssociationReplace(t *testing.T) {
	db.Transaction(func(tx *gorm.DB) error {
		var user User
		err := tx.Take(&user, "id =?", "1").Error
		assert.Nil(t, err)

		wallet := Wallet{ID: "01", UserID: user.ID, Balance: 2000000}
		err = tx.Model(&user).Association("Wallet").Replace(&wallet)
		assert.Nil(t, err)
		return nil
	})
}

func TestAssociationDelete(t *testing.T) {
	var user User
	err := db.Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)

	var product Product
	err = db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)

	err = db.Model(&product).Association("LikedByUsers").Delete(&user)
	assert.Nil(t, err)
}

func TestAssociationClear(t *testing.T) {
	var product Product
	err := db.Take(&product, "id = ?", "P001").Error
	assert.Nil(t, err)
	err = db.Model(&product).Association("LikedByUsers").Clear()
	assert.Nil(t, err)
}

// preloading condition
func TestPreloadingWithCondition(t *testing.T) {
	var user User
	err := db.Preload("Wallet", "balance > ?", 1000000).Take(&user, "id = ?", "1").Error
	assert.Nil(t, err)
}

// nested preloading
func TestNestedPreloading(t *testing.T) {
	var wallet Wallet
	err := db.Preload("User.Address").Take(&wallet, "id = ?", "2").Error
	assert.Nil(t, err)
	fmt.Print(wallet)
	fmt.Print(wallet.User)
	fmt.Print(wallet.User.Address)
}

// preload all
func TestPreloadingAll(t *testing.T) {
	var user User
	err := db.Preload(clause.Associations).Take(&user, "id  = ?", "2").Error
	assert.Nil(t, err)
}

// joins
func TestJoinQuery(t *testing.T) {
	var users []User
	err := db.Joins("join wallets on wallets.user_id = users.id").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users)) //inner join

	users = []User{}
	err = db.Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 22, len(users))
}

// join condition
func TestJoinWithCondition(t *testing.T) {
	var users []User
	err := db.Joins("join wallets on wallets.user_id = users.id AND wallets.balance > ?", 200).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 5, len(users))

	users = []User{}
	err = db.Joins("Wallet").Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 22, len(users))
}

// query aggregation
func TestCount(t *testing.T) {
	var count int64
	err := db.Model(&User{}).Joins("Wallet").Where("Wallet.balance > ?", 500000).Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(3), count)
}

type AggregationResult struct {
	TotalBalance int64
	MinBalance   int64
	MaxBalance   int64
	AvgBalance   float64
}

func TestAggregation(t *testing.T) {
	var result AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance", "min(balance)  as min_balance", "max(balance) as max_balance",
		"avg(balance) as avg_balance").Take(&result).Error
	assert.Nil(t, err)

	assert.Equal(t, int64(5000000), result.TotalBalance)
	assert.Equal(t, int64(1000000), result.MinBalance)
	assert.Equal(t, int64(1000000), result.MaxBalance)
	assert.Equal(t, float64(1000000), result.AvgBalance)

}

func TestAggregationHavingGroupBy(t *testing.T) {
	var result []AggregationResult
	err := db.Model(&Wallet{}).Select("sum(balance) as total_balance", "min(balance) as min_balance", "max(balance) as max_balance", "avg(balance) as avg_balance").
		Joins("User").Group("User.id").Having("sum(balance) > ?", 1000000).Find(&result).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}
