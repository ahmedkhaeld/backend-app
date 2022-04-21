package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"io"
	"log"
	"net/http"
	"strings"
)

// movies hold the data we are querying
var movies []*models.Movie

// fields graphql schema definition that describe the fields used in GQL
// this tells gql what kind of queries we are getting, which in here there are two
// individual movie, and list of movies
var fields = graphql.Fields{
	"movie": &graphql.Field{
		Type:        movieType,
		Description: "Get movie by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(int)
			if ok {
				for _, movie := range movies {
					if movie.ID == id {
						return movie, nil
					}
				}
			}
			return nil, nil
		},
	},
	"list": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Get all movies",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return movies, nil
		},
	},
	"search": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Search movies by title",
		Args: graphql.FieldConfigArgument{
			"titleContains": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var theList []*models.Movie
			search, ok := params.Args["titleContains"].(string)
			if ok {
				for _, currentMovie := range movies {
					if strings.Contains(currentMovie.Title, search) {
						log.Println("Found one")
						theList = append(theList, currentMovie)
					}
				}
			}
			return theList, nil
		},
	},
}

// movieType describe the kind of info(fields) in movies variable to expose to the end users
var movieType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Movie",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"year": &graphql.Field{
				Type: graphql.Int,
			},
			"release_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"runtime": &graphql.Field{
				Type: graphql.Int,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"mpaa_rating": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"poster": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func (app *application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	// populate the query data to var movies
	movies, _ = app.models.DB.All()

	// read the request body, which comes in json format
	q, _ := io.ReadAll(r.Body)
	query := string(q) // type conversion to string

	log.Println(query)

	// set up the schema
	// 1. specify the  root query
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	// 2. create a schema configuration
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	// 3. make the schema, and hand it the configuration
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		app.errorJSON(w, errors.New("failed to create schema"))
		log.Println(err)
		return
	}
	// now we have a schema,
	//at this point extract the params in the request
	// hand the gql request, so we can get a response
	params := graphql.Params{Schema: schema, RequestString: query}
	resp := graphql.Do(params)
	if len(resp.Errors) > 0 {
		app.errorJSON(w, errors.New(fmt.Sprintf("failed: %+v", resp.Errors)))
	}
	// now we have a response, we can send back in form of json

	j, _ := json.MarshalIndent(resp, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
