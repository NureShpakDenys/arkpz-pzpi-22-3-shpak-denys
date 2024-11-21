// Bad example
id, err := storage.Save(
	context.Background(),
	primitive.NewObjectID(),
	"My Post Title",
    "Post content goes here...",
	"header_image_url",
	[]string{"tag1", "tag2"}
)

func (s *Storage) Save(
	ctx context.Context,
	userId primitive.ObjectID,
	title string,
	content string,
	headerImage string,
	tags []string) 
(primitive.ObjectID, error) {
	const op = "storage.mongodb.Save"

	collection := s.db.Database("DevHubDB").Collection("posts")
	post := &models.Post{
		User:        userId,
		Title:       title,
		Content:     content,
		CreatedAt:   time.Now(),
		Likes:       0,
		Dislikes:    0,
		HeaderImage: headerImage,
		Comments:    []primitive.ObjectID{},
		Tags:        tags,
	}

	insertResult, err := collection.InsertOne(ctx, post)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%s: %w", op, err)
	}

	oid := insertResult.InsertedID.(primitive.ObjectID)

	return oid, nil
}

// Після рефакторингу
data := PostData{
    UserID:      primitive.NewObjectID(),
    Title:       "My Post Title",
    Content:     "Post content goes here...",
    HeaderImage: "header_image_url",
    Tags:        []string{"tag1", "tag2"},
}

id, err := storage.Save(context.Background(), data)

func (s *Storage) Save(ctx context.Context, data PostData) (primitive.ObjectID, error) {
	const op = "storage.mongodb.Save"

	collection := s.db.Database("DevHubDB").Collection("posts")

	post := &models.Post{
		User:        data.UserID,
		Title:       data.Title,
		Content:     data.Content,
		CreatedAt:   time.Now(),
		Likes:       0,
		Dislikes:    0,
		HeaderImage: data.HeaderImage,
		Comments:    []primitive.ObjectID{},
		Tags:        data.Tags,
	}

	insertResult, err := collection.InsertOne(ctx, post)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%s: %w", op, err)
	}

	oid := insertResult.InsertedID.(primitive.ObjectID)

	return oid, nil
}
