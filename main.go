package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	Name      string          `json:"name"`
	Email     string          `json:"email" gorm:"unique"`
	Password  string          `json:"-"`
	Role      string          `json:"role" gorm:"default:user"`
	DeletedAt gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
}

func initDatabase() {
	godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	database.AutoMigrate(&User{})
	db = database
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateJWT(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}


func RegisterUser(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	hashedPassword, err := hashPassword(data["password"])
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	role := "user"
	if data["role"] == "admin" {
		role = "admin"
	}

	user := User{Name: data["name"], Email: data["email"], Password: hashedPassword, Role: role}
	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.JSON(user)
}


func LoginUser(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user User
	if err := db.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if !checkPasswordHash(data["password"], user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
	}

	token, err := generateJWT(user.Email, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func AdminMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Admins only"})
	}

	return c.Next()
}


func GetUsers(c *fiber.Ctx) error {
	var users []User
	db.Find(&users)
	return c.JSON(users)
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	user.Name = data["name"]
	user.Email = data["email"]
	db.Save(&user)

	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}


func main() {
	app := fiber.New()
	initDatabase()

	app.Post("/register", RegisterUser)
	app.Post("/login", LoginUser)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	app.Get("/users", GetUsers)
	app.Get("/users/:id", GetUserByID)
	app.Put("/users/:id", UpdateUser)
	app.Delete("/users/:id", jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}), AdminMiddleware, DeleteUser)	

	log.Fatal(app.Listen(":3000"))
}