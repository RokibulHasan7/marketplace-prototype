# Marketplace Prototype

A microservices-based application for managing and deploying applications on Kubernetes and VM environments.

## Table of Contents
- [User Guide](#user-guide)
    - [Installation](#installation)
    - [Workflow Example](#workflow-example)
- [Developer Guide](#developer-guide)
    - [Architecture](#architecture)
    - [Project Structure](#project-structure)

## User Guide

### Installation

#### Prerequisites
- Docker
- PostgreSQL
- Redis
- Go (1.19+)

#### Setup

After cloning this repository run the following commands to start the required services:

```sh
export DATABASE_URL="postgres://admin:secret@localhost:5433/marketplace"

docker run --name marketplace-db -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=marketplace -p 5433:5432 -d postgres

docker run --name marketplace-redis -p 6370:6379 -d redis

go run main.go
```

or you can run:

```sh
export DATABASE_URL="postgres://admin:secret@localhost:5433/marketplace"

make run
```

## Workflow Example

Follow this step-by-step guide to use the Marketplace Prototype APIs:

### 1. Create Two Users

First, create two users using the following API:

```sh
curl -X POST http://localhost:3000/api/users/new \
  -H "Content-Type: application/json" \
  -d '{
    "name": "User 1"
  }'

curl -X POST http://localhost:3000/api/users/new \
  -H "Content-Type: application/json" \
  -d '{
    "name": "User 2"
  }'
```

### 2. Create an Application for User 1

Next, create an application for User 1:
```shell
curl -X POST http://localhost:3000/api/apps/new \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kubernetes App",
    "description": "This is a Kubernetes-based application",
    "publisher_id": 1,
    "hourly_rate": 1.1,
    "deployment" :{
      "type": "k8s",
      "repoURL": "https://charts.bitnami.com/bitnami",
      "chartName": "nginx",
      "image": "",
      "cpu": "",
      "memory": ""
    }
  }'
```

### 3. Create a Project for User 2
Now, create a project under User 2:

```shell
curl -X POST http://localhost:3000/api/user/project/new \
  -H "Content-Type: application/json" \
  -d '{"name": "p1", "user_id": 2}'
```

### 4. Deploy an Application for User 1 in Project 1

Deploy the previously created application (App 1) under Project 1:
```shell
curl -X POST http://localhost:3000/api/deployments/install \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id": 2,
    "application_id": 1,
    "project_id": 1
  }'
```

### 5. Get the billing info by user id and deployment id

We have a background task which update billing records in every 5 min. So after a deployment if you call this api you will see the amount you charged for this deployment.
```shell
curl -X GET http://localhost:3000/api/billing/user/2/deployment/1 \
  -H "Content-Type: application/json"
```

### 6. Delete a deployment
```shell
curl -X DELETE http://localhost:3000/api/deployments/1 \
  -H "Content-Type: application/json"
```

There are also some others apis to Get the details of application, List application, Delete application, Get Deployment info, List Deployments etc.
You can see the `/internal/handlers/hendlers.go` file to see the api endpoints.

## Developer Guide

### Architecture
The marketplace-prototype project is designed to handle user applications and deployments in a cloud environment. It features a publisher-consumer mechanism, a background billing system, and supports both installation and uninstallation of deployments asynchronously using Redis queues. The system is built with Go, uses PostgreSQL as its database, and leverages Redis for task queuing.
#### Key Components
1️⃣ Deployment Management

    Implements two interfaces:
        Installer: Handles application deployments.
        Cleaner: Handles application uninstallations.
    Supports multiple deployment types (e.g., Kubernetes, Virtual Machines).
    Uses dependency injection for flexible deployment management.

2️⃣ Publisher-Consumer Mechanism (Async Processing)

    Uses a Redis queue to asynchronously process deployment and uninstallation requests.
    A publisher adds deployment/uninstallation tasks to the queue.
    A consumer worker picks up the tasks and executes them in the background.

3️⃣ Queue System (Redis)

    Redis is used as a message queue to decouple request handling from execution.
    Enables non-blocking API responses.
    Queues:
        installer_queue: Handles application installations.
        uninstaller_queue: Handles application uninstallations.

4️⃣ Database Layer (PostgreSQL)

    PostgreSQL is used to store:
        Users
        Projects
        Applications
        Deployments
        Billing records
    Ensures data consistency and persistence.

5️⃣ Billing System

    A background task runs periodically to calculate usage-based billing.
    Fetches deployment durations and applies hourly rates to generate cost records.
    Provides APIs to query user-specific and deployment-specific billing records.

6️⃣ REST API

    Built using Go (Golang) with the chi router for handling HTTP requests.
    Implements CRUD operations for:
        Users
        Projects
        Applications
        Deployments
        Billing records


This architecture ensures scalability, asynchronous processing, and separation of concerns, making the system flexible and efficient. 🚀

## Project Structure

```shell


```