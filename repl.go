package main

import (
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/tcnksm/go-input"
	"github.com/urfave/cli"
)

// PDIRepl type
type PDIRepl struct {
	prefix  string
	logined bool
	client  *PDIClient
	ui      *input.UI
}

func (r *PDIRepl) completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	s := []prompt.Suggest{}
	if w == "" {
		return []prompt.Suggest{}
	}
	if r.logined {
		s = []prompt.Suggest{}
	} else {
		s = []prompt.Suggest{
			{Text: "login", Description: "login to netweaver"},
		}
	}

	return prompt.FilterHasPrefix(s, w, true)
}

// RunRepl logic
func (r *PDIRepl) RunRepl() {
	for {
		in := strings.TrimSpace(prompt.Input(r.prefix, r.completer))

		switch in {
		case "login":
			r.client.username, _ = r.ui.Ask("username", &input.Options{
				Default:  "",
				Required: true,
				Loop:     true,
			})
			r.client.password, _ = r.ui.Ask("password", &input.Options{
				Default:  "",
				Required: true,
				Loop:     true,
				Mask:     true,
			})
			r.client.hostname, _ = r.ui.Ask("hostname", &input.Options{
				Default:  "",
				Required: true,
				Loop:     true,
			})
			r.client.login()
		}
	}
}

// NewPDIRepl instance
func NewPDIRepl() *PDIRepl {
	return &PDIRepl{
		prefix:  "> ",
		logined: false,
		client:  &PDIClient{exitCode: 0},
		ui: &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		},
	}
}

var commandRepl = cli.Command{
	Name:  "repl",
	Usage: "start repl",
	Action: func(ctx *cli.Context) {
		NewPDIRepl().RunRepl()
	},
}
