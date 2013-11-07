package main

import (
  "fmt"
  "os"
  "flag"
  "strings"
  "html/template"
)

const (
  VERSION = "0.1.7a"
)

type command struct {
  Run func(cmd *command, args []string)
  UsageLine string
  Short string
  Long string
}

type Gloss struct {
  pos []string
  meaning []string
}

type Word struct {
  kana string
  kanji []string
  gloss map[int]*Gloss
}

var usageTemplate = `msyu is a japanese learning tool.

Usage:

        mysu <command> [arguments]

The commands are:
{{range .}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "msyu help [command]" for more information about a command.
`

var helpTemplate = `usage: msyu {{.UsageLine}}

{{.Long}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(text string, data interface{}) {
  t := template.New("top")
  template.Must(t.Parse(text))
  if err := t.Execute(os.Stderr, data); err != nil {
    panic(err)
  }
}

func (c *command) Name() string {
  name := c.UsageLine
  i := strings.Index(name, " ")

  if i >= 0 {
    name = name[:i]
  }

  return name
}

func (c *command) Usage() {
  fmt.Println("usage: %s\n", c.UsageLine)
  fmt.Println("%s", strings.TrimSpace(c.Long))
}

func help(args []string) {
  if len(args) == 0 {
    tmpl(usageTemplate, commands)
    return
  }

  arg := args[0]

  for _, cmd := range commands {
    if cmd.Name() == arg {
      tmpl(helpTemplate, cmd)
      return
    }
  }

  fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'msyu help' for a list of valid commands.\n", arg)
}

func main() {
  flag.Parse()
  args := flag.Args()

  if len(args) < 1 {
    tmpl(usageTemplate, commands)
    os.Exit(1)
  }

  if args[0] == "help" {
    help(args[1:])
    os.Exit(2)
  }

  DB_init()
  for _, cmd := range commands {
    if args[0] == cmd.Name() {
      cmd.Run(&cmd, args[1:])
    }
  }
  DB_close()

  return
}
