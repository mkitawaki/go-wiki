package main

import (
	"github.com/gin-gonic/gin"
	w "github.com/thinxer/go-word2vec"
	"net/http"
	"os"
	"strconv"
)

type Distance struct {
	SpecifiedWord string   `json:"specified_word"`
	Words         []w.Pair `json:"words"`
}

type Similarity struct {
	SpecifiedWords []*string `json:"specified_words"`
	Similary       float32   `json:"similary"`
}

var Word2Vec *w.Model

func initWord2Vec(dic string) {
	var err error
	Word2Vec, err = w.Load(dic)
	if err != nil {
		panic(err)
	}
}

func main() {

	initWord2Vec("wikipedia.bin")

	router := gin.Default()

	router.GET("/distance/:word/:count", distance)
	router.GET("/similarity/:x/:y", similarity)

	router.Run(":" + os.Getenv("PORT"))
}

func distance(c *gin.Context) {
	word := c.Param("word")
	count := c.Param("count")
	positives := []string{word}
	negatives := []string{}
	cnt, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(http.StatusOK, &Distance{
			SpecifiedWord: word,
			Words:         []w.Pair{},
		})
		return
	}
	pairs, err := Word2Vec.MostSimilar(positives, negatives, cnt)
	if err != nil {
		c.JSON(http.StatusOK, &Distance{
			SpecifiedWord: word,
			Words:         []w.Pair{},
		})
		return
	}
	c.JSON(http.StatusOK, &Distance{
		SpecifiedWord: word,
		Words:         pairs,
	})
}

func similarity(c *gin.Context) {
	x := c.Param("x")
	y := c.Param("y")
	similary, err := Word2Vec.Similarity(x, y)
	if err != nil {
		c.JSON(http.StatusOK, &Similarity{
			SpecifiedWords: []*string{&x, &y},
			Similary:       0,
		})
		return
	}
	c.JSON(http.StatusOK, &Similarity{
		SpecifiedWords: []*string{&x, &y},
		Similary:       similary,
	})
}
