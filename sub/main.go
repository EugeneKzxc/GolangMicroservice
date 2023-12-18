package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
}
