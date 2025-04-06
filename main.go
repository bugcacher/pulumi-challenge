package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	pulumi "github.com/bugcacher/open-feature-pulumi-esc-provider/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

func main() {
	fmt.Println("Starting application")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	accessToken := os.Getenv("PULUMI_ACCESS_TOKEN")

	// Pulumi ESC environment configuration
	orgName := os.Getenv("PULUMI_ORG_NAME")
	projectName := os.Getenv("PULUMI_PROJECT_NAME")
	envName := os.Getenv("PULUMI_ENV_NAME")

	// (Optional) Custom backend URL if you're self-hosting Pulumi ESC
	backendURL, _ := url.Parse("https://api.pulumi.com")

	// Initialize the OpenFeature provider backed by Pulumi ESC
	provider, err := pulumi.NewPulumiESCProvider(
		orgName,
		projectName,
		envName,
		accessToken,
		pulumi.WithCustomBackendUrl(*backendURL), // Optional ProviderOption
	)
	if err != nil {
		fmt.Printf("Failed to initialize Pulumi ESC provider: %v\n", err)
		return
	}

	// Register the provider globally
	if err := openfeature.SetProviderAndWait(provider); err != nil {
		fmt.Printf("Failed to set provider: %v\n", err)
		return
	}

	// Create a new OpenFeature client
	client := openfeature.NewClient("example-app")
	ctx := context.Background()

	// Fetch secret stored in AWS Secrets Manager via Pulumi ESC
	awsSecret, _ := client.StringValueDetails(ctx,
		"aws.secrets.GITHUB_ACCESS_TOKEN",
		"default-github-token", openfeature.EvaluationContext{},
	)
	fmt.Println("GitHub Access Token (AWS Secrets Manager):", awsSecret.Value)

	// Fetch parameter from AWS Parameter Store via Pulumi ESC
	awsParam, _ := client.StringValueDetails(ctx,
		"aws.params.GOOGLE_API_KEY",
		"default-google-api-key", openfeature.EvaluationContext{},
	)
	fmt.Println("Google API Key (AWS Parameter Store):", awsParam.Value)

	// Fetch string config value from Pulumi ESC
	dbURL, _ := client.StringValueDetails(ctx,
		"configs.USERS_DB_MONGO_URL",
		"mongodb://localhost:27017", openfeature.EvaluationContext{},
	)
	fmt.Println("Users DB URL (Pulumi ESC config value):", dbURL.Value)

	// Fetch integer config value from Pulumi ESC
	maxConnections, _ := client.IntValueDetails(ctx,
		"configs.MAX_CONNECTIONS",
		50, openfeature.EvaluationContext{},
	)
	fmt.Println("Max Connections (Pulumi ESC config value):", maxConnections.Value)

	// Fetch boolean config value from Pulumi ESC
	debugMode, _ := client.BooleanValueDetails(ctx,
		"configs.DEBUG_MODE",
		false, openfeature.EvaluationContext{},
	)
	fmt.Println("Debug Mode Enabled (Pulumi ESC config value):", debugMode.Value)

	// Fetch float config value from Pulumi ESC
	cpuThreshold, _ := client.FloatValueDetails(ctx,
		"configs.CPU_THRESHOLD",
		0.75, openfeature.EvaluationContext{},
	)
	fmt.Println("CPU Threshold (Pulumi ESC config value):", cpuThreshold.Value)

	// Fetch secret value encrypted within Pulumi ESC
	secretValue, _ := client.StringValueDetails(ctx,
		"configs.OPENAI_API_KEY",
		"sk-12345", openfeature.EvaluationContext{},
	)
	fmt.Println("Encrypted Secret (Pulumi ESC secret value):", secretValue.Value)
	isSecret, _ := secretValue.FlagMetadata.GetBool("secret")
	fmt.Printf("Is secret value: %t\n", isSecret)

	<-sigChan
	fmt.Println("Stopping application")
}
