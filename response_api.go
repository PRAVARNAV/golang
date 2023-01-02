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
	"time"
	"strconv"
	"golang.org/x/exp/slices"
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
	// url_endpoint = "https://login.microsoftonline.com/eab3c9cb-cfc6-4a90-b01d-41c77afef7f3/oauth2/v2.0/token"
	subscription_id = "8524df34-d0be-4424-9c6d-50bd515836fd"
	usageStart_date = "2022-11-01T00:00:00.0000000Z"
	usageEnd_date = "2022-12-31T23:59:59.0000000Z"
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

	_, err = db.Exec("CREATE TABLE  IF  NOT EXISTS azure_billing_new_re(id VARCHAR(300),Resource_id  character varying(400), Resource_name varchar(300),Resource_group varchar(300),tendenid varchar(300) ,consumedService VARCHAR(300), subscriptionid VARCHAR(300),BillingCurrency  VARCHAR(300),usageQuantity   NUMERIC,meterId VARCHAR(300), BillingPeriodStartDate VARCHAR(300) , BillingPeriodEndDate VARCHAR(300),UnitPrice float ,Location varchar(40),SkuId varchar(40),ServiceName varchar(40),UnitOfMeasure VARCHAR(40),ArmSkuName VARCHAR(40),Cost numeric,months int );")
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

var url_endpoint string 
fmt.Println("enter your url endpoint: ")
fmt.Scanln(&url_endpoint)




	// url := url_endpoint
	method := "POST"
	client_details:= fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&scope=%s",grant_type,client_id,client_secret,scope)
	fmt.Println("client details for bearer token = ",client_details)
	payload := strings.NewReader(client_details)
	client1 := &http.Client{}
	req1, err := http.NewRequest(method, url_endpoint, payload)

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
	sub_url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/providers/Microsoft.Consumption/usageDetails?$filter=properties/usageStart+eq+'%s'+and+properties/usageEnd+eq+'%s'&metric=usage&api-version=2021-10-01" ,subscription_id,usageStart_date,usageEnd_date)
	// fmt.Println(sub_url,"comsumption")
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
// fmt.Println(string(body),"consumption response body=")
	

/* consumption api response in json format*/	
var response Consumptn 
json.Unmarshal(body, &response)

//------------------------------------
var Year_year  int
fmt.Println("please enter year:")
fmt.Scanln(&Year_year)

var Month_month  time.Month
fmt.Println("please enter month:")
fmt.Scanln(&Month_month)

var Date_date  int
fmt.Println("please enter date:")
fmt.Scanln(&Date_date)



pro_dev_p := []string{}
	
