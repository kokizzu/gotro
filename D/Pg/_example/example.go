package main

import (
	"fmt"
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
)

func main() {
	// Get Connection
	pg_conn := Pg.NewConn(`gotro_example`, `gotro_example`)

	// Create Table
	pg_conn.CreateBaseTable(`users`)

	// Insert 3 rows
	fmt.Println("---Insert three rows---")
	var row_ids [3]int64
	row_ids[0] = pg_conn.DoInsert(1, `users`, M.SX{
		`data`: M.ToJson(M.SX{
			`name`: `alpha`,
			`age`:  12,
		}),
	})
	row_ids[1] = pg_conn.DoInsert(1, `users`, M.SX{
		`data`: M.ToJson(M.SX{
			`name`: `beta`,
			`age`:  23,
		}),
	})
	row_ids[2] = pg_conn.DoInsert(1, `users`, M.SX{
		`data`: M.ToJson(M.SX{
			`name`: `gamma`,
			`age`:  15,
		}),
	})

	// Query all rows
	fmt.Println("\n---Query all rows---")
	rows := pg_conn.QMapArray(`SELECT * FROM users WHERE is_deleted = false`)
	for _, row := range rows {
		for k, v := range row {
			fmt.Printf("%s: %v\n", k, v)
		}
		fmt.Println()
	}

	// Update the third row that we've just inserted
	fmt.Println("---Update the third row that we've just inserted---")
	pg_conn.DoUpdate(1, `users`, row_ids[2], M.SX{
		`data`: M.ToJson(M.SX{
			`name`: `delta`,
			`age`:  18,
		}),
	})

	// Query the updated row
	fmt.Println("\n---Query the updated row---")
	rows = pg_conn.QMapArray(`SELECT * FROM users WHERE id = ` + S.ZI(row_ids[2]) + ` AND is_deleted = false`)
	for _, row := range rows {
		for k, v := range row {
			fmt.Printf("%s: %v\n", k, v)
		}
		fmt.Println()
	}

	// Delete all rows
	fmt.Println("\n---Delete all three rows that we've just inserted---")
	ids := pg_conn.QIntArr(`SELECT id FROM users WHERE is_deleted = false`)
	for _, id := range ids {
		fmt.Printf("Delete row with id = %d\n", id)
		pg_conn.DoDelete(1, `users`, id)
	}

	// Query all rows that are not deleted
	fmt.Println("\n---Query all rows---")
	rows = pg_conn.QMapArray(`SELECT * FROM users WHERE is_deleted = false`)
	for _, row := range rows {
		for k, v := range row {
			fmt.Printf("%s: %v\n", k, v)
		}
		fmt.Println()
	}
}
