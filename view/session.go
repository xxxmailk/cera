package view

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	SessionExpired   = errors.New("session has been expired, logout")
	SessionNoThisKey = errors.New("has not key in session data")
)

// stored sessions
var sessions = make(map[[16]byte]*Session)

type Session struct {
	id        [16]byte               // session uuid
	data      map[string]interface{} // session data of this user
	expire    time.Time              // expire time of this session
	zeroTime  time.Time              // new session time
	lifeCycle int64                  // life cycle of session, unit: minute
}

// update session start time
func (s *Session) updateTime() {
	s.zeroTime = time.Now()
	s.expire = s.zeroTime.Add(time.Duration(s.lifeCycle) * time.Minute)
}

func (s *Session) isExpired() bool {
	if s.expire.After(time.Now()) {
		return true
	}
	return false
}

// set session data
func (s *Session) Set(key string, val interface{}) {
	s.data[key] = val
}

// check has key in session data or not
func (s *Session) HasKey(key string) bool {
	if _, ok := s.data[key]; ok {
		return true
	}
	return false
}

// get session , Get() return value and error.
// if time expire, Get() will return timeout error
func (s *Session) Get(key string) (val interface{}, err error) {
	if s.isExpired() {
		return nil, SessionExpired
	}
	if !s.HasKey(key) {
		return nil, SessionNoThisKey
	}
	return s.data[key], nil
}

// create new session, user id could not be repeat
func CreateSessionId(username string, userId interface{}) ([16]byte, error) {
	var buf []byte
	var b = new(bytes.Buffer)
	buf = append(buf, []byte(username)...)
	gb := gob.NewEncoder(b)
	err := gb.Encode(userId)
	if err != nil {
		return [16]byte{}, err
	}
	return uuid.Must(uuid.FromBytes(buf)), nil
}

func NewSession(username string, userId interface{}, lifeCycle int64) (*Session, error) {
	id, err := CreateSessionId(username, userId)
	if err != nil {
		return nil, err
	}
	s := new(Session)
	s.data = make(map[string]interface{})
	s.zeroTime = time.Now()
	s.expire = s.zeroTime.Add(time.Duration(lifeCycle) * time.Minute)
	s.lifeCycle = lifeCycle
	s.id = id

	// add session to session map
	sessions[id] = s
	return s, nil
}
