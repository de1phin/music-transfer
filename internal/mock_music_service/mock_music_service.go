package mockmusicservice

type mockMusicService struct {
	callbackURL string
}

func NewMockMusicService(callbackURL string) *mockMusicService {
	return &mockMusicService{callbackURL}
}

func (service *mockMusicService) Name() string {
	return "Mock"
}

func (service *mockMusicService) URLName() string {
	return "mock"
}
