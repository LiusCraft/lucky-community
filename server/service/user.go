package services

import (
	"errors"

	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

type UserService struct {
}

// get user information
func (*UserService) GetUserById(id int) *model.Users {

	user := userDao.QueryUser(&model.Users{ID: id})
	user.Avatar = utils.BuildFileUrl(user.Avatar)
	user.Password = ""
	user.InviteCode = 0
	return user
}

// get user information
func (*UserService) GetUserSimpleById(id int) *model.UserSimple {
	user := userDao.QueryUser(&model.Users{ID: id})
	role := "none"
	model.InviteCode().Where("code = ?", user.InviteCode).
		Joins("JOIN member_infos ON member_infos.id = invite_codes.member_id").
		Select("member_infos.name").Limit(1).Scan(&role)
	return &model.UserSimple{
		UId:     user.ID,
		UName:   user.Name,
		UDesc:   user.Desc,
		UAvatar: utils.BuildFileUrl(user.Avatar),
		Role:    role,
		Account: user.Account,
	}
}

// update user information
func (*UserService) UpdateUser(user *model.Users) {

	userDao.UpdateUser(user)

}

func (*UserService) ListByIdsSelectEmail(id ...int) []string {
	return userDao.ListByIds(id...)
}

func (s *UserService) ListByIdsToMap(ids []int) map[int]model.Users {

	m := make(map[int]model.Users)
	users := userDao.ListByIdsSelectIdName(ids)
	for i := range users {
		user := users[i]
		m[user.ID] = user
	}
	return m
}

func Login(account, pswd string) (*model.Users, error) {

	user := userDao.QueryUser(&model.Users{Account: account, Password: pswd})
	if user.ID == 0 {
		return nil, errors.New("登录失败！账号或密码错误")
	}

	user.Avatar = utils.BuildFileUrl(user.Avatar)
	return user, nil
}

func Register(account, pswd, name string, inviteCode int) error {

	if err := utils.NotBlank(account, pswd, name, inviteCode); err != nil {
		return err
	}

	// query codeDao
	if !codeDao.Exist(inviteCode) {
		return errors.New("验证码不存在")
	}

	// 查询账户
	user := userDao.QueryUser(&model.Users{Account: account})
	if user.ID > 0 {
		return errors.New("账户已存在,换一个吧")
	}

	// 保存用户
	userDao.CreateUser(account, name, pswd, inviteCode)
	// 修改code状态
	var c CodeService
	c.SetState(inviteCode)

	return nil
}

type UserMenu struct {
	Path       string                 `json:"path"`
	Redirect   string                 `json:"redirect,omitempty"`
	Name       string                 `json:"name"`
	Components string                 `json:"components,omitempty"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
	Children   []UserMenu             `json:"children,omitempty"`
}

func (t *UserService) GetUserMenu() []*UserMenu {
	rootMenu := typeDao.List(0)
	parentIds := make([]int, len(rootMenu))
	userMenu := make(map[int]*UserMenu)
	for i, item := range rootMenu {
		parentIds[i] = item.ID
		userMenu[int(item.ID)] = &UserMenu{
			Path:     "/article/" + item.FlagName,
			Name:     item.FlagName,
			Children: []UserMenu{},
			Meta: map[string]interface{}{
				"locale":       item.Title,
				"requiresAuth": true,
				"icon":         "icon-dashboard",
				"order":        1,
			},
		}
	}
	children := []model.Types{}
	model.Type().Where("parent_id in (?)", parentIds).Find(&children)
	for _, item := range children {
		um := userMenu[int(item.ParentId)]
		um.Children = append(um.Children, UserMenu{
			Path: "/article/" + item.FlagName,
			Name: item.FlagName,
			Meta: map[string]interface{}{
				"locale":       item.Title,
				"requiresAuth": true,
				"icon":         "icon-dashboard",
			},
		})
	}
	result := []*UserMenu{}

	for _, um := range userMenu {
		result = append(result, um)
	}
	return result
}

func (s *UserService) CheckCodeUsed(code int) bool {
	var count int64
	model.User().Where("invite_code = ?", code).Count(&count)
	return count == 1
}

func (s *UserService) PageUsers(p, limit int) (users []model.Users, count int64) {
	model.User().Offset(limit).Limit((p - 1) * limit).Find(&users)
	model.User().Count(&count)
	return users, count
}

func (s *UserService) ListByNameSelectEmailAndId(usernames []string) (emails []string, id []int) {
	var users []model.Users
	model.User().Where("name in ? ", usernames).Select("account", "id").Find(&users)

	for i := range users {
		u := users[i]
		emails = append(emails, u.Account)
		id = append(id, u.ID)
	}

	return emails, id
}

func (s *UserService) ListBySelect(user model.Users) (users []model.Users) {
	model.User().Where(user).Find(&users)
	return
}

func (s *UserService) Statistics(userId int) (m map[string]interface{}) {
	m = make(map[string]interface{})
	var articleS ArticleService
	// 获取被点赞次数,获取用户发布的所有文章
	ids := articleS.PublishArticlesSelectId(userId)
	likeCount := articleS.ArticlesLikeCount(ids)
	// 获取发布文章
	count := len(ids)

	m["articleCount"] = count
	m["likeCount"] = likeCount
	return
}
