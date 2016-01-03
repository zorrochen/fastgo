package main

import (
	"encoding/xml"
	"log"
)

const (
	DIR_MODULE = "./module/"
	DIR_MYSQL = "./mysql/"
)

const (
	FILE_MYSQL_FMT = "mysql.xml"
	FILE_MYSQL_OUT = "mysql.go"
)

type parse_and_write_t struct {
}


func (obj *parse_and_write_t) Parse() error {
	log.Printf("parse start!")

	fileList, err := FilesInDirection(*moduleDir)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	for _, fileName := range fileList {
		err := parse(fileName)
		if err != nil {
			log.Printf("Error:%v", err)
			return err
		}
	}

	err = parseMysql(FILE_MYSQL_FMT)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	log.Printf("parse end!")
	return nil
}

func parse(fileName string) error {
	srcDat, err := readFile(DIR_MODULE + fileName)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	var data module_fmt_t
	err = xml.Unmarshal(srcDat, &data)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	srcFileStr, err := genModuleStr(&data)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	err = writeFile(data.Name, srcFileStr)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}
	return nil
}


func parseMysql(fileName string) error {
	srcDat, err := readFile(DIR_MYSQL + fileName)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	var data tables_fmt_t
	err = xml.Unmarshal(srcDat, &data)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}

	srcMysqlStr := ""
	for _, tableInfo := range data.TablesInfo {
		str, err := global_server_ref.tranSqlObj.genMysql(&tableInfo)
		if err != nil {
			log.Printf("Error:%v", err)
			return err
		}
		srcMysqlStr += str
		srcMysqlStr += "\n"
	}

	err = writeFile(FILE_MYSQL_OUT, srcMysqlStr)
	if err != nil {
		log.Printf("Error:%v", err)
		return err
	}
	return nil
}