for i, p := range response.Consumptions {

	

	//_______________________________________________________________________________________________________________________________
	
        start := time.Date(Year_year,Month_month,Date_date,0, 00, 0, 0, time.UTC)
// start := (time.Date(p.Properties.BillingPeriodStartDate) , time.UTC)
	
	// fmt.Println(time.Now())
	
        // calculate years, month, days and time betwen dates
        year, month, day, hour, min, sec := diff(start, time.Now())

months_diff := ((year * 12) + (month * 1))
		fmt.Println(month,"month difference")
        fmt.Printf("difference %d years, %d months, %d days, %d hours, %d mins and %d seconds.", year, month, day, hour, min, sec)
        // fmt.Printf("")

        // calculate total number of days
	// duration := time.Now().Sub(start)
	// fmt.Printf("difference %d days", int(duration.Hours()/24) )
//_________________________________________________________________________________
str1 := p.Properties.BillingPeriodEndDate
 
    // using ParseInt method
    int1, err := strconv.ParseInt(str1, 0, 0)

	str2 := p.Properties.BillingPeriodStartDate
    // using ParseInt method
    int2, err := strconv.ParseInt(str2, 0, 0)

pro_dev := int1 - int2
fmt.Println("months difference=",pro_dev )



	//________________________________________________________________________________________________________________________________
/*replacing spaces with '+' in instancelocation*/

placeholder := strings.Replace(p.Properties.ResourceLocation ," ","+",1)
fmt.Println("Instance Location = ",placeholder)

/*rate card api url*/
url := fmt.Sprintf("https://prices.azure.com/api/retail/prices?&currencyCode='%s'&$filter=meterId+eq+'%s'+and+location+eq+'%s'",p.Properties.BillingCurrency,p.Properties.MeterId,placeholder)
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
//   fmt.Println(string(body1))


/* rate card  api response in json format*/	
var respons Rate_card
json.Unmarshal(body1, &respons)

	
for j, s := range respons.Ratecard {
	
	
        fmt.Println("Ratecard rec count =", (j + 1), ":", s.CurrencyCode, s.UnitPrice)

		fmt.Println("Consumptions rec count = ", (i + 1), ":", p.Properties.BillingCurrency, p.Properties.ResourceGroup)
		
/* here we are using split method to get Resource group  and resource name */
		// str := p.Properties.InstanceId
		// split := strings.Split(str, "/")
		stringa := p.Id
		split1 := strings.Split(stringa, "/")
		// var cost_sum float64
// for k, p := range response.Consumptions {
  	
		// fmt.Println(p.Properties.Cost,"siva developer")
		// cost_sum := ((float64(i))+p.Properties.Cost)
		// cost_sum1 := ((p.Properties.Cost))

		// caluculation  := (p.Properties.Cost)
		// for i=0 ; i < 100; i++{
		//    caluculation += (float64(i))
		// }

var caluculation float64

		
/* Here we are inserting into the azure billing table*/
	_, err = db.Exec("INSERT INTO azure_billing_new_re(Id,Resource_id,Resource_name,Resource_group,tendenid,ConsumedService ,Subscriptionid , BillingCurrency ,UsageQuantity ,MeterId , BillingPeriodStartDate , BillingPeriodEndDate,Cost ,months) VALUES ($1, (Lower($2)), $3, (Lower($4)),$5, $6, (Lower($7)), $8,$9, (Lower($10)),$11,$12,$13,$14) ",
	p.Id, p.Properties.ResourceId,p.Properties.ResourceName, p.Properties.ResourceGroup, split1[10], p.Properties.ConsumedService, p.Properties.SubscriptionId, p.Properties.BillingCurrency, p.Properties.Quantity, p.Properties.MeterId, p.Properties.BillingPeriodStartDate, p.Properties.BillingPeriodEndDate,caluculation,months_diff)
		
	 


	if err != nil {
		// fmt.Println("record already exists")
	}
	
//  pro_dev_p := []{}
if slices.Contains(pro_dev_p,strings.ToLower(p.Properties.ResourceId)) {
	fmt.Println("already exists")
} else {
	pro_dev_p = append(pro_dev_p,strings.ToLower(p.Properties.ResourceId))

}



	
/* here we are updating azure billing table based on meterids of consumption & ratecard */
	sqlStatement :=
	 `UPDATE azure_billing_new_re
	  SET UnitPrice = $1, Location = $2 ,SkuId =$3,ServiceName= $4,UnitOfMeasure=$5 ,ArmSkuName =$6 ,Cost =$7
	  WHERE (MeterId = $8) OR resource_id = $9   `
	  caluculation  = (p.Properties.Cost)
	  for i=0 ; i < 100; i++{
		 caluculation += (float64(i))
	  }
	 
	_, err = db.Exec(sqlStatement, s.UnitPrice, s.Location,s.SkuId,s.ServiceName,s.UnitOfMeasure,s.ArmSkuName,caluculation,s.MeterId,p.Properties.ResourceId)
	if err != nil {  
		panic(err)
	}	
// fmt.Println("next page link--- ", respons.NextPageLink)
	
		}
	}
	


	fmt.Println(" data inserted into the azure_consuptn")
	fmt.Println(pro_dev_p,"pro dev----------------")

}



//___________________________________________________________________________________________________
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}


//__________________________________________________________________________________________________________

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
	BillingPeriodStartDate      string
	BillingPeriodEndDate        string
	SubscriptionId              string
	MeterId                     string
	Cost                        float64  `json:"cost`
	Quantity                    float64
	BillingCurrency             string `json:"billingCurrency"`
	ResourceLocation            string
	ConsumedService             string
	ResourceId                  string
	ResourceName                string
	ResourceGroup               string
	
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

