package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func readToken() string {
	token, err := ioutil.ReadFile("token.txt")
	if err != nil {
		panic(err)
	}
	return string(token)
}

func main() {
	var query struct {
		Organization struct {
			Repositories struct {
				TotalCount int
				PageInfo   struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
				Edges []struct {
					Node struct {
						Name string
					}
				}
			} `graphql:"repositories(first: 100, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}

	var loginFlag = flag.String("org", "", "Github Organization Name")
	flag.Parse()

	variables := map[string]interface{}{
		"cursor": (*githubv4.String)(nil), // Null after argument to get first page.
		"login":  githubv4.String(*loginFlag),
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: readToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := githubv4.NewClient(tc)

	var sb strings.Builder

	sb.WriteString("GITHUB_REPOS=")

	for {
		err := client.Query(ctx, &query, variables)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, v := range query.Organization.Repositories.Edges {
			fmt.Printf("%s \n", v.Node.Name)
			sb.WriteString("Blueshoe/")
			sb.WriteString(v.Node.Name)
			sb.WriteString(",")
		}
		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Repositories.PageInfo.EndCursor)
	}

	str := sb.String()

	str = str[:len(str)-1]

	ioutil.WriteFile(".env", []byte(str), 0644)

}
