package main 
 
import ( 
	"fmt" 
	"github.com/gocolly/colly"
	"database/sql" 
	"strings"
	"strconv"
	 _ "github.com/mattn/go-sqlite3"
) 


const dictUrl = "http://www.nplg.gov.ge/gwdict/index.php?a=list&d=9&p="

func main() { 
	const file string = "dict.db"
	db, err := sql.Open("sqlite3", file)
	var count int = 0
	var page int = 1
	var offset int = 0

	if err != nil {
  		fmt.Println("Database opening error ", err) 
  		return
	}

	//db.Exec("DELETE FROM words;")

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) { 
		fmt.Println("Visiting: ", r.URL) 
	}) 
	 
	c.OnError(func(_ *colly.Response, err error) { 
		fmt.Println("Something went wrong: ", err) 
	}) 
	 
	c.OnResponse(func(r *colly.Response) { 
		fmt.Println("Page visited: ", r.Request.URL) 
	}) 
	 
	c.OnHTML("dt.termpreview", func(e *colly.HTMLElement) { 		
		text := strings.Trim(e.ChildText("a"), " ")
		text = strings.Trim(text, ",")
		id := e.Index + offset
		fmt.Printf("%v %v: \n", id, text) 
		db.Exec("INSERT INTO words (id, ru) VALUES(?, ?);", id, text)
	}) 
	 
	c.OnHTML("dd.defnpreview", func(e *colly.HTMLElement) { 
		var wtype int = 0
		text := strings.Trim(e.Text, " ")
		id := e.Index + offset
		db.Exec("UPDATE words SET ge = ?, wtype = ? WHERE id = ?;", text, wtype, id)
		count++
	}) 

	c.OnScraped(func(r *colly.Response) { 
		fmt.Println(r.Request.URL, " scraped!") 
		page++;
		if (page == 2289) {
			db.Close()
		} else {			
			offset += count
			count = 0
			c.Visit(dictUrl + strconv.Itoa(page))
		}
	})

	c.Visit(dictUrl + strconv.Itoa(page))

}