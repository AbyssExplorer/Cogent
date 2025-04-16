package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func NewClient(cfg *aws.Config) *cognitoidentityprovider.Client {
	return cognitoidentityprovider.NewFromConfig(*cfg)
}

func LoadAWSConfig(ctx context.Context, region string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load default SDK configuration: %w", err)
	}
	return &cfg, nil
}

func ListUserPools(ctx context.Context, client *cognitoidentityprovider.Client) ([]string, error) {
	var poolIDs []string
	var nextToken *string
	const maxResults = 50
	for {
		output, err := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int32(maxResults),
			NextToken:  nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list user pools: %w", err)
		}
		for _, pool := range output.UserPools {
			poolIDs = append(poolIDs, *pool.Id)
		}
		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return poolIDs, nil
}

func ListAppClients(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID string) ([]types.UserPoolClientDescription, error) {
	var clients []types.UserPoolClientDescription
	var nextToken *string
	const maxResults = 50
	for {
		output, err := client.ListUserPoolClients(ctx, &cognitoidentityprovider.ListUserPoolClientsInput{
			UserPoolId: aws.String(userPoolID),
			MaxResults: aws.Int32(maxResults),
			NextToken:  nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list app clients for pool %s: %w", userPoolID, err)
		}
		clients = append(clients, output.UserPoolClients...)
		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return clients, nil
}

func GetAppClientSecret(ctx context.Context, client *cognitoidentityprovider.Client, userPoolId *string, clientId *string) (string, error) {
	output, err := client.DescribeUserPoolClient(ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
		UserPoolId: userPoolId,
		ClientId:   clientId,
	})
	if err != nil {
		return "", fmt.Errorf("failed to describe user pool client %s in pool %s: %w", aws.ToString(clientId), aws.ToString(userPoolId), err)
	}
	if output.UserPoolClient != nil && output.UserPoolClient.ClientSecret != nil {
		return *output.UserPoolClient.ClientSecret, nil
	}
	return "", fmt.Errorf("client secret not found for client ID: %s in pool: %s", aws.ToString(clientId), aws.ToString(userPoolId))
}
