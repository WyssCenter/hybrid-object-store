package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	service "github.com/gigantum/hoss-service"

	"github.com/gigantum/hoss-core/pkg/config"
	"github.com/gigantum/hoss-core/pkg/database"
	"github.com/gigantum/hoss-core/pkg/store"
	"github.com/gin-gonic/gin"
)

func SyncAllUserGroups(config *config.Configuration, db *database.Database, stores map[string]store.ObjectStore) {
	logrus.Infof("Starting user sync goroutine (%d minutes interval)", config.Server.SyncFrequencyMinutes)
	// endless loop
	for {
		// wait 5 minutes between syncs
		time.Sleep(time.Minute * time.Duration(config.Server.SyncFrequencyMinutes))

		// get tokens for use with auth service
		tokens, err := service.GetServiceJWT(config.Server.AuthService)
		if err != nil {
			log.Fatalf("Unable to get tokens: %s\n", err.Error())
			return
		}

		limit := 100
		offset := 0
		for {
			// list users
			users, err := db.ListUsers(limit, offset)
			if err != nil {
				log.Printf("Unable to list users for group membership syncing: %s\n", err.Error())
				return
			}
			if len(users) == 0 {
				break
			}

			for _, user := range users {

				if user.Username == "HOSS-Service" {
					continue
				}

				// for each user, get a list of their groups from auth
				req, err := http.NewRequest("GET", config.Server.AuthService+"/user/"+user.Username, nil)
				if err != nil {
					log.Fatalf("Unable to sync user %s: %s\n", user.Username, err.Error())
					continue
				}

				req.Header.Set("Origin", os.Getenv("EXTERNAL_HOSTNAME"))
				req.Header.Set("Authorization", "Bearer "+tokens.IDToken)

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalf("Unable to sync user %s: %s\n", user.Username, err.Error())
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					log.Printf("Status Code: %d", resp.StatusCode)
					continue
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Unable to sync user %s: %s\n", user.Username, err.Error())
					continue
				}

				var authUser database.User
				if err := json.Unmarshal(body, &authUser); err != nil {
					log.Printf("Unable to sync user %s: %s\n", user.Username, err.Error())
					continue
				}

				var groups []string
				for _, membership := range authUser.Memberships {
					groups = append(groups, membership.Group.GroupName)
				}

				// validate against their current memberships, fix any differences, and re-render their policies
				err = SyncUserGroups(user, groups, db, stores)
				if err != nil {
					log.Printf("Unable to sync user %s: %s\n", user.Username, err.Error())
					continue
				}
			}
			offset += limit
		}
	}
}

// @Summary Sync a user's group details
// @Schemes
// @Description Primarily an internal endpoint that is used to trigger a synchronization of the user's groups with the auth service.
// @Tags User
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /user/sync [get]
func SyncCurrentUserGroups(c *gin.Context) {
	_, db := getAppConfig(c)

	userInfo := getUserInfo(c)

	// create user if does not exist yet
	user, err := db.GetOrCreateUser(userInfo.Username)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = SyncUserGroups(user, userInfo.Groups, db, getStores(c))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func SyncUserGroups(user *database.User, groups []string, db *database.Database, stores map[string]store.ObjectStore) error {

	// if there are any groups the user is no longer part of, remove them
	for _, membership := range user.Memberships {

		membershipOutdated := true
		for _, groupName := range groups {
			if membership.Group.GroupName == groupName {
				membershipOutdated = false
				break
			}
		}

		if membershipOutdated {
			err := db.RemoveGroupMembership(user.Username, membership.Group.GroupName)
			if err != nil {
				return err
			}
		}
	}

	// for all the groups user is part of, create if they don't exist yet and make sure user is added to them
	for _, groupName := range groups {
		if _, err := db.GetOrCreateGroup(groupName); err != nil {
			return err
		}
		err := db.UpdateGroupMembership(user.Username, groupName)
		if err != nil {
			return err
		}
	}

	// rerender policy for user in each store
	offset := 0
	limit := 25
	for {
		objStores, err := db.ListObjectStores(limit, offset)
		if err != nil {
			return err
		}

		if len(objStores) == 0 {
			return nil
		}

		for _, objStore := range objStores {
			currentStore, err := getStoreByName(stores, objStore.Name)
			if err != nil {
				return err
			}

			// re-render the user's policies
			perms, err := db.GetPermissionsByUser(objStore, user.Username, false)
			if err != nil {
				return err
			}

			err = currentStore.SetUserPolicy(user.Username, perms)
			if err != nil {
				return err
			}
		}
		offset += limit
	}
}
