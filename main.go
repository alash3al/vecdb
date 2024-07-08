package main

import (
	"fmt"
	"github.com/alash3al/vecdb/internals/store"
	_ "github.com/alash3al/vecdb/internals/store/drivers/v1"
	"github.com/alash3al/vecdb/internals/vector"
)

func main() {
	db, err := store.Open("v1", "./v1.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vectors := []vector.Vec{
		{32.4, 74.1, 3.2},
		{15.1, 19.2, 15.8},
		{0.16, 1.2, 3.8},
		{75.1, 67.1, 29.9},
		{58.8, 6.7, 3.4},
	}

	for i, v := range vectors {
		if err := db.Put("playground", fmt.Sprintf("vec_%02d", i+1), v); err != nil {
			panic(err.Error())
		}
	}

	result, err := db.Query(store.VectorQueryInput{
		Bucket:                      "playground",
		Vector:                      []float64{54.8, 5.5, 3.1},
		MinCosineSimilarityDistance: 0.6,
		MaxResultCount:              10,
	})
	if err != nil {
		panic(err.Error())
	}

	for _, v := range result.Items {
		fmt.Println(v.Key, " <===> ", v.Distance)
	}
}
