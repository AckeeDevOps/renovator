package main

// TokenConfig contains configuration for managed tokens
type TokenConfig struct {
	Tokens []TokenDetails `json:"tokens"`
}

// TokenDetails represents token's details
type TokenDetails struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}
