package models

// OAuthUser chứa thông tin chung của user từ các provider OAuth
type OAuthUser struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	Picture       string
	Provider      string
}
