package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func (w *word) Print() {
	if w.kanji[0] != "" {
		fmt.Printf("%s (%s)\n", w.kana, strings.Join(w.kanji, ", "))
	} else {
		fmt.Printf("%s\n", w.kana)
	}

	for _, g := range w.gloss {
		if g.pos[0] != "" {
			fmt.Printf("    t: %s\n", strings.Join(g.pos, ", "))
		}
		fmt.Printf("        * %s \n", strings.Join(g.meaning, ", "))
	}
}

func isLatin(s string) bool {
	runes := make([]rune, len(s))
	copy(runes, []rune(s))
	for _, r := range runes {
		if !unicode.Is(unicode.Latin, r) && !unicode.Is(unicode.White_Space, r) {
			return false
		}
	}
	return true
}

func isJapanese(r rune) bool {
	return unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) || (r >= '\u4E00' && r <= '\u9FA5')
}

func isJapaneseString(s string) bool {
	for _, r := range []rune(s) {
		if !isJapanese(r) {
			return false
		}
	}
	return true
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	cmd = exec.Command("cmd", "/c", "cls & clear")
	cmd.Run()
}
