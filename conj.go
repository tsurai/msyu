package main

import (
	"unicode/utf8"
	"fmt"
)

type conjunction struct {
  Exec func(*word, bool,bool) (string, string)
  Name string
  Rule string
}

var conjunctions = []conjunction {
  {
    Exec:       (*word).ToPresent,
    Name:       "Present Tense",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToPast,
    Name:       "Past Tense",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToTeForm,
    Name:       "Te Form",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToConditional,
    Name:       "Conditional",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToProvisional,
    Name:       "Provisional",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToPassiveAndPotentional,
    Name:       "Passive & Potentional",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToCausative,
    Name:       "Causative",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToCausativePassive,
    Name:       "Causative Passive",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToConjectural,
    Name:       "Conjectural",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToAlternative,
    Name:       "Alternative",
    Rule:       ``,
  },
  {
    Exec:       (*word).ToImperative,
    Name:       "Imperative",
    Rule:       ``,
  },
}

// conjunction ------------
func (w *word) ToPresent(positive bool, formal bool) (string, string) {
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
  
  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToPast(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ました"
    } else {
      if w.gloss[0].pos[0] == "v1" {
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

  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToTeForm(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string  
  
  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "まして"
    } else {
      if w.gloss[0].pos[0] == "v1" {
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
  
  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToConditional(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string
  
  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ましたら"
    } else {
      if w.gloss[0].pos[0] == "v1" {
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

  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToProvisional(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ますなら(ば)"
    } else {
      kana, kanji = w.ToIzenkei()
      ending = "ば"
    }
  } else {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "ませんなら(ば)"
    } else {
      kana, kanji = w.ToRenyoukei()
      ending = "なければ"
    }
  }

  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToPassiveAndPotentional(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "れます"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "れる"
    }
  } else {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "れません"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "れない"
    }
  }

  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToCausative(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "せます"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "せる"
    }
  } else {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "せません"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "せない"
    }
  }

  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToCausativePassive(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "せられます"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "せられる"
    }
  } else {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "せられません"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "せられない"
    }
  }

  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToConjectural(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToRentaikei()
      ending = "でしょう"
    } else {
      kana, kanji = w.ToRentaikei()
      ending = "だろう"
    }
  } else {
    if formal {
      kana, kanji = w.ToMizenkei()
      ending = "ないでしょう"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "なかっただろう"
    }
  }

  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToAlternative(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ましたり"
    } else {
      if w.gloss[0].pos[0] == "v1" {
        kana, kanji = w.ToRenyoukei()
        ending = "たり"
      } else {
        kana, kanji = w.ToStem()
        switch(w.kana[len(kana):]) {
          case "く":
            ending = "いたり"
          case "ぐ":
            ending = "いだり"
          case "ぬ":
            fallthrough
          case "ぶ":
            fallthrough
          case "む":
            ending = "んだり"
          case "う":
            fallthrough
          case "つ":
            fallthrough
          case "る":
            ending = "ったり"
          default:
            kana, kanji = w.ToRenyoukei()
            ending = "たり"
        }
      }
    }
  } else {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "ませんでしたり"
    } else {
      kana, kanji = w.ToMizenkei()
      ending = "なかったり"
    }
  }
  
  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

func (w *word) ToImperative(positive bool, formal bool) (string, string) {
  var kana, kanji, ending string

  if positive {
    if formal {
      kana, kanji = w.ToRenyoukei()
      ending = "なさい"
    } else {
      kana, kanji = w.ToMeireikei()
      ending = "でしょう"
    }
  } else {
    if formal {
      kana, kanji = w.ToRentaikei()
      ending = "な"
    } else {
      kana, kanji = w.ToRenyoukei()
      ending = "なさるな"
    }
  }
  
  if(kanji != "") {
    return kana + ending, kanji + ending
  } else {
    return kana + ending, kanji
  }
}

// Inflection bases -------------
func (w *word) ToStem() (string, string) {
  kanji := ""
  _, size := utf8.DecodeLastRuneInString(w.kana)

  if(w.kanji[0] != "") {
    kanji = w.kanji[0][:len(w.kanji[0])-size]
  }

  return w.kana[:len(w.kana)-size], kanji
}

func (w *word) ToMizenkei() (string, string) {
  pos := w.gloss[0].pos[0]
  stem, kstem := w.ToStem()
  
  if w.gloss[0].pos[0] == "v1" {
    return stem, kstem
  } else if pos == "v5aru" || pos == "v5k-s" || pos == "v5r-i" ||
       pos == "v5u-s" || pos == "v5uru" {
    return "?", "?"
  } else {
    ending := change_vovel_sound(w.kana[len(stem):], "あ")
    
    if(kstem != "") {
      return stem + ending, kstem + ending
    } else {
      return stem + ending, kstem
    }
  }
}

func (w *word) ToRenyoukei() (string, string) {
  pos := w.gloss[0].pos[0]
  stem, kstem := w.ToStem()
  
  if w.gloss[0].pos[0] == "v1" {
    return stem, kstem
  } else if pos == "v5aru" || pos == "v5k-s" || pos == "v5r-i" ||
       pos == "v5u-s" || pos == "v5uru" {
    return "?", "?"
  } else {
    ending := change_vovel_sound(w.kana[len(stem):], "い")
    
    if(kstem != "") {
      return stem + ending, kstem + ending
    } else {
      return stem + ending, kstem
    }
  }
}

func (w *word) ToRentaikei() (string, string) {
  return w.kana, w.kanji[0]
}

func (w *word) ToIzenkei() (string, string) {
  stem, kstem := w.ToStem()
  ending := change_vovel_sound(w.kana[len(stem):], "え")
  
  if(kstem != "") {
    return stem + ending, kstem + ending
  } else {
    return stem + ending, kstem
  }
}

func (w *word) ToMeireikei() (string, string) {
  if w.gloss[0].pos[0] == "v1" {
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

func (w *word) PrintConjTable() {
  var kana, kanji string
  
  /* make a proper class for the conjunctions with proper building rules */
  conj := make(map[string]func(bool, bool) (string, string))
  conj["Present"] = w.ToPresent
  conj["Past"] = w.ToPast
  conj["-te Form"] = w.ToTeForm
  conj["Conditional"] = w.ToConditional
  conj["Provisional"] = w.ToProvisional
  conj["Passive & Potentional"] = w.ToPassiveAndPotentional 
  conj["Causative"] = w.ToCausative
  conj["Causative Passive"] = w.ToCausativePassive
  conj["Conjectural"] = w.ToConjectural
  conj["Alternative"] = w.ToAlternative
  conj["Imperative"] = w.ToImperative

  fmt.Println("--------------------")
  w.Print()
  fmt.Println("")
  for n, f := range conj {
	  fmt.Printf("%s (pos)\n", n)
	  kana, kanji = f(true, false)
    fmt.Printf("\tinformal: \t%s  %s\n", kanji, kana)
    kana, kanji = f(true, true)
	  fmt.Printf("\tformal: \t%s  %s\n", kanji, kana)
	  fmt.Printf("%s (neg)\n", n)
    kana, kanji = f(false, false)
	  fmt.Printf("\tinformal: \t%s  %s\n", kanji, kana)
    kana, kanji = f(false, true)
	  fmt.Printf("\tformal: \t%s  %s\n", kanji, kana)
  }
}
