// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package alexa

import (
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/uuid"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"net/http"
	"strings"
	"sync"
)

var (
	log = common.MustGetLogger("alexa")
)

// Alexa ...
type Alexa struct {
	engine        *gin.Engine
	isStarted     *atomic.Bool
	addressPort   *string
	server        *http.Server
	skillLock     *sync.Mutex
	skills        map[int64]Skill
	adaptors      *adaptors.Adaptors
	appConfig     *config.AppConfig
	token         *atomic.String
	scriptService *scripts.ScriptService
	core          *core.Core
	gateClient    *gate_client.GateClient
}

// NewAlexa ...
func NewAlexa(adaptors *adaptors.Adaptors,
	appConfig *config.AppConfig,
	scriptService *scripts.ScriptService,
	core *core.Core,
	gateClient *gate_client.GateClient) *Alexa {
	return &Alexa{
		isStarted:     atomic.NewBool(false),
		adaptors:      adaptors,
		skillLock:     &sync.Mutex{},
		skills:        make(map[int64]Skill),
		appConfig:     appConfig,
		token:         atomic.NewString(""),
		scriptService: scriptService,
		core:          core,
		gateClient:    gateClient,
	}
}

// Start ...
func (a *Alexa) Start() {

	if a.isStarted.Load() {
		return
	}
	a.isStarted.Store(true)

	a.init()

	a.engine = gin.New()
	a.engine.POST("/*any", a.Auth, a.handlerFunc)

	port := "3033"
	a.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: a.engine,
	}

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err.Error())
		}
	}()

	go func() {
		a.gateClient.SetAlexaApiEngine(a.engine)
	}()

	log.Infof("Serving server at http://[::]:%s", port)
}

func (a *Alexa) init() {

	if err := a.getSettings(); err != nil {
		log.Error(err.Error())
	}

	list, err := a.adaptors.AlexaSkill.ListEnabled(999, 0)
	if err != nil {
		log.Error(err.Error())
		return
	}

	a.skillLock.Lock()
	defer a.skillLock.Unlock()
	for _, skill := range list {
		a.skills[skill.Id] = NewWorker(skill, a.adaptors, a.scriptService, a.core)
	}
}

// Stop ...
func (a *Alexa) Stop() {
	if !a.isStarted.Load() {
		return
	}
	a.isStarted.Store(false)

	if a.server != nil {
		a.server.Close()
	}
}

func (a *Alexa) handlerFunc(ctx *gin.Context) {

	log.Info("new request")

	req := &Request{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Error(err.Error())
		ctx.AbortWithError(400, err)
		return
	}

	resp := NewResponse()

	if req.GetRequestType() == "LaunchRequest" {
		a.onLaunchHandler(ctx, req, resp)
	} else if req.GetRequestType() == "IntentRequest" {
		a.onIntentHandle(ctx, req, resp)
	} else if req.GetRequestType() == "SessionEndedRequest" {
		a.onSessionEndedHandler(ctx, req, resp)
	} else if strings.HasPrefix(req.GetRequestType(), "AudioPlayer.") {
		a.onAudioPlayerHandler(ctx, req, resp)
	} else {
		http.Error(ctx.Writer, "Invalid request.", http.StatusBadRequest)
	}

	ctx.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")

	b, _ := resp.String()
	ctx.Writer.Write(b)
}

func (a *Alexa) onLaunchHandler(ctx *gin.Context, req *Request, resp *Response) {
	a.skillLock.Lock()
	defer a.skillLock.Unlock()

	for _, skill := range a.skills {
		if skill.GetAppID() != req.Context.System.Application.ApplicationID {
			continue
		}
		if skill.OnLaunch != nil {
			skill.OnLaunch(ctx, req, resp)
		}
	}
}

func (a *Alexa) onIntentHandle(ctx *gin.Context, req *Request, resp *Response) {
	a.skillLock.Lock()
	defer a.skillLock.Unlock()

	for _, skill := range a.skills {
		if skill.GetAppID() != req.Context.System.Application.ApplicationID {
			continue
		}
		if skill.OnIntent != nil {
			skill.OnIntent(ctx, req, resp)
		}
	}
}

func (a *Alexa) onSessionEndedHandler(ctx *gin.Context, req *Request, resp *Response) {
	a.skillLock.Lock()
	defer a.skillLock.Unlock()

	for _, skill := range a.skills {
		if skill.GetAppID() != req.Context.System.Application.ApplicationID {
			continue
		}
		if skill.OnSessionEnded != nil {
			skill.OnSessionEnded(ctx, req, resp)
		}
	}
}

func (a *Alexa) onAudioPlayerHandler(ctx *gin.Context, req *Request, resp *Response) {
	a.skillLock.Lock()
	defer a.skillLock.Unlock()

	for _, skill := range a.skills {
		if skill.GetAppID() != req.Context.System.Application.ApplicationID {
			continue
		}
		if skill.OnAudioPlayerState != nil {
			skill.OnAudioPlayerState(ctx, req, resp)
		}
	}
}

func (a *Alexa) genAccessToken() {
	a.token.Store(uuid.NewV4().String())
}

func (a *Alexa) getSettings() (err error) {

	var variable *m.Variable
	if variable, err = a.adaptors.Variable.GetByName("alexa_token"); err != nil || variable == nil {
		a.genAccessToken()
		variable = &m.Variable{
			Name:     "alexa_token",
			Value:    a.token.Load(),
			Autoload: false,
		}
		err = a.adaptors.Variable.Add(variable)
	}
	a.token.Store(variable.Value)
	return
}

func (a *Alexa) updateSettings() (err error) {
	err = a.adaptors.Variable.Update(&m.Variable{
		Name:     "alexa_token",
		Value:    a.token.Load(),
		Autoload: false,
	})
	return
}

// Auth ...
func (a Alexa) Auth(ctx *gin.Context) {

	//accessToken := ctx.Request.URL.Query().Get("alexa_token")
	//
	//if accessToken == "" || accessToken != a.token.Load() {
	//	ctx.AbortWithError(401, errors.New("access token invalid"))
	//	return
	//}
	if !IsValidAlexaRequest(ctx.Writer, ctx.Request) {
		ctx.AbortWithError(401, errors.New("invalid request"))
		return
	}
}

// Add ...
func (a *Alexa) Add(skill *m.AlexaSkill) {
	a.skillLock.Lock()
	defer a.skillLock.Unlock()
	if _, ok := a.skills[skill.Id]; !ok {
		a.skills[skill.Id] = NewWorker(skill, a.adaptors, a.scriptService, a.core)
	}
}

// Update ...
func (a *Alexa) Update(skill *m.AlexaSkill) {
	a.skillLock.Lock()
	defer a.skillLock.Unlock()
	a.skills[skill.Id] = NewWorker(skill, a.adaptors, a.scriptService, a.core)
}

// Delete ...
func (a *Alexa) Delete(skill *m.AlexaSkill) {
	a.skillLock.Lock()
	defer a.skillLock.Unlock()
	if _, ok := a.skills[skill.Id]; !ok {
		delete(a.skills, skill.Id)
	}
}
