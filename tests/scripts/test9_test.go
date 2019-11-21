package scripts

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/notify"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test9(t *testing.T) {

	//var state string
	//store = func(i interface{}) {
	//	state = fmt.Sprintf("%v", i)
	//}

	const path = "conf/notify.json"

	Convey("clear db", t, func(ctx C) {
		_ = container.Invoke(func(migrations *migrations.Migrations) {

			// clear database
			// ------------------------------------------------
			err := migrations.Purge()
			So(err, ShouldBeNil)
		})
	})

	Convey("send sms", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			notifyService *notify.Notify) {

			// add templates
			// ------------------------------------------------
			templates := []*m.Template{
				{
					Name:       "sms_body",
					Content:    "[code:block]",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: nil,
				},
				{
					Name:       "code",
					Content:    "[code:content] [code]",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: common.String("sms_body"),
				},
				{
					Name: "template2",
					Content: fmt.Sprintf(`{
 "items": [
   "code"
 ],
 "title": "",
 "fields": [
	{
     "name": "code",
     "value": "Activate code:"
   }
]
}`),
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeTemplate,
					ParentName: nil,
				},
			}

			for _, template := range templates {
				err := adaptors.Template.UpdateOrCreate(template)
				So(err, ShouldBeNil)
			}


			// ------------------------------------------------
			render, err := adaptors.Template.Render("template2", map[string]interface{}{
				"code": 12345,
			})
			So(err, ShouldBeNil)
			So(render.Subject, ShouldEqual, "")
			So(render.Body, ShouldEqual, "Activate code: 12345")

			//read config file
			// ------------------------------------------------
			//var file []byte
			//file, err = ioutil.ReadFile(path)
			//So(err, ShouldBeNil)
			//
			//conf := &notify.NotifyConfig{}
			//err = json.Unmarshal(file, &conf)
			//So(err, ShouldBeNil)
			//
			//notifyService.UpdateCfg(conf)
			//notifyService.Restart()
			//
			//// scripts
			//// ------------------------------------------------
			//storeRegisterCallback(scriptService)
			//
			//scripts := GetScripts(ctx, scriptService, adaptors, 24)
			//
			//engine24, err := scriptService.NewEngine(scripts["script24"])
			//So(err, ShouldBeNil)
			//err = engine24.Compile()
			//So(err, ShouldBeNil)
			//
			//_, err = engine24.DoCustom("main")
			//So(err, ShouldBeNil)
			//
			//time.Sleep(time.Second * 5)
		})
	})
}