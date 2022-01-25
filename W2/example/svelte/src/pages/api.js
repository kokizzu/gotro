// can be hit using with /api/[ApiName]
export const LastUpdatedAt = 1643081754
export const APIs = {
	Health: {
		in: {
		}, out: {
		}, read: [
		], write: [
		], stat: [
		], deps: [
		], err: []
	},
	StoreCartItemsAdd: {
		in: {
			productId: 0, // uint64
			deltaQty: 0, // int64
			sessionToken: '', //string | user login token
		}, out: {
			cartItems: [{
				id:  '', // uint64
				createdAt:  0, // int64
				createdBy:  '', // uint64
				updatedAt:  0, // int64
				updatedBy:  '', // uint64
				deletedAt:  0, // int64
				deletedBy:  '', // uint64
				isDeleted:  false, // bool
				restoredAt:  0, // int64
				restoredBy:  '', // uint64
				ownerId:  '', // uint64
				invoiceId:  '', // uint64
				productId:  '', // uint64
				nameCopy:  '', // string
				priceCopy:  0, // int64
				qty:  0, // int64
				discount:  0, // uint64
				subTotal:  0, // int64
				info:  '', // string
			}],
			total: 0, // uint32
			isOverflow: false, // bool
		}, read: [
			"Auth.Sessions",
			"Store.Products",
		], write: [
			"Store.CartItems",
		], stat: [
		], deps: [
		], err: [
			[400, `cannot add more`],
			[400, `cannot remove more`],
			[400, `invalid session token`],
			[400, `missing session token`],
			[400, `session missing from database, wrong env?`],
			[400, `token expired`],
			[403, `must login`],
			[403, `session expired or logged out`],
			[404, `cart item not found`],
			[404, `product not found`],
			[500, `failed add/remove item on cart`],
			[500, `failed insert to cart`],
		]
	},
	StoreInvoice: {
		in: {
			invoiceId: 0, // uint64
			recalculate: false, // bool
			doPurchase: false, // bool
			sessionToken: '', //string | user login token
		}, out: {
			cartItems: [{
				id:  '', // uint64
				createdAt:  0, // int64
				createdBy:  '', // uint64
				updatedAt:  0, // int64
				updatedBy:  '', // uint64
				deletedAt:  0, // int64
				deletedBy:  '', // uint64
				isDeleted:  false, // bool
				restoredAt:  0, // int64
				restoredBy:  '', // uint64
				ownerId:  '', // uint64
				invoiceId:  '', // uint64
				productId:  '', // uint64
				nameCopy:  '', // string
				priceCopy:  0, // int64
				qty:  0, // int64
				discount:  0, // uint64
				subTotal:  0, // int64
				info:  '', // string
			}],
			invoice: {
				id:  '', // uint64
				createdAt:  0, // int64
				createdBy:  '', // uint64
				updatedAt:  0, // int64
				updatedBy:  '', // uint64
				deletedAt:  0, // int64
				deletedBy:  '', // uint64
				isDeleted:  false, // bool
				restoredAt:  0, // int64
				restoredBy:  '', // uint64
				ownerId:  '', // uint64
				totalWeight:  0, // uint64
				totalPrice:  0, // uint64
				totalDiscount:  0, // uint64
				deliveryMethod:  0, // uint64
				deliveryPrice:  0, // uint64
				totalPaid:  0, // uint64
				paidAt:  0, // uint64
				paymentMethod:  0, // uint64
				deadlineAt:  0, // uint64
				promoRuleIds:  '', // string
			},
		}, read: [
			"Auth.Sessions",
			"Store.CartItems",
			"Store.Products",
			"Store.Promos",
		], write: [
			"Store.CartItems",
			"Store.Invoices",
		], stat: [
		], deps: [
		], err: [
			[400, `invalid session token`],
			[400, `missing session token`],
			[400, `session missing from database, wrong env?`],
			[400, `token expired`],
			[403, `must login`],
			[403, `session expired or logged out`],
		]
	},
	StoreProducts: {
		in: {
			limit: 0, // uint32
			offset: 0, // uint32
		}, out: {
			limit: 0, // uint32
			offset: 0, // uint32
			total: 0, // uint32
			products: [{
				id:  '', // uint64
				createdAt:  0, // int64
				createdBy:  '', // uint64
				updatedAt:  0, // int64
				updatedBy:  '', // uint64
				deletedAt:  0, // int64
				deletedBy:  '', // uint64
				isDeleted:  false, // bool
				restoredAt:  0, // int64
				restoredBy:  '', // uint64
				sku:  '', // string
				name:  '', // string
				price:  0, // uint64
				inventoryQty:  0, // uint64
				weightGram:  0, // uint64
			}],
		}, read: [
			"Store.Products",
		], write: [
		], stat: [
		], deps: [
		], err: []
	},
	UserChangeEmail: {
		in: {
		}, out: {
		}, read: [
		], write: [
		], stat: [
		], deps: [
		], err: []
	},
	UserChangePassword: {
		in: {
			password: '', // string
			newPassword: '', // string
			sessionToken: '', //string | user login token
		}, out: {
			updatedAt: 0, // int64
		}, read: [
			"Auth.Sessions",
		], write: [
			"Auth.Users",
		], stat: [
		], deps: [
		], err: [
			[400, `invalid session token`],
			[400, `missing session token`],
			[400, `session missing from database, wrong env?`],
			[400, `token expired`],
			[401, `wrong password`],
			[403, `session expired or logged out`],
			[500, `cannot encrypt password`],
			[500, `failed to change password`],
			[500, `user not found`],
		]
	},
	UserConfirmEmail: {
		in: {
		}, out: {
		}, read: [
		], write: [
		], stat: [
		], deps: [
		], err: []
	},
	UserExternalLogin: {
		in: {
			provider: '', // string
		}, out: {
			link: '', // string
		}, read: [
		], write: [
		], stat: [
		], deps: [
		], err: [
			[400, "Invalid host for oauth"],
		]
	},
	UserForgotPassword: {
		in: {
			email: '', // string
			changePassCallback: '', // string
		}, out: {
			ok: false, // bool
		}, read: [
		], write: [
			"Auth.Users",
		], stat: [
		], deps: [
		], err: [
			[400, `email not found`],
			[500, `failed to update row on database`],
		]
	},
	UserList: {
		in: {
			limit: 0, // uint32
			offset: 0, // uint32
		}, out: {
			limit: 0, // uint32
			offset: 0, // uint32
			total: 0, // uint32
			users: [{
				id:  '', // uint64
				email:  '', // string
				password:  '', // string
				createdAt:  0, // int64
				createdBy:  '', // uint64
				updatedAt:  0, // int64
				updatedBy:  '', // uint64
				deletedAt:  0, // int64
				deletedBy:  '', // uint64
				isDeleted:  false, // bool
				restoredAt:  0, // int64
				restoredBy:  '', // uint64
				passwordSetAt:  0, // int64
				secretCode:  '', // string
				secretCodeAt:  0, // int64
				verificationSentAt:  0, // int64
				verifiedAt:  0, // int64
				lastLoginAt:  0, // int64
			}],
		}, read: [
			"nsync.NamedMutex",
			"Auth.Users",
		], write: [
		], stat: [
		], deps: [
		], err: []
	},
	UserLogin: {
		in: {
			email: '', // string
			password: '', // string
		}, out: {
			sessionToken: '', //string | login token
		}, read: [
			"Auth.Users",
		], write: [
			"Auth.Sessions",
		], stat: [
		], deps: [
		], err: [
			[401, `wrong email or password`],
			[401, `wrong password`],
			[500, `cannot create session`],
		]
	},
	UserLogout: {
		in: {
		}, out: {
			loggedOut: false, // bool
			sessionToken: '', //string | login token
		}, read: [
		], write: [
			"Auth.Sessions",
		], stat: [
		], deps: [
		], err: []
	},
	UserOauth: {
		in: {
			provider: '', // string
		}, out: {
			link: '', // string
		}, read: [
		], write: [
		], stat: [
		], deps: [
		], err: [
			[400, "Invalid host for oauth"],
		]
	},
	UserProfile: {
		in: {
			sessionToken: '', //string | user login token
		}, out: {
			user: {
				id:  '', // uint64
				email:  '', // string
				password:  '', // string
				createdAt:  0, // int64
				createdBy:  '', // uint64
				updatedAt:  0, // int64
				updatedBy:  '', // uint64
				deletedAt:  0, // int64
				deletedBy:  '', // uint64
				isDeleted:  false, // bool
				restoredAt:  0, // int64
				restoredBy:  '', // uint64
				passwordSetAt:  0, // int64
				secretCode:  '', // string
				secretCodeAt:  0, // int64
				verificationSentAt:  0, // int64
				verifiedAt:  0, // int64
				lastLoginAt:  0, // int64
			},
		}, read: [
			"Auth.Sessions",
			"Auth.Users",
		], write: [
		], stat: [
		], deps: [
		], err: [
			[400, `invalid session token`],
			[400, `missing session token`],
			[400, `session missing from database, wrong env?`],
			[400, `token expired`],
			[403, `session expired or logged out`],
			[404, `user does not exists on database: `],
		]
	},
	UserRegister: {
		in: {
			userName: '', // string
			email: '', // string
			password: '', // string
		}, out: {
			createdAt: 0, // int64
			userId: '', // uint64
		}, read: [
		], write: [
			"Auth.Users",
		], stat: [
		], deps: [
		], err: [
			[400, `email must not be empty`],
			[400, `userName must not be empty`],
			[451, `failed to register this user: `],
			[451, `user already exists: `],
			[500, `cannot encrypt password`],
		]
	},
	UserResetPassword: {
		in: {
			password: '', // string
			secretCode: '', // string
			hash: '', // string
		}, out: {
			ok: false, // bool
		}, read: [
		], write: [
			"Auth.Users",
		], stat: [
		], deps: [
		], err: [
			[400, `cannot find user, wrong env?`],
			[400, `invalid hash`],
			[400, `invalid secret code`],
			[400, `secret code expired`],
			[500, `cannot encrypt password`],
			[500, `failed to update user password`],
		]
	},
}
