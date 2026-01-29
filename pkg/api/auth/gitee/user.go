package gitee

import (
	"encoding/json"
	"strconv"

	"github.com/aide-family/rabbit/pkg/api/auth"
	"github.com/aide-family/rabbit/pkg/config"
)

var _ auth.User = (*GiteeUser)(nil)

type GiteeUser struct {
	AvatarURL         string `json:"avatar_url"`
	Bio               string `json:"bio"`
	Blog              string `json:"blog"`
	CreatedAt         string `json:"created_at"`
	Email             string `json:"email"`
	EventsURL         string `json:"events_url"`
	Followers         uint32 `json:"followers"`
	FollowersURL      string `json:"followers_url"`
	Following         uint32 `json:"following"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	HTMLURL           string `json:"html_url"`
	ID                uint32 `json:"id"`
	Login             string `json:"login"`
	Name              string `json:"name"`
	OrganizationsURL  string `json:"organizations_url"`
	PublicGists       uint32 `json:"public_gists"`
	PublicRepos       uint32 `json:"public_repos"`
	ReceivedEventsURL string `json:"received_events_url"`
	Remark            string `json:"remark"`
	ReposURL          string `json:"repos_url"`
	Stared            uint32 `json:"stared"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	UpdatedAt         string `json:"updated_at"`
	URL               string `json:"url"`
	Watched           uint32 `json:"watched"`
	Weibo             string `json:"weibo"`
}

// GetAPP implements [auth.User].
func (g *GiteeUser) GetAPP() config.OAuth2_APP {
	return config.OAuth2_GITEE
}

// GetAvatar implements [auth.User].
func (g *GiteeUser) GetAvatar() string {
	return g.AvatarURL
}

// GetEmail implements [auth.User].
func (g *GiteeUser) GetEmail() string {
	return g.Email
}

// GetName implements [auth.User].
func (g *GiteeUser) GetName() string {
	return g.Login
}

// GetOpenID implements [auth.User].
func (g *GiteeUser) GetOpenID() string {
	return strconv.FormatUint(uint64(g.ID), 10)
}

// GetRaw implements [auth.User].
func (g *GiteeUser) GetRaw() []byte {
	raw, _ := json.Marshal(g)
	return raw
}

// GetRemark implements [auth.User].
func (g *GiteeUser) GetRemark() string {
	return g.Remark
}

// GetNickname implements [auth.User].
func (g *GiteeUser) GetNickname() string {
	return g.Name
}
