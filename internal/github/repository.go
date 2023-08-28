package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	git "github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

var client *http.Client
var token string

var page_size int

type repository struct {
	Name     string `json:"name"`
	Fullname string `json:"full_name"`
	CloneUrl string `json:"clone_url"`
}

func List(org string, prefix string) error {
	repositories, rerr := repositories(org, prefix)
	if rerr != nil {
		return rerr
	}

	fmt.Println("")
	fmt.Println(org)

	for _, v := range repositories {
		fmt.Printf("|-- %s\n", v.Name)
	}

	fmt.Println("")
	fmt.Printf("%d repositories\n", len(repositories))

	return nil
}

func Clone(org string, prefix string) error {
	repositories, rerr := repositories(org, prefix)
	if rerr != nil {
		return rerr
	}

	path, wderr := os.Getwd()
	if wderr != nil {
		return wderr
	}

	fmt.Println("")
	fmt.Println(org)

	for _, r := range repositories {
		cpath := fmt.Sprintf("%s/%s", path, r.Name)

		_, cerr := git.PlainClone(cpath, false, &git.CloneOptions{
			Auth: &githttp.BasicAuth{
				Username: "git",
				Password: token,
			},
			URL:      r.CloneUrl,
			Progress: os.Stdout,
		})
		if cerr != nil {
			if cerr == git.ErrRepositoryAlreadyExists {
				fmt.Printf("|-- %s (exists)\n", r.Name)
			} else {
				fmt.Printf("|-- %s (error)\n", r.Name)
				fmt.Printf("|  -- %s\n", cerr.Error())
			}
		} else {
			fmt.Printf("|-- %s (cloned)\n", r.Name)
		}
	}

	fmt.Println("")
	fmt.Printf("%d repositories\n", len(repositories))

	return nil
}

func Pull() error {
	path, wderr := os.Getwd()
	if wderr != nil {
		return wderr
	}

	fmt.Println("")
	fmt.Println("Pulling latest")

	entries, derr := os.ReadDir(path)
	if derr != nil {
		fmt.Printf("|-- %s (error)\n", path)
		fmt.Printf("|  -- %s\n", derr.Error())
		return derr
	}

	//TODO properly handle errors
	for _, e := range entries {
		if e.IsDir() {
			r, oerr := git.PlainOpen(fmt.Sprintf("%s/%s", path, e.Name()))
			if oerr == nil {
				w, werr := r.Worktree()
				if werr == nil {
					s, serr := w.Status()
					if serr != nil {
						fmt.Printf("|-- %s (status error)\n", path)
						fmt.Printf("|  -- %s\n", serr.Error())
					} else {
						if s.IsClean() {
							err := w.Pull(&git.PullOptions{
								Auth: &githttp.BasicAuth{
									Username: "git",
									Password: token,
								},
							})
							if err == nil {
								fmt.Printf("|-- %s (updated)\n", e.Name())
							} else {
								if err == git.NoErrAlreadyUpToDate {
									fmt.Printf("|-- %s (already up to date)\n", e.Name())
								} else {
									fmt.Printf("|-- %s (error)\n", e.Name())
									fmt.Printf("|  -- %s\n", err.Error())
								}
							}
						} else {
							fmt.Printf("|-- %s (not clean, pull skipped)\n", path)
						}
					}
				} else {
					fmt.Printf("|-- %s (error)\n", e.Name())
					fmt.Printf("|  -- %s\n", werr.Error())
				}
			} else {
				fmt.Printf("|-- %s (error)\n", e.Name())
				fmt.Printf("|  -- %s\n", oerr.Error())
			}
		}
	}

	return nil
}

func init() {
	page_size = 100

	token = os.Getenv("GH_TOKEN")
	if token == "" {
		fmt.Println("GH_TOKEN env not found. This is required to authenticate to GitHub.")
		os.Exit(1)
	}

	client = &http.Client{}
}

func repositories(org string, prefix string) ([]repository, error) {
	var err error
	repositories := []repository{}

	has_next_page := true
	page := 1

	for has_next_page {
		req, rerr := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/orgs/%s/repos", org), nil)
		if rerr != nil {
			err = rerr
			break
		}

		req.Header.Add("Accept", "application/vnd.github+json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

		q := req.URL.Query()
		q.Add("per_page", strconv.Itoa(page_size))
		q.Add("page", strconv.Itoa(page))
		q.Add("sort", "full_name")
		req.URL.RawQuery = q.Encode()

		//TODO handle non 200 codes
		resp, derr := client.Do(req)
		if derr != nil {
			err = derr
			break
		}

		var respbody []repository
		resperr := json.NewDecoder(resp.Body).Decode(&respbody)
		if resperr != nil {
			err = resperr
			break
		}

		if len(respbody) > 0 {
			if prefix != "" {
				for _, r := range respbody {
					if strings.HasPrefix(r.Name, prefix) {
						repositories = append(repositories, r)
					}
				}
			} else {
				repositories = append(repositories, respbody...)
			}
		}

		if len(respbody) == page_size {
			page++
		} else {
			has_next_page = false
		}
	}

	return repositories, err
}
