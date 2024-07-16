package daos

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gm/global"
	"gm/middlewares"
	"gm/model"
	"gm/model/gm"
	"gm/request"
	"gm/response"
	"gm/utils/ecode"
	"image"
	"image/png"
	"strings"
	"time"
)

// 登录
func Login(req request.Login) (res response.Login, err error) {
	db := global.DB
	obj := new(gm.User)

	user, err := obj.GetUserByParams(db, req.Id, req.Username, "id,username,name,password,is_bind")

	if err != nil {
		return
	}

	options := &password.Options{16, 100, 32, sha512.New}

	pwdSlice := strings.Split(user.Password, "$")

	if !password.Verify(req.Password, pwdSlice[1], pwdSlice[2], options) {
		err = ecode.PasswordCheckFailed
		return
	}

	//if !CheckGoogleCaptcha(user.Id, req.Captcha) {
	//	err = ecode.CaptchaVerifyFailed
	//
	//	global.Logger["err"].Errorf("captcha verify failed")
	//	return
	//}

	channelIds := make([]int, 0)

	channelIds, err = GetAdminUserChannelIds(user.Id)
	if err != nil {
		return
	}

	hour := time.Duration(24) * 3

	claims := middlewares.CustomClaims{
		ID:             user.Id,
		Username:       user.Username,
		Name:           user.Name,
		RoleChannelIds: channelIds,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(hour * time.Hour).Unix()},
	}

	j := middlewares.NewJwt()
	token, err := j.CreateToken(claims)
	if err != nil {
		return
	}

	res = response.Login{
		UserId:   user.Id,
		Username: user.Username,
		Name:     user.Name,
		Token:    token,
		IsBind:   user.IsBind,
	}

	//获取用户权限菜单列表
	tree, err := UserRoleMenuTree(user.Id)

	if err != nil {
		return
	}

	res.Menu = tree

	//获取用户权限菜单列表
	mp, err := UserRoleMenuMp(user.Id)

	if err != nil {
		return
	}

	res.MenuMp = mp

	return
}

