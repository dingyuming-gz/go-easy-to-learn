package main

import (
	"context"
	"errors"
	"time"

	"github.com/reactivex/rxgo/v2"
	"github.com/sirupsen/logrus"
)

// ----------- 定义任务与用户的抽象

type Task struct {
	userID int
}

// GetUserID 获取这个task关联的用户ID
func (t *Task) GetUserID() int {
	return t.userID
}

type User struct {
	id   int
	role string
}

// GetID 获取用户的ID
func (u *User) GetID() int {
	return u.id
}

// GetRole 获取用户的角色
func (u *User) GetRole() string {
	return u.role
}

// ----------- 定义鉴权系统的入参

// checkInput 传入鉴权系统的入参，它是一个future，使用GetRet等待系统对其进行鉴权
type checkInput struct {
	t  *Task
	u  *User
	ch chan bool
}

// GetRet 阻塞并等待鉴权的结果
func (r *checkInput) GetRet() bool {
	select {
	case v := <-r.ch:
		return v
	case <-time.After(time.Second):
		return false
	}
}

// Error 这里利用rxgo return error提前结束流程的特性，这个参数可以作为error提前返回
func (checkInput) Error() string {
	return ""
}

func newReqStruct(u *User, t *Task) checkInput {
	return checkInput{
		u:  u,
		t:  t,
		ch: make(chan bool),
	}
}

// ----------- 把各种鉴权操作抽象为函数

func CheckTaskUser(task interface{ GetUserID() int }, user interface{ GetID() int }) bool {
	return task.GetUserID() == user.GetID()
}

func IsAdmin(user interface{ GetRole() string }) bool {
	return user.GetRole() == "admin"
}

// ----------- 定义鉴权流程，流程持续运行

func Check(req checkInput) bool {
	checkCh <- req
	return req.GetRet()
}

var checkCh = make(chan checkInput, 1)

func init() {
	// 这里只是演示所以硬编码，实际上每个节点是通过配置文件加载的，流程可以动态变更与配置
	rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		for req := range checkCh {
			next <- rxgo.Of(req)
		}
	}}).
		Map(func(c context.Context, i interface{}) (interface{}, error) {
			// 第1个节点，检查是否是管理员
			req := i.(checkInput)
			if IsAdmin(req.u) {
				// 如果是管理员，使用return error提前终止流程
				return nil, req
			}
			return req, nil
		}).
		Map(func(c context.Context, i interface{}) (interface{}, error) {
			// 第2个节点，检查是否是自己的任务
			req := i.(checkInput)
			if CheckTaskUser(req.t, req.u) {
				// 如果是自己的任务，使用return error提前终止流程
				return nil, req
			}
			return req, nil
		}).
		ForEach(func(i interface{}) {
			v := i.(checkInput)
			// 如果返回值仍然不为空，则说明所有流程都执行完毕，检查失败
			if i != nil {
				v.ch <- false
			}
		}, func(e error) {
			req := checkInput{}
			// 如果错误类型是req类型，说明是提前返回的中断流程信号，对此次检查设置为成功
			if errors.As(e, &req) {
				req.ch <- true
			}
			// 因为要持续运行，发现中断错误则继续运行
		}, func() {}, rxgo.WithErrorStrategy(rxgo.ContinueOnError))
}

// ----------- 开始鉴权
func main() {
	task := &Task{
		userID: 1,
	}
	task2 := &Task{
		userID: 2,
	}
	adminUser := &User{
		id:   1,
		role: "admin",
	}
	otherUser := &User{
		id:   2,
		role: "normal",
	}
	// 如果是管理员，且是自己的任务，通过
	t1 := Check(newReqStruct(adminUser, task))
	logrus.Infof("t1:%v", t1) // true

	// 如果不是管理员，且不是自己的任务，不通过
	t2 := Check(newReqStruct(otherUser, task)) // false

	logrus.Infof("t2:%v", t2)

	// 如果是管理员，且不是自己的任务，通过
	t3 := Check(newReqStruct(adminUser, task2))
	logrus.Infof("t3:%v", t3) // true
}
