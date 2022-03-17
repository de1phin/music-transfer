package mockmusicservice

import "fmt"

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

func (service *mockMusicService) URLName() string {
	return "mock"
}

func (service *mockMusicService) GetAuthURL(id int64, url string) string {
	return fmt.Sprintf("%s?id=%d", url+"/"+service.URLName(), id)
}
