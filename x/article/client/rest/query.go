package rest

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	articleTypes "github.com/quokki/quokki/x/article"

	"github.com/gorilla/mux"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc(
		"/article/{id:[0-9]+}",
		QueryArticleRequestHandlerFn(storeName, cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/article/{parentId:[0-9]+}/{id:[0-9]+}",
		QueryArticleRequestHandlerFn(storeName, cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/article/root",
		QueryArticleRootRequestHandlerFn(storeName, cdc, cliCtx),
	)
	r.HandleFunc(
		"/articles/{page:[0-9]+}",
		QueryArticlesRequestHandlerFn(storeName, cdc, cliCtx),
	)
	r.HandleFunc(
		"/articles/{parentId:[0-9]+}/{page:[0-9]+}",
		QueryArticlesRequestHandlerFn(storeName, cdc, cliCtx),
	)
}

func QueryArticleRequestHandlerFn(storeName string, cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		parentId := -1
		if vars["parentId"] != "" {
			parentId, err = strconv.Atoi(vars["parentId"])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		}

		bz := make([]byte, 8)
		binary.BigEndian.PutUint64(bz, uint64(id))
		if parentId >= 0 {
			pbz := make([]byte, 8)
			binary.BigEndian.PutUint64(pbz, uint64(parentId))
			bz = append(pbz, bz...)
		}

		res, err := cliCtx.QueryStore(append([]byte("article"), bz...), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query article. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// decode the value
		article := articleTypes.Article{}
		err = cdc.UnmarshalBinaryBare(res, &article)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't parse article result. Result: %s. Error: %s", res, err.Error())))
			return
		}

		// print out whole account
		output, err := cdc.MarshalJSON(article)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't marshall article result. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

func QueryArticleRootRequestHandlerFn(storeName string, cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryStore([]byte("article"), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query article. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// decode the value
		article := articleTypes.Article{}
		err = cdc.UnmarshalBinaryBare(res, &article)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't parse article result. Result: %s. Error: %s", res, err.Error())))
			return
		}

		// print out whole account
		output, err := cdc.MarshalJSON(article)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't marshall article result. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

type Articles []articleTypes.Article

func (articles Articles) Len() int {
	return len(articles)
}

func (articles Articles) Less(i, j int) bool {
	return bytes.Compare(articles[i].Id, articles[j].Id) < 0
}

func (articles Articles) Swap(i, j int) {
	articles[i], articles[j] = articles[j], articles[i]
}

func QueryArticlesRequestHandlerFn(storeName string, cdc *wire.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		page, err := strconv.Atoi(vars["page"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if page <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Negative page"))
			return
		}

		parentId := -1
		if vars["parentId"] != "" {
			parentId, err = strconv.Atoi(vars["parentId"])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		}

		v := r.URL.Query()
		perPage := 10
		if v.Get("per_page") != "" {
			var err error
			perPage, err = strconv.Atoi(v.Get("per_page"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			if perPage > 30 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Too many per page"))
				return
			}
		}

		pId := []byte("article")
		if parentId >= 0 {
			bz := make([]byte, 8)
			binary.BigEndian.PutUint64(bz, uint64(parentId))
			pId = append(pId, bz...)
		}

		res, err := cliCtx.QueryStore(pId, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query article. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// decode the value
		article := articleTypes.Article{}
		err = cdc.UnmarshalBinaryBare(res, &article)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't parse article result. Result: %s. Error: %s", res, err.Error())))
			return
		}

		numOfArticles := int(article.Sequence)
		start := numOfArticles - (page * perPage)
		end := numOfArticles - ((page - 1) * perPage)

		articles := make(Articles, 0, perPage)
		count := 0

		for id := end - 1; id >= start && id >= 0; id-- {
			count++
		}
		chs := make(chan articleTypes.Article, count)
		defer func() {
			close(chs)
		}()

		// TODO: Use iterator. They don't support query by iterator yet...
		for id := end - 1; id >= start && id >= 0; id-- {
			go func(id int) {
				bz := make([]byte, 8)
				binary.BigEndian.PutUint64(bz, uint64(id))
				res, err := cliCtx.QueryStore(append(pId, bz...), storeName)
				if err != nil {
					chs <- articleTypes.Article{}
				}

				article := articleTypes.Article{}
				err = cdc.UnmarshalBinaryBare(res, &article)
				if err != nil {
					chs <- articleTypes.Article{}
				}

				chs <- article
			}(id)
		}

		for i := 0; i < count; i++ {
			if article, ok := <-chs; ok {
				if len(article.Id) > 0 {
					articles = append(articles, article)
				}
			}
		}

		sort.Sort(sort.Reverse(articles))

		// print out whole account
		output, err := cdc.MarshalJSON(articles)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't marshall article result. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}
