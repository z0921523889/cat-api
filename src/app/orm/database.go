package orm

type Users struct {
	Id               int64  `xorm:"int64 autoincr pk 'id'"`
	Name             string `xorm:"varchar(25) notnull unique 'usr_name'"`
	Phone            string `xorm:"varchar(15) notnull unique 'phone'"`
	Password         string `xorm:"varchar(25) notnull unique 'password'"`
	SecurityPassword string `xorm:"varchar(25) notnull unique 'security_password'"`
	CreateAt         string `xorm:"notnull created 'create_at'"`
	UpdateAt         string `xorm:"notnull updated 'update_at'"`
	DeleteAt         string `xorm:"notnull deleted 'delete_at'"`
}

type Cat struct {
	Id int64
	Name             string `xorm:"varchar(25) notnull unique 'name'"`

}

type CatThumbnail struct {
	Id int64
}
