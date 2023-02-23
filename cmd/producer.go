package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uminac/go-pb-stuff/internal/producer"
)

var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "produces messages",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		p := producer.NewMQTTProducer()

		if err := p.Run(); err != nil {
			fmt.Printf("error: producer ended (%s)\n", err.Error())
		} else {
			fmt.Println("producer ended.")
		}
	},
}

func init() {
	rootCmd.AddCommand(producerCmd)
}
