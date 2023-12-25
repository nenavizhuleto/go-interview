package megaplan

// Megaplan API
type MegaPlan struct {
	Url         string // Megaplan API url. example "http://megaplan.company.com"
	Responsible string // ID of who is going to be responsible at the time of task creation
	Token       Token  // Access token
}

// Creates new megaplan API instance. Must be authenticated after
func New(url string, authOpt *AuthOpt) *MegaPlan {
	mp := &MegaPlan{
		Url:         url,
		Responsible: authOpt.Responsible,
	}

	return mp
}

// Authorization options
type AuthOpt struct {
	Username       string // Username to authenticate megaplan
	Password       string // Password to authenticate megaplan
	GrantType      string // Type of authentication. Defaults to 'password'
	AccessTokenUrl string // URL to get access token from. Defaults to '/auth/access_token'
	Responsible    string // ID of who is going to be responsible at the time of task creation
}

// NewAuthOpt constructs authorization options with default values
func NewAuthOpt(username, password, responsible string) *AuthOpt {
	return &AuthOpt{
		Username:       username,
		Password:       password,
		GrantType:      "password",
		AccessTokenUrl: "/auth/access_token",
		Responsible:    responsible,
	}
}

// Access Token. Required to communicate with Megaplan API
type Token struct {
	AccessToken  string `json:"access_token"`  // Access token body
	ExpiresIn    int    `json:"expires_in"`    // Time when access token becomes invalid
	TokenType    string `json:"token_type"`    // Type of token. Defaults to 'bearer'
	Scope        string `json:"scope"`         // To which scope access token is related
	RefreshToken string `json:"refresh_token"` // Different token to refresh current one
}

// getToken returns token string formatted the way to satisfy headers format
func (mp *MegaPlan) getToken() string {
	return "Bearer " + mp.Token.AccessToken
}

// General Response from the Megaplan API
type Response struct {
	Meta Meta          `json:"meta"`
	Data []interface{} `json:"data"` // Actual Data returned from the Megaplan API
}

// Meta information of the response. aka. Status, Errors, Pagination etc.
type Meta struct {
	Status     int        `json:"status"`     // Status of the response
	Errors     []string   `json:"errors"`     // Array of errors if any
	Pagination Pagination `json:"pagination"` // Current pagination
}

type Pagination struct {
	Count int `json:"count"` // Count
	Limit int `json:"limit"` // Limit
}
