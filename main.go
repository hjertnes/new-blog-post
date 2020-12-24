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
	now := time.Now()
	dateString := now.Format("2006-01-02")
	dateTimeString := now.Format("2006-01-02T15:04:05-07:00")
	path := getPath()
	filename, err := textinput.Run("Filename", "Filename", fmt.Sprintf("%s/%s-", path, dateString))
	if err != nil{
		panic(err)
	}

	title, err := textinput.Run("Title", "Title", "")
	if err != nil{
		panic(err)
	}

	if utils.FileExist(filename){
		fmt.Println("File exists, pick another a filename")
		os.Exit(1)
	}

	data := []string{
		"---",
		fmt.Sprintf("date: \"%s\"", dateTimeString),
		fmt.Sprintf("title: \"%s\"", title),
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
