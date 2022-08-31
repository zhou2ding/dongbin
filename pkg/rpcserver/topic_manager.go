package rpcserver

import (
	"blog/pkg/l"
	"blog/pkg/rpcservice"
	"encoding/json"
	"go.uber.org/zap"
	"sync"
	"time"
)

type subscriber interface {
	Id() int
	Valid() bool
	SendMessage(id int, value json.RawMessage)
}

type notifyMsg interface {
	GetDstId() int
	GetTopicName() string
	GetMsgValue() []byte
}

type topicManager struct {
	topics   map[string]*topic
	mutex    sync.Mutex
	notifyCh chan interface{}
	once     sync.Once
}

var gTopicManager *topicManager = &topicManager{
	topics:   make(map[string]*topic),
	notifyCh: make(chan interface{}, 100),
}

func TopicManagerInstance() *topicManager {
	return gTopicManager
}

func (t *topicManager) filterRPCSubscribe(key string, value json.RawMessage, subscriber subscriber) (bool, int, error) {
	var filtered = false
	var id = -1
	var err error

	serviceName, action := rpcservice.GetRPCServiceMgr().FilterAndPreprocess(key)
	if action == "attach" {
		filtered = true
		id, err = t.register(serviceName, subscriber)
	} else if action == "detach" {
		filtered = true
		err = t.unregister(serviceName, subscriber)
	}

	/*
		开启一个协程，定时检查订阅者的有效性
		如此操作的必要性在于，RPCSession的Close被调用和开启的go s.doExecute存在竞态问题
		即Close被调用时，可能恰有go s.doExecute在执行，如果是RPC订阅请求，则TopicManager在过滤时可能不能判断出
		RPCSession已经Close，从而开启了一个无效的RPC订阅。
	*/
	if filtered {
		t.once.Do(func() {
			rpcservice.GetRPCServiceMgr().SetNotifyChan(t.notifyCh)
			go t.checkValid()
			go t.notify()
		})
	}

	return filtered, id, err
}

func (t *topicManager) register(topicName string, subscriber subscriber) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var err error
	var notifyId int
	tp, ok := t.topics[topicName]
	if !ok {
		//启动相关服务组件的订阅功能
		l.GetLogger().Info("topicManager register start", zap.String("topicName", topicName), zap.Int("Id", subscriber.Id()))
		err = rpcservice.GetRPCServiceMgr().StartTopicListen(topicName)
		if err == nil {
			tp = &topic{}
			notifyId = tp.add(subscriber)
			t.topics[topicName] = tp
		}
	} else {
		l.GetLogger().Info("topicManager register add subscriber", zap.String("topicName", topicName), zap.Int("Id", subscriber.Id()))
		notifyId = tp.add(subscriber)
	}

	return notifyId, err
}

func (t *topicManager) unregister(topicName string, subscriber subscriber) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	tp, ok := t.topics[topicName]
	if ok {
		tp.remove(subscriber)
		l.GetLogger().Info("topicManager unregister remove", zap.String("topicName", topicName), zap.Int("Id", subscriber.Id()))
		if tp.size() == 0 {
			//已没有订阅者，关闭相关服务组件的订阅功能
			l.GetLogger().Info("topicManager unregister stop topic", zap.String("topicName", topicName))
			rpcservice.GetRPCServiceMgr().StopTopicListen(topicName)

			delete(t.topics, topicName)
		}
	}

	return nil
}

func (t *topicManager) remove(subscriber subscriber) {
	//fmt.Println("TopicManager::Remove")
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for topicName, tp := range t.topics {
		tp.remove(subscriber)
		l.GetLogger().Info("topicManager remove remove", zap.String("topicName", topicName), zap.Int("Id", subscriber.Id()))
		if tp.size() == 0 {
			//已没有订阅者，关闭相关服务组件的订阅功能
			l.GetLogger().Info("topicManager remove stop topic", zap.String("topicName", topicName))
			rpcservice.GetRPCServiceMgr().StopTopicListen(topicName)

			delete(t.topics, topicName)
		}
	}
}

func (t *topicManager) checkValid() {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		t.mutex.Lock()
		for topicName, tp := range t.topics {
			tp.checkValid()
			if tp.size() == 0 {
				//已没有订阅者，关闭相关服务组件的订阅功能
				l.GetLogger().Info("topicManager checkValid stop topic", zap.String("topicName", topicName))
				rpcservice.GetRPCServiceMgr().StopTopicListen(topicName)

				delete(t.topics, topicName)
			}
		}
		t.mutex.Unlock()
	}
}

func (t *topicManager) notify() {
	for i := range t.notifyCh {
		if msg, ok := i.(notifyMsg); ok {
			t.mutex.Lock()
			if tp, ok := t.topics[msg.GetTopicName()]; ok {
				tp.notify(msg.GetDstId(), msg.GetMsgValue())
			}
			t.mutex.Unlock()
		}
	}
}
