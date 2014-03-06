package main

import (
  "fmt"
  "log"
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
  VERB = 10
  ADJ = 11
  NOUN = 12
)

var database *sql.DB = nil

func DB_init() {
  var err error
  database, err = sql.Open("sqlite3", "JMdict.db")
  if err != nil {
    log.Fatal("A database error has occured:", err)
  }
}

func DB_close() {
  database.Close()
}

func db_parse_results(rows *sql.Rows) ([]*word, int) {
  var rvalue sql.NullString
  var kvalue sql.NullString
  var pos sql.NullString
  var meaning sql.NullString
  var words []*word

  id := 0
  lastId := 0
  for rows.Next() {
    var g []*gloss
    rows.Scan(&id, &rvalue, &pos, &meaning, &kvalue)

    if(lastId != id) {
      lastId = id
      g = append(g, &gloss{strings.Split(pos.String, ","), strings.Split(meaning.String, ",")})
      words = append(words, &word{rvalue.String, strings.Split(kvalue.String, ","), g})
    } else {
      words[len(words)-1].gloss = append(words[len(words)-1].gloss, &gloss{strings.Split(pos.String, ","), strings.Split(meaning.String, ",")})
    }
  }

  return words, len(words)
}

func DB_search_word(w string, mode int, filter int) *word {
  if w == "" {
    fmt.Println("Missing parameter")
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
    query = fmt.Sprintf("SELECT r_ele.fk, r_ele.value, " +
            "GROUP_CONCAT(DISTINCT entity.entity), " +
            "GROUP_CONCAT(DISTINCT gloss.value), " +
            "GROUP_CONCAT(DISTINCT k_ele.value) FROM r_ele, k_ele, gloss, sense " +
            "LEFT OUTER JOIN pos ON sense.id = pos.fk " +
            "LEFT OUTER JOIN entity ON pos.entity = entity.id " +
            "WHERE r_ele.id IN (SELECT r_ele.id FROM r_ele, sense, pos, entity WHERE " + sqlfilter +
            "AND r_ele.fk = sense.fk AND sense.id = pos.fk AND pos.entity = entity.id) " +
            "AND (r_ele.value LIKE '%%%s%%' OR k_ele.value LIKE '%%%s%%') " +
            "AND sense.fk = k_ele.fk AND r_ele.fk = sense.fk AND gloss.fk = sense.id " +
            "GROUP BY sense.id ORDER BY length(r_ele.value), r_ele.fk", w, w)

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
            "ORDER BY length(r_ele.value)", w)
  default:
    panic("Unknown search mode")
  }

  rows, err := database.Query(query)
  if err != nil {
    log.Fatal("A database error has occured:", err)
  }

  words, num := db_parse_results(rows)

  if num > 1 {
    offset := 0

    for {
      fmt.Println("--------------------")
      for i := offset; i < offset+step_size && i < num; i++ {
        w := words[i]

        fmt.Printf("%d: ", i+1)
        w.Print()
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

func DB_get_random_verbs(n int) []*word {
  if n <= 0 {
    fmt.Println("Invalid parameter")
    return nil
  }

  query := fmt.Sprintf("SELECT r_ele.fk, r_ele.value, " +
            "GROUP_CONCAT(DISTINCT entity.entity), " +
            "GROUP_CONCAT(DISTINCT gloss.value), " +
            "GROUP_CONCAT(DISTINCT k_ele.value) FROM r_ele, k_ele, gloss, sense " +
            "LEFT OUTER JOIN pos ON sense.id = pos.fk " +
            "LEFT OUTER JOIN entity ON pos.entity = entity.id " +
            "WHERE r_ele.id IN (SELECT r_ele.id FROM r_ele, sense, pos, entity WHERE (entity.entity = 'v1' OR entity.entity LIKE 'v5%%')  " +
            "AND r_ele.fk = sense.fk AND sense.id = pos.fk AND pos.entity = entity.id ORDER BY RANDOM() LIMIT %d) " +
            "AND sense.fk = k_ele.fk AND r_ele.fk = sense.fk AND gloss.fk = sense.id " +
            "GROUP BY sense.id, pos.fk ORDER BY r_ele.fk", n)

  rows, err := database.Query(query)
  if err != nil {
    fmt.Println("A database error has occured:", err)
    return nil
  }
  defer rows.Close()

  w, _ := db_parse_results(rows)

  return w
}
