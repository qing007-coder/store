package gateway

import (
	"github.com/gin-gonic/gin"
	"store/pkg/config"
)

type Router struct {
	user   *UserApi
	engine *gin.Engine
}

func NewRouter(conf *config.GlobalConfig) *Router {
	r := new(Router)
	r.init(conf)

	return r
}

func (r *Router) init(conf *config.GlobalConfig) {
	r.engine = gin.Default()
	srv := NewService(conf)
	r.user = NewUserApi(srv)

	root := r.engine.Group("api")
	{
		user := root.Group("user")
		{
			user.POST("update_personal_info", r.user.UpdatePersonalInfo)
			user.POST("modify_password", r.user.ModifyPassword)
			user.POST("add_receiver_address", r.user.AddReceiverAddress)
			user.POST("delete_receiver_address", r.user.DeleteReceiverAddress)
			user.POST("update_receiver_address", r.user.UpdateReceiverAddress)
			user.GET("get_receiver_address", r.user.GetReceiverAddress)
			user.POST("add_favourites", r.user.AddFavourites)
			user.POST("delete_favourites", r.user.DeleteFavourites)
			user.GET("get_favourites_list", r.user.GetFavouritesList)
			user.POST("add_footprint", r.user.AddFootprint)
			user.POST("delete_footprint", r.user.DeleteFootprint)
			user.GET("get_footprint_list", r.user.GetFootprintList)
			user.POST("follow_merchant", r.user.FollowMerchant)
			user.POST("cancel_follow", r.user.CancelFollow)
			user.GET("get_follow_list", r.user.GetFollowList)
		}
	}
}

func (r *Router) Run() {

}
