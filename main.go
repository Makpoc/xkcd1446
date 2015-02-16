package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	REMOTE_HOST          = "http://xkcd1446.org"
	LIST_JSON            = "/list.json"
	REMOTE_IMAGES_FOLDER = "/img/"
	LOCAL_IMAGES_FOLDER  = "./images/"
)

var AllImages struct {
	Images []string
}

func getImageNames() {
	fmt.Println("Getting all image names...")

	resp, err := http.Get(REMOTE_HOST + LIST_JSON)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(body, &AllImages.Images)
	if err != nil {
		panic(err.Error())
	}

}

func fetchImages() {
	fmt.Println("Fetching new images...")

	for _, name := range AllImages.Images {
		localFilePath := LOCAL_IMAGES_FOLDER + name
		if fi, err := os.Stat(localFilePath); err == nil && fi.Size() > 0 {
			//fmt.Printf("File %s already exists. Skipping...\n", localFilePath)
			continue
		}
		imageUrl := REMOTE_HOST + REMOTE_IMAGES_FOLDER + name
		fmt.Printf("Now fetching %s\n", imageUrl)
		resp, err := http.Get(imageUrl)
		if err != nil {
			fmt.Printf("Skipping %s because of error: %s\n", name, err.Error())
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		ioutil.WriteFile(LOCAL_IMAGES_FOLDER+name, body, 0666)
	}
}

func dividableBy(args ...interface{}) bool {
	if len(args) != 2 {
		return false
	}

	x := args[0].(int)
	y := args[1].(int)

	if x%y == 0 {
		return true
	}

	return false
}

func main() {
	go func() {
		for {
			getImageNames()
			fetchImages()
			time.Sleep(1 * time.Minute)
		}
	}()

	fmt.Println("All done. Starting server at http://localhost:8080")

	t := template.Must(template.New("index.html").Funcs(template.FuncMap{"dividableBy": dividableBy}).ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := t.Execute(w, AllImages)
		if err != nil {
			panic(err)
		}
	})
	http.Handle("/images/", http.FileServer(http.Dir("")))
	http.ListenAndServe(":8080", nil)
}
