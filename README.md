# Draft

```

docker run --rm --name my-redis -p 6379 redis:5.0-rc
docker run --rm -e REDIS=my-redis:6379 --link=my-redis redrabbit:0


GET /info
GET /{scope}/queues
GET /{scope}/queues/{qid}/count
GET /{scope}/queues/{qid}?page=0

POST /{scope}/queues/ -> qid
POST /{scope}/queues/{qid} -> jid
GET /{scope}/queues/{qid}/{jid}

POST /{scope}/queues/{qid}/{jid}?ack
POST /{scope}/queues/{qid}/{jid}?nack

```
