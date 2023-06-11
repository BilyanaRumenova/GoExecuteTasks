# GoExecuteTasks

The project is a HTTP job processing service written in Golang.
A job is a collection of tasks, where each task has a name and a shell command. Tasks may
depend on other tasks and require that those are executed beforehand. The service takes care
of sorting the tasks to create a proper execution order.

Additionally, the service is able to return a bash script representation directly thus allowing us to run the commands directly from shell:
$ curl -d "@mytasks.json" -X POST http://localhost:4000/process_tasks | bash
