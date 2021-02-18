package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var podcastCollection *mongo.Collection
var episodeCollection *mongo.Collection

func main() {
	clientOptions :=
		options.Client().ApplyURI("mongodb://localhost")
	client, mongoErr := mongo.Connect(context.TODO(), clientOptions)
	if mongoErr != nil {
		panic(mongoErr)
	}
	defer client.Disconnect(context.TODO())
	podcastCollection = client.Database("podgo").Collection("podcasts")
	episodeCollection = client.Database("podgo").Collection("episodes")

	r := gin.Default()

	r.GET("/podcast/random", getPodcastRandom)
	r.GET("/podcast/all/:page", getPodcastAll)
	r.GET("/podcast/single/:podlisturl", getPodcastSingle)

	r.GET("/episode/all/:page", getEpisodeAll)
	r.GET("/episode/podcast/:podcasturl/:page", getEpisodePodcast)
	r.GET("/episode/single/:podcasturl/:podlisturl", getEpisodeSingle)

	ginErr := r.Run(":3007")
	if ginErr != nil {
		log.Fatal(ginErr)
	}
}

func getPodcastRandom(c *gin.Context) {
	var podcast Podcast
	opts1 := options.Count().SetMaxTime(5 * time.Second)
	count, err := podcastCollection.CountDocuments(context.TODO(), bson.D{}, opts1)
	if err != nil {
		log.Fatal(err)
	}
	opts2 := options.Find().SetSkip(rand.Int63n(count)).SetLimit(1)
	cursor, err := podcastCollection.Find(context.TODO(), bson.D{}, opts2)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		defer cursor.Close(context.TODO())
		for cursor.Next(context.TODO()) {
			if err = cursor.Decode(&podcast); err != nil {
				log.Fatal(err)
			}
		}
		c.JSON(200, podcast)
	}
}

func getPodcastAll(c *gin.Context) {
	pageSize := int64(12)
	page, _ := strconv.ParseInt(c.Param("page"), 10, 32)
	count, _ := podcastCollection.CountDocuments(context.TODO(), bson.D{})
	lastPage := int64(math.Ceil(float64(count) / float64(pageSize)))
	var podcasts []Podcast
	opts := options.Find().SetSort(bson.D{{"podlistUrl", 1}}).SetSkip(pageSize*page - pageSize).SetLimit(pageSize)
	cursor, err := podcastCollection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		defer cursor.Close(context.TODO())
		if cursor.RemainingBatchLength() < 1 {
			c.JSON(500, gin.H{
				"error": "Page does not exist!",
			})
			return
		}
		for cursor.Next(context.TODO()) {
			var podcast Podcast
			if err = cursor.Decode(&podcast); err != nil {
				log.Fatal(err)
			}
			podcasts = append(podcasts, podcast)
		}
		c.JSON(200, PodcastListPage{Page: page, LastPage: lastPage, PageSize: pageSize, AllCount: count, Podcasts: podcasts})
	}
}

func getPodcastSingle(c *gin.Context) {
	podlistURL := c.Param("podlisturl")
	var podcast Podcast
	err := podcastCollection.FindOne(context.TODO(), bson.M{"podlistUrl": podlistURL}).Decode(&podcast)
	if !checkErrorOne(err, c) {
		c.JSON(200, podcast)
	}
}

func getEpisodeAll(c *gin.Context) {
	pageSize := int64(12)
	page, _ := strconv.ParseInt(c.Param("page"), 10, 32)
	count, _ := episodeCollection.CountDocuments(context.TODO(), bson.D{})
	lastPage := int64(math.Ceil(float64(count) / float64(pageSize)))
	var episodes []Episode
	opts := options.Find().SetSort(bson.D{{"published", -1}}).SetSkip(pageSize*page - pageSize).SetLimit(pageSize)
	cursor, err := episodeCollection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		defer cursor.Close(context.TODO())
		if cursor.RemainingBatchLength() < 1 {
			c.JSON(500, gin.H{
				"error": "Page does not exist!",
			})
			return
		}
		for cursor.Next(context.TODO()) {
			var episode Episode
			if err = cursor.Decode(&episode); err != nil {
				log.Fatal(err)
			}
			episodes = append(episodes, episode)
		}
		c.JSON(200, EpisodeListPage{Page: page, LastPage: lastPage, PageSize: pageSize, AllCount: count, Episodes: episodes})
	}
}

func getEpisodePodcast(c *gin.Context) {
	pageSize := int64(12)
	page, _ := strconv.ParseInt(c.Param("page"), 10, 32)
	podcastUrl := c.Param("podcasturl")
	count, _ := episodeCollection.CountDocuments(context.TODO(), bson.M{"podcastUrl": podcastUrl})
	lastPage := int64(math.Ceil(float64(count) / float64(pageSize)))
	var episodes []Episode
	opts := options.Find().SetSort(bson.D{{"published", -1}}).SetSkip(pageSize*page - pageSize).SetLimit(pageSize)
	cursor, err := episodeCollection.Find(context.TODO(), bson.M{"podcastUrl": podcastUrl}, opts)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		defer cursor.Close(context.TODO())
		cursor.RemainingBatchLength()
		if cursor.RemainingBatchLength() < 1 {
			c.JSON(500, gin.H{
				"error": "Page does not exist!",
			})
			return
		}
		for cursor.Next(context.TODO()) {
			var episode Episode
			if err = cursor.Decode(&episode); err != nil {
				log.Fatal(err)
			}
			episodes = append(episodes, episode)
		}
		c.JSON(200, EpisodeListPage{Page: page, LastPage: lastPage, PageSize: pageSize, AllCount: count, Episodes: episodes})
	}
}

func getEpisodeSingle(c *gin.Context) {
	podlistUrl := c.Param("podlisturl")
	podcastUrl := c.Param("podcasturl")
	var episode Episode
	err := episodeCollection.FindOne(context.TODO(), bson.M{"podlistUrl": podlistUrl, "podcastUrl": podcastUrl}).Decode(&episode)
	if !checkErrorOne(err, c) {
		c.JSON(200, episode)
	}
}

func checkErrorOne(err error, c *gin.Context) bool {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		return true
	} else {
		return false
	}
}
