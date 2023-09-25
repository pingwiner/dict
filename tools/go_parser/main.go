package main
  
import (
    "bufio"
    "fmt"
    "log"
    "strings"
    "os"
    "sort"
    "regexp"
    "strconv"
)

type Theme struct {
    id int
    name string
}

type Item struct {
	ru string
	ge string
    theme int
}

type ByGe []Item

func (this ByGe) Len() int {
    return len(this)
}
func (this ByGe) Less(i, j int) bool {
    return this[i].ge < this[j].ge
}
func (this ByGe) Swap(i, j int) {
    this[i], this[j] = this[j], this[i]
}

var items []Item
var themes []Theme

const enChars string = "abcdefghijklmnopqrstuvwxyz"
const ruChars string = "абвгдежзийклмнопрстуфхцчшщьыъэюя"
const geChars string = "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ" 

func isCharFrom(ch rune, charSet string) bool {
    r := []rune(charSet)
    for i := 0; i < len(r); i++ {
        if (r[i] == ch) {
            return true
        }
    }
    return false
} 

func getLangText(text string, charSet string, startFrom int) (string, int) {
    r := []rune(text)

    for i := startFrom; i < len(r); i++ {
        if (isCharFrom(r[i], charSet)) {
            return strings.Trim(string(r[startFrom:i]), " "), i
        }
    }
    if (startFrom == 0) {
        return "", 0
    } 
    return strings.Trim(string(r[startFrom:]), " "), len(r)
} 


func tailString(text string, i int) string {
    r := []rune(text)    
    return string(r[i:])
}

func itemsContains(key string) bool {
    for _, v := range items {
        if v.ge == key {
            return true
        }
    }
    return false
}

func printItems() {
    for _, v := range items {
        fmt.Printf("aw('%s','%s', %d);\n", v.ge, v.ru, v.theme)
    }
}

func printThemes() {
    for _, v := range themes {
        fmt.Printf("th(%d,'%s');\n", v.id, v.name)
    }    
}

func parseTheme(text string) int {
    r := []rune(text)
    i := 0
    
    for r[i] != '.' {
        i++
    }
    
    index, err := strconv.Atoi(string(r[:i]))
    if err != nil {
        panic(err)
    }

    themeName := strings.Trim(string(r[i+1:]), " ")
    theme := Theme {index, themeName}
    themes = append(themes, theme)

    return index
}

func loadInput1() {
    file, err := os.Open("input.txt")

    if err != nil {
        log.Fatalf("failed to open")
  
    }
  
    scanner := bufio.NewScanner(file)  
    scanner.Split(bufio.ScanLines)
  
    r, _ := regexp.Compile("[0-9]+\\.")
    var currentTheme int = 0

    for scanner.Scan() {
        text := strings.ToLower(scanner.Text())

        if (r.MatchString(text)) {
            currentTheme = parseTheme(text)
        }

        ruText, i := getLangText(text, geChars, 0)
        if (i > 0) {
            geText, _ := getLangText(text, enChars, i)
            item := Item {ruText, geText, currentTheme}
            items = append(items, item)
        }       
    }
  
    file.Close()    
}

func loadInput2() {
    file, err := os.Open("input2.txt")

    if err != nil {
        log.Fatalf("failed to open")
  
    }
  
    scanner := bufio.NewScanner(file)  
    scanner.Split(bufio.ScanLines)
  
    for scanner.Scan() {
        text := strings.ToLower(scanner.Text())
        geText, i := getLangText(text, ruChars, 0)
        if (i > 0) {
            ruText := tailString(text, i)
            item := Item {ruText, geText, 0}
            if (itemsContains(geText)) {
                //fmt.Println("duplicate: ", geText)
            } else {
                items = append(items, item)
            }
        }
    }
  
    file.Close()    
}


func main() {  
    loadInput1()
    loadInput2()
    sort.Sort(ByGe(items))
    printThemes()
    printItems()
}
