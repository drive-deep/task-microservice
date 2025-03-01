#!/bin/bash

# URL of the task creation endpoint
URL="http://localhost:8080/tasks"

# Loop to create 10 tasks
for i in {1..10}
do
    # JSON payload for the task
    TASK=$(cat <<EOF
{
    "title": "Sample Task $i",
    "description": "This is the description for task $i.",
    "status": "pending",
    "priority": $i,
    "assignee": "user$i"
}
EOF
)

    # Send POST request to create the task
    curl -X POST $URL \
             -H "Content-Type: application/json" \
             -d "$TASK"

    echo "Created task $i"
done

echo "All tasks created successfully"

# URL of the get all tasks endpoint with query parameters
GET_URL1="http://localhost:8080/tasks?status=pending&priority=5&assignee=user1"
GET_URL2="http://localhost:8080/tasks?status=completed&priority=3"
GET_URL3="http://localhost:8080/tasks?assignee=user2"
GET_URL4="http://localhost:8080/tasks?status=pending"

# Send GET request to retrieve all tasks with the specified query parameters
curl -X GET $GET_URL1 \
        -H "Content-Type: application/json"
echo "Retrieved all tasks with status=pending, priority=5, assignee=user1"

curl -X GET $GET_URL2 \
        -H "Content-Type: application/json"
echo "Retrieved all tasks with status=completed, priority=3"

curl -X GET $GET_URL3 \
        -H "Content-Type: application/json"
echo "Retrieved all tasks with assignee=user2"

curl -X GET $GET_URL4 \
        -H "Content-Type: application/json"
echo "Retrieved all tasks with status=pending"


# URL of the get task by ID endpoint
GET_TASK_BY_ID_URL="http://localhost:8080/tasks/1"

# Send GET request to retrieve the task by ID
curl -X GET $GET_TASK_BY_ID_URL \
    -H "Content-Type: application/json"
echo "Retrieved task with ID=1"

# URL of the update task endpoint
UPDATE_TASK_URL="http://localhost:8080/tasks/1"

# JSON payload for the updated task
UPDATED_TASK=$(cat <<EOF
{
    "title": "Updated Task 1",
    "description": "This is the updated description for task 1.",
    "status": "completed",
    "priority": 10,
    "assignee": "user1"
}
EOF
)

# Send PUT request to update the task
curl -X PUT $UPDATE_TASK_URL \
    -H "Content-Type: application/json" \
    -d "$UPDATED_TASK"
echo "Updated task with ID=1"

# Send GET request to retrieve the updated task by ID
curl -X GET $GET_TASK_BY_ID_URL \
    -H "Content-Type: application/json"
echo "Retrieved updated task with ID=1"

# URL of the delete task endpoint
DELETE_TASK_URL="http://localhost:8080/tasks/1"

# Send DELETE request to delete the task
curl -X DELETE $DELETE_TASK_URL \
    -H "Content-Type: application/json"
echo "Deleted task with ID=1"

# Send GET request to retrieve the task by ID to confirm deletion
curl -X GET $GET_TASK_BY_ID_URL \
    -H "Content-Type: application/json"
echo "Attempted to retrieve deleted task with ID=1"

# URL of the Kafka topic endpoint
KAFKA_URL="http://localhost:8080/kafka/create_task"

# JSON payload for the task to be sent to Kafka
KAFKA_TASK=$(cat <<EOF
{   "id": "12",
    "title": "Kafka Task",
    "description": "This is a task sent to Kafka.",
    "status": "pending",
    "priority": 1,
    "assignee": "user_kafka"
}
EOF
)

# Send POST request to send the task to the Kafka topic
curl -X POST $KAFKA_URL \
    -H "Content-Type: application/json" \
    -d "$KAFKA_TASK"
echo "Sent task to Kafka topic"

# URL of the get task by ID endpoint
GET_KAFKA_TASK_BY_ID_URL="http://localhost:8080/tasks/12"

# Send GET request to retrieve the task by ID
curl -X GET $GET_KAFKA_TASK_BY_ID_URL \
    -H "Content-Type: application/json"
echo "Retrieved task with ID=1 from Kafka"

# URL of the Kafka topic endpoint for updating a task
KAFKA_UPDATE_URL="http://localhost:8080/kafka/update_task"

# JSON payload for the updated task to be sent to Kafka
KAFKA_UPDATED_TASK=$(cat <<EOF
{
    "id": "12",
    "title": "Updated Kafka Task",
    "description": "This is an updated task sent to Kafka.",
    "status": "completed",
    "priority": 2,
    "assignee": "user_kafka_updated"
}
EOF
)

# Send POST request to send the updated task to the Kafka topic
curl -X POST $KAFKA_UPDATE_URL \
    -H "Content-Type: application/json" \
    -d "$KAFKA_UPDATED_TASK"
echo "Sent updated task to Kafka topic"

# Send GET request to retrieve the updated task by ID from Kafka
curl -X GET $GET_KAFKA_TASK_BY_ID_URL \
    -H "Content-Type: application/json"
echo "Retrieved updated task with ID=12 from Kafka"

# URL of the Kafka topic endpoint for deleting a task
KAFKA_DELETE_URL="http://localhost:8080/kafka/delete_task"

# JSON payload for the task to be deleted to be sent to Kafka
KAFKA_DELETE_TASK=$(cat <<EOF
{
    "id": "12"
}
EOF
)

# Send POST request to send the delete task command to the Kafka topic
curl -X POST $KAFKA_DELETE_URL \
    -H "Content-Type: application/json" \
    -d "$KAFKA_DELETE_TASK"
echo "Sent delete task command to Kafka topic"

# Send GET request to retrieve the task by ID to confirm deletion from Kafka
curl -X GET $GET_KAFKA_TASK_BY_ID_URL \
    -H "Content-Type: application/json"
echo "Attempted to retrieve deleted task with ID=12 from Kafka"