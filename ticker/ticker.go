package ticker

import (
	"fmt"
	"html/template"
	"htmx/utils"
	"math/rand"
	"net/http"
	"time"
)

type Price struct {
	Dollars int
	Cents   int
}

func (p Price) String() string {
	return fmt.Sprintf("$%d.%d", p.Dollars, p.Cents)
}

type Ticker struct {
	Symbol   string
	Price    Price
	Quantity int
}

type TickerList struct {
	broadcasters map[string]utils.Broadcaster[Ticker]
	quitters     map[string]chan bool
}

func (tickers TickerList) AddTicker(symbol string) {
	tickers.broadcasters[symbol] = utils.NewBroadcaster[Ticker](3)
}

func (tickers TickerList) RemoveTicker(symbol string) {
	tickers.StopTicker(symbol)
	delete(tickers.broadcasters, symbol)
}

func (tickers TickerList) StartTicker(symbol string) {
	var quit = make(chan bool)
	tickers.quitters[symbol] = quit
	go TickerLoopUpdater(tickers.broadcasters[symbol], symbol, quit)
}

func (tickers TickerList) StopTicker(symbol string) {
	// Quit if quit channel exists
	if quit, ok := tickers.quitters[symbol]; ok {
		quit <- true
		close(quit)
		delete(tickers.quitters, symbol)
	}
}

func (tickers TickerList) StartAllTickers() {
	for symbol := range tickers.broadcasters {
		tickers.StartTicker(symbol)
	}
}

func (tickers TickerList) StopAllTickers() {
	for symbol := range tickers.broadcasters {
		tickers.StopTicker(symbol)
	}
}

func (tickers TickerList) GetTicker(symbol string) utils.Broadcaster[Ticker] {
	return tickers.broadcasters[symbol]
}

func CreateTickerListFromArray(symbols []string) TickerList {
	var tickers = TickerList{
		broadcasters: make(map[string]utils.Broadcaster[Ticker]),
		quitters:     make(map[string]chan bool),
	}
	for _, symbol := range symbols {
		tickers.broadcasters[symbol] = utils.NewBroadcaster[Ticker](3)
	}
	return tickers
}

func TickerLoopUpdater(broadcaster utils.Broadcaster[Ticker], symbol string, quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			dollars := rand.Intn(100)
			cents := rand.Intn(100)
			quantity := rand.Intn(500) + 100
			interval := rand.Intn(1000) + 200
			ticker := Ticker{
				Symbol:   symbol,
				Price:    Price{Dollars: dollars, Cents: cents},
				Quantity: quantity,
			}
			broadcaster.Submit(ticker)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
}

func CreateTickerListener(tickers TickerList, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")
		symbols := r.URL.Query()["symbol"]

		fmt.Println("Client connected to ticker(s):", symbols)

		ch := make(chan Ticker)
		for _, symbol := range symbols {
			tickers.GetTicker(symbol).Register(ch)
			defer tickers.GetTicker(symbol).Unregister(ch)
		}

		for {
			select {
			case <-r.Context().Done():
				fmt.Println("Client disconnected from ticker(s):", symbols)
				return
			case ticker := <-ch:
				element := utils.ExecuteTemplateToString(templates, "ticker", ticker)
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ticker.Symbol, element)
				w.(http.Flusher).Flush()
			}
		}
	}
}