// 新增用户
func CreateUser(form request.AddUser) (err error) {
	db := global.DB
	userObj := new(gm.User)
	createUser, err := userObj.GetUserByParams(db, form.CreatorId, "", "id,name")

	if err != nil {
		return
	}

	user := gm.User{
		Base: model.Base{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: form.Username,
		Name:     form.Name,
		Phone:    form.Phone,
		Email:    form.Email,
		Password: GenderPwd(form.Password),
		Status:   form.Status,
		Creator: model.Creator{
			CreatorId: createUser.Id,
			Creator:   createUser.Name,
		},
	}

	err = db.Model(&gm.User{}).Save(&user).Error

	if err != nil {
		return
	}

	roleObj := new(gm.UserRole)

	err = roleObj.DelAndCreateRoleMenu(db, user.Id, form.RoleList)
	if err != nil {
		return
	}
	return
}

// ChangeUser 修改用户信息
func ChangeUser(req request.ChangeUser) (err error) {

	db := global.DB
	user := new(gm.User)
	mp := make(map[string]interface{}, 0)

	//更新用户密码
	if req.Password != "" {
		mp["password"] = GenderPwd(req.Password)
	}

	//更新用户名称
	if req.Name != "" {
		mp["name"] = req.Name
	}

	//更新账号
	if req.Username != "" {
		mp["username"] = req.Username
	}

	//更新邮箱
	if req.Email != "" {
		mp["email"] = req.Email
	}

	//更新手机号
	if req.Phone != "" {
		mp["phone"] = req.Phone
	}

	//更新用户状态
	if req.Status > 0 {
		mp["status"] = req.Status
	}

	err = user.UpdateById(db, req.Id, mp)

	if err != nil {
		return
	}

	roleObj := new(gm.UserRole)

	err = roleObj.DelAndCreateRoleMenu(db, req.Id, req.RoleList)
	if err != nil {
		return
	}

	return
}

// 获取用户列表
func UserList(req request.UserList) (result response.UserList, err error) {
	db := global.DB
	user := new(gm.User)

	var (
		users []gm.User
		total int64
	)

	localDb := db

	if req.Id > 0 {
		localDb.Where("id = ?", req.Id)
	}

	if req.Username != "" {
		localDb.Where("username like ?", "%"+req.Username+"%")
	}

	if req.Name != "" {
		localDb.Where("name like ?", "%"+req.Name+"%")
	}

	if req.Status > 0 {
		localDb.Where("status = ?", req.Status)
	}

	total, err = user.Count(localDb)

	if err != nil {
		return
	}

	result.Total = total

	users, err = user.GetPageList(localDb, req)

	if err != nil {
		return
	}

	userIds := make([]int, 0, len(users))

	for _, u := range users {
		userIds = append(userIds, u.Id)
	}

	//用户角色列表
	roleObj := new(gm.UserRole)

	mp := make(map[int][]int)

	mp, err = roleObj.MapOfRoleIdsByUserId(db, userIds)

	if err != nil {
		return
	}

	list := make([]response.User, 0, len(users))

	for _, u := range users {
		tmp := response.User{
			Id:        u.Id,
			Username:  u.Username,
			Name:      u.Name,
			Phone:     u.Phone,
			Email:     u.Email,
			Status:    u.Status,
			CreatorId: u.CreatorId,
			Creator:   u.Creator.Creator,
			Secret:    u.Secret,
			IsBind:    u.IsBind,
		}

		roleList, ok := mp[u.Id]

		if !ok {
			roleList = make([]int, 0)
		}

		tmp.RoleList = roleList

		list = append(list, tmp)
	}

	result.List = list

	return
}

// 获取用户权限菜单列表
func UserRoleMenuTree(uid int) (tree []response.Group, err error) {
	db := global.DB

	userRole := new(gm.UserRole)

	//用户所有角色id
	roleIds, err := userRole.GetRoleIdsByUserId(db, uid)
	if err != nil {
		return
	}

	//角色菜单ids
	roleMenu := new(gm.RoleMenu)
	menuIds, err := roleMenu.GetMenuIdsByRoleIds(db, roleIds)
	if err != nil {
		return
	}

	var menuList []gm.Menu

	mp := make(map[int]gm.Menu)

	mp, err = MapOfAllMenuList()

	if err != nil {
		return
	}

	visited := make(map[int]bool)

	CollectMenus(mp, menuIds, &menuList, visited)

	menuGroup := make([]response.Group, 0, len(menuList))

	for _, m := range menuList {
		menuGroup = append(menuGroup, response.Group{
			Id:       m.Id,
			Name:     m.Name,
			Path:     m.Path,
			Label:    m.Label,
			Icon:     m.Icon,
			Url:      m.Url,
			ParentId: m.ParentId,
			Children: []response.Group{},
		})
	}

	tree = buildTree(menuGroup, 0)

	return
}

func CollectMenus(mp map[int]gm.Menu, ids []int, list *[]gm.Menu, visited map[int]bool) {
	for _, id := range ids {
		if menu, ok := mp[id]; ok && !visited[id] {
			visited[id] = true // 标记为已访问
			*list = append(*list, menu)

			// 递归查找父级
			if menu.ParentId != 0 {
				CollectMenus(mp, []int{menu.ParentId}, list, visited)
			}
		}
	}
}

func UserRoleMenuMp(uid int) (mp map[string]struct{}, err error) {
	db := global.DB

	userRole := new(gm.UserRole)

	//用户所有角色id
	roleIds, err := userRole.GetRoleIdsByUserId(db, uid)
	if err != nil {
		return
	}

	//角色菜单ids
	roleMenu := new(gm.RoleMenu)
	menuIds, err := roleMenu.GetMenuIdsByRoleIds(db, roleIds)
	if err != nil {
		return
	}

	menu := new(gm.Menu)

	menuList, err := menu.GetListByIds(db, menuIds)

	if err != nil {
		return
	}

	mp = make(map[string]struct{})
	for _, m := range menuList {
		mp[m.Url] = struct{}{}
	}

	return
}

func GetAdminUserChannelIds(userId int) (channelIds []int, err error) {
	db := global.DB

	//查询用户角色
	userRoleObj := new(gm.UserRole)

	roleIds := make([]int, 0)

	roleIds, err = userRoleObj.GetRoleIdsByUserId(db, userId)
	if err != nil {
		global.Logger["err"].Infof("查询用户角色失败:" + err.Error())
		return
	}
	//查询角色渠道
	roleChanObj := new(gm.RoleChannel)

	channelIds, err = roleChanObj.GetChannelIdsByRoleIds(db, roleIds)

	if err != nil {
		global.Logger["err"].Infof("查询角色渠道数据失败:" + err.Error())
		return
	}

	return
}

// GenderPwd 生成密码
func GenderPwd(pwd string) string {
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(pwd, options)
	return fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
}

// buildTree 生成树形结构
func buildTree(groups []response.Group, parentId int) []response.Group {
	var tree []response.Group
	for _, group := range groups {
		if group.ParentId == parentId {
			children := buildTree(groups, group.Id)
			if children != nil {
				group.Children = children
			}
			tree = append(tree, group)
		}
	}
	return tree
}

func UpdateUserGoogleCaptcha(req request.Uid) (base64Img string, err error) {
	var (
		img image.Image
		g   *otp.Key
	)

	// 生成一个随机的密钥
	g, err = totp.Generate(totp.GenerateOpts{
		Issuer:      "zw",
		AccountName: "zheanhuyu@gmail.com",
		Period:      60,
	})

	if err != nil {
		global.Logger["err"].Errorf("totp.Generate key failed,err:[%v]", err.Error())
		return
	}

	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err = g.Image(200, 200)
	if err != nil {
		global.Logger["err"].Errorf("g.Image failed,err:[%v]", err.Error())
		return
	}

	err = png.Encode(&buf, img)
	if err != nil {
		global.Logger["err"].Errorf("png.Encode failed,err:[%v]", err.Error())
		return
	}

	// 将bytes.Buffer中的数据编码为Base64字符串
	pngBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// 构造data:image/png;base64,开头的字符串
	base64Img = fmt.Sprintf("data:image/png;base64,%s", pngBase64)

	secret := g.Secret()

	db := global.DB

	obj := new(gm.User)

	err = obj.UpdateById(db, req.Uid, map[string]interface{}{"secret": secret, "is_bind": gm.IsBindNo})

	if err != nil {
		global.Logger["err"].Errorf("update user secret failed,err:[%v]", err.Error())
		return
	}

	return
}

func GetCaptchaQrBySecret(req request.GetCaptchaQr) (base64Img string, err error) {
	var (
		g   *otp.Key
		img image.Image
		u   gm.User
	)

	if req.Secret == "" {
		db := global.DB

		obj := new(gm.User)

		u, err = obj.GetOneById(db, req.Uid)

		if err != nil {
			global.Logger["err"].Errorf("GetCaptchaQrBySecret obj.GetOneById failed,err:[%v]", err.Error())
			return
		}

		req.Secret = u.Secret
	}

	g, err = otp.NewKeyFromURL(fmt.Sprintf(`otpauth://totp/zheanhuyu@gmail.com?secret=%s&issuer=zw`, req.Secret))

	if err != nil {
		global.Logger["err"].Errorf("totp.Generate key failed,err:[%v]", err.Error())
		return
	}

	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err = g.Image(200, 200)
	if err != nil {
		global.Logger["err"].Errorf("g.Image failed,err:[%v]", err.Error())
		return
	}

	err = png.Encode(&buf, img)
	if err != nil {
		global.Logger["err"].Errorf("png.Encode failed,err:[%v]", err.Error())
		return
	}

	// 将bytes.Buffer中的数据编码为Base64字符串
	pngBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// 构造data:image/png;base64,开头的字符串
	base64Img = fmt.Sprintf("data:image/png;base64,%s", pngBase64)
	return
}

func CheckGoogleCaptcha(uid int, code string) bool {
	db := global.DB

	obj := new(gm.User)

	user, err := obj.GetOneById(db, uid)
	if err != nil {
		global.Logger["err"].Errorf("select user failed,err:[%s]", err.Error())

		return false
	}

	return totp.Validate(code, user.Secret)
}

func BindGoogleCaptcha(req request.CheckGoogleCaptcha) (err error) {
	db := global.DB

	obj := new(gm.User)

	user, err := obj.GetOneById(db, req.Uid)
	if err != nil {
		global.Logger["err"].Errorf("select user failed,err:[%s]", err.Error())

		return
	}

	ok := totp.Validate(req.Code, user.Secret)

	if !ok {
		err = errors.New("google captcha Validate failed")
		global.Logger["err"].Errorf("totp.Validate failed,err:[%s]", err.Error())
		return
	}

	err = obj.UpdateById(db, user.Id, map[string]interface{}{"is_bind": gm.IsBindYes})
	if err != nil {
		global.Logger["err"].Errorf("obj.UpdateById failed,err:[%s]", err.Error())
		return
	}

	return
}
