package rpcserver

import (
	"blog/pkg/rand"
)

type object struct {
	subscriber subscriber
	id         int
}

type topic struct {
	objects []*object
}

func (t *topic) size() int {
	return len(t.objects)
}

func (t *topic) add(subscriber subscriber) int {
	var notifyId int
	found := false
	for _, o := range t.objects {
		if subscriber == o.subscriber {
			found = true
			notifyId = o.id //如果订阅者已订阅过，返回之前的订阅id
			break
		}
	}

	if !found {
		obj := &object{
			subscriber: subscriber,
			id:         rand.GetIdGeneratorInstance().GetId(),
		}
		t.objects = append(t.objects, obj)
		notifyId = obj.id
	}

	return notifyId
}

func (t *topic) remove(subscriber subscriber) {
	var idx int
	for idx = 0; idx < len(t.objects); idx++ {
		if t.objects[idx].subscriber == subscriber { //found
			break
		}
	}

	if idx == len(t.objects) { //not found
		return
	}

	t.objects = append(t.objects[:idx], t.objects[idx+1:]...)
}

func (t *topic) checkValid() {
	//remove those invalid subscribes
	k := 0
	for _, obj := range t.objects {
		if obj.subscriber.Valid() { //filter
			t.objects[k] = obj
			k++
		}
	}
	t.objects = t.objects[:k]
}

func (t *topic) notify(dstId int, value []byte) {
	for _, obj := range t.objects {
		if obj.subscriber.Id() == dstId || dstId == -1 { //只发给对应客户端，或-1表示发给订阅该主题的所有客户端
			obj.subscriber.SendMessage(obj.id, value)
		}
	}
}
