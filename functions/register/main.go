package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"hash/fnv"
	"math"
	"os"
	"strconv"

	"github.com/apex/go-apex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sdming/gosnow"
	"time"
)

// To change these settings for DynamoDB, deploy with a different environment variable.
// apex deploy -s AUTH_DB_REGION=us-west-1
var authDBRegion = os.Getenv("AUTH_DB_REGION")
var authDBTable = os.Getenv("AUTH_DB_TABLE")

// DynamoDB service
var svc = dynamodb.New(session.New(&aws.Config{
	Region: aws.String(authDBRegion),
}))

// The JSON message passd to the Lambda (should include email, password, etc.)
type message struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRecord ...
type UserRecord struct {
	Version   int64     `json:"version"`
	ID        int64     `json:"id"`
	IDStr     string    `json:id_str"`
	Email     string    `json:"email"`
	Updated   time.Time `json:"updated"`
	Created   time.Time `json:"created"`
	LastLogin time.Time `json:"last_login"`
}

func main() {
	// If not set for some reason, use us-east-1 by default.
	if authDBRegion == "" {
		authDBRegion = "us-east-1"
	}

	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message
		//var u UserRecord

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		// Register the user
		u, regErr := register(m, ctx)
		if regErr != nil {
			return nil, regErr
		}
		return u, nil
	})
}

func register(m message, ctx *apex.Context) (UserRecord, error) {
	result := UserRecord{}

	// First make sure the user doesn't already exist
	u, err := getUserByEmail(m.Email)
	if err != nil {
		return result, err
	}
	// If so, just return that.
	if u.Email != "" && u.Email == m.Email {
		return u, nil
	}

	// v, err := gosnow.Default()
	// Because there could be multiple workers running at the same time, create a worker id based on the Lambda request ID.
	v, err := gosnow.NewSnowFlake(generateWorkerID(ctx.RequestID))
	if err != nil {
		return result, err
	}
	id, err := v.Next()
	if err != nil {
		return result, err
	}
	// We'll use a string for the id in DynamoDB
	sID := strconv.FormatUint(id, 10)

	// Hash the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return result, err
	}
	// Note: Comparing the password with the hash
	// err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	// fmt.Println(err) // nil means it is a match

	// User creation and last update time
	created := time.Now()
	createdInt := created.UnixNano()
	createdString := strconv.FormatInt(createdInt, 10)

	params := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(sID),
			},
		},
		TableName: aws.String(authDBTable),
		// KEY and VALUE are reserved words so the query needs to dereference them
		// ExpressionAttributeNames: map[string]*string{
		// 	//"#k": aws.String("key"),
		// 	"#v": aws.String("value"),
		// 	"#u": aws.String("user"),
		// },
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id_str": {
				S: aws.String(sID),
			},
			":email": {
				S: aws.String(m.Email),
			},
			":password": {
				S: aws.String(string(hashedPassword)),
			},

			":created": {
				N: aws.String(createdString),
			},
			":updated": {
				N: aws.String(createdString),
			},
			":last_login": {
				N: aws.String(createdString),
			},
			// user data held in S3 allows for an unlimited amount of user data in any format
			// TODO: create a folder in S3 under whatever bucket based on the user's id.
			// In fact, if that's the case, this data URL may not even be needed... Because the data will be held in a predictable place.
			// ":data": {
			// 	S: aws.String(""),
			// },

			// TODO: Think about also holding a phone number in this user record for SMS recovery (reset) or even 2 factor auth.

			// record version increment just to keep track of how many times the user was updated maybe useful maybe not
			// TODO: maybe one day even look into keeping a history around? ie. "You can't change your password to a password you used previously."
			":i": {
				N: aws.String("1"),
			},
		},
		// ReturnConsumedCapacity:      aws.String("TOTAL"),
		// ReturnItemCollectionMetrics: aws.String("ReturnItemCollectionMetrics"),
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET id_str = :id_str, email = :email, password = :password, updated = :updated, created = :created, last_login = :last_login ADD version :i"),
	}

	response, err := svc.UpdateItem(params)
	if err != nil {
		return result, err
	}

	// The newly set values
	if val, ok := response.Attributes["id"]; ok {
		result.Version, _ = strconv.ParseInt(*response.Attributes["version"].N, 10, 64)
		result.ID, _ = strconv.ParseInt(*val.N, 10, 64)
		result.IDStr = *response.Attributes["id_str"].S
		result.Email = *response.Attributes["email"].S
		uN, _ := strconv.ParseInt(*response.Attributes["updated"].N, 10, 64)
		result.Updated = time.Unix(0, uN)
		cN, _ := strconv.ParseInt(*response.Attributes["created"].N, 10, 64)
		result.Created = time.Unix(0, cN)
		lN, _ := strconv.ParseInt(*response.Attributes["last_login"].N, 10, 64)
		result.LastLogin = time.Unix(0, lN)
	}

	return result, nil
}

