## GG-Seq

A distributed id generation system implemented in go language.

## Clone code 

git clone https://github.com/gouez/gg-seq.git

## config db

path /app/configs/config.yaml

## run 

cd /app/cmd/server 
go run .

## rest api

```shell
curl 127.0.0.1:8000/v1/getId/{bizType}
```