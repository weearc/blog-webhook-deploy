package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

func reLaunch() {
	cmd := exec.Command("load")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
}
func index(w http.ResponseWriter, r *http.Request) {

	signature := r.Header.Get("X-Hub-Signature")
	if len(signature) <= 0 {
		return
	}
	payload, _ := ioutil.ReadAll(r.Body)
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write(payload)
	expectedMac := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature[5:]), []byte(expectedMac)) {
		io.WriteString(w, "<h1>401 Signature is error!</h1>")
		return
	}
	io.WriteString(w, "<h1>200 Deploy server is running!</h1>")
	reLaunch()
}
var port int
var secret string
func main() {


	app := cli.NewApp()
	app.Name = "Deploy tool"
	app.Usage = "For set up hexo blog on web server"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "port to listen to",
			Required:    true,
			Value:       44477,
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "key, k",
			Usage:       "your key",
			Required:    true,
			Value:       "",
			Destination: &secret,
		},
	}
	app.Action = func(c *cli.Context) error {

		return nil
	}
	app.After = func(c *cli.Context) error {
		portForward := ":" + strconv.Itoa(port)
		fmt.Printf(secret + "\n")
		fmt.Printf(portForward + "\n")
		fmt.Printf("Listening to %d ...", port)
		http.HandleFunc("/", index)
		http.ListenAndServe(portForward, nil)
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
