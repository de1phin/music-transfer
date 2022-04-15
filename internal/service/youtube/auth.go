package youtube

import (
	"fmt"

	"github.com/de1phin/music-transfer/internal/api/youtube"
)

func (yt *youtubeService) GetAuthURL(userID int64) (string, error) {
	return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&response_type=code&scope=%s&redirect_uri=%s&state=%d", yt.config.ClientID, yt.config.Scopes, yt.config.RedirectURI, userID), nil
}

func (yt *youtubeService) OnGetTokens(userID int64, tokens youtube.Credentials) {
	yt.tokenStorage.Put(userID, tokens)
}

func (yt *youtubeService) Authorized(userID int64) (bool, error) {
	tokens := yt.tokenStorage.Get(userID)
	return yt.api.Authorized(tokens)
}
