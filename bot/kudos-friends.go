package bot

import (
	"log"
	"strconv"

	"strings"

	"github.com/stevegore/stravaKudos/parser"
)

func (s *Strava) kudosFriend(c *parser.Client, friendId int64) {
	var headers = map[string]string{}
	headers["authorization"] = "access_token " + s.authToken

	var myFolowersUrl = strings.ReplaceAll(s.MapUrls["kudos_url"], "{ACTIVITIES-ID}", strconv.Itoa(int(friendId)))

	_, statusCode := c.MakeRequest(myFolowersUrl, "POST", "", headers)

	if statusCode != 201 {
		log.Fatalf("Status from kudosFriend request no HTTP_CREATED | statusCode => %d", statusCode)
	}
}
