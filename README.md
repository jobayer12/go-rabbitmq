# Setup RabbitMQ locally Using Docker 

## Run RabbitMQ docker image in background
```shell
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management
```

## Add new RabbitMQ user
```shell
docker exec rabbitmq rabbitmqctl add_user [username] [password]
```
```shell
Example: docker exec rabbitmq rabbitmqctl add_user jobayer jobayer
```

## Set RabbitMq user permission
```shell
docker exec rabbitmq rabbitmqctl set_user_tags [username] [permission]
```
```shell
Example: docker exec rabbitmq rabbitmqctl set_user_tags jobayer administrator
```

## Delete rabbitmq user
```shell
docker exec rabbitmq rabbitmqctl delete_user [username]
```

## Add virtual host
```shell
docker exec rabbitmq rabbitmqctl add_vhost [vhost_name]
```

## Create Binding
```shell
docker exec rabbitmq rabbitmqadmin declare exchange --vhost=[virtual host name] name=[binding name] type=topic -u [username] -p [password] durable=true
```
## Add binding permission
```shell
docker exec rabbitmq rabbitmqctl set_topic_permissions -p [vhost] [username] [bindingname] "^customers.*" "^customers.*"
```
