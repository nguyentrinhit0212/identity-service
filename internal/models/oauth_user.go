package models

type OAuthUser struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	Picture       string
	Provider      string
}
