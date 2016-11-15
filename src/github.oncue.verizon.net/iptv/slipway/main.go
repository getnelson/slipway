package main

import (
  "os"
  "fmt"
  "time"
  // "regexp"
  // "strings"
  "strconv"
  "gopkg.in/urfave/cli.v1"
  // "github.com/google/go-github/github"
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

  // github := github.NewClient(nil)

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
        cli.StringFlag {
          Name:   "dir, d",
          Value:  "",
          Usage:  "directory of .deployable.yml files to upload",
          Destination: &userGithubTag,
        },
      },
      Action:  func(c *cli.Context) error {
        credentials, err := loadGithubCredentials();
        if err == nil {
          fmt.Println(credentials);
        } else {
          return cli.NewExitError("Unable to load github credentials. Please ensure you have a valid properties file at $HOME/.github", 1)
        }
        return nil
      },
    },
  }

  // run it!
  app.Run(os.Args)
}