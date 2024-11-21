// До рефакторингу
func (s *service) AddEmotionalAnalysis(emotionalAnalysis models.EmotionalAnalysis) (int, error) {
	ctx, cancel := s.createContextWithTimeout(5 * time.Second)
	defer cancel()

	result, err := s.executeInsertQuery(ctx, emotionalAnalysis)
	if err != nil {
		return s.handleInsertError(err)
	}

	id, err := s.getLastInsertID(result)
	if err != nil {
		return s.handleIDError()
	}

	return s.convertIDToInt(id)
}

func (s *service) createContextWithTimeout(duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), duration)
}

func (s *service) executeInsertQuery(ctx context.Context, emotionalAnalysis models.EmotionalAnalysis) (sql.Result, error) {
	query := `INSERT INTO emotionalanalysis (emotionalstate, emotionalIcon) VALUES (?, ?)`
	return s.db.ExecContext(ctx, query, emotionalAnalysis.EmotionalState, emotionalAnalysis.EmotionalIcon)
}

func (s *service) handleInsertError(err error) (int, error) {
	log.Printf("Error executing insert query: %v", err)
	return 0, fmt.Errorf("failed to execute insert query: %w", err)
}

func (s *service) getLastInsertID(result sql.Result) (int64, error) {
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}
	return id, nil
}

func (s *service) handleIDError() (int, error) {
	return 0, fmt.Errorf("error retrieving ID")
}

func (s *service) convertIDToInt(id int64) (int, error) {
	return int(id), nil
}

// Після рефакторингу
func (s *service) AddEmotionalAnalysis(emotionalAnalysis models.EmotionalAnalysis) (int, error) {
	query := `INSERT INTO emotionalanalysis (emotionalstate, emotionalIcon) VALUES (?, ?)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	result, err := s.db.ExecContext(
		ctx,
		query, 
		emotionalAnalysis.EmotionalState, 
		emotionalAnalysis.EmotionalIcon,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}