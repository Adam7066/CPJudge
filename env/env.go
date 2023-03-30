package env

import (
	"os"

	"github.com/Adam7066/golang/log"
	"github.com/joho/godotenv"
)

var HWInfo = getHWEnv()

func getHWEnv() map[string]string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error.Fatal("Error loading .env file")
	}
	ret := make(map[string]string)
	ret["HWZip"] = os.Getenv("HWZip")
	return ret
}
