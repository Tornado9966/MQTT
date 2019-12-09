package main

import (
        mqtt "github.com/eclipse/paho.mqtt.golang"
	"encoding/json" 
        "crypto/sha1"
	"encoding/hex"
        "time"
        "fmt"
)

type payload struct {
	Time   time.Time `json:"time"`
	Wisdom string    `json:"wisdom"`
	Secret string    `json:"secret"`
	Team   string    `json:"team"`
}

func main() {
	var done bool = false
        var result [4]string

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:1883", "localhost"))
	client := mqtt.NewClient(opts)
	client.Connect().Wait()
	

	go func() {
		client.Subscribe("/test/inception", 0, func(client mqtt.Client, message mqtt.Message) {
			var data payload
			if err := json.Unmarshal(message.Payload(), &data); err != nil {
				return
			}

                        secret := hex.EncodeToString(sha1.New().Sum([]byte(data.Secret + "a,b - opimisti")))

			switch data.Wisdom[:1] {
	                    case "1":
		                result[0] = data.Wisdom[2:]
	                    case "2":
		                result[1] = data.Wisdom[2:]
	                    case "3":
		                result[2] = data.Wisdom[2:]
	                    case "4":
		                result[3] = data.Wisdom[2:]
	                }

			if result[0] != "" && 
                           result[1] != "" && 
                           result[2] != "" &&
                           result[3] != "" {
	                        data, _ := json.Marshal(&payload{
                                    Time: time.Now(),
                                    Wisdom: result[0] + result[1] + result[2] + result[3],
	                            Secret: secret,
	                            Team: "a,b - opimisti",
                                })
			        client.Publish("/test/result", 0, false, data)
				done = true
			}
		})
	}()
	for !done {}
        fmt.Printf("Done\n")
}
