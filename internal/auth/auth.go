package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/AbyssExplorer/Cogent/internal/cognito"
	"github.com/AbyssExplorer/Cogent/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

const tokenUrl = "https://%s.auth.%s.amazoncognito.com/oauth2/token"

func GenerateJWTToken(ctx context.Context, poolId, clientId, clientSecret string, conf *aws.Config) (string, error) {

	client := cognito.NewClient(conf)

	describeOutput, err := client.DescribeUserPool(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: aws.String(poolId),
	})
	if err != nil {
		return "", fmt.Errorf("failed to describe user pool: %w", err)
	}

	if describeOutput.UserPool == nil || describeOutput.UserPool.Domain == nil {
		return "", fmt.Errorf("user pool domain not found for Id: %s", poolId)
	}

	domain := aws.ToString(describeOutput.UserPool.Domain)

	tokenURL := fmt.Sprintf(tokenUrl, domain, conf.Region)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned non-200 status: %d", resp.StatusCode)
	}

	var credentials models.Credentials
	if err := json.NewDecoder(resp.Body).Decode(&credentials); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return credentials.AccessToken, nil
}
