package dao

import (
	"CDcoding2333/scaffold/constant"
	"CDcoding2333/scaffold/errs"
	"time"

	log "github.com/sirupsen/logrus"
)

// User ...
type User struct {
	ID       uint       `gorm:"column:id;primary_key"`
	UID      string     `gorm:"column:uid;type:varchar(128);not null;unique_index"`
	Name     string     `gorm:"column:name;type:varchar(256);not null;unique_index"`
	Alias    string     `gorm:"column:alias;type:varchar(256);not null"`
	Avatar   string     `gorm:"column:avatar;type:varchar(256);default:''"`
	Age      int        `gorm:"column:age;type:int(8);default:0"`
	Sex      int        `gorm:"column:sex;type:int(2);default:0"`
	Email    string     `gorm:"column:email;type:varchar(128);default:''"`
	Tel      string     `gorm:"column:tel;type:varchar(128);default:''"`
	Addr     string     `gorm:"column:addr;type:varchar(256);default:''"`
	Password string     `gorm:"column:password;type:varchar(256);default:''"`
	Salt     string     `gorm:"column:salt;type:varchar(256);default:''"`
	Birth    *time.Time `gorm:"column:birth;type:timestamp;not null"`
	State    int        `gorm:"column:state;type:int(4);default:0;index"`
	Model
}

// TableName ...
func (u *User) TableName() string {
	return "users"
}

// Create ...
func (u *User) Create() error {
	if err := db.Create(u).Error; err != nil {
		log.WithError(err).Error("user.Create")
		return errs.New(errs.ErrDatabase, err.Error())
	}
	return nil
}

// GetUsers ...
func (u *User) GetUsers(id ...uint) ([]*User, error) {
	var users []*User
	return users, db.Model(User{}).Where("id in (?)", id).Find(&users).Error
}

// DelUsers ...
func (u *User) DelUsers(id ...uint) error {
	return db.Where("id in (?)", id).Update("state", constant.UserStateDelete).Error
}

// GetUserByName ...
func (u *User) GetUserByName() error {
	return db.Model(User{}).Where("name = ? ", u.Name).First(u).Error
}
