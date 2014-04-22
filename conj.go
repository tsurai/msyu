package main

import (
	"fmt"
	"unicode/utf8"
)

type conjugation struct {
	Exec func(*word, bool, bool) (string, string)
	Name string
	Rule map[string]string
}

var (
	baseRules = `  　　         一段　　　        五段
  * 語幹:　  remove last る    remove last syllable
  * 未然形:  語幹              replace last vowel with あ-vowel
  * 連用形:  語幹              replace last vowel with い-vowel
  * 連体形:  - 　              -
  * 已然形:  語幹+る 　        replace last vowel with え-vowel
  * 命令形:  語幹   　　       replace last vowel with え-vowel`
)

var conjugations = []conjugation{
	{
		Exec: (*word).ToPresent,
		Name: "Present Tense",
		Rule: map[string]string{
			"v1": `* Positive Plain:  連体形 
* Positive Polite: 連用形 + ます
* Negative Plain:  未然形
* Negative Polite: 連用形 + ません`,
			"v5": `* Positive Plain:  連体形 
* Positive Polite: 連用形 + ます
* Negative Plain:  未然形
* Negative Polite: 連用形 + ません`,
		},
	},
	{
		Exec: (*word).ToPast,
		Name: "Past Tense",
		Rule: map[string]string{
			"v1": `* Positive Plain:  連用形 + た
* Positive Polite: 連用形 + ました
* Negative Plain:  未然形 + なかった
* Negative Polite: 連用形 + ませんでした`,
			"v5": `* Positive Plain:  連用形 + た *
* Positive Polite: 連用形 + ました
* Negative Plain:  未然形 + なかった
* Negative Polite: 連用形 + ませんでした`,
		},
	},
	{
		Exec: (*word).ToTeForm,
		Name: "Te Form",
		Rule: map[string]string{
			"v1": `* Positive Plain:  連用形 + て
* Positive Polite: 連用形 + まして
* Negative Plain:  未然形 + ないで
* Negative Polite: 連用形 + ませんで`,
			"v5": `* Positive Plain:  連用形 + て *
* Positive Polite: 連用形 + まして
* Negative Plain:  未然形 + ないで
* Negative Polite: 連用形 + ませんで`,
		},
	},
	{
		Exec: (*word).ToConditional,
		Name: "Conditional",
		Rule: map[string]string{
			"v1": `* Positive Plain:  連用形 + たら
* Positive Polite: 連用形 + ましたら
* Negative Plain:  未然形 + なかったら
* Negative Polite: 連用形 + ませんでしたら`,
			"v5": `* Positive Plain:  連用形 + たら *
* Positive Polite: 連用形 + ましたら
* Negative Plain:  未然形 + なかったら
* Negative Polite: 連用形 + ませんでしたら`,
		},
	},
	{
		Exec: (*word).ToProvisional,
		Name: "Provisional",
		Rule: map[string]string{
			"v1": `* Positive Plain:  已然形 + ば
* Positive Polite: 連体形 + なら
* Negative Plain:  未然形 + なければ 
* Negative Polite: 連用形 + ませんなら`,
			"v5": `* Positive Plain:  已然形 + ば
* Positive Polite: 連体形 + なら
* Negative Plain:  未然形 + なければ 
* Negative Polite: 連用形 + ませんなら`,
		},
	},
	{
		Exec: (*word).ToPassiveAndPotentional,
		Name: "Passive & Potentional",
		Rule: map[string]string{
			"v1": `* Positive Plain:  未然形 + れる
* Positive Polite: 未然形 + れます
* Negative Plain:  未然形 + れない
* Negative Polite: 未然形 + れません`,
			"v5": `* Positive Plain:  未然形 + れる
* Positive Polite: 未然形 + れます
* Negative Plain:  未然形 + れない
* Negative Polite: 未然形 + れません`,
		},
	},
	{
		Exec: (*word).ToCausative,
		Name: "Causative",
		Rule: map[string]string{
			"v1": `* Positive Plain:  未然形 + せる
* Positive Polite: 未然形 + させます
* Negative Plain:  未然形 + させない
* Negative Polite: 未然形 + させません`,
			"v5": `* Positive Plain:  未然形 + させる
* Positive Polite: 未然形 + せます
* Negative Plain:  未然形 + せない
* Negative Polite: 未然形 + せません`,
		},
	},
	{
		Exec: (*word).ToCausativePassive,
		Name: "Causative Passive",
		Rule: map[string]string{
			"v1": `* Positive Plain:  未然形 + させられる
* Positive Polite: 未然形 + させられます
* Negative Plain:  未然形 + させられない 
* Negative Polite: 未然形 + させられません`,
			"v5": `* Positive Plain:  未然形 + せられる
* Positive Polite: 未然形 + せられます
* Negative Plain:  未然形 + せられない
* Negative Polite: 未然形 + せられません`,
		},
	},
	{
		Exec: (*word).ToConjectural,
		Name: "Conjectural",
		Rule: map[string]string{
			"v1": `* Positive Plain:  連体形 + だろう
* Positive Polite: 連体形 + でしょう
* Negative Plain:  未然形 + ないだろう
* Negative Polite: 未然形 + ないでしょう`,
			"v5": `* Positive Plain:  連体形 + だろう
* Positive Polite: 連体形 + でしょう
* Negative Plain:  未然形 + ないだろう
* Negative Polite: 未然形 + ないでしょう`,
		},
	},
	{
		Exec: (*word).ToAlternative,
		Name: "Alternative",
		Rule: map[string]string{
			"v1": `* Positive Plain:  連用形 + たり
* Positive Polite: 連用形 + ましたり
* Negative Plain:  未然形 + なかったり
* Negative Polite: 連用形 + ませんでしたり`,
			"v5": `* Positive Plain:  連用形 + たり *
* Positive Polite: 連用形 + ましたり
* Negative Plain:  未然形 + なかったり
* Negative Polite: 連用形 + ませんでしたり`,
		},
	},
	{
		Exec: (*word).ToImperative,
		Name: "Imperative",
		Rule: map[string]string{
			"v1": `* Positive Plain:  命令形 + ろ
* Positive Polite: 連用形 + なさい
* Negative Plain:  連体形 + な
* Negative Polite: 連用形 + なさるな`,
			"v5": `* Positive Plain:  命令形
* Positive Polite: 連用形 + なさい
* Negative Plain:  連体形 + な
* Negative Polite: 連用形 + なさるな`,
		},
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

	if kanji != "" {
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
				switch w.kana[len(kana):] {
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

	if kanji != "" {
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
				switch w.kana[len(kana):] {
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

	if kanji != "" {
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
				switch w.kana[len(kana):] {
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

	if kanji != "" {
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
			ending = "ますなら"
		} else {
			kana, kanji = w.ToIzenkei()
			ending = "ば"
		}
	} else {
		if formal {
			kana, kanji = w.ToMizenkei()
			ending = "ませんなら"
		} else {
			kana, kanji = w.ToRenyoukei()
			ending = "なければ"
		}
	}

	if kanji != "" {
		return kana + ending, kanji + ending
	} else {
		return kana + ending, kanji
	}
}

func (w *word) ToPassiveAndPotentional(positive bool, formal bool) (string, string) {
	var ending string
	kana, kanji := w.ToMizenkei()

	if positive {
		if formal {
			ending = "れます"
		} else {
			ending = "れる"
		}
	} else {
		if formal {
			ending = "れません"
		} else {
			ending = "れない"
		}
	}

	if kanji != "" {
		return kana + ending, kanji + ending
	} else {
		return kana + ending, kanji
	}
}

func (w *word) ToCausative(positive bool, formal bool) (string, string) {
	var ending string
	kana, kanji := w.ToMizenkei()

	if positive {
		if formal {
			ending = "せます"
		} else {
			ending = "せる"
		}
	} else {
		if formal {
			ending = "せません"
		} else {
			ending = "せない"
		}
	}

	if kanji != "" {
		return kana + ending, kanji + ending
	} else {
		return kana + ending, kanji
	}
}

func (w *word) ToCausativePassive(positive bool, formal bool) (string, string) {
	var ending string
	kana, kanji := w.ToMizenkei()

	if positive {
		if formal {
			ending = "せられます"
		} else {
			ending = "せられる"
		}
	} else {
		if formal {
			ending = "せられません"
		} else {
			ending = "せられない"
		}
	}

	if kanji != "" {
		return kana + ending, kanji + ending
	} else {
		return kana + ending, kanji
	}
}

func (w *word) ToConjectural(positive bool, formal bool) (string, string) {
	var kana, kanji, ending string

	if positive {
		kana, kanji = w.ToRentaikei()
		if formal {
			ending = "でしょう"
		} else {
			ending = "だろう"
		}
	} else {
		kana, kanji = w.ToMizenkei()
		if formal {
			ending = "ないでしょう"
		} else {
			ending = "なかっただろう"
		}
	}

	if kanji != "" {
		return kana + ending, kanji + ending
	} else {
		return kana + ending, kanji
	}
}

func (w *word) ToAlternative(positive bool, formal bool) (string, string) {
	var ending string
	kana, kanji := w.ToRenyoukei()

	if positive {
		if formal {
			ending = "ましたり"
		} else {
			if w.gloss[0].pos[0] == "v1" {
				ending = "たり"
			} else {
				kana, kanji = w.ToStem()
				switch w.kana[len(kana):] {
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
					ending = "たり"
				}
			}
		}
	} else {
		if formal {
			ending = "ませんでしたり"
		} else {
			kana, kanji = w.ToMizenkei()
			ending = "なかったり"
		}
	}

	if kanji != "" {
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
			if w.gloss[0].pos[0] == "v1" {

			} else {
				kana, kanji = w.ToIzenkei()
			}
		}
	} else {
		if formal {
			kana, kanji = w.ToRentaikei()
			ending = "なさるな"
		} else {
			kana, kanji = w.ToRenyoukei()
			ending = "な"
		}
	}

	if kanji != "" {
		return kana + ending, kanji + ending
	} else {
		return kana + ending, kanji
	}
}

// Inflection bases -------------
func (w *word) ToStem() (string, string) {
	kanji := ""
	_, size := utf8.DecodeLastRuneInString(w.kana)

	if w.kanji[0] != "" {
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
		ending := changeVovelSound(w.kana[len(stem):], "あ")

		if kstem != "" {
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
		ending := changeVovelSound(w.kana[len(stem):], "い")

		if kstem != "" {
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
	ending := changeVovelSound(w.kana[len(stem):], "え")

	if kstem != "" {
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
func changeVovelSound(vovel string, sound string) string {
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

	// make a proper class for the conjugations with proper building rules
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
