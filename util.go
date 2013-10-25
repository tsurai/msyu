package main

import (
  "fmt"
  "strings"
  "unicode"
)

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

func isJapanese(s string) bool {
  for _, r := range []rune(s) {
    if !unicode.Is(unicode.Hiragana, r) && !unicode.Is(unicode.Katakana, r) && !(r >= '\u4E00' && r <= '\u9FA5') {
      return false
    }
  }
  return true
}

func (w *Word) Print() {
  if w.kanji != "" {
    fmt.Printf("%s(%s)(%s) - %s\n", w.kanji, strings.Join(w.typ, ", "), w.kana, w.meaning)
  } else {
    fmt.Printf("%s(%s) - %s\n", w.kana, strings.Join(w.typ, ", "), w.meaning)
  }
}
