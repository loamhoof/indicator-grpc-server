package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/conformal/gotk3/gtk"
	"google.golang.org/grpc"

	pb "github.com/loamhoof/indicator"
)

var (
	port             int
	logFile, iconDir string
)

func init() {
	flag.IntVar(&port, "port", 15000, "Bind port")
	flag.StringVar(&logFile, "log", "", "Log file")
	flag.StringVar(&iconDir, "icon", "", "Icons directory")

	flag.Parse()
}

func main() {
	go serveGtk()

	serveGrpc()
}

func serveGtk() {
	gtk.Init(nil)
	gtk.Main()
}

func serveGrpc() {
	var logger *log.Logger
	if logFile == "" {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		f, err := os.OpenFile(logFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			logger.Fatalln(err)
		}
		defer f.Close()

		logger = log.New(f, "", log.LstdFlags)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShepherdServer(grpcServer, newIndicatorShepherd(iconDir, logger))
	grpcServer.Serve(lis)
}
