package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"example.com/m/4_11/github"
	git "gopl.io/ch4/github"
)

func search(query []string) {
	result, err := git.SearchIssues(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d тем:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}

func create(owner,repo string){
	
	editorPath := "C:/Program Files/Git/usr/bin/vim.exe"
	
	tempfile, err := ioutil.TempFile("", "*.json")
	if err != nil {
		log.Fatal(err)
	}
	defer tempfile.Close()
	defer os.Remove(tempfile.Name())
	cmd := &exec.Cmd{
		Path:   editorPath,
		Args:   []string{"", tempfile.Name()},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	tempfile.Seek(0, 0)
	fmt.Scan()
	body, err := ioutil.ReadAll(tempfile)
	if err != nil{
		log.Fatal(err)
		os.Exit(1)
	}
	fields := make(map[string]string)
	err = json.Unmarshal(body,&fields)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Print(fields)
	v := url.Values{}
	for key,value := range fields {
		v.Add(key,value)
	}
	resp, err := http.PostForm(github.IssuesURL+owner+repo,v)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Print(resp)
}

func read(owner,repo,number string){
	resp,err := http.Get(fmt.Sprintf("https://api.github.com/repos%s/%s/issues/%s",owner,repo,number))
	if err != nil {
			log.Fatal(err)
			os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatalf("search query failed: %s", resp.Status)
		os.Exit(1)
	}
	var issue git.Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		panic(err)
	}
	fmt.Printf("repo: %s/%s\nnumber: %s\nuser: %s\ntitle: %s\n\n%s",
		owner, repo, number, issue.User.Login, issue.Title, issue.Body)
}

func main() {
	if len(os.Args) < 2 {
		messageErr()
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	if cmd == "search" {
		if len(args) < 1 {
			messageErr()
		}
		search(args)
		os.Exit(0)
	}
	owner := args[0]
	repo := args[1]
	if len(args) > 2{
		num := args[2]
		switch {
		case cmd == "read":
			read(owner,repo,num)
		}
	}
	switch {
	case cmd == "create":
		create(owner,repo)
	
	default:
		messageErr()
		os.Exit(1)
	}

}

func messageErr() {
	fmt.Print("usage: search QUERY[read|edit|close|create] OWNER REPO ISSUE_NUMBER")
	os.Exit(1)
}
