package service

import (
	"github.com/spf13/cobra"
)

// NewServiceCmd creates the service command
func NewServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Aliases: []string{"services", "svc"},
		Short: "Manage microservices",
		Long:  "Create, list, update, and delete managed microservices (Redis, RabbitMQ, PostgreSQL, etc.)",
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewDescribeCmd())
	cmd.AddCommand(NewScaleCmd())
	cmd.AddCommand(NewLogsCmd())

	return cmd
}
