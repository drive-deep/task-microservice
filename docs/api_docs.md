Collecting workspace informationSure, here is the generated API documentation for your task microservice:

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
```

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
You can use the test.sh script to test the service endpoints and Kafka messaging. The script uses `curl` to send HTTP requests and Kafka CLI tools to send messages to Kafka topics.

To run the tests, follow these steps:
1. Make the script executable:
    ```sh
    chmod +x test.sh
    ```

2. Execute the script:
    ```sh
    ./test.sh
    ```

You can see the result of test.sh execution in test_exec_result.md.