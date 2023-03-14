package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// Parcel represents a row in the parcels table
type Parcel struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}

// GetParcelCoords retrieves the lat and long for a given address from the database
func GetParcelCoords(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	conn, err := pgx.Connect(context.Background(), os.Getenv("ASSESSOR_DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT saddr as address, center_lat as lat, center_lon as long FROM parcels WHERE saddr = $1", address)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var parcel Parcel
		if err := rows.Scan(&parcel.Address, &parcel.Lat, &parcel.Long); err != nil {
			log.Fatal(err)
		}
		parcels = append(parcels, parcel)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(parcels)
}


// Healthz returns a 200 OK status code to indicate that the service is healthy and ready to serve requests
func Healthz(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.Background(),  os.Getenv("ASSESSOR_DB_URL"))
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer conn.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}


func main() {
	http.HandleFunc("/getcoords", GetParcelCoords)
	http.HandleFunc("/healthz", Healthz)
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}