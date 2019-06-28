package orm

import (
	"github.com/jinzhu/gorm"
	"time"
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
	//state:        待放養 : 0 /預約中 : 1 /繁殖中 : 2 /收養中 : 3
	Status         int64 `gorm:"type:integer;not null;column:status"`
	CatThumbnailId uint  `gorm:"type:integer;column:cat_thumbnail_id"`
}

type CatThumbnails struct {
	gorm.Model
	Data    []byte `gorm:"type:bytea;column:data"`
	CatList []Cat  `gorm:"foreignkey:CatThumbnailId"`
}

type AdoptionTimePeriod struct {
	gorm.Model
	StartAt time.Time `gorm:"not null;column:start_time"`
	EndAt   time.Time `gorm:"not null;column:end_time"`
}

type Sessions struct {
	Token  string    `gorm:"type:text;primary_key;column:token"`
	Data   []byte    `gorm:"type:bytea;not null;column:data"`
	Expiry time.Time `gorm:"not null;column:expiry"`
}
