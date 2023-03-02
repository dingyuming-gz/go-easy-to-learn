/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	age     int
	like    []string
	address string
)

var hiCmd = &cobra.Command{
	Use:   "hi",
	Short: "sss",
	Long:  "lll",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
		Print()
	},
}

func init() {
	rootCmd.AddCommand(hiCmd)
	hiCmd.Flags().StringVarP(&name, "name", "", "world", "")
	hiCmd.Flags().IntVarP(&age, "age", "", 1, "")
	hiCmd.Flags().StringSliceVarP(&like, "like", "", []string{}, "")
	hiCmd.Flags().StringVarP(&address, "address", "", "Beijing", "")
}

func Print() {
	fmt.Println("name:", name)
	fmt.Println("age:", age)
	fmt.Println("like:", like)
	fmt.Println("address:", address)
}
