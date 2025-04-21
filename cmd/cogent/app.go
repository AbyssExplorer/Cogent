package app

import (
	"context"
	"flag"
	"fmt"

	"github.com/AbyssExplorer/Cogent/internal/auth"
	"github.com/AbyssExplorer/Cogent/internal/cognito"
	"github.com/AbyssExplorer/Cogent/internal/prompt"
	"github.com/AbyssExplorer/Cogent/internal/sysutils"
)

func Execute() error {

	var health bool
	flag.BoolVar(&health, "testcheck", false, "Perform health check")

	var region string
	flag.StringVar(&region, "region", "us-east-1", "AWS Region to use")
	flag.Parse()

	if health {
		fmt.Println("I am OK!")
		return nil
	}

	ctx := context.Background()
	conf, err := cognito.LoadAWSConfig(ctx, region)
	if err != nil {
		return err
	}

	client := cognito.NewClient(conf)
	pools, err := cognito.ListUserPools(ctx, client)
	if err != nil {
		return err
	}
	if len(pools) == 0 {
		fmt.Println("No Cognito User Pools found in the current region.")
		return nil
	}
	poolId, err := prompt.ForUserPool(pools)
	if err != nil {
		return err
	}

	fmt.Printf("Selected User Pool ID: %s\n", poolId)
	appClients, err := cognito.ListAppClients(ctx, client, poolId)
	if err != nil {
		return err
	}
	if len(appClients) == 0 {
		fmt.Printf("No App Clients found in User Pool: %s\n", poolId)
		return nil
	}

	clientId, clientSecret, err := prompt.ForAppClient(ctx, client, appClients)
	if err != nil {
		return err
	}
	fmt.Printf("Selected App Client ID: %s\n", clientId)

	token, err := auth.GenerateJWTToken(ctx, poolId, clientId, clientSecret, conf)
	if err != nil {
		return err
	}
	fmt.Println("Generated Client Credentials JWT Token:")
	fmt.Println(token)

	sysutils.CopyToClipboard(token)
	fmt.Println("\nToken copied to clipboard.")

	return nil
}
