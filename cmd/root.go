package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:   "gaterun",
		Short: "gaterun is an api gateway",
		Long: `
			this is an api gateway which lets you map services and routes to multiple microservices 
			and have full control over security and access management concepts.	
			gaterun offer rate-limiting, authentication, load-balancing, monitoring.
		`,
	}

	cmd.AddCommand(NewStartGatewayCmd(ctx))

	return cmd
}
