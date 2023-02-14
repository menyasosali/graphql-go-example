package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
	"time"
)

type Film struct {
	ID       int
	Title    string
	Author   Author
	Year     int
	Comments []Comment
}

type Author struct {
	Name  string
	Films []int
}

type Comment struct {
	AuthorComment AuthorComment
	Body          string
	Time          time.Time
}

type AuthorComment struct {
	Name  string
	Email string
}

func populate() []Film {
	author1 := &Author{Name: "Elliot Forbes", Films: []int{1}}
	authorcomment1 := &AuthorComment{Name: "Fake user 1", Email: "fakeuser1@mail.ru"}
	film1 := Film{
		ID:     1,
		Title:  "Go GraphQL Tutorial",
		Author: *author1,
		Comments: []Comment{
			Comment{AuthorComment: *authorcomment1, Body: "First Comment", Time: time.Now()},
		},
	}
	author2 := &Author{Name: "Mark Spancer", Films: []int{2}}
	authorcomment2 := &AuthorComment{Name: "Fake user 2", Email: "fakeuser2@mail.ru"}
	film2 := Film{
		ID:     2,
		Title:  "Go World",
		Author: *author2,
		Comments: []Comment{
			Comment{AuthorComment: *authorcomment2, Body: "First Comment", Time: time.Now()},
		},
	}

	var films []Film
	films = append(films, film1, film2)

	return films
}

func main() {
	films := populate()
	var authorCommentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "AuthorComment",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"email": &graphql.Field{
					Type: graphql.String,
				},
			},
		})

	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"authorComment": &graphql.Field{
					Type: authorCommentType,
				},
				"body": &graphql.Field{
					Type: graphql.String,
				},
				"time": &graphql.Field{
					Type: graphql.DateTime,
				},
			},
		})

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"author": &graphql.Field{
					Type: graphql.String,
				},
				"films": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		})

	var filmType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Film",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: authorType,
				},
				"year": &graphql.Field{
					Type: graphql.Int,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		})

	fileds := graphql.Fields{
		"film": &graphql.Field{
			Type:        filmType,
			Description: "Get Film by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"]
				if ok {
					for _, film := range films {
						if int(film.ID) == id {
							return film, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(filmType),
			Description: "Get film list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return films, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fileds}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create schema, error: %v", err)
	}

	query := `
{
		list {
			id
			title
			year
			}
}
`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJson, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJson)
}
