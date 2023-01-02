
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
	"reflect"
	"sort"
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
	usageStart_date = "2022-11-01T00:00:00.0000000Z"     /* consumption api billing start date*/
	usageEnd_date = "2022-12-31T23:59:59.0000000Z"
)



func checkError(err error) {
	if err != nil {
		panic(err)
	}
}


type void struct{}

func main() {
/*creating connection for database*/
	
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", HOST, USER, PASSWORD, DATABASE)
	db, err := sql.Open("postgres", connectionString)
	checkError(err)
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to Postgres database \n")

/*creating table for azure billing*/

	_, err = db.Exec("CREATE TABLE  IF  NOT EXISTS azure_billing_new_coster(id VARCHAR(300),Resource_id  character varying(400), Resource_name varchar(300),Resource_group varchar(300),tendenid varchar(300) ,consumedService VARCHAR(300), subscriptionid VARCHAR(300),BillingCurrency  VARCHAR(300),usageQuantity   NUMERIC,meterId VARCHAR(300), BillingPeriodStartDate VARCHAR(300) , BillingPeriodEndDate VARCHAR(300),UnitPrice float ,Location varchar(40),SkuId varchar(40),ServiceName varchar(40),UnitOfMeasure VARCHAR(40),ArmSkuName VARCHAR(40),Cost NUMERIC , month_res varchar(90));")
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
 fmt.Println(string(body),"consumption response body=")
	

/* consumption api response in json format*/	
var response Consumptn 
json.Unmarshal(body, &response)

/* dictionary creation*/
pro_dev_p := []string{}
Month_cnt := []string{}
Month_cnter := []string{}

n := make(map[string]float32)
m:=make(map[string]string)


for i, p := range response.Consumptions {



	// m := {
	// 	p.Properties.ResourceId: p.Properties.BillingPeriodStartDate,
		
	//   }
	
	//   fmt.Println("ditionary  {{{{{{{{{{{{{{{{{{{{}}}}}}}}}}}}}}}}}}}}}}}]" ,m)

//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

if slices.Contains(Month_cnter,(p.Properties.BillingPeriodStartDate) ) {
	fmt.Println("billing period start date")
} else {
	Month_cnter = append(Month_cnter,(p.Properties.BillingPeriodStartDate))

}
fmt.Println(Month_cnter,"list #$#****]]]]]]]]]]]]]]]]]]]]]]]]]]]]]start date*******=========================*******")


	if slices.Contains(Month_cnt,(p.Properties.BillingPeriodEndDate) ) {
		fmt.Println("billing period end date")
	} else {
		Month_cnt = append(Month_cnt,(p.Properties.BillingPeriodEndDate))
	
	}
	fmt.Println(Month_cnt,"list #$#????????????????????????????????????????#end date#####&&&&&&===============")
	
	
	
	// fmt.Println(Month_cnter[0:],"slicing ______________________________________________________________--")

// month_one := (Month_cnter[0:])
// month_two := (Month_cnt[0:])
stringArray := Month_cnt
	justString := fmt.Sprint(stringArray)

	fmt.Println("value \t=", justString, "\ntype \t=", reflect.TypeOf(justString))

	stringArray1 := Month_cnter
	justString1 := fmt.Sprint(stringArray1)

	fmt.Println("value \t=", justString1, "\ntype \t=", reflect.TypeOf(justString1))
	//____________________________________+++++++++++++++++++++________________________________________
	// fmt.Println(Reverse(justString))
	// fmt.Println(Reverse(justString1))
	fmt.Printf("%s\n", strings.Join(reverse(strings.Split(justString, " ")), " "))
	fmt.Printf("%s\n", strings.Join(reverse(strings.Split(justString1, " ")), " "))
	
	


//___________________________________________________________________________________________________________________

	// str1 := p.Properties.BillingPeriodstartDate
	// fmt.Println("BillingPeriodEndDate",p.Properties.BillingPeriodstartDate)
 
    // // using ParseInt method
    int4, err := strconv.Atoi(justString[0:4])
    int5, err := strconv.Atoi(justString[5:7])
    int6, err := strconv.Atoi(justString[8:10])

	// str2 := p.Properties.BillingPeriodStartDate
	// fmt.Println("BillingPeriodendDate",p.Properties.BillingPeriodendDate)
    // // using ParseInt method
    int1, err := strconv.Atoi(justString1[0:4])
    int2, err := strconv.Atoi(justString1[5:7])
    int3, err := strconv.Atoi(justString1[8:10])


	// fmt.Println(int4)

	/*__________________________________________________________________________________________________*/

	firstDate :=  time.Date(int(int1),time.Month(int2),int(int3), 0,0,0,0, time.UTC)
    secondDate :=  time.Date(int(int4),time.Month(int5),int(int6), 24,0,0,0, time.UTC)

	difference := (secondDate.Sub(firstDate))
	fmt.Println(difference.Hours() / 24)
	siva := int(difference.Hours() / 24 / 30)
	fmt.Println(siva)
	// if difference == 31 {
	// 	difference := 1
	// 	fmt.Println(difference, "siva---------")
	// }
	fmt.Printf("difference = %v\n", difference)


	



    // difference := firstDate.Sub(secondDate)

    // fmt.Printf("Years: %d\n", int(difference.Hours()/24/365))
    // fmt.Printf("Months: %d\n", int(difference.Hours()/24/30), int(difference.Hours()/24/28), int(difference.Hours()/24/31))
    // fmt.Printf("Weeks: %d\n", int(difference.Hours()/24/7))
    // fmt.Printf("Days: %d\n", int(difference.Hours()/24))
    
// month_conunt1 := int(difference.Hours()/24/30)
//  month_conunt2 := int(difference.Hours()/24/28)    
//   month_conunt3 :=int(difference.Hours()/24/31)

//  if month_conunt1 == 1 ; month_conunt2 == 1  ; month_conunt3 == 1 {
// 	fmt.Println(month_conunt,"month count")
//  }
	//_________________________________________________________________________________________________________________

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
	
	
        fmt.Println("Ratecard rec count =", (j + 1), ":", s.CurrencyCode, s.UnitPrice,"\n")

		fmt.Println("Consumptions rec count = ", (i + 1), ":", p.Properties.BillingCurrency, p.Properties.ResourceGroup,"\n")
		
/* here we are using split method to get Resource group  and resource name */
		// str := p.Properties.InstanceId
		// split := strings.Split(str, "/")
		stringa := p.Id
		split1 := strings.Split(stringa, "/")
		


if slices.Contains(pro_dev_p,strings.ToLower(p.Properties.ResourceId)) {
	fmt.Println("resourceid in lower case")
} else {
	pro_dev_p = append(pro_dev_p,strings.ToLower(p.Properties.ResourceId))

}

for i, k := range pro_dev_p {
	
	if k == strings.ToLower(p.Properties.ResourceId) {
		n[k] =n[k]+ float32(p.Properties.Cost)
		m[k] = (p.Properties.BillingPeriodStartDate) 


		

		fmt.Println(i)
		
	} else {
		fmt.Println("not matching with resource id  \n")
	}

/*________________________________________________________________________________________________*/



	
	
	

//________________________________________________________________________________________________

		
/* Here we are inserting into the azure billing table*/
	_, err = db.Exec("INSERT INTO azure_billing_new_coster(Id,Resource_id,Resource_name,Resource_group,tendenid,ConsumedService ,Subscriptionid , BillingCurrency ,UsageQuantity ,MeterId , BillingPeriodStartDate , BillingPeriodEndDate) VALUES ($1, (Lower($2)), $3, (Lower($4)),$5, $6, (Lower($7)), $8,$9, (Lower($10)),$11,$12) ",
	p.Id, p.Properties.ResourceId,p.Properties.ResourceName, p.Properties.ResourceGroup, split1[10], p.Properties.ConsumedService, p.Properties.SubscriptionId, p.Properties.BillingCurrency, p.Properties.Quantity, p.Properties.MeterId, p.Properties.BillingPeriodStartDate, p.Properties.BillingPeriodEndDate)
		

	if err != nil {
		// fmt.Println("record already exists")
	}
	
	
/* here we are updating azure billing table based on resource ids  of consumption  */
	sqlStatement :=
	 `UPDATE azure_billing_new_coster
	  SET UnitPrice = $1, Location = $2 ,SkuId =$3,ServiceName= $4,UnitOfMeasure=$5 ,ArmSkuName =$6 
	  WHERE MeterId = $7   ` 
	_, err = db.Exec(sqlStatement, s.UnitPrice, s.Location,s.SkuId,s.ServiceName,s.UnitOfMeasure,s.ArmSkuName,s.MeterId)
	if err != nil {  
		panic(err)
	}	

//_________________________________________________________________________________________________________________________________
for key, value := range n {
	
 sqlStatements :=
	 `UPDATE azure_billing_new_coster
	  SET Cost =$2 ,month_res = $3
	  WHERE resource_id= $1  `
	  
	  
	_, err = db.Exec(sqlStatements, key,value,justString)
	if err != nil {  
		panic(err)
	}	
}

// }	


}
	}
	
	
}
	fmt.Println(" data inserted into the azure_consuptn")

	fmt.Println("dictionary =================" , m)

	set := make(map[string]void)
	for _, element := range m {
		set[element] = void{}
	}

	fmt.Println(reflect.ValueOf(set).MapKeys())

	
	
	
	
	
}
//_____________________________________________________________________________________________________________

















func Reverse(input string) string {
	s := strings.Split(input, " ")
	sort.Sort(sort.Reverse(sort.StringSlice(s)))
	return strings.Join(s, " ")
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
//_________________________________________________________________________________________________________________



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
	Cost                        float32  `json:"cost`
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






// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	start := time.Date(2022, 11, 01, 0, 0, 0, 0, time.UTC)
// 	end := time.Date(2022, 12, 31, 24, 00, 00, 00, time.UTC)

// 	difference := (end.Sub(start))
// 	fmt.Println(difference.Hours() / 24)
// 	siva := int(difference.Hours() / 24 / 30)
// 	fmt.Println(siva)
// 	if difference == 31 {
// 		difference := 1
// 		fmt.Println(difference, "siva---------")
// 	}
// 	fmt.Printf("difference = %v\n", difference)

// }