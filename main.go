package main

import (
	"fmt"
	"github.com/hjertnes/new-blog-post/textinput"
	"github.com/hjertnes/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getEditor() string{
	editor, isSet := os.LookupEnv("EDITOR")

	if !isSet{
		return "emacs"
	}

	return editor
}

func getPath() string{
	path, isSet := os.LookupEnv("HUGO_POST_COLLECTION")

	if !isSet{
		return utils.ExpandTilde("~/Code/blog/content/posts")
	}

	return path
}

func main(){
	statusOpt := false
	for i := range os.Args{
		if os.Args[i] == "--status"{
			statusOpt = true
		}
	}
	now := time.Now()
	dateString := now.Format("2006-01-02")
	dateTimeString := now.Format("2006-01-02T15:04:05-07:00")
	path := getPath()
	filenamePre := fmt.Sprintf("%s/%s-", path, dateString)
	filename, err := textinput.Run(fmt.Sprintf("Filename: %s", filenamePre))
	if err != nil{
		panic(err)
	}

	filename = fmt.Sprintf("%s%s.md", filenamePre, filename)

	title := ""

	if !statusOpt{
		title, err = textinput.Run("Title: ")
		if err != nil{
			panic(err)
		}

		if title == " "{
			title = ""
		}
	}

	
	if utils.FileExist(filename){
		fmt.Println("File exists, pick another a filename")
		os.Exit(1)
	}

	data := []string{
		"---",
		fmt.Sprintf("date: \"%s\"", dateTimeString),
		fmt.Sprintf("title: \"%s\"", title),
		fmt.Sprintf("type: \"post\""),
		"---",
		"",
	}


	err = ioutil.WriteFile(filename, []byte(strings.Join(data, "\n")), 0600)
	if err != nil{
		panic(err)
	}

	editor := getEditor()
	cmd := exec.Command(editor, filename)

	err = cmd.Start()
	if err != nil{
		panic(err)
	}
}
