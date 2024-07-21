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
		CreateBucket(S3Client(os.Args[3]), os.Args[2], os.Args[3])
	case "--delbkt":
		DeleteBucket(S3Client(os.Args[3]), os.Args[2], os.Args[3])
	case "--del":
		DeleteObject(S3Client(os.Args[3]), os.Args[2], os.Args[3])
	case "--listbkts":
		ListBuckets(S3Client(os.Args[2]))
	case "--list":
		ListObjects(S3Client("us-east-1"), os.Args[2], os.Args[3])
	case "--upload":
		// awsUploadObject(os.Args[2], os.Args[3])
	case "--download":
		// awsDownloadObject(os.Args[2], os.Args[3])
	default:
		fmt.Println("Parâmetro Inválido!")
		DocHelp()
	}

}

func SdkAuth(bucketRegion string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(bucketRegion))

	if err != nil {
		fmt.Println("Erro ao tentar autenticar:", err)
	}

	return cfg
}

func S3Client(bucketRegion string) *s3.Client {
	s3_client := s3.NewFromConfig(SdkAuth(bucketRegion))

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
		--del / --upload / --download --> go s3_no_auth.go --parâmetro nome_do_bucket path_com_nome_do_arquivo.txt região
		--list --> go s3_no_auth.go --parâmetro nome_do_bucket path_diretorio região
		--listbkt --> s3_no_auth.go --parâmetro região
	`

	fmt.Println(doc)
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

func ListBuckets(s3_client *s3.Client) {
	resposta, err := s3_client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		fmt.Println("Falha ao listar os buckets", err)
	}

	for _, bucket := range resposta.Buckets {
		fmt.Println(*bucket.Name)
	}
}

func ListObjects(s3_client *s3.Client, bucketName string, prefixName string) {
	resposta, err := s3_client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefixName),
	})

	if err != nil {
		fmt.Println("Falha ao listar os objetos em ", bucketName+"/"+prefixName)
		fmt.Println(err)
	} else {
		for _, objeto := range resposta.Contents {
			fmt.Println(*objeto.Key)
		}
	}

}

func DeleteObject(s3_client *s3.Client, bucketName string, objectName string) {
	fmt.Println("Em desenvolvimento...")
}
