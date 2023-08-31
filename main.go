package bbmessagequeue

import (
	"sync"

	"github.com/gunnarsutter/blackboard"
)

// Queue of messages in internal format
type MessageQueue struct {
	mutex sync.Mutex
	queue []QueueMessage
}

// Message in internal format
type QueueMessage struct {
	senderID    string
	messageType int
	content     string
}

// Sets values of the struct varables
func (qm *QueueMessage) Init(s_id string, m_t int, c string) {
	qm.senderID = s_id
	qm.messageType = m_t
	qm.content = c
}

// Converts a gRPC message to an internal message
func (qm *QueueMessage) New(message *blackboard.Message) {
	qm.senderID = message.SenderID
	qm.messageType = int(message.MessageType)
	qm.content = message.Content
}

// Adds an internal message to the end of the internal message queue
func (mq *MessageQueue) Push(m *blackboard.Message) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	var qm QueueMessage
	qm.New(m)
	mq.queue = append(mq.queue, qm)
}

// Adds an internal message to the front of the internal message queue
func (mq *MessageQueue) BbmPushFront(m *blackboard.Message) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	var qm QueueMessage
	qm.New(m)
	mq.queue = append([]QueueMessage{qm}, mq.queue...)
}

// Adds a blackboard message to the front of the internal message queue
func (mq *MessageQueue) PushFront(qm *QueueMessage) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	mq.queue = append([]QueueMessage{*qm}, mq.queue...)
}

// Returns and removes an internal message from the front of the queue
func (mq *MessageQueue) Pop() *QueueMessage {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()
	ret_message := mq.queue[0]
	mq.queue = mq.queue[1:]
	return &ret_message
}

// Returns the length of the internal message queue
func (mq *MessageQueue) Length() int {
	return len(mq.queue)
}

// Returns a gRPC message from an internal message
func (qm *QueueMessage) ToBlackboardMessage(msg *blackboard.Message) {
	msg.Content = qm.content //DETTA CRASHAR!!! Fortfarande?
	msg.MessageType = blackboard.MessageType(qm.messageType)
	msg.SenderID = qm.senderID
}
