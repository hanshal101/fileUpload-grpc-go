package routes

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	pb "github.com/hanshal101/fileUpload/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client pb.UploadClient

func UploadPOST(res *gin.Context) {
	conn, err := grpc.Dial("localhost:9876", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	file, err := res.FormFile("file")
	if err != nil {
		res.JSON(403, gin.H{"message": "Error in File Upload"})
		return
	}
	fileName := file.Filename

	fileData, err := file.Open()
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer fileData.Close()

	Client = pb.NewUploadClient(conn)

	buf := make([]byte, 1024*1024*2)
	stream, err := Client.FileUpload(context.TODO())
	if err != nil {
		log.Fatalf("Error in Stream Client")
	}

	for {
		var num int
		num, err = fileData.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		chunk := buf[:num]
		stream.Send(&pb.UploadRequest{Chunks: chunk, FileName: fileName})
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}
	res.JSON(201, gin.H{"message": "File Uploaded Succesfully"})
}
