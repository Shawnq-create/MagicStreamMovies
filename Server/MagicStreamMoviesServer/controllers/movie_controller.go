package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Shawnq-create/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/Shawnq-create/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies", database.Client)

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
