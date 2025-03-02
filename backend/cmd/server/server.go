package main

import (
	"fmt"
	"log"
	"sportbetsim/internal/auth"
	"sportbetsim/internal/config"
	"sportbetsim/internal/handlers"
	"sportbetsim/internal/models"
	"sportbetsim/internal/repository"
	"sportbetsim/internal/router"
	"sportbetsim/internal/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite" // Alterado de postgres para sqlite
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Falha ao carregar as configurações")
	}

	dbDSN := cfg.Database.DSN

	// Conexão com o SQLite
	db, err := gorm.Open(sqlite.Open(dbDSN), &gorm.Config{})
	if err != nil {
		panic("Falha ao conectar ao banco de dados")
	}

	// db.Migrator().DropTable(&models.User{})

	// Migrar tabelas
	db.AutoMigrate(&models.User{}, &models.Match{}, &models.Bet{})

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	result := db.Where(models.User{Username: "admin"}).FirstOrCreate(&models.User{Username: "admin", Email: "admin@example.com", DisplayName: "Administrator", PasswordHash: string(passwordHash), Coins: 999999, Role: "admin"})
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	fmt.Printf("Usuário criado com sucesso - %d\n", result.RowsAffected)

	// Inicializar serviços e repositórios
	userRepo := repository.NewUserRepository(db)
	tokenService := auth.NewTokenService(cfg)
	authService := service.NewAuthService(userRepo, tokenService)
	authHandler := handlers.NewAuthHandler(authService)

	r := router.SetupRouter(authHandler, tokenService)

	// Iniciar servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
