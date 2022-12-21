package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

    "github.com/fatih/color"
    "github.com/go-git/go-git/v5"
)

var appURL string

func doNew(appName string) {
    appName = strings.ToLower(appName)
    appURL = appName

    // sanitize the application name ( convert url to single word )
    if strings.Contains(appName, "/") {
        exploded := strings.SplitAfter(appName, "/")
        appName = exploded[len(exploded)-1]
    }

    log.Println("App name is", appName)

    // git clone the skeleton application
    color.Green("\tCloning repository...")
    _, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
        URL:      "https://Terminon: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@github.com/Terminon/celeritas-app",
        Progress: os.Stdout,
        Depth:    1,
    })
    if err != nil {
        exitGracefully(err)
    }

    // remove the .git directory
    err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
    if err != nil {
        exitGracefully(err)
    }
    // create a ready to go .env file
    color.Yellow("\tCreating .env file...")
    data, err := templateFS.ReadFile("templates/env.txt")
    if err != nil {
        exitGracefully(err)
    }
    env := string(data)
    env = strings.ReplaceAll(env, "${APP_NAME}", appName)
    env = strings.ReplaceAll(env, "${KEY}", cel.RandomString(32))

    err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
    if err != nil {
        exitGracefully(err)
    }

    // update the go.mod file
    color.Yellow("\tCreating go.mod file...")
    _ = os.Remove("./" + appName + "/go.mod")

    data, err = templateFS.ReadFile("templates/go.mod.txt")
    if err != nil {
        exitGracefully(err)
    }
    mod := string(data)
    mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)
    err = copyDataToFile([]byte(mod), "./"+appName+"/go.mod")
    if err != nil {
        exitGracefully(err)
    }

    // update existing .go files with correct /name/imports
    color.Yellow("\tUpdating source files...")
    os.Chdir("./" + appName)
    updateSource()

    // run go mod tidy in the project directory
    color.Yellow("\tRunning go mod tidy")
    cmd := exec.Command("go", "mod", "tidy")
    err = cmd.Start()
    if err != nil {
        exitGracefully(err)
    }
    color.Green("Done building " + appURL)

}
