package auth

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "app-setup",
			Usage:  "print directions on how to set up an app to speak with the graph api",
			Action: AppSetup,
		},
		{
			Name:   "login",
			Usage:  "initiate an oauth2 login with the graph api",
			Action: Login,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Usage: "the port to listen to locally for the oauth response",
				},
			},
		},
	}
}

func AppSetup(c *cli.Context) {
	fmt.Printf("To set up a new app with the Microsoft Graph API, follow these instructions:\n")
	fmt.Printf("\thttps://developer.microsoft.com/en-us/graph/docs/concepts/auth_v2_user#1-register-your-app\n")
	fmt.Printf("After that is set up, provide your application's clientID as an environment variable:\n")
	fmt.Printf("\tMSGRAPH_CLIENT_ID=6731de76-14a6-49ae-97bc-6eba6914391e\n")
	fmt.Printf("Finally, execute:\n")
	fmt.Printf("\tmsgraph auth login\n")
	fmt.Printf("And sign-in with your Microsoft account authenticated with the information you want to access\n")
	fmt.Printf("This will create a file at ~/.msgraphcli/auth.json where the token state is stored and used by the cli for future requests\n")
}

func Login(c *cli.Context) {
	clientID := os.Getenv("MSGRAPH_CLIENT_ID")
	if clientID == "" {
		fmt.Printf("MSGRAPH_CLIENT_ID not provided")
		os.Exit(1)
	}
	port := "8537"
	if c.String("port") != "" {
		port = c.String("port")
	}
	queryString := url.Values{}
	queryString.Add("client_id", clientID)
	queryString.Add("redirect_uri", fmt.Sprintf("http://localhost:%v/login", port))
	queryString.Add("response_mode", "query")
	queryString.Add("response_type", "code")
	queryString.Add("scope", "offline_access user.read mail.read")
	queryString.Add("state", "12345")
	openURL := fmt.Sprintf("https://login.microsoftonline.com/common/oauth2/v2.0/authorize?%v", queryString.Encode())
	var err error
	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", openURL).Start()
	case "linux":
		err = exec.Command("xdg-open", openURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", openURL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
