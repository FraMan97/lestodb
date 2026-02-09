# lestodb

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![BoltDB](https://img.shields.io/badge/BoltDB-003B57?style=for-the-badge&logo=databricks&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-GPL_v3-blue?style=for-the-badge&logo=gnu&logoColor=white)

`lestodb` is an ultra-lightweight NoSQL Key-Value database written in **Go**. It is designed to provide high-performance read and write operations through an **in-memory sharding** system, with optional disk persistence powered by **BoltDB**. It is ideal for scenarios requiring a fast cache with granular backup and restore capabilities, easily deployable via **Docker**.

---

## Table of Contents

* [Key Features](#key-features)
* [Environment Variables](#environment-variables)
* [Technology Stack](#technology-stack)
* [Installation and Setup](#installation-and-setup)
    * [Prerequisites](#prerequisites)
    * [Docker Deployment](#docker-deployment)
    * [How to Use](#how-to-use)

---

## Key Features

* **In-Memory Sharding**: Uses a **consistent hashing** system (FNV-1a) to distribute keys across multiple independent memory shards. This drastically reduces lock contention and enables high-speed concurrent operations.
* **Selective Persistence**: Implements a **Repository Pattern** to interface memory with **BoltDB**. Users can decide when to persist specific keys or the entire dataset using `BACKUP` and `RESTORE` commands.
* **TTL Management**: Supports automatic key expiration. A background worker constantly monitors shards to remove expired records, optimizing memory usage.
* **Batch Processing**: The TCP server supports sending multiple commands in a single request (separated by `;`), improving throughput for bulk operations.
* **Decoupled Architecture**: Storage logic is separated from the physical database through clear interfaces, making the system easily extensible to other persistence backends.

---

### Environment Variables

The `docker-compose.yaml` file allows you to customize the database behavior through several environment variables:

* **`SHARDING_COUNT`**: Defines the number of memory shards to use for concurrent access (default: `36`).
* **`PORT`**: The TCP port the server will listen on inside and outside the container (default: `2001`).
* **`URL`**: The network interface for the server, typically set to `0.0.0.0` for Docker environments.
* **`FILE_PATH`**: The internal container path where the BoltDB persistence file is stored (default: `/app/database/lesto.db`).

---

## Technology Stack

* **Core Language**: Go 1.23+
* **Storage Engine**: 
    * **In-Memory**: Sharded maps protected by `RWMutex`.
    * **BoltDB**: ACID-compliant key-value store for disk persistence.
* **Containerization**: Docker & Docker Compose.

---

## Installation and Setup

> [!NOTE]
> ### Prerequisites
> * Docker and Docker Compose installed.
> * A free TCP port (default `2001`).

### Docker Deployment

The fastest way to start `lestodb` is via Docker Compose. This will automatically configure volumes for data persistence on your local host.

## How to Use

```bash
# Clone the repository
git clone https://github.com/FraMan97/lestodb.git
cd lestodb

# Start the container
docker-compose up --build -d

# Connect through telnet
telnet localhost 2001

# SET: Save a key with TTL (in seconds) and value
SET user:1 3600 {"name":"Francesco"}

# GET: Retrieve the value
GET user:1

# DEL: Delete the key from memory
DEL user:1

# BACKUP: Save a specific key or the entire DB to BoltDB
BACKUP user:1
BACKUP ALL

# RESTORE: Retrieve the value from disk and restore it to RAM with a new TTL
RESTORE user:1 3600
RESTORE ALL 7200

SET a 60 valA; SET b 60 valB; BACKUP ALL

```