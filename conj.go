package main

import (
	"unicode/utf8"
	"fmt"
)

// conjunction ------------
func (w *Word) ToPresent(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string
  
  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ます"
    } else {
      return w.ToRentaikei()
    }
  } else {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ません"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "ない"
    }
  }
  return kana + ending, kanji + ending
}

func (w *Word) ToPast(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ました"
    } else {
      if w.typ[0] == "v1" {
        kana, kanji = w.ToRenyoukei()
        ending = "た"
      } else {
        kana, kanji = w.ToStem()
        switch(w.kana[len(kana):]) {
          case "く":
            ending = "いた"
          case "ぐ":
            ending = "いだ"
          case "ぬ":
            fallthrough
          case "ぶ":
            fallthrough
          case "む":
            ending = "んだ"
          case "う":
            fallthrough
          case "つ":
            fallthrough
          case "る":
            ending = "った"
          default:
            kana, kanji = w.ToRenyoukei()
            ending = "た"
        }
      }
    }
  } else {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "ませんでした"
    } else {
      kana, kanji = w.ToRenyoukei()
      ending = "なかった" 
    }
  }
  return kana + ending, kanji + ending
}

func (w *Word) ToTeForm(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string  
  
  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "まして"
    } else {
      if w.typ[0] == "v1" {
        kana, kanji = w.ToRenyoukei()
        ending = "て"
      } else {
        kana, kanji = w.ToStem()
        switch(w.kana[len(kana):]) {
          case "く":
            ending = "いて"
          case "ぐ":
            ending = "いで"
          case "ぬ":
            fallthrough
          case "ぶ":
            fallthrough
          case "む":
            ending = "んで"
          case "う":
            fallthrough
          case "つ":
            fallthrough
          case "る":
            ending = "って"
          default:
            kana, kanji = w.ToRenyoukei()
            ending = "て"
        }
      }
    }
  } else {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ませんで" 
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "ないで"
    }
  }
  return kana + ending, kanji + ending
}

func (w *Word) ToConditional(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string
  
  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ましたら"
    } else {
      if w.typ[0] == "v1" {
        kana, kanji = w.ToRenyoukei()
        ending = "たら"
      } else {
        kana, kanji = w.ToStem()
        switch(w.kana[len(kana):]) {
          case "く":
            ending = "いたら"
          case "ぐ":
            ending = "いだら"
          case "ぬ":
            fallthrough
          case "ぶ":
            fallthrough
          case "む":
            ending = "んだら"
          case "う":
            fallthrough
          case "つ":
            fallthrough
          case "る":
            ending = "ったら"
          default:
            kana, kanji = w.ToRenyoukei()
            ending = "たら"
        }
      }
    }
  } else {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ませんでしたら" 
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "なかったら"
    }
  }
  return kana + ending, kanji + ending
}

// Inflection bases -------------
func (w *Word) ToStem() (string, string) {
  _, size := utf8.DecodeLastRuneInString(w.kana)
  return w.kana[:len(w.kana)-size], w.kanji[:len(w.kanji)-size]
}

func (w *Word) ToMizenkei() (string, string) {
  stem, kstem := w.ToStem()
  if w.typ[0] == "v1" {
    return stem, kstem
  } else if w.typ[0] == "v5aru" || w.typ[0] == "v5k-s" || w.typ[0] == "v5r-i" ||
       w.typ[0] == "v5u-s" || w.typ[0] == "v5uru" {
    return "?", "?"
  } else {
    ending := change_vovel_sound(w.kana[len(stem):], "あ")
    return stem + ending, kstem + ending 
  }
}

func (w *Word) ToRenyoukei() (string, string) {
  stem, kstem := w.ToStem()
  if w.typ[0] == "v1" {
    return stem, kstem
  } else if w.typ[0] == "v5aru" || w.typ[0] == "v5k-s" || w.typ[0] == "v5r-i" ||
       w.typ[0] == "v5u-s" || w.typ[0] == "v5uru" {
    return "?", "?"
  } else {
    ending := change_vovel_sound(w.kana[len(stem):], "い")
    return stem + ending, kstem + ending
  }
}

func (w *Word) ToRentaikei() (string, string) {
  return w.kana, w.kanji
}

func (w *Word) ToIzenkei() (string, string) {
  stem, kstem := w.ToStem()
  ending := change_vovel_sound(w.kana[len(stem):], "え")
  return stem + ending, kstem + ending
}

