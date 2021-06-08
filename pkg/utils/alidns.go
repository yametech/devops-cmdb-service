package utils

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"log"
	"strconv"
)

type AliDnsClient struct {
	RegionId        string
	AccessKeyId     string
	AccessKeySecret string
	*alidns.Client
}

func NewAliDnsClient(regionId string, accessKeyId string, accessKeySecret string) (client *AliDnsClient) {
	client = &AliDnsClient{RegionId: regionId, AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}
	c, err := alidns.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		log.Println(err.Error())
	}
	client.Client = c
	return
}

func (c *AliDnsClient) DomainList() (domains []alidns.DomainInDescribeDomains) {
	domains = make([]alidns.DomainInDescribeDomains, 0)
	request := alidns.CreateDescribeDomainsRequest()
	request.PageSize = "100"
	request.PageNumber = "1"
	request.AcceptFormat = "json"

	response, err := c.Client.DescribeDomains(request)
	if err != nil {
		log.Println(err.Error())
	}
	if response == nil {
		log.Println("response is nil")
		return
	}

	domains = append(domains, response.Domains.Domain...)

	// continue
	totalPage := response.TotalCount / response.PageSize
	if response.TotalCount%response.PageSize > 0 {
		totalPage += 1
	}
	for totalPage-1 > 0 {
		totalPage -= 1
		pageNumber, _ := request.PageNumber.GetValue()
		pageNumber += 1
		request.PageNumber = requests.Integer(strconv.Itoa(pageNumber))

		response, err = c.Client.DescribeDomains(request)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if response == nil {
			log.Println("response is nil")
			continue
		}
		domains = append(domains, response.Domains.Domain...)
	}
	return
}

func (c *AliDnsClient) DomainLogs() (domainLogs []alidns.DomainLog) {
	domainLogs = make([]alidns.DomainLog, 0)
	request := alidns.CreateDescribeDomainLogsRequest()
	request.PageSize = "100"
	request.PageNumber = "1"
	request.AcceptFormat = "json"
	//request.Type = "domain"

	response, err := c.Client.DescribeDomainLogs(request)
	if err != nil {
		log.Println(err.Error())
	}
	if response == nil {
		log.Println("response is nil")
		return
	}

	domainLogs = append(domainLogs, response.DomainLogs.DomainLog...)
	totalPage := response.TotalCount / response.PageSize
	if response.TotalCount%response.PageSize > 0 {
		totalPage += 1
	}
	for totalPage-1 > 0 {
		totalPage -= 1
		pageNumber, _ := request.PageNumber.GetValue()
		pageNumber += 1
		request.PageNumber = requests.Integer(strconv.Itoa(pageNumber))

		response, err = c.Client.DescribeDomainLogs(request)
		if err != nil {
			log.Println(err.Error())
		}
		if response == nil {
			log.Println("response is nil")
			continue
		}
		domainLogs = append(domainLogs, response.DomainLogs.DomainLog...)
	}
	return
}

func (c *AliDnsClient) DomainParsingRecords(domainName string) (records []alidns.Record) {
	records = make([]alidns.Record, 0)

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.PageSize = "100"
	request.PageNumber = "1"
	request.AcceptFormat = "json"
	request.DomainName = domainName

	response, err := c.Client.DescribeDomainRecords(request)
	if err != nil {
		log.Println(err.Error())
	}
	if response == nil {
		log.Println("response is nil")
		return
	}

	records = append(records, response.DomainRecords.Record...)
	totalPage := response.TotalCount / response.PageSize
	if response.TotalCount%response.PageSize > 0 {
		totalPage += 1
	}
	for totalPage-1 > 0 {
		totalPage -= 1
		pageNumber, _ := request.PageNumber.GetValue()
		pageNumber += 1
		request.PageNumber = requests.Integer(strconv.Itoa(pageNumber))

		response, err = c.Client.DescribeDomainRecords(request)
		if err != nil {
			log.Println(err.Error())
		}
		if response == nil {
			log.Println("response is nil")
			continue
		}
		records = append(records, response.DomainRecords.Record...)
	}

	return
}
