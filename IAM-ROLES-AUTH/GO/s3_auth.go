package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	msg := "Número inválido de parâmetros!"

	switch os.Args[1] {
	case "--help":
		DocHelp()
	case "--create":
		if len(os.Args) != 3 {
			fmt.Println(msg)
		}

		AwsCreateBucket(S3Client(), os.Args[2], os.Args[3])
	case "--delbkt":
		if len(os.Args) != 3 {
			fmt.Println(msg)
		}
		// awsDeleteBucket(os.Args[2], os.Args[3])
	case "--del":
		if len(os.Args) != 3 {
			fmt.Println(msg)
		}
		// awsDeleteObject(os.Args[2], os.Args[3])
	case "--listbkts":
		ListBuckets(S3Client())
	case "--list":
		if len(os.Args) == 2 {
			// awsListObjects(os.Args[2])
		} else {
			fmt.Println(msg)
		}
	case "--upload":
		if len(os.Args) != 2 {
			fmt.Println(msg)
		}
		// awsUploadObject(os.Args[2], os.Args[3])
	case "--download":
		if len(os.Args) != 3 {
			fmt.Println(msg)
		}
		// awsDownloadObject(os.Args[2], os.Args[3])
	default:
		fmt.Println("Parâmetro Inválido!")
		DocHelp()
	}
}

func SdkAuth() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		fmt.Println("Erro ao autenticar", err)
	}

	return cfg
}

func S3Client() *s3.Client {
	s3_client := s3.NewFromConfig(SdkAuth())

	return s3_client
}

func DocHelp() {
	doc := `	Dependências:
		go get github.com/aws/aws-sdk-go-v2
		go get github.com/aws/aws-sdk-go-v2/config
		go get github.com/aws/aws-sdk-go-v2/service/s3
	Parâmetros:
		--create    --> Cria o Bucket com o nome e região especificados
		--delbkt    --> Deleta um Bucket vazio
		--del       --> Deleta objetos no Bucket
		--listbkts  --> Lista todos os Buckets
		--list      --> Lista os objetos do Buckets
		--upload    --> Faz upload de objetos para o Bucket
		--download  --> Faz o download de objetos do Bucket

	Sintaxe:
		--create / --delbkt --> go s3_no_auth.go --parâmetro nome_do_bucket região
		--del / --upload / --download --> go s3_no_auth.go --parâmetro nome_do_bucket path_com_nome_do_arquivo.txt
		--list --> go s3_no_auth.go --parâmetro nome_do_bucket path_diretorio(opcional)
		--listbkt --> s3_no_auth.go --parâmetro
	`

	fmt.Println(doc)
}

func ListBuckets(s3_client *s3.Client) {
	resposta, err := s3_client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		fmt.Println("Falha ao listar os buckets", err)
	}

	for _, bucket := range resposta.Buckets {
		fmt.Println(*bucket.Name)
	}
}

func AwsCreateBucket(s3_client *s3.Client, bucketName string, bucketRegion string) {
	resposta, err := s3_client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		// CreateBucketConfiguration: &s3.CreateBucketConfiguration{
		// 	LocationConstraint: s3.BucketLocationConstraint(bucketRegion),
		// },
	})

	if err != nil {
		fmt.Println("Falha ao criar o bucket: ", err)
	}

	fmt.Println(resposta)
}

// func AwsDeleteBucket(bucketName, bucketRegion string) {
//
// }

// func AwsDeleteObject(bucketName, objectName string) {
//
// }
