/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var scrFile string

// scriptCmd represents the script command
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("script called")
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("script called")
		logonname := viper.GetString("LOGONNAME")
		pw := viper.GetString("PW")
		// fmt.Println(hosts, logonname, pw)
		if logonname == "" {
			return fmt.Errorf("no logon name")
		}
		if pw == "" {
			return fmt.Errorf("no password")
		}
		if scrFile == "" {
			return fmt.Errorf("no commands file")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scriptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scriptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scriptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scriptCmd.Flags().StringVar(&scrFile, "cmd", "", "script commands file")

}
