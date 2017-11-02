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
	"gopkg.in/magiconair/properties.v1"
	"io/ioutil"
	"os"
	"regexp"
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
	re, _ := regexp.Compile("^.*(-[0-9.]+)$")

	arr := strings.Split(ctr, ":")
	image := arr[0]
	tag := arr[len(arr)-1]
	name := strings.Split(image, "/")
	last := name[len(name)-1]
	fev := re.FindStringSubmatch(last)[1]
	sanitzed := strings.Replace(last, fev, "", 1)

	return sanitzed, tag
}
