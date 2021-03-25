Implement a distributed cache


start:
```
go build

./distributed-cache
```




http
```
curl -v 127.0.0.1:8080/PUT  -XPUT  -d  '{"key":"testkey", "value": "1233"}'
curl 127.0.0.1:8080/INFO 
curl '127.0.0.1:8080/GET?key=testkey'
```


tcp 
```
cd ./client/run && go build

./run -h localhost -c get -k "k1"
./run -h localhost -c set  -k "k1" -v "v1"
./run -h localhost -c del  -k "k1"
```