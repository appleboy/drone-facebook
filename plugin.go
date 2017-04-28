package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/appleboy/drone-facebook/template"
	"github.com/paked/messenger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"
)

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Tag      string
		Event    string
		Number   int
		Commit   string
		Message  string
		Branch   string
		Author   string
		Email    string
		Status   string
		Link     string
		Started  float64
		Finished float64
	}

	// Config for the plugin.
	Config struct {
		PageToken   string
		VerifyToken string
		Verify      bool
		MatchEmail  bool
		To          []string
		Message     []string
		Image       []string
		Audio       []string
		Video       []string
		File        []string
		Port        int
		AutoTLS     bool
		Host        []string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

var (
	// ReceiveCount is receive notification count
	ReceiveCount int64
	// SendCount is send notification count
	SendCount int64
)

func init() {
	// Support metrics
	m := NewMetrics()
	prometheus.MustRegister(m)
}

func trimElement(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.Trim(value, " ")
		if len(value) == 0 {
			continue
		}
		newKeys = append(newKeys, value)
	}

	return newKeys
}

func parseTo(to []string, authorEmail string, matchEmail bool) []int64 {
	var emails []int64
	var ids []int64
	attachEmail := true

	for _, value := range trimElement(to) {
		idArray := trimElement(strings.Split(value, ":"))

		// check id
		id, err := strconv.ParseInt(idArray[0], 10, 64)
		if err != nil {
			continue
		}

		// check match author email
		if len(idArray) > 1 {
			if email := idArray[1]; email != authorEmail {
				continue
			}

			emails = append(emails, id)
			attachEmail = false
			continue
		}

		ids = append(ids, id)
	}

	if matchEmail == true && attachEmail == false {
		return emails
	}

	for _, value := range emails {
		ids = append(ids, value)
	}

	return ids
}

// Handler is http handler.
func (p Plugin) Handler(client *messenger.Messenger) http.Handler {
	// Setup a handler to be triggered when a message is received
	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		p, err := client.ProfileByID(m.Sender.ID)
		if err != nil {
			fmt.Println("Something went wrong!", err)
		}

		ReceiveCount++
		r.Text(fmt.Sprintf("Hello, %v!", p.FirstName))
	})

	// Setup a handler to be triggered when a message is delivered
	client.HandleDelivery(func(d messenger.Delivery, r *messenger.Response) {
		SendCount++
		fmt.Println("Delivered at:", d.Watermark().Format(time.UnixDate))
	})

	// Setup a handler to be triggered when a message is read
	client.HandleRead(func(m messenger.Read, r *messenger.Response) {
		fmt.Println("Read at:", m.Watermark().Format(time.UnixDate))
	})

	return client.Handler()
}

// Webhook support line callback service.
func (p Plugin) Webhook() error {
	client, err := p.Bot()
	if err != nil {
		return err
	}

	mux := p.Handler(client)

	if p.Config.Port != 443 && !p.Config.AutoTLS {
		log.Println("Line Webhook Server Listin on " + strconv.Itoa(p.Config.Port) + " port")
		if err := http.ListenAndServe(":"+strconv.Itoa(p.Config.Port), mux); err != nil {
			log.Fatal(err)
		}
	}

	if p.Config.AutoTLS && len(p.Config.Host) != 0 {
		log.Println("Line Webhook Server Listin on 443 port, hostname: " + strings.Join(p.Config.Host, ", "))
		return http.Serve(autocert.NewListener(p.Config.Host...), mux)
	}

	return nil
}

func (p Plugin) serveMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		promhttp.Handler().ServeHTTP(w, req)
	})

	// Setup HTTP Server for receiving requests from LINE platform
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Welcome to facebook webhook page.")
	})

	return mux
}

// Bot is new Line Bot clien.
func (p Plugin) Bot() (*messenger.Messenger, error) {
	if len(p.Config.PageToken) == 0 || len(p.Config.VerifyToken) == 0 {
		log.Println("missing facebook config")

		return nil, errors.New("missing facebook config")
	}

	return messenger.New(messenger.Options{
		Verify:      p.Config.Verify,
		Token:       p.Config.PageToken,
		VerifyToken: p.Config.VerifyToken,
		WebhookURL:  "callback",
		Mux:         p.serveMux(),
	}), nil
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	client, err := p.Bot()
	if err != nil {
		return err
	}

	var message []string
	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	ids := parseTo(p.Config.To, p.Build.Email, p.Config.MatchEmail)

	// send message.
	for _, user := range ids {
		To := messenger.Recipient{
			ID: user,
		}

		// send text notification
		for _, value := range trimElement(message) {
			txt, err := template.RenderTrim(value, p)
			if err != nil {
				return err
			}

			client.Send(To, txt)
		}

		// send image notification
		for _, value := range trimElement(p.Config.Image) {
			client.Attachment(To, messenger.ImageAttachment, value)
		}

		// send audio notification
		for _, value := range trimElement(p.Config.Audio) {
			client.Attachment(To, messenger.AudioAttachment, value)
		}

		// send video notification
		for _, value := range trimElement(p.Config.Video) {
			client.Attachment(To, messenger.VideoAttachment, value)
		}

		// send file notification
		for _, value := range trimElement(p.Config.File) {
			client.Attachment(To, messenger.FileAttachment, value)
		}
	}

	return nil
}

// Message is plugin default message.
func (p Plugin) Message(repo Repo, build Build) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		build.Status,
		build.Link,
		build.Branch,
		build.Message,
		build.Author,
	)}
}
