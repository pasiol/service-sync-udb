package internal

import (
	"context"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	mssqlutils "github.com/pasiol/go-mssql-utils"
	pq "github.com/pasiol/gopq"

	"github.com/pasiol/service-sync-udb/configs"
)

var (
	debug = false
)

func SyncIds() {

	var (
		succeed = 0
		failed  = 0
	)

	studenstWithoutPersonalId := getStudentsWithoutIds()
	total := len(studenstWithoutPersonalId)
	log.Printf("founded %d students without personal ids", total)
	sqlDb := mssqlutils.ConnectOrDie(configs.SQLServer, configs.SQLPort, configs.SQLUser, configs.SQLPassword, configs.SQLDb, true, true)
	defer sqlDb.Close()
	ctx := context.Background()
	err := sqlDb.PingContext(ctx)
	if err != nil {
		log.Fatalf("sql connection error: %s", err)
	}
	err = pq.UpdatePQ(configs.GetPrimusConfig().PrimusHost, configs.GetPrimusConfig().PrimusPort)
	if err != nil {
		log.Printf("PQ update failed: %s", err)
	}
	for _, s := range studenstWithoutPersonalId {
		personalId, personalEmail, studentId, err := lookupPrimus(s.Id)
		if err != nil {
			log.Printf("getting student %s primus data failed: %s", studentId, err)
			failed = failed + 1
		} else {
			if len(personalId) > 0 {
				tsql := configs.GetUpdateStatementSQL(s.Id, personalId, personalEmail, studentId)
				_, err := sqlDb.QueryContext(ctx, tsql)
				if err != nil {
					log.Printf("updating personal id of student %d failed: %s", s.Id, err)
					failed = failed + 1
				} else {
					log.Printf("updating personal id of student %d succeed", s.Id)
					succeed = succeed + 1
				}
			} else {
				log.Printf("student %d has no personal id", s.Id)
			}
		}
	}
	log.Printf("sync ended total count: %d, succeed count: %d, failed count: %d", total, succeed, failed)
}

func lookupPrimus(id int64) (string, string, string, error) {

	query := configs.StudentQuery(id)
	output, err := pq.ExecuteAndRead(query, 30)
	if err != nil {
		return "", "", "", err
	}
	rows := strings.Split(output, "\n")
	if len(rows) == 2 {
		updatedData := strings.Split(rows[0], ";")
		if len(updatedData) == 4 {
			return updatedData[1], updatedData[2], updatedData[3], nil
		}
	}
	return "", "", "", nil
}

func getStudentsWithoutIds() []Student {
	var students = []Student{}
	sqlDb := mssqlutils.ConnectOrDie(configs.SQLServer, configs.SQLPort, configs.SQLUser, configs.SQLPassword, configs.SQLDb, true, true)
	defer sqlDb.Close()
	ctx := context.Background()
	err := sqlDb.PingContext(ctx)
	if err != nil {
		log.Fatalf("sql connection error: %s", err)
	}
	tsql := configs.GetStudentsWithoutIdsSQL()
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
