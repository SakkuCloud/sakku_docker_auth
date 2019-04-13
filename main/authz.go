package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cesanta/docker_auth/auth_server/authz"
	"github.com/sakku_docker_auth/util"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	RepositoryTypeKeyWord = "Repository"
)

type DockerRegAPIRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type DockerRegADIResponse struct {
	Code    string                     `json:"code"`
	Message string                     `json:"message"`
	Error   string                     `json:"error"`
	Result  DockerRegADIResponseResult `json:"result"`
}

type DockerRegADIResponseResult struct {
	Actions string `json:"actions"`
}

func isAuthorized(reqActions []string, rspActions string) bool{
	for _, action := range reqActions {
		if !strings.Contains(rspActions, action) {
			return false
		}
	}
	return true
}

func main() {
	text := util.ReadStdIn()

	var authReqInfo authz.AuthRequestInfo
	err := json.Unmarshal([]byte(text), &authReqInfo)
	if err != nil {
		fmt.Println("Error Code 1: Cannot parse the input")
		os.Exit(util.ErrorExitCode)
	}

	dockerRegReq := DockerRegAPIRequest{}
	dockerRegReq.Type = RepositoryTypeKeyWord
	dockerRegReq.Email = authReqInfo.Account
	dockerRegReq.Name = authReqInfo.Name

	values, err := json.Marshal(dockerRegReq)
	if err != nil{
		fmt.Println("Error Code 2: Cannot parse the input")
		os.Exit(util.ErrorExitCode)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", util.SakkuDockerRegServiceAddr, bytes.NewBuffer(values))
	if err != nil{
		fmt.Println("Error Code 3: Cannot create request to Sakku docker reg authorization server")
		os.Exit(util.ErrorExitCode)
	}

	req.Header.Add("Content-Type","Application/JSON")
	req.Header.Add("service",util.SakkuDockerRegServiceName)
	req.Header.Add("service-key",util.SakkuDockerRegServiceKey)
	rsp, err := client.Do(req)
	if err != nil{
		fmt.Println("Error Code 4: Cannot connect to Sakku docker reg authorization server")
		os.Exit(util.ErrorExitCode)
	}

	rspData, err := ioutil.ReadAll(rsp.Body)
	if err != nil{
		fmt.Println("Error Code 4: Cannot parse data from Sakku docker reg authorization server")
		os.Exit(util.ErrorExitCode)
	}

	dockerRegADIResponse := DockerRegADIResponse{}
	if err := json.Unmarshal(rspData, &dockerRegADIResponse); err != nil {
		fmt.Println("Error Code 5: Cannot parse data from Sakku docker reg authorization server")
		os.Exit(util.ErrorExitCode)
	}

	if isAuthorized(authReqInfo.Actions, dockerRegADIResponse.Result.Actions) {
		os.Exit(util.SuccessExitCode)
	} else {
		os.Exit(util.ErrorExitCode)
	}
}
