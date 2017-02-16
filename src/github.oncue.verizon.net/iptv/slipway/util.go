package main

import (
	"gopkg.in/magiconair/properties.v1"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Credentials struct {
	Username string
	Token    string
}

func loadGithubCredentials() (cred Credentials, err []error) {
	location := os.Getenv("HOME") + "/.github"

	if _, err := os.Stat(location); err == nil {
		p := properties.MustLoadFile(location, properties.UTF8)
		user := p.MustGetString("github.login")
		tokn := p.MustGetString("github.token")

		c := Credentials{user, tokn}
		return c, nil
	} else {
		return Credentials{"unknown", "unknown"}, []error{err}
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func findDeployableFilesInDir(path string) ([]string, error) {
	if _, e := os.Stat(path); e == nil {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}

		desired := []string{}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".deployable.yml") && file.IsDir() == false {
				desired = append(desired, path+"/"+file.Name())
			}
		}

		return desired, nil

	} else {
		return nil, e
	}
}

/*
 * this function is pretty basic, and really only splits
 * the following case:
 * your.docker.com/foo/bar:1.2.3
 *
 * The author fully reconizes that this is fucking janky
 * and does not cover all the possible use cases for a docker
 * container name.
 */
func getUnitNameFromDockerContainer(ctr string) (a string, b string) {
	arr := strings.Split(ctr, ":")
	image := arr[0]
	tag := arr[len(arr)-1]
	name := strings.Split(image, "/")
	last := name[len(name)-1]
	return last, tag
}
