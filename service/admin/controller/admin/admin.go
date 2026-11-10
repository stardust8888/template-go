package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"template-go/db/mysql"
	"template-go/library/ginx"
	"template-go/library/jwt"
	"template-go/library/security"
	"template-go/logs"
	"template-go/model"
	"template-go/utils"
	"template-go/utils/constant"
	"template-go/utils/x/timex"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token       string       `json:"token"`
	ExpiresAt   timex.Time   `json:"expiresAt"`
	Permissions []string     `json:"permissions"`
	User        *model.Admin `json:"user"`
}

// Login 管理员登录
func Login(c *gin.Context) {
	ctx := c.Request.Context()
	req := LoginReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ginx.ResponseParamError(c)
		return
	}

	var admin model.Admin
	err = mysql.FindOne(ctx, mysql.GetAdminDB(), &admin, map[string]any{"username": req.Username})
	if err != nil {
		ginx.ResponseServerError(c, err)
		return
	}

	if admin.Status != 1 {
		ginx.ResponseError(c, constant.CodeUserDisabled)
		return
	}

	fingerprint := utils.MD5(req.Password + "")
	if fingerprint != admin.Password {
		ginx.ResponseError(c, constant.CodePasswordError)
		return
	}

	jt, err := jwt.New("", timex.Day)
	if err != nil {
		ginx.ResponseServerError(c, err)
		return
	}

	simple := constant.SimpleUser{
		ID:      admin.ID,
		IsAdmin: admin.IsSuper,
	}
	token, err := jt.GenerateToken(simple)
	if err != nil {
		ginx.ResponseServerError(c, err)
		return
	}

	err = mysql.Save(ctx, mysql.GetAdminDB(), &admin)
	if err != nil {
		logs.Error(c, "更新管理员登录信息失败", zap.Error(err))
		//ginx.ResponseError(c, constant.CodeServerInternal)
		//return
	}
	admin.Password = constant.MaskedString

	loginLog := model.AdminLoginLog{
		AdminID:   admin.ID,
		LoginIP:   c.ClientIP(),
		LoginTime: timex.GetBeijingTime(),
	}
	err = mysql.Create(ctx, mysql.GetAdminDB(), &loginLog)
	if err != nil {
		logs.Error(c, "更新管理员登录日志失败", zap.Error(err))
	}

	response := LoginResp{
		Token:     token,
		ExpiresAt: timex.GetBeijingTime().Add(timex.Day),
		User:      &admin,
	}
	ginx.ResponseSuccess(c, response)
}

// Generate 随机生成 4 个 32 位密钥（JWT/INF/PWD/CFG），用 KMS 子密钥加密后写入 key_center。
// 按 id 逐条写入，已存在的不覆盖，避免正在使用的密钥被替换导致旧数据无法解密。
func Generate(c *gin.Context) {
	ctx := c.Request.Context()
	key, err := security.DeriveKey(security.PurposeKMSCrypt)
	if err != nil {
		logs.Error(ctx, "加载密钥失败", zap.Error(err))
		ginx.ResponseServerError(c, err)
		return
	}

	if _, err = mysql.Count(ctx, mysql.GetUserDB(), &model.KeyCenter{}, map[string]any{}); err != nil {
		ginx.ResponseServerError(c, err)
		return
	}

	keyCenters := []*model.KeyCenter{
		{GormModel: model.GormModel{ID: 1}, Desc: "JWT"},
		{GormModel: model.GormModel{ID: 2}, Desc: "INF"},
		{GormModel: model.GormModel{ID: 3}, Desc: "PWD"},
		{GormModel: model.GormModel{ID: 4}, Desc: "CFG"},
	}

	// 按 id 逐条处理：已存在的跳过（不覆盖正在使用的密钥），不存在的生成后写入
	for _, kc := range keyCenters {
		exist, err := mysql.Exist(ctx, mysql.GetUserDB(), &model.KeyCenter{}, map[string]any{"id": kc.ID})
		if err != nil {
			ginx.ResponseServerError(c, err)
			continue
		}
		if exist {
			continue
		}

		secret, err := utils.AESGCMEncryptHex([]byte(utils.RandString(32)), key)
		if err != nil {
			ginx.ResponseServerError(c, err)
			return
		}
		kc.Secret = secret

		if err = mysql.Create(ctx, mysql.GetUserDB(), kc); err != nil {
			ginx.ResponseServerError(c, err)
			return
		}
	}

	ginx.ResponseSuccess(c, nil)
}
