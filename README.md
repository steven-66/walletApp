# Wallet Application - README

## Overview

This repository contains the implementation of a **Wallet Application** backend written in **Go**. The application provides functionality for managing user wallets, including depositing money, withdrawing money, transferring money, checking balances, and viewing transaction history. The backend is built with **GORM** for database interactions, **PostgreSQL** for persistent storage, and **Docker** for containerization.

---

## Key Decisions Made

### Use of GORM
- GORM was chosen as the ORM to simplify database interactions and schema migrations.

### Dockerized Environment
- The application is fully containerized using Docker, with dependencies like PostgreSQL and Redis managed via Docker Compose.

### Table-Driven Tests
- Unit tests were written in a **table-driven format** for scalability and readability.

### CLI used for visualization
- A CLI tool was developed to visualize interactions with the application.

---

## How to Set Up and Run the Code

### Prerequisites
1. **Install Docker**: Ensure Docker is installed and running on your machine.
2. **Install Go**: Ensure Go is installed (version 1.19 or higher).

### Steps to Run
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/steven-66/walletApp.git
   cd walletApp
   ```

2. **Build and Start Services**:
   - Start PostgreSQL, Redis, and the application using Docker Compose:
     ```bash
     docker-compose up --build -d
     ```

3. **Run the Application**:
   - Interact with the application via the CLI:
     ```bash
     docker exec -it wallet_cli_app ./wallet-cli
     ```

4. **Run Unit Tests**:
   - Run unit tests directly on your local machine:
     ```bash
     go test ./... -v
     ```

## How Should Reviewers View the Code?

1. **Architecture**
    - The application is structured as the **MVC-alike pattern**:
        - **model**: Define the database schema.
        - **dto**: For data transferring between client and server.
        - **server/handler**: Handle business logic.
        - **storage**: Abstract database operations as dao layers.

2. **Unit Tests**
   - Each component has its own unit tests and each dependency has been properly mocked.
   - Table-driven format for scalability and readability.
   
3. **Error Handling**
    - Error handling is implemented throughout the application.

---

## Areas to Be Improved

1. **Performance Optimization**
   - Setup a message queue (e.g., Kafka) to handle the transfer of money between wallets. This will help to decouple the balance table from the transfer process, reducing the load on the balance table.
   - Make sure the idempotency of the transfer operation, so that the same transfer operation can be retried for a single request.
   - Implement write-through Redis caching for strong data consistency and reduce load on DB. i.e if the user frequently checks their balance, the balance should be cached for a certain period of time.

4. **Data Consistency**
   - For the transfer/deposit/withdraw operation, ensure that the entire operation is atomic. We can either implement a database transaction with `For Update` sql statement to achieve row-level lock, or use a distributed lock to ensure that the transfer is atomic.
   - Consider implementing a pessimistic read-write lock to handle concurrent access to the balance table if the concurrency of read is high. For example, adding a version number to the balance table and checking the version number before updating the balance, or compare the original amount when updating the balance using query `update balance set amount = amount + <amount> where id = <id> and amount = <original_amount>`.
   
2. **Logging + metrics**
    - Add structured logging (e.g., using `logrus`) for better observability and debugging.
    - Implement metrics for monitoring and alerting when certain db operations are not working as expected.

3. **Database Schema**
    - Add indexes on frequently queried columns (e.g., `user_id` in the `balance` table) for improved query performance.



5. **Testing**
    - Add integration tests to validate the end-to-end functionality of the application.

6. **Production Readiness**
    - Consider converting the application to a gRPC service for better performance and scalability.
    - Add health checks and readiness probes for Kubernetes deployments.

---

## Time Spent on the Test

I spent approximately **12 hours** on the test, broken down as follows:
- **Implementing core functionality**: ~8 hours
- **Writing unit tests**: ~2 hours
- **Containerizing the application**: ~1 hours
- **Debugging and refining**: ~1 hours

---

## Features Chosen Not to Implement

1. **Authentication**
    - User authentication and authorization were not included in the submission to focus on core wallet functionalities.
    - User signup and login were not implemented for simplicity.

2. **Integration Tests**
    - Only unit tests were implemented due to time constraints.

3. **Advanced Caching**
    - Features like Redis cache expiration and invalidation were not implemented as they were not part of the initial requirements.