func (w *Word) ToMeireikei() (string, string) {
  if w.typ[0] == "v1" {
    return w.ToStem()
  } else {
    return w.ToIzenkei()
  }
}

// helper functions
func change_vovel_sound(vovel string, sound string) string {
  //lastVovel, _ := utf8.DecodeLastRuneInString(vovel)
  lastVovel := vovel
  if sound == "あ" {
    switch lastVovel {
      case "う":
        // WARNING: could be very wrong, not quiet sure. Maybe move it to the construction of the past tense
        return "わ"
      case "る":
        return "ら"
      case "す":
        return "さ"
      case "く":
        return "か"
      case "ぐ":
        return "が"
      case "む":
        return "ま"
      case "ぶ":
        return "ば"
      case "ぬ":
        return "な"
      case "つ":
        return "た"
    }
  } else if sound == "い" {
    switch lastVovel {
      case "う":
        return "い"
      case "る":
        return "り"
      case "す":
        return "し"
      case "く":
        return "き"
      case "ぐ":
        return "ぎ"
      case "む":
        return "み"
      case "ぶ":
        return "び"
      case "ぬ":
        return "に"
      case "つ":
        return "ち"
    }
  } else if sound == "え" {
    switch lastVovel {
      case "う":
        return "え"
      case "る":
        return "れ"
      case "す":
        return "せ"
      case "く":
        return "け"
      case "ぐ":
        return "げ"
      case "む":
        return "め"
      case "ぶ":
        return "べ"
      case "ぬ":
        return "ね"
      case "つ":
        return "て"
    }
  }
  return ""
}

func (w *Word) PrintConjTable() {
  var kana, kanji string
	
  w.Print()
	fmt.Printf("Present (pos)\n")
	kana, kanji = w.ToPresent(true, false)
  fmt.Printf("\tinformal: \t%s  (%s)\n", kanji, kana)
  kana, kanji = w.ToPresent(true, true)
	fmt.Printf("\tformal: \t%s  (%s)\n", kanji, kana)
	fmt.Printf("Present (neg)\n")
  kana, kanji = w.ToPresent(false, false)
	fmt.Printf("\tinformal: \t%s  (%s)\n", kanji, kana)
  kana, kanji = w.ToPresent(false, true)
	fmt.Printf("\tformal: \t%s  (%s)\n", kanji, kana)
	fmt.Printf("Past (pos)\n")
	kana, kanji = w.ToPast(true, false)
	fmt.Printf("\tinformal: \t%s  (%s)\n", kanji, kana)
	kana, kanji = w.ToPast(true, true)
	fmt.Printf("\tformal: \t%s  (%s)\n", kanji, kana)
	fmt.Printf("Past (neg)\n")
	kana, kanji = w.ToPast(false, false)
	fmt.Printf("\tinformal: \t%s  (%s)\n", kanji, kana)
	kana, kanji = w.ToPast(false, true)
	fmt.Printf("\tformal: \t%s  (%s)\n", kanji, kana)
	fmt.Printf("-te form (pos)\n")
  kana, kanji = w.ToTeForm(true, false)
	fmt.Printf("\tinformal: \t%s  (%s)\n", kanji, kana)
  kana, kanji = w.ToTeForm(true, true)
	fmt.Printf("\tformal: \t%s  (%s)\n", kanji, kana)
	fmt.Printf("-te form (neg)\n")
  kana, kanji = w.ToTeForm(false, false)
	fmt.Printf("\tinformal: \t%s  (%s)\n", kanji, kana)
  kana, kanji = w.ToTeForm(false, true)
	fmt.Printf("\tformal: \t%s  (%s)\n", kanji, kana)
	fmt.Printf("Conditional (pos)\n")
  kana, kanji = w.ToConditional(true, false)
	fmt.Printf("\tinformal: \t%s  (%s)\n", kanji, kana)
  kana, kanji = w.ToConditional(true, true)
	fmt.Printf("\tformal: \t%s  (%s)\n", kanji, kana)
	fmt.Printf("Conditional (neg)\n")
  kana, kanji = w.ToConditional(false, false)
	fmt.Printf("\tinformal: \t%s  (%s)\n", kanji, kana)
  kana, kanji = w.ToConditional(false, true)
	fmt.Printf("\tformal: \t%s  (%s)\n", kanji, kana)
}
