package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ronenniv/s3dir/cli"
	"github.com/ronenniv/s3dir/s3client"
)

func initClient() *s3client.S3Client {
	client := s3client.New()
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {
	fc, err := cli.CheckArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", fc)

	// TODO: execute the command based on the arguments (ls/cd)

	client := initClient()

	buckets, _ := client.ListBuckets()
	buckets.PrintShort(os.Stdout)
	buckets.PrintLong(os.Stdout)

	obj, err := client.ListObjects("tess-qa-checks")
	if err != nil {
		log.Fatal(err)
	}
	obj.PrintShort(os.Stdout)
	obj.PrintLong(os.Stdout)
	obj, err = client.ListObjects("h1-reports-for-slack")
	if err != nil {
		log.Fatal(err)
	}
	obj.PrintShort(os.Stdout)
	obj.PrintLong(os.Stdout)

}
