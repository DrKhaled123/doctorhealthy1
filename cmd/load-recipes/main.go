package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"api-key-generator/internal/services"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		dbPath     = flag.String("db", "data/apikeys.db", "Database path")
		recipesDir = flag.String("dir", ".", "Directory containing recipe JSON files")
		initDB     = flag.Bool("init", false, "Initialize database schema")
		stats      = flag.Bool("stats", false, "Show recipe statistics")
	)
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Error closing database: %v", cerr)
		}
	}()

	loader := services.NewRecipeLoader(db)

	if *initDB {
		fmt.Println("Initializing database schema...")
		if err := loader.InitializeDatabase(); err != nil {
			log.Printf("Failed to initialize database: %v", err)
			os.Exit(1)
		}
		fmt.Println("Database schema initialized successfully")
	}

	if *stats {
		stats, err := loader.GetRecipeStats()
		if err != nil {
			log.Printf("Failed to get stats: %v", err)
			os.Exit(1)
		}

		fmt.Println("Recipe Statistics:")
		fmt.Printf("Total recipes: %d\n", stats["total"])
		fmt.Printf("Arabian Gulf: %d\n", stats["arabian_gulf"])
		fmt.Printf("Shami: %d\n", stats["shami"])
		fmt.Printf("Egyptian: %d\n", stats["egyptian"])
		fmt.Printf("Moroccan: %d\n", stats["moroccan"])
		return
	}

	fmt.Printf("Loading recipes from directory: %s\n", *recipesDir)
	if err := loader.LoadAllRecipes(*recipesDir); err != nil {
		log.Printf("Failed to load recipes: %v", err)
		os.Exit(1)
	}

	fmt.Println("Recipe loading completed successfully!")
}
