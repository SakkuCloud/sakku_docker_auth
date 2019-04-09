package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sakku_docker_auth/util"
	"net/http"
	"os"
	"strings"
)

type AuthAPIRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func main() {
	credentials := strings.Split(util.ReadStdIn(), " ")
	if len(credentials) != 2 {
		fmt.Println("Error Code 1: Cannot parse the input")
		os.Exit(util.ErrorExitCode)
	}

	authReq := AuthAPIRequest{}
	authReq.Email = credentials[0]
	authReq.Token = credentials[1]

	values, err := json.Marshal(authReq)
	if err != nil{
		fmt.Println("Error Code 2: Cannot parse the input")
		os.Exit(util.ErrorExitCode)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.sakku.cloud/service/users/authenticate", bytes.NewBuffer(values))
	if err != nil{
		fmt.Println("Error Code 3: Cannot create request to Sakku authenticate server")
		os.Exit(util.ErrorExitCode)
	}

	req.Header.Add("Content-Type","Application/JSON")
	req.Header.Add("service","123")
	req.Header.Add("service-key","123")
	rsp, err := client.Do(req)
	if err != nil{
		fmt.Println("Error Code 4: Cannot connect to Sakku authenticate server")
		os.Exit(util.ErrorExitCode)
	}

	if rsp.StatusCode == http.StatusOK {
		os.Exit(util.SuccessExitCode)
	} else {
		os.Exit(util.ErrorExitCode)
	}

}
