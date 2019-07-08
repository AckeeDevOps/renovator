package client

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"

	"encoding/json"

	"gopkg.in/resty.v1"
)

// Client contains Resty client and Vault address
type Client struct {
	VaultAddress string
	RestClient   *resty.Client
	Debug        bool
}

// TokenLookupData contains all Vault token details
type TokenLookupData struct {
	Accessor     string `json:"accessor"`
	CreationTime int64  `json:"creation_time"`
	CreationTTL  int64  `json:"creation_ttl"`
	DisplayName  string `json:"display_name"`
	ExpireTime   string `json:"expire_time"`
	IssueTime    string `json:"issue_time"`
	Renewable    bool   `json:"renewable"`
	TTL          int64  `json:"ttl"`
}

// RenewalRequest contains a new TTL
type RenewalRequest struct {
	Increment int64 `json:"increment"`
}

// TokenLookupResponse wraps TokenLookupData
type TokenLookupResponse struct {
	Data TokenLookupData `json:"data"`
}

// NewClient creates a new Client
func NewClient(vaultAddress string, insecure bool, debug bool) (*Client, error) {
	c := new(Client)

	u, err := url.Parse(vaultAddress)
	if err != nil {
		return nil, fmt.Errorf("could not parse Vault address: %s", err)
	}

	c.VaultAddress = u.Scheme + "://" + u.Host

	// create a new Resty client
	client := resty.New()

	// disable TLS verification if needed
	if insecure {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	c.RestClient = client
	c.Debug = debug
	return c, nil
}

// LookupSelf gets token's data using the actual Vault token
// so no super-powered token is needed
func (c *Client) LookupSelf(token string) (*TokenLookupData, error) {

	// print debug info if needed
	if c.Debug {
		log.Printf("obtaining details for token '%s...'", token[0:8])
	}

	// get token details from Vault
	resp, err := c.RestClient.R().
		SetHeader("X-Vault-Token", token).
		Get(c.VaultAddress + "/v1/auth/token/lookup-self")
	if err != nil {
		return nil, err
	}

	// print debug info if needed
	if c.Debug {
		log.Printf("code: %d; body: %s", resp.StatusCode(), resp.Body())
	}

	// check response code
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode())
	}

	// parse response data
	lookupReponse := TokenLookupResponse{}
	err = json.Unmarshal(resp.Body(), &lookupReponse)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal Vault response: %s", err)
	}

	return &lookupReponse.Data, nil
}

// Renew adds increment to the token's TTL
func (c Client) Renew(token string, increment int64) (*TokenLookupData, error) {

	// print debug info if needed
	if c.Debug {
		log.Printf("sending renewal request for token '%s' with TTL %d ...", token, increment)
	}

	resp, err := c.RestClient.R().
		SetHeader("X-Vault-Token", token).
		SetBody(RenewalRequest{Increment: increment}).
		Post(c.VaultAddress + "/v1/auth/token/renew-self")
	if err != nil {
		return nil, err
	}

	// print debug info if needed
	if c.Debug {
		log.Printf("code: %d; body: %s", resp.StatusCode(), resp.Body())
	}

	// check response code
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode())
	}

	// get token details again
	tokenDetails, err := c.LookupSelf(token)
	if err != nil {
		return nil, fmt.Errorf("could not get token details: %s", err)
	}

	return tokenDetails, nil
}

// CheckOrRenew gets the current TTL and renews the token if needed
func (c Client) CheckOrRenew(token string, threshold int64, increment int64) (*TokenLookupData, error) {
	tokenDetails, err := c.LookupSelf(token)
	if err != nil {
		return nil, fmt.Errorf("Could not get details for token: %s... ", token[0:8])
	}

	// do not renew tokens with sufficient TTL
	if tokenDetails.TTL <= threshold {

		// print debug info if needed
		if c.Debug {
			log.Printf("Current TTL of token '%s...' is %d days ...", token[0:8], (tokenDetails.TTL / 60 / 60 / 24))
		}

		// it seems that increment is not increment but total TTL
		// hence we're appending increment to the current TTL
		// tested with Vault v0.11.4
		tokenDetailsNew, err := c.Renew(token, tokenDetails.TTL+increment)
		if err != nil {
			return nil, fmt.Errorf("could not renew token %s... ", token[0:8])
		}

		if tokenDetailsNew.TTL > tokenDetails.TTL {

			// print debug info if needed
			if c.Debug {
				log.Printf("Renewed; Token: %s; new TTL: %d; old TTL: %d", token[0:8], tokenDetailsNew.TTL, tokenDetails.TTL)
			}
			return tokenDetailsNew, nil
		}

		return nil, fmt.Errorf("TTL for token %s... was not increased", token[0:8])
	}

	// print debug info if needed
	if c.Debug {
		aboveThreshold := (tokenDetails.TTL - threshold) / 60 / 60 / 24
		log.Printf("renewal not needed for token %s... (%d days above the threshold)", token[0:8], aboveThreshold)
	}

	return tokenDetails, nil
}
