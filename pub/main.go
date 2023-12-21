package main

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/nats-io/stan.go"
	"golang.org/x/oauth2"
)

const clusterID = "test-cluster"
const clientID = "producer-client"
const subject = "zxc"

func main() {
	configURL := "http://host.docker.internal:9080/auth/realms/zxc-realm"
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, configURL)
	if err != nil {
		log.Fatalf("failed to get provider: %v", err)
	}

	clientID := "confidential-client"
	clientSecret := "v5j235h5v25236j26g2h3v3"
	redirectURL := "http://localhost:8081/demo/callback"

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	state := "somestate"

	oidcConfig := &oidc.Config{ClientID: clientID}
	verifier := provider.Verifier(oidcConfig)

	var tmpl = template.Must(template.New("order").Parse(htmlTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Проверка наличия куки с ID токеном
		if _, err := r.Cookie("id_token"); err == nil {
			http.Redirect(w, r, "/hello", http.StatusFound)
			return
		}
		http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("id_token")
		if err != nil {
			log.Println("No ID token in cookie")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		_, err = verifier.Verify(ctx, tokenCookie.Value)
		if err != nil {
			log.Println("Failed to verify ID Token:", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form", http.StatusInternalServerError)
				return
			}
			id := r.FormValue("id")

			order, err := loadOrderAndUpdateUID(id) // Загрузка и обновление заказа
			if err != nil {
				http.Error(w, "Error loading order", http.StatusInternalServerError)
				return
			}
			jsonData, err := json.Marshal(order) // Сериализация обновлённого заказа в JSON
			if err != nil {
				http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
				return
			}

			err = publishToNATS(jsonData) // Отправка данных в NATS
			if err != nil {
				http.Error(w, "Error publishing to NATS", http.StatusInternalServerError)
				return
			}

			log.Println("sended message with id:", order.OrderUID)
		}

		// Вывод страницы при GET-запросе
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/demo/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		idToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "id_token",
			Value: idToken,
			Path:  "/",
		})

		http.Redirect(w, r, "/hello", http.StatusFound)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func loadOrderAndUpdateUID(orderUID string) (Order, error) {
	data, err := ioutil.ReadFile("/ord.json") // Убедитесь, что путь к файлу указан верно
	if err != nil {
		return Order{}, err
	}

	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		return Order{}, err
	}

	order.OrderUID = orderUID      // Обновление OrderUID согласно введённому пользователем значению
	order.DateCreated = time.Now() // Обновление даты создания заказа

	return order, nil
}

func publishToNATS(data []byte) error {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://nats-streaming:4222"))
	if err != nil {
		return err
	}
	defer sc.Close()

	err = sc.Publish(subject, data)
	if err != nil {
		return err
	}
	return nil
}
