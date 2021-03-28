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


cluster
````
mac config:
one node:
sudo ifconfig lo0 alias 1.1.1.1
sudo ifconfig lo0 alias 1.1.1.2
sudo ifconfig lo0 alias 1.1.1.3


 ./distributed-cache -node 1.1.1.1
 ./distributed-cache -node 1.1.1.2 -cluster 1.1.1.1
 ./distributed-cache -node 1.1.1.3 -cluster 1.1.1.1

 ./client -h 1.1.1.1 -c set -k keya - v a
 ./client -h 1.1.1.2 -c set -k keyb - v b
 ./client -h 1.1.1.3 -c set -k keyc - v c
````