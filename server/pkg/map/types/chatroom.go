package types

import "sync"

type ChatRoom struct {
	Id       string
	Name     string
	Admin    User
	Chatters *users
}

func NewChatRoom(id, name string, admin User) *ChatRoom {
	cr := &ChatRoom{
		Id:       id,
		Name:     name,
		Admin:    admin,
		Chatters: newUsers(),
	}

	cr.AddChatter(admin)

	return cr
}

func (cr *ChatRoom) AddChatter(u User) {
	cr.Chatters.add(u)
}

func (cr *ChatRoom) RemoveChatter(u User) {
	cr.Chatters.remove(u)
}

func (cr *ChatRoom) GetChatters() []User {
	return cr.Chatters.users
}

func (cr *ChatRoom) UserLen() int {
	return cr.Chatters.len()
}

type users struct {
	lock  sync.RWMutex
	users []User
}

func newUsers() *users {
	return &users{
		users: make([]User, 0),
	}
}

func (u *users) len() int {
	return len(u.users)
}

func (u *users) add(user User) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.users = append(u.users, user)
}

func (u *users) remove(user User) {
	u.lock.Lock()
	defer u.lock.Unlock()
	i := u.find(user)
	if i == -1 {
		return
	}
	u.users = append(u.users[:i], u.users[i+1:]...)
}

func (u *users) find(user User) int {
	for i, v := range u.users {
		if user.Id == v.Id {
			return i
		}
	}
	return -1
}
