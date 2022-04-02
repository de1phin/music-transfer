package spotify

import (
	"fmt"

	"github.com/de1phin/music-transfer/internal/api/spotify"
)

func (ss *spotifyService) GetAuthURL(userID int64) string {
	return fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&"+
			"response_type=code&redirect_uri=%s/spotify&scope=%s&state=%d",
		ss.client.ID,
		ss.redirectURI,
		ss.scopes,
		userID,
	)
}

func (ss *spotifyService) OnGetTokens(userID int64, credentials spotify.Credentials) {
	ss.storage.Put(userID, credentials)
}
