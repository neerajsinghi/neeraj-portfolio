package blog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound  = errors.New("blog post not found")
	ErrConflict  = errors.New("blog post was changed by another request")
	ErrForbidden = errors.New("operation is not permitted for this role")
)

type Store interface {
	ListPublished(context.Context, int) ([]Post, error)
	GetPublishedBySlug(context.Context, string) (Post, error)
	ListAll(context.Context, int) ([]Post, error)
	Create(context.Context, WriteInput, Principal) (Post, error)
	Update(context.Context, string, WriteInput, Principal) (Post, error)
	Delete(context.Context, string, Principal) error
}

type MongoStore struct {
	client    *mongo.Client
	posts     *mongo.Collection
	revisions *mongo.Collection
	audit     *mongo.Collection
}

type postDocument struct {
	MongoID primitive.ObjectID `bson:"_id,omitempty"`
	Post    `bson:",inline"`
}

type revisionDocument struct {
	PostID    primitive.ObjectID `bson:"post_id"`
	Version   int                `bson:"version"`
	Snapshot  Post               `bson:"snapshot"`
	ChangedBy Author             `bson:"changed_by"`
	CreatedAt time.Time          `bson:"created_at"`
}

type auditDocument struct {
	PostID    primitive.ObjectID `bson:"post_id"`
	Action    string             `bson:"action"`
	Actor     Author             `bson:"actor"`
	CreatedAt time.Time          `bson:"created_at"`
}

func NewMongoStore(ctx context.Context, uri, database string) (*MongoStore, error) {
	if uri == "" {
		return nil, errors.New("MONGODB_URI is required")
	}
	if database == "" {
		database = "neeraj_portfolio"
	}
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(5).
		SetMinPoolSize(0).
		SetMaxConnIdleTime(30 * time.Second).
		SetConnectTimeout(5 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetSocketTimeout(15 * time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connect to MongoDB: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}

	db := client.Database(database)
	store := &MongoStore{
		client: client, posts: db.Collection("blog_posts"),
		revisions: db.Collection("blog_revisions"), audit: db.Collection("blog_audit"),
	}
	if err := store.ensureIndexes(ctx); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}
	return store, nil
}

func (store *MongoStore) ensureIndexes(ctx context.Context) error {
	_, err := store.posts.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "slug", Value: 1}}, Options: options.Index().SetUnique(true).SetName("unique_slug")},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "published_at", Value: -1}}, Options: options.Index().SetName("published_feed")},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "scheduled_at", Value: 1}}, Options: options.Index().SetName("scheduled_release")},
	})
	if err != nil {
		return fmt.Errorf("create blog indexes: %w", err)
	}
	_, err = store.revisions.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "post_id", Value: 1}, {Key: "version", Value: -1}},
		Options: options.Index().SetUnique(true).SetName("post_version"),
	})
	if err != nil {
		return fmt.Errorf("create revision index: %w", err)
	}
	return nil
}

func (store *MongoStore) ListPublished(ctx context.Context, limit int) ([]Post, error) {
	return store.list(ctx, publicFilter(time.Now().UTC()), limit, bson.D{{Key: "published_at", Value: -1}})
}

func (store *MongoStore) ListAll(ctx context.Context, limit int) ([]Post, error) {
	return store.list(ctx, bson.M{}, limit, bson.D{{Key: "updated_at", Value: -1}})
}

func (store *MongoStore) list(ctx context.Context, filter any, limit int, sort any) ([]Post, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	cursor, err := store.posts.Find(ctx, filter, options.Find().SetLimit(int64(limit)).SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	posts := make([]Post, 0)
	for cursor.Next(ctx) {
		var document postDocument
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}
		document.Post.ID = document.MongoID.Hex()
		posts = append(posts, document.Post)
	}
	return posts, cursor.Err()
}

func (store *MongoStore) GetPublishedBySlug(ctx context.Context, slug string) (Post, error) {
	var document postDocument
	filter := publicFilter(time.Now().UTC())
	filter["slug"] = slug
	err := store.posts.FindOne(ctx, filter).Decode(&document)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Post{}, ErrNotFound
	}
	if err != nil {
		return Post{}, err
	}
	document.Post.ID = document.MongoID.Hex()
	return document.Post, nil
}

