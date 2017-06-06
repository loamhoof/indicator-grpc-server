package main

import (
	"flag"
	"log"

	pb "github.com/loamhoof/indicator"
	"github.com/loamhoof/indicator/client"
)

var (
	port            int
	id, label, icon string
	active          bool
)

func init() {
	flag.IntVar(&port, "port", 15000, "Grpc server port")

	flag.StringVar(&id, "id", "", "Indicator's ID")
	flag.StringVar(&label, "label", "", "Indicator's label")
	flag.StringVar(&icon, "icon", "", "Indicator's icon")
	flag.BoolVar(&active, "active", true, "Is indicator active?")

	flag.Parse()
}

func main() {
	c, err := client.NewShepherdClient(port)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer c.Close()

	req := &pb.Request{
		Id:     id,
		Label:  label,
		Icon:   icon,
		Active: active,
	}

	if _, err := c.Update(req); err != nil {
		log.Println(err)
	}
}
