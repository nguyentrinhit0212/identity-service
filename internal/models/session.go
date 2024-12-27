package models

type AuthToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

type AuthSession struct {
	User          User               `json:"user"`
	CurrentTenant Tenant             `json:"currentTenant"`
	Tenants       []UserTenantAccess `json:"tenants"`
	Tokens        AuthToken          `json:"tokens"`
}
