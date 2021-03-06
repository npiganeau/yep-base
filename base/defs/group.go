// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package defs

import (
	"github.com/npiganeau/yep/pool"
	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/models/security"
)

func initGroups() {
	models.NewModel("Group")
	group := pool.Group()
	group.AddCharField("GroupID", models.StringFieldParams{Required: true})
	group.AddCharField("Name", models.StringFieldParams{Required: true, Translate: true})

	group.Methods().Create().Extend("",
		func(rs pool.GroupSet, data *pool.GroupData) pool.GroupSet {
			if rs.Env().Context().HasKey("GroupForceCreate") {
				return rs.Super().Create(data)
			}
			log.Panic("Trying to create a security group")
			panic("Unreachable")
		})

	group.Methods().Write().Extend("",
		func(rs pool.GroupSet, data *pool.GroupData, fieldsToUnset ...models.FieldNamer) {
			log.Panic("Trying to modify a security group")
		})

	group.AddMethod("ReloadGroups",
		`ReloadGroups populates the Group table with groups from the security.Registry
		and refresh all memberships.`,
		func(rs pool.GroupSet) {
			log.Debug("Reloading groups")
			// Sync groups
			pool.Group().NewSet(rs.Env()).FetchAll().Unlink()
			for _, group := range security.Registry.AllGroups() {
				rs.WithContext("GroupForceCreate", true).Create(&pool.GroupData{
					GroupID: group.ID,
					Name:    group.Name,
				})
			}
			// Sync memberships
			for _, user := range pool.User().NewSet(rs.Env()).FetchAll().Records() {
				secGroups := security.Registry.UserGroups(user.ID())
				grpIds := make([]string, len(secGroups))
				i := 0
				for grp := range secGroups {
					grpIds[i] = grp.ID
					i++
				}
				groups := pool.Group().Search(rs.Env(), pool.Group().GroupID().In(grpIds))
				user.SetGroups(groups)
			}
		})
}
