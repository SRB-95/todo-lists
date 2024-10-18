# todo-lists 
## curl commands
- create Task: curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"name":"Split_Wise_service","deadline":"2024-10-18T17:00:00+05:30","tag":"high"}'
- get all tasks: curl -X GET http://localhost:8080/tasks
- get task by id: curl -X GET http://localhost:8080/tasks/1
- get task by tag: curl -X GET http://localhost:8080/tasks/tag/high
- update task: curl -X PUT http://localhost:8080/tasks/8 -H "Content-Type: application/json" -d '{"name":"Todo_lists_service","deadline":"2024-10-18T17:00:00+05:30","tag":"high"}'
- search task by name: curl -X GET "http://localhost:8080/tasks/search?keyword=new"
- filter tasks by date-range: curl -X GET "http://localhost:8080/tasks/filter?start=2024-01-01&end=2024-12-31"
- delete task by id: curl -X DELETE "http://localhost:8080/tasks/{id}"