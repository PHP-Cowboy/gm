package daos

import (
	"gm/global"
	"gm/model"
	"gm/model/gm"
	"gm/request"
	"gm/response"
)

func CreateMenu(req request.AddMenu) (err error) {
	db := global.DB

	menu := &gm.Menu{
		Name:     req.Name,
		Path:     req.Path,
		Label:    req.Label,
		Icon:     req.Icon,
		Url:      req.Url,
		ParentId: *req.ParentId,
		Level:    req.Level,
		Creator: model.Creator{
			CreatorId: req.CreatorId,
			Creator:   req.Creator,
		},
	}

	err = menu.Save(db)

	if err != nil {
		return
	}
	return
}

// 修改菜单信息
func ChangeMenu(req request.ChangeMenu) (err error) {
	db := global.DB

	mp := make(map[string]interface{}, 0)
	//更新菜单描述
	if req.Desc != "" {
		mp["desc"] = req.Desc
	}

	//更新菜单名称
	if req.Name != "" {
		mp["name"] = req.Name
	}

	//菜单路由
	if req.Path != "" {
		mp["path"] = req.Path
	}

	//label
	if req.Label != "" {
		mp["label"] = req.Label
	}

	//前端组件url
	if req.Url != "" {
		mp["url"] = req.Url
	}

	//更新菜单状态
	if req.Status > 0 {
		mp["status"] = req.Status
	}

	if req.Level > 0 {
		mp["level"] = req.Level
	}

	if req.ParentId != nil {
		mp["parent_id"] = *req.ParentId
	}

	obj := new(gm.Menu)

	return obj.UpdateMenuById(db, req.Id, mp)
}

// all
func MapOfAllMenuList() (mp map[int]gm.Menu, err error) {
	db := global.DB

	menuObj := new(gm.Menu)

	var dataList []gm.Menu

	dataList, err = menuObj.GetList(db)
	if err != nil {
		global.Logger["err"].Infof("查询全部菜单列表失败：" + err.Error())
		return
	}

	mp = make(map[int]gm.Menu)

	for _, d := range dataList {
		mp[d.Id] = d
	}

	return
}

// 获取角色列表
func MenuList(req request.MenuList) (result response.MenuList, err error) {
	db := global.DB.Model(&gm.Menu{})

	var (
		menus []gm.Menu
	)

	if req.Name != "" {
		db.Where("name like ?", "%"+req.Name+"%")
	}

	if req.Status > 0 {
		db.Where("status = ?", req.Status)
	}

	var total int64

	err = db.Count(&total).Error

	if err != nil {
		return
	}

	result.Total = total

	err = db.Find(&menus).Error

	if err != nil {
		return
	}

	list := make([]response.Menu, 0, len(menus))

	for _, m := range menus {
		list = append(list, response.Menu{
			Id:            m.Id,
			Name:          m.Name,
			ComponentName: m.ComponentName,
			Path:          m.Path,
			Label:         m.Label,
			Icon:          m.Icon,
			Url:           m.Url,
			Status:        m.Status,
			ParentId:      m.ParentId,
			Level:         m.Level,
		})
	}

	result.List = list

	return
}

// MenuTree
func MenuTree() (tree []response.Group, err error) {
	db := global.DB

	menu := new(gm.Menu)

	menuList, err := menu.GetList(db)
	if err != nil {
		return
	}

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
			Value:    m.Id,
		})
	}

	tree = buildTree(menuGroup, 0)

	return
}

// 上级菜单列表
func LevelList(req request.LevelList) (list []response.ParentList, err error) {
	db := global.DB

	menuObj := new(gm.Menu)

	dataList, err := menuObj.GetListByLevel(db, req.Level)
	if err != nil {
		return
	}

	list = make([]response.ParentList, 0, len(dataList))

	for _, menu := range dataList {
		list = append(list, response.ParentList{
			Id:    menu.Id,
			Label: menu.Label,
		})
	}

	return
}

// 获取用户角色列表
func RoleMenuList(req request.RoleMenuList) (list []response.Menu, err error) {
	db := global.DB

	var (
		roleMenus []gm.RoleMenu
		menus     []gm.Menu
		menuIds   []int
	)

	err = db.Model(&gm.RoleMenu{}).Where("role_id = ?", req.RoleId).Find(&roleMenus).Error

	if err != nil {
		return
	}

	for _, rm := range roleMenus {
		menuIds = append(menuIds, rm.MenuId)
	}

	err = db.Model(&gm.Menu{}).Where("id in (?)", menuIds).Find(&menus).Error

	if err != nil {
		return
	}

	result := make([]response.Menu, 0, len(menus))

	for _, m := range menus {
		result = append(result, response.Menu{
			Id:       m.Id,
			Name:     m.Name,
			Path:     m.Path,
			Label:    m.Label,
			Icon:     m.Icon,
			Url:      m.Url,
			Status:   m.Status,
			ParentId: m.ParentId,
		})
	}
	return
}

// 角色新增菜单
func AddRoleMenu(form request.AddRoleMenu) (err error) {
	db := global.DB
	//校验角色是否存在
	role := new(gm.Role)
	err = role.CheckStatusById(db, form.RoleId)
	if err != nil {
		return
	}

	//校验菜单是否存在
	menu := new(gm.Menu)
	err = menu.CheckStatusById(db, form.RoleId)
	if err != nil {
		return
	}

	roleMenu := &gm.RoleMenu{
		RoleId: form.RoleId,
		MenuId: form.MenuId,
		Status: 1,
	}

	return roleMenu.Save(db)
}
