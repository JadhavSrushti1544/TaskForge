package main
import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)
func main() {
	godotenv.Load()
	fmt.Printf("DSN: host=%q port=%q user=%q password=%q dbname=%q\n",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
}
