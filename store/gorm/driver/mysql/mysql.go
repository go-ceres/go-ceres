package mysql

import (
	"github.com/go-ceres/go-ceres/store/gorm"
	"gorm.io/driver/mysql"
)

func init() {
	gorm.Register("mysql", func(dns string) gorm.Dialector {
		return mysql.Open(dns)
	})
}
