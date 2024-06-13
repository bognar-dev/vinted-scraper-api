package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	vinted_scraper "vinted-scraper/internal/vinted-scraper"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// Exec executes a SQL query with the provided arguments and returns the result.
	// It is safe against SQL injection when used with parameter placeholders.
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// BeginTx starts a new database transaction with the specified options.
	// The transaction is bound to the context passed.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// QueryRow executes a query that is expected to return at most one row.
	// The result is scanned into the provided destination variables.
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Query executes a query that returns multiple rows.
	// It is safe against SQL injection when used with parameter placeholders.
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// Prepare creates a prepared statement for repeated use.
	// A prepared statement takes parameters and is safe against SQL injection.
	Prepare(ctx context.Context, query string) (*sql.Stmt, error)

	AddItems(items []vinted_scraper.Item, topic string) error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) AddItems(items []vinted_scraper.Item, topic string) error {
	// Begin a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Insert the topic into the Topic table (if it doesn't already exist)
	var topicID int
	err = s.db.QueryRow("INSERT INTO Topic (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = $1 RETURNING id", topic).Scan(&topicID)
	if err != nil {
		// If topic insertion fails, rollback the transaction and return the error
		tx.Rollback()
		return fmt.Errorf("error inserting topic %s: %v", topic, err)
	}

	// Loop through each item and insert photos and thumbnails
	for _, item := range items {
		// Insert photos
		photoID, err := s.insertPhoto(tx, item.Photo)
		if err != nil {
			// If photo insertion fails, rollback the transaction and return the error
			tx.Rollback()
			return fmt.Errorf("error inserting photo for item %d: %v", item.ID, err)
		}

		// Insert thumbnails for each photo
		for _, thumbnail := range item.Photo.Thumbnails {
			_, err := tx.Exec("INSERT INTO Thumbnails (Type, URL, Width, Height, photo_id) VALUES ($1, $2, $3, $4, $5)",
				thumbnail.Type, thumbnail.URL, thumbnail.Width, thumbnail.Height, item.Photo.ID)
			if err != nil {
				// If thumbnail insertion fails, rollback the transaction and return the error
				tx.Rollback()
				return fmt.Errorf("error inserting thumbnail for photo %d: %v", item.Photo.ID, err)
			}
		}

		// Insert item into Item table and into Item_Topic
		_, err = tx.Exec(`INSERT INTO Item (
			id, title, price, is_visible, discount, currency, brand_title,
			user_id, url, promoted, photo_id, favourite_count, is_favourite,
			badge, conversion, service_fee, total_item_price, total_item_price_rounded,
			view_count, size_title, content_source, status, icon_badges, search_tracking_params,topic_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
		) ON CONFLICT (id) DO UPDATE SET title = $2, price = $3, is_visible = $4, discount = $5, currency = $6, brand_title = $7, user_id = $8, url = $9, promoted = $10, photo_id = $11, favourite_count = $12, is_favourite = $13, badge = $14, conversion = $15, service_fee = $16, total_item_price = $17, total_item_price_rounded = $18, view_count = $19, size_title = $20, content_source = $21, status = $22, icon_badges = $23, search_tracking_params = $24, topic_id = $25`,
			item.ID, item.Title, item.Price, item.IsVisible, item.Discount, item.Currency, item.BrandTitle,
			item.User.ID, item.URL, item.Promoted, photoID, item.FavouriteCount, item.IsFavourite,
			item.Badge, item.Conversion, item.ServiceFee, item.TotalItemPrice, item.TotalItemPriceRounded,
			item.ViewCount, item.SizeTitle, item.ContentSource, item.Status, nil, nil, topicID)
		if err != nil {
			// If item insertion fails, rollback the transaction and return the error
			tx.Rollback()
			return fmt.Errorf("error inserting item %d: %v", item.ID, err)
		}
	}

	// Commit the transaction if everything is successful
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Helper function to insert a photo and return its ID
func (s *service) insertPhoto(tx *sql.Tx, photo vinted_scraper.Photo) (int, error) {
	var photoID int
	fmt.Println("Inserting photo:", photo)
	// Insert photo into the Photos table
	err := tx.QueryRow("INSERT INTO Photos (id, ImageNo, Width, Height, DominantColor, DominantColorOpaque, URL, IsMain,  IsSuspicious, FullSizeURL, IsHidden) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (id) DO UPDATE SET id = $1 RETURNING id",
		photo.ID, photo.ImageNo, photo.Width, photo.Height, photo.DominantColor, photo.DominantColorOpaque, photo.URL, photo.IsMain, photo.IsSuspicious, photo.FullSizeURL, photo.IsHidden).Scan(&photoID)
	if err != nil {
		return 0, err
	}

	return photoID, nil
}

// Exec executes a SQL query with the given arguments within the provided context.
// It returns the result of the execution, such as the number of affected rows.
func (s *service) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

// BeginTx starts a new transaction with the given transaction options within the provided context.
// It returns a transaction handle to be used for executing statements and committing or rolling back.
func (s *service) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}

// QueryRow executes a SQL query that is expected to return at most one row,
// scanning the result into the provided destination variables.
func (s *service) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

// Query executes a SQL query with the given arguments within the provided context.
// It returns a result set containing multiple rows, which must be iterated over.
func (s *service) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

// Prepare creates a new prepared statement for the given query within the provided context.
// Prepared statements can be reused and are safe against SQL injection.
func (s *service) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.db.PrepareContext(ctx, query)
}
