package main
import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
)
func main() {
	dsn := "host=postgres user=postgres password=evoai_dev_password dbname=evo_community port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { fmt.Println(err); return }
	var agents []models.Agent
	db.Find(&agents)
	for _, a := range agents {
		fmt.Printf("Agent: %s, Config: %s\n", a.Name, string(a.Config))
	}
}
