package manage

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"main.go/global"
	"main.go/model/manage"
	manageReq "main.go/model/manage/request"
	"main.go/utils"
)

type ManageAdminUserService struct {
}

// CreateMallAdminUser 创建MallAdminUser记录
func (m *ManageAdminUserService) CreateMallAdminUser(mallAdminUser manage.MallAdminUser) (err error) {
	if !errors.Is(global.GVA_DB.Where("login_user_name = ?", mallAdminUser.LoginUserName).First(&manage.MallAdminUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同用户名")
	}
	err = global.GVA_DB.Create(&mallAdminUser).Error
	return err
}

// UpdateMallAdminName 更新MallAdminUser昵称
func (m *ManageAdminUserService) UpdateMallAdminName(token string, req manageReq.MallUpdateNameParam) (err error) {
	var adminUserToken manage.MallAdminUserToken
	err = global.GVA_DB.Where("token =? ", token).First(&adminUserToken).Error
	if err != nil {
		return errors.New("不存在的用户")
	}
	err = global.GVA_DB.Where("admin_user_id = ?", adminUserToken.AdminUserId).Updates(&manage.MallAdminUser{
		LoginUserName: req.LoginUserName,
		NickName:      req.NickName,
	}).Error
	return err
}

func (m *ManageAdminUserService) UpdateMallAdminPassWord(token string, req manageReq.MallUpdatePasswordParam) (err error) {
	var adminUserToken manage.MallAdminUserToken
	err = global.GVA_DB.Where("token =? ", token).First(&adminUserToken).Error
	if err != nil {
		return errors.New("用户未登录")
	}
	var adminUser manage.MallAdminUser
	err = global.GVA_DB.Where("admin_user_id =?", adminUserToken.AdminUserId).First(&adminUser).Error
	if err != nil {
		return errors.New("不存在的用户")
	}
	if adminUser.LoginPassword != req.OriginalPassword {
		return errors.New("原密码不正确")
	}
	adminUser.LoginPassword = req.NewPassword

	err = global.GVA_DB.Where("admin_user_id=?", adminUser.AdminUserId).Updates(&adminUser).Error
	return
}

// GetMallAdminUser 根据id获取MallAdminUser记录
func (m *ManageAdminUserService) GetMallAdminUser(c *gin.Context) (err error, mallAdminUser manage.MallAdminUser) {
	userName := c.GetString("userName")
	// if errors.Is(global.GVA_DB.Where("token =?", token).First(&adminToken).Error, gorm.ErrRecordNotFound) {
	// 	return errors.New("不存在的用户"), mallAdminUser
	// }
	err = global.GVA_DB.Where("login_user_name = ?", userName).First(&mallAdminUser).Error
	return err, mallAdminUser
}

// AdminLogin 管理员登陆
func (m *ManageAdminUserService) AdminLogin(params manageReq.MallAdminLoginParam) (err error, mallAdminUser manage.MallAdminUser, adminToken string) {

	err = global.GVA_DB.Where("login_user_name=? AND login_password=?", params.UserName, params.PasswordMd5).First(&mallAdminUser).Error

	if mallAdminUser != (manage.MallAdminUser{}) {
		// token := getNewToken(time.Now().UnixNano()/1e6, int(mallAdminUser.AdminUserId))
		// global.GVA_DB.Where("admin_user_id", mallAdminUser.AdminUserId).First(&adminToken)
		// nowDate := time.Now()
		// 48小时过期
		// expireTime, _ := time.ParseDuration("48h")
		// expireDate := nowDate.Add(expireTime)
		// 生成 jwt
		// generate jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userName": mallAdminUser.LoginUserName,
			"userId":   mallAdminUser.AdminUserId,
			"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		})
		tokenString, err2 := token.SignedString([]byte(global.GVA_CONFIG.Redis.SecretKey))
		if err2 != nil {
			global.GVA_LOG.Error("redis 生成token失败")
			return
		}
		key := "auth_" + mallAdminUser.LoginUserName
		err1 := global.GVA_REDIS.Set(key, tokenString, time.Hour).Err()
		global.GVA_LOG.Info("token: " + tokenString)
		global.GVA_LOG.Info("key: " + key)
		if err1 != nil {
			global.GVA_LOG.Error(fmt.Sprintf("redis set error, 用户名id：%d", mallAdminUser.AdminUserId))
			return
		}
		// 没有token新增，有token 则更新
		// if adminToken == (manage.MallAdminUserToken{}) {
		// 	adminToken.AdminUserId = mallAdminUser.AdminUserId
		// 	adminToken.Token = token
		// 	adminToken.UpdateTime = nowDate
		// 	adminToken.ExpireTime = expireDate
		// 	if err = global.GVA_DB.Create(&adminToken).Error; err != nil {
		// 		return
		// 	}
		// } else {
		// 	adminToken.Token = token
		// 	adminToken.UpdateTime = nowDate
		// 	adminToken.ExpireTime = expireDate
		// 	if err = global.GVA_DB.Save(&adminToken).Error; err != nil {
		// 		return
		// 	}
		// }
		return err, mallAdminUser, tokenString
	}
	return err, mallAdminUser, ""

}

func getNewToken(timeInt int64, userId int) (token string) {
	var build strings.Builder
	build.WriteString(strconv.FormatInt(timeInt, 10))
	build.WriteString(strconv.Itoa(userId))
	build.WriteString(utils.GenValidateCode(6))
	return utils.MD5V([]byte(build.String()))
}
