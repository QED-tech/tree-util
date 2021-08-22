package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

type File struct {
	Name     string
	Type     string
	Size     int64
	Children []File
}

func (f *File) isDir() bool {
	return f.Type == "dir"
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	fileCollections, err := recursiveRead(path)
	if err != nil {
		return err
	}
	prefix := []string{}
	View := recursiveView(fileCollections, prefix, 0, printFiles)

	fmt.Println(View)
	return nil
}

func recursiveView(files []File, prefixes []string, depth int, printFiles bool) string {
	var View string
	graphic := strings.Join(prefixes, "")
	length := len(files)

	for i, file := range files {
		isLast := i == length-1
		isDir := file.isDir()

		if isDir && isLast {
			View += fmt.Sprintf("%v└───%v\n", graphic, file.Name)
			View += recursiveView(file.Children, append(prefixes, "\t"), depth+1, printFiles)
			continue
		}

		if isDir {
			View += fmt.Sprintf("%v├───%v\n", graphic, file.Name)
			View += recursiveView(file.Children, append(prefixes, "│\t"), depth+1, printFiles)
			continue
		}

		if isLast {
			View += fmt.Sprintf("%v└───%v\n", graphic, file.Name)
			continue
		}

		View += fmt.Sprintf("%v├───%v\n", graphic, file.Name)
	}

	return View
}

func recursiveRead(filepath string) ([]File, error) {

	files, err := ioutil.ReadDir(filepath)
	var filesCollection []File

	if err != nil {
		return filesCollection, err
	}

	for _, file := range files {
		fileName := file.Name()
		fullPath := filepath + string(os.PathSeparator) + fileName
		fi, err := os.Stat(fullPath)

		if err != nil {
			return filesCollection, err
		}

		if file.IsDir() {
			children, err := recursiveRead(fullPath)
			if err != nil {
				return filesCollection, err
			}
			fileData := File{fileName, "dir", fi.Size(), children}
			filesCollection = append(filesCollection, fileData)
			continue
		}

		fileData := File{fileName, "file", fi.Size(), nil}
		filesCollection = append(filesCollection, fileData)
	}

	return filesCollection, nil
}
