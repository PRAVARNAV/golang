package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	_ "github.com/lib/pq"
)

const (
	// Initialize connection constants.
	HOST     = "34.131.237.18"
	DATABASE = "siva_pg_db_azure"
	USER     = "dev_team"
	PASSWORD = "-X1-LhpP(5Th%gaq"

	/*client details*/
	grant_type  = "client_credentials"
	client_id = "36161062-7913-4b80-884c-28e45ac71540"
	client_secret = "Lrx8Q~u0cK5iig8GspTVlS1rYctdQQE6qqrOPc8A"
	scope  ="https%3A%2F%2Fmanagement.azure.com%2F.default"
	url_endpoint = "https://login.microsoftonline.com/eab3c9cb-cfc6-4a90-b01d-41c77afef7f3/oauth2/v2.0/token"
	subscription_id = "8524df34-d0be-4424-9c6d-50bd515836fd"
)



func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
/*creating connection for database*/
	
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", HOST, USER, PASSWORD, DATABASE)
	db, err := sql.Open("postgres", connectionString)
	checkError(err)
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to Postgres database \n")

/*creating table for azure billing*/

	_, err = db.Exec("CREATE TABLE  IF  NOT EXISTS azure_billing_new(id VARCHAR(300),Resource_id  character varying(400), Resource_name varchar(300),Resource_group varchar(300),tendenid varchar(300) ,consumedService VARCHAR(300), subscriptionGuid VARCHAR(300),currency  VARCHAR(300),usageQuantity   NUMERIC ,meterId VARCHAR(300), usageStart VARCHAR(300) , usageEnd VARCHAR(300),UnitPrice float ,Location varchar(40),SkuId varchar(40),ServiceName varchar(40),UnitOfMeasure VARCHAR(40),ArmSkuName VARCHAR(40));")
	checkError(err)
	fmt.Println("Finished creating table azure_billing_res \n")
	
/*Here we are passing constant values(dynamically)*/
// var client_id string 
// fmt.Println("enter your client id: ")
// fmt.Scanln(&client_id)

// var client_secret string 
// fmt.Println("enter your client Secret: ")
// fmt.Scanln(&client_secret)

// var subscription_id string 
// fmt.Println("enter your subscription_id: ")
// fmt.Scanln(&subscription_id)



	url := url_endpoint
	method := "POST"
	client_details:= fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&scope=%s",grant_type,client_id,client_secret,scope)
	fmt.Println("client details for bearer token = ",client_details)
	payload := strings.NewReader(client_details)
	client1 := &http.Client{}
	req1, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req1.Header.Add("Cookie", "fpc=An0PeRGH9pBLi4eUtOSP2LvgfLQZAQAAALEvGtsOAAAANrU7NwEAAABPMBrbDgAAAA; stsservicecookie=estsfd; x-ms-gateway-slice=estsfd")

	res1, err := client1.Do(req1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res1.Body.Close()

	consumptn_body, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

/*Here removing extra characters from bearer response*/

	str1 := string(consumptn_body)
	re, err := regexp.Compile(`[:}"]`)
	if err != nil {
		log.Fatal(err)
	}
	str1 = re.ReplaceAllString(str1, " ")
	rmv_extra := str1
	// fmt.Println(str1[113:])

/*Here we have to pass Bearer keyword into bearer token*/

	myslice := []string{"Bearer", string(rmv_extra[79:])}
	result := strings.Join(myslice, " ")
	
/* Here We are  passing  subscription id into azure consumption url*/

	client := &http.Client{}
	sub_url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/providers/Microsoft.Consumption/usageDetails?api-version=2019-01-01" ,subscription_id)
	req, err := http.NewRequest("GET",sub_url  ,nil)
	req.Header.Add("Authorization", string(result))
	res, err := client.Do(req)
	
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	

/* consumption api response in json format*/	
var response Consumptn 
json.Unmarshal(body, &response)

	
for i, p := range response.Consumptions {

/*replacing spaces with '+' in instancelocation*/

placeholder := strings.Replace(p.Properties.InstanceLocation ," ","+",1)
fmt.Println("Instance Location = ",placeholder)

/*rate card api url*/
url := fmt.Sprintf("https://prices.azure.com/api/retail/prices?&currencyCode='%s'&$filter=meterId+eq+'%s'+and+location+eq+'%s'",p.Properties.Currency,p.Properties.MeterId,placeholder)
fmt.Println("rate card url = ",url) 
method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("Cookie", "ARRAffinity=40237ffdc57de1390eeff374e782e979bae0af189a51754a0bc4cff0e861cdf3; ARRAffinitySameSite=40237ffdc57de1390eeff374e782e979bae0af189a51754a0bc4cff0e861cdf3")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body1, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body1))


