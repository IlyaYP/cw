/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/IlyaYP/cw/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// shUsersCmd represents the users command
var shUsersCmd = &cobra.Command{
	Use:   "shusers",
	Short: "Shows users local account on devices",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("users called")
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("shusers called")
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
	rootCmd.AddCommand(shUsersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type UserName struct {
	Key   string
	Value int
}

type UserNameList []UserName

func (p UserNameList) Len() int           { return len(p) }
func (p UserNameList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p UserNameList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func getusers(logonname, pw string, hosts []string) error {
	m := make(map[string]int)

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

	// log.Print("doing")

	// здесь fanIn
	for v := range pkg.FanIn(chs...) {
		// fmt.Println(v)
		m[v] += 1
	}

	p := make(UserNameList, len(m))

	i := 0
	for k, v := range m {
		p[i] = UserName{k, v}
		i++
	}

	sort.Sort(sort.Reverse(p))

	for _, k := range p {
		fmt.Printf("%s: %v\n", k.Key, k.Value)
	}

	// log.Print("done")

	err := g.Wait()
	if err != nil {
		// log.Print(err)
		return err
	}

	// log.Print("all done no errors")
	return nil
}
