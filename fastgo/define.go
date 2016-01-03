package main

const MYSQL_FMT_QUERY =
`func table%sQuery() ([]%s_t, error) {
	sqlcmd := "%s"

	rows, err := mysql.db.Query(sqlcmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	retList := []%s_t{}
	for rows.Next() {
		var retInfo %s_t
		err = rows.Scan(%s)
		if err != nil {
			return nil, err
		}
		retList = append(retList, retInfo)
	}

	return retList, nil
}

`

const MYSQL_FMT_INSERT =
`func table%sInsert(obj %s_t) (error) {
	sqlcmd := "%s"
	_, err := mysql.db.Exec(sqlcmd)
	if err != nil {
		return err
	}

	return nil
}

`
