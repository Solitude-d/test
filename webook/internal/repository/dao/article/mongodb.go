package article

import (
	"context"
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDAO struct {
	//client *mongo.Client
	//webook 的 database
	//database *mongo.Database
	//制作库
	col *mongo.Collection
	//线上库
	colLive *mongo.Collection
	node    *snowflake.Node
}

func (m *MongoDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Utime = now
	art.Ctime = now
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	return id, err
}

func (m *MongoDAO) UpdateById(ctx context.Context, art Article) error {
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.D{bson.E{"$set", bson.M{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   time.Now().UnixMilli(),
	}}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func (m *MongoDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDAO) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	art.Id = id
	now := time.Now().UnixMilli()
	art.Utime = now
	filter := bson.M{"id": art.Id}
	//update := bson.E{"$set", art}
	//upsert := bson.E{"$setOnInsert", bson.D{bson.E{"ctime", now}}}
	update := bson.M{
		"$set":         art,
		"$setOnInsert": bson.M{"ctime": now},
	}
	_, err = m.colLive.UpdateOne(ctx, filter,
		//bson.D{update, upsert},
		update,
		options.Update().SetUpsert(true))
	return id, err
}

func (m *MongoDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func NewMongoDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDAO{
		col:     db.Collection("articles"),
		colLive: db.Collection("publish_articles"),
		node:    node,
	}
}

func InitCollections(db mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	//创建索引
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1}},
		},
	}
	_, err := db.Collection("articles").Indexes().CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_article").Indexes().CreateMany(ctx, index)
	return err
}
