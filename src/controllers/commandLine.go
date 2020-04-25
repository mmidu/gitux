package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"../models"
	"../utils"
)

// JSONStub is the default json structure of the git-credentials.json file
const JSONStub = "{\n\t\"users\": [\n\t\t{\n\t\t\t\"identifier\": \"github-work\",\n\t\t\t\"credentials\": {\n\t\t\t\t\"username\": \"username\",\n\t\t\t\t\"password\": \"password1\",\n\t\t\t\t\"domain\": \"github.com\"\n\t\t\t}\n\t\t},\n\t\t{\n\t\t\t\"identifier\": \"two\",\n\t\t\t\"credentials\": {\n\t\t\t\t\"username\": \"uname2\",\n\t\t\t\t\"password\": \"password2\",\n\t\t\t\t\"domain\": \"domain\"\n\t\t\t}\n\t\t}\n\t]\n}"

// Init initializes the JSON stub file
func Init() {
	path := fmt.Sprintf("%s/git-credentials.json", utils.GetHomeDir())
	if !utils.FileExists(path) {
		file, err := os.Create(path)
		utils.Check(err)
		defer file.Close()

		w := bufio.NewWriter(file)
		_, err = w.WriteString(JSONStub)
		utils.Check(err)
		w.Flush()
	} else {
		fmt.Println(fmt.Sprintf("%s/git-credentials.json file already exists.", utils.GetHomeDir()))
	}

	os.Exit(0)
}

// ManageArguments get arguments from the command line and executes the right functions
func ManageArguments() string {
	var identifier string
	for i := 0; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-i", "--init":
			Init()
		case "-s", "--set":
			if i+1 < len(os.Args) {
				identifier = os.Args[i+1]
			}
		case "-lu", "--listUsers":
			printCredentials()
		case "-c", "--current":
			getCurrent()
		case "-h", "--help":
			fmt.Println("HELP STRING")
			os.Exit(0)
		}
	}
	return identifier
}

func printCredentials() {
	users := getCredentials()

	fmt.Println(users.List())

	os.Exit(0)
}

func getCredentials() models.Users {
	jsonFile, err := os.Open(fmt.Sprintf("%s/git-credentials.json", utils.GetHomeDir()))
	utils.Check(err)
	defer jsonFile.Close()

	var users models.Users

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &users)

	return users
}

func getCurrent() {
	credentialsFile, err := os.Open(fmt.Sprintf("%s/.git-credentials", utils.GetHomeDir()))
	utils.Check(err)
	defer credentialsFile.Close()

	reader := bufio.NewReader(credentialsFile)
	data, err := reader.ReadString('\n')
	utils.Check(err)

	fl := data[8:]

	resA := strings.SplitN(fl, ":", -1)

	username := resA[0]

	resB := strings.SplitN(resA[1], "@", -1)

	password, _ := url.QueryUnescape(resB[0])

	domain := strings.TrimSuffix(resB[1], "\n")

	users := getCredentials()

	for _, user := range users.Users {
		if user.Credentials.Username == username && user.Credentials.Password == password && user.Credentials.Domain == domain {
			fmt.Println(user.Identifier)
			os.Exit(0)
		}
	}

	os.Exit(1)
}
