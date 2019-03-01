package main

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	input "github.com/tcnksm/go-input"
	"github.com/urfave/cli"
)

// PDIRepl type
type PDIRepl struct {
	prefix   string
	logined  bool
	client   *PDIClient
	ui       *input.UI
	commands []*ReplCommand
}

// ReplCommand command
type ReplCommand struct {
	// command name
	name string
	// description text
	description string
	// paramNames
	paramNames []string
	// command action
	action func(p *PDIRepl, params ...string)
	// can suggest
	canSuggest func(p *PDIRepl) bool
}

var replCommandLogin = &ReplCommand{
	name:        "login",
	description: "login to netweaver",
	paramNames:  []string{"username", "password", "hostname"},
	action: func(p *PDIRepl, params ...string) {
		p.client.username = params[0]
		p.client.password = params[1]
		p.client.hostname = params[2]
		p.client.login()
		p.logined = true
	},
	canSuggest: func(p *PDIRepl) bool {
		return !p.logined
	},
}

var replCommandLogout = &ReplCommand{
	name:        "logout",
	description: "logout system",
	paramNames:  []string{},
	action: func(p *PDIRepl, params ...string) {
		p.logined = false
	},
	canSuggest: func(p *PDIRepl) bool {
		return p.logined
	},
}

func (r *PDIRepl) completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	s := []prompt.Suggest{}
	for _, c := range r.commands {
		if c.canSuggest(r) {
			s = append(s, prompt.Suggest{
				Text:        c.name,
				Description: c.description,
			})
		}
	}
	return prompt.FilterHasPrefix(s, w, true)
}

// RunRepl logic
func (r *PDIRepl) RunRepl() {
	for {
		in := strings.TrimSpace(prompt.Input(r.prefix, r.completer))

		foundCommand := false
		for _, c := range r.commands {
			if c.name == in {
				foundCommand = true
				params := []string{}
				for _, paramName := range c.paramNames {
					param, _ := r.ui.Ask(fmt.Sprintf("input %s:", paramName), &input.Options{
						Default:  "",
						Required: true,
						Loop:     true,
					})
					params = append(params, param)
				}
				c.action(r, params...)
				break
			}
		}

		if !foundCommand {
			fmt.Printf("not found command: %s", in)
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
		commands: []*ReplCommand{
			replCommandLogin,
			replCommandLogout,
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
