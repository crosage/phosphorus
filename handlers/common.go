package handlers

import (
	"chat/structs"
	"chat/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"time"
)

func InitHandlers(app *fiber.App) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	app.Post("/api/user", userRegister)
	app.Post("/api/user/login", userLogin)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(utils.JwtKey),
		},
		ContextKey: "user", // 将用户信息存储在 ctx.locals["user"]
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// 鉴权失败的处理
			return sendCommonResponse(ctx, 401, "过期或非法JWT", map[string]interface{}{
				"path": ctx.Path(),
			})
		},
	}))

	// 需要 JWT 鉴权的路由
	app.Post("/api/messages/private", sendPrivateMessage) // 发送私人消息
	app.Post("/api/messages/group", sendGroupMessage)     // 发送群聊消息

	// 获取私人消息和群聊消息
	app.Get("/api/messages/private/:user_id", getPrivateMessages) // 获取私人消息
	app.Get("/api/messages/group/:group_id", getGroupMessages)    // 获取群聊消息

	// 获取用户信息
	app.Get("/api/user/:user_id", getUserInfo) // 获取用户信息

	// 搜索用户
	app.Get("/api/users/search", searchUsers) // 搜索用户

	// 群聊操作
	app.Post("/api/groups", createGroup)                // 创建群组
	app.Get("/api/groups/:group_id", getGroupInfo)      // 获取群组信息
	app.Post("/api/groups/:group_id/join", joinGroup)   // 加入群组
	app.Post("/api/groups/:group_id/leave", leaveGroup) // 退出群组
}

// JWT 生成函数
func generateJWT(userID int) (string, error) {
	// 创建 JWT claims
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // 设置过期时间为 72 小时
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用秘钥签名生成 token
	return token.SignedString([]byte(utils.JwtKey)) // utils.JwtKey 是密钥
}
func sendCommonResponse(ctx *fiber.Ctx, code int, message string, data map[string]interface{}) error {
	response := structs.Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	json, err := jsoniter.Marshal(response)
	if err != nil {
		// THIS SHOULD NOT HAPPEN
		// If this happens, just stop the server and wait for further investigation

	}
	return ctx.Status(code).Send(json)
}
