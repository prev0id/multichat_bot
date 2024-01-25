package service

import (
	"errors"
	"sync"
	"time"
)

const (
	minDurationAfterLastJoin = time.Second
)

//var (
//	errNotExist = errors.New("this chat does not exit in the system")
//)

type chats struct {
	chats map[string]*chatInfo
	m     sync.RWMutex
}

type chatInfo struct {
	isJoined     bool
	latJoinTry   time.Time
	joinCancelFn func(error)
}

func newChats() chats {
	return chats{
		chats: make(map[string]*chatInfo),
	}
}

//
//func (c *chats) updateToConnected(chat string) error {
//	c.m.Lock()
//	defer c.m.Unlock()
//
//	info, isExist := c.chats[chat]
//	if !isExist {
//		return errNotExist
//	}
//
//	info.isJoined = true
//	return nil
//}
//
//func (c *chats) updateToDisconnected(chat string) error {
//	c.m.Lock()
//	defer c.m.Unlock()
//
//	info, isExist := c.chats[chat]
//	if !isExist {
//		return errNotExist
//	}
//
//	info.isJoined = false
//	info.disconnectedAt = time.Now()
//	return nil
//}
//
//func (c *chats) cancelJoinRequest(chat string, err error) error {
//	c.m.Lock()
//	defer c.m.Unlock()
//
//	info, isExist := c.chats[chat]
//	if !isExist {
//		return errNotExist
//	}
//
//	info.joinCancelFn(err)
//	return nil
//}

func (c *chats) processJoinRequest(chat string, cancelFn func(error)) error {
	c.m.Lock()
	defer c.m.Unlock()

	prevInfo, isExist := c.chats[chat]
	if !isExist {
		c.chats[chat] = &chatInfo{
			isJoined:     false,
			latJoinTry:   time.Now(),
			joinCancelFn: cancelFn,
		}
		return nil
	}

	if err := validateExistingRecord(prevInfo); err != nil {
		return err
	}

	prevInfo.latJoinTry = time.Now()

	return nil
}

func validateExistingRecord(info *chatInfo) error {
	if info.isJoined {
		return errors.New("already connected to the chat")
	}

	nextRequestThreshold := info.latJoinTry.Add(minDurationAfterLastJoin)
	if time.Now().Before(nextRequestThreshold) {
		return errors.New("to many request, please try later")
	}

	return nil
}
