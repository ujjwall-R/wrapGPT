package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Userid string `json:"userid"`
	Prompt string `json:"prompt"`
}

func DB() *mongo.Collection {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	coll := client.Database("sample_mflix").Collection("movies")
	return coll
}

func UpdateUser(coll *mongo.Collection, user User, updatedPrompt string) error {
	filter := bson.D{{"userid", user.Userid}}
	update := bson.D{{"$set", bson.D{{"prompt", updatedPrompt}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func SearchUser(coll *mongo.Collection, userid string) (User, error) {
	var result User
	err := coll.FindOne(context.TODO(), bson.D{{"userid", userid}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			fmt.Printf("No document was found with the userid %v\n", userid)
			return result, err
		}
	}

	return result, nil
}

func AddUser(coll *mongo.Collection, user User) (User, error) {
	_, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
