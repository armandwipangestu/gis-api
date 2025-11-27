<h1 align="center">An application backend or RESTful API server for <a href="https://github.com/armandwipangestu/gis-ui">https://github.com/armandwipangestu/gis-ui</a></h1>

<div align="center">

![Golang](https://img.shields.io/badge/-Golang-F9FAFC?style=for-the-badge&logo=go)&nbsp;
![MariaDB](https://img.shields.io/badge/-MariaDB-003545?style=for-the-badge&logo=mariadb&logoColor=white)&nbsp;
![Docker](https://img.shields.io/badge/-Docker-F9FAFC?style=for-the-badge&logo=docker)&nbsp;
![Bruno](https://img.shields.io/badge/-Bruno-F9FAFC?style=for-the-badge&logo=bruno)&nbsp;

</div>

<p align="center">A simple RESTful API for GIS Application built using <b>Golang</b>, <b>Gin</b>, and <b>GORM</b></p>

---

## Table of Contents

-   [Features](#features)
-   [Requirements](#requirements)
-   [Running the Application](#running-the-application)
    -   [Development Mode](#development-mode)
    -   [Compile Manual (Build Binary)](#compile-manual-build-binary)
    -   [Running with Docker](#running-with-docker)
    -   [Running with Docker Compose](#running-with-docker-compose)

---

## Features

-   Simple Migration & Seeder
-   Simple clean architecture: routes -> controller -> struct -> helper -> models
-   Support build manual, binary release, and Docker image

## Requirements

-   Go 1.25+
-   MariaDB 12.0.2
-   Git
-   Docker & Docker Compose (optional)

## Running the Application

### Development Mode

1. Clone Repostory & Install dependencies

```bash
git clone https://github.com/armandwipangestu/gis-api && cd gis-api
go mod tidy
```

2. Setup Environment Variable

```bash
cp .env.example .env
```

Fill with your own configuration

```bash
# App Configuration
APP_PORT=3000

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=
DB_NAME=golang_gis

# You can generate the random string using this tool https://jwtsecrets.com/#generator
JWT_SECRET=<random_string>
```

3. Create Database

```sql
CREATE DATABASE golang_gis;
```

4. Install air-verse (hot reload)

```bash
go install github.com/air-verse/air@latest
```

5. Running the application

> [!NOTE]
> Access the API at http://localhost:3000

```bash
air
```

### Compile Manual (Build Binary)

1. Compile to make executable file

> [!TIP]
> To compile for difference architecture (like Linux AMD64)
>
> ```bash
> GOOS=linux GOARCH=amd64 go build -o dist/gis-api ./main.go
> ```

```bash
go build -o dist/gis-api ./main.go
```

2. Setup Environment Variable

```bash
cp .env.example .env
```

Fill with your own configuration

```bash
# App Configuration
APP_PORT=3000

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=
DB_NAME=golang_gis

# You can generate the random string using this tool https://jwtsecrets.com/#generator
JWT_SECRET=<random_string>
```

3. Running the executable file

```bash
./dist/gis-api
```

### Running with Docker

1. Build the image

```bash
docker build -t gis-api .
```

2. Run the image

```bash
docker run -p 3000:3000 --env-file .env gis-api
```

### Running with Docker Compose

1. Copy the `.env.example` and `.env.example.mysql`

```bash
cp .env.example .env
cp .env.example.mysql .env.mysql
```

2. Fill the value of `.env` and `.env.mysql` with your own configuration

```bash
# App Configuration
APP_PORT=3000

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=
DB_NAME=golang_gis

# You can generate the random string using this tool https://jwtsecrets.com/#generator
JWT_SECRET=<random_string>
```

```bash
MYSQL_ROOT_PASSWORD="<your_password>"
MYSQL_DATABASE="golang_gis"
MYSQL_USER="root"
MYSQL_PASSWORD="<your_password>"
```

3. Runing the application using compose

```bash
docker compose up -d
```
