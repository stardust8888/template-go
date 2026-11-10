package user

import (
	"github.com/gin-gonic/gin"

	"template-go/db/mysql"
	"template-go/library/ginx"
	"template-go/model"
	"template-go/utils/common"
)

// GetListReq 用户列表请求
type GetListReq struct {
	common.PageReq
	Username    string `form:"username"` // 用户名模糊搜索（可选）
	PhoneSuffix string `form:"phone_suffix"`
	Status      *int   `form:"status"` // 状态过滤 0:正常 1:禁用（可选，不传则不过滤）
}

type GetListResp struct {
	User        *model.User `json:"user"`
	Phone       string      `json:"phone"`
	PhoneSuffix string      `json:"phone_suffix"`
	Wx          string      `json:"wx"`
}

// GetList 用户列表，直接查 user 表
func GetList(c *gin.Context) {
	ctx := c.Request.Context()
	req := GetListReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.ResponseParamError(c)
		return
	}
	req.BindDefault()

	users := make([]*model.User, 0)
	where := map[string]any{
		mysql.Order: "id desc",
	}
	if req.Username != "" {
		where["username LIKE ?"] = mysql.LikeQuery(mysql.Like, req.Username)
	}
	if req.PhoneSuffix != "" {
		where["phone_suffix"] = req.PhoneSuffix
	}
	if req.Status != nil {
		where["status"] = *req.Status
	}
	mysql.SetPage(where, req.Page, req.PageSize)

	total, err := mysql.PageQuery(ctx, mysql.GetUserDB(), &users, where)
	if err != nil {
		ginx.ResponseServerError(c, err)
		return
	}

	resp := make([]GetListResp, 0, len(users))
	for _, user := range users {
		resp = append(resp, GetListResp{
			User:        user,
			Phone:       user.Phone,
			PhoneSuffix: user.PhoneSuffix,
			Wx:          user.Wx,
		})
	}

	ginx.ResponsePageSuccess(c, resp, req.Page, req.PageSize, total)
}
