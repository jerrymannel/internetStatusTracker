package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/tatsushid/go-fastping"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func check(_err error) {
	if _err != nil {
		fmt.Println(_err)
		os.Exit(1)
	}
}

type internetStatus struct {
	Time      time.Time `json:"time"`
	Active    bool      `json:"active"`
	Bandwidth string    `json:"bandwidth"`
	Host      string    `json:"host"`
	UpSince   string    `json:"upSinse"`
}

func main() {
	t := time.Now()

	mongoURL := flag.String("d", "localhost:27017", "MongoDB URL")
	ipToPing := flag.String("i", "1.1.1.1", "IP to ping")

	flag.Parse()

	clientOptions := options.Client().ApplyURI(*mongoURL)

	client, err := mongo.Connect(context.Background(), clientOptions)
	check(err)

	err = client.Ping(context.Background(), nil)
	check(err)

	collection := client.Database("18294").Collection("internetStatus")

	hostname, err := os.Hostname()
	check(err)

	upSince, err := exec.Command("uptime", "-s").Output()
	check(err)

	data := internetStatus{t, false, "NA", hostname, string(upSince)}

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", *ipToPing)
	check(err)

	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		data.Active = true
	}

	p.OnIdle = func() {
		fmt.Println(*ipToPing, t.Format(time.RFC1123), data.Active)
		_, err = collection.InsertOne(context.Background(), data)
		check(err)
	}

	err = p.Run()
	check(err)
}
