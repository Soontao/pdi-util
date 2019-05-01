package main

import "github.com/urfave/cli"

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "username, u",
		EnvVar: "PDI_USER",
		Usage:  "The PDI Development User",
	},
	cli.StringFlag{
		Name:   "password, p",
		EnvVar: "PDI_PASSWORD",
		Usage:  "The PDI Development User Password",
	},
	cli.StringFlag{
		Name:   "hostname, t",
		EnvVar: "PDI_TENANT_HOST",
		Usage:  "The PDI Tenant host",
	},
	cli.StringFlag{
		Name:   "release,r",
		EnvVar: "TENANT_RELEASE",
		Value:  "1905",
		Usage:  "The tenant release version, e.g. 1905",
	},
}
