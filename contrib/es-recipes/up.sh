#!/bin/bash

set -e
set -o pipefail

if [ ! -d docker-elk ]; then
  git clone https://github.com/deviantony/docker-elk
fi
(
    cd docker-elk
    docker-compose up -d
)

if [ ! -d brewman ]; then
    git clone https://github.com/ripx80/brewman
    (
        cd brewman
        go build -o brewman main.go && \
        chmod +x brewman && \
        brewman recipe down
    )
fi

until [ $(curl --user elastic:changeme -Is http://localhost:9200/ | head -n 1 | wc -l) -eq 1 ]
do
   sleep 5
   echo "connecting..."
done

# create index
# curl --user elastic:changeme -X PUT "localhost:9200/recipe?pretty" -H 'Content-Type: application/json' -d'
# {}
# '

# create mapping
curl --user elastic:changeme -X PUT "localhost:9200/recipes/" -H 'Content-Type: application/json' -d @recipe-m3-mapping.json

# push recipes
(
    cd brewman/recipes
    for i in *;do
        if [ -d $i ]; then continue; fi
        curl --user elastic:changeme -XPOST http://localhost:9200/recipes/_doc -H "Content-Type: application/json" -d @$i
    done
)
echo """
exposed:
    5000: Logstash TCP input
    9200: Elasticsearch HTTP
    9300: Elasticsearch TCP transport
    5601: Kibana

change password:
    docker-compose exec -T elasticsearch bin/elasticsearch-setup-passwords auto --batch

cleanup:
    (cd docker-elk;docker-compose down -v)
"""

exit 0
