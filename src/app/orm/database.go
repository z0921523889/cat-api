package orm

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Sessions struct {
	Token  string    `gorm:"type:text;primary_key;column:token"`
	Data   []byte    `gorm:"type:bytea;not null;column:data"`
	Expiry time.Time `gorm:"not null;column:expiry"`
}

type ApplicationConfigs struct {
	Key   string `gorm:"type:text;primary_key;column:key"`
	Value string `gorm:"type:text;column:value"`
}

type Admins struct {
	gorm.Model
	Account                     string                     `gorm:"type:varchar(25);not null;unique;column:account"`
	Password                    string                     `gorm:"type:varchar(25);not null;column:password"`
	AdminTimePeriodTemplateList []AdminTimePeriodTemplates `gorm:"foreignkey:AdminId"`
}

type AdminProfiles struct {
	gorm.Model
	AdminId     uint   `gorm:"type:integer;not null;column:admin_id"`
	permissions []byte `gorm:"type:bytea;not null;column:permissions"`
}

type Users struct {
	gorm.Model
	Phone            string `gorm:"type:varchar(15);not null;unique;column:phone"`
	UserName             string `gorm:"type:varchar(25);not null;column:usr_name"`
	Password         string `gorm:"type:varchar(25);not null;column:password"`
	SecurityPassword string `gorm:"type:varchar(25);not null;column:security_password"`
	IdentifiedCode   string `gorm:"type:varchar(25);not null;column:identified_code"`
	//state:       有效的 : 1 / 無效的 : 2/ 被封鎖的 : 3
	Status int64 `gorm:"type:integer;not null;column:status"`
}

type Cats struct {
	gorm.Model
	Name             string `gorm:"type:varchar(25);not null;column:name"`
	Level            string `gorm:"type:varchar(25);not null;column:level"`
	Price            int64  `gorm:"type:integer;not null;column:price"`
	PetCoin          int64  `gorm:"type:integer;not null;column:pet_coin"`
	ReservationPrice int64  `gorm:"type:integer;not null;column:reservation_price"`
	AdoptionPrice    int64  `gorm:"type:integer;not null;column:adoption_price"`
	ContractDays     int64  `gorm:"type:integer;not null;column:contract_days"`
	ContractBenefit  int64  `gorm:"type:integer;not null;column:contract_benefit"`
	//state:       系統代售中 : 1 / 待售中 : 2/預約中 : 3 /確認交易 :4 / 待交貨 : 5 /收養中 : 6
	Status              int64                 `gorm:"type:integer;not null;column:status"`
	CatThumbnailId      uint                  `gorm:"type:integer;column:cat_thumbnail_id"`
	AdoptionTimePeriods []AdoptionTimePeriods `gorm:"many2many:adoption_time_period_cat_pivot"`
}

type CatThumbnails struct {
	gorm.Model
	Data []byte `gorm:"type:bytea;column:data"`
}

type AdminTimePeriodTemplates struct {
	gorm.Model
	AdminId uint      `gorm:"type:integer;column:admin_id"`
	StartAt time.Time `gorm:"not null;column:start_time"`
	EndAt   time.Time `gorm:"not null;column:end_time"`
}

type AdoptionTimePeriods struct {
	gorm.Model
	StartAt time.Time `gorm:"not null;column:start_time"`
	EndAt   time.Time `gorm:"not null;column:end_time"`
	Cats    []Cats    `gorm:"many2many:adoption_time_period_cat_pivots"`
}

type AdoptionTimePeriodCatPivots struct {
	gorm.Model
	CatsId                uint `gorm:"type:integer;column:cats_id"`
	AdoptionTimePeriodsId uint `gorm:"type:integer;column:adoption_time_periods_id"`
}

type CatUserReservations struct {
	gorm.Model
	CatsId  uint `gorm:"type:integer;column:cats_id"`
	UsersId uint `gorm:"type:integer;column:users_id"`
}
