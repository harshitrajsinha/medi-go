# MediGo - A Hospital Management System written in Go

![frontend](./assets/images/frontend.png)

### 🚀 A Golang web application that consists of a receptionist portal & doctor portal which performs the following tasks:

- A single login page for both portals.
- Receptionists can register a new patient & perform CRUD operations.
- Doctors can view registered patient-related details and diagnose based on symptoms
- Gemini AI to generate diagnosis based on symptoms.

## Deployment links

Frontend: https://medigo-frontend.vercel.app \
Backend: https://medigo-7u9l.onrender.com/

## 👨‍🚀 Postman Collection Documentation

[Postman Documentation URL](https://documenter.getpostman.com/view/40689865/2sB34kDdn6)

## ✨Features

- 🔗 **Dependency Injection** for modularity
- ⌛ **API Versioning** for backward compatibility
- 💾 **Persistant storage** using PostgreSQL
- ⚡**Caching** - To reduce server load
- 🗐 **Pagination** - To efficiently handle and deliver large datasets
- 🚧 **Rate Limit** - To protect server resources
- 🔒 **JWT Authentication** for security

## 📦 Tech Stack

- **Backend** : Golang
- **AI Assitant**: Gemini
- **Frontend** : HTML/CSS/JavaScript
- **Database** : PostgreSQL
- **Caching** : Redis
- **Containerization** : Docker
- **Deployment** : Vercel (Frontend) + Render (Backend)

## How to run this application locally

### 1. Prerequisites

Make sure you have the `Docker Desktop` installed on your system:

### 2. Clone the Repository

```bash
git https://github.com/harshitrajsinha/medi-go.git
cd medi-go
```

### 3. Set up environment variable

Create .env file in root directory

```bash
# For postgres docker image
POSTGRES_USER=postgres
POSTGRES_PASSWORD=yourstrongpassword
POSTGRES_DB=yourfavouritedbname

# Data storage
DB_USER=postgres
DB_NAME=yourfavouritedbname
DB_PASS=yourstrongpassword
DB_PORT=5432
DB_HOST=db

JWT_KEY = secretkeyword

REDIS_HOST=redis
REDIS_PORT=6379
```

### 4. Run the application

Run the following command in your bash terminal

```bash
docker-compose up --build
```

## Receptionist Dashboard

![Screenshot](./assets/images/Screenshot%202025-07-19%20185445.png)

## Doctor Diagnosis (using Gemini AI)

![Screenshot](./assets/images/Screenshot%202025-07-20%20145456.png)

## Patient Report

![Screenshot](./assets/images/Screenshot%202025-07-20%20145535.png)
