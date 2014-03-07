package main

import (
  "fmt"
  "os"
  "log"
  "math/big"
  "strings"
  "strconv"
  "crypto/rand"
  "github.com/nsf/termbox-go"
)

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
  {
    Run:        test,
    UsageLine:  "test [name]",
    Short:      "starts an interactive test",
    Long:       `Starts a new interactive test.

    Available tests:
      conj`,
  },
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
  fmt.Printf("usage: %s\n", c.UsageLine)
  fmt.Println(strings.TrimSpace(c.Long))
}

func version(cmd *command, args []string) {
  fmt.Println("msyu version", VERSION)
  fmt.Println("Copyright (C) 2014 Cristian Kubis")
}

/* TODO: list multiple results and let the user choose */
func conj(cmd *command, args []string) {
  var word *word = nil

  if len(args) < 1 {
    word = DB_get_random_verbs(1)[0]
  } else {
    arg := args[0]
    
    if isJapaneseString(arg) {
      word = DB_search_word(arg, JAP, VERB)
    } else if isLatin(arg) {
      word = DB_search_word(arg, EN, VERB)
    }
  }
 
  if word == nil {
    log.Fatal("Could not find word '%s'\n", args[0])
  }
  word.PrintConjTable()
}

func test(cmd *command, args []string) {
  if len(args) < 1 {
    cmd.Usage()
    os.Exit(2)
  }

  n := -1
  if len(args) > 1 {
    n, _ = strconv.Atoi(args[1])
  }
  
  if n < 0 {
    n = 5
  }

  switch(args[0]) {
  case "conj":
    test_conj(n)
  }
}

func test_conj(n int) {
  words := DB_get_random_verbs(n)

  if(words == nil) {
    panic("no verbs found")
  }

  for _, word := range words {
    var sPolite string
    var sPositive string

    n, _ := rand.Int(rand.Reader, big.NewInt(2/*int64(len(conjunctions))*/))
    conj := conjunctions[n.Int64()]
    polite := n.Int64() % 2 == 0
    positive := n.Int64() % 3 == 0
    
    if polite {
      sPolite = "Polite"
    } else {
      sPolite = "Plain"
    }

    if positive {
      sPositive = "Positiv"
    } else {
      sPositive = "Negative"
    }

    const coldef = termbox.ColorDefault
    kana, kanji := conj.Exec(word, positive, polite)
    tbprint(1, 1, coldef, coldef, fmt.Sprintf("%s - %s / %s", conj.Name, sPositive, sPolite))
    tbprint(1, 3, coldef, coldef, fmt.Sprintf("%s (%s)", word.kana, strings.Join(word.kanji, ", ")))
    tbprint(1, 6, coldef, coldef, "Answer: ")
    termbox.SetCursor(1+len("Answer: "), 6)
    termbox.Flush()

    entry := getString(1+len("Answer: "), 6, coldef, coldef)

    termbox.Flush()
    termbox.Clear(coldef, coldef)

    if entry == kana || entry == kanji {
      tbprint(1, 1, coldef, coldef, "Correct")
      tbprint(1, 3, coldef, coldef, fmt.Sprintf("%s - %s / %s", conj.Name, sPositive, sPolite))
      tbprint(1, 5, coldef, coldef, fmt.Sprintf("%s (%s)", kana, kanji))
    } else {
      tbprint(1, 1, coldef, coldef, "Wrong!")
      tbprint(1, 3, coldef, coldef, fmt.Sprintf("%s (%s)", kana, kanji))
      tbprint(1, 5, coldef, coldef, fmt.Sprintf("%s %s %s", conj.Name, sPositive, sPolite))
    }

    tbprint(1, 8, coldef, coldef, fmt.Sprint(conj.Rule))
    termbox.Flush()

    getString(-100, -100, coldef, coldef)
    termbox.Clear(coldef, coldef)
  }
}
