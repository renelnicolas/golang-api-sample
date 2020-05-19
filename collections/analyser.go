package collections

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Analyser :
type Analyser struct {
	ID                  primitive.ObjectID     `json:"id" bson:"_id"`
	WorkID              string                 `json:"work_id" bson:"workId"`
	CompanyID           int64                  `json:"company_id" bson:"companyId"`
	At                  int64                  `json:"at" bson:"at"`
	URL                 string                 `json:"url" bson:"url"`
	URLHash             string                 `json:"url_hash" bson:"urlHash"`
	Cookies             []Cookie               `json:"cookies" bson:"cookies"`
	Domains             []string               `json:"domains" bson:"domains"`
	PostData            interface{}            `json:"post_data" bson:"postData"`
	Headers             map[string]interface{} `json:"headers" bson:"headers"`
	RemoteAddress       map[string]interface{} `json:"remote_address" bson:"remoteAddress"`
	FromCache           bool                   `json:"from_cache" bson:"fromCache"`
	Status              int                    `json:"status" bson:"status"`
	StatusText          string                 `json:"status_text" bson:"statusText"`
	Ok                  bool                   `json:"ok" bson:"ok"`
	IsDetached          bool                   `json:"is_detached" bson:"isDetached"`
	FrameID             string                 `json:"frame_id" bson:"frameId"`
	IsNavigationRequest bool                   `json:"is_navigationRequest" bson:"isNavigationRequest"`
	Method              string                 `json:"method" bson:"method"`
	ResourceType        string                 `json:"resource_type" bson:"resourceType"`
	RequestID           string                 `json:"request_id" bson:"requestId"`
	Timing              map[string]interface{} `json:"timing" bson:"timing"`
	Extras              map[string]interface{} `json:"-" bson:"extras"`
	// Content             string                 `json:"content" bson:"content"`
}

// AnalyserCollection :
type AnalyserCollection struct {
}

// FindAll :
func (c AnalyserCollection) FindAll(filters map[string]interface{}) (documents []*Analyser, err error) {
	collection := getConn().Database(dbName).Collection("analysers")

	findOptions := options.Find()

	cur, err := collection.Find(context.TODO(), filters, findOptions)
	if err != nil {
		return nil, fmt.Errorf("AnalyserCollection FindAll Find error: %s", err.Error())
	}

	for cur.Next(context.TODO()) {
		var elem Analyser

		err := cur.Decode(&elem)
		if err != nil {
			return nil, fmt.Errorf("AnalyserCollection FindAll Decode error: %s", err.Error())
		}

		documents = append(documents, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("AnalyserCollection FindAll cur.Err: %s", err.Error())
	}

	cur.Close(context.TODO())

	return documents, nil
}
