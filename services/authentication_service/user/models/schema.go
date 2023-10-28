package models

type AuthenticationSchema struct {
	AccessToken           string  `json:"accessToken"`
	ExpiresIn             int     `json:"expiresIn"`
	RefreshTokenExpiresIn int     `json:"refreshTokenExpiresIn"`
	RefreshToken          string  `json:"refreshToken"`
	TokenType             string  `json:"tokenType"`
	IDToken               *string `json:"idToken,omitempty"`
	NotBeforePolicy       int     `json:"notBeforePolicy"`
	SessionState          string  `json:"sessionState"`
	Scope                 string  `json:"scope"`
}
