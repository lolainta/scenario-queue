# Scenario Queue

## Architecure

### Manager

The manager is responsible for managing the scenario queue, including adding, removing, and updating scenarios. It interacts with the database to store and retrieve scenario information.

- Writting in Rust
- Using SeaORM as ORM
- Exposing REST API using Actix-web
- Storing data in PostgreSQL

### Worker

The worker is responsible for executing the scenarios in the queue. It fetches scenarios from the manager and runs them using the specified simulator and sampler.
The worker is expected to be spreaded in a SLURM cluster.

- Writting in Python
- Communicating with the manager via REST API

## Usage

### Manager

1. Configure the database metadata in both `db/.env` and `manager/.env`.
2. Start the PostgreSQL database using Docker Compose:
   ```bash
   cd db
   docker-compose up -d
   ```
3. Start the manager server:
   ```bash
   cargo run --release --bin manager
   ```

### Worker

To be implemented.
