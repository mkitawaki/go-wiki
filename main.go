package main

import (
	"github.com/gin-gonic/gin"
	w "github.com/thinxer/go-word2vec"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Distance struct {
	SpecifiedWord string   `json:"specified_word"`
	Words         []w.Pair `json:"words"`
}

type Analogy struct {
	Is    string   `json:"is"`
	To    string   `json:"to"`
	What  string   `json:"what"`
	Words []w.Pair `json:"words"`
}

type MostSimilarity struct {
	PositiveWords []string `json:"positive_words"`
	NegativeWords []string `json:"negative_words"`
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
	router.GET("/analogy/:is/:to/:what/:count", analogy)
	router.GET("/mostSimilarity/:positives/:negatives/:count", mostSimilarity)
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

func analogy(c *gin.Context) {
	is := c.Param("is")
	to := c.Param("to")
	what := c.Param("what")
	count := c.Param("count")
	positives := []string{what, to}
	negatives := []string{is}
	cnt, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(http.StatusOK, &Analogy{
			Is:    is,
			To:    to,
			What:  what,
			Words: []w.Pair{},
		})
		return
	}
	pairs, err := Word2Vec.MostSimilar(positives, negatives, cnt)
	if err != nil {
		c.JSON(http.StatusOK, &Analogy{
			Is:    is,
			To:    to,
			What:  what,
			Words: []w.Pair{},
		})
		return
	}
	c.JSON(http.StatusOK, &Analogy{
		Is:    is,
		To:    to,
		What:  what,
		Words: pairs,
	})
}

func mostSimilarity(c *gin.Context) {
	p := c.Param("positives")
	n := c.Param("negatives")
	count := c.Param("count")
	positives := strings.Split(p, "+")
	negatives := strings.Split(n, "+")
	cnt, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(http.StatusOK, &MostSimilarity{
			PositiveWords: positives,
			NegativeWords: negatives,
			Words:         []w.Pair{},
		})
		return
	}
	pairs, err := Word2Vec.MostSimilar(positives, negatives, cnt)
	if err != nil {
		c.JSON(http.StatusOK, &MostSimilarity{
			PositiveWords: positives,
			NegativeWords: negatives,
			Words:         []w.Pair{},
		})
		return
	}
	c.JSON(http.StatusOK, &MostSimilarity{
		PositiveWords: positives,
		NegativeWords: negatives,
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
