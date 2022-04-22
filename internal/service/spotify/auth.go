package spotify

import (
	"fmt"
	"log"

	"github.com/de1phin/music-transfer/internal/api/spotify"
)

func (ss *spotifyService) GetAuthURL(userID int64) (string, error) {
	return fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&"+
			"response_type=code&redirect_uri=%s/spotify&scope=%s&state=%d",
		ss.client.ID,
		ss.redirectURI,
		ss.scopes,
		userID,
	), nil
}

func (ss *spotifyService) OnGetTokens(userID int64, tokens spotify.Credentials) {
	log.Println("PUT", userID, tokens)
	ss.tokenStorage.Put(userID, tokens)
}

func (ss *spotifyService) Authorized(userID int64) (bool, error) {
	tokens, err := ss.tokenStorage.Get(userID)
	if err != nil {
		return false, err
	}
	return ss.api.Authorized(tokens)
}
