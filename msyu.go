package main

import (
  "fmt"
  "os"
  "flag"
  "strings"
  "html/template"
)

const (
  VERSION = "0.1.5a"
)

type command struct {
  Run func(cmd *command, args []string)
  UsageLine string
  Short string
  Long string
}

type Word struct {
  kana string
  kanji string
  typ []string
  meaning string
}

var commands = []command {
  { 
    Run:        version,
    UsageLine:  "version",
    Short:      "prints msyu version",
    Long:       `Prints the currently version of msyu.`,
  },
  {
    Run:        conj,
    UsageLine:  "conj [word]",
    Short:      "prints conjunction table",
    Long:       `Prints the conjunction table of a given word. Uses a random verb instead if no word is supplied.`,
  },
/*  {
    Run:        test,
    UsageLine:  "test [name]",
    Short:      "starts a test",
    Long:       `Starts a new interactive test.`,
  },*/
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

func version(cmd *command, args []string) {
  fmt.Println("msyu version", VERSION)
  fmt.Println("Copyright (C) 2013 Cristian Kubis")
}

/* TODO: list multiple results and let the user choose */
func conj(cmd *command, args []string) {
  var word *Word = nil

  if(len(args) < 1) {
    word = DB_get_random_verbs(1)[0]
  } else {
    arg := args[0]
    
    if isJapanese(arg) {
      word = DB_search_in_jap(arg)
    } else if isLatin(arg) {
      word = DB_search_in_eng(arg)
    }
  }
 
  if(word == nil) {
    fmt.Fprintf(os.Stderr, "Could not find word '%s'\n", args[0])
    os.Exit(2)
  }
  word.PrintConjTable()
}

/*
func test(cmd *command, args []string) {
  db, err := sql.Open("sqlite3", "JMdict.db")
  if err != nil {
    fmt.Fprintln(os.Stderr, "An error has occured:", err)
    os.Exit(1)
  }
  defer db.Close()
 
  query := "SELECT r_ele.value, group_concat(DISTINCT entity.entity), group_concat(gloss.value, ','), k_ele.value FROM r_ele "+
      "LEFT JOIN sense ON r_ele.fk=sense.fk "+
      "LEFT JOIN pos ON sense.id=pos.fk "+
      "LEFT JOIN entity ON pos.entity=entity.id "+
      "LEFT JOIN gloss ON sense.id=gloss.fk "+
      "LEFT JOIN k_ele ON r_ele.fk = k_ele.fk "+
      "WHERE entity.entity LIKE '%iv%' OR " +
      "entity.entity LIKE '%v1%' OR entity.entity LIKE '%vz%' OR " +
      "entity.entity LIKE '%vi%' OR entity.entity LIKE '%vk%' OR " +
      "entity.entity LIKE '%vn%' OR entity.entity LIKE '%v5%' OR " +
      "entity.entity LIKE '%vr%' " +
      "GROUP BY r_ele.id ORDER BY RANDOM() LIMIT 100"

  rows, err := db.Query(query)
  if err != nil {
    fmt.Fprintln(os.Stderr, "An error has occured:", err)
    os.Exit(1)
  }
  defer rows.Close()

  i := 0
  
  var rvalue string
  var kvalue string
  var pos string
  var gloss string
 
  for rows.Next() {
    i++
    rows.Scan(&rvalue, &pos, &gloss, &kvalue)

    w := &Word{rvalue, kvalue, strings.Split(pos, ","), gloss}
    w.Print()
  }
}
*/
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
