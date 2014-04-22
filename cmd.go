package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var commands = []command{
	{
		Run:       version,
		UsageLine: "version",
		Short:     "prints msyu version",
		Long:      `Prints the currently version of msyu.`,
	},
	{
		Run:       conj,
		UsageLine: "conj [word]",
		Short:     "prints conjugation table",
		Long:      `Prints the conjugation table of a given word. Uses a random verb instead if no word is supplied.`,
	},
	{
		Run:       test,
		UsageLine: `test [name] [n]`,
		Short:     "starts an interactive test with n items asked",
		Long: `Starts a new interactive test with n items asked.

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
		n = 25
	}

	switch args[0] {
	case "conj":
		test_conj(n)
	}
}

func test_conj(n int) {
	words := DB_get_random_verbs(n)

	if words == nil {
		panic("no verbs found")
	}

	for _, word := range words {
		var sPolite string
		var sPositive string

		randomBytes := make([]byte, 3)
		_, err := rand.Read(randomBytes)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		positive := (int(randomBytes[0]) % 2) == 0
		polite := (int(randomBytes[1]) % 2) == 0
		conj := conjugations[int(int(randomBytes[2])%len(conjugations))]

		if polite {
			sPolite = "Polite"
		} else {
			sPolite = "Plain"
		}

		if positive {
			sPositive = "Positive"
		} else {
			sPositive = "Negative"
		}

		clear()

		kana, kanji := conj.Exec(word, positive, polite)
		fmt.Printf("%s - %s / %s\n\n", conj.Name, sPositive, sPolite)
		fmt.Printf("%s (%s)\n\n", word.kana, strings.Join(word.kanji, ", "))
		fmt.Printf("Answer: ")

		var input string
		fmt.Scanf("%s", &input)

		clear()

		if input == kana || input == kanji {
			fmt.Printf("Correct Answer !\n\n")
			fmt.Printf("%s - %s / %s\n\n", conj.Name, sPositive, sPolite)
			fmt.Printf("%s (%s)\n", kana, kanji)
		} else {
			fmt.Printf("Wrong Answer !\n\n")
			fmt.Printf("%s - %s / %s\n\n", conj.Name, sPositive, sPolite)
			fmt.Printf("Entered: %s\n", input)
			fmt.Printf("Correct: %s (%s)\n\n", kana, kanji)
			fmt.Println("Conjugation Rules:")

			// this will make sense in the future, I promise
			if strings.HasPrefix(word.gloss[0].pos[0], "v1") {
				fmt.Printf("%s\n", conj.Rule["v1"])
			} else if strings.HasPrefix(word.gloss[0].pos[0], "v5") {
				fmt.Printf("%s\n", conj.Rule["v5"])
			}
			fmt.Println("\nBase Rules:\n", baseRules)
		}

		fmt.Printf("\n<Enter> -> Next")
		fmt.Scanf("%s", &input)
		clear()
	}
}
