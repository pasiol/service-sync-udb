package internal

import (
	"context"
	"log"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	mssqlutils "github.com/pasiol/go-mssql-utils"
	"github.com/pasiol/pq"

	"github.com/pasiol/service-sync-udb/configs"
)

var (
	debug = false
)

func SyncIds() {
	table := configs.SQLTable
	debug, debugExists := os.LookupEnv("APP_DEBUG")
	if debugExists {
		if debug == "true" {
			table = configs.SQLDebugTable
		}
	}
	studenstWithoutPersonalId := getStudentsWithoutIds(table)
	log.Printf("founded %d students without personal ids", len(studenstWithoutPersonalId))
	sqlDb := mssqlutils.ConnectOrDie(configs.SQLServer, configs.SQLPort, configs.SQLUser, configs.SQLPassword, configs.SQLDb, true, true)
	defer sqlDb.Close()
	ctx := context.Background()
	err := sqlDb.PingContext(ctx)
	if err != nil {
		log.Fatalf("sql connection error: %s", err)
	}
	for _, s := range studenstWithoutPersonalId {
		personalId, personalEmail, studentId := lookupPrimus(s.Id)
		if len(personalId) > 0 {
			tsql := configs.GetUpdateStatementSQL(s.Id, personalId, personalEmail, studentId)
			_, err := sqlDb.QueryContext(ctx, tsql)
			if err != nil {
				log.Printf("updating personal id of student %d failed: %s", s.Id, err)
			} else {
				log.Printf("updating personal id of student %d succeed", s.Id)
			}
		} else {
			log.Printf("student %d has no personal id", s.Id)
		}
	}
}

func lookupPrimus(id int64) (string, string, string) {
	pq.Debug = true

	query := configs.StudentQuery(id)
	output := pq.ExecuteAndRead(query, 30)
	rows := strings.Split(output, "\n")
	if len(rows) == 2 {
		updatedData := strings.Split(rows[0], ";")
		if len(updatedData) == 4 {
			return updatedData[1], updatedData[2], updatedData[3]
		}
	}
	return "", "", ""
}

func getStudentsWithoutIds(t string) []Student {
	var students = []Student{}
	sqlDb := mssqlutils.ConnectOrDie(configs.SQLServer, configs.SQLPort, configs.SQLUser, configs.SQLPassword, configs.SQLDb, true, true)
	defer sqlDb.Close()
	ctx := context.Background()
	err := sqlDb.PingContext(ctx)
	if err != nil {
		log.Fatalf("sql connection error: %s", err)
	}
	tsql := configs.StudentsWithoutIdsSQL
	rows, err := sqlDb.QueryContext(ctx, tsql)
	if err != nil {
		if debug {
			log.Fatalf("tsql: %s", tsql)
		}
		log.Fatal(err.Error())
	}
	for rows.Next() {
		var s Student
		err := rows.Scan(&s.Id, &s.PersonalId, &s.personalEmail)
		if err != nil {
			log.Fatalf("reading sql row failed: %s", err)
		}
		students = append(students, s)
		s = Student{}
	}
	return students
}