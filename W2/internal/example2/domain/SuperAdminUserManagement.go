package domain

import (
	"example2/model/mAuth"
	"example2/model/mAuth/rqAuth"
	"example2/model/mAuth/wcAuth"
	"example2/model/zCrud"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file SuperAdminUserManagement.go
//go:generate replacer -afterprefix "Id\" form" "Id,string\" form" type SuperAdminUserManagement.go
//go:generate replacer -afterprefix "json:\"id\"" "json:\"id,string\"" type SuperAdminUserManagement.go
//go:generate replacer -afterprefix "By\" form" "By,string\" form" type SuperAdminUserManagement.go
//go:generate farify doublequote --file SuperAdminUserManagement.go

type (
	SuperAdminUserManagementIn struct {
		RequestCommon
		Cmd      string        `json:"cmd" form:"cmd" query:"cmd" long:"cmd" msg:"cmd"`
		User     rqAuth.Users  `json:"user" form:"user" query:"user" long:"user" msg:"user"`
		WithMeta bool          `json:"withMeta" form:"withMeta" query:"withMeta" long:"withMeta" msg:"withMeta"`
		Pager    zCrud.PagerIn `json:"pager" form:"pager" query:"pager" long:"pager" msg:"pager"`
	}
	SuperAdminUserManagementOut struct {
		ResponseCommon
		Pager zCrud.PagerOut `json:"pager" form:"pager" query:"pager" long:"pager" msg:"pager"`
		Meta  *zCrud.Meta    `json:"meta" form:"meta" query:"meta" long:"meta" msg:"meta"`
		User  *rqAuth.Users  `json:"user" form:"user" query:"user" long:"user" msg:"user"`
		// listing
		Users [][]any `json:"users" form:"users" query:"users" long:"users" msg:"users"`
	}
)

const (
	SuperAdminUserManagementAction = `superAdmin/userManagement`
	ErrUserIdNotFound              = `user id not found`
	ErrTenantCodeNotFound          = `tenant code is not allready`
	ErrInvalidSegment              = `invalid segment`
	ErrUserSaveFailed              = `user save failed`
	ErrUsersEmailDuplicate         = `email already by another user`
)

var SuperAdminUserManagementMeta = zCrud.Meta{
	Fields: []zCrud.Field{
		{
			Name:      mAuth.Id,
			Label:     "ID",
			DataType:  zCrud.DataTypeInt,
			InputType: zCrud.InputTypeHidden,
		},
		{
			Name:      mAuth.TenantCode,
			Label:     "Tenant Code",
			DataType:  zCrud.DataTypeString,
			InputType: zCrud.InputTypeCombobox,
		},
		{
			Name:      mAuth.Email,
			Label:     "Email",
			DataType:  zCrud.DataTypeString,
			InputType: zCrud.InputTypeText,
		},
		{
			Name:      mAuth.FullName,
			Label:     "Full Name",
			DataType:  zCrud.DataTypeString,
			InputType: zCrud.InputTypeText,
		},
		{
			Name:      mAuth.Role,
			Label:     "Role",
			DataType:  zCrud.DataTypeString,
			InputType: zCrud.InputTypeCombobox,
			Ref: []string{
				UserSegment, EntryUserSegment, TenantAdminSegment, ReportViewerSegment,
			},
		},
		{
			Name:      mAuth.CreatedAt,
			Label:     `Created At`,
			ReadOnly:  true,
			DataType:  zCrud.DataTypeInt,
			InputType: zCrud.InputTypeDateTime,
		},
		{
			Name:      mAuth.UpdatedAt,
			Label:     `Updated At`,
			ReadOnly:  true,
			DataType:  zCrud.DataTypeInt,
			InputType: zCrud.InputTypeDateTime,
		},
		{
			Name:      mAuth.DeletedAt,
			Label:     `Deleted At`,
			ReadOnly:  true,
			DataType:  zCrud.DataTypeInt,
			InputType: zCrud.InputTypeDateTime,
		},
		{
			Name:      mAuth.VerifiedAt,
			Label:     `Verified At`,
			ReadOnly:  true,
			DataType:  zCrud.DataTypeInt,
			InputType: zCrud.InputTypeDateTime,
		},
		{
			Name:      mAuth.LastLoginAt,
			Label:     `Last Login At`,
			ReadOnly:  true,
			DataType:  zCrud.DataTypeInt,
			InputType: zCrud.InputTypeDateTime,
		},
	},
}

func (d *Domain) SuperAdminUserManagement(in *SuperAdminUserManagementIn) (out SuperAdminUserManagementOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)
	out.refId = in.User.Id

	sess := d.MustSuperAdmin(in.RequestCommon, &out.ResponseCommon)
	if sess == nil {
		return
	}

	if in.WithMeta {
		out.Meta = &SuperAdminUserManagementMeta
	}

	switch in.Cmd {
	case zCrud.CmdForm:
		if in.User.Id <= 0 {
			out.Meta = &SuperAdminUserManagementMeta
		}

		user := rqAuth.NewUsers(d.AuthOltp)
		user.Id = in.User.Id
		if !user.FindById() {
			out.SetError(400, ErrUserIdNotFound)
			return
		}
		user.CensorFields()
		out.User = user
	case zCrud.CmdUpsert, zCrud.CmdDelete, zCrud.CmdRestore:
		user := wcAuth.NewUsersMutator(d.AuthOltp)
		user.Id = in.User.Id
		if user.Id > 0 {
			if !user.FindById() {
				out.SetError(400, ErrUserIdNotFound)
				return
			}

			if in.Cmd == zCrud.CmdDelete {
				if user.DeletedAt == 0 {
					user.SetDeletedAt(in.UnixNow())
				}
			} else if in.Cmd == zCrud.CmdRestore {
				if user.DeletedAt > 0 {
					user.SetDeletedAt(0)
				}
			}
		} else {
			user.SetCreatedAt(in.UnixNow())
		}

		if user.SetEmail(in.User.Email) {
			user.SetVerifiedAt(0)
			dup := rqAuth.NewUsers(d.AuthOltp)
			dup.Email = in.User.Email
			if dup.FindByEmail() && dup.Id != user.Id {
				out.SetError(400, ErrUsersEmailDuplicate)
				return
			}
		}
		user.SetFullName(in.User.FullName)

		if user.SetRole(in.User.Role) {
			if in.User.Role != UserSegment &&
				in.User.Role != EntryUserSegment &&
				in.User.Role != TenantAdminSegment &&
				in.User.Role != ReportViewerSegment {
				out.SetError(400, ErrInvalidSegment)
				return
			}
		}

		if in.User.TenantCode != "" {
			tenant := rqAuth.NewTenants(d.AuthOltp)
			tenant.TenantCode = in.User.TenantCode
			if !tenant.FindByTenantCode() {
				out.SetError(400, ErrTenantCodeNotFound)
				return
			}
			user.SetTenantCode(in.User.TenantCode)
		}

		if user.HaveMutation() {
			user.SetUpdatedAt(in.UnixNow())
			user.SetUpdatedBy(sess.UserId)
			if user.Id == 0 {
				user.SetCreatedAt(in.UnixNow())
			}
		}

		if !user.DoUpsert() {
			out.SetError(500, ErrUserSaveFailed)
		}

		user.CensorFields()
		out.User = &user.Users

		if in.Pager.Page == 0 {
			break
		}

		fallthrough
	case zCrud.CmdList:
		r := rqAuth.NewUsers(d.AuthOltp)
		out.Users = r.FindByPagination(&SuperAdminUserManagementMeta, &in.Pager, &out.Pager)
	}
	return
}
