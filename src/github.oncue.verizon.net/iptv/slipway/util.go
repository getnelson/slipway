package main

import (
  // "errors"
  "os"
  "time"
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
