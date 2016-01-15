package main
import (
	"fmt"
	"strings"
	"encoding/xml"
)

type tran_sql_t struct {

}

type tables_fmt_t struct {
	TablesInfo []table_fmt_t `xml:"tables>table"`
}

type table_fmt_t struct {
	Name string `xml:"name"`
	Columns []column_fmt_t `xml:"columns>column"`
	Pkey []string `xml:"primaryKey"`
}

type column_fmt_t struct {
	Name string `xml:"name,attr"`
	ColType string `xml:"type,attr"`
}


func (owner *tran_sql_t) tranXml(sqlStr string) (string, error) {
	maybeTbs := strings.Split(sqlStr, ";")

	tsf := tables_fmt_t{}
	for _, mt := range maybeTbs {
		if !strings.Contains(mt, "CREATE TABLE") {
			continue
		}

		a := strings.Split(strings.TrimSpace(mt), "` (")[0]
		b := strings.Split(strings.TrimSpace(a), "CREATE TABLE `")[1]

		var tf table_fmt_t
		tf.Name = strings.TrimSpace(b)

		c := strings.Split(strings.TrimSpace(mt), "` (")[1]
		d := strings.Split(strings.TrimSpace(c), "\n)")[0]
		e := strings.Split(strings.TrimSpace(d), ",\n")
		for _, f := range e {
			if strings.TrimSpace(f) == "" {
				continue
			}

			fs := strings.Split(strings.TrimSpace(f), " ")

			var cf column_fmt_t
			cf.Name = strings.TrimSpace(fs[0])
			cf.ColType = strings.TrimSpace(fs[1])
			tf.Columns = append(tf.Columns, cf)
		}

		tsf.TablesInfo = append(tsf.TablesInfo, tf)
	}

	tsfXmlBytes, err := xml.Marshal(tsf)
	if err != nil {
		return "", err
	}

	return string(tsfXmlBytes), nil
}

func (owner *tran_sql_t) genStruct(t *table_fmt_t) (string, error) {
	ret := ""
	ret += fmt.Sprintf("type %s_t struct{\n%s\n}", t.Name, genStructOfColumn(t.Columns))
	return ret, nil
}


func (owner *tran_sql_t) genQuery(t *table_fmt_t) (string, error) {
	cmd := fmt.Sprintf("select %s from %s where %s", genColumnStr(t.Columns), t.Name, genWhereStr(t.Pkey))
	rowsScanStr := genRowsScanStr(t.Columns)
	retFmt := MYSQL_FMT_QUERY

	return fmt.Sprintf(retFmt, t.Name, t.Name, cmd, t.Name, t.Name, rowsScanStr), nil
}

func (owner *tran_sql_t) genInsert(t *table_fmt_t) (string, error) {
	cmd := fmt.Sprintf("insert into %s (%s) values (%s)", t.Name, genColumnStr(t.Columns), genValueFmtStr(t.Columns))

	retFmt := MYSQL_FMT_INSERT
	return fmt.Sprintf(retFmt, t.Name, t.Name, cmd), nil
}

//func genUpdate(t *table_fmt_t) (string, error) {
//
//}
//
//func genDelete(t *table_fmt_t) (string, error) {
//
//}

func (owner *tran_sql_t) genMysql(t *table_fmt_t) (string, error) {
	ret := ""

	tmpStr := ""
	tmpStr, _ = owner.genStruct(t)
	ret += tmpStr
	ret += "\n"

	tmpStr, _ = owner.genQuery(t)
	ret += tmpStr
	ret += "\n"

	tmpStr, _ = owner.genInsert(t)
	ret += tmpStr
	ret += "\n"

	return ret, nil
}


//================================sub===================================

func genStructOfColumn(cols []column_fmt_t) (string) {
	ret := ""
	isFirstParam := true
	for _, col := range cols {
		if isFirstParam {
			isFirstParam = false
		} else {
			ret += "\n"
		}
		ret += fmt.Sprintf("\t%s %s", col.Name, col.ColType)
	}
	return ret
}

func genColumnStr(cols []column_fmt_t) (string) {
	ret := ""
	isFirstParam := true
	for _, col := range cols {
		if isFirstParam {
			isFirstParam = false
		} else {
			ret += ", "
		}
		ret += fmt.Sprintf("%s", col.Name)
	}
	return ret
}

func genWhereStr(pKey	[]string) (string) {
	if len(pKey) == 1 {
		return pKey[0]
	}

	return "NULL"
}

func genValueFmtStr(cols []column_fmt_t) (string) {
	ret := ""
	for _, col := range cols {
		fmtStr := ""
		if col.ColType == "bool" {
			fmtStr = "%t"
		} else if col.ColType == "string"{
			fmtStr = "'%s'"
		} else if col.ColType == "int"{
			fmtStr = "%d"
		} else if col.ColType == "float"{
			fmtStr = "%f"
		}
		ret += fmtStr
		ret += ","
	}
	return ret
}

func genRowsScanStr(cols []column_fmt_t) (string) {
	ret := ""
	isFirstParam := true
	for _, col := range cols {
		if isFirstParam {
			isFirstParam = false
		} else {
			ret += ", "
		}
		ret += fmt.Sprintf("&%s", col.Name)
	}
	return ret
}
