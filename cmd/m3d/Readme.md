# m3d - maischemalzundmehr crawler

Get all recipes from m3 and convert it to brewman recipe json format. If you have downloaded recipes before they will be skipped if they are in the same dir.

## Roadmap

Get more recipes and convert to brewman format

```plain
https://beersmithrecipes.com/ (very big community! own format bsmx (just xml), must login for download recipes, no api )
    https://beersmithrecipes.com/viewrecipe/2589508 id counting, and export
    https://beersmithrecipes.com/download.php?id=2589508

https://beerandbrewing.com/beer-recipes/ (Beersmith and BeerXML version. Must subscribe )
https://beerrecipes.org/ (4495 recipes)
    https://beerrecipes.org/Recipe/1167 count id ;-), no export -.-

https://www.brewersfriend.com/
    https://www.brewersfriend.com/homebrew/recipe/view/1633 id counting,
    https://www.brewersfriend.com/homebrew/recipe/beerxml1.0/1633 export


https://www.brewerydb.com/developers
https://www.brewerydb.com/developers/docs/ (API)
https://github.com/homebrewing/tapline
http://www.malt.io/ (xml export, all recipes, get links and add beerxml to download xml)
    http://www.malt.io/users/david/recipes/the-duvel-s-in-the-details-belgian-golden-strong/beerxml but xml is very bad. no water aso

```

## Broken json

- Parsing Amount error: Hopfen_VWH_1_Menge strconv.ParseFloat: parsing " 90": invalid syntax
- json: invalid use of ,string struct tag, trying to unmarshal "" into float64
- Rest value is empty: Infusion_Rastzeit4
- Parsing Amount error: Infusion_Rastzeit2 strconv.Atoi: parsing "7.5": invalid syntax

## Notices

### Elasticsearch Stuff

insert into recipes-m3 index. change your user:password arg

```bash
curl --user elastic:changeme -XPOST http://node1:9200/recipes-m3/doc -H "Content-Type: application/json" -d @400_Meraner_Weizen.json
for i in $(ls):; do curl --user elastic:changeme -XPOST http://node1:9200/recipes-m3/doc -H "Content-Type: application/json" -d @$i; done

#or with bulk
PUT hockey/_bulk?refresh
{"index":{"_id":1}}
{"first":"johnny","last":"gaudreau","goals":[9,27,1],"assists":[17,46,0],"gp":[26,82,1],"born":"1993/08/13"}
{"index":{"_id":2}}
```

Found a lot of interesting misspelled stuff in recipes. So this will be corrected.
Then convert the Data like: "Pilsener" and "Pilsener Malz" zu "Pilsener" in mapping from m3 to recipe struct

- convert Pilsner Malz, Pilsner, Pilsener, Pilsenermalz, "Pilsenermalz, hell" to Pilsener Malz with data pipelines in es
- Pale Ale, Pale ale Malz, Best Pale Ale Malz-> Pale Ale Malz
- Münchener, Münchner Malz -> Münchener Malz
- Wiener, Wienermalz -> Wiener Malz
- Weizenmalz hell, "Weizenmalz, hell", Weizenmalz Hell -> Weizenmalz
- Carahell, CaraHell, Cara Hell
- CaraMünch II, Cara II
- CaraAmber, Cara Amber
- Amber Malz, Amber Malt
- Best Röstmalz -> Röstmalz
- Best Red X -> Red X

find duplicates with elastic search and remove them: https://www.elastic.co/de/blog/how-to-find-and-remove-duplicate-documents-in-elasticsearch
