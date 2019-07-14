package orm

func (Session) TableName() string {
	return "session"
}

func (ApplicationConfig) TableName() string {
	return "application_config"
}

func (Admin) TableName() string {
	return "admin"
}

func (AdminProfile) TableName() string {
	return "admin_profile"
}

func (User) TableName() string {
	return "user"
}

func (Cat) TableName() string {
	return "cat"
}

func (CatThumbnail) TableName() string {
	return "cat_thumbnail"
}

func (AdminTimePeriodTemplate) TableName() string {
	return "admin_time_period_template"
}

func (AdoptionTimePeriod) TableName() string {
	return "adoption_time_period"
}

func (AdoptionTimePeriodCatPivot) TableName() string {
	return "adoption_time_period_cat_pivot"
}

func (CatUserReservation) TableName() string {
	return "cat_user_reservation"
}

func (CatUserTransfer) TableName() string {
	return "cat_user_transfer"
}

func (CatUserAdoption) TableName() string {
	return "cat_user_adoption"
}

func (Wallet) TableName() string {
	return "wallet"
}

func (Banner) TableName() string {
	return "banner"
}

func (ExecuteTask) TableName() string {
	return "execute_task"
}
