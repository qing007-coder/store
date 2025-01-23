package gateway

import (
	"github.com/gin-gonic/gin"
	"store/pkg/config"
)

type Router struct {
	user        *UserApi
	merchandise *MerchandiseApi
	engine      *gin.Engine
	middleware  *Middleware
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
	r.merchandise = NewMerchantApi(srv)
	r.middleware = NewMiddleware()

	r.engine.Use(r.middleware.GetAuthorizationHeader())
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

		merchandise := root.Group("merchandise")
		{
			merchandise.POST("put_away_merchandise", r.merchandise.PutAwayMerchandise)
			merchandise.POST("remove_merchandise", r.merchandise.RemoveMerchandise)
			merchandise.POST("update_merchandise", r.merchandise.UpdateMerchandise)
			merchandise.GET("get_merchandise_details", r.merchandise.GetMerchandiseDetails)
			merchandise.GET("search", r.merchandise.Search)
			merchandise.GET("search_by_category", r.merchandise.SearchByCategory)
			merchandise.POST("add_merchandise_style", r.merchandise.AddMerchandiseStyle)
			merchandise.POST("remove_merchandise_style", r.merchandise.RemoveMerchandiseStyle)
			merchandise.POST("update_merchandise_style", r.merchandise.UpdateMerchandiseStyle)
			merchandise.GET("get_merchandise_style_list", r.merchandise.GetMerchandiseStyleList)
			merchandise.GET("get_merchandise_style_details", r.merchandise.GetMerchandiseStyleDetails)
		}
	}
}

func (r *Router) Run() error {
	return r.engine.Run(":8080")
}
