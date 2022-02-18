package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const ALLOWED_SPECIAL_CHARS = "!#$%&'()*+,-./:;<=>?@[]^_`{|}~"

type ChangePasswordRequest struct {
	// Current is the current password
	Current string `json:"current" binding:"required"`
	// New is the new password to be set
	New string `json:"new" binding:"required"`
}

// internalLdapIsEnabled checks if the ldap service is in the list of running services
func internalLdapIsEnabled() bool {
	services := os.Getenv("SERVICES")
	if strings.Contains(services, "ldap") {
		return true
	} else {
		return false
	}
}

// ChangePasswordSupported is a simple endpoint to check if you can change
// change a user's password (i.e. are running the internal LDAP auth provider)
// @Summary Check if password changing is supported
// @Schemes
// @Tags Password
// @Description This endpoint is used to check if changing a user's password via this API is
// @Description supported. Internal password management is only supported if this Hoss server
// @Description is using the internal LDAP auth provider.
// @Accept json
// @Produce json
// @Success 200 {object} object{changePasswordIsSupported=bool}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /password [get]
func (a *Auth) ChangePasswordSupported(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"changePasswordIsSupported": internalLdapIsEnabled()})
}

// ChangePassword is an endpoint for a user to change their own password
// if the internal LDAP auth provider is in use.
// @Summary Change a password
// @Schemes
// @Tags Password
// @Description This endpoint is used to change the authenticated user's password.
// @Description Internal password management is only supported if this Hoss server is using the internal LDAP auth provider.
// @Description If the provided password fails the password policy, a 400 will be returned with an error message.
// @Description If the current password is incorrect, a 403 will be returned with an error message.
// @Accept json
// @Produce json
// @Param	changePasswordInput		body	api.ChangePasswordRequest	true	"Change Password Input"
// @Success 204
// @Failure 400 {object} object{error=string}
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /password [put]
func (a *Auth) ChangePassword(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check to make sure that the internal `ldap` service is running.
	// If it is not, we currently assume user accounts will be managed via
	// some external means and we should not be changing passwords in the
	// Hoss UI
	if !internalLdapIsEnabled() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Users are managed via an external system." +
			" You cannot change your password directly via the Hoss API"})
		return
	}

	var pass ChangePasswordRequest
	err = c.BindJSON(&pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check password policy
	// If the user's new provided password does not meet the password policy as set by the
	// system admin via the service config file, return a Bad Request response and a message
	if len(pass.New) < a.Settings.PasswordPolicy.MinLength {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("New password must be at least %v characters.", a.Settings.PasswordPolicy.MinLength)})
		return
	}
	if a.Settings.PasswordPolicy.RequireUppercase {
		hasUpper := false
		for _, ch := range pass.New {
			if unicode.IsUpper(ch) {
				hasUpper = true
				break
			}
		}

		if !hasUpper {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": "New password must include an uppercase character."})
			return
		}
	}
	if a.Settings.PasswordPolicy.RequireSpecial {
		hasSpecial := false
		for _, ch := range pass.New {
			if strings.ContainsRune(ALLOWED_SPECIAL_CHARS, ch) {
				hasSpecial = true
				break
			}
		}

		if !hasSpecial {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": fmt.Sprintf("New password must include a special character (%s).", ALLOWED_SPECIAL_CHARS)})
			return
		}
	}

	// Run command to change password
	// We use the OpenLDAP utility `ldappasswd` to change the user's password.
	// It will bind to the ldap server as the user, with their current password,
	// and then run a series of commands to hash the provided "new" password and
	// set it as the user's current password. If this command completes successfully
	// the password has changed. If it does not, the old password will continue to work.
	domainParts := strings.Split(os.Getenv("LDAP_DOMAIN"), ".")
	dnPart := strings.Join(domainParts, ",dc=")
	dn := fmt.Sprintf("cn=%s,ou=People,dc=%s", userInfo.Username, dnPart)

	cmd := exec.Command("ldappasswd", "-H", "ldap://ldap", "-x",
		"-D", dn,
		"-w", pass.Current,
		"-a", pass.Current,
		"-s", pass.New)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password."})
		return
	}

	if err := cmd.Wait(); err != nil {
		if cmd.ProcessState.ExitCode() == 49 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect current password."})
			return
		}

		if cmd.ProcessState.ExitCode() != 0 {
			logrus.Warningf("An error occurred while changing the password for user %s: %v - Stdout: %s - Stderr: %s",
				userInfo.Username, err.Error(), stdout.String(), stderr.String())
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": "An error occurred while trying to change the password. Try again or contact your administrator."})
			return
		}
	}

	c.Status(http.StatusNoContent)
}
