package auth

import (
	"log"
	"strconv"

	"github.com/mhoc/msgoraph/client"
	"github.com/mhoc/msgoraph/scopes"
	"github.com/mhoc/msgraph-cli/authstate"
	"github.com/urfave/cli"
)

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "login",
			Usage:  "initiate an oauth2 login with the graph api",
			Action: Login,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "client-id",
					Usage: "the client id of the configured msgraph application",
				},
				cli.StringFlag{
					Name:  "client-secret",
					Usage: "the client secret of the configured msgraph application",
				},
				cli.StringFlag{
					Name:  "port",
					Usage: "the port to listen to locally for the oauth response",
				},
				cli.StringFlag{
					Name:  "scopes",
					Usage: "a space-separate list of permissions to ask for during the authorization",
				},
			},
		},
	}
}

func Login(c *cli.Context) {
	clientID := c.String("client-id")
	if clientID == "" {
		log.Fatal("-client-id not provided")
	}
	clientSecret := c.String("client-secret")
	if clientSecret == "" {
		log.Fatal("-client-secret not provided")
	}
	portS := c.String("port")
	if portS == "" {
		log.Fatal("-port not provided")
	}
	port, err := strconv.Atoi(portS)
	if err != nil {
		log.Fatal(err.Error())
	}
	scopez := c.String("scopes")
	if scopez == "" {
		log.Fatalf("-scopes not provided")
	}
	webClient := client.NewWeb(clientID, clientSecret, port, scopes.Resolve(scopez, scopes.PermissionTypeDelegated))
	err = webClient.InitializeCredentials()
	if err != nil {
		panic(err)
	}
	err = authstate.Dump(webClient)
	if err != nil {
		panic(err)
	}
}
