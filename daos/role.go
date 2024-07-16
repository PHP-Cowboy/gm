package daos

import (
	"gm/global"
	"gm/model"
	"gm/model/gm"
	"gm/request"
	"gm/response"
)

func CreateRole(req request.CreateRole) (err error) {
	db := global.DB

	role := gm.Role{
		Name: req.Name,
		Desc: req.Desc,
		Creator: model.Creator{
			CreatorId: req.CreatorId,
			Creator:   req.Creator,
		},
	}

	err = db.Model(&gm.Role{}).Save(&role).Error

	if err != nil {
		global.Logger["err"].Infof("保存角色数据失败:" + err.Error())
		return
	}

	roleMenuObj := new(gm.RoleMenu)

	err = roleMenuObj.DelAndCreateRoleMenu(db, role.Id, req.MenuList)
	if err != nil {
		global.Logger["err"].Infof("保存角色菜单数据失败:" + err.Error())
		return err
	}

	roleChanObj := new(gm.RoleChannel)

	err = roleChanObj.DelAndCreateRoleMenu(db, role.Id, req.ChannelList)

	if err != nil {
		global.Logger["err"].Infof("保存角色渠道数据失败:" + err.Error())
		return err
	}

	return
}

func ChangeRole(req request.ChangeRole) (err error) {
	db := global.DB
	mp := make(map[string]interface{}, 0)
	//更新用户密码
	if req.Desc != "" {
		mp["desc"] = req.Desc
	}

	//更新用户名称
	if req.Name != "" {
		mp["name"] = req.Name
	}

	//更新用户状态
	if req.Status > 0 {
		mp["status"] = req.Status
	}

	roleMenuObj := new(gm.RoleMenu)

	err = roleMenuObj.DelAndCreateRoleMenu(db, req.Id, req.MenuList)
	if err != nil {
		global.Logger["err"].Infof("保存角色菜单数据失败:" + err.Error())
		return err
	}

	err = new(gm.Role).UpdateById(db, req.Id, mp)

	if err != nil {
		global.Logger["err"].Infof("更新角色数据失败:" + err.Error())
		return
	}

	roleChanObj := new(gm.RoleChannel)

	err = roleChanObj.DelAndCreateRoleMenu(db, req.Id, req.ChannelList)

	if err != nil {
		global.Logger["err"].Infof("保存角色渠道数据失败:" + err.Error())
		return err
	}

	return
}

func RoleList(req request.RoleList) (result response.RoleList, err error) {
	db := global.DB

	var (
		roles []gm.Role
		total int64
	)

	localDb := db.Model(&gm.Role{})

	if req.Id > 0 {
		localDb.Where("id = ?", req.Id)
	}

	if req.Name != "" {
		localDb.Where("name like ?", "%"+req.Name+"%")
	}

	if req.Status > 0 {
		localDb.Where("status = ?", req.Status)
	}

	err = localDb.Count(&total).Error

	if err != nil {
		return
	}

	result.Total = total

	err = localDb.Scopes(model.Paginate(req.Page, req.Size)).Find(&roles).Error

	if err != nil {
		return
	}

	roleIds := make([]int, 0, len(roles))

	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}

	roleMenuObj := new(gm.RoleMenu)
	mp := make(map[int][]int)
	mpChannel := make(map[int][]int)

	mp, err = roleMenuObj.MapOfMenuIdsByRoleIds(db, roleIds)
	if err != nil {
		return response.RoleList{}, err
	}

	roleChanObj := new(gm.RoleChannel)

	mpChannel, err = roleChanObj.MapOfMenuIdsByRoleIds(db, roleIds)
	if err != nil {
		return response.RoleList{}, err
	}

	list := make([]response.Role, 0, len(roles))

	for _, r := range roles {

		tmp := response.Role{
			Id:     r.Id,
			Name:   r.Name,
			Desc:   r.Desc,
			Status: r.Status,
		}

		menuIds, ok := mp[r.Id]

		if !ok {
			menuIds = make([]int, 0)
		}

		tmp.MenuList = menuIds

		channelList, cOk := mpChannel[r.Id]

		if !cOk {
			channelList = make([]int, 0)
		}

		tmp.ChannelList = channelList

		list = append(list, tmp)
	}

	result.List = list
	return
}

func UserRoleList(req request.UserRoleList) (result []response.Role, err error) {
	db := global.DB

	var (
		roles     []gm.Role
		userRoles []gm.UserRole
		roleIds   []int
	)

	err = db.Model(&gm.UserRole{}).Where("user_id = ?", req.UserId).Find(&userRoles).Error

	if err != nil {
		return
	}

	for _, ur := range userRoles {
		roleIds = append(roleIds, ur.RoleId)
	}

	err = db.Model(&gm.Role{}).Where("id in (?)", roleIds).Find(&roles).Error

	if err != nil {
		return
	}

	result = make([]response.Role, 0, len(roles))

	for _, r := range roles {
		result = append(result, response.Role{
			Id:     r.Id,
			Name:   r.Name,
			Desc:   r.Desc,
			Status: r.Status,
		})
	}

	return
}

func AddUserRole(req request.AddUserRole) (err error) {
	db := global.DB
	user := new(gm.User)
	//校验用户是否存在
	err = user.CheckStatusById(db, req.UserId)
	if err != nil {
		return
	}

	//校验角色是否存在
	err = new(gm.Role).CheckStatusById(db, req.RoleId)
	if err != nil {
		return
	}

	userRole := &gm.UserRole{
		UserId: req.UserId,
		RoleId: req.RoleId,
		Status: 1,
	}

	return userRole.Save(db)
}
