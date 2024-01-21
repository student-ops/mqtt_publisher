package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	fmt.Fprintf(os.Stderr, "*** 開始 ***\n")

	nodeIP := os.Getenv("NODE_IP") // 環境変数からNodeのIPを取得
	if nodeIP == "" {
		log.Fatal("NODE_IP 環境変数が設定されていません")
	}

	mqttBroker := fmt.Sprintf("tcp://%s:1883", nodeIP) // MQTTブローカーのアドレス

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
	cc := mqtt.NewClient(opts)

	if token := cc.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}
	type Surroundings struct {
		Timestamp   time.Time `json:"timestamp"`
		Rssi        float64   `json:"rssi"`
		Tempreture  float64   `json:"tempreture"`
		Moisuture   float64   `json:"moisuture"`
		AirPressure float64   `json:"airPressure"`
	}

	type RequestPayload struct {
		Surroundings []Surroundings `json:"surroundings,omitempty"`
	}
	tempreture := 25.8
	moisture := 60.3
	airPressure := 1024.5

	// var msg RequestPayload
	for it := 0; it < 1000; it++ {

		go func(it int) {
			tempreture += rand.Float64()*0.5 - 0.25 // 約1%の変動
			moisture += rand.Float64()*0.5 - 0.25
			airPressure += rand.Float64()*5 - 2.5

			payload := Surroundings{
				Timestamp:   time.Now(),
				Rssi:        float64(it),
				Tempreture:  tempreture,
				Moisuture:   moisture,
				AirPressure: airPressure,
			}
			fmt.Println(it)

			jsonMsg, err := json.Marshal(payload)
			if err != nil {
				log.Fatal(err)
			}
			token := cc.Publish("paper_wifi/test", 0, false, jsonMsg)
			token.Wait()
		}(it)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(3 * time.Second)

	cc.Disconnect(250)

	fmt.Println("Complete publish")
	fmt.Fprintf(os.Stderr, "*** 終了 ***\n")
}
