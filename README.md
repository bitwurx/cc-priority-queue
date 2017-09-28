# Concord PQ

Concord PQ is a priority queue for queueing tasks for a specific resource.  The PQ priority is time based with lowest runtimes tasks taking prescedence over tasks with high runtimes.

### JSON-RPC 2.0 HTTP API - Method Reference

-
#### get(key) : get a queue by key
-

#### Parameters:

key - (*String*) the queue key.

#### Returns:
(*Object*) the queue with the associated key

-
#### getAll() : get all queues
-

#### Returns:
(*Array*) the list of all existing queues 

-
#### peek(key) : return the next task from the queue
-

#### Parameters:

key - (*String*) the queue key.

#### Returns:
(*Object*) the next task in queue

-
#### pop(key) : remove and return the next task from the queue
-

#### Parameters:

key - (*String*) the queue key.

#### Returns:
(*Object*) the next task in queue

-
#### push(key, id, priority) : add a task to a queue
-

#### Parameters:

key - (*String*) the resource key for the task.

id - (*String*) the id of the task.

priority - (*Number*) the priority value for the task.  
<sub><sup>*Lower values have highest priority*</sup></sub>.

#### Returns:	
(*Number*)

-
#### remove(key, id) - remove a task from a queue
-

#### Parameters:

key - (*String*) the queue key.

id - (*String*) the id of the task.

#### Returns:
(*Number*)