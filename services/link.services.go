package services

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/fredeom/shortlinks/db"
)

func NewServicesLink(formData FormData, lStore db.LinkStore) *ServicesLink {
	data := NewData()
	return &ServicesLink{
		FormData:  formData,
		Data:      data,
		LinkStore: lStore,
	}
}

type ServicesLink struct {
	FormData  FormData
	Data      Data
	LinkStore db.LinkStore
}

type Link struct {
	ID        int       `json:"id"`
	Full      string    `json:"full"`
	Short     string    `json:"short"`
	Hits      string    `json:"hits"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Visitor struct {
	ID        int       `json:"id"`
	LinkId    int       `json:"link_id"`
	Agent     string    `json:"agent"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (sl *ServicesLink) GetFormData() (FormData, error) {
	return sl.FormData, nil
}

func (sl *ServicesLink) NewLink(full, short string) Link {
	stmt, _ := sl.LinkStore.Db.Prepare("INSERT INTO links(full, short)VALUES('" + full + "', '" + short + "')")
	stmt.Exec()
	defer stmt.Close()
	return Link{
		Full:      full,
		Short:     short,
		ID:        sl.HasLink(full),
		Hits:      "0",
		CreatedAt: time.Now(),
	}
}

func (sl *ServicesLink) GetLink(id int) Link {
	stmt, _ := sl.LinkStore.Db.Prepare("SELECT * FROM links WHERE ID=" + strconv.Itoa(id))
	rows, _ := stmt.Query()
	defer rows.Close()
	if rows.Next() {
		var link Link
		rows.Scan(&link.ID, &link.Full, &link.Short, &link.Hits, &link.CreatedAt)
		return link
	}
	return Link{}
}

type Links = []Link
type Visitors = []Visitor

type Data struct {
	Links Links
}

func (sl *ServicesLink) GetData() (Data, error) {
	stmt, _ := sl.LinkStore.Db.Prepare("SELECT * from links WHERE 1=1 LIMIT 20")
	rows, _ := stmt.Query()
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		rows.Scan(&link.ID, &link.Full, &link.Short, &link.Hits, &link.CreatedAt)
		links = append(links, link)
	}

	var data = Data{Links: links}
	return data, nil
}

func NewData() Data {
	return Data{
		Links: []Link{},
	}
}

func (sl *ServicesLink) HasLink(s string) int {
	var ID int
	stmt, _ := sl.LinkStore.Db.Prepare("SELECT ID from links WHERE full LIKE '" + s + "' or short LIKE '" + s + "' LIMIT 1")
	rows, _ := stmt.Query()
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&ID)
		return ID
	}
	return -1
}

func (sl *ServicesLink) DeleteLink(link Link) {
	stmt, _ := sl.LinkStore.Db.Prepare("DELETE FROM links WHERE ID=" + strconv.Itoa(link.ID))
	stmt.Exec()
	defer stmt.Close()
}

func (sl *ServicesLink) UpdateHits(id int) {
	sql := "UPDATE links SET hits=hits+1 WHERE ID=" + strconv.Itoa(id) + ";"
	stmt, _ := sl.LinkStore.Db.Prepare(sql)
	stmt.Exec()
	defer stmt.Close()
}

func (sl *ServicesLink) AddVisitor(id int, userAgent string) {
	stmt, _ := sl.LinkStore.Db.Prepare("INSERT INTO visitors(link_id, agent) VALUES (" + strconv.Itoa(id) + ", '" + userAgent + "')")
	stmt.Exec()
	stmt.Close()
}

func (sl *ServicesLink) GetVisitors(id int) Visitors {
	stmt, _ := sl.LinkStore.Db.Prepare("SELECT * FROM visitors WHERE link_id=" + strconv.Itoa(id) + " LIMIT 10;")
	rows, _ := stmt.Query()
	defer rows.Close()

	var visitors Visitors
	for rows.Next() {
		var visitor Visitor
		rows.Scan(&visitor.ID, &visitor.LinkId, &visitor.Agent, &visitor.CreatedAt)

		visitors = append(visitors, visitor)
	}
	return visitors
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func NewFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
