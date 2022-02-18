package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gigantum/hoss-auth/pkg/state"
	userinfo "github.com/gigantum/hoss-auth/pkg/userinfo"
)

// Ping returns simple response to show the service is reachable
// @Summary Check if the server is alive
// @Schemes
// @Tags Ping
// @Description This open endpoint returns a simple response if the server is running
// @Accept json
// @Produce json
// @Success 200 {object} object{alive=bool}
// @Router /ping [get]
func (a *Auth) Ping(c *gin.Context) {
	log.Println("Ping")

	c.JSON(200, gin.H{
		"alive": "true",
	})
}

type PatResponse struct {
	Id          int64  `json:"id"`
	Description string `json:"description"`
	Token       string `json:"token,omitempty"`
}

type PatRequest struct {
	Description string `json:"description" binding:"required"`
}

// JWTForPAT exchanges a JWT token for a new personal access token
// @Summary Generate a new PAT
// @Schemes
// @Tags Tokens
// @Description This endpoint is used to create a new Personal Access Token (PAT) for the authorized user.
// @Description The resulting PAT is only accessible once in the response.
// @Description The `Authorization` header should contain a valid JWT with the format `Bearer <id_token>`.
// @Accept json
// @Produce json
// @Success 201 {object} PatResponse
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /pat [post]
func (a *Auth) JWTForPAT(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var patRequest PatRequest
	err = c.BindJSON(&patRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	patRecord, err := a.Database.CreatePAT(userInfo.Username, patRequest.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	patStruct := PatResponse{
		Id:          patRecord.Id,
		Description: patRecord.Description,
		Token:       patRecord.PAT,
	}

	c.JSON(http.StatusCreated, patStruct)
}

// PATforJWT exchanges a JWT token for a personal access token
// @Summary Exchange a PAT for JWT tokens
// @Schemes
// @Tags Tokens
// @Description This endpoint is used to exchange a PAT sent in the header of the request for
// @Description JWTs that can then be used with the rest of the system to make API calls.
// @Description For most requests, the resulting `id_token` value should be used as the Bearer token.
// @Description The PAT should be sent in the `Authorization` header as `Bearer hp_xxxxxxx`
// @Accept json
// @Produce json
// @Success 201 {object} state.Tokens
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /pat/exchange/jwt [post]
func (a *Auth) PATforJWT(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokens, err := a.CreateGigantumTokens(userInfo, state.Scopes{
		ScopeOpenID:  true,
		ScopeProfile: true,
		ScopeEmail:   true,
		ScopeHOSS:    true,
	}, "", false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

// ListPAT lists all personal access tokens for a user
// @Summary List PATs
// @Schemes
// @Tags Tokens
// @Description This endpoint will list all Personal Access Tokens (PAT) for the authorized user.
// @Description Note, this will not return the actual PAT values, as they are only available at create time.
// @Description The `Authorization` header should contain a valid JWT with the format `Bearer <id_token>`.
// @Accept json
// @Produce json
// @Success 200 {array} object{id=string,description=string}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /pat/ [get]
func (a *Auth) ListPAT(c *gin.Context) {
	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	pats, err := a.Database.ListPAT(userInfo.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if pats == nil {
		// Explicitly return an empty array if no PATs exist instead of null
		c.JSON(http.StatusOK, []string{})
		return
	}

	patList := []PatResponse{}
	for _, pat := range pats {
		patList = append(patList, PatResponse{Id: pat.Id, Description: pat.Description})
	}
	c.JSON(http.StatusOK, patList)
}

// DeletePAT deletes a user's personal access token
// @Summary Delete a PAT
// @Schemes
// @Tags Tokens
// @Description This endpoint will delete a Personal Access Token (PAT) by its ID.
// @Description The `Authorization` header should contain a valid JWT with the format `Bearer <id_token>`.
// @Accept json
// @Produce json
// @Param	patId   path      int  true  "PAT ID"
// @Success 204
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /pat/{patId} [delete]
func (a *Auth) DeletePAT(c *gin.Context) {

	patIdStr := c.Param("id")
	patId, _ := strconv.ParseInt(patIdStr, 10, 64)

	userInfo, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = a.Database.DeletePAT(patId, userInfo.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SATforJWT exchanges a service access token for a JWT
// @Summary Exchange a SAT for JWT tokens
// @Schemes
// @Tags Tokens
// @Description This endpoint is used to exchange a Service Account Token (SAT) sent in the header of the request for
// @Description Service Account JWTs that can then be used with the rest of the system to make API calls **as the service account**.
// @Description The SAT should be sent in the `Authorization` header as `Bearer hsvc_xxxxxxx`
// @Accept json
// @Produce json
// @Success 201 {object} state.Tokens
// @Failure 401 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerToken
// @Router /sat/exchange/jwt [post]
func (a *Auth) SATforJWT(c *gin.Context) {
	// If a service account token, in the authorization middleware the token will be stored in the `sat` key.
	_, ok := c.Get("sat")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "SAT to JWT endpoint only available to service accouts"})
		return
	}

	tokens, err := a.CreateGigantumTokens(
		userinfo.UserInfo{Subject: "HOSS-Service"},
		state.Scopes{ScopeOpenID: true},
		"",
		true,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}
