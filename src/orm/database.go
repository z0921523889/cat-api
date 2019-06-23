package orm


type User struct {
	Id   int64  `xorm:"varchar(25) notnull unique 'usr_name'"`
	Name string  `xorm:"varchar(25) notnull unique 'usr_name'"`
}

type Cat struct {
	Id   int64
	Name string  `xorm:"varchar(25) notnull unique 'usr_name'"`
}

type CatThumbnail struct {
	Id   int64
	Name string  `xorm:"varchar(25) notnull unique 'usr_name'"`
}