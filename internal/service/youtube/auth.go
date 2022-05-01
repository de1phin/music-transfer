package youtube

import (
	"fmt"

	"github.com/de1phin/music-transfer/internal/api/youtube"
	"github.com/de1phin/music-transfer/internal/mux"
)

func (yt *youtubeService) GetAuthURL(userID int64) (string, error) {
	return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&response_type=code&scope=%s&redirect_uri=%s&state=%d", yt.api.ClientID, yt.api.Scopes, yt.api.RedirectURI, userID), nil
}

func (yt *youtubeService) BindOnAuthorized(OnAuthorized mux.OnAuthorized) {
	yt.OnAuthorizedNotify = OnAuthorized
}

func (yt *youtubeService) OnGetTokens(userID int64, tokens youtube.Credentials) error {
	err := yt.tokenStorage.Set(userID, tokens)
	if err != nil {
		return err
	}

	if yt.OnAuthorizedNotify != nil {
		yt.OnAuthorizedNotify(yt, userID)
	}
	return nil
}

func (yt *youtubeService) Authorized(userID int64) (bool, error) {
	exist, err := yt.tokenStorage.Exist(userID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, err
	}
	tokens, err := yt.tokenStorage.Get(userID)
	if err != nil {
		return false, err
	}
	return yt.api.Authorized(tokens)
}
