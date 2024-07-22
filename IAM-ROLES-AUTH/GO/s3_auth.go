package main

import (
	"context"
	"fmt"
	"io"
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
		CreateBucket(S3Client("us-east-1"), os.Args[2], os.Args[3])
	case "--delbkt":
		DeleteBucket(S3Client("us-east-1"), os.Args[2], os.Args[3])
	case "--del":
		DeleteObject(S3Client("us-east-1"), os.Args[2], os.Args[3])
	case "--listbkts":
		ListBuckets(S3Client("us-east-1"))
	case "--list":
		ListObjects(S3Client("us-east-1"), os.Args[2], os.Args[3])
	case "--upload":
		UploadObject(S3Client("us-east-1"), os.Args[2], os.Args[3], os.Args[4])
	case "--download":
		DownloadObject(S3Client("us-east-1"), os.Args[2], os.Args[3], os.Args[4])
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
	_, err := s3_client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})

	if err != nil {
		fmt.Printf("Falha ao deletar o objeto %s no Bucket %s: %s\n", objectName, bucketName, err)
	} else {
		fmt.Printf("O objeto %s foi deletado com sucesso no Bucket %s\n", objectName, bucketName)
	}
}

func UploadObject(s3_client *s3.Client, bucketName string, objectPath string, fileName string) {
	fileContent, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
	}

	_, err = s3_client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectPath + fileName),
		Body:   fileContent,
	})

	if err != nil {
		fmt.Printf("Falha ao carregar o objeto %s no Bucket %s: %s\n", objectPath+fileName, bucketName, err)
	} else {
		fmt.Printf("O objeto %s foi carregado com sucesso no Bucket %s\n", objectPath+fileName, bucketName)
	}
}

func DownloadObject(s3_client *s3.Client, bucketName string, objectPath string, fileName string) error {
	ctx := context.TODO()
	objInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectPath + fileName),
	}

	resp, err := s3_client.GetObject(ctx, objInput)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("O arquivo %s foi baixado com sucesso!\n", objectPath+fileName)

	return nil
}
