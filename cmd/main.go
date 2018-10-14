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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/google/go-github/github"
	"gopkg.in/urfave/cli.v1"

	nelson "github.com/getnelson/slipway/nelson"
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
	app.Copyright = "© " + strconv.Itoa(year) + " Nelson Team"
	app.Usage = "generate metadata and releases compatible with Nelson"
	app.EnableBashCompletion = true

	const NLDP string = "nldp"
	const YAML string = "yml"

	// switches for the cli
	var (
		userDirectory       string
		userGithubHost      string
		userGithubTag       string
		userGithubRepoSlug  string
		credentialsLocation string
		targetBranch        string
		genEncodingMode     string
		isDryRun            bool
	)

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
				cli.StringFlag{
					Name:        "format, f",
					Value:       YAML,
					Usage:       "Encoding format to use; present options are '" + YAML + "' or '" + NLDP + "'",
					Destination: &genEncodingMode,
				},
			},
			Action: func(c *cli.Context) error {
				ctr := strings.TrimSpace(c.Args().First())
				if len(ctr) <= 0 {
					return cli.NewExitError("You must specify the name of a docker container in order to generate Nelson deployable yml.", 1)
				}

				if genEncodingMode != YAML && genEncodingMode != NLDP {
					return cli.NewExitError("When specifying an encoding, only '"+YAML+"' or '"+NLDP+"' are allowed. The '"+genEncodingMode+"' type is not supported.", 1)
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

				var outputPath string
				var encoded []byte

				if genEncodingMode == YAML {
					outputPath = canonicalDir + name + ".deployable." + YAML
					yaml := "---\n" +
						"name: " + name + "\n" +
						"version: " + tag + "\n" +
						"output:\n" +
						"  kind: docker\n" +
						"  image: " + ctr
					encoded = []byte(yaml)
				} else if genEncodingMode == NLDP {
					outputPath = canonicalDir + name + ".deployable." + NLDP

					deployable, e := newProtoDeployable(ctr, name, tag)
					if e != nil {
						printTerminalErrors(e)
						return cli.NewExitError("Fatal error whilst generating NLDP format.", 1)
					}
					data, ex := proto.Marshal(deployable)
					if ex != nil {
						return cli.NewExitError("Unable to encode deplyoable in binary format.", 1)
					}

					encoded = data
				}

				fmt.Println("Writing to " + outputPath + "...")
				ioutil.WriteFile(outputPath, encoded, 0644)

				return nil
			},
		},
		{
			Name:  "release",
			Usage: "create a Github release for the given repo + tag",
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
					Usage:       "the repository in question, e.g. getnelson/nelson",
					EnvVar:      "TRAVIS_REPO_SLUG",
					Destination: &userGithubRepoSlug,
				},
				cli.StringFlag{
					Name:        "tag, t",
					Value:       "",
					Usage:       "Git tag to use for this release",
					Destination: &userGithubTag,
				},
				cli.StringFlag{
					Name:        "dir, d",
					Value:       "",
					Usage:       "Path to the directory where *.deployable.yml files to upload can be found",
					EnvVar:      "PWD",
					Destination: &userDirectory,
				},
				cli.StringFlag{
					Name:        "creds, c",
					Usage:       "GitHub credentials file",
					Destination: &credentialsLocation,
				},
				cli.StringFlag{
					Name:        "branch, b",
					Value:       "",
					Usage:       "Branch to base release off from",
					Destination: &targetBranch,
				},
				cli.BoolFlag{
					Name:        "dry",
					Usage:       "Is this a dry run or not",
					Destination: &isDryRun,
				},
			},
			Action: func(c *cli.Context) error {
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

				deployablePaths, direrr := findDeployableFilesInDir(userDirectory, YAML)

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

				credentials, errs := getRuntimeCredentials(credentialsLocation)
				if errs != nil {
					return cli.NewMultiError(errs...)
				}

				gh := buildGithubClient(userGithubHost, credentials)

				name := GenerateRandomName()
				isDraft := true

				// release structure
				r := github.RepositoryRelease{
					TagName:         &userGithubTag,
					TargetCommitish: &targetBranch,
					Name:            &name,
					Draft:           &isDraft,
				}

				if isDryRun {
					fmt.Println("The following release payload would be sent to Github:")
					j, _ := json.Marshal(r)
					fmt.Println(string(j))
				} else {
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
						fmt.Println("Error encountered, cleaning up release...")
						gh.Repositories.DeleteRelease(owner, reponame, *release.ID)
						return cli.NewExitError("Unable to promote this release to an offical release. Please ensure that the no other release references the same tag.", 1)
					}
				}
				return nil
			},
		},
		{
			Name:  "deploy",
			Usage: "create a Github deployment for a given repository",
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
					Usage:       "the repository in question, e.g. getnelson/nelson",
					EnvVar:      "TRAVIS_REPO_SLUG",
					Destination: &userGithubRepoSlug,
				},
				cli.StringFlag{
					Name:        "ref, t, s",
					Value:       "",
					Usage:       "Git tag to use for this release",
					Destination: &userGithubTag,
				},
				cli.StringFlag{
					Name:        "dir, d",
					Value:       "",
					Usage:       "Path to the directory where *.deployable.yml files to upload can be found",
					EnvVar:      "PWD",
					Destination: &userDirectory,
				},
				cli.StringFlag{
					Name:        "creds, c",
					Usage:       "GitHub credentials file",
					Destination: &credentialsLocation,
				},
				cli.BoolFlag{
					Name:        "dry",
					Usage:       "Is this a dry run or not",
					Destination: &isDryRun,
				},
			},
			Action: func(c *cli.Context) error {
				if len(userGithubTag) <= 0 {
					return cli.NewExitError("You must specifiy a `--ref`, `-s` or `-t` with a git references (SHA, tag name or branch name).", 1)
				}
				if len(userGithubRepoSlug) <= 0 {
					return cli.NewExitError("You must specifiy a `--repo` or a `-r` denoting the repository to deploy from.", 1)
				}

				splitarr := strings.Split(userGithubRepoSlug, "/")
				if len(splitarr) != 2 {
					return cli.NewExitError("The specified repository name was not of the format 'foo/bar'", 1)
				}

				owner := splitarr[0]
				reponame := splitarr[1]

				deployablePaths, direrr := findDeployableFilesInDir(userDirectory, NLDP)

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

				credentials, errs := getRuntimeCredentials(credentialsLocation)
				if errs != nil {
					return cli.NewMultiError(errs...)
				}

				gh := buildGithubClient(userGithubHost, credentials)

				// here we're taking any of the deployable files that we can find
				// in the specified directory and then encoding them as base64,
				// before packing them into a JSON array (as required by the Github
				// deployment api), with the intention that Nelson gets this payload
				// and can decode the deployables as-is, using its existing decoders.
				encodedDeployables := []*nelson.Deployable{}
				for _, path := range deployablePaths {
					cnts, err := ioutil.ReadFile(path)
					if err != nil {
						return cli.NewExitError("Could not read deployable at "+path, 1)
					}
					o := &nelson.Deployable{}
					if err := proto.Unmarshal(cnts, o); err != nil {
						return cli.NewExitError("Could not unmarshal deployable at '"+path+"'", 1)
					}
					encodedDeployables = append(encodedDeployables, o)
				}

				// package the deployables into our superstructure
				// analog of a Seq[Deployable] in protobuf
				all := &nelson.Deployables{Deployables: encodedDeployables}

				data, ex := proto.Marshal(all)
				if ex != nil {
					return cli.NewExitError("Unable to encode deplyoable in binary format.", 1)
				}

				// take our protobuf byte array and encode it into base64
				// so that we can pretend its a JSON string
				encoded := base64.StdEncoding.EncodeToString(data)

				task := "deploy"
				// yes my pretty, a JSON string you will be
				payload := string(encoded)

				r := github.DeploymentRequest{
					Ref:     &userGithubTag,
					Task:    &task,
					Payload: &payload,
				}

				if isDryRun {
					fmt.Println("The following payload would be sent to Github:")
					j, _ := json.Marshal(r)
					fmt.Println(string(j))
				} else {
					deployment, _, errors := gh.Repositories.CreateDeployment(owner, reponame, &r)
					if errors != nil {
						printTerminalErrors([]error{errors})
						return cli.NewExitError("Fatal error encountered whilst creating deployment.", 1)
					}
					fmt.Println("Created deployment " + strconv.Itoa(*deployment.ID) + " on " + owner + "/" + reponame)
				}
				return nil
			},
		},
	}

	// run it!
	app.Run(os.Args)
}

func getRuntimeCredentials(credentialsLocation string) (Credentials, []error) {
	var credentials Credentials
	envUser := os.Getenv("GITHUB_USERNAME")
	envToken := os.Getenv("GITHUB_TOKEN")

	// if the user did not explictly tell us where the credentials file
	// is located, and we have a GITHUB_TOKEN in the environment, lets use
	// the GITHUB_TOKEN and GITHUB_USERNAME
	if len(envUser) > 0 &&
		len(envToken) > 0 &&
		len(credentialsLocation) < 1 {
		credentials = Credentials{
			Username: envUser,
			Token:    envToken,
		}
	} else if len(credentialsLocation) > 0 {
		loaded, err := loadGithubCredentials(credentialsLocation)
		if err == nil {
			credentials = loaded
		}
	} else {
		errs := []error{}
		errs = append(errs, cli.NewExitError("Slipway requires credentials either in the environment (GITHUB_USERNAME and GITHUB_TOKEN) or specified with a file path using the -c flag.", 1))
		return credentials, errs
	}

	return credentials, nil
}
