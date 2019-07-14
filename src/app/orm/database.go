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
	Status int64  `gorm:"type:integer;not null;column:status"`
	Hint   string `gorm:"type:varchar(25);column:hint"`
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
	//state:       系統掛售中 : 1 / 期約到期掛售中 : 2/轉讓中 : 3 /領養增值中 :4 / 等待裂變中 :5
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
	StartTime  time.Time `gorm:"not null;column:start_time"`
	EndTime    time.Time `gorm:"not null;column:end_time"`
	Cats       []Cat     `gorm:"many2many:adoption_time_period_cat_pivot"`
}

type AdoptionTimePeriodCatPivot struct {
	gorm.Model           `gorm:"embedded"`
	CatId                uint                 `gorm:"column:cat_id"`
	AdoptionTimePeriodId uint                 `gorm:"column:adoption_time_period_id"`
	UserId               uint                 `gorm:"column:user_id"`
	Cat                  Cat                  `gorm:"foreignkey:CatId"`
	AdoptionTimePeriod   AdoptionTimePeriod   `gorm:"foreignkey:AdoptionTimePeriodId"`
	User                 User                 `gorm:"foreignkey:UserId"`
	CatUserReservations  []CatUserReservation `gorm:"foreignkey:AdoptionTimePeriodCatPivotsId"`
}

type CatUserReservation struct {
	gorm.Model                    `gorm:"embedded"`
	AdoptionTimePeriodCatPivotsId uint                       `gorm:"column:adoption_time_period_cat_pivot_id"`
	UserId                        uint                       `gorm:"column:user_id"`
	User                          User                       `gorm:"foreignkey:UserId"`
	Coin                          int64                      `gorm:"type:integer;not null;column:coin"`
	AdoptionTimePeriodCatPivot    AdoptionTimePeriodCatPivot `gorm:"foreignkey:AdoptionTimePeriodCatPivotId"`
}

type CatUserTransfer struct {
	gorm.Model                    `gorm:"embedded"`
	AdoptionTimePeriodCatPivotId uint                       `gorm:"column:adoption_time_period_cat_pivot_id"`
	CatUserAdoptionId             uint                       `gorm:"column:cat_user_adoption_id"`
	UserId                        uint                       `gorm:"column:user_id"`
	AdoptionTimePeriodCatPivot    AdoptionTimePeriodCatPivot `gorm:"foreignkey:AdoptionTimePeriodCatPivotId"`
	CatUserAdoption               CatUserAdoption            `gorm:"foreignkey:CatUserAdoptionId"`
	User                          User                       `gorm:"foreignkey:UserId"`
	//state:       待交易: 1 / 買家已上傳憑證 : 2/已完成 :3 /已取消
	Status      int64     `gorm:"type:integer;not null;column:status"`
	Certificate []byte    `gorm:"type:bytea;column:certificate"`
	StartTime   time.Time `gorm:"column:start_time"`
	EndTime     time.Time `gorm:"column:end_time"`
}

type CatUserAdoption struct {
	gorm.Model `gorm:"embedded"`
	CatId      uint      `gorm:"column:cat_id"`
	UserId     uint      `gorm:"column:user_id"`
	StartTime  time.Time `gorm:"not null;column:start_time"`
	EndTime    time.Time `gorm:"not null;column:end_time"`
	//state:       增值中: 1 / 已完成 :2 / 已售出 : 3
	Status int64 `gorm:"type:integer;not null;column:status"`
	Cat    Cat   `gorm:"foreignkey:CatId"`
	User   User  `gorm:"foreignkey:UserId"`
}

type Wallet struct {
	gorm.Model `gorm:"embedded"`
	Coin       int64 `gorm:"type:integer;column:coin"`
	PetCoin    int64 `gorm:"type:integer;column:pet_coin"`
}

type Banner struct {
	gorm.Model `gorm:"embedded"`
	Data       []byte `gorm:"type:bytea;column:data"`
	Sort       int64  `gorm:"type:integer;column:sort"`
}

type ExecuteTask struct {
	gorm.Model  `gorm:"embedded"`
	ExecuteType int64     `gorm:"type:integer;not null;column:execute_type"`
	ExecuteTime time.Time `gorm:"not null;column:execute_time"`
	Data        []byte    `gorm:"type:bytea;column:data"`
	Done        bool      `gorm:"not null;column:done"`
}
