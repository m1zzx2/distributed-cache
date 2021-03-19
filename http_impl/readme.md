Implement a distributed cache


####start:
```
go build

./distributed-cache
```




#curl -v 127.0.0.1:8080/PUT  -XPUT  -d  '{"key":"testkey", "value": "1233"}'
#curl 127.0.0.1:8080/INFO 
#curl '127.0.0.1:8080/GET?key=testkey'