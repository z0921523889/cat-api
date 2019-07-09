package orm

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Session struct {
	Token  string    `gorm:"type:text;primary_key;column:token"`
	Data   []byte    `gorm:"type:bytea;not null;column:data"`
	Expiry time.Time `gorm:"not null;column:expiry"`
}

type ApplicationConfig struct {
	Key   string `gorm:"type:text;primary_key;column:key"`
	Value string `gorm:"type:text;column:value"`
}

type Admin struct {
	gorm.Model               `gorm:"embedded"`
	Account                  string                    `gorm:"type:varchar(25);not null;unique;column:account"`
	Password                 string                    `gorm:"type:varchar(25);not null;column:password"`
	UserId                   uint                      `gorm:"column:user_id"`
	User                     User                      `gorm:"foreignkey:UserId"`
	AdminTimePeriodTemplates []AdminTimePeriodTemplate `gorm:"foreignkey:AdminId"`
}

type AdminProfile struct {
	gorm.Model  `gorm:"embedded"`
	AdminId     uint   `gorm:"not null;column:admin_id"`
	permissions []byte `gorm:"type:bytea;not null;column:permissions"`
}

type User struct {
	gorm.Model       `gorm:"embedded"`
	Phone            string `gorm:"type:varchar(15);not null;unique;column:phone"`
	UserName         string `gorm:"type:varchar(25);not null;column:usr_name"`
	Password         string `gorm:"type:varchar(25);not null;column:password"`
	SecurityPassword string `gorm:"type:varchar(25);not null;column:security_password"`
	IdentifiedCode   string `gorm:"type:varchar(25);not null;column:identified_code"`
	WalletId         uint   `gorm:"column:wallet_id"`
	Wallet           Wallet `gorm:"foreignkey:WalletId"`
	//state:       有效的 : 1 / 無效的 : 2/ 被封鎖的 : 3
	Status int64 `gorm:"type:integer;not null;column:status"`
}

type Cat struct {
	gorm.Model       `gorm:"embedded"`
	Name             string `gorm:"type:varchar(25);not null;column:name"`
	Level            string `gorm:"type:varchar(25);not null;column:level"`
	Price            int64  `gorm:"type:integer;not null;column:price"`
	Deposit          int64  `gorm:"type:integer;not null;column:deposit"`
	PetCoin          int64  `gorm:"type:integer;not null;column:pet_coin"`
	ReservationPrice int64  `gorm:"type:integer;not null;column:reservation_price"`
	AdoptionPrice    int64  `gorm:"type:integer;not null;column:adoption_price"`
	ContractDays     int64  `gorm:"type:integer;not null;column:contract_days"`
	ContractBenefit  int64  `gorm:"type:integer;not null;column:contract_benefit"`
	//state:       系統代售中 : 1 / 待售中 : 2/預約中 : 3 /確認交易 :4 / 待交貨 : 5 /收養中 : 6
	Status              int64                `gorm:"type:integer;not null;column:status"`
	CatThumbnailId      uint                 `gorm:"column:cat_thumbnail_id"`
	CatThumbnail        CatThumbnail         `gorm:"foreignkey:CatThumbnailId"`
	AdoptionTimePeriods []AdoptionTimePeriod `gorm:"many2many:adoption_time_period_cat_pivot"`
}

type CatThumbnail struct {
	gorm.Model `gorm:"embedded"`
	Data       []byte `gorm:"type:bytea;column:data"`
}

type AdminTimePeriodTemplate struct {
	gorm.Model `gorm:"embedded"`
	AdminId    uint      `gorm:"column:admin_id"`
	StartAt    time.Time `gorm:"not null;column:start_time"`
	EndAt      time.Time `gorm:"not null;column:end_time"`
	Admin      Admin     `gorm:"foreignkey:AdminId"`
}

type AdoptionTimePeriod struct {
	gorm.Model `gorm:"embedded"`
	StartAt    time.Time `gorm:"not null;column:start_time"`
	EndAt      time.Time `gorm:"not null;column:end_time"`
	Cats       []Cat     `gorm:"many2many:adoption_time_period_cat_pivot"`
}

type AdoptionTimePeriodCatPivot struct {
	gorm.Model           `gorm:"embedded"`
	CatId                uint               `gorm:"column:cat_id"`
	AdoptionTimePeriodId uint               `gorm:"column:adoption_time_period_id"`
	Cat                  Cat                `gorm:"foreignkey:CatId"`
	AdoptionTimePeriod   AdoptionTimePeriod `gorm:"foreignkey:AdoptionTimePeriodId"`
}

type CatUserReservation struct {
	gorm.Model                    `gorm:"embedded"`
	AdoptionTimePeriodCatPivotsId uint                       `gorm:"column:adoption_time_period_cat_pivot_id"`
	UserId                        uint                       `gorm:"column:user_id"`
	User                          User                       `gorm:"foreignkey:UserId"`
	AdoptionTimePeriodCatPivot    AdoptionTimePeriodCatPivot `gorm:"foreignkey:AdoptionTimePeriodCatPivotId"`
}

type Wallet struct {
	gorm.Model `gorm:"embedded"`
	Coin       int64 `gorm:"type:integer;column:coin"`
	PetCoin    int64 `gorm:"type:integer;column:pet_coin"`
}

type Banner struct {
	gorm.Model `gorm:"embedded"`
	Data       []byte `gorm:"type:bytea;column:data"`
	Sort      int64  `gorm:"type:integer;column:sort"`
}
