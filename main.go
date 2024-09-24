package main

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	podcastCollection *mongo.Collection
	episodeCollection *mongo.Collection
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	podcastCollection = client.Database("podgo").Collection("podcasts")
	episodeCollection = client.Database("podgo").Collection("episodes")

	r := gin.Default()

	r.GET("/podcast/random", getPodcastRandom)
	r.GET("/podcast/all", getPodcastAll)
	r.GET("/podcast/single/:podlisturl", getPodcastSingle)

	r.GET("/episode/all", getEpisodeAll)
	r.GET("/episode/podcast/:podcasturl", getEpisodePodcast)
	r.GET("/episode/single/:podcasturl/:podlisturl", getEpisodeSingle)

	if err := r.Run(":3007"); err != nil {
		log.Fatal(err)
	}
}

func getPodcastRandom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{{{"$sample", bson.D{{"size", 1}}}}}
	cursor, err := podcastCollection.Aggregate(ctx, pipeline)
	if handleError(err, c) {
		return
	}
	defer cursor.Close(ctx)

	var podcasts []Podcast
	if err = cursor.All(ctx, &podcasts); err != nil {
		handleError(err, c)
		return
	}
	if len(podcasts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Podcast not found"})
		return
	}
	c.JSON(http.StatusOK, podcasts[0])
}

func getPodcastAll(c *gin.Context) {
	pageSize := int64(12)
	page := getPage(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := podcastCollection.CountDocuments(ctx, bson.D{})
	if handleError(err, c) {
		return
	}
	lastPage := int64(math.Ceil(float64(count) / float64(pageSize)))
	if page > lastPage {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page does not exist"})
		return
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "podlistUrl", Value: 1}}).
		SetSkip(pageSize * (page - 1)).
		SetLimit(pageSize)
	cursor, err := podcastCollection.Find(ctx, bson.D{}, opts)
	if handleError(err, c) {
		return
	}
	defer cursor.Close(ctx)

	var podcasts []Podcast
	if err = cursor.All(ctx, &podcasts); err != nil {
		handleError(err, c)
		return
	}
	c.JSON(http.StatusOK, PodcastListPage{
		Page:     page,
		LastPage: lastPage,
		PageSize: pageSize,
		AllCount: count,
		Podcasts: podcasts,
	})
}

func getPodcastSingle(c *gin.Context) {
	podlistURL := c.Param("podlisturl")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var podcast Podcast
	err := podcastCollection.FindOne(ctx, bson.M{"podlistUrl": podlistURL}).Decode(&podcast)
	if handleError(err, c) {
		return
	}
	c.JSON(http.StatusOK, podcast)
}

func getEpisodeAll(c *gin.Context) {
	pageSize := int64(12)
	page := getPage(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := episodeCollection.CountDocuments(ctx, bson.D{})
	if handleError(err, c) {
		return
	}
	lastPage := int64(math.Ceil(float64(count) / float64(pageSize)))
	if page > lastPage {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page does not exist"})
		return
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "published", Value: -1}}).
		SetSkip(pageSize * (page - 1)).
		SetLimit(pageSize)
	cursor, err := episodeCollection.Find(ctx, bson.D{}, opts)
	if handleError(err, c) {
		return
	}
	defer cursor.Close(ctx)

	var episodes []Episode
	if err = cursor.All(ctx, &episodes); err != nil {
		handleError(err, c)
		return
	}
	c.JSON(http.StatusOK, EpisodeListPage{
		Page:     page,
		LastPage: lastPage,
		PageSize: pageSize,
		AllCount: count,
		Episodes: episodes,
	})
}

func getEpisodePodcast(c *gin.Context) {
	pageSize := int64(12)
	page := getPage(c)
	podcastUrl := c.Param("podcasturl")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"podcastUrl": podcastUrl}
	count, err := episodeCollection.CountDocuments(ctx, filter)
	if handleError(err, c) {
		return
	}
	lastPage := int64(math.Ceil(float64(count) / float64(pageSize)))
	if page > lastPage {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page does not exist"})
		return
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "published", Value: -1}}).
		SetSkip(pageSize * (page - 1)).
		SetLimit(pageSize)
	cursor, err := episodeCollection.Find(ctx, filter, opts)
	if handleError(err, c) {
		return
	}
	defer cursor.Close(ctx)

	var episodes []Episode
	if err = cursor.All(ctx, &episodes); err != nil {
		handleError(err, c)
		return
	}
	c.JSON(http.StatusOK, EpisodeListPage{
		Page:     page,
		LastPage: lastPage,
		PageSize: pageSize,
		AllCount: count,
		Episodes: episodes,
	})
}

func getEpisodeSingle(c *gin.Context) {
	podlistUrl := c.Param("podlisturl")
	podcastUrl := c.Param("podcasturl")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"podlistUrl": podlistUrl, "podcastUrl": podcastUrl}
	var episode Episode
	err := episodeCollection.FindOne(ctx, filter).Decode(&episode)
	if handleError(err, c) {
		return
	}
	c.JSON(http.StatusOK, episode)
}

func getPage(c *gin.Context) int64 {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func handleError(err error, c *gin.Context) bool {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return true
	}
	return false
}
