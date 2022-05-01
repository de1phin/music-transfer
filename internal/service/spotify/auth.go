package spotify

import (
	"fmt"

	"github.com/de1phin/music-transfer/internal/api/spotify"
	"github.com/de1phin/music-transfer/internal/mux"
)

func (ss *spotifyService) GetAuthURL(userID int64) (string, error) {
	return fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&"+
			"response_type=code&redirect_uri=%s&scope=%s&state=%d",
		ss.client.ID,
		ss.redirectURI,
		ss.scopes,
		userID,
	), nil
}

func (spotify *spotifyService) BindOnAuthorized(OnAuthorized mux.OnAuthorized) {
	spotify.OnAuthorizedNotify = OnAuthorized
}

func (ss *spotifyService) OnGetTokens(userID int64, tokens spotify.Credentials) error {
	err := ss.tokenStorage.Set(userID, tokens)
	if err != nil {
		return err
	}
	if ss.OnAuthorizedNotify != nil {
		ss.OnAuthorizedNotify(ss, userID)
	}
	return nil
}

func (ss *spotifyService) Authorized(userID int64) (bool, error) {
	exist, err := ss.tokenStorage.Exist(userID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	tokens, err := ss.tokenStorage.Get(userID)
	if err != nil {
		return false, err
	}
	return ss.api.Authorized(tokens)
}
