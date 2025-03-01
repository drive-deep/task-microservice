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
- **Kafka**: Used for asynchronous messaging to handle `task_create`, `task_update`, and `task_delete` events, which helps in scaling the service.

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
    #### Get Tasks with Sorting, Filtering, and Pagination
    - **URL**: `/tasks`
    - **Method**: `GET`
    - **Query Parameters**:
        - `sort_by` (optional): Field to sort by (e.g., `title`, `priority`, `created_at`)
        - `order` (optional): Sort order (`asc` for ascending, `desc` for descending)
        - `status` (optional): Filter by task status (e.g., `Pending`, `Completed`)
        - `priority` (optional): Filter by task priority (e.g., `1`, `2`)
        - `page` (optional): Page number (default is `1`)
        - `page_size` (optional): Number of tasks per page (default is `10`)

    - **Example Request**:
        ```
        GET /tasks?sort_by=priority&order=asc&status=Pending&page=1&page_size=5
        ```

    - **Response**:
        ```json
        {
            "tasks": [
                {
                    "id": "1",
                    "title": "Sample Task",
                    "description": "This is a sample task",
                    "status": "Pending",
                    "priority": 1,
                    "created_at": "2025-02-28T00:00:00Z",
                    "updated_at": "2025-02-28T00:00:00Z"
                }
            ],
            "page": 1,
            "page_size": 5,
            "total_pages": 1,
            "total_tasks": 1
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

## Kafka Message Queue

### Overview
The service uses Kafka for asynchronous messaging to handle `task_create`, `task_update`, and `task_delete` events. This helps in scaling the service by decoupling the task processing logic from the main application flow.

### Message Format
The messages sent to Kafka topics are in the following format:
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
### Example Usage

#### Creating a Task
When a message is sent to the `task_create` topic, a new task will be created:
```json
{
    "id": "1",
    "title": "New Task",
    "description": "This is a new task",
    "status": "Pending",
    "priority": 1,
    "created_at": "2025-02-28T00:00:00Z",
    "updated_at": "2025-02-28T00:00:00Z"
}
```

#### Updating a Task
When a message is sent to the `task_update` topic, the specified task will be updated:
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

#### Deleting a Task
When a message is sent to the `task_delete` topic, the specified task will be deleted:
```json
{
    "id": "1",
    "title": "New Task",
    "description": "This is a new task",
    "status": "Pending",
    "priority": 1,
    "created_at": "2025-02-28T00:00:00Z",
    "updated_at": "2025-02-28T00:00:00Z"
}
```