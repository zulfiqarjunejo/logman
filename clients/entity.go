package clients

type Client struct {
	ApiKey      string `bson:"api_key" json:"-"`
	DisplayName string `bson:"display_name" json:"name"`
	Id          string `bson:"client_id" json:"id"`
}
