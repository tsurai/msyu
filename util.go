package main

import (
  "fmt"
  "strings"
  "unicode"
  "github.com/nsf/termbox-go"
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

func (w *word) Print() {
  if(w.kanji[0] != "") {
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

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) int {
  offset := 0

  for _, r := range msg {
    termbox.SetCell(x + offset, y, r, fg, bg)
    if isJapanese(r) {
      offset += 2
    } else {
      offset++
    }
  }
  return offset
}

func getString(x, y int, fg, bg termbox.Attribute) string {
  var ret []rune
  ret = nil

loop:
  for {
    switch ev := termbox.PollEvent(); ev.Type {
    case termbox.EventKey:
      switch ev.Key {
      case termbox.KeyEnter:
        termbox.Flush()
        break loop
      case termbox.KeyBackspace:
        fallthrough
      case termbox.KeyBackspace2:
        if len(ret) > 0 {
          tbprint(x, y, fg, bg, "                                                                  ")
          termbox.Flush()

          ret = ret[:len(ret)-1]
          written := tbprint(x, y, fg, bg, string(ret))
          termbox.SetCursor(x+written, y)
          termbox.Flush()
        }
      default:
        ret = append(ret, ev.Ch)
        written := tbprint(x, y, fg, bg, string(ret))
        termbox.SetCursor(x+written, y)
        termbox.Flush()
      }
    }
  }
  termbox.HideCursor()
  return string(ret)
}