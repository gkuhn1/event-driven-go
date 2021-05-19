
RUNNING
=======

docker-compose up -d kafka zookeeper

docker-compose up members_service orders_service payments_service shipping_service

APIs
-----------

Create member
```
curl --location --request POST 'localhost:8580/members' \
    --form 'name="test"'
```


Create order
```
curl --location --request POST 'localhost:8581/orders'
```