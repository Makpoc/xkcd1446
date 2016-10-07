package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/donovanhide/eventsource"
)

const (
	port                = "8080"
	intervalEnvKey      = "XKCD1446_INTERVAL"
	defaultPushInterval = 2
	landing             = "landing"

	imageListFile = "./images.txt"
)

var (
	pushInterval      time.Duration = defaultPushInterval
	currentImageIndex               = 0
)

var allImageUrls = make(map[int]image)

type image struct {
	seqN  int
	url   string
	image []byte
}

type imageEvent struct {
	id    int
	image image
}

func (ie imageEvent) Id() string    { return fmt.Sprintf("%d", ie.image.seqN) }
func (ie imageEvent) Event() string { return "landing" }
func (ie imageEvent) Data() string  { return base64.StdEncoding.EncodeToString(ie.image.image) }

// loadImageUrls loads all the image urls from the file system
func loadImageUrls() error {
	file, err := os.Open(imageListFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Stat(); os.IsNotExist(err) {
		return errors.New(imageListFile + " not found!")
	}

	curLine := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allImageUrls[curLine] = image{seqN: curLine, url: scanner.Text()}
		curLine += 1
	}

	if curLine == 0 {
		return errors.New(imageListFile + " is empty!")
	}

	return nil
}

// pushImage regulary pushes the next image to the server
func pushImage(srv *eventsource.Server, channel string) error {
	for {
		if currentImageIndex >= len(allImageUrls)-1 {
			currentImageIndex = 0
		}

		image, err := loadImage(currentImageIndex)
		if err != nil {
			return err
		}
		srv.Publish([]string{channel}, imageEvent{id: time.Now().Nanosecond(), image: image})
		time.Sleep(pushInterval)
		currentImageIndex += 1
	}
}

func loadImage(index int) (image, error) {
	img, ok := allImageUrls[index]
	if !ok {
		return image{}, errors.New(fmt.Sprintf("Image does not exist on index %d", index))
	}
	resp, err := http.Get(img.url)
	if err != nil {
		return image{}, errors.New("Failed to get image from " + img.url + ". " + err.Error())
	}
	defer resp.Body.Close()
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return image{}, errors.New("Failed to read image from response body. " + err.Error())
	}

	img.image = imageData
	return img, nil
}

// starts the listener
func startServer() {
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln("Serve failed!", err)
	}
}

// setPushInterval sets the interval to push next image based on the environment or the default value
func setPushInterval() {
	if os.Getenv(intervalEnvKey) != "" {
		interval, err := strconv.Atoi(os.Getenv(intervalEnvKey))
		if err != nil {
			if err != nil {
				log.Printf("Set interval is not valid. Using default... (%d)\n", defaultPushInterval)
			}
			return
		}
		pushInterval = time.Duration(interval) * time.Second
	}
}

func main() {
	var err error
	setPushInterval()

	srv := eventsource.NewServer()
	srv.Gzip = true
	defer srv.Close()

	//	l, err := net.Listen("tcp", ":"+port)
	//	if err != nil {
	//		log.Fatalln("Failed to start listener on port "+port, err)
	//	}
	//	defer l.Close()
	http.HandleFunc("/"+landing, srv.Handler(landing))
	http.Handle("/", http.FileServer(http.Dir("")))
	go startServer()

	err = loadImageUrls()
	if err != nil {
		log.Fatalln("Failed to load list with image urls!", err.Error())
	}

	err = pushImage(srv, landing)
	if err != nil {
		log.Fatalln("Error while pushing new image.", err)

	}
}
