package main

import (
  "os"
  "fmt"
  "time"
  // "regexp"
  // "strings"
  "strconv"
  "gopkg.in/urfave/cli.v1"
)

var globalBuildVersion string

func CurrentVersion() string {
  if len(globalBuildVersion) == 0 {
    return "devel"
  } else {
    return "v"+globalBuildVersion
  }
}

func main() {
  year, _, _ := time.Now().Date()
  app := cli.NewApp()
  app.Name = "slipway"
  app.Version = CurrentVersion()
  app.Copyright = "© "+strconv.Itoa(year)+" Verizon Labs"
  app.Usage = "generate metadata and releases compatible with Nelson"
  app.EnableBashCompletion = true

  // pi   := ProgressIndicator()

  // switches for the cli
  var userGithubToken string
  var userDirectory string
  var userGithubHost string
  var userGithubTag string

  app.Commands = []cli.Command {
    ////////////////////////////// DEPLOYABLE //////////////////////////////////
    {
      Name:    "gen",
      Usage:   "generate deployable metdata for units",
      Flags: []cli.Flag {
        cli.StringFlag{
          Name:   "dir, d",
          Value:  "",
          Usage:  "location to output the YAML file",
          Destination: &userDirectory,
        },
      },
      Action:  func(c *cli.Context) error {
        fmt.Println("testing")
        return nil
      },
    },
    {
      Name:    "release",
      Usage:   "generate deployable metdata for units",
      Flags: []cli.Flag {
        cli.StringFlag {
          Name:   "auth, a",
          Value:  "",
          Usage:  "your github personal access token",
          EnvVar: "GITHUB_TOKEN",
          Destination: &userGithubToken,
        },
        cli.StringFlag {
          Name:   "endpoint, x",
          Value:  "",
          Usage:  "host of the github api endpoint",
          EnvVar: "GITHUB_ADDR",
          Destination: &userGithubHost,
        },
        cli.StringFlag {
          Name:   "tag, t",
          Value:  "",
          Usage:  "host of the github api endpoint",
          Destination: &userGithubTag,
        },
      },
      Action:  func(c *cli.Context) error {
        fmt.Println("testing")
        return nil
      },
    },
  }

  // run it!
  app.Run(os.Args)
}