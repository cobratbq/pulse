package pulse

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// Dial establishes a new connection to a pulse database.
func Dial(connectString string, db string, collection string) (*Connection, error) {
	session, err := mgo.Dial(connectString)
	if err != nil {
		return nil, err
	}
	c := session.DB(db).C(collection)
	return &Connection{session, c}, nil
}

// Connection represents the connection to the pulse database.
type Connection struct {
	session *mgo.Session
	c       *mgo.Collection
}

// Record records a new pulse in the database.
func (conn *Connection) Record(namespace string) error {
	return conn.c.Insert(&Pulse{Namespace: namespace, Time: bson.Now()})
}

// Get retrieves all pulses for the given namespace.
func (conn *Connection) Get(namespace string) ([]Pulse, error) {
	q := conn.c.Find(bson.M{"namespace": namespace}).Sort("time")
	var pulses []Pulse
	err := q.All(&pulses)
	if err != nil {
		return nil, err
	}
	return pulses, nil
}

// Close closes the connection.
func (conn *Connection) Close() {
	conn.session.Close()
}

// Pulse represents a single pulse.
type Pulse struct {
	Namespace string
	Time      time.Time
}
