package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

var (
	SCYLLADB_URL = "172.17.0.3:9042"
	CLUSTER_KEYSPACE = os.Getenv("CLUSTER_KEYSPACE")
)

func InitDatabase() {
	CreateKeypace()
	CreateTables()
}

func CreateKeypace() {
	cluster := gocql.NewCluster(SCYLLADB_URL)
	session, _ := cluster.CreateSession()
	defer session.Close()
	
	session.Query(`DROP KEYSPACE IF EXISTS hermes`)
	session.Query(`CREATE KEYSPACE hermes WITH replication = {'class' : 'SimpleStrategy','replication_factor' : 1}`).Exec()

}

func CreateTables() {
	cluster := gocql.NewCluster(SCYLLADB_URL)
	cluster.Keyspace = CLUSTER_KEYSPACE
	session, err := cluster.CreateSession()

	if err != nil {
		log.Println("Couldn't connect to ScyllaDB")
		return
	}
	defer session.Close()

	PUBLISHER_TABLE := "CREATE TABLE IF NOT EXISTS hermes.publisher (id text PRIMARY KEY)"

	if err := session.Query(PUBLISHER_TABLE).Exec(); err != nil {
		log.Printf("[DATABASE] - %v - Error while creating hermes.publisher table\n", time.Now())
	}

	PAGEVIEW_TABLE := "CREATE TABLE IF NOT EXISTS hermes.pageView (id UUID PRIMARY KEY, data text)"

	if err := session.Query(PAGEVIEW_TABLE).Exec(); err != nil  {
		log.Printf("[DATABASE] - %v - Error while creating hermes.pageView table\n", time.Now())
		return
	}
}

func InsertObject(topic string, body gjson.Result) {
	cluster := gocql.NewCluster(SCYLLADB_URL)
	cluster.Keyspace = CLUSTER_KEYSPACE
	session, _ := cluster.CreateSession()
	defer session.Close()

	QUERY := fmt.Sprintf("INSERT INTO %s (id, data) VALUES (?,?)", topic)

	if err := session.Query(QUERY, gocql.TimeUUID(), body.Raw).Exec(); err != nil  {
		log.Printf("[DATABASE] - %v - Error while executing insert event to database\n", time.Now())
		return
	}
}

func InsertPublisher(id uuid.UUID) {
	cluster := gocql.NewCluster(SCYLLADB_URL)
	cluster.Keyspace =CLUSTER_KEYSPACE
	session, _ := cluster.CreateSession()
	defer session.Close()

	QUERY := "INSERT INTO hermes.publisher (id) VALUES (?)"

	if err := session.Query(QUERY, id.String()).Exec(); err != nil {
		log.Printf("[DATABASE] - %v - Error while inserting new publisher to database\n", time.Now())
		return
	}
}

func FindPublisherByUUID(id string) bool {
	cluster := gocql.NewCluster(SCYLLADB_URL)
	cluster.Keyspace = CLUSTER_KEYSPACE
	session, _ := cluster.CreateSession()
	defer session.Close()

	QUERY := "SELECT id FROM hermes.publisher WHERE id = ?"
	iter := session.Query(QUERY, id).Iter()

	return iter.NumRows() != 0

}