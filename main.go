package main

import (
	"fmt"
	"groupie/functions"
	"html/template" //text
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// hfgwhkuef

type Final struct {
	ID        int
	Image     string
	Artist    string
	Members   functions.Members
	AlbumYear int
	Album1    string
	Locations []string
}


func main() {
	// handler functions
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.HandleFunc("/homepagesearch", homepagesearch)
	http.HandleFunc("/result", result)
	http.HandleFunc("/mappage", mappage)
	http.HandleFunc("/", homepage)
	http.ListenAndServe(":8080", nil)
}
//function for waitingn pagge between google maps annd result
func mappage(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {

	city := r.FormValue("location")
	fmt.Println("city:",city)
cityinfo := strings.Split(city,":")
places := strings.Split(string(cityinfo[0]),", ")
place := strings.Split(places[1],":")
cit := strings.Title(places[0])
count := strings.Title(place[0])
fmt.Println("cit:",cit)
fmt.Println("count:",count)

mapdata,err := functions.MapsApiResp(cit,count)
if err != nil {
	fmt.Println("err:",err)
}

fmt.Println("mapdata:",mapdata)
  tmpl,_ := template.ParseFiles("mappage.html")
  tmpl.Execute(w, mapdata)
  return

	}
}

func homepagesearch(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	// where form value is collected for artist name annd fed innto the relavent funnction
	artIdsug := r.FormValue("suggestions[]")
	var ids []int
	character, _ := functions.LoadData("https://groupietrackers.herokuapp.com/api/artists")
	charData, err := functions.LoadUrelles("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		http.Error(w, "Failed to load relations data", http.StatusInternalServerError)
		return
	}


	// var finalResults []Final
	var locations []string

	for i := 0; i < len(charData); i++ {
		for location := range charData[i].DatesLocations {
			locations = append(locations, location)
		}
	}
		// for i := 0; i < len(character); i++ {
		// 	final := Final{
		// 		ID:        character[i].ID,
		// 		Image:     character[i].Image,
		// 		Artist:    character[i].Artist,
		// 		AlbumYear: character[i].AlbumYear,
		// 		Album1:    character[i].Album1,
		// 		Members:   character[i].Members,
		// 		Locations: locations,
		// 	}
		// 	finalResults = append(finalResults, final)
		// }

		if strings.Contains(artIdsug, "-") {
			artIdsug = (strings.Split(artIdsug, ":")[1])
			fmt.Println(artIdsug)
			id, _ := strconv.Atoi(artIdsug)
			ids = append(ids, id)
		} else {
			artIdsug = strings.ToLower(artIdsug)

			for _, data := range character {

				if strings.Contains(string(data.Artist), strings.Title(artIdsug)) {
					fmt.Print(string(data.Artist))
					ids = append(ids, data.ID)
				}
				if strings.Contains(string(data.Album1), strings.Title(artIdsug)) {
					fmt.Print(string(data.Album1))
					ids = append(ids, data.ID)
				}
				if strings.Contains(strconv.Itoa(data.AlbumYear), strings.Title(artIdsug)) {
					fmt.Print(string(data.AlbumYear))
					ids = append(ids, data.ID)
				}

				for i := 0; i < len(data.Members); i++ {
					if strings.Contains(string(data.Members[i]), strings.Title(artIdsug)) {
						fmt.Print(string(data.Members[i]))
						ids = append(ids, data.ID)
					}
				}
			}

			for _, dat := range charData {
				for location, _ := range dat.DatesLocations {
					if strings.Contains(location, artIdsug) {
						fmt.Printf("Location matched: %s\n", location)
						ids = append(ids, dat.ID)
					}
				}
			}
		}

		ids = removeDuplicates(ids)
fmt.Println("ids;",ids)

		var FFFinal []Final
		for i := 0; i < len(ids); i++ {
			final := Final{
				ID:        character[ids[i]-1].ID,
				Image:     character[ids[i]-1].Image,
				Artist:    character[ids[i]-1].Artist,
				Members:   character[ids[i]-1].Members,
				AlbumYear: character[ids[i]-1].AlbumYear,
				Album1:    character[ids[i]-1].Album1,
				Locations: locations,
			}
			FFFinal = append(FFFinal, final)
		}
fmt.Println("FFinal:",FFFinal)
		t, err := template.ParseFiles("indexs.html")
		if err != nil {
			http.Error(w, "Error parsing html", http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, FFFinal)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}


func homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderErrorPage(w, 404)
		return
	}
	character, _ := functions.LoadData("https://groupietrackers.herokuapp.com/api/artists")
	charData, _ := functions.LoadUrelles("https://groupietrackers.herokuapp.com/api/relation")

	var locations []string

	var finalResults []Final
	for i := 0; i < len(charData); i++ {
		for location := range charData[i].DatesLocations {
			locations = append(locations, location)
		}
	}
		for i := 0; i < len(character); i++ {
			final := Final{
				ID:        character[i].ID,
				Image:     character[i].Image,
				Artist:    character[i].Artist,
				Members:   character[i].Members,
				AlbumYear: character[i].AlbumYear,
				Album1:    character[i].Album1,
				Locations: locations,
			}
			finalResults = append(finalResults, final)
		}

		t, err := template.ParseFiles("index.html")
		if err != nil {
			renderErrorPage(w, 500)
			return
		}

		err = t.Execute(w, finalResults)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}


