package main

import (
	"github.com/codepix/imersao/codepix-go/application/grpc"
	"github.com/codepix/imersao/codepix-go/infra/db"
	"github.com/jinzhu/gorm"
	"os"
)

var database *gorm.DB

func main() {
	database = db.ConnectDB(os.Getenv("env"))
	grpc.StartGrpcServer(database, 50051)
}