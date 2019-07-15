package main

import (
	"bufio"
	"fmt"
	"io"

	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"

	"github.com/ripx80/brewman/pkgs/recipe"
)

/*
ES stuff
curl --user elastic:changeme -XPOST http://node1:9200/recipes-m3/doc -H "Content-Type: application/json" -d @400_Meraner_Weizen.json
for i in $(ls):; do curl --user elastic:changeme -XPOST http://node1:9200/recipes-m3/doc -H "Content-Type: application/json" -d @$i; done

PUT hockey/_bulk?refresh
{"index":{"_id":1}}
{"first":"johnny","last":"gaudreau","goals":[9,27,1],"assists":[17,46,0],"gp":[26,82,1],"born":"1993/08/13"}
{"index":{"_id":2}}


https://beerandbrewing.com/beer-recipes/
https://beerrecipes.org/
https://www.kaggle.com/jtrofe/beer-recipes (https://www.brewersfriend.com/search/)
https://www.brewerydb.com/developers
https://www.brewerydb.com/developers/docs/ (API)
https://github.com/homebrewing/tapline
http://www.malt.io/

Parsing Amount error: Hopfen_VWH_1_Menge strconv.ParseFloat: parsing " 90": invalid syntax
json: invalid use of ,string struct tag, trying to unmarshal "" into float64
Rest value is empty: Infusion_Rastzeit4
Parsing Amount error: Infusion_Rastzeit2 strconv.Atoi: parsing "7.5": invalid syntax


get term aggs for Malts. Then convert the Data like: "Pilsener" and "Pilsener Malz" zu "Pilsener"

- convert Pilsner Malz, Pilsner, Pilsener, Pilsenermalz, (Pilsenermalz, hell) to Pilsener Malz with data pipelines in es
- Pale Ale, Pale ale Malz, Best Pale Ale Malz-> Pale Ale Malz
- Münchener, Münchner Malz -> Münchener Malz
- Wiener, Wienermalz -> Wiener Malz
- Weizenmalz hell, (Weizenmalz, hell), Weizenmalz Hell -> Weizenmalz
- Carahell, CaraHell, Cara Hell
- CaraMünch II, Cara II
- CaraAmber, Cara Amber
- Amber Malz, Amber Malt
- Best Röstmalz -> Röstmalz
- Best Red X -> Red X




find duplicates with elastic search and remove them: https://www.elastic.co/de/blog/how-to-find-and-remove-duplicate-documents-in-elasticsearch






*/

func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

func DecodeHTMLBody(body io.Reader, charset string) (io.Reader, error) {
	if charset == "" {
		charset = detectContentCharset(body)
	}
	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}
	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}
	return body, nil
}

func rmChar(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}

	return strings.Map(filter, input)

}

func inSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getUserComments(recipeID int) *[]recipe.RecipeComment {
	var comments []recipe.RecipeComment

	resp, err := http.Get(fmt.Sprintf("https://www.maischemalzundmehr.de/index.php?id=%d&inhaltmitte=rezept", recipeID))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := DecodeHTMLBody(resp.Body, "")
	if err != nil {
		fmt.Println("error")
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	doc.Find("div.userkommentare").Each(func(i int, s *goquery.Selection) {
		child := s.Find("p").Contents()
		name := strip.StripTags(child.Eq(1).Text())
		datetime := strings.Split(rmChar(strip.StripTags(child.Eq(2).Text()), " \n"), "-")
		comment := rmChar(strip.StripTags(s.Find("p").Next().Text()), "\n")
		comments = append(comments, recipe.RecipeComment{
			Name:    name[5:],
			Date:    datetime[0],
			Comment: comment,
		})
	})

	return &comments
}

func main() {

	files, err := ioutil.ReadDir("recipes/")
	if err != nil {
		log.Fatal(err)
	}
	var keymap []int
	var id int
	for _, f := range files {
		id, err = strconv.Atoi(strings.Split(f.Name(), "_")[0])
		if err == nil {
			keymap = append(keymap, id)
		}
	}

	exportUrl := "https://www.maischemalzundmehr.de/export.php?id="

	resp, err := http.Get("https://www.maischemalzundmehr.de/index.php")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	re := regexp.MustCompile(`<a class="rezeptlink" href="index.php\?id=(?P<id>\d+)\&inhaltmitte=rezept">`)
	match := re.FindStringSubmatch(string(body))

	lid, err := strconv.Atoi(match[1])
	cnt := 0
	errcnt := 0

	for lid > 0 {
		if inSlice(lid, keymap) {
			lid--
			continue
		}

		//get json export
		resp, err = http.Get(fmt.Sprintf("%s%d", exportUrl, lid))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if len(body) != 471 {
			cnt++
			pp := &recipe.RecipeM3{}

			recipe, err := pp.Load(string(body))
			if err != nil {
				errcnt++
				ioutil.WriteFile(fmt.Sprintf("broken/%d.json", lid), body, 0644)
				fmt.Println(err)
				lid--
				continue
			}

			fn := strings.ReplaceAll(recipe.Global.Name, " ", "_")
			fn = rmChar(fn, "?%$&#-`'().")

			//get comments
			recipe.Comment = *getUserComments(lid)

			recipe.SavePretty(fmt.Sprintf("recipes/%d_%s.json", lid, fn))
			fmt.Println(recipe.Global.Name)
		}

		lid--
	}
	fmt.Printf("Fetch %d documents\nError on Fetch %d\n", cnt, errcnt)

}
