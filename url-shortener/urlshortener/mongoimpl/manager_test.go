package mongoimpl

import (
	"context"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"url-shortener/urlshortener"
)

var ctx = context.Background()

func TestManager(t *testing.T) {
	suite.Run(t, &ManagerSuite{mongoAddr: "mongodb://localhost:27017"})
}

type ManagerSuite struct {
	suite.Suite

	mongoAddr   string
	mongoClient *mongo.Client

	manager urlshortener.Manager
}

func (s *ManagerSuite) SetupSuite() {
	s.manager = NewManager(s.mongoAddr)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(s.mongoAddr))
	s.Require().NoError(err)
	s.mongoClient = mongoClient
}

func (s *ManagerSuite) SetupTest() {
	s.Require().NoError(s.mongoClient.Database(dbName).Drop(ctx))
}

func (s *ManagerSuite) TearDownTest() {
	s.Require().NoError(s.mongoClient.Database(dbName).Drop(ctx))
}

func (s *ManagerSuite) TestNotFound() {
	// when:
	_, err := s.manager.ResolveShortcut(ctx, "there-is-no-such-key")

	// then:
	s.Require().ErrorIs(err, urlshortener.ErrNotFound)
}

func (s *ManagerSuite) TestCreateResolve() {
	// given:
	const fullURL = "https://google.de"

	// when:
	key, err := s.manager.CreateShortcut(ctx, fullURL)

	// then:
	s.Require().NoError(err)
	s.Require().NotEmpty(key)

	// when:
	resolved, err := s.manager.ResolveShortcut(ctx, key)

	// then:
	s.Require().NoError(err)
	s.Require().Equal(fullURL, resolved)
}
