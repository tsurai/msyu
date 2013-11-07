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
  // lang
  EN = 0
  JAP = 1
  // filter
  VERB = 2
  ADJ = 3
  NOUN = 4
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

func db_parse_results(rows *sql.Rows) (map[int]*Word, int) {
  var id int
  var rvalue sql.NullString
  var kvalue sql.NullString
  var pos sql.NullString
  var meaning sql.NullString
  words := make(map[int]*Word)

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
      words[nRows] = &Word{rvalue.String, strings.Split(kvalue.String, ","), gloss}
      nRows = nRows + 1
    } else {
      nGloss = nGloss + 1
      words[nRows-1].gloss[nGloss] = &Gloss{strings.Split(pos.String, ","), strings.Split(meaning.String, ",")}
    }
  }

  return words, nRows
}

func DB_search_word(word string, mode int, filter int) *Word {
  if word == "" {
    fmt.Fprintln(os.Stderr, "Missing parameter")
    return nil
  }

  step_size := 5
  sqlfilter := ""
  var query string

  switch(filter) {
  case VERB:
    sqlfilter = " (entity.entity = 'v1' OR entity.entity LIKE 'v5%%') "
  }

  switch(mode) {
  case JAP:
    query = fmt.Sprintf("SELECT sense.fk, r_ele.value, " +
            "group_concat(DISTINCT entity.entity), " +
            "group_concat(DISTINCT gloss.value), "  +
            "group_concat(DISTINCT k_ele.value) FROM r_ele, gloss, sense " +
            "LEFT JOIN k_ele ON sense.fk = k_ele.fk " +
            "LEFT JOIN pos ON sense.id = pos.fk " +
            "LEFT JOIN entity ON pos.entity = entity.id " +
            "WHERE gloss.fk = sense.id AND sense.fk = r_ele.fk AND " + sqlfilter +
            "AND sense.fk IN (SELECT r_ele.fk FROM r_ele, k_ele WHERE r_ele.fk = k_ele.fk AND " +
            "(r_ele.value LIKE '%%%s%%' OR k_ele.value LIKE '%%%s%%')) " +
            "GROUP BY sense.id ORDER BY length(r_ele.value)", word, word)
  case EN:
    query = fmt.Sprintf("SELECT sense.fk, r_ele.value, " +
            "group_concat(DISTINCT entity.entity), " +
            "group_concat(DISTINCT gloss.value), " +
            "group_concat(DISTINCT k_ele.value) FROM r_ele, gloss, sense " +
            "LEFT JOIN k_ele ON sense.fk = k_ele.fk " +
            "LEFT JOIN pos ON sense.id = pos.fk " +
            "LEFT JOIN entity ON pos.entity = entity.id " +
            "WHERE gloss.fk = sense.id AND sense.fk = r_ele.fk AND " + sqlfilter +
            "AND sense.fk IN (SELECT sense.fk FROM sense, gloss WHERE gloss.value LIKE '%%%s%%' " +
            "AND gloss.fk = sense.id) GROUP BY sense.id, pos.fk " +
            "ORDER BY length(r_ele.value)", word)
  default:
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
      fmt.Println("--------------------")
      for i := offset; i < offset+step_size && i < num; i++ {
        word := words[i]

        fmt.Printf("%d: ", i+1)
        word.Print()
      }

      valid := false
      for !valid {
        entry := ""
        info := "%d-%d of %d"
        
        if(offset + 5 < num) {
          info += " | <n> for next"
        }
        if(offset >= 5) {
          info += " | <p> for previous"
        }
        fmt.Printf(info+"\nSelect an Entry : ", offset + 1, offset + step_size, num)
        fmt.Scanf("%s", &entry)

        if entry == "n" {
          if(offset + step_size < num) {
            offset = offset + step_size
            valid = true
          }
        } else if entry == "p" {
          if((offset - step_size) >= 0) {
            offset = offset - step_size
            valid = true
          }
        } else if i, err := strconv.Atoi(entry); err == nil {
            return words[i-1]
        }

        if(!valid) {
          fmt.Println("Invalid input. Try again\n--------------------")
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
  
  query := fmt.Sprintf("SELECT sense.fk, r_ele.value, " +
              "group_concat(DISTINCT entity.entity), " +
              "group_concat(DISTINCT gloss.value), " +
              "group_concat(DISTINCT k_ele.value) FROM r_ele, gloss, sense " +
              "LEFT JOIN k_ele ON sense.fk = k_ele.fk " +
              "LEFT JOIN pos ON sense.id = pos.fk " +
              "LEFT JOIN entity ON pos.entity = entity.id " +
              "WHERE gloss.fk = sense.id AND sense.fk = r_ele.fk AND sense.fk IN " +
              "(SELECT sense.fk FROM sense, pos, entity WHERE (entity.entity = 'v1' OR entity.entity LIKE 'v5%%') " +
              "AND entity.id = pos.entity AND pos.fk = sense.id) GROUP BY sense.id, pos.fk " +
              "ORDER BY RANDOM() LIMIT %d", n)
  
  rows, err := database.Query(query)
  if err != nil {
    fmt.Fprintln(os.Stderr, "A database error has occured:", err)
    return nil
  }
  defer rows.Close()

  word, _ := db_parse_results(rows)

  return word
}
