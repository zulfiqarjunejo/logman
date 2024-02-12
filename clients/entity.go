package clients

type Client struct {
	ApiKey      string `bson:"api_key"`
	DisplayName string `bson:"display_name"`
	Id          string `bson:"client_id"`
}
