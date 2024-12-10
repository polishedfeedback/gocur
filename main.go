package main

import (
  "fmt"
  "flag"
  "log"
  "encoding/json"
  "os"

  "github.com/everapihq/freecurrencyapi-go"
  "github.com/joho/godotenv"
)


func main(){

  err := godotenv.Load()
  if err != nil {
    log.Println("Error loading .env file : %v", err)
  }

  apiKey := os.Getenv("GOCUR_API_KEY")
  if apiKey == "" {
    log.Fatal("API_KEY not set in the environment")
  }

  freecurrencyapi.Init(apiKey)

  from := flag.String("from", "GBP", "Currency to convert from (e.g., GBP)")
  to := flag.String("to", "USD", "Currency to convert to (e.g., USD)")
  cur := flag.Float64("cur", 1.0, "Amount to convert (e.g., 1000)")
  flag.Parse()

  rates := freecurrencyapi.Latest(map[string] string{
    "base_currency": *from,
    "currencies": *to,
  })

  // fmt.Println(string(rates))

  var result = make(map[string]interface{})

  if err := json.Unmarshal(rates, &result); err != nil {
    log.Fatalf("error decoding JSON: %v", err)
  }

  if data, ok := result["data"].(map[string]interface{}); ok {
      usdToFrom, fromExists := data[*from].(float64) 
      usdToTo, toExists := data[*to].(float64)     

      if fromExists && toExists {
          rate := usdToTo / usdToFrom
          fmt.Printf("Exchange rate from %v to %v: %.4f\n", *from, *to, (*cur)*rate)
      } else {
          log.Fatalf("Exchange rate for %v or %v not found", *from, *to)
      }
  } else {
      log.Fatalf("Invalid response format")
  }
}
