# Ontrek Backend

This repository contains the backend for the Ontrek application.

## Requirements

- Docker installed on your machine

## How to run the backend with Docker

1. **Build the Docker image**

```bash
sudo docker build -t ontrek-backend .
```

2. **Run the container**

```bash
sudo docker run --name ontrek -p 3000:8080 ontrek-backend
```
