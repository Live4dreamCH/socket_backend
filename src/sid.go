package main

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

// sessionID数据结构及管理
type sID struct {
	m map[string]int
	l sync.RWMutex
}

// 添加一个用户
// 若已在线, 返回原有sid
// 否则, 生成新sid, 添加并返回
func (s *sID) set(uid int) string {
	s.l.RLock()
	for sid, u := range s.m {
		if u == uid {
			s.l.RUnlock()
			return sid
		}
	}

	sid := newSID()
	_, ok := sess.m[sid]
	for ok {
		sid = newSID()
		_, ok = sess.m[sid]
	}
	s.l.RUnlock()
	s.l.Lock()
	sess.m[sid] = uid
	s.l.Unlock()
	return sid
}

// 用sid查询登录状态
// 已登录, 返回uid, nil
// 未登录, 返回0, err
func (s *sID) get(sid string) (int, error) {
	return 0, nil
}

// 随机生成一个sid
func newSID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
