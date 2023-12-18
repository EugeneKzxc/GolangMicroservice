package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/stan.go"
)

const clusterID = "test-cluster"
const clientID = "producer-client"
const subject = "zxc"

func main() {
	var tmpl = template.Must(template.New("order").Parse(htmlTemplate))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl.Execute(w, nil) // Отобразить форму, если метод не POST

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
	})

	log.Println("Server started at http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
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
