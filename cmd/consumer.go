package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uminac/go-pb-stuff/internal/consumer"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "consumes messages",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		c := consumer.NewMQTTConsumer()

		if err := c.Run(); err != nil {
			fmt.Printf("error: consumer ended (%s)\n", err.Error())
		} else {
			fmt.Println("consumer ended.")
		}
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)
}
