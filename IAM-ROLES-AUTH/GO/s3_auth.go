package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	switch os.Args[1] {
	case "--help":
		DocHelp()
	case "--create":
		CreateBucket(S3Client(), os.Args[2], os.Args[3])
	case "--delbkt":
		DeleteBucket(S3Client(), os.Args[2], os.Args[3])
	case "--del":
		DeleteObject(S3Client(), os.Args[2], os.Args[3])
	case "--listbkts":
		ListBuckets(S3Client())
	case "--list":
		// awsListObjects(os.Args[2])
	case "--upload":
		// awsUploadObject(os.Args[2], os.Args[3])
	case "--download":
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

func CreateBucket(s3_client *s3.Client, bucketName string, bucketRegion string) {
	_, err := s3_client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(bucketRegion),
		},
	})

	if err != nil {
		fmt.Println("Falha ao criar o bucket: ", err)
	} else {
		fmt.Println("Bucket", bucketName, "criado com sucesso!")
	}

}

func DeleteBucket(s3_client *s3.Client, bucketName string, bucketRegion string) {
	_, err := s3_client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		fmt.Println("Falha ao deletar o bucket: ", err)
	} else {
		fmt.Println("Bucket", bucketName, "deletado com sucesso!")
	}
}

func DeleteObject(s3_client *s3.Client, bucketName string, objectName string) {
	fmt.Println("Em desenvolvimento...")
}
