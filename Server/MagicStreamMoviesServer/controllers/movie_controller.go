// Gin expects a handler of type `gin.HandlerFunc`, i.e. `func(*gin.Context)`.
// By writing `func AddMovie() gin.HandlerFunc { return func(c *gin.Context){...} }`
//We create a *factory function* that returns the actual request handler.
//    - Outer function (AddMovie): used for setup or dependency injection
//      (e.g., pass DB connection, config, etc.)
//    - Inner anonymous func: the real handler that runs for each request.

package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Shawnq-create/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/Shawnq-create/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies", database.Client)
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		// Logic to fetch movies from the database will go here
		var movies []models.Movie
		// Example: Fetch movies from MongoDB collection and populate the 'movies' slice

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movies from database"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding movies from database"})
			return
		}
		c.JSON(http.StatusOK, movies)
	}

}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()
		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieId}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie from database"})
			return
		}

		c.JSON(http.StatusOK, movie)

	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logic to add a new movie to the database will go here

		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting movie into database"})
			return
		}
		c.JSON(http.StatusOK, result)

	}
}
