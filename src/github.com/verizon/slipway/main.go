//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//:
//:   Licensed under the Apache License, Version 2.0 (the "License");
//:   you may not use this file except in compliance with the License.
//:   You may obtain a copy of the License at
//:
//:       http://www.apache.org/licenses/LICENSE-2.0
//:
//:   Unless required by applicable law or agreed to in writing, software
//:   distributed under the License is distributed on an "AS IS" BASIS,
//:   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//:   See the License for the specific language governing permissions and
//:   limitations under the License.
//:
//: ----------------------------------------------------------------------------
package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var globalBuildVersion string

func currentVersion() string {
	if len(globalBuildVersion) == 0 {
		return "devel"
	} else {
		return "v" + globalBuildVersion
	}
}

func main() {
	year, _, _ := time.Now().Date()
	app := cli.NewApp()
	app.Name = "slipway"
	app.Version = currentVersion()
	app.Copyright = "© " + strconv.Itoa(year) + " Verizon Labs"
	app.Usage = "generate metadata and releases compatible with Nelson"
	app.EnableBashCompletion = true

	// switches for the cli
	var userDirectory string
	var userGithubHost string
	var userGithubTag string
	var userGithubRepoSlug string

	app.Commands = []cli.Command{
		////////////////////////////// DEPLOYABLE //////////////////////////////////
		{
			Name:  "gen",
			Usage: "generate deployable metdata for units",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "dir, d",
					Value:       "",
					Usage:       "location to output the YAML file",
					Destination: &userDirectory,
				},
			},
			Action: func(c *cli.Context) error {
				ctr := strings.TrimSpace(c.Args().First())
				if len(ctr) <= 0 {
					return cli.NewExitError("You must specify the name of a docker container in order to generate Nelson deployable yml.", 1)
				}

				if len(userDirectory) <= 0 {
					pwd, _ := os.Getwd()
					if len(pwd) <= 0 {
						return cli.NewExitError("You must specify a '--dir' or '-d' flag with the destination directory for the deployable yml file.", 1)
					} else {
						userDirectory = pwd
						fmt.Println("No destination folder specified. Assuming the current working directory.")
					}
				} else {
					if _, err := os.Stat(userDirectory); err != nil {
						return cli.NewExitError("The specified directory "+userDirectory+" does not exist or cannot be accessed.", 1)
					}
				}

				var canonicalDir string
				if strings.HasSuffix(userDirectory, "/") {
					canonicalDir = userDirectory
				} else {
					canonicalDir = userDirectory + "/"
				}

				name, tag := getUnitNameFromDockerContainer(ctr)

				yaml := "---\n" +
					"name: " + name + "\n" +
					"version: " + tag + "\n" +
					"output:\n" +
					"  kind: docker\n" +
					"  image: " + ctr

				outputPath := canonicalDir + name + ".deployable.yml"

				fmt.Println("Writing to " + outputPath + "...")
				ioutil.WriteFile(outputPath, []byte(yaml), 0644)

				return nil
			},
		},
		{
			Name:  "release",
			Usage: "generate deployable metdata for units",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "endpoint, x",
					Value:       "",
					Usage:       "domain of the github api endpoint",
					EnvVar:      "GITHUB_ADDR",
					Destination: &userGithubHost,
				},
				cli.StringFlag{
					Name:        "repo, r",
					Value:       "",
					Usage:       "the repository in question, e.g. verizon/knobs",
					EnvVar:      "TRAVIS_REPO_SLUG",
					Destination: &userGithubRepoSlug,
				},
				cli.StringFlag{
					Name:        "tag, t",
					Value:       "",
					Usage:       "host of the github api endpoint",
					Destination: &userGithubTag,
				},
				cli.StringFlag{
					Name:        "dir, d",
					Value:       "",
					Usage:       "directory of .deployable.yml files to upload",
					Destination: &userDirectory,
				},
			},
			Action: func(c *cli.Context) error {
				// deployables =

				if len(userGithubTag) <= 0 {
					return cli.NewExitError("You must specifiy a `--tag` or a `-t` to create releases.", 1)
				}
				if len(userGithubRepoSlug) <= 0 {
					return cli.NewExitError("You must specifiy a `--repo` or a `-r` to create releases.", 1)
				}

				splitarr := strings.Split(userGithubRepoSlug, "/")
				if len(splitarr) != 2 {
					return cli.NewExitError("The specified repository name was not of the format 'foo/bar'", 1)
				}

				owner := splitarr[0]
				reponame := splitarr[1]

				deployablePaths, direrr := findDeployableFilesInDir(userDirectory)

				if len(userDirectory) != 0 {
					// if you specified a dir, but it was not readable or it didnt exist
					if direrr != nil {
						return cli.NewExitError("Unable to read from "+userDirectory+"; check the location exists and is readable.", 1)
					}
					// if you specify a dir, and it was readable, but there were no deployable files
					if len(deployablePaths) <= 0 {
						return cli.NewExitError("Readable directory "+userDirectory+" contained no '.deployable.yml' files.", 1)
					}
				}

				credentials, err := loadGithubCredentials()
				if err == nil {
					gh := buildGithubClient(userGithubHost, credentials)

					name := GenerateRandomName()
					commitish := "master"
					isDraft := true

					// release structure
					r := github.RepositoryRelease{
						TagName:         &userGithubTag,
						TargetCommitish: &commitish,
						Name:            &name,
						Draft:           &isDraft,
					}

					// create the release
					release, _, e := gh.Repositories.CreateRelease(owner, reponame, &r)

					if e != nil {
						fmt.Println(e)
						return cli.NewExitError("Encountered an unexpected error whilst calling the specified Github endpint. Does Travis have permissions to your repository?", 1)
					} else {
						fmt.Println("Created release " + strconv.Itoa(*release.ID) + " on " + owner + "/" + reponame)
					}

					// upload the release assets
					for _, path := range deployablePaths {
						slices := strings.Split(path, "/")
						name := slices[len(slices)-1]
						file, _ := os.Open(path)

						fmt.Println("Uploading " + name + " as a release asset...")

						opt := &github.UploadOptions{Name: name}
						gh.Repositories.UploadReleaseAsset(owner, reponame, *release.ID, opt, file)
					}

					fmt.Println("Promoting release from a draft to offical release...")

					// mutability ftw?
					isDraft = false
					r2 := github.RepositoryRelease{
						Draft: &isDraft,
					}

					_, _, xxx := gh.Repositories.EditRelease(owner, reponame, *release.ID, &r2)

					if xxx != nil {
						fmt.Println(xxx)
						return cli.NewExitError("Unable to promote this release to an offical release. Please ensure that the no other release references the same tag.", 1)
					}

				} else {
					fmt.Println(err)
					return cli.NewExitError("Unable to load github credentials. Please ensure you have a valid properties file at $HOME/.github", 1)
				}

				return nil
			},
		},
	}

	// run it!
	app.Run(os.Args)
}
