package hubspot

import (
	"hubspot-api-go-client/config"

	"github.com/entanglesoftware/hubspot-api-go/configuration"
	"github.com/entanglesoftware/hubspot-api-go/hubspot"
)

func NewClient(apiKey string) *hubspot.Client {

	baseURL := config.GetHubSpotBaseURL()

	config := configuration.Configuration{
		AccessToken:            apiKey,
		BasePath:               baseURL,
		NumberOfAPICallRetries: 3,
	}

	hsClient := hubspot.NewClient(config)
	return hsClient
}
