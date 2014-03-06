package main

import (
  "fmt"
  "os"
  "log"
  _ "reflect"
  "math/big"
  "strings"
  "strconv"
  "crypto/rand"
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
    
    if isJapanese(arg) {
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
    n = 25
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

    n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(conjunctions))))
    conj := conjunctions[n.Int64()]
    polite := n.Int64() % 2 == 0
    positive := n.Int64() % 3 == 0
    
    if polite {
      sPolite = "polite"
    } else {
      sPolite = "plain"
    }

    if positive {
      sPositive = "positiv"
    } else {
      sPositive = "negative"
    }

    kana, kanji := conj.Exec(word, positive, polite)
    word.Print()
    fmt.Printf("\n%s   %s   %s\n\n", conj.Name, sPositive, sPolite)

    entry := ""
    fmt.Scanf("%s", &entry)

    if entry == kana || entry == kanji {
      fmt.Println("correct!")
    } else {
      fmt.Println(kanji)
    }
  }
}
