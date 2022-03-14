package mockmusicservice

type mockMusicService struct {
}

type credentials struct {
}

func NewMockMusicService() *mockMusicService {
	return new(mockMusicService)
}

func (service *mockMusicService) Name() string {
	return "Mock"
}