/* rate card  api response in json format*/	
var respons Rate_card
json.Unmarshal(body1, &respons)

	
for j, s := range respons.Ratecard {
	
	
        fmt.Println("Ratecard rec count =", (j + 1), ":", s.CurrencyCode, s.UnitPrice)

		fmt.Println("Consumptions rec count = ", (i + 1), ":", p.Properties.Currency, p.Properties.InstanceName)
		
/* here we are using split method to get Resource group  and resource name */
		str := p.Properties.InstanceId
		split := strings.Split(str, "/")
		string := p.Id
		split1 := strings.Split(string, "/")
		
/* Here we are inserting into the azure billing table*/
	_, err = db.Exec("INSERT INTO azure_billing_new(Id,Resource_id,Resource_name,Resource_group,tendenid,ConsumedService ,SubscriptionGuid , Currency ,UsageQuantity  ,MeterId , UsageStart , UsageEnd) VALUES ($1, $2, $3, $4,$5, $6, $7, $8,$9, $10,$11,$12) ",
	p.Id, p.Properties.InstanceId,p.Properties.InstanceName, split[4], split1[10], p.Properties.ConsumedService, p.Properties.SubscriptionGuid, p.Properties.Currency, p.Properties.UsageQuantity, p.Properties.MeterId, p.Properties.UsageStart, p.Properties.UsageEnd)
			
	if err != nil {
		// fmt.Println("record already exists")
	}
	


/* here we are updating azure billing table based on meterids of consumption & ratecard */
	sqlStatement :=
	 `UPDATE azure_billing_new
	 SET UnitPrice = $2, Location = $3 ,SkuId =$4,ServiceName= $5,UnitOfMeasure=$6 ,ArmSkuName =$7
	 WHERE MeterId = $1;`
	_, err = db.Exec(sqlStatement, s.MeterId, s.UnitPrice, s.Location,s.SkuId,s.ServiceName,s.UnitOfMeasure,s.ArmSkuName)
	if err != nil {  
		panic(err)
	}	
fmt.Println("next page link--- ", respons.NextPageLink)
	
		}
	}
	
	fmt.Println(" data inserted into the azure_consuptn")
}

/*Consumption api struct*/
type Consumptn struct {
	// BillingCurrency    string
	// CustomerEntityId   string
	// CustomerEntityType string
	Consumptions []Consumptions `json:"value"`
}

type Consumptions struct {
	Id         string
	Properties Properties
}
type Properties struct {
	ConsumedService  string
	InstanceName     string
	SubscriptionGuid string
	InstanceId       string
	Currency         string `json:"currency"`
	UsageQuantity    float64
	MeterId          string
	InstanceLocation string
	UsageStart       string
	UsageEnd         string
}
/*Rate card api struct*/
type Rate_card struct {
	NextPageLink   string
	// CustomerEntityId   string
	// CustomerEntityType string
	Ratecard []Ratecard `json:"Items"`
}

type Ratecard struct {
	CurrencyCode       string
	UnitPrice          float64
	ArmRegionName      string
	Location           string
	EffectiveStartDate string
	MeterId            string
	MeterName          string
	SkuId              string
	SkuName            string
	ServiceName        string
	ServiceId          string
	UnitOfMeasure      string
	ArmSkuName         string
}

