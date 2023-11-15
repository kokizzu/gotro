package domain

type (
	SuperAdminDashboardIn struct {
		RequestCommon
	}

	SuperAdminDashboardOut struct {
		ResponseCommon
	}
)

const (
	SuperAdminDashboardAction = `superAdmin/dashboard`
)

func (d *Domain) SuperAdminDashboard(in *SuperAdminDashboardIn) (out SuperAdminDashboardOut) {
	defer d.InsertActionLog(&in.RequestCommon, &out.ResponseCommon)

	sess := d.MustSuperAdmin(in.RequestCommon, &out.ResponseCommon)
	if sess == nil {
		return
	}

	// TODO: implement superadmin dashboard
	return
}
