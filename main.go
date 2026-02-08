package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Model with Unique Phone constraint
type Invitado struct {
	gorm.Model
	Nombre    string `json:"nombre"`
	Apellido  string `json:"apellido"`
	Phone     string `json:"phone" gorm:"uniqueIndex"`
	Respuesta string `json:"respuesta"`
}

var db *gorm.DB

func main() {
	// 1. Load configuration from Environment Variables
	// Use "PORT" if defined, otherwise default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "invitados.db"
	}

	// 2. Initialize Database
	var err error
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Critical: Could not connect to DB: %v", err)
	}
	db.AutoMigrate(&Invitado{})

	// 3. Setup Gin
	r := gin.Default()

	// Serve Frontend
	// r.Static("/static", "./public")

	// API Routes
	api := r.Group("/api")
	{
		api.GET("/search", searchInvitado)
		api.PATCH("/invitados/:id", updateRespuesta)
		api.POST("/upload", uploadCSV)
	}

		// 2. Serve the main HTML file at the root "/"
	// StaticFile is specific and does not use a wildcard, so it won't panic
	r.StaticFile("/", "./public/index.html")

	// 3. Serve other folders specifically if they exist (e.g., css, js, images)
	// This maps http://localhost:8080/css/style.css to ./public/css/style.css
	r.Static("/css", "./public/css")
	r.Static("/js", "./public/js")

	// 4. OPTIONAL: Catch-all for other files in public
	// If you have files directly in /public (like favicon.ico), 
    // use NoRoute to serve them without a wildcard conflict.
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Check if the requested file exists in the public folder
		if _, err := os.Stat("./public" + path); err == nil {
			c.File("./public" + path)
			return
		}
		// If nothing matches, send them back to index (Standard for Single Page Apps)
		c.File("./public/index.html")
	})


	log.Printf("🚀 Server starting on http://localhost:%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Critical: Server failed: %v", err)
	}
}

// --- HANDLERS ---

func searchInvitado(c *gin.Context) {
	nombre := c.Query("nombre")
	var invitados []Invitado
	db.Where("nombre LIKE ?", "%"+nombre+"%").Find(&invitados)
	c.JSON(http.StatusOK, invitados)
}

func updateRespuesta(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Respuesta string `json:"respuesta" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Respuesta is required"})
		return
	}

	var invitado Invitado
	if err := db.First(&invitado, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guest not found"})
		return
	}

	db.Model(&invitado).Update("respuesta", input.Respuesta)
	c.JSON(http.StatusOK, invitado)
}

func uploadCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	f, _ := file.Open()
	defer f.Close()

	reader := csv.NewReader(f)
	if _, err := reader.Read(); err != nil { // Skip Header
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV"})
		return
	}

	var added, skipped int64
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		invitado := Invitado{
			Nombre: record[0], Apellido: record[1],
			Phone: record[2], Respuesta: record[3],
		}

		// Insert or ignore if Phone exists
		result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&invitado)

		if result.RowsAffected > 0 {
			added++
		} else {
			skipped++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"added":   added,
		"skipped": skipped,
	})
}
