package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/imroc/req"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

// PDIAction wrapper
func PDIAction(action func(pdiClient *pdiutil.PDIClient, c *cli.Context)) func(c *cli.Context) error {
	return func(c *cli.Context) error {

		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				log.Println("FATAL error happened, so terminated")
				os.Exit(1)
			}
		}()

		// set long timeout
		req.SetTimeout(time.Hour * 12)

		// set long tls handshake timeout
		req.Client().Transport.(*http.Transport).TLSHandshakeTimeout = (time.Minute * 30)

		// overwrite here
		username := c.GlobalString("username")
		password := c.GlobalString("password")
		hostname := c.GlobalString("hostname")
		release := c.GlobalString("release")

		hostname = strings.TrimPrefix(hostname, "https://") // remove hostname schema
		hostname = strings.TrimPrefix(hostname, "http://")  // remove hostname schema

		pdiClient, err := pdiutil.NewPDIClient(username, password, hostname, release)

		if err != nil {
			return err
		}

		log.Printf("Login to %v as %v", hostname, username)

		action(pdiClient, c) // do process

		if pdiClient.GetExitCode() > 0 {
			// if error happened, change exit code
			return fmt.Errorf("finished %s with error", c.Command.FullName())
		}
		return nil
	}
}
