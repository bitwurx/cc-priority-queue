package main

import (
	"encoding/json"
	"testing"

	"github.com/bitwurx/jrpc2"
)

func TestApiV1Get(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	result, errObj := api.Push([]byte(`{"key": "get", "id": "abc123", "priority": 2.3}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	queue, errObj := api.Get([]byte(`{"key": "get"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	m := make(map[string]interface{})
	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if m["_key"].(string) != "get" {
		t.Fatal("expected key to be 'get'")
	}
	if m["count"].(float64) != 1 {
		t.Fatal("expected count to be 1")
	}
}

func TestApiV1GetAll(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	_, errObj := api.Push([]byte(`{"key": "k1", "id": "abc123", "priority": 2.3}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	_, errObj = api.Push([]byte(`{"key": "k2", "id": "abc123", "priority": 2.3}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	queues, errObj := api.GetAll([]byte(`{"key": "getAll"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	var q []map[string]interface{}
	data, err := json.Marshal(queues)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &q); err != nil {
		t.Fatal(err)
	}
	keys := make([]string, 0)
	keys = append(keys, q[0]["_key"].(string))
	keys = append(keys, q[1]["_key"].(string))

	for _, k := range keys {
		if k != "k1" && k != "k2" {
			t.Fatal("got unexpected queue key")
		}
	}
}

func TestApiV1Peek(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	result, errObj := api.Push([]byte(`{"key": "abc", "id": "111", "priority": 0.5}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	result, errObj = api.Peek([]byte(`{"key": "abc"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	var task map[string]interface{}
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &task); err != nil {
		t.Fatal(err)
	}
	if task["priority"].(float64) != 0.5 {
		t.Fatal("expected task priority to be 0.5")
	}
	api.queues["abc"].Pop()
}

func TestApiV1Pop(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	result, errObj := api.Push([]byte(`{"key": "key-abc", "id": "111", "priority": 3.9}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	result, errObj = api.Pop([]byte(`{"key": "key-abc"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	var task map[string]interface{}
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &task); err != nil {
		t.Fatal(err)
	}
	if task["priority"].(float64) != 3.9 {
		t.Fatal("expected task priority to be 3.9")
	}
}

func TestApiV1Push(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	result, err := api.Push([]byte(`{"key": "test1", "id": "abc123", "priority": 2.3}`))
	if err != nil {
		t.Fatal(err)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
}

func TestApiV1Remove(t *testing.T) {
	api := NewApiV1(&MockModel{}, jrpc2.NewServer("", ""))
	result, errObj := api.Push([]byte(`{"key": "test1", "id": "abc123", "priority": 1.3}`))
	if errObj != nil {
		t.Fatal(errObj)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	result, errObj = api.Push([]byte(`{"key": "test1", "id": "abc321", "priority": 9.3}`))
	if errObj != nil {
		t.Fatal(errObj)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	result, errObj = api.Push([]byte(`{"key": "test1", "id": "abcxyz", "priority": 5.3}`))
	if errObj != nil {
		t.Fatal(errObj)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	result, errObj = api.Remove([]byte(`{"key": "test1", "id": "9g49g44"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != -1 {
		t.Fatal("expected result to be -1")
	}
	result, errObj = api.Remove([]byte(`{"key": "test1", "id": "abc321"}`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != 0 {
		t.Fatal("expected result to be -1")
	}
	for _, task := range api.queues["test1"].List() {
		if task.Id == "abc321" {
			t.Fatal("expected task with id 'abc321' to be removed")
		}
	}
}
