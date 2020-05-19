package collections

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Resume :
type Resume struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	ScheduleID        string             `json:"schedule_id" bson:"scheduleId"`
	ScheduleHistoryID string             `json:"schedule_history_id" bson:"scheduleHistoryId"`
	WorkID            string             `json:"work_id" bson:"workId"`
	CompanyID         int64              `json:"company_id" bson:"companyId"`
	At                int64              `json:"at" bson:"at"`
	URL               string             `json:"url" bson:"url"`
	URLHash           string             `json:"url_hash" bson:"urlHash"`
	RequestCount      int32              `json:"request_count" bson:"requestCount"`
	LoadTime          float64            `json:"load_time" bson:"loadTime"`
	Cookies           []Cookie           `json:"cookies" bson:"cookies"`
	Domains           []string           `json:"domains" bson:"domains"`
}

// Cookie :
type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	SameSite *int64  `json:"samesite"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	MaxAge   string  `json:"maxage"`
	Expires  float64 `json:"expires"`
	Secure   bool    `json:"secure"`
	HTTPOnly bool    `json:"httponly"`
	Size     *int64  `json:"size"`
	Priority *string `json:"priority"`
}

// ResumeCollection :
type ResumeCollection struct {
}

// FindAll :
func (c ResumeCollection) FindAll(filters map[string]interface{}) (documents []*Resume, err error) {
	collection := getConn().Database(dbName).Collection("resumes")

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "at", Value: 1}})
	findOptions.SetLimit(int64(10))

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filters, findOptions)
	if err != nil {
		return nil, fmt.Errorf("ResumeCollection FindAll Find error: %s", err.Error())
	}

	for cur.Next(context.TODO()) {
		var elem Resume

		err := cur.Decode(&elem)
		if err != nil {
			return nil, fmt.Errorf("ResumeCollection FindAll Decode error: %s", err.Error())
		}

		documents = append(documents, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("ResumeCollection FindAll cur.Err: %s", err.Error())
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return documents, nil
}

// FindOne :
func (c ResumeCollection) FindOne(filter map[string]interface{}) (Resume, error) {
	collection := getConn().Database(dbName).Collection("resumes")

	var elem Resume

	// get only the last
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{primitive.E{Key: "at", Value: 1}})

	err := collection.FindOne(context.TODO(), filter, findOptions).Decode(&elem)
	if err != nil {
		return elem, fmt.Errorf("ResumeCollection FindOne error: %s", err.Error())
	}

	return elem, nil
}
