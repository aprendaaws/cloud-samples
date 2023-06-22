# Feito com Boto 1.21.21 e Python 3.6.9

import boto3
from sys import argv


# Client para interação
s3_client = boto3.client("s3")


def banner(nums=15):
    mensagem = "-=" * nums + "-"
    print(mensagem)


def aws_create_bucket(bucket_name, bucket_region):
    response_create = s3_client.create_bucket(
        Bucket=bucket_name,
        CreateBucketConfiguration={
            'LocationConstraint': bucket_region,
        },
    )


    response_http_status_code = response_create["ResponseMetadata"]["HTTPStatusCode"]
    if response_http_status_code == 200:
        banner()
        print(f"Bucket {bucket_name} criado com sucesso em {bucket_region}")
    else:
        banner()
        print(f"Falha ao criar o Bucket {bucket_name}")


def aws_delete_bucket(bucket_name, bucket_region):
    response_delete = s3_client.delete_bucket(
        Bucket=bucket_name
    )

    response_http_status_code = response_delete["ResponseMetadata"]["HTTPStatusCode"]
    if response_http_status_code == 204:
        banner()
        print(f"Bucket {bucket_name} deletado com sucesso em {bucket_region}")
    else:
        banner()
        print(f"Falha ao deletar o Bucket {bucket_name}")
    

def aws_delete_object(bucket_name, object_name):
    response_delete_object = s3_client.delete_object(
        Bucket=bucket_name,
        Key=object_name
    )

    response_http_status_code = response_delete_object["ResponseMetadata"]["HTTPStatusCode"]
    if response_http_status_code == 204:
        banner()
        print(f"O objeto {object_name} foi deletado com sucesso no Bucket {bucket_name}")
    else:
        banner()
        print(f"Falha ao deletar o objeto {object_name} no Bucket {bucket_name}")


def aws_list_objects(bucket_name, prefix_name=""):
    response_list_objects = s3_client.list_objects(
        Bucket=bucket_name,
        Prefix=prefix_name
    )

    response_http_status_code = response_list_objects["ResponseMetadata"]["HTTPStatusCode"]

    mensagem_de_listagem = "Listando os Objetos"
    len_mensagem_de_listagem =len(mensagem_de_listagem)

    if response_http_status_code == 200:
        banner(len_mensagem_de_listagem)
        print(f"++++++++  {mensagem_de_listagem} ++++++++")
        banner(len_mensagem_de_listagem)

        objetos = response_list_objects["Contents"]
        num_objetos = len(objetos)
        for indice in range(num_objetos):
            print(objetos[indice]["Key"])
        
    else:
        print(f"Falha ao listar o Bucket {bucket_name}")


def aws_upload_object(bucket_name, object_name):
    # Extraindo o nome do arquivo
    file_to_upload = object_name.split("/")[-1]

    # Lendo o arquivo para enviar como objeto
    with open(file_to_upload, "rb") as object_data:
        file_content = object_data.read()
        object_data.close()

    response_upload_object = s3_client.put_object(
        Bucket=bucket_name,
        Key=object_name,
        Body=file_content
    )

    response_http_status_code = response_upload_object["ResponseMetadata"]["HTTPStatusCode"]
    if response_http_status_code == 200:
        banner()
        print(f"O arquivo {file_to_upload} foi enviado com sucesso!")
    else:
        banner()
        print(f"Falha ao enviar o arquivo {file_to_upload}!")


def aws_download_object(bucket_name, object_name):
    banner()
    file_name = object_name.split("/")[-1]
    
    # Lendo o arquivo para enviar como objeto
    with open(file_name, "wb") as object_data:
        try:
            s3_client.download_fileobj(bucket_name, object_name, object_data)
            print(f"Arquivo {file_name} baixado com sucesso!")
        except:
            print(f"Falha ao fazer download do arquivo {file_name}")
        finally:
            object_data.close()


def main():
    """
    Pré Requisitos:
        Dependências:
            * Instalar as dependênciais do requirements.txt
              --> pip install -r requirements.txt
        
        Definir essas variáveis de ambiente:
            * RoleArn="arn:aws:iam::123456:role/NOME_DA_ROLE"
            * RoleSessionName="Qualquer nome de identificação"
            * ExternalId="NOME_DO_ID"
        
        RoleArn: Identificador único de um recurso na AWS

        RoleSessionName: Um nome qualquer para identificar
        a sessão que está sendo criada.

        ExternalID: É definido quando criamos a role
        e server como uma camada extra de segurança
        para poder validar a autenticação.

    Parâmetros:
        --create    --> Cria o Bucket com o nome e região especificados
        --delbkt    --> Deleta um Bucket vazio
        --del       --> Deleta objetos no Bucket
        --list      --> Lista os objetos do Buckets
        --upload    --> Faz upload de objetos para o Bucket
        --download  --> Faz o download de objetos do Bucket
        
    Sintaxe:
        --create/--delbkt
            --> python3 s3_auth.py --parâmetro nome_do_bucket região
        --del/--upload/--download
            --> python3 s3_auth.py --parâmetro nome_do_bucket path_com_nome_do_arquivo.txt
        --list
            --> python3 s3_auth.py --parâmetro nome_do_bucket path_diretorio(opcional)
    
    """

    # Parâmetros:
    if argv[1] == "--help":
        print(main.__doc__)
    elif argv[1] == "--create":
        aws_create_bucket(argv[2], argv[3])
    elif argv[1] == "--delbkt":
        aws_delete_bucket(argv[2], argv[3])
    elif argv[1] == "--del":
        aws_delete_object(argv[2], argv[3])
    elif argv[1] == "--list":
        if len(argv) == 3:
            aws_list_objects(argv[2])
        elif len(argv) == 4:
            aws_list_objects(argv[2], argv[3])
        else:
            print("Parâmetro Inválido!")
            print("--list espera dois argumentos")
    elif argv[1] == "--upload":
        aws_upload_object(argv[2], argv[3])
    elif argv[1] == "--download":
        aws_download_object(argv[2], argv[3])
    else:
        print("Parâmetro Inválido!")
        print(main.__doc__)

main()
