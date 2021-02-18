package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Podcast struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title,omitempty" json:"title"`
	Categories  []string           `bson:"categories,omitempty" json:"categories"`
	Link        string             `bson:"link,omitempty" json:"link"`
	Description string             `bson:"description,omitempty" json:"description"`
	Subtitle    string             `bson:"subtitle,omitempty" json:"subtitle"`
	Owner       PodcastOwner       `bson:"owner,omitempty" json:"owner"`
	Author      string             `bson:"author,omitempty" json:"author"`
	Image       string             `bson:"image,omitempty" json:"image"`
	Feed        string             `bson:"feed,omitempty" json:"feed"`
	PodlistUrl  string             `bson:"podlistUrl,omitempty" json:"podlistUrl"`
	Updated     time.Time          `bson:"updated,omitempty" json:"updated"`
}

type PodcastListPage struct {
	Podcasts []Podcast `json:"podcasts"`
	AllCount int64     `json:"allCount"`
	PageSize int64     `json:"pageSize"`
	Page     int64     `json:"page"`
	LastPage int64     `json:"lastPage"`
}

type Episode struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PodlistUrl   string             `bson:"podlistUrl,omitempty" json:"podlistUrl"`
	PodcastId    primitive.ObjectID `bson:"podcastId,omitempty" json:"podcastId"`
	PodcastUrl   string             `bson:"podcastUrl,omitempty" json:"podcastUrl"`
	PodcastTitle string             `bson:"podcastTitle,omitempty" json:"podcastTitle"`
	PodcastImage string             `bson:"podcastImage,omitempty" json:"podcastImage"`
	Guid         string             `bson:"guid,omitempty" json:"guid"`
	Title        string             `bson:"title,omitempty" json:"title"`
	Published    time.Time          `bson:"published,omitempty" json:"published"`
	Duration     string             `bson:"Duration,omitempty" json:"duration"`
	Summary      string             `bson:"summary,omitempty" json:"summary"`
	Subtitle     string             `bson:"subtitle,omitempty" json:"subtitle"`
	Description  string             `bson:"description,omitempty" json:"description"`
	Image        string             `bson:"image,omitempty" json:"image"`
	Content      string             `bson:"content,omitempty" json:"content"`
	Enclosure    EpisodeEnclosure   `bson:"enclosure,omitempty" json:"enclosure"`
}

type EpisodeListPage struct {
	Episodes []Episode `json:"episodes"`
	AllCount int64     `json:"allCount"`
	PageSize int64     `json:"pageSize"`
	Page     int64     `json:"page"`
	LastPage int64     `json:"lastPage"`
}

type PodcastOwner struct {
	Name  string `bson:"name,omitempty" json:"name"`
	Email string `bson:"email,omitempty" json:"email"`
}

type EpisodeEnclosure struct {
	Filesize string `bson:"filesize,omitempty" json:"filesize"`
	Filetype string `bson:"filetype,omitempty" json:"filetype"`
	Url      string `bson:"url,omitempty" json:"url"`
}
