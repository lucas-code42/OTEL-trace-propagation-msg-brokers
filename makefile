run:
	docker compose up --build -d &&\
	echo "waiting for rabbitmq server up..." &&\
	sleep 15 &&\
	go run main.go

