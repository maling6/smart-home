package api

import (
	"fmt"
	"time"
	"github.com/astaxie/beego/validation"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"./routers"
	"./filters"
	"./core"
	"./cron"
	"../api/log"
)

func Initialize() {
	//log.Println("Api initialize...")

	// site base
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_host := beego.AppConfig.String("db_host")
	db_name := beego.AppConfig.String("db_name")
	db_port := beego.AppConfig.String("db_port")
	// parseTime https://github.com/go-sql-driver/mysql#parsetime
	db := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", db_user, db_pass, db_host, db_port, db_name)
	orm.RegisterDataBase("default", "mysql", db, 30, 30)
	// Timezone http://beego.me/docs/mvc/model/orm.md#timezone-config
	orm.DefaultTimeLoc, _ = time.LoadLocation("Asia/Novosibirsk")

	routers.Initialize()

	beego.Info("AppPath:", beego.AppPath)
	if(beego.BConfig.RunMode == "dev") {
		beego.Info("Development mode enabled")
		// orm debug mode
		if orm_debug, _ := beego.AppConfig.Bool("orm_debug"); orm_debug {
			orm.Debug = true
		}

		//beego.BConfig.WebConfig.DirectoryIndex = true
		//beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		//beego.SetStaticPath("/admin/static", "static_source/admin")
	} else {
		beego.Info("Product mode enabled")
		//beego.SetStaticPath("/admin/static", "www/admin")
	}

	// register access filters
	filters.RegisterFilters()

	validation.SetDefaultMessage(map[string]string{
		"Required":     "Должно быть заполнено",
		"Min":          "Минимально допустимое значение %d",
		"Max":          "Максимально допустимое значение %d",
		"Range":        "Должно быть в диапазоне от %d до %d",
		"MinSize":      "Минимально допустимая длина %d",
		"MaxSize":      "Максимально допустимая длина %d",
		"Length":       "Длина должна быть равна %d",
		"Alpha":        "Должно состоять из букв",
		"Numeric":      "Должно состоять из цифр",
		"AlphaNumeric": "Должно состоять из букв или цифр",
		"Match":        "Должно совпадать с %s",
		"NoMatch":      "Не должно совпадать с %s",
		"AlphaDash":    "Должно состоять из букв, цифр или символов (-_)",
		"Email":        "Должно быть в правильном формате email",
		"IP":           "Должен быть правильный IP адрес",
		"Base64":       "Должно быть представлено в правильном формате base64",
		"Mobile":       "Должно быть правильным номером мобильного телефона",
		"Tel":          "Должно быть правильным номером телефона",
		"Phone":        "Должно быть правильным номером телефона или мобильного телефона",
		"ZipCode":      "Должно быть правильным почтовым индексом",
	})

	// cron
	cron.Initialize()

	// core
	if err := core.Initialize(); err != nil {
		log.Fatal(err.Error())
	}
}
