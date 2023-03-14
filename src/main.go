package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the Assessor API! Made by Kyler Chin https://github.com/kylerchin")
}

func main() {
// Connect  to database
// urlExample := "postgres://username:password@localhost:5432/database_name"
conn, err := pgx.Connect(context.Background(), os.Getenv("ASSESSOR_DB_URL"))
if err != nil {
	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	os.Exit(1)
}
defer conn.Close(context.Background())


    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homeLink)

    log.Fatal(http.ListenAndServe(":8080", router))
}