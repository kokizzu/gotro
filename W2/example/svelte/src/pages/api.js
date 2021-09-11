// can be hit using with /api/[ApiName]
export const LastUpdatedAt = 1630271387
export const APIs = {
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
			walletId: '', // string
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
