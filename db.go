package main

import (
  "fmt"
  "os"
  "strings"
  "database/sql"
	_ "github.com/mattn/go-sqlite3"
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

/* this does almost the same as DB_search_in_jap. make a function
   that takes a query as argument and returns the word to keep it simple */
func DB_search_in_eng(word string) *Word {
  if word == "" {
    fmt.Fprintln(os.Stderr, "Missing parameter")
    return nil
  }

  /* FIXME: fetches wrong pos infos */
  query := fmt.Sprintf("SELECT R.value, group_concat(E.entity), group_concat(DISTINCT G.value), K.value " +
           "FROM r_ele as R, entity as E, gloss as G, pos as P, sense as S, k_ele as K " +
           "WHERE R.fk = S.fk AND S.id = P.fk AND E.id = P.entity AND S.id = G.fk AND R.fk = K.fk AND " +
           "G.value = '%s' " +
           "GROUP BY G.id ORDER BY R.fk LIMIT 10", word)

  var rvalue string
  var kvalue string
  var pos string
  var gloss string
  words := make(map[int]*Word)

  rows, err := database.Query(query)
  if err != nil {
    fmt.Fprintln(os.Stderr, "A database error has occured:", err)
    os.Exit(1)
  }

  nRows := 0
  for rows.Next() {
    rows.Scan(&rvalue, &pos, &gloss, &kvalue)
    words[nRows] = &Word{rvalue, kvalue, strings.Split(pos, ", "), gloss}
    nRows = nRows + 1
  }

  if nRows > 1 {
    fmt.Println("Select an entry:")
    for n := 0; n < nRows; n++ {
      fmt.Printf("%d: ", n+1)
      words[n].Print()
    }

    entry := -1
    fmt.Print("Entry number: ");
    fmt.Scanf("%d", &entry)

    return words[entry-1]
  }

  return words[0] 
}

func DB_search_in_jap(word string) *Word {
  if word == "" {
    fmt.Fprintln(os.Stderr, "Missing parameter")
    return nil
  }

  query := fmt.Sprintf("SELECT R.value, group_concat(E.entity), group_concat(DISTINCT G.value), K.value " +
           "FROM r_ele as R, entity as E, gloss as G, pos as P, sense as S, k_ele as K " +
           "WHERE R.fk = S.fk AND S.id = P.fk AND E.id = P.entity AND S.id = G.fk AND R.fk = K.fk AND " +
           "(R.value = '%s' OR K.value = '%s') " +
           "ORDER BY R.id LIMIT 1", word, word)

  var rvalue string
  var kvalue string
  var pos string
  var gloss string

  err := database.QueryRow(query).Scan(&rvalue, &pos, &gloss, &kvalue)
  if err != nil { 
    return nil
  }

  return &Word{rvalue, kvalue, strings.Split(pos, ", "), gloss}
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
