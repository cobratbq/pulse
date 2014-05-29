package pulse

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// Dial establishes a new connection to a pulse database.
func Dial(connectString string) (*Connection, error) {
	session, err := mgo.Dial(connectString)
	if err != nil {
		return nil, err
	}
	return &Connection{session}, nil
}

// Connection represents the connection to the pulse database.
type Connection struct {
	session *mgo.Session
}

// Record records a new pulse in the database.
func (conn *Connection) Record(namespace string) error {
	c := conn.session.DB("pulse").C("pulses")
	return c.Insert(&Pulse{Namespace: namespace, Time: bson.Now()})
}

// Get retrieves all pulses for the given namespace.
func (conn *Connection) Get(namespace string) ([]Pulse, error) {
	c := conn.session.DB("pulse").C("pulses")
	q := c.Find(bson.M{"namespace": namespace}).Sort("time")
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
