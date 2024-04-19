package service

import (
	"errors"
	"sync"
	"time"
)

const (
	minDurationAfterLastJoin = time.Second
)

var (
	errNotExist = errors.New("this chat does not exit in the system")
)

type chats struct {
	chats map[string]*chatInfo
	m     sync.RWMutex
}

type chatInfo struct {
	isJoined     bool
	lastJoinTry  time.Time
	joinCancelFn func(error)
}

func newChats() chats {
	return chats{
		chats: make(map[string]*chatInfo),
	}
}

func (c *chats) updateToJoined(chat string) error {
	c.m.Lock()
	defer c.m.Unlock()

	info, isExist := c.chats[chat]
	if !isExist {
		return errNotExist
	}

	info.isJoined = true
	info.joinCancelFn(nil)

	return nil
}

func (c *chats) processJoinRequest(chat string, cancelFn func(error)) error {
	c.m.Lock()
	defer c.m.Unlock()

	prevInfo, isExist := c.chats[chat]
	if !isExist {
		c.chats[chat] = &chatInfo{
			isJoined:     false,
			lastJoinTry:  time.Now(),
			joinCancelFn: cancelFn,
		}
		return nil
	}

	if err := validateExistingRecord(prevInfo); err != nil {
		return err
	}

	prevInfo.lastJoinTry = time.Now()
	return nil
}

func validateExistingRecord(info *chatInfo) error {
	if info.isJoined {
		return errors.New("already connected to the chat")
	}

	nextRequestThreshold := info.lastJoinTry.Add(minDurationAfterLastJoin)
	if time.Now().Before(nextRequestThreshold) {
		return errors.New("to many request, please try later")
	}

	return nil
}
