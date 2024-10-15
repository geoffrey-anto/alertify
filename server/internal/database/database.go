package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"server/internal/models"
	"strconv"
	"time"

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

	// SaveUser saves a user to the database.
	SaveUser(email, pass, fname, lname string) error

	// GetUser retrieves a user from the database.
	GetUser(email string) (models.User, error)

	// SaveReminder saves a reminder to the database.
	SaveReminder(userId int, name, status, description, category, reminderInterval, reminderEnd string) error

	// GetReminderById retrieves a reminder from the database by its ID.
	GetReminderById(id int) (models.Reminder, error)

	// GetAllRemindersForUser retrieves all reminders for a user from the database.
	GetAllRemindersForUser(userId int) ([]models.Reminder, error)

	// GetAllReminders retrieves all reminders from the database.
	GetAllReminders() ([]models.Reminder, error)
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func initDB(db *sql.DB) {
	fmt.Println("Initializing database...")
	// create users table if it doesn't exist
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		pass TEXT NOT NULL,
		fname TEXT NOT NULL,
		lname TEXT NOT NULL
	)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reminders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		name TEXT NOT NULL,
		status TEXT NOT NULL,
		description TEXT NOT NULL,
		category TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		reminder_interval TEXT NOT NULL,
		reminder_end TEXT NOT NULL
	)`)

	if err != nil {
		log.Fatal(err)
	}
}

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

	initDB(db)

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

func (s *service) SaveUser(email, pass, fname, lname string) error {
	_, err := s.db.Exec("INSERT INTO users (email, pass, fname, lname) VALUES ($1, $2, $3, $4)", email, pass, fname, lname)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUser(email string) (models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Pass, &user.Fname, &user.Lname)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *service) SaveReminder(userId int, name, status, description, category, reminderInterval, reminderEnd string) error {
	_, err := s.db.Exec("INSERT INTO reminders (user_id, name, status, description, category, reminder_interval, reminder_end) VALUES ($1, $2, $3, $4, $5, $6, $7)", userId, name, status, description, category, reminderInterval, reminderEnd)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetReminderById(id int) (models.Reminder, error) {
	var reminder models.Reminder
	err := s.db.QueryRow("SELECT * FROM reminders WHERE id = $1", id).Scan(&reminder.ID, &reminder.UserID, &reminder.Name, &reminder.Status, &reminder.Description, &reminder.Category, &reminder.CreatedAt, &reminder.UpdatedAt, &reminder.ReminderInterval, &reminder.ReminderEnd)
	if err != nil {
		return reminder, err
	}
	return reminder, nil
}

func (s *service) GetAllRemindersForUser(userId int) ([]models.Reminder, error) {
	rows, err := s.db.Query("SELECT * FROM reminders WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []models.Reminder = []models.Reminder{}
	for rows.Next() {
		var reminder models.Reminder
		err := rows.Scan(&reminder.ID, &reminder.UserID, &reminder.Name, &reminder.Status, &reminder.Description, &reminder.Category, &reminder.CreatedAt, &reminder.UpdatedAt, &reminder.ReminderInterval, &reminder.ReminderEnd)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	return reminders, nil
}

func (s *service) GetAllReminders() ([]models.Reminder, error) {
	rows, err := s.db.Query("SELECT * FROM reminders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []models.Reminder = []models.Reminder{}
	for rows.Next() {
		var reminder models.Reminder
		err := rows.Scan(&reminder.ID, &reminder.UserID, &reminder.Name, &reminder.Status, &reminder.Description, &reminder.Category, &reminder.CreatedAt, &reminder.UpdatedAt, &reminder.ReminderInterval, &reminder.ReminderEnd)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	return reminders, nil
}
