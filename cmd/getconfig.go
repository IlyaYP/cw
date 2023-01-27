/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/IlyaYP/cw/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// getconfigCmd represents the getconfig command
var getconfigCmd = &cobra.Command{
	Use:   "getconfig",
	Short: "Gets running config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("getconfig called")
	// },
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("getconfig called")
		logonname := viper.GetString("LOGONNAME")
		pw := viper.GetString("PW")
		// fmt.Println(hosts, logonname, pw)
		if logonname == "" {
			return fmt.Errorf("no logon name")
		}
		if pw == "" {
			return fmt.Errorf("no password")
		}
		if len(HOSTS) == 0 {
			return fmt.Errorf("no hosts")
		}
		return getconfig(logonname, pw, HOSTS)
	},
}

func init() {
	rootCmd.AddCommand(getconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getconfig(logonname, pw string, hosts []string) error {

	// send the commands
	commands := []string{
		"terminal length 0",
		"show running-config",
		"exit",
	}

	g := new(errgroup.Group)
	for _, hostname := range hosts {
		hostname := hostname // https://go.dev/doc/faq#closures_and_goroutines
		g.Go(func() error {
			return pkg.GetConfig(hostname, logonname, pw, commands)
		})
	}

	log.Print("doing")
	err := g.Wait()
	if err != nil {
		log.Print(err)
		return err
	}

	log.Print("all done no errors")
	return nil
}
