package bbmessagequeue

import (
	"sync"

	"github.com/gunnarsutter/blackboard"
)

// Queue of messages in internal format
type MessageQueue struct {
	Mutex sync.Mutex
	Queue []QueueMessage
}

// Message in internal format
type QueueMessage struct {
	SenderID    string
	MessageType int
	Content     string
}

// Sets values of the struct varables
func (qm *QueueMessage) Init(s_id string, m_t int, c string) {
	qm.SenderID = s_id
	qm.MessageType = m_t
	qm.Content = c
}

// Populates an internal message using a blackboard message
func (qm *QueueMessage) New(message *blackboard.Message) {
	qm.SenderID = message.SenderID
	qm.MessageType = int(message.MessageType)
	qm.Content = message.Content
}

// Adds a blackboard message to the end of the internal message queue
func (mq *MessageQueue) Push(m *blackboard.Message) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	var qm QueueMessage
	qm.New(m)
	mq.Queue = append(mq.Queue, qm)
}

// Adds an internal message to the end of the internal message queue
func (mq *MessageQueue) PushQM(qm *QueueMessage) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	mq.Queue = append(mq.Queue, *qm)
}

// Adds a blackboard message to the front of the internal message queue
func (mq *MessageQueue) PushFront(m *blackboard.Message) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	var qm QueueMessage
	qm.New(m)
	mq.Queue = append([]QueueMessage{qm}, mq.Queue...)
}

// Adds an internal message to the front of the internal message queue
func (mq *MessageQueue) PushFrontQM(qm *QueueMessage) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	mq.Queue = append([]QueueMessage{*qm}, mq.Queue...)
}

// Returns and removes an internal message from the front of the queue
func (mq *MessageQueue) Pop() *QueueMessage {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	ret_message := mq.Queue[0]
	mq.Queue = mq.Queue[1:]
	return &ret_message
}

// Returns the length of the internal message queue
func (mq *MessageQueue) Length() int {
	return len(mq.Queue)
}

// Populates a blackboard message from an internal message
func (qm *QueueMessage) ToBlackboardMessage(msg *blackboard.Message) {
	msg.Content = qm.Content //DETTA CRASHAR!!! Fortfarande?
	msg.MessageType = blackboard.MessageType(qm.MessageType)
	msg.SenderID = qm.SenderID
}
