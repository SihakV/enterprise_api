package env

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	Hostname      string
	Port          string
	Name          string
	Password      string
	DB            string
	CloudAccess   string
	CloudSecret   string
	CloudRegion   string
	CloudEndPoint string
	CloudBucket   string
	JwtSecret     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load environment variables")
	}

	//db credential
	Hostname = os.Getenv("HOSTNAME")
	Name = os.Getenv("NAME")
	DB = os.Getenv("DB")
	Password = os.Getenv("PASSWORD")
	Port = os.Getenv("PORT")

	//object storage credential
	CloudAccess = os.Getenv("CLOUD_ACCESS")
	CloudSecret = os.Getenv("CLOUD_SECRET")
	CloudEndPoint = os.Getenv("CLOUD_END_POINT")
	CloudRegion = os.Getenv("CLOUD_REGION")
	CloudBucket = os.Getenv("CLOUD_BUCKET")

	//jwt secret
	JwtSecret = os.Getenv("JWT_SECRET_KEY")
}
