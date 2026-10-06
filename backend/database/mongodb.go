package database

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"tracker/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var DB *mongo.Database

func Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

var credRe = regexp.MustCompile(`://[^@/\s]+@`)

func redact(s string) string {
	return strings.ReplaceAll(credRe.ReplaceAllString(s, "://***@"), config.C.MongoURI, "<redacted-uri>")
}

func hint(e string) string {
	e = strings.ToLower(e)
	switch {
	case strings.Contains(e, "auth"):
		return "wrong database user/password. Check Atlas > Database Access; URL-encode special characters in the password (@ -> %40, : -> %3A, / -> %2F, # -> %23, % -> %25)."
	case strings.Contains(e, "no such host") || strings.Contains(e, "srv"):
		return "DNS could not resolve the cluster (SRV lookup). Check internet/VPN/DNS, or use Atlas's non-SRV 'standard connection string'."
	default:
		return "usually Atlas Network Access does not include your current IP. Atlas > Security > Network Access > Add IP Address > 'Add Current IP Address' (or 0.0.0.0/0 for local dev). Also check the cluster is not paused and that your network allows outbound port 27017."
	}
}

func ConnectMongoDB() (*mongo.Client, *mongo.Database, error) {
	opts := options.Client().ApplyURI(config.C.MongoURI).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetServerSelectionTimeout(10 * time.Second).SetConnectTimeout(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid connection settings: %s", redact(err.Error()))
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		client.Disconnect(context.Background())
		return nil, nil, fmt.Errorf("ping failed: %s\n  hint: %s", redact(err.Error()), hint(err.Error()))
	}
	Client, DB = client, client.Database(config.C.MongoDB)
	return Client, DB, nil
}

func InitIndexes() {
	type spec struct {
		col    string
		keys   bson.D
		unique bool
	}
	asc := func(f ...string) bson.D {
		d := bson.D{}
		for _, k := range f {
			d = append(d, bson.E{Key: k, Value: 1})
		}
		return d
	}
	specs := []spec{
		{"users", asc("email"), true},
		{"tasks", asc("userId", "status"), false}, {"tasks", asc("userId", "dueDate"), false},
		{"habits", asc("userId"), false},
		{"habit_checkins", asc("userId", "habitId", "date"), true},
		{"events", asc("userId", "date"), false}, {"reminders", asc("userId", "date"), false},
		{"diary", asc("userId", "date"), false},
		{"notes", bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}, false},
	}
	for _, s := range specs {
		ctx, cancel := Ctx()
		_, err := DB.Collection(s.col).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: s.keys, Options: options.Index().SetUnique(s.unique)})
		cancel()
		if err != nil {
			log.Printf("index on %s failed: %v", s.col, err)
		}
	}
}
