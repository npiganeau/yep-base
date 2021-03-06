// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package base

import (
	"encoding/base64"
	"io/ioutil"
	"path"

	// Import this module's defs
	_ "github.com/npiganeau/yep-base/base/defs"
	"github.com/npiganeau/yep/pool"
	"github.com/npiganeau/yep/yep/actions"
	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/models/security"
	"github.com/npiganeau/yep/yep/models/types"
	"github.com/npiganeau/yep/yep/server"
	"github.com/npiganeau/yep/yep/tools/generate"
	"github.com/npiganeau/yep/yep/tools/logging"
)

const (
	MODULE_NAME string = "base"
	SEQUENCE    uint8  = 100
	NAME        string = "Base"
	VERSION     string = "0.1"
	CATEGORY    string = "Hidden"
	DESCRIPTION string = `
The kernel of YEP, needed for all installation
==============================================
	`
	AUTHOR     string = "NDP Systèmes"
	MAINTAINER string = "NDP Systèmes"
	WEBSITE    string = "http://www.ndp-systemes.fr"
)

var log *logging.Logger

func init() {
	log = logging.GetLogger("base")
	server.RegisterModule(&server.Module{
		Name: MODULE_NAME,
		PostInit: func() {
			err := models.ExecuteInNewEnvironment(security.SuperUserID, func(env models.Environment) {

				mainCompany := pool.Company().Search(env, pool.Company().ID().Equals(1))
				if mainCompany.IsEmpty() {
					mainCompany = pool.Company().Create(env, &pool.CompanyData{
						ID:   1,
						Name: "Your Company",
					})
					env.Cr().Execute("ALTER SEQUENCE company_id_seq RESTART WITH 2")
				}

				adminPartner := pool.Partner().Search(env, pool.Partner().ID().Equals(1))
				if adminPartner.IsEmpty() {
					adminPartner = pool.Partner().Create(env, &pool.PartnerData{
						ID:       1,
						Lang:     "en_US",
						Name:     "Administrator",
						Function: "IT Manager",
					})
					env.Cr().Execute("ALTER SEQUENCE partner_id_seq RESTART WITH 2")
				}

				avatarImg, _ := ioutil.ReadFile(path.Join(generate.YEPDir, "yep", "server", "static", "base", "src", "img", "avatar.png"))

				adminUser := pool.User().Search(env, pool.User().ID().Equals(security.SuperUserID))
				ActionID := actions.MakeActionRef("base_action_res_users")
				if adminUser.IsEmpty() {
					pool.User().Create(env, &pool.UserData{
						ID:         security.SuperUserID,
						Name:       "Administrator",
						Active:     true,
						Company:    mainCompany,
						Login:      "admin",
						LoginDate:  types.DateTime{},
						Password:   "admin",
						Partner:    adminPartner,
						ActionID:   ActionID,
						ImageSmall: base64.StdEncoding.EncodeToString(avatarImg),
					})
					env.Cr().Execute("ALTER SEQUENCE user_id_seq RESTART WITH 2")
				}

				pool.Group().NewSet(env).ReloadGroups()
			})
			if err != nil {
				log.Panic("Error while initializing", "error", err)
			}
		},
	})
}
