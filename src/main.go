package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/jackc/pgx/v5"
)

// Config represents the JSON configuration file structure
type Config struct {
	PGHost     string `json:"pg_host"`
	PGUser     string `json:"pg_user"`
	PGPassword string `json:"pg_password"`
	PGDatabase string `json:"pg_database"`
	PGPort     string `json:"pg_port"`
}

// Parcel represents a row in the parcels table
type Parcel struct {
	Address string  `json:"address"`
	City    string  `json:"city"`
	Zip     string  `json:"zip"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}

// GetParcelCoords retrieves the lat and long for a given address from the database
func GetParcelCoords(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	conn, err := pgx.Connect(context.Background(), GetPGConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT saddr as address, taxratecity as city, zip, center_lat as lat, center_lon as long FROM parcels WHERE saddr = $1", address)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var parcel Parcel
		if err := rows.Scan(&parcel.Address, &parcel.City, &parcel.Zip, &parcel.Lat, &parcel.Long); err != nil {
			log.Fatal(err)
		}
		parcels = append(parcels, parcel)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(parcels)
}

// GetPGConfig returns a pgx.ConnConfig object with PostgreSQL connection information
func GetPGConfig() string {
	config := ReadConfig()

	fmt.Println("PGHost", config.PGHost)

	return fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
	config.PGHost, config.PGUser, config.PGPassword, config.PGPort, config.PGDatabase)
}

// ReadConfig reads the configuration from a JSON file
func ReadConfig() Config {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	var config Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Config", config)

	return config
}

// Healthz returns a 200 OK status code to indicate that the service is healthy and ready to serve requests
func Healthz(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.Background(),  GetPGConfig())
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

func homeLink(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the Assessor API! Made by Kyler Chin https://github.com/kylerchin")
}

func main() {

	// Attempt to connect to the database with retry logic
	var conn *pgx.Conn
	var err error
	maxAttempts := 5
	attempt := 1
	for {
		conn, err = pgx.Connect(context.Background(), GetPGConfig())
		if err == nil {
			break
		}
		if attempt == maxAttempts {
			log.Fatalf("Could not connect to database after %d attempts: %v", maxAttempts, err)
		}
		log.Printf("Failed to connect to database on attempt %d: %v. Retrying in 5 seconds...", attempt, err)
		time.Sleep(5 * time.Second)
		attempt++
	}

	defer conn.Close(context.Background())

	http.HandleFunc("/", homeLink)
	http.HandleFunc("/getcoords", GetParcelCoords)
	http.HandleFunc("/healthz", Healthz)
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}