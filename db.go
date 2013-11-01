package main

import (
  "fmt"
  "os"
  "strconv"
  "strings"
  "database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
  EN = 0
  JAP = 1
)
var database *sql.DB = nil

func DB_init() {
  var err error
  database, err = sql.Open("sqlite3", "JMdict.db")
  if err != nil {
    fmt.Fprintln(os.Stderr, "A database error has occured:", err)
    os.Exit(1)
  }
}

func DB_close() {
  database.Close()
}

type Gloss struct {
  pos []string
  meaning []string
}

type Wordd struct {
  kana string
  kanji string
  gloss map[int]*Gloss
}

func db_parse_results(rows *sql.Rows) (map[int]*Wordd, int) {
  var id int
  var rvalue sql.NullString
  var kvalue sql.NullString
  var pos sql.NullString
  var meaning sql.NullString
  words := make(map[int]*Wordd)

  lastId := 0
  nGloss := 0
  nRows := 0
  for rows.Next() {
    rows.Scan(&id, &rvalue, &pos, &meaning, &kvalue)
    if(lastId != id) {
      nGloss = 0
      lastId = id
      gloss := make(map[int]*Gloss)
      gloss[0] = &Gloss{strings.Split(pos.String, ","), strings.Split(meaning.String, ",")}
      words[nRows] = &Wordd{rvalue.String, kvalue.String, gloss}
      nRows = nRows + 1
    } else {
      nGloss = nGloss + 1
      words[nRows-1].gloss[nGloss] = &Gloss{strings.Split(pos.String, ","), strings.Split(meaning.String, ",")}
    }
  }

  return words, nRows
}

func DB_search_word(word string, mode int) *Wordd {
  if word == "" {
    fmt.Fprintln(os.Stderr, "Missing parameter")
    return nil
  }

  var query string

  if mode == JAP {
     query = fmt.Sprintf("SELECT sense.fk, r_ele.value, " +
              "group_concat(DISTINCT entity.entity), " +
              "group_concat(DISTINCT gloss.value), "  +
              "group_concat(DISTINCT k_ele.value) FROM r_ele, gloss, sense " +
              "LEFT JOIN k_ele ON sense.fk = k_ele.fk " +
              "LEFT JOIN pos ON sense.id = pos.fk " +
              "LEFT JOIN entity ON pos.entity = entity.id " +
              "WHERE gloss.fk = sense.id AND sense.fk = r_ele.fk AND sense.fk IN " +
              "(SELECT r_ele.fk FROM r_ele, k_ele WHERE r_ele.fk = k_ele.fk AND " +
              "(r_ele.value LIKE '%%%s%%' OR k_ele.value LIKE '%%%s%%')) " +
              "GROUP BY sense.id ORDER BY length(r_ele.value)", word, word)
  } else if mode == EN {
     query = fmt.Sprintf("SELECT sense.fk, r_ele.value, " +
              "group_concat(DISTINCT entity.entity), " +
              "group_concat(DISTINCT gloss.value), " +
              "group_concat(DISTINCT k_ele.value) FROM r_ele, gloss, sense " +
              "LEFT JOIN k_ele ON sense.fk = k_ele.fk " +
              "LEFT JOIN pos ON sense.id = pos.fk " +
              "LEFT JOIN entity ON pos.entity = entity.id " +
              "WHERE gloss.fk = sense.id AND sense.fk = r_ele.fk AND sense.fk IN " +
              "(SELECT sense.fk FROM sense, gloss WHERE gloss.value LIKE '%%%s%%' " +
              "AND gloss.fk = sense.id) GROUP BY sense.id, pos.fk " +
              "ORDER BY length(r_ele.value)", word)
  } else {
    panic("Unknown search mode")
  }

  rows, err := database.Query(query)
  if err != nil {
    fmt.Fprintln(os.Stderr, "A database error has occured:", err)
    os.Exit(1)
  }

  words, num := db_parse_results(rows)

  if num > 1 {
    offset := 0

    for {
      for i := offset; i < offset+5; i++ {
        word := words[i]

        fmt.Printf("%d: ", i+1)
        word.Print()
      }

      for {
        entry := ""
        fmt.Printf("%d-%d of %d | <n> for next | <p> for previous\n", offset+1, offset+5, num)
        fmt.Print("Select an Entry: ");
        fmt.Scanf("%s", &entry)

        if entry == "n" {
          if((offset + 5) <= num) {
            offset = offset + 5
            fmt.Println("")
          }
        } else if entry == "p" {
          if((offset - 5) >= 0) {
            offset = offset - 5
            fmt.Println("")
          }
        } else {
          if i, err := strconv.Atoi(entry); err == nil {
            return words[i]
          } else {
            fmt.Println("Invalid input. Try again\n")
          }
        }
      }
    }
  }
  return words[0] 
}

func DB_get_random_verbs(n int) map[int]*Word {
  if n <= 0 {
    fmt.Fprintln(os.Stderr, "Invalid parameter")
    return nil
  }
  
  query := fmt.Sprintf("SELECT R.value, group_concat(DISTINCT E.entity), group_concat(G.value, ','), K.value " +
           "FROM r_ele as R, entity as E, gloss as G, pos as P, sense as S, k_ele as K " +
           "WHERE R.fk = S.fk AND S.id = P.fk AND E.id = P.entity AND S.id = G.fk AND R.fk = K.fk AND " +
           "(E.entity = 'v1' OR E.entity = 'v5%%')" +
           "GROUP BY R.id ORDER BY RANDOM() LIMIT %d", n)

  rows, err := database.Query(query)
  if err != nil {
    fmt.Fprintln(os.Stderr, "A database error has occured:", err)
    return nil
  }
  defer rows.Close()

  i := 0
  var words = make(map[int]*Word)
  var rvalue string
  var kvalue string
  var pos string
  var gloss string

  for rows.Next() {
    rows.Scan(&rvalue, &pos, &gloss, &kvalue)

    words[i] = &Word{rvalue, kvalue, strings.Split(pos, ", "), gloss}
    i = i + 1
  }

  return words
}
