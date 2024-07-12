package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-adodb"
)

type SageArtikel struct {
	Id            int
	Artikelnummer string
	Suchbegriff   string
	Preis         float64
}

type AccessArtikel struct {
	Id            int
	Artikelnummer string
	Artikeltext   string
	Preis         float64
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var sage_server = os.Getenv("SAGE_SERVER")
	var sage_password = os.Getenv("SAGE_PASSWORD")
	var sage_user = os.Getenv("SAGE_USER")
	var sage_db = os.Getenv("SAGE_DB")
	var access_db = os.Getenv("ACCESS_DB")

	sage := readSage(sage_server, sage_user, sage_password, sage_db)
	label := readAccessDb(access_db)
	syncDb(sage, label, access_db)
	fmt.Println("Fertig")
}

func syncDb(sage []SageArtikel, label []AccessArtikel, access_db string) {
	var update []AccessArtikel
	var create []AccessArtikel

	for i := range sage {
		var found bool
		found = false
		for x := range label {
			if sage[i].Id == label[x].Id {
				found = true
				break
			}
		}
		var art AccessArtikel
		art.Id = sage[i].Id
		art.Artikelnummer = sage[i].Artikelnummer
		art.Preis = sage[i].Preis
		art.Artikeltext = sage[i].Suchbegriff
		if found {
			update = append(update, art)
		} else {
			create = append(create, art)
		}
	}

	insert(create, access_db)
	updatefunc(update, access_db)

}

func updatefunc(create []AccessArtikel, access_db string) {
	conn, err := sql.Open("adodb", fmt.Sprintf("Provider=Microsoft.ACE.OLEDB.12.0;Data Source=%s;", access_db))
	if err != nil {
		log.Fatal("syncDB: Connection failed: ", err)
	}
	defer conn.Close()

	stmt, err := conn.Prepare("UPDATE Artikel SET Artikelnummer=?, Artikeltext=?, Preis=? where ID=?")
	if err != nil {
		log.Fatal("syncDb: Insert Prepare failed: ", err)
	}
	defer stmt.Close()

	for x := range create {
		if _, err := stmt.Exec(create[x].Artikelnummer, strings.ReplaceAll(create[x].Artikeltext, "'", "\""), create[x].Preis, create[x].Id); err != nil {
			log.Fatal("syncDb: Insert Exec failed: ", err)
		}

	}
}

func insert(create []AccessArtikel, access_db string) {
	conn, err := sql.Open("adodb", fmt.Sprintf("Provider=Microsoft.ACE.OLEDB.12.0;Data Source=%s;", access_db))
	if err != nil {
		log.Fatal("syncDB: Connection failed: ", err)
	}
	defer conn.Close()

	stmt, err := conn.Prepare("INSERT INTO Artikel (ID, Artikelnummer, Artikeltext, Preis) VALUES (?,?,?,?)")
	if err != nil {
		log.Fatal("syncDb: Insert Prepare failed: ", err)
	}
	defer stmt.Close()

	for x := range create {
		if _, err := stmt.Exec(create[x].Id, create[x].Artikelnummer, strings.ReplaceAll(create[x].Artikeltext, "'", "\""), create[x].Preis); err != nil {
			log.Fatal("syncDb: Insert Exec failed: ", err)
		}

	}
}

func readAccessDb(access_db string) []AccessArtikel {
	conn, err := sql.Open("adodb", fmt.Sprintf("Provider=Microsoft.ACE.OLEDB.12.0;Data Source=%s;", access_db))
	if err != nil {
		log.Fatal("readAccessDb: Connection failed: ", err)
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT ID, Artikelnummer, Artikeltext, Preis FROM Artikel")
	if err != nil {
		log.Fatal("readAccessDb: Query failed: ", err)
	}
	defer rows.Close()

	var artikel []AccessArtikel

	for rows.Next() {
		var art AccessArtikel
		if err := rows.Scan(&art.Id, &art.Artikelnummer, &art.Artikeltext, &art.Preis); err != nil {
			log.Fatal("readAccessDb: Scan failed: ", err)
		}
		artikel = append(artikel, art)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("readAccessDb: Row error: ", err)
	}

	return artikel
}

func readSage(sage_server string, sage_user string, sage_password string, sage_db string) []SageArtikel {
	sage_port, err := strconv.ParseInt(os.Getenv("SAGE_PORT"), 0, 64)
	if err != nil {
		log.Fatal("readSage: SAGE_PORT not in .env: ", err)
	}
	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d", sage_server, sage_db, sage_user, sage_password, sage_port)

	conn, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("readSage: Connection failed: ", err)
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT sg_auf_artikel.SG_AUF_ARTIKEL_PK, sg_auf_artikel.ARTNR, sg_auf_artikel.SUCHBEGRIFF, sg_auf_vkpreis.PR01 FROM sg_auf_artikel INNER JOIN sg_auf_vkpreis ON (sg_auf_artikel.SG_AUF_ARTIKEL_PK = sg_auf_vkpreis.SG_AUF_ARTIKEL_FK)")
	if err != nil {
		log.Fatal("readSage: Query failed: ", err)
	}
	defer rows.Close()

	var artikel []SageArtikel

	for rows.Next() {
		var art SageArtikel
		var Artikelnummer sql.NullString
		var Suchbegriff sql.NullString
		var Price sql.NullFloat64

		if err := rows.Scan(&art.Id, &Artikelnummer, &Suchbegriff, &Price); err != nil {
			log.Fatal("readSage: Scan failed: ", err)
		}
		if Artikelnummer.Valid && Suchbegriff.Valid && Price.Valid {
			art.Artikelnummer = Artikelnummer.String
			art.Suchbegriff = Suchbegriff.String
			art.Preis = Price.Float64
			artikel = append(artikel, art)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("readSage: Row Error: ", err)
	}
	return artikel
}
