package store

import (
	"context"

	"github.com/Hassall/transit/pkg/request"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

// URLStore interface for URLs
type URLStore interface {
	Connect() error
	Close()
	StoreUrls([]request.URLRequest)
	RetrieveUrls() []request.URLRequest
}

// DB implements URLStore
type DB struct {
	ctx  context.Context
	conn *pgx.Conn
}

// Connect connects to db
func (db *DB) Connect() error {
	log.Info("Connecting to URLStore...")
	db.ctx = context.Background()
	connStr := "postgres://postgres:root@localhost:5432/postgres"
	var err error
	db.conn, err = pgx.Connect(db.ctx, connStr)
	return err
}

// Close closes db connection
func (db DB) Close() {
	log.Info("Closing URLStore.")
	db.conn.Close(db.ctx)
}

// StoreUrls stores url to database
func (db DB) StoreUrls(urls []request.URLRequest) {
	insertQuery := "INSERT INTO url_requests VALUES ($1);"
	for _, url := range urls {
		_, err := db.conn.Exec(db.ctx, insertQuery, url.URL)
		if err != nil {
			log.Error("Failed to insert item into DB", err)
		} else {
			log.Debug("Succesfully added item to DB: ", url)
		}
	}
}

// RetrieveUrls from database
func (db DB) RetrieveUrls() []request.URLRequest {
	var urls []request.URLRequest
	sqlRequestQuery := "SELECT * FROM url_requests;"
	rows, err := db.conn.Query(db.ctx, sqlRequestQuery)
	if err != nil {
		log.Error("Failed to execute queery on database: ", sqlRequestQuery, err)
	}

	defer rows.Close()

	for rows.Next() {
		var url request.URLRequest
		if err := rows.Scan(&url.URL); err != nil {
			log.Error("Failed to read row from database", err)
		}
		urls = append(urls, url)
	}

	if rows.Err() != nil {
		log.Error("Failed to read rows from database", err)
	}

	return urls
}
