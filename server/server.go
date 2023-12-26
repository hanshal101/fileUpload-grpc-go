package main

import (
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	pb "github.com/hanshal101/fileUpload/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedUploadServer
}

func main() {
	lis, err := net.Listen("tcp", ":9876")
	if err != nil {
		log.Fatalf("Error in listen")
	}

	c := grpc.NewServer()
	pb.RegisterUploadServer(c, &Server{})
	reflection.Register(c)

	if err := c.Serve(lis); err != nil {
		log.Fatalf("Error in Serve")
	}
}

func (s *Server) FileUpload(stream pb.Upload_FileUploadServer) error {
	filesize := 0
	var fileBytes []byte
	var fileName string
	// firstMessage, err := stream.Recv()
	// if err != nil {
	// 	log.Fatalf("Error receiving file name: %v", err)
	// }
	// fileName = firstMessage.GetFileName()
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			// fileName = request.GetFileName()
			break
		}
		if err != nil {
			log.Fatalf("Error in file server upload")
		}

		if fileName == "" && request.GetFileName() != "" {
			fileName = request.GetFileName()
		}
		chunks := request.GetChunks()
		fileBytes = append(fileBytes, chunks...)
		filesize += int(len(chunks))
	}
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		err := os.Mkdir("uploads", os.ModePerm)
		if err != nil {
			return err
		}
	}

	timeStamp := time.Now()
	fileploadName := timeStamp.Local().Format("2006-01-02-15-04-05") + fileName
	dstPath := filepath.Join("./uploads/", fileploadName)
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err2 := f.Write(fileBytes)

	if err2 != nil {
		return err2
	}

	return stream.SendAndClose(&pb.UploadResponse{FileName: fileName, FileSize: strconv.Itoa(filesize)})
}
