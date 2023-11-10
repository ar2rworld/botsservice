package bots

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DimaDay   = "DimaDay"
	DimaMonth = "DimaMonth"
)

const NoDocuments   = "mongo: no documents in result"

type Word struct {
	ID 	   string `json:"_id" bson:"_id,omitempty"`
	Phrase string `json:"phrase" bson:"phrase,omitempty"`
	Index  int    `json:"index" bson:"index,omitempty"`
	Uses   int    `json:"uses" bson:"uses,omitempty"`
}

func GetRandomPhrase(ctx context.Context, db *mongo.Database) (string, error) {
	c1, err := db.Collection(DimaMonth).Find(context.TODO(), bson.D{})
	if err != nil {
		return "", fmt.Errorf("Error fetching DimaMonth: %v", err)
	}
	c2, err := db.Collection(DimaDay).Find(context.TODO(), bson.D{})
	if err != nil {
		return "", fmt.Errorf("Error fetching DimaDay: %v", err)
	}

	var days []Word
	var months []Word
	err = c1.All(context.TODO(), &months)
	if err != nil {
		return "", err
	}
	err = c2.All(context.TODO(), &days)
	if err != nil {
		return "", err
	}

	nMonths := len(months)
	nDays   := len(days)
	if nMonths == 0 || nDays == 0 {
		return "", fmt.Errorf("Could not find words")
	}

	rand.Shuffle(nMonths, func(i, j int) {
		months[i], months[j] = months[j], months[i]
	})
	rand.Shuffle(nDays, func(i, j int) {
		days[i], days[j] = days[j], days[i]
	})

	day   := days[0]
	month := months[0]

	err = UpdateRandomPhraseUses(ctx, db, day, month)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", day.Phrase, month.Phrase), nil
}

func UpdateRandomPhraseUses(ctx context.Context, db *mongo.Database, day Word, month Word) error {
	id, err := primitive.ObjectIDFromHex(day.ID)
	if err != nil {
		return err
	}
	_, err = db.Collection(DimaDay).UpdateByID(
		ctx,
		id,
		bson.D{{
			Key: "$inc", Value: bson.D{{
				Key: "uses", Value: 1,
			}},
		}}, 
	)
	if err != nil {
		return err
	}

	id, err = primitive.ObjectIDFromHex(month.ID)
	if err != nil {
		return err
	}
	_, err = db.Collection(DimaMonth).UpdateByID(
		ctx,
		id,
		bson.D{{
			Key: "$inc", Value: bson.D{{
				Key: "uses", Value: 1,
			}},
		}}, 
	)
	return err
}

func GetWhoAreYouByDima(ctx context.Context, db *mongo.Database, args string) (string, error) {
	args = strings.TrimLeft(args, " ")
	s := strings.Split(args, " ")
	if len(s) != 2 {
		return "Please provide /whoAreYouByDima dayIndex monthIndex", nil
	}

	day, err := strconv.Atoi(s[0])
	if err != nil {
		return fmt.Sprintf("Invalid input: %v", err), nil
	}
	month, err := strconv.Atoi(s[1])
	if err != nil {
		return fmt.Sprintf("Invalid input: %v", err), nil
	}

	var whoDay Word
	var whoMonth Word
	err = db.Collection(DimaDay).FindOne(ctx, bson.D{{ Key: "index", Value: day }}).Decode(&whoDay)
	if err != nil && err.Error() != NoDocuments {
		return "", err
	}
	if err != nil && err.Error() == NoDocuments {
		return fmt.Sprintf("Could not find a Day assosiated with this number: %v", day), nil
	}
	err = db.Collection(DimaMonth).FindOne(ctx, bson.D{{ Key: "index", Value: month }}).Decode(&whoMonth)
	if err != nil && err.Error() != NoDocuments {
		return "", err
	}
	if err != nil && err.Error() == NoDocuments {
		return fmt.Sprintf("Could not find a Month assosiated with this number: %v", month), nil
	}

	err = UpdateRandomPhraseUses(ctx, db, whoDay, whoMonth)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", whoDay.Phrase, whoMonth.Phrase), nil
}