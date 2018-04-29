package authstate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/mhoc/msgoraph/client"
)

type AuthState struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
	ClientID             string
	ClientSecret         string
	RefreshToken         string
}

func Directory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, ".msgraphcli"), nil
}

func Path() (string, error) {
	dir, err := Directory()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "auth.json"), nil
}

func Dump(client *client.Web) error {
	authState := AuthState{
		AccessToken:          client.Credentials().AccessToken,
		AccessTokenExpiresAt: client.Credentials().AccessTokenExpiresAt,
		ClientID:             client.ApplicationID,
		ClientSecret:         client.ApplicationSecret,
		RefreshToken:         client.RefreshToken,
	}
	b, err := json.Marshal(authState)
	if err != nil {
		return err
	}
	authStateFilePath, err := Directory()
	if err != nil {
		return err
	}
	err = os.Mkdir(authStateFilePath, 0777)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return err
	}
	authStateFileName, err := Path()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(authStateFileName, b, 0777)
	if err != nil {
		return err
	}
	return nil
}

func Load() (*client.Web, error) {
	authStateFileName, err := Path()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(authStateFileName)
	if err != nil {
		return nil, err
	}
	as := AuthState{}
	err = json.Unmarshal(b, as)
	if err != nil {
		return nil, err
	}
	c := &client.Web{
		ApplicationID:     as.ClientID,
		ApplicationSecret: as.ClientSecret,
		RefreshToken:      as.RefreshToken,
		RequestCredentials: &client.RequestCredentials{
			AccessToken:          as.AccessToken,
			AccessTokenExpiresAt: as.AccessTokenExpiresAt,
		},
	}
	err = c.RefreshCredentials()
	if err != nil {
		return nil, err
	}
	err = Dump(c)
	return c, err
}
