## GG-Seq

A distributed id generation system implemented in go language.

## Clone code 

```shell
git clone https://github.com/gouez/gg-seq.git
```


## run server

```shell
cd server && go run .

```

## rest api

```shell
curl 127.0.0.1:8000/get?bizType=test
curl 127.0.0.1:8000/get/batch?bizType=test&size=100
```