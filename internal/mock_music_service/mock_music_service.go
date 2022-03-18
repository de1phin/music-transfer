package mockmusicservice

import "fmt"

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

func (service *mockMusicService) GetAuthURL(id int64) string {
	return fmt.Sprintf("%s?id=%d", service.callbackURL+"/"+service.URLName(), id)
}
