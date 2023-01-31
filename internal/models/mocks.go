package models

type MockMovieModel struct {
	MoviesDB map[int64]Movie
}

func DefaultMockMovieModel() *MockMovieModel {
	mm := MockMovieModel{MoviesDB: map[int64]Movie{}}
	mm.MoviesDB[1] = Movie{
		ID:      1,
		Title:   "mock movie",
		Year:    2000,
		Runtime: 100,
		Genres:  []string{"mock genre1", "mock genre2"},
		Version: 1,
	}
	return &mm
}

func (m *MockMovieModel) Insert(movie *Movie) error {
	m.MoviesDB[movie.ID] = *movie
	return nil
}

func (m *MockMovieModel) Get(id int64) (*Movie, error) {
	if int(id) < 0 {
		return nil, ErrRecordNotFound
	}

	movie, ok := m.MoviesDB[id]
	if !ok {
		return nil, ErrRecordNotFound
	}

	return &movie, nil
}

func (m *MockMovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	return nil, Metadata{}, nil
}

func (m *MockMovieModel) Update(movie *Movie) error {
	if int(movie.ID) < 0 {
		return ErrRecordNotFound
	}

	movie.Version++
	m.MoviesDB[movie.ID] = *movie
	return nil
}

func (m *MockMovieModel) Delete(id int64) error {
	if int(id) < 0 {
		return ErrRecordNotFound
	}

	var found bool
	for k := range m.MoviesDB {
		if k == id {
			found = true
		}
	}
	if !found {
		return ErrRecordNotFound
	}

	delete(m.MoviesDB, id)

	return nil
}
