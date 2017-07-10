package sql

import (
	"errors"
	"fmt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/W"
)

const (
	ERR_001_MUST_LOGIN             = `001 must be logged in`
	ERR_002_INVALID_UPLOAD_KEY     = `002 Invalid Upload Key: `
	ERR_004_INVALID_RESET_LINK     = `004 Invalid reset password link`
	ERR_005_INVALID_RESET_EMAIL    = `005 Invalid reset password link e-mail`
	ERR_006_INVALID_RESET_ASSOC    = `006 Invalid reset password link and e-mail association`
	ERR_007_UPLOAD_ERROR           = `007 Upload error: `
	ERR_008_INVALID_CONTENT_TYPE   = `008 Invalid file content-type: `
	ERR_009_RECORD_NOT_FOUND       = `009 Record not found, ID: `
	ERR_010_PASSWORD_TOO_SHORT     = `010 Password too short`
	ERR_011_INCORRECT_OLD_PASSWORD = `011 Incorrect old password`
	ERR_012_INCORRECT_PASSWORD     = `012 Incorrect password`
	ERR_013_FAILED_RESET_PASSWORD  = `013 Failed resetting password`

	ERR_121_CONVERT_FILE = `121 Error resizing/converting image`
	ERR_122_STAT_FILE    = `122 Error stating file`

	ERR_201_FAILED_OAUTH_EXCHANGE      = `201 Failed OAuth Exchange: `
	ERR_206_MISSING_OAUTH_PROVIDER     = `206 Missing OAuth Provider `
	ERR_207_FB_AK_TOKEN_EXCHANGE_ERROR = `207 FBAK Token Exchange Error: `
	ERR_208_FB_AK_USER_INFO_ERROR      = `208 FBAK User Info Error: `
	ERR_209_PHONE_LOGIN_FAILED         = `209 Phone Login Failed`

	ERR_301_WRONG_USERNAME_OR_PASSWORD       = `301 Wrong username or password; username is case sensitive, probably phone not registered, contact the administrator`
	ERR_302_NAME_EMAIL_COMBINATION_NOT_FOUND = `302 Name and e-mail combination not found`
	ERR_303_TOO_SOON_RESET                   = `303 Reset password link has been sent before, please retry in `
	ERR_304_FAILED_SEND_RESET_EMAIL          = `304 Failed to send reset password link, please contact`
	ERR_305_EMAIL_NOT_REGISTERED             = `305 E-Mail not found on database: `
	ERR_306_CSRF_STATE                       = `306 Invalid CSRF State: `

	ERR_402_NO_ACTION         = `402 Invalid action: `
	ERR_403_MUST_LOGIN_HIGHER = `403 This page can only be accessed by user with higher privilege`

	ERR_501_CANNOT_CREATE_DIR  = `501 Error creating directory: `
	ERR_502_CANNOT_CREATE_FILE = `502 File creation error: `
	ERR_503_CANNOT_MOVE_FILE   = `503 Move file error: `
)

func ErrorReport(skip int, title string) {
	if err := recover(); err != nil {
		//L.Panic(errors.New(`Internal Server Error`), ``, err)
		err2, ok := err.(error)
		if !ok {
			err2 = errors.New(fmt.Sprintf("%# v", err))
		}
		err_str := err2.Error()
		L.LOG.Errorf(err_str)
		str := L.StackTrace(skip)
		L.LOG.Criticalf("StackTrace: %s", str)
		str = S.Replace("\n", `<br/>`, str)
		W.Mailers[``].SendBCC([]string{DEBUGGER_EMAIL}, title, err_str+"\n"+str)
	}
}
