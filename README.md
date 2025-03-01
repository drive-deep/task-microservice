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

## Testing the Service

### Running Tests
You can use the `test.sh` script to test the service endpoints and Kafka messaging. The script uses `curl` to send HTTP requests and Kafka CLI tools to send messages to Kafka topics.

### Running Tests
You can use the `test.sh` script to test the service endpoints and Kafka messaging. The script uses `curl` to send HTTP requests and Kafka CLI tools to send messages to Kafka topics.

To run the tests, follow these steps:
1. Make the script executable:
    ```sh
    chmod +x test.sh
    ```

2. Execute the script:
    ```sh
    ./test.sh
    ```

You can see the result of `./test.sh` execution in `test_exec_result.md`.


```diff
- {"id":"1","title":"Sample Task 1","description":"This is the description for task 1.","status":"pending","priority":1,"created_at":"2025-03-01T05:52:26.042198Z","updated_at":"2025-03-01T05:52:26.042198Z"}
- Created task 1
- {"id":"2","title":"Sample Task 2","description":"This is the description for task 2.","status":"pending","priority":2,"created_at":"2025-03-01T05:52:26.062808Z","updated_at":"2025-03-01T05:52:26.062808Z"}
- Created task 2
- {"id":"3","title":"Sample Task 3","description":"This is the description for task 3.","status":"pending","priority":3,"created_at":"2025-03-01T05:52:26.078623Z","updated_at":"2025-03-01T05:52:26.078623Z"}
- Created task 3
- {"id":"4","title":"Sample Task 4","description":"This is the description for task 4.","status":"pending","priority":4,"created_at":"2025-03-01T05:52:26.096484Z","updated_at":"2025-03-01T05:52:26.096484Z"}
- Created task 4
- {"id":"5","title":"Sample Task 5","description":"This is the description for task 5.","status":"pending","priority":5,"created_at":"2025-03-01T05:52:26.117998Z","updated_at":"2025-03-01T05:52:26.117998Z"}
- Created task 5
- {"id":"6","title":"Sample Task 6","description":"This is the description for task 6.","status":"pending","priority":6,"created_at":"2025-03-01T05:52:26.142219Z","updated_at":"2025-03-01T05:52:26.142219Z"}
- Created task 6
- {"id":"7","title":"Sample Task 7","description":"This is the description for task 7.","status":"pending","priority":7,"created_at":"2025-03-01T05:52:26.163838Z","updated_at":"2025-03-01T05:52:26.163838Z"}
- Created task 7
- {"id":"8","title":"Sample Task 8","description":"This is the description for task 8.","status":"pending","priority":8,"created_at":"2025-03-01T05:52:26.180094Z","updated_at":"2025-03-01T05:52:26.180094Z"}
- Created task 8
- {"id":"9","title":"Sample Task 9","description":"This is the description for task 9.","status":"pending","priority":9,"created_at":"2025-03-01T05:52:26.197867Z","updated_at":"2025-03-01T05:52:26.197867Z"}
- Created task 9
- {"id":"10","title":"Sample Task 10","description":"This is the description for task 10.","status":"pending","priority":10,"created_at":"2025-03-01T05:52:26.218282Z","updated_at":"2025-03-01T05:52:26.218282Z"}
- Created task 10
- All tasks created successfully
- [{"id":"5","title":"Sample Task 5","description":"This is the description for task 5.","status":"pending","priority":5,"created_at":"2025-03-01T05:52:26.117998Z","updated_at":"2025-03-01T05:52:26.117998Z"}]
- Retrieved all tasks with status=pending, priority=5, assignee=user1
- []
- Retrieved all tasks with status=completed, priority=3
- [{"id":"1","title":"Sample Task 1","description":"This is the description for task 1.","status":"pending","priority":1,"created_at":"2025-03-01T05:52:26.042198Z","updated_at":"2025-03-01T05:52:26.042198Z"},{"id":"2","title":"Sample Task 2","description":"This is the description for task 2.","status":"pending","priority":2,"created_at":"2025-03-01T05:52:26.062808Z","updated_at":"2025-03-01T05:52:26.062808Z"},{"id":"3","title":"Sample Task 3","description":"This is the description for task 3.","status":"pending","priority":3,"created_at":"2025-03-01T05:52:26.078623Z","updated_at":"2025-03-01T05:52:26.078623Z"},{"id":"4","title":"Sample Task 4","description":"This is the description for task 4.","status":"pending","priority":4,"created_at":"2025-03-01T05:52:26.096484Z","updated_at":"2025-03-01T05:52:26.096484Z"},{"id":"5","title":"Sample Task 5","description":"This is the description for task 5.","status":"pending","priority":5,"created_at":"2025-03-01T05:52:26.117998Z","updated_at":"2025-03-01T05:52:26.117998Z"},{"id":"6","title":"Sample Task 6","description":"This is the description for task 6.","status":"pending","priority":6,"created_at":"2025-03-01T05:52:26.142219Z","updated_at":"2025-03-01T05:52:26.142219Z"},{"id":"7","title":"Sample Task 7","description":"This is the description for task 7.","status":"pending","priority":7,"created_at":"2025-03-01T05:52:26.163838Z","updated_at":"2025-03-01T05:52:26.163838Z"},{"id":"8","title":"Sample Task 8","description":"This is the description for task 8.","status":"pending","priority":8,"created_at":"2025-03-01T05:52:26.180094Z","updated_at":"2025-03-01T05:52:26.180094Z"},{"id":"9","title":"Sample Task 9","description":"This is the description for task 9.","status":"pending","priority":9,"created_at":"2025-03-01T05:52:26.197867Z","updated_at":"2025-03-01T05:52:26.197867Z"},{"id":"10","title":"Sample Task 10","description":"This is the description for task 10.","status":"pending","priority":10,"created_at":"2025-03-01T05:52:26.218282Z","updated_at":"2025-03-01T05:52:26.218282Z"}]
- Retrieved all tasks with assignee=user2
- [{"id":"1","title":"Sample Task 1","description":"This is the description for task 1.","status":"pending","priority":1,"created_at":"2025-03-01T05:52:26.042198Z","updated_at":"2025-03-01T05:52:26.042198Z"},{"id":"2","title":"Sample Task 2","description":"This is the description for task 2.","status":"pending","priority":2,"created_at":"2025-03-01T05:52:26.062808Z","updated_at":"2025-03-01T05:52:26.062808Z"},{"id":"3","title":"Sample Task 3","description":"This is the description for task 3.","status":"pending","priority":3,"created_at":"2025-03-01T05:52:26.078623Z","updated_at":"2025-03-01T05:52:26.078623Z"},{"id":"4","title":"Sample Task 4","description":"This is the description for task 4.","status":"pending","priority":4,"created_at":"2025-03-01T05:52:26.096484Z","updated_at":"2025-03-01T05:52:26.096484Z"},{"id":"5","title":"Sample Task 5","description":"This is the description for task 5.","status":"pending","priority":5,"created_at":"2025-03-01T05:52:26.117998Z","updated_at":"2025-03-01T05:52:26.117998Z"},{"id":"6","title":"Sample Task 6","description":"This is the description for task 6.","status":"pending","priority":6,"created_at":"2025-03-01T05:52:26.142219Z","updated_at":"2025-03-01T05:52:26.142219Z"},{"id":"7","title":"Sample Task 7","description":"This is the description for task 7.","status":"pending","priority":7,"created_at":"2025-03-01T05:52:26.163838Z","updated_at":"2025-03-01T05:52:26.163838Z"},{"id":"8","title":"Sample Task 8","description":"This is the description for task 8.","status":"pending","priority":8,"created_at":"2025-03-01T05:52:26.180094Z","updated_at":"2025-03-01T05:52:26.180094Z"},{"id":"9","title":"Sample Task 9","description":"This is the description for task 9.","status":"pending","priority":9,"created_at":"2025-03-01T05:52:26.197867Z","updated_at":"2025-03-01T05:52:26.197867Z"},{"id":"10","title":"Sample Task 10","description":"This is the description for task 10.","status":"pending","priority":10,"created_at":"2025-03-01T05:52:26.218282Z","updated_at":"2025-03-01T05:52:26.218282Z"}]
- Retrieved all tasks with status=pending
- {"id":"1","title":"Sample Task 1","description":"This is the description for task 1.","status":"pending","priority":1,"created_at":"2025-03-01T05:52:26.042198Z","updated_at":"2025-03-01T05:52:26.042198Z"}
- Retrieved task with ID=1
- {"id":"1","title":"Updated Task 1","description":"This is the updated description for task 1.","status":"completed","priority":10,"created_at":"0001-01-01T00:00:00Z","updated_at":"2025-03-01T05:52:26.370933373Z"}
- Updated task with ID=1
- {"id":"1","title":"Updated Task 1","description":"This is the updated description for task 1.","status":"completed","priority":10,"created_at":"0001-01-01T00:00:00Z","updated_at":"2025-03-01T05:52:26.370933373Z"}
- Retrieved updated task with ID=1
- Deleted task with ID=1
- record not found
- Attempted to retrieve deleted task with ID=1
- 404 page not found
- Sent task to Kafka topic
- record not found
- Retrieved task with ID=1 from Kafka
- 404 page not found
- Sent updated task to Kafka topic
- record not found
- Retrieved updated task with ID=12 from Kafka
- 404 page not found
- Sent delete task command to Kafka topic
- record not found
- Attempted to retrieve deleted task with ID=12 from Kafka
```