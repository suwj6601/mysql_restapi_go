package db

import (
	"database/sql"
	"fmt"
	"go-mysql-restapi/config"
	"log"
)

func InitDB() *sql.DB {
	var err error
	var db *sql.DB

	// Initialize the database
	dsn := config.GetDsnDB()
	fmt.Printf("MYSQL Connection: %s\n", dsn)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	CreateInitTable(db)
	fmt.Println("Connected to MySQL database")
	return db
}

func GetDbConnection() (*sql.DB, error) {
	dsn := config.GetDsnDB()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}

func CreateInitTable(db *sql.DB) error {
	// Create the users table
	_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS users (
					id INT AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(100) NOT NULL,
					age INT NOT NULL,
					email VARCHAR(100) NOT NULL UNIQUE,
					password VARCHAR(100) NOT NULL,
					balance INT NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
	`)

	if err != nil {
		log.Fatal("Error when creating init table Users", err)
		return err
	}

	// Create the user_role table
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS user_role (
					id INT AUTO_INCREMENT PRIMARY KEY,
					user_id INT NOT NULL,
					is_admin BOOLEAN NOT NULL,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			);
	`)
	if err != nil {
		log.Fatal("Error when creating init table User_role", err)

		return err
	}

	// check user admin exist

	isAdminExist, err := CheckAdminExist(db)
	if err != nil {
		log.Fatal("Error when checking if admin exists", err)
		return err
	}

	if isAdminExist {
		log.Println("Admin already exists")
		return nil
	} else {

		_, err = db.Exec(`
			INSERT INTO users (id, name, age, email, password, balance, created_at)
			VALUES (1, "admin", 0, "admin@gmail.com", "$2a$10$ILqpfcd6EIqezw0.Xxqn4eHatBip6Qv68EXn0fC3qm3Xa3lMtgnGm", 1000, NOW());
`)
		if err != nil {
			// Handle the error
			log.Fatalf("Error inserting user: %v", err)
		}

		// Second insert
		_, err = db.Exec(`
			INSERT INTO user_role (id, user_id, is_admin)
			VALUES (1, 1, 1);
`)
		if err != nil {
			// Handle the error
			log.Fatalf("Error inserting user role: %v", err)
		}

	}

	if err != nil {
		log.Fatal("Error when creating init table Users", err)
		return err
	}

	// Create the products table
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS products (
					id INT AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(100) NOT NULL,
					price INT NOT NULL,
					description TEXT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
	`)
	if err != nil {
		log.Fatal("Error when creating init table Products", err)

		return err
	}

	// Create the orders table
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS orders (
					id INT AUTO_INCREMENT PRIMARY KEY,
					user_id INT NOT NULL,
					product_id INT NOT NULL,
					quantity INT NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					is_done BOOLEAN NOT NULL,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
					FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
			);
	`)
	if err != nil {
		log.Fatal("Error when creating init table Orders", err)

		return err
	}

	log.Default().Println("Init table created successfully")
	return nil
}

func CheckAdminExist(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE name = 'admin'").Scan(&count)
	if err != nil {
		log.Fatal("Error when checking if admin exists", err)
		return false, err
	}
	return count > 0, nil
}
