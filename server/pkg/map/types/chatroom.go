package types

import "sync"

// ChatRoom represents a chat room with an ID, name, admin, and chatters
type ChatRoom struct {
    Id       string
    Name     string
    Admin    *User
    Chatters *users
}

// NewChatRoom creates a new chat room with the given ID, name, and admin
func NewChatRoom(id, name string, admin *User) *ChatRoom {
    cr := &ChatRoom{
        Id:       id,
        Name:     name,
        Admin:    admin,
        Chatters: newUsers(),
    }

    cr.AddChatter(admin)

    return cr
}

// AddChatter adds a user to the chat room
func (cr *ChatRoom) AddChatter(u *User) {
    cr.Chatters.add(u)
}

// RemoveChatter removes a user from the chat room and returns true if the room is empty
func (cr *ChatRoom) RemoveChatter(u *User) bool {
    cr.Chatters.remove(u)
    return cr.UserLen() == 0
}

// GetChatters returns the list of users in the chat room
func (cr *ChatRoom) GetChatters() []*User {
    return cr.Chatters.users
}

// UserLen returns the number of users in the chat room
func (cr *ChatRoom) UserLen() int {
    return cr.Chatters.len()
}

// users represents a thread-safe list of users
type users struct {
    lock  sync.RWMutex
    users []*User
}


// newUsers creates a new users instance
func newUsers() *users {
    return &users{
        users: make([]*User, 0),
    }
}

// len returns the number of users
func (u *users) len() int {
    return len(u.users)
}

// add adds a user to the users list
func (u *users) add(user *User) {
    u.lock.Lock()
    defer u.lock.Unlock()
    u.users = append(u.users, user)
}

// remove removes a user from the users list
func (u *users) remove(user *User) {
    u.lock.Lock()
    defer u.lock.Unlock()
    i := u.find(user)
    if i == -1 {
        return
    }
    u.users = append(u.users[:i], u.users[i+1:]...)
}

// find finds the index of a user in the users list
func (u *users) find(user *User) int {
    for i, v := range u.users {
        if user.Id == v.Id {
            return i
        }
    }
    return -1
}