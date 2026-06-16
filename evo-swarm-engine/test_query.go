package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Agent struct {
	ID   string
	Name string
	Config string
}

func main() {
	dsn := "host=localhost user=postgres password=evoai_dev_password dbname=evo_community port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var agents []Agent
	db.Table("evo_core_agents").Select("id, name, config::text as config").Find(&agents)
	
	fmt.Printf("Found %d agents\n", len(agents))
	for _, a := range agents {
		fmt.Printf("Agent ID: %s, Name: %s\nConfig: %s\n\n", a.ID, a.Name, a.Config)
	}
}
