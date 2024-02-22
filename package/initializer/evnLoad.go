package initializer

import (
	"fmt"

	"github.com/joho/godotenv"
)

func EnvLoad() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error to load env..............")
	}
}
