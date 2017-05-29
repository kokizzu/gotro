package model

const (
	ERR_001_MUST_LOGIN             = `001 must be logged in`
	ERR_004_INVALID_RESET_LINK     = `004 Invalid reset password link`
	ERR_005_INVALID_RESET_EMAIL    = `005 Invalid reset password link e-mail`
	ERR_006_INVALID_RESET_ASSOC    = `006 Invalid reset password link and e-mail association`
	ERR_007_UPLOAD_ERROR           = `007 Upload error: `
	ERR_008_INVALID_CONTENT_TYPE   = `008 Invalid file content-type: `
	ERR_009_RECORD_NOT_FOUND       = `009 Record not found, ID: `
	ERR_010_PASSWORD_TOO_SHORT     = `010 Password too short`
	ERR_011_INCORRECT_OLD_PASSWORD = `011 Incorrect old password`
	ERR_012_INCORRECT_PASSWORD     = `012 Incorrect password`

	ERR_201_FAILED_OAUTH_EXCHANGE      = `201 Failed OAuth Exchange: `
	ERR_206_MISSING_OAUTH_PROVIDER     = `206 Missing OAuth Provider `
	ERR_207_FB_AK_TOKEN_EXCHANGE_ERROR = `207 FBAK Token Exchange Error: `
	ERR_208_FB_AK_USER_INFO_ERROR      = `208 FBAK User Info Error: `
	ERR_209_PHONE_LOGIN_FAILED         = `209 Phone Login Failed`

	ERR_301_WRONG_USERNAME_OR_PASSWORD       = `301 Wrong username or password; username is case sensitive`
	ERR_302_NAME_EMAIL_COMBINATION_NOT_FOUND = `302 Name and e-mail combination not found`
	ERR_303_TOO_SOON_RESET                   = `303 Reset password link has been sent before, please retry in `
	ERR_304_FAILED_SEND_RESET_EMAIL          = `304 Failed to send reset password link, please contact`
	ERR_305_EMAIL_NOT_REGISTERED             = `305 E-Mail not found on database: `
	ERR_306_CSRF_STATE                       = `306 Invalid CSRF State: `

	ERR_402_INVALID_ACTION = `402 Invalid action: `
)
