package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"store/pkg/config"
	"store/pkg/errors"
	"store/pkg/model"
)

func NewClient(conf *config.GlobalConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		conf.Mysql.Name,
		conf.Mysql.Password,
		conf.Mysql.Address,
		conf.Mysql.Port,
		conf.Mysql.Database,
		conf.Mysql.Conf,
	)
	// 初始化数据库时的高级配置
	client, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		return nil, errors.DBInitFailed
	}

	if err := client.AutoMigrate(
		&model.UserRole{},
		&model.ArticleJudgeRecord{},
		&model.User{},
	); err != nil {
		return nil, errors.AutoMigrateFailed
	}

	return client, nil
}
