package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

func main() {
	terms := os.Args[1:]
	fmt.Println(terms)
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d тем:\n", result.TotalCount)
	var (
		lessMonth []string
		lessYear []string
		moreYear []string
	)
	func (){
		for _, item := range result.Items {
			data := fmt.Sprintf("#%-5d %9.9s %.55s\n",item.Number, item.User.Login, item.Title)
			if item.CreatedAt.After(time.Now().AddDate(0,-1,0)){
				lessMonth = append(lessMonth, data)
			} else if item.CreatedAt.After(time.Now().AddDate(-1,0,0)){
				lessYear = append(lessYear, data)
			} else {
				moreYear = append(moreYear, data)
			}
		}
	}()
	
	fmt.Println("Less than month ago created: ")
	for _,item := range lessMonth {
		fmt.Println(item)
	}
	fmt.Println("Less than year ago created: ")
	for _,item := range lessYear {
		fmt.Println(item)
	}
	fmt.Println("More than year ago created: ")
	for _,item := range moreYear {
		fmt.Println(item)
	}
}


