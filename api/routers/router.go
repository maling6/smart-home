package routers

import (
	"github.com/e154/smart-home/api/controllers"
	"github.com/astaxie/beego"
	"github.com/e154/smart-home/api/stream"
)

func Initialize() {
	beego.ErrorController(&controllers.ErrorController{})

	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/signin", &controllers.AuthController{}, "get:SignIn"),
			beego.NSRouter("/signout", &controllers.AuthController{}, "post:SignOut"),
			beego.NSRouter("/recovery", &controllers.AuthController{}, "post:Recovery"),
			beego.NSRouter("/reset", &controllers.AuthController{}, "post:Reset"),
			beego.NSRouter("/access_list", &controllers.AuthController{}, "get:AccessList"),

			beego.NSRouter("/ws)", &stream.StreamCotroller{}, "get:Get"),
			beego.NSRouter("/ws/*)", &stream.StreamCotroller{}, "get:Get"),

			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "get:GetOne"),
			beego.NSRouter("/node", &controllers.NodeController{}, "get:GetAll"),
			beego.NSRouter("/node", &controllers.NodeController{}, "post:Post"),
			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "put:Put"),
			beego.NSRouter("/node/:id([0-9]+)", &controllers.NodeController{}, "delete:Delete"),

			beego.NSRouter("/device/:id([0-9]+)", &controllers.DeviceController{}, "get:GetOne"),
			beego.NSRouter("/device", &controllers.DeviceController{}, "get:GetAll"),
			beego.NSRouter("/device", &controllers.DeviceController{}, "post:Post"),
			beego.NSRouter("/device/:id([0-9]+)", &controllers.DeviceController{}, "put:Put"),
			beego.NSRouter("/device/:id([0-9]+)", &controllers.DeviceController{}, "delete:Delete"),
			beego.NSRouter("/device/group", &controllers.DeviceController{}, "get:GetGroup"),
			beego.NSRouter("/device/:id([0-9]+)/actions", &controllers.DeviceController{}, "get:GetActions"),
			beego.NSRouter("/device/search", &controllers.DeviceController{}, "get:Search"),

			beego.NSRouter("/workflow/:id([0-9]+)", &controllers.WorkflowController{}, "get:GetOne"),
			beego.NSRouter("/workflow", &controllers.WorkflowController{}, "get:GetAll"),
			beego.NSRouter("/workflow", &controllers.WorkflowController{}, "post:Post"),
			beego.NSRouter("/workflow/:id([0-9]+)", &controllers.WorkflowController{}, "put:Put"),
			beego.NSRouter("/workflow/:id([0-9]+)", &controllers.WorkflowController{}, "delete:Delete"),
			beego.NSRouter("/workflow/:id([0-9]+)/update_scenario", &controllers.WorkflowController{}, "put:UpdateScenario"),

			beego.NSRouter("/flow/:id([0-9]+)", &controllers.FlowController{}, "get:GetOne"),
			beego.NSRouter("/flow/:id([0-9]+)/full", &controllers.FlowController{}, "get:GetOneFull"),
			beego.NSRouter("/flow/:id([0-9]+)/redactor", &controllers.FlowController{}, "get:GetOneRedactor"),
			beego.NSRouter("/flow/:id([0-9]+)/redactor", &controllers.FlowController{}, "put:UpdateRedactor"),
			beego.NSRouter("/flow", &controllers.FlowController{}, "get:GetAll"),
			beego.NSRouter("/flow", &controllers.FlowController{}, "post:Post"),
			beego.NSRouter("/flow/:id([0-9]+)", &controllers.FlowController{}, "put:Put"),
			beego.NSRouter("/flow/:id([0-9]+)", &controllers.FlowController{}, "delete:Delete"),
			beego.NSRouter("/flow/:id([0-9]+)/workers", &controllers.FlowController{}, "get:GetWorkers"),
			beego.NSRouter("/flow/search", &controllers.FlowController{}, "get:Search"),

			beego.NSRouter("/device_action/:id([0-9]+)", &controllers.DeviceActionController{}, "get:GetOne"),
			beego.NSRouter("/device_action", &controllers.DeviceActionController{}, "get:GetAll"),
			beego.NSRouter("/device_action", &controllers.DeviceActionController{}, "post:Post"),
			beego.NSRouter("/device_action/:id([0-9]+)", &controllers.DeviceActionController{}, "put:Put"),
			beego.NSRouter("/device_action/:id([0-9]+)", &controllers.DeviceActionController{}, "delete:Delete"),
			beego.NSRouter("/device_action/search", &controllers.DeviceActionController{}, "get:Search"),
			beego.NSRouter("/device_action/get_by_device/:id([0-9]+)", &controllers.DeviceActionController{}, "get:GetByDevice"),

			beego.NSRouter("/worker/:id([0-9]+)", &controllers.WorkerController{}, "get:GetOne"),
			beego.NSRouter("/worker", &controllers.WorkerController{}, "get:GetAll"),
			beego.NSRouter("/worker", &controllers.WorkerController{}, "post:Post"),
			beego.NSRouter("/worker/:id([0-9]+)", &controllers.WorkerController{}, "put:Put"),
			beego.NSRouter("/worker/:id([0-9]+)", &controllers.WorkerController{}, "delete:Delete"),

			beego.NSRouter("/script/:id([0-9]+)", &controllers.ScriptController{}, "get:GetOne"),
			beego.NSRouter("/script", &controllers.ScriptController{}, "get:GetAll"),
			beego.NSRouter("/script", &controllers.ScriptController{}, "post:Post"),
			beego.NSRouter("/script/:id([0-9]+)", &controllers.ScriptController{}, "put:Put"),
			beego.NSRouter("/script/:id([0-9]+)", &controllers.ScriptController{}, "delete:Delete"),
			beego.NSRouter("/script/:id([0-9]+)/exec", &controllers.ScriptController{}, "post:Exec"),
			beego.NSRouter("/script/search", &controllers.ScriptController{}, "get:Search"),

			beego.NSRouter("/log/:id([0-9]+)", &controllers.LogController{}, "get:GetOne"),
			beego.NSRouter("/log", &controllers.LogController{}, "get:GetAll"),
			beego.NSRouter("/log", &controllers.LogController{}, "post:Post"),
			beego.NSRouter("/log/:id([0-9]+)", &controllers.LogController{}, "put:Put"),
			beego.NSRouter("/log/:id([0-9]+)", &controllers.LogController{}, "delete:Delete"),

			beego.NSRouter("/email/template/:name([\\w]+)", &controllers.EmailTemplateController{}, "get:GetOne"),
			beego.NSRouter("/email/templates", &controllers.EmailTemplateController{}, "get:GetAll"),
			beego.NSRouter("/email/template", &controllers.EmailTemplateController{}, "post:Post"),
			beego.NSRouter("/email/template/:name([\\w]+)", &controllers.EmailTemplateController{}, "put:Put"),
			beego.NSRouter("/email/template/:name([\\w]+)", &controllers.EmailTemplateController{}, "delete:Delete"),
			beego.NSRouter("/email/template/search", &controllers.EmailTemplateController{}, "get:Search"),
			beego.NSRouter("/email/preview", &controllers.EmailTemplateController{}, "post:Preview")		,
			beego.NSRouter("/email/item/:name([\\w]+)", &controllers.EmailItemController{}, "get:GetOne"),
			beego.NSRouter("/email/items", &controllers.EmailItemController{}, "get:GetAll"),
			beego.NSRouter("/email/item", &controllers.EmailItemController{}, "post:Post"),
			beego.NSRouter("/email/item/:name([\\w]+)", &controllers.EmailItemController{}, "put:Put"),
			beego.NSRouter("/email/item/:name([\\w]+)", &controllers.EmailItemController{}, "delete:Delete"),
			beego.NSRouter("/email/items/tree", &controllers.EmailItemController{}, "get:GetTree"),
			beego.NSRouter("/email/items/tree", &controllers.EmailItemController{}, "post:UpdateTree"),
			beego.NSRouter("/notifr/:id([0-9]+)", &controllers.NotifrController{}, "get:GetOne"),
			beego.NSRouter("/notifr", &controllers.NotifrController{}, "get:GetAll"),
			beego.NSRouter("/notifr", &controllers.NotifrController{}, "post:Post"),
			beego.NSRouter("/notifr/:id([0-9]+)", &controllers.NotifrController{}, "delete:Delete"),

			beego.NSRouter("/map/:id([0-9]+)", &controllers.MapController{}, "get:GetOne"),
			beego.NSRouter("/map/:id([0-9]+)/full", &controllers.MapController{}, "get:GetFull"),
			beego.NSRouter("/map", &controllers.MapController{}, "get:GetAll"),
			beego.NSRouter("/map", &controllers.MapController{}, "post:Post"),
			beego.NSRouter("/map/:id([0-9]+)", &controllers.MapController{}, "put:Put"),
			beego.NSRouter("/map/:id([0-9]+)", &controllers.MapController{}, "delete:Delete"),
			beego.NSRouter("/map_layer/:id([0-9]+)", &controllers.MapLayerController{}, "get:GetOne"),
			beego.NSRouter("/map_layer", &controllers.MapLayerController{}, "get:GetAll"),
			beego.NSRouter("/map_layer", &controllers.MapLayerController{}, "post:Post"),
			beego.NSRouter("/map_layer/:id([0-9]+)", &controllers.MapLayerController{}, "put:Put"),
			beego.NSRouter("/map_layer/:id([0-9]+)", &controllers.MapLayerController{}, "delete:Delete"),
			beego.NSRouter("/map_layer/sort", &controllers.MapLayerController{}, "put:Sort"),
			beego.NSRouter("/map_element/:id([0-9]+)", &controllers.MapElementController{}, "get:GetOne"),
			beego.NSRouter("/map_element", &controllers.MapElementController{}, "get:GetAll"),
			beego.NSRouter("/map_element", &controllers.MapElementController{}, "post:Post"),
			beego.NSRouter("/map_element/:id([0-9]+)", &controllers.MapElementController{}, "put:Put"),
			beego.NSRouter("/map_element/:id([0-9]+)", &controllers.MapElementController{}, "delete:Delete"),
			beego.NSRouter("/map_element/:id([0-9]+)/element_only", &controllers.MapElementController{}, "put:PutElementOnly"),
			beego.NSRouter("/map_element/sort", &controllers.MapElementController{}, "put:Sort"),

			beego.NSRouter("/device_state/:id([0-9]+)", &controllers.DeviceStateController{}, "get:GetOne"),
			beego.NSRouter("/device_state", &controllers.DeviceStateController{}, "get:GetAll"),
			beego.NSRouter("/device_state", &controllers.DeviceStateController{}, "post:Post"),
			beego.NSRouter("/device_state/:id([0-9]+)", &controllers.DeviceStateController{}, "put:Put"),
			beego.NSRouter("/device_state/:id([0-9]+)", &controllers.DeviceStateController{}, "delete:Delete"),
			beego.NSRouter("/device_state/get_by_device/:id([0-9]+)", &controllers.DeviceStateController{}, "get:GetByDevice"),

			beego.NSRouter("/image/:id([0-9]+)", &controllers.ImageController{}, "get:GetOne"),
			beego.NSRouter("/image", &controllers.ImageController{}, "get:GetAll"),
			beego.NSRouter("/image", &controllers.ImageController{}, "post:Post"),
			beego.NSRouter("/image/upload", &controllers.ImageController{}, "post:Upload"),
			beego.NSRouter("/image/:id([0-9]+)", &controllers.ImageController{}, "put:Put"),
			beego.NSRouter("/image/:id([0-9]+)", &controllers.ImageController{}, "delete:Delete"),

			beego.NSRouter("/dashboard/:id([0-9]+)", &controllers.DashboardController{}, "get:GetOne"),
			beego.NSRouter("/dashboard", &controllers.DashboardController{}, "get:GetAll"),
			beego.NSRouter("/dashboard", &controllers.DashboardController{}, "post:Post"),
			beego.NSRouter("/dashboard/:id([0-9]+)", &controllers.DashboardController{}, "put:Put"),
			beego.NSRouter("/dashboard/:id([0-9]+)", &controllers.DashboardController{}, "delete:Delete"),

			beego.NSRouter("/user/:id([0-9]+)", &controllers.UserController{}, "get:GetOne"),
			beego.NSRouter("/user", &controllers.UserController{}, "get:GetAll"),
			beego.NSRouter("/user", &controllers.UserController{}, "post:Post"),
			beego.NSRouter("/user/:id([0-9]+)", &controllers.UserController{}, "put:Put"),
			beego.NSRouter("/user/:id([0-9]+)/update_status", &controllers.UserController{}, "put:UpdateStatus"),
			beego.NSRouter("/user/:id([0-9]+)", &controllers.UserController{}, "delete:Delete"),
			beego.NSRouter("/user/search", &controllers.UserController{}, "get:Search"),
			beego.NSRouter("/role/:name([\\w]+)", &controllers.RoleController{}, "get:GetOne"),
			beego.NSRouter("/role", &controllers.RoleController{}, "get:GetAll"),
			beego.NSRouter("/role", &controllers.RoleController{}, "post:Post"),
			beego.NSRouter("/role/:name([\\w]+)", &controllers.RoleController{}, "put:Put"),
			beego.NSRouter("/role/:name([\\w]+)", &controllers.RoleController{}, "delete:Delete"),
			beego.NSRouter("/role/search", &controllers.RoleController{}, "get:Search"),
			beego.NSRouter("/role/:name([\\w]+)/access_list", &controllers.RoleController{}, "get:GetAccessList"),
			beego.NSRouter("/role/:name([\\w]+)/access_list", &controllers.RoleController{}, "put:PutAccessList"),

			beego.NSRouter("/scenario/:id([0-9]+)", &controllers.ScenarioController{}, "get:GetOne"),
			beego.NSRouter("/scenario", &controllers.ScenarioController{}, "get:GetAll"),
			beego.NSRouter("/scenario", &controllers.ScenarioController{}, "post:Post"),
			beego.NSRouter("/scenario/:id([0-9]+)", &controllers.ScenarioController{}, "put:Put"),
			beego.NSRouter("/scenario/:id([0-9]+)", &controllers.ScenarioController{}, "delete:Delete"),
			beego.NSRouter("/scenario/search", &controllers.ScenarioController{}, "get:Search"),
			beego.NSRouter("/scenario_script/:id([0-9]+)", &controllers.ScenarioScriptController{}, "get:GetOne"),
			beego.NSRouter("/scenario_script", &controllers.ScenarioScriptController{}, "get:GetAll"),
			beego.NSRouter("/scenario_script", &controllers.ScenarioScriptController{}, "post:Post"),
			beego.NSRouter("/scenario_script/:id([0-9]+)", &controllers.ScenarioScriptController{}, "put:Put"),
			beego.NSRouter("/scenario_script/:id([0-9]+)", &controllers.ScenarioScriptController{}, "delete:Delete"),
	),
	)
	beego.AddNamespace(ns)
}
