package conf

import (
	"cat-api/src/app/orm"
	"github.com/jinzhu/gorm"
	"log"
)

var DefaultConfig map[string]string

func init() {
	DefaultConfig = make(map[string]string)
	DefaultConfig["TimePeriodMaxCatAmount"] = "5"
	DefaultConfig["TimePeriodTemplateAdminId"] = "1"
}

func CheckDataBaseConfig() {
	for key, value := range DefaultConfig {
		var config orm.ApplicationConfig
		err := orm.Engine.First(&config, "key = ?", key).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				config = orm.ApplicationConfig{
					Key:   key,
					Value: value,
				}
				if err := orm.Engine.Create(&config).Error; err != nil {
					log.Fatalln("fail to create config data err : " + err.Error())
				}
			} else {
				log.Fatalln("fail to find config data err : " + err.Error())
			}
		}
		log.Println("config data already exist")
	}

	var configs []orm.ApplicationConfig
	err := orm.Engine.Find(&configs).Error
	if err != nil {
		log.Fatalln("fail to find config data err : " + err.Error())
	}
	for _, configValue := range configs {
		DefaultConfig[configValue.Key] = configValue.Value
	}
}
