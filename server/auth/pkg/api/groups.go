package api

import (
	"net/http"

	"github.com/gigantum/hoss-auth/pkg/userinfo"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type CreateGroupRequest struct {
	// GroupName is the unique name to identify the group
	GroupName string `json:"name" binding:"required"`
	// Description is a short description of the group's purpose
	Description *string `json:"description" binding:"required"`
}

// CreateGroup creates a new group
// @Summary Create a group
// @Schemes
// @Tags Groups
// @Description This endpoint is used to create a new group. The authorized user must have the `admin` or `privileged` roles.
// @Accept json
// @Produce json
// @Param	createGroupInput		body	api.CreateGroupRequest	true	"Create Group Input"
// @Success 201 {object} database.Group
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /group/ [post]
func (a *Auth) CreateGroup(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if privileged := validatePrivileged(userInfo.Role); !privileged {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		return
	}

	var req CreateGroupRequest
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create database entry
	group, err := a.Database.CreateGroup(req.GroupName, *req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add user to the group
	_, err = a.Database.UpdateGroupMembership(userInfo.Username, group.GroupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, group)
}

// DeleteGroup deletes an existing group
// @Summary Delete a group
// @Schemes
// @Tags Groups
// @Description This endpoint is used to delete an existing group. The authorized user must have the `admin` or
// @Description the `privileged` role and be a member of the group.
// @Accept json
// @Produce json
// @Param   groupName   path      string  true  "Group Name"
// @Success 204
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /group/{groupName} [delete]
func (a *Auth) DeleteGroup(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupName := c.Param("groupname")

	if roleValid := a.validateAdminOrPrivilegedMember(userInfo.Role, userInfo.Username, groupName); !roleValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		return
	}

	err = a.Database.DeleteGroup(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetGroup gets an existing group
// @Summary Get a group
// @Schemes
// @Tags Groups
// @Description This endpoint is used to fetch details about an existing group. If the authorized user does not have the
// @Description `admin` or `privileged` roles, they must be a member of the group.
// @Accept json
// @Produce json
// @Param   groupName   path      string  true  "Group Name"
// @Success 200 {object} database.Group
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /group/{groupName} [get]
func (a *Auth) GetGroup(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupName := c.Param("groupname")

	if privileged := validatePrivileged(userInfo.Role); !privileged {
		// if user is not privileged, they must be a member
		if member := a.Database.ValidateMembership(userInfo.Username, groupName); !member {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
			return
		}
	}

	group, err := a.Database.GetGroup(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// UpdateGroupUser updates membership for a user in a group
// @Summary Modify a user's group membership
// @Schemes
// @Tags Groups
// @Description This endpoint updates/adds the membership of a user to a group. The authorized user must have the `admin` or
// @Description the `privileged` role and be a member of the group.
// @Accept json
// @Produce json
// @Param   groupName   path      string  true  "Group Name to modify"
// @Param   username   path      string  true  "Username to add"
// @Success 200 {object} database.Membership
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /group/{groupName}/user/{username} [put]
func (a *Auth) UpdateGroupUser(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupName := c.Param("groupname")
	username := c.Param("username")

	if roleValid := a.validateAdminOrPrivilegedMember(userInfo.Role, userInfo.Username, groupName); !roleValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		return
	}

	membership, err := a.Database.UpdateGroupMembership(username, groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, membership)
}

// RemoveGroupUser removes membership for a user from a group
// @Summary Remove a user's group membership
// @Schemes
// @Tags Groups
// @Description This endpoint removes the membership of a user from a group. The authorized user must have the `admin` or
// @Description the `privileged` role and be a member of the group.
// @Accept json
// @Produce json
// @Param   groupName   path      string  true  "Group Name to modify"
// @Param   username   path      string  true  "Username to remove"
// @Success 204
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /group/{groupName}/user/{username} [delete]
func (a *Auth) RemoveGroupUser(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupName := c.Param("groupname")
	username := c.Param("username")

	if roleValid := a.validateAdminOrPrivilegedMember(userInfo.Role, userInfo.Username, groupName); !roleValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
		return
	}

	if groupName == "public" && userInfo.Role != ROLE_ADMIN {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only administrators can remove users from the 'public' group."})
		return
	}

	err = a.Database.RemoveGroupMembership(username, groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUser fetches the information about a user
// @Summary Get a user
// @Schemes
// @Tags Groups
// @Description This endpoint returns data about the user and their group memberships
// @Accept json
// @Produce json
// @Param   username   path      string  true  "Username"
// @Success 200 {object} database.User
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /user/{username} [get]
func (a *Auth) GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := a.Database.GetUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsernames fetches usernames based on something else
// @Summary Lookup usernames
// @Schemes
// @Tags Groups
// @Description This endpoint looks up usernames based on other data. Currently only searching by email
// @Description is supported. Since emails are not required to be unique, multiple usernames could be returned.
// @Accept json
// @Produce json
// @Param   email   query      string  true  "User's Email Address"
// @Success 200 {object} object{usernames=[]string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /usernames [get]
func (a *Auth) GetUsernames(c *gin.Context) {
	email := c.Query("email")

	users, err := a.Database.GetUsersByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var usernames []string
	for _, u := range *users {
		usernames = append(usernames, u.Username)
	}

	c.JSON(http.StatusOK, gin.H{"usernames": usernames})
}

func (a *Auth) addToAutoGroups(u userinfo.UserInfo) error {
	// If an admin, add to the admin group
	if u.Role == ROLE_ADMIN {
		_, err := a.Database.UpdateGroupMembership(u.Username, "admin")
		if err != nil {
			return errors.Wrap(err, "failed to add user to admin group")
		}
	}

	// Make sure the user is in the public group
	_, err := a.Database.UpdateGroupMembership(u.Username, "public")
	if err != nil {
		return errors.Wrap(err, "failed to add user to public group")
	}

	return nil
}
