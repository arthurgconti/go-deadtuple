package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"

	_ "github.com/lib/pq"
)

type tuple_query struct {
	relname    string
	dead_tuple string
	live_tuple string
}

func main() {
	fmt.Println("hello")
	//	deadTupleSql := "select relname,n_dead_tup,n_live_tup,last_vacuum,last_analyze, last_autovacuum from pg_stat_all_tables order by 2 desc;"
	tuples, err := getTuples()
	if err != nil {
		panic(err)
	}
	segregateMetrics(tuples)
}

func OpenConn() (*sql.DB, error) {
	dbConn := "host=" + os.Getenv("DB_HOST") + " " + "port=" + os.Getenv("DB_PORT") + " " + "user=" + os.Getenv("DB_USER") + " " + "password=" + os.Getenv("DB_PASS") + " " + "dbname=" + os.Getenv("DB_NAME") + " sslmode=disable"
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	return db, err
}

func getTuples() (tp []tuple_query, err error) {
	//deadTupleSql := "select relname,n_dead_tup,n_live_tup from pg_stat_all_tables order by 2 desc;"
	deadTupleSql := "select relname,n_dead_tup,n_live_tup from pg_stat_all_tables where n_dead_tup <> 0 and n_live_tup <>0 order by 2 desc;"
	conn, err := OpenConn()
	if err != nil {
		return
	}
	defer conn.Close()
	rows, err := conn.Query(deadTupleSql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t tuple_query
		err = rows.Scan(&t.relname, &t.dead_tuple, &t.live_tuple)
		if err != nil {
			continue
		}
		tp = append(tp, t)
	}
	return
}

// PUT metrics to prometheus/CloudwatchLogs
func putMetrics() {}

// Kind sort algorithm to separate the major pg offensor dead_tuples and major user table offensor
func segregateMetrics(tp []tuple_query) {
	regex := regexp.MustCompile(`^pg_.*`)
	var pg_counter, counter int = 0, 0
	var highTuples []tuple_query

	for i, s := range tp {
		fmt.Println(i+1, s.relname)
		if regex.MatchString(s.relname) && pg_counter == 0 {
			pg_counter++
			highTuples = append(highTuples, s)
		}
		if !regex.MatchString(s.relname) && counter == 0 {
			counter++
			highTuples = append(highTuples, s)
		}
	}

	fmt.Println("high:", highTuples)
}
