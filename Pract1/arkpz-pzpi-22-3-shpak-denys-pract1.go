// Форматування
// Поганий приклад
func PostHandler( w  http.ResponseWriter ,  r  * http.Request )  { 
	defer r .Body.Close()
	  postId ,_:=primitive.ObjectIDFromHex(id) 

 render.JSON(w,r, map[string]interface{}{ 
	"Status" : http.StatusOK,
	  "Message" :"Successfully deleted" } ) } 
// Гарний приклад
func PostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := chi.URLParam(r, "id")

	postId, _ := primitive.ObjectIDFromHex(id)

	render.JSON(w, r, map[string]interface{}{
		"Status":  http.StatusOK,
		"Message": "Successfully deleted",
	})
}

// Коментарі
// Поганий приклад
// Фунція приймає w та r
func phandl(w http.ResponseWriter, r *http.Request) {
	// Відкладення закриття тіла запиту
	defer r.Body.Close()

	// Отримання айді з юрл
	id := chi.URLParam(r, "id")

	// Перетворення айді на много об'єкт
	pId, _ := primitive.ObjectIDFromHex(id)

	// Виклик функції Delete
	_ = postRemover.Delete(
		context.TODO(),
		postId,
	)
	// Запис відповіді в форматі JSON
	render.JSON(w, r, map[string]interface{}{
		"Status":  http.StatusOK,
		"Message": "done",
	})
}
// Гарний приклад
// PostHandler - це обробник події видалення поста.
// Він викликає функцію Delete інтерфейсу PostRemover
func PostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := chi.URLParam(r, "id")

	postId, _ := primitive.ObjectIDFromHex(id)
	_ = postRemover.Delete(
		context.TODO(),
		postId,
	)

	render.JSON(w, r, map[string]interface{}{
		"Status":  http.StatusOK,
		"Message": "Successfully deleted",
	})
}

// Іменування
// Поганий приклад
package mongodbwithsavingseature

type storage struct {
	Db *mongo.Client
}

func new_storage(storagePath string) *storage {
	db, _ := mongo.Connect(
		options.Client().ApplyURI(storagePath),
	)

	return storage{Db: db}
}
// Гарний приклад
package mongodb

type Storage struct {
	db *mongo.Client
}

func New(storagePath string) *Storage {
	db, _ := mongo.Connect(
		options.Client().ApplyURI(storagePath),
	)

	return Storage{db: db}
}

// Функції
// Поганий приклад
func New(log *slog.Logger, postRemover PostRemover) http.HandlerFunc {
	func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		postId, _ := primitive.ObjectIDFromHex(id)

		render.JSON(w, r, map[string]interface{}{
			"Status":  http.StatusOK,
			"Message": "Successfully deleted",
		})
	}
}
// Гарний приклад
func New(log *slog.Logger, postRemover PostRemover) (http.HandlerFunc, error) {
	func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		id := chi.URLParam(r, "id")

		postId, _ := primitive.ObjectIDFromHex(id)

		render.JSON(w, r, map[string]interface{}{
			"Status":  http.StatusOK,
			"Message": "Successfully deleted",
		})
	}, nil
}

// Помилки
// Поганий приклад
func New(storagePath string) (*Storage, error) {
	db, _ := mongo.Connect(
		options.Client().ApplyURI(storagePath),
	)

	return &Storage{db: db}, nil
}
// Гарний приклад
func New(storagePath string) (*Storage, error) {
	const op = "storage.mongodb.New"

	db, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(storagePath),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Помилки: паніка
// Поганий приклад
func New(storagePath string,
	httpPort int,
	timeout time.Duration) *App {
	storage, err := mongodb.New(storagePath)
	if err != nil {
		fmt.Println("DB not started")
	}

	httpApp := httpapp.New(
		storage, storage, storage,
		httpPort, timeout,
	)

	return &App{
		HttpServer: httpApp,
	}
}
// Гарний приклад
func New(storagePath string,
	httpPort int,
	timeout time.Duration) *App {
	storage, err := mongodb.New(storagePath)
	if err != nil {
		panic(err)
	}

	httpApp := httpapp.New(
		storage, storage, storage,
		httpPort, timeout,
	)

	return &App{
		HttpServer: httpApp,
	}
}

// Контекст
// Поганий приклад
func (s *Storage) GetPostById(
	postId primitive.ObjectID,
) (*models.Post, error) {
	collection := s.db.Database("DevHubDB").Collection("posts")

	post := &models.Post{}
	filter := bson.M{"_id": postId}

	_ = collection.FindOne(context.TODO(), filter).Decode(post)

	return post, nil
}
// Гарний приклад
func (s *Storage) GetPostById(
	ctx context.Context,
	postId primitive.ObjectID,
) (*models.Post, error) {
	collection := s.db.Database("DevHubDB").Collection("posts")

	post := &models.Post{}
	filter := bson.M{"_id": postId}

	_ = collection.FindOne(ctx, filter).Decode(post)

	return post, nil
}