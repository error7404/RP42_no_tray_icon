package api

import (
	"encoding/json"
	"fmt"
	"github.com/alexandregv/RP42/pkg/oauth"
	"io/ioutil"
)

const URL = "https://api.intra.42.fr"

// fetch() queries an endpoint of the API.
func fetch(endpoint string) []byte {
	client := oauth.GetClient()

	resp, err := client.Get(fmt.Sprint(URL, endpoint))
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return body
	} else {
		panic(fmt.Sprintf("The API responded with a bad status code (%d): %s", resp.StatusCode, string(body)))
	}
}

// GetUser() returns an User, based on his login.
func GetUser(login string) *User {
	resp := fetch(fmt.Sprint("/v2/users/", login))

	user := User{}
	json.Unmarshal(resp, &user)

	return &user
}

// GetUserLastLocation returns the last Location of an user.
func GetUserLastLocation(login string) *Location {
	resp := fetch(fmt.Sprint("/v2/users/", login, "/locations?filter[active]=true"))

	locations := []Location{}
	json.Unmarshal(resp, &locations)

	if len(locations) > 0 {
		return &locations[len(locations)-1]
	} else {
		return nil
	}
}

// GetUserCoalition() returns the Coalition of an user.
func GetUserCoalition(user *User) *Coalition {
	resp := fetch(fmt.Sprint("/v2/coalitions_users/", "?user_id=", fmt.Sprint(user.ID), "&sort=-created_at"))
	coalition_users := []CoalitionUser{}
	json.Unmarshal(resp, &coalition_users)

	resp = fetch(fmt.Sprint("/v2/users/", user.Login, "/coalitions"))
	coalitions := []Coalition{}
	json.Unmarshal(resp, &coalitions)

	if len(coalitions) > 0 {
		for i, n := range coalitions {
			if n.ID == coalition_users[0].CoalitionID {
				return &coalitions[i]
			}
		}
	}
	return nil
}
