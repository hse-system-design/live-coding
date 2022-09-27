package mongoimpl

type urlItem struct {
	Key string `bson:"_id"`
	URL string `bson:"url"`
}
