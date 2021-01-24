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

func getSubCommand() string {
	subCommand := os.Args[1]

	if subCommand != "post" && subCommand != "status"{
		return "post"
	}

	return subCommand
}

type options struct {
	photo bool
}

func getOptions() *options{
	o := &options{
		photo: false,
	}

	for i := range os.Args{
		if os.Args[i] == "--photo"{
			o.photo = true
		}
	}

	return o
}

func help(){
	fmt.Println("blog is a utility for creating blog posts in hugo")
	fmt.Println("usage:")
	fmt.Println("  blog post: creates a post with a title")
	fmt.Println("  blog post: creates a post without a title")
	fmt.Println()
	fmt.Println("options:")
	fmt.Println("  --photo will include all the photos in current dir in the post + copy them to your static folder")
	fmt.Println()
}

func writeFile(filename, title, dateTimeString, content string) error{
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
		content,
	}


	err := ioutil.WriteFile(filename, []byte(strings.Join(data, "\n")), 0600)
	if err != nil{
		return err
	}

	editor := getEditor()
	cmd := exec.Command(editor, filename)

	err = cmd.Start()
	if err != nil{
		return err
	}

	return nil
}

func errorHandler(err error){
	if err != nil{
		panic(err)
	}
}

func main(){
	if len(os.Args) == 1{
		help()
		os.Exit(0)
	}

	subCommand := getSubCommand()

	o := getOptions()

	switch subCommand {
	case "post":
		errorHandler(post(o))
	case "status":
		errorHandler(status(o))
	default:
		help()

	}
}

func figureOutStatusFilename(dateString, path string) string{
	i := 1

	for {
		filename := fmt.Sprintf("%s/%s-%v.md", path, dateString, i)
		if !utils.FileExist(filename){
			return filename
		}

		i++
	}
}

func post(o *options) error {
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

	title, err := textinput.Run("Title: ")
	if err != nil{
		return err
	}

	content := ""
	if o.photo{
		content, err = photo()
		if err != nil{
			return err
		}
	}

	err = writeFile(filename, title, dateTimeString, content)
	return err
}

func status(o *options) error{
	now := time.Now()
	dateString := now.Format("2006-01-02")
	dateTimeString := now.Format("2006-01-02T15:04:05-07:00")

	path := getPath()

	filename := figureOutStatusFilename(dateString, path)

	content := ""

	if o.photo{
		var err error
		content, err = photo()
		if err != nil{
			return err
		}
	}

	err := writeFile(filename, "", dateTimeString, content)

	return err
}

func isPhoto(filename string) bool{
	fmt.Println(strings.ToLower(filename))
	if strings.HasSuffix(strings.ToLower(filename), ".jpeg"){
		return true
	}

	if strings.HasSuffix(strings.ToLower(filename), ".jpg"){
		return true
	}

	if strings.HasSuffix(strings.ToLower(filename), ".png"){
		return true
	}

	if strings.HasSuffix(strings.ToLower(filename), ".gif"){
		return true
	}

	return false

}

func getExt(filename string) string{
	elems := strings.Split(filename, ".")

	return elems[len(elems)-1]
}

func photo() (string, error){
	path := getPath()
	staticDir := fmt.Sprintf("%s/../../static", path)

	files, err := ioutil.ReadDir(".")
	if err != nil{
		return "", err
	}

	elements := make([]string, 0)

	for _, f := range files{
		if f.IsDir() {
			continue
		}

		if strings.HasPrefix(f.Name(), "."){
			continue
		}

		if !isPhoto(f.Name()){
			continue
		}

		ext := getExt(f.Name())

		dateString := f.ModTime().Format("2006-01-02")

		photoPath, photoName := generateImageName(staticDir, dateString, ext)
		err = copy(f.Name(), photoPath)
		if err != nil{
			return "", err
		}
		elements = append(elements, fmt.Sprintf("![](/%s)", photoName))
	}

	return strings.Join(elements, "\n"), nil
}

func generateImageName(staticDir, dateString, ext string) (string, string){
	c := 0
	for{
		photoName := fmt.Sprintf("%s-%v.%s", dateString, c, ext)
		photoPath := fmt.Sprintf("%s/%s-%v.%s", staticDir, dateString, c, ext)
		if !utils.FileExist(photoPath){
			return photoPath, photoName
		}

		c++
	}
}

func copy(from, to string) error{
	bytesRead, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(to, bytesRead, 0600)
	if err != nil {
		return err
	}

	return nil
}