func result(wr http.ResponseWriter, r *http.Request) {
	// For resultpage the request is always POST not GET
	if r.Method != http.MethodPost {
		renderErrorPage(wr, 405)
		return
	}

	// if url is not for result page error handle
	if r.URL.Path != "/result" {
		renderErrorPage(wr, 404)
		return
	}

	char, _ := functions.LoadData("https://groupietrackers.herokuapp.com/api/artists")
	artId := r.FormValue("artist")
	for _, ch := range artId {
		if ch != 10 && ch != 13 && (ch < 32 || ch > 126) {
			renderErrorPage(wr, 400)
			return
		}
	}
	artId = strings.Title(artId)

	if len(artId) > 2 {
		pattern := regexp.MustCompile(`:\d+`)
		artId = strings.TrimPrefix(pattern.FindString(artId), ":")
	}

	for i := 0; i < 52; i++ {
		if artId == char[i].Artist {
			artId = strconv.Itoa(char[i].ID)
		}
	}

	iint, err := strconv.Atoi(artId)
	if err != nil || iint <= 0 {
		renderErrorPage(wr, 500)
		return
	}
	i := iint - 1

	// Load artist data
	character, err := functions.LoadData("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(wr, "Failed to load artist data", http.StatusInternalServerError)
		return
	}

	if len(character) == 0 {
		http.Error(wr, "No artist data available", http.StatusInternalServerError)
		return
	}

	charData, err := functions.LoadUrelles("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		http.Error(wr, "Failed to load relations data", http.StatusInternalServerError)
		return
	}

	if len(charData) == 0 {
		http.Error(wr, "No data available", http.StatusInternalServerError)
		return
	}

	// members := "No members available"
	// if len(character[i].Members) > 0 {
	// 	members = strings.Join(character[i].Members, ", ")
	// }

	var cdata []string
	x := 1
	d := ""
	for location, date := range charData[i].DatesLocations {
		d =  strings.ReplaceAll(string(strings.Title(location)), "-", ", ") + ": " + strings.Join(date, ", ")
		d = strings.ReplaceAll(d, "_", " ")
		cdata = append(cdata, d)
		x++
	}

	FFinal := Final{
		ID:        character[i].ID,
		Image:     character[i].Image,
		Artist:    character[i].Artist,
		Members:   functions.Members{},
		AlbumYear: character[i].AlbumYear,
		Album1:    character[i].Album1,
		Locations: cdata,
	}

	MMember := []string{}
	MMember = append(MMember, character[i].Members...)

	FFinal.Members = MMember

	t, err := template.ParseFiles("result.html")
	if err != nil {
		renderErrorPage(wr, 500)
		return
	}

	errr := t.Execute(wr, FFinal)
	if errr != nil {
		log.Println(errr)
		http.Error(wr, "Error executing template", http.StatusInternalServerError)
	}
}

func renderErrorPage(w http.ResponseWriter, code int) {

	// Set the status coded
	if code == 500 {
		w.WriteHeader(http.StatusInternalServerError)
	} else if code == 404 {
		w.WriteHeader(http.StatusNotFound)
	} else if code == 400 {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Generate ASCII art for the error code with the "Standard" style

	// Parse and render the custom 404 template
	t, err := template.ParseFiles(fmt.Sprintf("style/%d.html", code))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing %d HTML", code), http.StatusInternalServerError)
		return
	}

	// Render the template with the result
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing %d template", code), http.StatusInternalServerError)
	}
}

func removeDuplicates(old []int) []int {
	var new []int

	if len(old) == 1 {
		return old
	}

	for i := 0; i < len(old); i++ {
		if !Seen(old[i], new) {
			new = append(new, old[i])
		}
	}
	return new
}

func Seen(n int, new []int) bool {
	for i := 0; i < len(new); i++ {
		if n == new[i] {
			return true
		}
	}
	return false
}
