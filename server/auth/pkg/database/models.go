package database

const (
	ROLE_ADMIN      = "admin"
	ROLE_PRIVILEGED = "privileged"
	ROLE_USER       = "user"
)

const (
	DEFAULT_GROUP_SUFFIX = "hoss-default-group"
)

// User object containing a user's username
type User struct {
	Id       int64  `json:"-"`
	Username string `json:"username"`

	FullName      string `json:"full_name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Subject       string `json:"-"`

	Role        string        `pg:"type:role" json:"role"`
	Memberships []*Membership `pg:"rel:has-many" json:"memberships,omitempty"`
}

// Group object containing a list of user's who are part of the group
type Group struct {
	Id          int64  `json:"-"`
	GroupName   string `json:"group_name"`
	Description string `json:"description"`

	Memberships []*Membership `pg:"rel:has-many" json:"memberships,omitempty"`
}

// Membership object tracks a user's membership in a group
type Membership struct {
	GroupId int64  `json:"-"`
	Group   *Group `pg:"rel:has-one" json:"group,omitempty"`

	UserId int64 `json:"-"`
	User   *User `pg:"rel:has-one" json:"user,omitempty"`
}

// Token object containing a token's information and the user who owns this token
type PersonalAccessToken struct {
	Id          int64  `json:"-"`
	PAT         string `json:"pat,omitempty"`
	Description string `json:"description"`

	OwnerId int64 `json:"-"`
	Owner   *User `pg:"rel:has-one" json:"-"`
}
