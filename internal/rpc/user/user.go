package user

import "store/internal/rpc/base"

type User struct {
	*Follow
	*Footprint
	*Favourites
	*ReceiverAddress
	*Personage
}

func NewUser(b *base.Base) *User {
	return &User{
		NewFollow(b),
		NewFootprint(b),
		NewFavourites(b),
		NewReceiverAddress(b),
		NewPersonage(b),
	}
}
