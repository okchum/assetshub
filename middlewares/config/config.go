package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

func LoadIni() *ini.File {
	config, err := ini.Load("./site.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return config
}

func GetKey(section string, key string) string {
	config := LoadIni()
	root := config.Section(section).Key(key).String()
	return root
}

func GetAuth() gin.Accounts {
	return gin.Accounts{
		GetKey("auth", "auth_username"): GetKey("auth", "auth_password"),
	}
}

func GetArray(section string, key string, delimiter string) []string {

	if len(delimiter) == 0 {
		delimiter = ","
	}
	string := GetKey(section, key)
	array := strings.Split(string, delimiter)

	return array
}

func GetPort() string {
	return GetKey("site", "port")
}

func GetRoot(key string, hasSlash bool) string {
	root := GetKey("site", key)
	if hasSlash {
		root = fmt.Sprintf("%s/", root)
	}
	return root
}
