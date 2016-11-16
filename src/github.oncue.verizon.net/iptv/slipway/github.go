package main

import (
  "net/url"
  "strings"
  "github.com/google/go-github/github"
)

func buildGithubClient(domain string, credentials Credentials) (gh *github.Client) {
  if domain == "" {
    return github.NewClient(nil)
  } else {
    tp := github.BasicAuthTransport {
      Username: strings.TrimSpace(credentials.Username),
      Password: strings.TrimSpace(credentials.Token),
    }

    client := github.NewClient(tp.Client())

    u, _ := url.Parse("https://"+domain+"/api/v3/")
    client.BaseURL = u

    uu, _ := url.Parse("https://"+domain+"/api/uploads/")
    client.UploadURL = uu

    return client
  }
}
