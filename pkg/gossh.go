package pkg

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"

	"golang.org/x/crypto/ssh"
)

func GetConfig(hostname, username, password string, commands []string) error {
	var reConf = regexp.MustCompile(`(?s)Current configuration .*end`)
	var reHost = regexp.MustCompile(`(?m)^hostname\s([-0-9A-Za-z_]+).?$`)

	out, err := doSSHCommands(hostname, username, password, commands)
	if err != nil {
		return err
	}

	if reHost.Match(out) {
		fname := string(reHost.FindSubmatch(out)[1])

		log.Print(fname)

		if reConf.Match(out) {
			config := reConf.FindAll(out, -1)[0]
			err := os.WriteFile(fname, config, 0644)
			if err != nil {
				log.Print(err)
				return err
			}
		} else {
			log.Print(hostname, "config not found")
		}
	} else {
		log.Print(hostname, "hostname not found")
	}
	// time.Sleep(5 * time.Second)
	return nil
}

func doSSHCommands(hostname, username, password string, commands []string) ([]byte, error) {
	port := "22"

	// SSH client config
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	config.KeyExchanges = append(
		config.KeyExchanges,
		"diffie-hellman-group-exchange-sha256",
		"diffie-hellman-group-exchange-sha1",
		"diffie-hellman-group1-sha1",
	)

	config.Ciphers = append(config.Ciphers, "aes128-cbc", "3des-cbc",
		"aes192-cbc", "aes256-cbc", "aes128-ctr", "aes192-ctr", "aes256-ctr")

	//////////////////////////////
	// Connect to host
	client, err := ssh.Dial("tcp", hostname+":"+port, config)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer client.Close()

	// Create sesssion
	sess, err := client.NewSession()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer sess.Close()

	// StdinPipe for commands
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Uncomment to store output in variable
	var b bytes.Buffer
	sess.Stdout = &b
	//sess.Stderr = &b

	// Enable system stdout
	// Comment these if you uncomment to store in variable
	// sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// Start remote shell
	err = sess.Shell()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	}

	// log.Print(hostname)
	// Wait for sess to finish
	err = sess.Wait()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Uncomment to store in variable
	// fmt.Println(b.String())
	out := b.Bytes()
	// fmt.Println(string(out))

	return out, nil

}

func GetUsers(hostname, username, password string, commands []string, ch chan string) error {
	defer close(ch)
	// var reConf = regexp.MustCompile(`(?s)Current configuration .*end`)
	// var reHost = regexp.MustCompile(`(?m)^hostname\s([-0-9A-Za-z_]+).?$`)
	var reUser = regexp.MustCompile(`(?m)^username\s([-0-9A-Za-z_]+)\s`)

	out, err := doSSHCommands(hostname, username, password, commands)
	if err != nil {
		return err
	}

	// if reHost.Match(out) {
	// 	fname := string(reHost.FindSubmatch(out)[1])

	// 	log.Print(fname)

	// 	if reConf.Match(out) {
	// 		config := reConf.FindAll(out, -1)[0]
	// 		err := os.WriteFile(fname, config, 0644)
	// 		if err != nil {
	// 			log.Print(err)
	// 			return err
	// 		}
	// 	} else {
	// 		log.Print(hostname, "config not found")
	// 	}
	// } else {
	// 	log.Print(hostname, "hostname not found")
	// }

	err = os.WriteFile(hostname+".log", out, 0644)
	if err != nil {
		log.Print(err)
		return err
	}

	if reUser.Match(out) {
		for _, sm := range reUser.FindAllSubmatch(out, -1) {
			uname := string(sm[1])
			// ch <- hostname + ":" + uname
			ch <- uname
		}
	} else {
		log.Print(hostname, "username not found")
	}

	return nil
}
