package zCrud

const (
	CmdList    = `list`    // retrieve rows for table
	CmdForm    = `form`    // retrieve 1 row for update
	CmdUpsert  = `upsert`  // insert if id=0, update if id>0
	CmdRestore = `restore` // same as upsert but also unset deletedAt, deletedBy
	CmdDelete  = `delete`  // same as upsert but also set deletedAt, deletedBy
)
