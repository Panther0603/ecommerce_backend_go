package database

import "errors"

var (
	// here defing some of teh general error

	ErrorCantFindProduct   = errors.New("can't find the producrs ")
	ErrorCantDecodeProduct = errors.New("can't able to decode the product")
	ErrUserIdIsNotValid    = errors.New("this user is not valid ")
	ErrCantUpdateUser      = errors.New("not able to update user")
	ErrCantRemoveCartItem  = errors.New("cannot able to remove this item from cart")
	ErrCantGetItem         = errors.New("was unable tp get the item, please try after sometime")
	ErrCantBuyCartItem     = errors.New("can not update purache order ")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func IntantBuy() {

}
