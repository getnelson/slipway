package main

import (
  "os"
  // "fmt"
  "time"
  "strings"
  // "errors"
  "io/ioutil"
  "gopkg.in/magiconair/properties.v1"
)

type Credentials struct {
  Username string
  Token string
}

func loadGithubCredentials() (cred Credentials, err []error) {
  location := os.Getenv("HOME") + "/.github"

  if _, err := os.Stat(location); err == nil {
    p := properties.MustLoadFile(location, properties.UTF8)
    user := p.MustGetString("github.login")
    tokn := p.MustGetString("github.token")

    c := Credentials { user, tokn }
    return c, nil
  } else {
    return Credentials { "unknown", "unknown" }, []error{ err }
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
      // fmt.Println(file.Name())
      if strings.HasSuffix(file.Name(), ".deployable.yml") && file.IsDir() == false {
        desired = append(desired, path+"/"+file.Name())
      }
    }

    return desired, nil

  } else {
    return nil, e
  }
}
