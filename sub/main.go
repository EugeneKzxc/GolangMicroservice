package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {

	orderCh := make(chan Order)

	var tmpl = template.Must(template.New("order").Parse(htmlTamplate))

	conn, err := connectToDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer conn.Close(context.Background())

	orderMap, err := getOrdersFromDB(conn)
	if err != nil {
		fmt.Println("Failed to select data:", err)
		return
	}

	sc, sub, err := connectToNatsStreaming(orderCh)
	if err != nil {
		fmt.Println("STAN connextion error:", err)
		return
	}
	defer sub.Unsubscribe()
	defer sc.Close()

	// Горутина, обрабатывающая заказы, приходящие через канал.
	go func() {
		for order := range orderCh {
			exists, err := CheckRecordExists(conn, order.OrderUID)
			if err != nil {
				log.Printf("Error on checking record`s existanse: %v", err)
				continue
			} else if !exists {
				err := insertOrderToDB(conn, order)
				if err != nil {
					log.Printf("Error on inserting data to DB: %v", err)
					continue
				}
				// Обновление локальной карты заказов.
				orderMap[order.OrderUID] = order
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}
		id := r.FormValue("id")

		order, found := orderMap[id]
		if !found {
			tmpl.Execute(w, nil)
			return
		}
		tmpl.Execute(w, order)
	})

	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Server started at http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
}