// Returns the number of records in the user table
func userCount() (int64, error) {
	c := int64(0)
	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(authDBTable),
	}
	response, err := svc.DescribeTable(params)
	if err == nil {
		c = *response.Table.ItemCount
	}
	return c, err
}

func getUserByEmail(email string) (UserRecord, error) {
	result := UserRecord{}

	// TODO: Fix this query.
	params := &dynamodb.QueryInput{
		TableName: aws.String(authDBTable),
		IndexName: aws.String("EmailIndex"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(email),
			},
		},
		KeyConditionExpression: aws.String("email = :email"),
		Limit: aws.Int64(1),

		// INDEXES | TOTAL | NONE (not required - not even sure if I need to worry about it)
		ReturnConsumedCapacity: aws.String("TOTAL"),
		// Important: This needs to be false so it returns results in descending order. If it's true (the default), it's sorted in the
		// order values were stored. So the first item stored for the key ever would be returned...But the latest item is needed.
		ScanIndexForward: aws.Bool(false),
		// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Query.html#DDB-Query-request-Select
		// Select: aws.String("SPECIFIC_ATTRIBUTES"),
		// ProjectionExpression: aws.String("id,id_str,version,email,created,updated,last_login"),
		// TODO: Not sure we need all attributes projected in the secondary index. it's going to use more space...
		// not that these values really take up that much storage space...but it'd be nice.
		// this could lead to two queries for certain things to avoid slower scans.
		// query by email, get id, and hen GetItem by id (primary key). two operations, but storage efficient and faster
		// than a scan.
	}
	response, err := svc.Query(params)
	if err != nil {
		return result, err
	}
	if len(response.Items) > 0 {
		if val, ok := response.Items[0]["id"]; ok {
			result.ID, _ = strconv.ParseInt(*val.N, 10, 64)
			result.IDStr = *response.Items[0]["id_str"].S
			result.Version, _ = strconv.ParseInt(*response.Items[0]["version"].N, 10, 64)
			result.Email = *response.Items[0]["email"].S
			uN, _ := strconv.ParseInt(*response.Items[0]["updated"].N, 10, 64)
			result.Updated = time.Unix(0, uN)
			cN, _ := strconv.ParseInt(*response.Items[0]["created"].N, 10, 64)
			result.Created = time.Unix(0, cN)
			lN, _ := strconv.ParseInt(*response.Items[0]["last_login"].N, 10, 64)
			result.LastLogin = time.Unix(0, lN)
		}
	}
	return result, nil
}

// Generates a worker id based on a string, like the ctx.RequestID
func generateWorkerID(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	n := h.Sum32()
	f := float64(n)
	// likely something like 3042975444 given the string length, but the snowflake max worker id is 1032
	// keep squaring it til it's low enough
	for f > 1032 {
		f = math.Sqrt(f)
	}
	return uint32(f)
}
