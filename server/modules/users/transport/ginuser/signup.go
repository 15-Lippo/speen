package ginuser

import (
	"context"
	"net/http"
	"nolan/spin-game/components/appctx"
	"nolan/spin-game/components/common"
	"nolan/spin-game/components/hasher"
	"nolan/spin-game/components/pubsub"
	userbiz "nolan/spin-game/modules/users/biz"
	usermodel "nolan/spin-game/modules/users/model"
	userrepo "nolan/spin-game/modules/users/repository"
	userstorage "nolan/spin-game/modules/users/storage"

	"github.com/gin-gonic/gin"
)

func SignUp(ctx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := ctx.GetMaiDBConnection()
		pb := ctx.GetPubsub()
		var data usermodel.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}
		data.Nonce = common.RandomNonce()

		if err := data.IsSignatureValid(); err != nil {
			panic(err)
		}

		md5 := hasher.NewMd5Hash()
		storage := userstorage.NewSQLStore(db)
		userRepo := userrepo.NewUserRepo(storage)

		biz := userbiz.NewSignupBiz(userRepo, md5)

		result, err := biz.Signup(c.Request.Context(), &data)
		if err != nil {
			panic(err)
		}

		pb.Publish(context.Background(), common.ChannelCreateUser, pubsub.NewMessage(result))

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(map[string]interface{}{
			"id": (*result).Id,
		}))
	}
}
