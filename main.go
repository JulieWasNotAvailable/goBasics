package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JulieWasNotAvailable/goBasics/models"
	"github.com/JulieWasNotAvailable/goBasics/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Beat struct {
	Author string 
	Title string
	LicenseName string
}

type Repository struct{
	DB *gorm.DB
}

func(r *Repository) CreateBeat (context *fiber.Ctx) error {
	beat:= Beat{}

	err := context.BodyParser(&beat)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request faield"})
		return err
	}

	err = r.DB.Create(&beat).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not create beat"})
			return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"beat has been added"})

	return nil
}

func (r *Repository) GetBeats(context *fiber.Ctx) error {
	beatModels := &[]models.Beats{}

	err:= r.DB.Find(beatModels).Error
	if err!=nil{
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not get books"})
			return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"beats fetched successfully",
		"data": beatModels})

	return nil
}

func (r *Repository) DeleteBeats (context *fiber.Ctx) error {
	//with fiber you can easily access params, request and response data
	//but you should be able to do stuff without fiber
	beatModel := models.Beats{}
	id := context.Params("id")

	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}
	

	err := r.DB.Delete(beatModel, id).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
		&fiber.Map{"message": "couldn't delete"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book deleted successfully"})

	return nil
}

func (r *Repository) GetBeatByID(context *fiber.Ctx) error {
	id := context.Params("id")
	beatModel :=  &models.Beats{}
	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"id cannot be empty",
		})
	}
	fmt.Println("the id is: ", id)

	err:= r.DB.Where("id = ?", id).First(beatModel).Error
	if err !=nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the beat"})
			return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "beat is fetched successfully",
		"data": beatModel,
	})

	return nil
}

func(r *Repository) SetupRoutes (app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_beat", r.CreateBeat)
	api.Delete("/delete_beat/:id", r.DeleteBeats)
	api.Get("/get_beat/:id", r.GetBeatByID)
	api.Get("/beats", r.GetBeats)

}

func main() {

	fmt.Println(os.Getwd())
	
	err:= godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	
	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv ("DB_PASS"),
		User: os.Getenv ("DB_USER"),
		SSLMode: os.Getenv ("DB_SSLMODE"),
		DBName: os.Getenv ("DB_NAME"),
	}

	db, err := storage.NewConnection(config) 

	if err != nil{
		log.Fatal("could not reach the DB")
	}

	err = models.MigrateBeats(db)
	if err != nil {
		log.Fatal("could not migrate")
	}

	repo := Repository{
		DB: db,
	}
	app := fiber.New()
	// app.Use(cors.New())
	repo.SetupRoutes(app)
	app.Listen(":8080")
}