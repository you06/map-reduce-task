package main

import (
	"github.com/spf13/cobra"
	"github.com/you06/map-reduce-task/generator"
)

var (
	rootCmd = &cobra.Command{
		Use:   "task",
		Short: "Map reduce task entrance",
	}

	generatorCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate test data",
		Run: func(cmd *cobra.Command, args []string) {
			generator.Run(generatorSize, generatorOutput, generatorPartition)
		},
	}
	generatorSize      string
	generatorOutput    string
	generatorPartition int
)

func init() {
	generatorCmd.Flags().StringVarP(&generatorSize, "size", "s", "100M", "Total size of generated files")
	generatorCmd.Flags().StringVarP(&generatorOutput, "output", "o", "data", "Output path")
	generatorCmd.Flags().IntVarP(&generatorPartition, "partition", "p", 10, "Partition amount of the generated files")
	rootCmd.AddCommand(generatorCmd)
}

func main() {
	rootCmd.Execute()
}
