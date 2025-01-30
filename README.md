# Go Kafka Project

This project demonstrates a microservices architecture where two services communicate via Apache Kafka. The project also integrates Grafana, Prometheus, and Kafka Exporter for monitoring, and includes an email server using SMTP with MailDev for intercepting and viewing emails.

## Overview

The project consists of the following components:

1. **Service 1 - User Register**: A Go service that produces messages to a Kafka topic.
2. **Service 2 - Notification**: A Go service that consumes messages from the Kafka topic and sends emails via SMTP.
3. **Kafka**: Apache Kafka is used as the message broker for communication between the two services.
4. **Kafka Exporter**: Exports Kafka metrics to Prometheus.
5. **Prometheus**: Collects and stores metrics from Kafka Exporter and other services.
6. **Grafana**: Visualizes the metrics collected by Prometheus.
7. **MailDev**: An SMTP server and web interface for intercepting and viewing emails sent by Service 2.
8. **Postgres**: Storage for users

## Prerequisites

Before running the project, ensure you have the following installed:

- Docker
- Docker Compose
- Go (for local development)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/AzizKHAN030/go-kafka.git
cd go-kafka
```

### 2. Docker composer

```bash
docker-composer up --build
```

### 3. Example query

- **Register user**
```bash
curl --request POST \
  --url http://localhost:8080/register \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "John",
	"email":"johndoe@gmail.com",
	"password":"JohnDoe123"
}'
```
##
- **Login user**
```bash
curl --request POST \
  --url http://localhost:8080/login \
  --header 'Content-Type: application/json' \
  --data '{
	"email":"johndoe@gmail.com",
	"password":"JohnDoe123"
}'
```
##

### 4. Check the messages

- **Reach port 1080 for MailDev**
