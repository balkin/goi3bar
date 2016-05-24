package command

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"os/exec"
	"bytes"
	"fmt"
	"log"
	"strings"
	"io/ioutil"
)

type Command struct {
	Label        string `json:"label"`
	Instance     string `json:"instance"`
	Command      string `json:"command"`
	Color        string `json:"color"`
	ClickCommand string `json:"onclick"`

	// Identifier for receiving click events
	Name         string
}

func (g Command) Generate() ([]i3.Output, error) {
	items := make([]i3.Output, 1)

	items[0].Name = g.Name
	cmd := exec.Command(g.Command)
	if g.Instance != "" {
		cmd.Env = []string{fmt.Sprintf("BLOCK_INSTANCE=%s", g.Instance)}
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Panicf("Failed to execute %s: %v", g.Command, err)
		items[0].FullText = "ERROR"
		items[0].Color = i3.DefaultColors.Crit
	} else {
		if g.Color == "" {
			items[0].Color = i3.DefaultColors.General
		} else {
			items[0].Color = g.Color
		}
		text := strings.TrimRight(out.String(), "\n\r")
		if g.Label == "" {
			items[0].FullText = fmt.Sprintf("%s %s", g.Label, text)
		} else {
			items[0].FullText = text
		}
	}
	items[0].Instance = g.Command
	items[0].Separator = true
	return items, nil
}

func (c *Command) Click(e i3.ClickEvent) error {
	if c.ClickCommand == "" {
		return nil
	}
	splitArgs := strings.Split(c.ClickCommand, " ")
	cmd := exec.Command(splitArgs[0], splitArgs[1:]...)
	cmd.Stdout = ioutil.Discard
	log.Println()
	err := cmd.Run()
	return err
}
