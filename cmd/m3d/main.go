package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/ripx80/brewman/pkgs/recipe"
)

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

func getLastID(baseURL string) (int, error) {
	resp, err := http.Get(baseURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err
	}
	link, ex := doc.Find(".rezeptlink").Eq(0).Attr("href")
	if !ex {
		return 0, fmt.Errorf("canot find the last id")
	}
	base, err := url.Parse(link)
	if err != nil {
		return 0, err
	}

	keys, ok := base.Query()["id"]
	if !ok {
		return 0, fmt.Errorf("no id extraction from index url")
	}
	lid, err := strconv.Atoi(keys[0])
	return lid, nil
}

func getUserComments(recipeID int) (*[]recipe.RecipeComment, error) {
	var comments []recipe.RecipeComment

	resp, err := http.Get(fmt.Sprintf("https://www.maischemalzundmehr.de/index.php?id=%d&inhaltmitte=rezept", recipeID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("div.userkommentare").Each(func(i int, s *goquery.Selection) {
		child := s.Find("p").Contents()
		name := strip.StripTags(child.Eq(1).Text())
		datetime := strings.Split(rmChar(strip.StripTags(child.Eq(2).Text()), " \n"), "-")
		comment := strings.TrimSpace(strings.ReplaceAll(strip.StripTags(s.Find("p").Next().Text()), "\n", " "))
		comments = append(comments, recipe.RecipeComment{
			Name:    name[5:],
			Date:    datetime[0],
			Comment: comment,
		})
	})

	return &comments, nil
}

const m3url = "https://www.maischemalzundmehr.de"
const endmsg = "\n\nFetch: %d doc\nBroken docs: %d\n"

func main() {
	outdir := flag.String("output", "recipes", "output dir. if not exists it will be created")
	flag.Parse()

	if _, err := os.Stat(*outdir); os.IsNotExist(err) {
		os.Mkdir(*outdir, 0755)
	}

	// simple check if some recipes downloaded before
	files, err := ioutil.ReadDir(*outdir)
	if err != nil {
		log.Fatal(err)
	}

	var keymap []int
	for _, f := range files {
		id, err := strconv.Atoi(strings.Split(f.Name(), "_")[0])
		if err == nil {
			keymap = append(keymap, id)
		}
	}

	lid, err := getLastID(fmt.Sprintf("%s/index.php", m3url))
	if err != nil {
		log.Fatal(err)
	}

	cnt := 0
	errcnt := 0
	var body []byte

	//signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Printf(endmsg, cnt, errcnt)
		os.Exit(0)
		return
	}()

	for ; lid > 0; lid-- {
		if inSlice(lid, keymap) {
			continue
		}

		//get json export
		resp, err := http.Get(fmt.Sprintf("%s/export.php?id=%d", m3url, lid))
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)

		//empty json
		if len(body) <= 471 {
			continue
		}

		cnt++
		recipe, err := (&recipe.RecipeM3{}).Load(string(body))
		if err != nil {
			errcnt++
			broken := path.Join(*outdir, "broken")
			if _, err := os.Stat(broken); os.IsNotExist(err) {
				os.Mkdir(broken, 0755)
			}
			ioutil.WriteFile(path.Join(broken, fmt.Sprintf("%d.json", lid)), body, 0644)
			log.Println(err)
			continue
		}

		fn := strings.ReplaceAll(recipe.Global.Name, " ", "_")
		fn = rmChar(fn, "?%$&#-`'().")

		//get comments
		c, err := getUserComments(lid)
		if err != nil {
			log.Fatal("usercomments: ", err)
		}
		recipe.Comment = *c

		recipe.SavePretty(path.Join(*outdir, fmt.Sprintf("%d_%s.json", lid, fn)))
		log.Println(recipe.Global.Name)

	}
	log.Printf(endmsg, cnt, errcnt)

}
