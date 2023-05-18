// @Author: zqy
// @File: manager.go
// @Date: 2023/5/18 14:37
// @Description todo

package Session

import (
	"GoCode/web/customize"
)

// Manager 提高用户体验
type Manager struct {
	Store
	Propagator
	CtxSessionKey string
}

// init get remove refresh

func (m *Manager) GetSession(ctx *customize.Context) (Session, error) {
	// 作为缓存功能 聊胜于无，收益较小
	// 也可以使用context进行传递，但是涉及到拷贝性能可能较差，withContext
	// 还有问题就是父context无法访问到子context
	//ctx.Req.Context().Value()
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	// 先从缓存里面进行读取
	ses, ok := ctx.UserValues[m.CtxSessionKey]
	if ok {
		return ses.(Session), nil
	}
	sessionID, err := m.Propagator.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	session, err := m.Store.Get(ctx.Req.Context(), sessionID)
	ctx.UserValues[m.CtxSessionKey] = session
	return session, err
}

// InitSession 初始化session，使用uuid作为key
func (m *Manager) InitSession(ctx *customize.Context, id string) (Session, error) {
	session, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	// 注册结束后 注入到resp里
	if err = m.Inject(id, ctx.Resp); err != nil {
		return nil, err
	}
	return session, nil
}

func (m *Manager) RemoveSession(ctx *customize.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), session.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}

// RefreshSession 更新session过期时间
func (m *Manager) RefreshSession(ctx *customize.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Refresh(ctx.Req.Context(), session.ID())
	return err
}