func (store *MongoStore) Create(ctx context.Context, input WriteInput, principal Principal) (Post, error) {
	if !principal.CanEdit() || (input.Status != "" && input.Status != StatusDraft && !principal.CanPublish()) {
		return Post{}, ErrForbidden
	}
	if err := input.NormalizeAndValidate(); err != nil {
		return Post{}, err
	}
	now := time.Now().UTC()
	post := Post{
		Slug: input.Slug, Title: input.Title, Description: input.Description,
		ContentMarkdown: input.ContentMarkdown, Tags: input.Tags, Status: input.Status,
		Author:    Author{Subject: principal.Subject, Email: principal.Email},
		CreatedAt: now, UpdatedAt: now, SchemaVersion: 1, Version: 1,
	}
	if post.Status == StatusPublished {
		post.PublishedAt = &now
	} else if post.Status == StatusScheduled {
		post.ScheduledAt = input.ScheduledAt
		post.PublishedAt = input.ScheduledAt
	}
	result, err := store.posts.InsertOne(ctx, post)
	if mongo.IsDuplicateKeyError(err) {
		return Post{}, ErrConflict
	}
	if err != nil {
		return Post{}, err
	}
	post.ID = result.InsertedID.(primitive.ObjectID).Hex()
	_ = store.writeAudit(ctx, result.InsertedID.(primitive.ObjectID), "create", principal)
	return post, nil
}

func (store *MongoStore) Update(ctx context.Context, id string, input WriteInput, principal Principal) (Post, error) {
	if !principal.CanEdit() {
		return Post{}, ErrForbidden
	}
	if err := input.NormalizeAndValidate(); err != nil {
		return Post{}, err
	}
	if input.Version < 1 {
		return Post{}, errors.Join(ErrInvalid, errors.New("version is required when updating a post"))
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Post{}, ErrNotFound
	}
	var existing postDocument
	if err := store.posts.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existing); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Post{}, ErrNotFound
		}
		return Post{}, err
	}
	if existing.Version != input.Version {
		return Post{}, ErrConflict
	}
	if !principal.CanPublish() && (existing.Status != StatusDraft || input.Status != StatusDraft) {
		return Post{}, ErrForbidden
	}

	now := time.Now().UTC()
	updated := existing.Post
	updated.Slug, updated.Title, updated.Description = input.Slug, input.Title, input.Description
	updated.ContentMarkdown, updated.Tags, updated.Status = input.ContentMarkdown, input.Tags, input.Status
	updated.UpdatedAt, updated.Version = now, existing.Version+1
	updated.ScheduledAt = nil
	if input.Status == StatusPublished {
		if existing.Status != StatusPublished || updated.PublishedAt == nil {
			updated.PublishedAt = &now
		}
	} else if input.Status == StatusScheduled {
		updated.ScheduledAt = input.ScheduledAt
		updated.PublishedAt = input.ScheduledAt
	} else {
		updated.PublishedAt = nil
	}
	result, err := store.posts.ReplaceOne(ctx, bson.M{"_id": objectID, "version": input.Version}, updated)
	if mongo.IsDuplicateKeyError(err) {
		return Post{}, ErrConflict
	}
	if err != nil {
		return Post{}, err
	}
	if result.MatchedCount == 0 {
		return Post{}, ErrConflict
	}
	_, _ = store.revisions.InsertOne(ctx, revisionDocument{
		PostID: objectID, Version: existing.Version, Snapshot: existing.Post,
		ChangedBy: Author{Subject: principal.Subject, Email: principal.Email}, CreatedAt: now,
	})
	_ = store.writeAudit(ctx, objectID, "update", principal)
	updated.ID = id
	return updated, nil
}

func publicFilter(now time.Time) bson.M {
	return bson.M{"$or": bson.A{
		bson.M{"status": StatusPublished},
		bson.M{"status": StatusScheduled, "scheduled_at": bson.M{"$lte": now}},
	}}
}

func (store *MongoStore) Delete(ctx context.Context, id string, principal Principal) error {
	if !principal.CanDelete() {
		return ErrForbidden
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}
	result, err := store.posts.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	_ = store.writeAudit(ctx, objectID, "delete", principal)
	return nil
}

func (store *MongoStore) writeAudit(ctx context.Context, postID primitive.ObjectID, action string, principal Principal) error {
	_, err := store.audit.InsertOne(ctx, auditDocument{
		PostID: postID, Action: action,
		Actor: Author{Subject: principal.Subject, Email: principal.Email}, CreatedAt: time.Now().UTC(),
	})
	return err
}
