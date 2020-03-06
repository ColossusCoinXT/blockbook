## To run application

1. Update following details inside ./build/blockchaincfg.json:
- `"rpc_url"`: _"RPC_URL"_
- `"rpc_user"`: _"RPC_USER_NAME"_
- `"rpc_pass"`: _"RPC_PASSWORD"_
- `"message_queue_binding"`: _"MESSAGE_QUEUE_BINDING"_

2. From `blockbook` project directory run the following command: `docker-compose up -d`

## To get all the running containers on the machine
- `docker ps -a`

## To check logs for a container
- `docker logs <container-name>`

## To stop the application
- From `blockbook` project directory run the following command: `docker-compose down`
