package orm

import (
	"github.com/jinzhu/gorm"
)

type Users struct {
	gorm.Model
	Name             string `gorm:"type:varchar(25);not null;unique;column:usr_name"`
	Phone            string `gorm:"type:varchar(15);not null;unique;column:phone"`
	Password         string `gorm:"type:varchar(25);not null;unique;column:password"`
	SecurityPassword string `gorm:"type:varchar(25);not null;unique;column:security_password"`
}

type Cat struct {
	gorm.Model
	Name             string `gorm:"type:varchar(25);not null;column:name"`
	Level            string `gorm:"type:varchar(25);not null;column:level"`
	Price            int64  `gorm:"type:integer;not null;column:price"`
	PetCoin          int64  `gorm:"type:integer;not null;column:pet_coin"`
	ReservationPrice int64  `gorm:"type:integer;not null;column:reservation_price"`
	AdoptionPrice    int64  `gorm:"type:integer;not null;column:adoption_price"`
	ContractDays     int64  `gorm:"type:integer;not null;column:contract_days"`
	ContractBenefit  int64  `gorm:"type:integer;not null;column:contract_benefit"`
	Status           int64  `gorm:"type:integer;not null;column:status"`
	CatThumbnailId   uint   `gorm:"type:integer;column:cat_thumbnail_id"`
}

type CatThumbnails struct {
	gorm.Model
	Data    []byte `gorm:"type:bytea;column:data"`
	CatList []Cat  `gorm:"foreignkey:CatThumbnailId"`
}
