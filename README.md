# Task Microservice

## Problem Breakdown and Design Decisions

### Problem Breakdown
The task microservice is designed to manage tasks with the following attributes:
- ID
- Title
- Description
- Status
- Priority
- CreatedAt
- UpdatedAt

The service provides CRUD operations for tasks and uses PostgreSQL for persistent storage and Redis for caching to improve performance.

### Design Decisions
- **PostgreSQL**: Chosen for its robustness and support for complex queries.
- **Redis**: Used as a caching layer to reduce database load and improve response times.
- **Docker Compose**: Used to orchestrate the microservice and its dependencies (PostgreSQL, Redis, Kafka, Zookeeper).
- **Gorilla Mux**: Used for routing HTTP requests.
- **Citus**: Used to scale out PostgreSQL horizontally.

## Instructions to Run the Service

### Prerequisites
- Docker
- Docker Compose

### Steps
1. Clone the repository:
    ```sh
    git clone https://github.com/your-repo/task-microservice.git
    cd task-microservice
    ```

2. Build and run the services using Docker Compose:
    ```sh
    docker-compose up --build
    ```

3. The service will be available at `http://localhost:8080`.

## API Documentation

### Endpoints

#### Create a Task
- **URL**: `/tasks`
- **Method**: `POST`
- **Request Body**:
    ```json
    {
        "id": "1",
        "title": "Sample Task",
        "description": "This is a sample task",
        "status": "Pending",
        "priority": 1,
        "created_at": "2025-02-28T00:00:00Z",
        "updated_at": "2025-02-28T00:00:00Z"
    }
    ```
- **Response**:
    ```json
    {
        "id": "1",
        "title": "Sample Task",
        "description": "This is a sample task",
        "status": "Pending",
        "priority": 1,
        "created_at": "2025-02-28T00:00:00Z",
        "updated_at": "2025-02-28T00:00:00Z"
    }
    ```

#### Get a Task by ID
- **URL**: `/tasks/{id}`
- **Method**: `GET`
- **Response**:
    ```json
    {
        "id": "1",
        "title": "Sample Task",
        "description": "This is a sample task",
        "status": "Pending",
        "priority": 1,
        "created_at": "2025-02-28T00:00:00Z",
        "updated_at": "2025-02-28T00:00:00Z"
    }
    ```

#### Get All Tasks
- **URL**: `/tasks`
- **Method**: `GET`
- **Response**:
    ```json
    [
        {
            "id": "1",
            "title": "Sample Task",
            "description": "This is a sample task",
            "status": "Pending",
            "priority": 1,
            "created_at": "2025-02-28T00:00:00Z",
            "updated_at": "2025-02-28T00:00:00Z"
        }
    ]
    ```

#### Update a Task
- **URL**: `/tasks/{id}`
- **Method**: `PUT`
- **Request Body**:
    ```json
    {
        "title": "Updated Task",
        "description": "This is an updated task",
        "status": "Completed",
        "priority": 2,
        "created_at": "2025-02-28T00:00:00Z",
        "updated_at": "2025-02-28T00:00:00Z"
    }
    ```
- **Response**:
    ```json
    {
        "id": "1",
        "title": "Updated Task",
        "description": "This is an updated task",
        "status": "Completed",
        "priority": 2,
        "created_at": "2025-02-28T00:00:00Z",
        "updated_at": "2025-02-28T00:00:00Z"
    }
    ```

#### Delete a Task
- **URL**: `/tasks/{id}`
- **Method**: `DELETE`
- **Response**: `204 No Content`

## Explanation of Microservices Concepts

### Scalability
- The service uses Citus to scale out PostgreSQL horizontally, allowing it to handle more load by distributing data across multiple nodes.

### Resilience
- The use of Redis as a caching layer helps to reduce the load on the database and improve response times, making the service more resilient to high traffic.

### Isolation
- Each component (PostgreSQL, Redis, Kafka, Zookeeper) runs in its own container, ensuring that they are isolated from each other and can be managed independently.

### Flexibility
- The service can be easily extended with new features or components by adding new services to the Docker Compose file.

### Observability
- Logs are printed to the console, and you can use tools like `docker logs` to monitor the output of each container.

### Conclusion
This task microservice demonstrates key microservices concepts such as scalability, resilience, isolation, flexibility, and observability. By using Docker Compose, it provides a simple way to orchestrate and manage the various components required for the service.