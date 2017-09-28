package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/bitwurx/jrpc2"
)

const (
	QueueNotFoundCode jrpc2.ErrorCode = -32002 // queue not found json rpc 2.0 error code.
)

const (
	QueueNotFoundMsg jrpc2.ErrorMsg = "Queue not found" // queue not found json rpc 2.0 error message.
)

// ApiV1 is the version 1 implementation of the rpc methods.
type ApiV1 struct {
	// model the priority queue database model.
	// queues A represetation of priority queues by key.
	model  Model
	queues map[string]*PriorityQueue
}

// GetParams contains the rpc parameters for the Get method.
type GetParams struct {
	// Key is the queue key.
	Key *string `json:"key"`
}

// FromPositional parses the key from the positional parameters.
func (params *GetParams) FromPositional(args []interface{}) error {
	if len(args) != 1 {
		return errors.New("key parameter is required")
	}
	key := args[0].(string)
	params.Key = &key

	return nil
}

// Get returns a queue by key.  An error is returned if the queue
//  does not exist.
func (api *ApiV1) Get(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(PushParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "queue key is required",
		}
	}
	queue, ok := api.queues[*p.Key]
	if !ok {
		return nil, &jrpc2.ErrorObject{
			Code:    QueueNotFoundCode,
			Message: QueueNotFoundMsg,
		}
	}
	return queue, nil
}

// GetAll returns all existing queues.
func (api *ApiV1) GetAll(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	queues := make([]*PriorityQueue, 0)
	for _, queue := range api.queues {
		queues = append(queues, queue)
	}
	return queues, nil
}

// PeekParams contains the rpc parameters for the Peek method.
type PeekParams struct {
	// Key is the queue key.
	Key *string `json:"key"`
}

// FromPositional parses the key from the positional parameters.
func (params *PeekParams) FromPositional(args []interface{}) error {
	if len(args) != 1 {
		return errors.New("key parameter is required")
	}
	key := args[0].(string)
	params.Key = &key

	return nil
}

// Peek returns the min node of the queue without deleting it.
func (api *ApiV1) Peek(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(PushParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task key is required",
		}
	}
	queue, ok := api.queues[*p.Key]
	if !ok {
		return nil, &jrpc2.ErrorObject{
			Code:    QueueNotFoundCode,
			Message: QueueNotFoundMsg,
		}
	}
	task := queue.Peek()
	if task != nil {
		return task, nil
	}
	return make(map[string]interface{}), nil
}

// PopParams contains the rpc parameters for the Pop method.
type PopParams struct {
	// Key is the queue key.
	Key *string `json:"key"`
}

// FromPositional parses the key from the positional parameters.
func (params *PopParams) FromPositional(args []interface{}) error {
	if len(args) != 1 {
		return errors.New("key parameter is required")
	}
	key := args[0].(string)
	params.Key = &key

	return nil
}

// Pop returns the min node of the queue and deletes it from the
// queue.
func (api *ApiV1) Pop(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(PopParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task key is required",
		}
	}
	queue, ok := api.queues[*p.Key]
	if !ok {
		return nil, &jrpc2.ErrorObject{
			Code:    QueueNotFoundCode,
			Message: QueueNotFoundMsg,
		}
	}
	task := queue.Pop()
	if task != nil {
		queue.Save(api.model)
		return task, nil
	}

	return make(map[string]interface{}), nil
}

// PushParams contains the rpc parameters fo the Push method.
type PushParams struct {
	// Key The resource key of the task.
	// Id the id of the task.
	// Priority the task priority value.
	Key      *string  `json:"key"`
	Id       *string  `json:"id"`
	Priority *float64 `json:"priority"`
}

// FromPositional parses the key, id, and priority from the
// positional parameters.
func (params *PushParams) FromPositional(args []interface{}) error {
	if len(args) != 3 {
		return errors.New("key, id, and priority parameters are required")
	}
	key := args[0].(string)
	id := args[1].(string)
	priority := args[2].(float64)
	params.Key = &key
	params.Id = &id
	params.Priority = &priority

	return nil
}

// Push adds the task to the queue with matching key. If the queue
// does not exist it will be created for insertion of the task.
func (api *ApiV1) Push(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(PushParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task key is required",
		}
	}
	if p.Id == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task id is required",
		}
	}
	if p.Priority == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task priority is required",
		}
	}

	var queue *PriorityQueue
	var ok bool

	if queue, ok = api.queues[*p.Key]; !ok {
		queue = NewPriorityQueue(*p.Key)
		api.queues[*p.Key] = queue
	}
	queue.Push(&Task{Id: *p.Id, Priority: *p.Priority})
	queue.Save(api.model)

	return 0, nil
}

// RemoveParams contains the rpc parameters for the Remove method
type RemoveParams struct {
	// Key is queue id.
	// Id the id of the task.
	Key *string `json:"key"`
	Id  *string `json:"id"`
}

// FromPositional parses the key and id from the positional
// parameters.
func (params *RemoveParams) FromPositional(args []interface{}) error {
	if len(args) != 2 {
		return errors.New("key, and id parameters are required")
	}
	key := args[0].(string)
	id := args[1].(string)
	params.Key = &key
	params.Id = &id

	return nil
}

// Remove removes the task from the queue
func (api *ApiV1) Remove(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(RemoveParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	if p.Key == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task key is required",
		}
	}
	if p.Id == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "task id is required",
		}
	}

	queue, ok := api.queues[*p.Key]
	if !ok {
		return nil, &jrpc2.ErrorObject{
			Code:    QueueNotFoundCode,
			Message: QueueNotFoundMsg,
		}
	}

	if err := queue.Remove(*p.Id); err != nil {
		return -1, nil
	}
	queue.Save(api.model)
	return 0, nil
}

// NewApiV1 returns a new api version 1 rpc api instance
func NewApiV1(model Model, s *jrpc2.Server) *ApiV1 {
	api := &ApiV1{model, make(map[string]*PriorityQueue)}
	queues, err := model.FetchAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, queue := range queues {
		v, _ := queue.(*PriorityQueue)
		api.queues[v.Key] = v
	}
	s.Register("get", jrpc2.Method{Method: api.Get})
	s.Register("getAll", jrpc2.Method{Method: api.GetAll})
	s.Register("peek", jrpc2.Method{Method: api.Peek})
	s.Register("pop", jrpc2.Method{Method: api.Pop})
	s.Register("push", jrpc2.Method{Method: api.Push})
	s.Register("remove", jrpc2.Method{Method: api.Remove})

	return api
}
