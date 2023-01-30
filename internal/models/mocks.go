package models

type MockMovieModel struct{
	Movies map[int64]Movie
}

func (m *MockMovieModel) Insert(movie *Movie) error {
	m.Movies[movie.ID] = *movie
	return nil
}

func (m *MockMovieModel) Get(id int64) (*Movie, error) {
	if int(id) < 0 {
		return nil, ErrRecordNotFound
	}

	movie, ok := m.Movies[id]
	if !ok {
		return nil, ErrRecordNotFound
	}

	return &movie, nil
}

func (m *MockMovieModel) Update(movie *Movie) error {
	if int(movie.ID) < 0 || int(movie.ID) > len(m.Movies) {
		return ErrRecordNotFound
	}

	movie.Version++
	m.Movies[movie.ID] = *movie
	return nil
}

func (m *MockMovieModel) Delete(id int64) error {
	if int(id) < 0 || int(id) > len(m.Movies) {
		return ErrRecordNotFound
	}

	delete(m.Movies, id)

	return nil
}
