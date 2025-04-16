package prompt

import (
	"context"
	"fmt"
	"strings"

	"github.com/AbyssExplorer/Cogent/internal/cognito"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/manifoldco/promptui"
)

func ForUserPool(pools []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select a Cognito User Pool",
		Items: pools,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}
	return result, nil
}

func ForAppClient(ctx context.Context, client *cognitoidentityprovider.Client, clients []types.UserPoolClientDescription) (string, string, error) {

	appClients := make([]string, len(clients))
	for i, client := range clients {
		appClients[i] = aws.ToString(client.ClientName)
	}

	prompt := promptui.Select{
		Label: "Select an App Client",
		Items: appClients,
		Searcher: func(input string, index int) bool {
			return strings.Contains(strings.ToLower(appClients[index]), strings.ToLower(input))
		},
		Size: 10,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", "", fmt.Errorf("prompt failed: %w", err)
	}
	selectedClient := clients[index]

	clientSecret, err := cognito.GetAppClientSecret(ctx, client, selectedClient.UserPoolId, selectedClient.ClientId)
	if err != nil {
		return "", "", fmt.Errorf("failed to get app client secret for client %s in pool %s: %w", *selectedClient.ClientId, *selectedClient.UserPoolId, err)
	}
	return aws.ToString(selectedClient.ClientId), clientSecret, nil
}
