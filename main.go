package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	w "github.com/thinxer/go-word2vec"
	"net/http"
	//"os"
	"strconv"
	"strings"
)

// Distance is struct to returned when distance api called
type Distance struct {
	SpecifiedWord string   `json:"specified_word"`
	Words         []w.Pair `json:"words"`
}

func newDistance(s string, w []w.Pair) *Distance {
	return &Distance{
		SpecifiedWord: s,
		Words:         w,
	}
}

func notFoundDistance(s string) *Distance {
	return newDistance(s, []w.Pair{})
}

// Analogy is struct to returned when analogy api called
type Analogy struct {
	Is    string   `json:"is"`
	To    string   `json:"to"`
	What  string   `json:"what"`
	Words []w.Pair `json:"words"`
}

func newAnalogy(i string, t string, h string, p []w.Pair) *Analogy {
	return &Analogy{
		Is:    i,
		To:    t,
		What:  h,
		Words: p,
	}
}

func notFoundAnalogy(i string, t string, h string) *Analogy {
	return newAnalogy(i, t, h, []w.Pair{})
}

// MostSimilarity is struct to returned when mostSimilarity api called
type MostSimilarity struct {
	PositiveWords []string `json:"positive_words"`
	NegativeWords []string `json:"negative_words"`
	Words         []w.Pair `json:"words"`
}

func newMostSimilarity(p []string, n []string, w []w.Pair) *MostSimilarity {
	return &MostSimilarity{
		PositiveWords: p,
		NegativeWords: n,
		Words:         w,
	}
}

func notFoundMostSimilarity(p []string, n []string) *MostSimilarity {
	return newMostSimilarity(p, n, []w.Pair{})
}

// Similarity is struct to returned when similarity api called
type Similarity struct {
	SpecifiedWords []string `json:"specified_words"`
	Similary       float32  `json:"similary"`
}

func newSimilarity(s []string, m float32) *Similarity {
	return &Similarity{
		SpecifiedWords: s,
		Similary:       m,
	}
}

func notFoundSimilarity(s []string) *Similarity {
	return newSimilarity(s, 0)
}

var word2Vec *w.Model

func initWord2Vec(dic string) {
	var err error
	word2Vec, err = w.Load(dic)
	if err != nil {
		panic(err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(
			"Access-Control-Allow-Origin",
			"http://localhost:8100")
		c.Writer.Header().Set(
			"Access-Control-Max-Age",
			"86400")
		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set(
			"Access-Control-Expose-Headers",
			"Content-Length")
		c.Writer.Header().Set(
			"Access-Control-Allow-Credentials",
			"true")

		if c.Request.Method == "OPTIONS" {
			//	fmt.Println("OPTIONS")
			c.AbortWithStatus(204)
		}
		c.Next()
	}
}

func main() {

	initWord2Vec("wikipedia.bin")

	router := gin.New()

	router.Use(corsMiddleware())
	router.GET("/distance/:word/:count", distance)
	router.GET("/analogy/:is/:to/:what/:count", analogy)
	router.GET("/mostSimilarity/:positives/:negatives/:count", mostSimilarity)
	router.GET("/similarity/:x/:y", similarity)

	router.Run(":8080") // + os.Getenv("PORT"))
}

func distance(c *gin.Context) {
	word := c.Param("word")
	count := c.Param("count")
	positives := []string{word}
	negatives := []string{}
	cnt, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(http.StatusOK, notFoundDistance(word))
		return
	}
	pairs, err := word2Vec.MostSimilar(positives, negatives, cnt)
	if err != nil {
		c.JSON(http.StatusOK, notFoundDistance(word))
		return
	}
	c.JSON(http.StatusOK, newDistance(word, pairs))
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
		c.JSON(http.StatusOK, notFoundAnalogy(is, to, what))
		return
	}
	pairs, err := word2Vec.MostSimilar(positives, negatives, cnt)
	if err != nil {
		c.JSON(http.StatusOK, notFoundAnalogy(is, to, what))
		return
	}
	c.JSON(http.StatusOK, newAnalogy(is, to, what, pairs))
}

func mostSimilarity(c *gin.Context) {
	p := c.Param("positives")
	n := c.Param("negatives")
	count := c.Param("count")
	positives := strings.Split(p, "+")
	negatives := strings.Split(n, "+")
	cnt, err := strconv.Atoi(count)
	if err != nil {
		c.JSON(http.StatusOK, notFoundMostSimilarity(positives, negatives))
		return
	}
	pairs, err := word2Vec.MostSimilar(positives, negatives, cnt)
	if err != nil {
		c.JSON(http.StatusOK, notFoundMostSimilarity(positives, negatives))
		return
	}
	c.JSON(http.StatusOK, newMostSimilarity(positives, negatives, pairs))
}

func similarity(c *gin.Context) {
	x := c.Param("x")
	y := c.Param("y")
	similary, err := word2Vec.Similarity(x, y)
	if err != nil {
		c.JSON(http.StatusOK, notFoundSimilarity([]string{x, y}))
		return
	}
	c.JSON(http.StatusOK, newSimilarity([]string{x, y}, similary))
}
