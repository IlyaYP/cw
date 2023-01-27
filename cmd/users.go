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

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("users called")
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("users called")
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
		return getusers(logonname, pw, HOSTS)
	},
}

func init() {
	showCmd.AddCommand(usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getusers(logonname, pw string, hosts []string) error {

	// send the commands
	commands := []string{
		"terminal length 0",
		"show running-config | include username",
		"exit",
	}

	n := len(hosts)
	chs := make([]chan string, 0, n)
	for i := 0; i < n; i++ {
		ch := make(chan string)
		chs = append(chs, ch)
	}

	g := new(errgroup.Group)
	for i, hostname := range hosts {
		hostname := hostname // https://go.dev/doc/faq#closures_and_goroutines
		ch := chs[i]
		g.Go(func() error {
			return pkg.GetUsers(hostname, logonname, pw, commands, ch)
		})
	}

	log.Print("doing")

	// здесь fanIn
	for v := range pkg.FanIn(chs...) {
		fmt.Println(v)
	}

	log.Print("done")

	err := g.Wait()
	if err != nil {
		log.Print(err)
		return err
	}

	log.Print("all done no errors")
	return nil
}
