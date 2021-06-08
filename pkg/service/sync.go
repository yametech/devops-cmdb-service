package service

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SyncService struct {
	ResourceService
	RelationService
	mutex sync.Mutex
}

type Idc struct {
	UUID string
	Idc  string
}

type SelfDomain struct {
	DomainUUID string `json:"domainUUID"`
	IdcUUID    string `json:"idcUUID"`
	Domain     string `db:"domain"`
	Idc        string `db:"idc"`
}

type SelfDomainParsingRecord struct {
	DomainUUID     string `json:"domainUUID"`
	IdcUUID        string `json:"idcUUID"`
	Idc            string `json:"idc" db:"idc"`
	Domain         string `json:"domain" db:"domain"`
	UUID           string `json:"uuid"`
	ParsingRecords string `json:"parsing_records" db:"host_record"`
	RecordType     string `json:"record_type" db:"record_type"`
	RecordValue    string `json:"record_value" db:"record_value"`
}

type Config struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func getRecordType(recordType string) string {
	typeMap := map[string]string{
		"A":     "0",
		"CNAME": "1",
		"AAAA":  "2",
		"NS":    "3",
		"MX":    "4",
		"SRV":   "5",
		"TXT":   "6",
		"CAA":   "7",
		"显性URL": "8",
		"隐性URL": "9",
	}
	return typeMap[recordType]
}

func formatTime(source string) string {
	if source != "" && len(source) >= 10 {
		source = strings.ReplaceAll(source, "T", " ")
		source = strings.ReplaceAll(source, "Z", "")
		source += ":00"
	}
	return source
}

func peelDomain(domain *alidns.DomainInDescribeDomains) map[string]string {
	var dateLeft int64 = 0
	var status = 1
	if domain.InstanceExpired {
		dateLeft = 0
		status = 2
	}
	//TODO 不足一天的忽略
	instanceEndTimeStr := formatTime(domain.InstanceEndTime)
	if instanceEndTimeStr != "" && len(instanceEndTimeStr) >= 10 {
		instanceEndTime, err := time.Parse("2006-01-02T15:04Z", domain.InstanceEndTime)
		if err == nil {
			day := (instanceEndTime.Unix() - time.Now().Unix()) / 86400
			if day >= 0 {
				dateLeft = day + 1
			}
		}
	}

	return map[string]string{
		"domain_name":              domain.DomainName,
		"domain_holder":            "",
		"domain_type":              "",
		"domain_audit_status":      "",
		"domain_group_id":          domain.GroupId,
		"domain_group_name":        domain.GroupName,
		"domain_curr_date_diff":    strconv.FormatInt(dateLeft, 10),
		"domain_status":            "",
		"domain_registrant_type":   "",
		"domain_expiration_Status": strconv.Itoa(status),
		"domain_instance_id":       domain.DomainId,
		"domain_remark":            domain.Remark,
		"domain_premium":           "",
		"domain_product_id":        domain.InstanceId,
		"domain_registration_date": "",
		"domain_expiration_date":   instanceEndTimeStr,
	}
}

func peelDomainLog(domainLog *alidns.DomainLog) map[string]string {

	return map[string]string{
		"details":              domainLog.Message,
		"domain_name":          domainLog.DomainName,
		"operation":            domainLog.Action,
		"operation_ip_address": domainLog.ClientIp,
		"result":               "",
		"time":                 formatTime(domainLog.ActionTime),
	}
}

func peelParsingRecord(record *alidns.Record) map[string]string {

	return map[string]string{
		"domain_id":   "",
		"domain_name": record.DomainName,
		"group_id":    "",
		"group_name":  "",
		"line":        record.Line,
		"locked":      strconv.FormatBool(record.Locked),
		"priority":    strconv.FormatInt(record.Priority, 10),
		"puny_code":   "",
		"rr":          record.RR,
		"record_id":   record.RecordId,
		//"request_id": response.RequestId,
		"status": record.Status,
		"ttl":    strconv.FormatInt(record.TTL, 10),
		"type":   record.Type,
		"value":  record.Value,
	}
}

func getDomainUUID(resources []store.Resource, domainName string) (uuid string) {
	for _, resource := range resources {
		for _, attributeGroup := range resource.AttributeGroupIns {
			for _, ins := range attributeGroup.AttributeIns {
				if ins.Uid == "domain_name" && ins.AttributeInsValue == domainName {
					return resource.UUID
				}
			}
		}
	}
	return ""
}

func (s *SyncService) processResource(operator string, mongoDb *utils.MongoDB, modelUid string, attributeGroupUid string, attributeMap map[string]string) (resource *store.Resource) {
	attributes := make([]map[string]interface{}, 0)
	for k, v := range attributeMap {
		attributes = append(attributes, map[string]interface{}{"uid": k, "attributeInsValue": v})
	}

	insVO := map[string]interface{}{
		"modelUid": modelUid,
		"attributeGroupIns": []interface{}{
			map[string]interface{}{
				"uid":          attributeGroupUid,
				"attributeIns": attributes,
			},
		},
	}

	body, _ := json.Marshal(insVO)
	resource, err := s.ResourceService.AddResource(string(body), operator)
	if err != nil {
		log.Println(err.Error())
		mongoDb.InsertOne("err_"+modelUid, bson.D{{"msg", err.Error()}})
		return
	}
	//mongoDb.InsertOne(modelUid, resource)
	return
}

func getResourceByAttributeMap(resources []*store.Resource, attributeMap map[string]string) *store.Resource {
	if attributeMap == nil || len(attributeMap) == 0 {
		return nil
	}

	for _, resource := range resources {
		attributeMapAll := make(map[string]string)
		for _, attributeGroupIns := range resource.AttributeGroupIns {
			for _, attribute := range attributeGroupIns.AttributeIns {
				attributeMapAll[attribute.Uid] = attribute.AttributeInsValue
			}
		}

		var exist = true
		for k, v := range attributeMap {
			value, ok := attributeMapAll[k]
			if !ok || value != v {
				exist = false
			}
		}
		if exist {
			return resource
		}
	}
	return nil
}

func resourceToAliyunAccountVO(resource *store.Resource) *common.AliyunAccountVO {
	aliyunAccountMap := map[string]string{}
	for _, attributeIns := range resource.AttributeGroupIns[0].AttributeIns {
		aliyunAccountMap[attributeIns.Uid] = attributeIns.AttributeInsValue
	}
	aliAccountVO := &common.AliyunAccountVO{}
	utils.SimpleConvert(aliAccountVO, aliyunAccountMap)
	aliAccountVO.Uuid = resource.UUID
	return aliAccountVO
}

func getLatestSelfDomain() ([]SelfDomain, error) {
	db := utils.GetDefaultDBConn()

	// domain
	rows, err := db.Query("select idc,domain from domain_roots")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Printf("Query failed,err:%v", err)
		return nil, err
	}
	selfDomains := make([]SelfDomain, 0)
	for rows.Next() {
		selfDomain := new(SelfDomain)
		err := rows.Scan(&selfDomain.Idc, &selfDomain.Domain)
		if err != nil {
			log.Printf("Scan failed,err:%v\n", err)
			return nil, err
		}
		selfDomains = append(selfDomains, *selfDomain)
	}
	return selfDomains, nil
}

func getLatestSelfDomainParsingRecord() ([]SelfDomainParsingRecord, error) {
	recordRow, err := utils.GetDefaultDBConn().Query("SELECT idc,domain,host_record,record_type,record_value FROM domain_records GROUP BY idc,domain,host_record,record_type,record_value")
	if err != nil {
		log.Printf("Query failed,err:%v", err)
		return nil, err
	}
	defer func() {
		if recordRow != nil {
			recordRow.Close()
		}
	}()
	records := make([]SelfDomainParsingRecord, 0)
	for recordRow.Next() {
		record := new(SelfDomainParsingRecord)
		err := recordRow.Scan(&record.Idc, &record.Domain, &record.ParsingRecords, &record.RecordType, &record.RecordValue)
		if err != nil {
			log.Printf("Scan failed,err:%v\n", err)
			return nil, err
		}
		records = append(records, *record)
	}
	return records, nil
}

func (s *SyncService) getExistIdc() []Idc {
	mongodb := utils.GetDefaultMongoDB()
	// 机房实例信息
	roomQuery := "MATCH (a:Resource {modelUid:'room'})-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns {uid:'idc'}) RETURN a.uuid,c.attributeInsValue"
	roomResult, err := s.ResourceService.ManualQueryRaw(roomQuery, nil)
	if err != nil {
		mongodb.InsertOne("err_"+"built_domain", bson.D{{"msg", err.Error()}})
	}
	idcs := make([]Idc, 0)
	for _, room := range roomResult {
		idc := Idc{UUID: room[0].(string), Idc: room[1].(string)}
		idcs = append(idcs, idc)
	}
	return idcs
}

func (s *SyncService) getExistSelfDomainAndIdc() ([]SelfDomain, []Idc) {
	mongodb := utils.GetDefaultMongoDB()
	// 已有自建域名实例
	query := "MATCH (a:Resource {modelUid:'built_domain'})-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns {uid:'name'}) " +
		"MATCH (d:Resource {modelUid:'room'})-[:Relation {uid:'room_correlation_built_domain'}]-(a) " +
		"RETURN a.uuid,c.attributeInsValue,d.uuid"
	domainResult, err := s.ResourceService.ManualQueryRaw(query, nil)
	if err != nil {
		mongodb.InsertOne("err_"+"built_domain", bson.D{{"msg", err.Error()}})
	}

	idcs := s.getExistIdc()
	domains := make([]SelfDomain, 0)
	for _, domain := range domainResult {
		selfDomain := SelfDomain{
			DomainUUID: domain[0].(string),
			Domain:     domain[1].(string),
			IdcUUID:    domain[2].(string),
		}

		if selfDomain.IdcUUID != "" {
			for _, idc := range idcs {
				if idc.UUID == selfDomain.IdcUUID {
					selfDomain.Idc = idc.Idc
					break
				}
			}
		}

		domains = append(domains, selfDomain)
	}
	return domains, idcs
}

func getSelfDomainByIdcAndDomain(idc, name string, domains []SelfDomain) *SelfDomain {
	for _, domain := range domains {
		if domain.Idc == idc && domain.Domain == name {
			return &domain
		}
	}
	return nil
}

func getIdcByName(idcName string, idcs []Idc) *Idc {
	for _, idc := range idcs {
		if idc.Idc == idcName {
			return &idc
		}
	}
	return nil
}

func parsingResultToVO(result [][]interface{}) []*SelfDomainParsingRecord {
	records := make([]*SelfDomainParsingRecord, 0)
	parsingMap := map[string]*SelfDomainParsingRecord{}
	for _, row := range result {
		rowMap := map[string]string{}
		rowMap["uuid"] = row[0].(string)
		rowMap[row[1].(string)] = row[2].(string)
		rowMap["domainUUID"] = row[3].(string)
		rowMap["domain"] = row[4].(string)
		rowMap["idcUUID"] = row[5].(string)
		rowMap["idc"] = row[6].(string)

		exit, ok := parsingMap[row[0].(string)]
		if !ok {
			exit = &SelfDomainParsingRecord{}
			utils.SimpleConvert(&exit, rowMap)
			parsingMap[row[0].(string)] = exit
			records = append(records, exit)
		} else {
			//TODO reflect.ValueOf(exit).MethodByName(methodName).Call(inputs)
			switch row[1].(string) {
			case "parsing_records":
				exit.ParsingRecords = row[2].(string)
			case "record_type":
				exit.RecordType = row[2].(string)
			case "record_value":
				exit.RecordValue = row[2].(string)
			}
		}
	}
	return records
}

func (s *SyncService) SyncSelfDomain(operator string) (bool, error) {
	// 同步之前清空mongodb相关数据
	mongodb := utils.GetDefaultMongoDB()
	mongodb.Drop("err_built_domain", "err_built_domain_parsing")

	latestDomains, err := getLatestSelfDomain()
	if err != nil {
		return false, err
	}

	latestRecords, err := getLatestSelfDomainParsingRecord()
	if err != nil {
		return false, err
	}

	existSelfDomains, idcs := s.getExistSelfDomainAndIdc()
	for _, d := range latestDomains {
		// 获取idc
		existDomain := getSelfDomainByIdcAndDomain(d.Idc, d.Domain, existSelfDomains)
		idc := getIdcByName(d.Idc, idcs)
		if existDomain == nil && idc != nil {
			resource := s.processResource(operator, mongodb, "built_domain", "built_domain_info", map[string]string{"name": d.Domain})
			if resource != nil {
				// relate
				_, err := s.RelationService.AddResourceRelation(idc.UUID, resource.UUID, "room_correlation_built_domain")
				if err != nil {
					mongodb.InsertOne("err_"+"built_domain", bson.D{{"msg", err.Error()}})
				}
				newExist := SelfDomain{
					DomainUUID: resource.UUID,
					Domain:     d.Domain,
					IdcUUID:    idc.UUID,
					Idc:        idc.Idc,
				}

				existSelfDomains = append(existSelfDomains, newExist)
			}
		}
	}

	query := "MATCH p=(a:Resource {modelUid:'built_domain_parsing'})-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns) " +
		"MATCH (d:Resource {modelUid:'built_domain'})-[]-(a) " +
		"MATCH (d)-[]-(d1:AttributeGroupIns)-[]-(d2:AttributeIns {uid:'name'}) " +
		"MATCH (e:Resource {modelUid:'room'})-[]-(d) " +
		"MATCH (e)-[]-(e1:AttributeGroupIns)-[]-(e2:AttributeIns {uid:'idc'}) " +
		"RETURN a.uuid,c.uid,c.attributeInsValue,d.uuid,d2.attributeInsValue,e.uuid,e2.attributeInsValue"
	result, err := s.ResourceService.ManualQueryRaw(query, nil)
	if err != nil {
		mongodb.InsertOne("err_"+"built_domain_parsing", bson.D{{"msg", err.Error()}})
	}

	records := parsingResultToVO(result)
	recordMap := map[string]int{}
	for i := 0; i < len(records); i++ {
		r := records[i]
		recordMap[r.Idc+r.RecordType+r.Domain+r.RecordValue+r.ParsingRecords] = 1
	}
	for _, r := range latestRecords {
		existDomain := getSelfDomainByIdcAndDomain(r.Idc, r.Domain, existSelfDomains)
		if existDomain == nil {
			continue
		}
		recordType := getRecordType(r.RecordType)
		if recordMap[r.Idc+recordType+r.Domain+r.RecordValue+r.ParsingRecords] == 0 {
			attributeMap := map[string]string{
				"parsing_records": r.ParsingRecords,
				"record_type":     recordType,
				"record_value":    r.RecordValue,
			}
			resource := s.processResource(operator, mongodb, "built_domain_parsing", "built_domain_parsing_info", attributeMap)
			if resource != nil {
				_, err := s.RelationService.AddResourceRelation(existDomain.DomainUUID, resource.UUID, "built_domain_including_built_domain_parsing")
				if err != nil {
					mongodb.InsertOne("err_"+"built_domain_parsing", bson.D{{"msg", err.Error()}, {"value", r}})
				}
			}
		}
	}
	return true, nil
}

// 同步阿里域名实例信息
func (s *SyncService) SyncAliDomain(operator string) (bool, error) {
	accounts := s.getAliyunAccounts()
	if accounts == nil || len(accounts) == 0 {
		return false, fmt.Errorf("请创建阿里云账号")
	}

	// 同步之前清空mongodb相关数据
	mongodb := utils.GetDefaultMongoDB()
	mongodb.Drop("err_ali_domain", "err_ali_domain_log", "err_ali_parsing_records")
	for _, account := range accounts {
		s.SyncAliDomainByAliyunAccount(&account, operator)
	}
	return true, nil
}

func (s *SyncService) SyncAliDomainByResource(resource *store.Resource, operator string) {
	mongodb := utils.GetDefaultMongoDB()
	mongodb.Drop("err_ali_domain", "err_ali_domain_log", "err_ali_parsing_records")
	s.SyncAliDomainByAliyunAccount(resourceToAliyunAccountVO(resource), operator)
}

func (s *SyncService) SyncAliDomainByAliyunAccount(aliyunAccount *common.AliyunAccountVO, operator string) {
	// 同步之前清空mongodb相关数据
	mongodb := utils.GetDefaultMongoDB()
	//mongodb.Drop("err_ali_domain", "err_ali_domain_log", "err_ali_parsing_records")

	domainResources := make([]store.Resource, 0)
	aliDnsClient := utils.NewAliDnsClient(aliyunAccount.Area, aliyunAccount.AccessKeyId, aliyunAccount.AccessSecret)
	domains := aliDnsClient.DomainList()
	b, _ := json.Marshal(domains)
	log.Println(string(b))
	for _, domain := range domains {
		// 新的数据需要新增
		existQuery := "MATCH (a:Resource)<-[:GroupBy]-(b:AttributeGroupIns)<-[:GroupBy]-(c:AttributeIns) " +
			"WHERE a.modelUid='ali_domain' and  c.uid = 'domain_name' AND c.attributeInsValue = $domainName RETURN *"
		result, err := s.ResourceService.ManualQueryRaw(existQuery, map[string]interface{}{"domainName": domain.DomainName})
		if err != nil {
			mongodb.InsertOne("err_"+"ali_domain", bson.D{{"msg", err.Error()}})
			continue
		}
		if result != nil {
			// 缓存起来后面用
			for _, row := range result {
				resource := utils.GetResourceFromNeo4jRow(row)
				domainResources = append(domainResources, *resource)
			}
			continue
		}

		resource := s.processResource(operator, mongodb, "ali_domain", "ali_domain_info", peelDomain(&domain))
		if resource != nil {
			domainResources = append(domainResources, *resource)
			_, err := s.RelationService.AddResourceRelation(aliyunAccount.Uuid, resource.UUID, "aliyun_account_including_ali_domain")
			if err != nil {
				mongodb.InsertOne("err_"+"ali_domain", bson.D{{"msg", err.Error()}, {"value", resource}})
			}
		}
	}

	// DomainLog
	domainLogs := aliDnsClient.DomainLogs()
	b, _ = json.Marshal(domainLogs)
	log.Println(string(b))
	// 新的数据需要新增
	// TODO 一次查询
	query := "MATCH p=(a:Resource)-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns) WHERE a.modelUid = 'ali_domain_log' RETURN a,b,c"
	result, err := s.ResourceService.ManualQueryRaw(query, nil)
	if err != nil {
		mongodb.InsertOne("err_"+"ali_domain", bson.D{{"msg", err.Error()}})
	}
	logResources := utils.GetResourceFromNeo4jResult(result)

	for _, domainLog := range domainLogs {
		attributeMap := map[string]string{
			"domain_name": domainLog.DomainName,
			"operation":   domainLog.Action,
			"time":        formatTime(domainLog.ActionTime),
		}
		exist := getResourceByAttributeMap(logResources, attributeMap)

		if exist == nil {
			resource := s.processResource(operator, mongodb, "ali_domain_log", "ali_domain_log_info", peelDomainLog(&domainLog))
			if resource != nil {
				uuid := getDomainUUID(domainResources, domainLog.DomainName)
				if uuid != "" {
					_, err := s.RelationService.AddResourceRelation(uuid, resource.UUID, "ali_domain_including_ali_domain_log")
					if err != nil {
						mongodb.InsertOne("err_"+"ali_domain_log", bson.D{{"msg", err.Error()}, {"value", resource}})
					}
				} else {
					mongodb.InsertOne("err_"+"ali_domain_log", bson.D{{"msg", "can't find domain by " + domainLog.DomainName}})
				}
			}
		}
	}

	// DomainParsingRecords
	for _, domain := range domains {
		records := aliDnsClient.DomainParsingRecords(domain.DomainName)
		b, _ := json.Marshal(records)
		log.Println(string(b))
		for _, record := range records {
			existQuery := "MATCH (a:Resource)<-[:GroupBy]-(b:AttributeGroupIns)<-[:GroupBy]-(c:AttributeIns) " +
				"WHERE a.modelUid='ali_parsing_records' and  c.uid = 'record_id' AND c.attributeInsValue = $recordId RETURN *"
			result, err := s.ResourceService.ManualQueryRaw(existQuery, map[string]interface{}{"recordId": record.RecordId})
			if err != nil {
				mongodb.InsertOne("err_"+"ali_parsing_records", bson.D{{"msg", err.Error()}})
				continue
			}
			if result == nil {
				resource := s.processResource(operator, mongodb, "ali_parsing_records", "ali_parsing_records_info", peelParsingRecord(&record))
				if resource != nil {
					domainUUID := getDomainUUID(domainResources, record.DomainName)
					if domainUUID != "" {
						_, err := s.RelationService.AddResourceRelation(domainUUID, resource.UUID, "ali_domain_including_ali_parsing_records")
						if err != nil {
							mongodb.InsertOne("err_"+"ali_parsing_records", bson.D{{"msg", err.Error()}, {"value", resource}})
						}
					} else {
						mongodb.InsertOne("err_"+"ali_parsing_records", bson.D{{"msg", "can't find domain by " + record.DomainName}})
					}
				}
			}
		}
	}
}

func (s *SyncService) getAliyunAccounts() []common.AliyunAccountVO {
	queryAliAccountStr := "MATCH p=(a:Resource)-[r1]-(b:AttributeGroupIns)-[r2]-(c:AttributeIns) WHERE a.modelUid ='aliyun_account' RETURN p"
	rows, err := s.ResourceService.ManualQueryRaw(queryAliAccountStr, nil)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// 组装进Resource
	resources := make([]*store.Resource, 0)
	resourceMap := map[string]*store.Resource{}
	for _, row := range rows {
		path := row[0].(*gogm.PathWrap)
		resource, ok := resourceMap[strconv.FormatInt(path.Nodes[0].Id, 10)]
		if !ok {
			resource = &store.Resource{}
			resourceMap[strconv.FormatInt(path.Nodes[0].Id, 10)] = resource
			resources = append(resources, resource)
		}
		utils.SimpleConvert(resource, path.Nodes[0].Props)
		groupIns := &store.AttributeGroupIns{}
		utils.SimpleConvert(groupIns, path.Nodes[1].Props)
		attributeIns := &store.AttributeIns{}
		utils.SimpleConvert(attributeIns, path.Nodes[2].Props)
		groupIns.AddAttributeIns(attributeIns)
		resource.AddAttributeGroupIns(groupIns)
	}

	// 转换成VO
	accounts := make([]common.AliyunAccountVO, 0)
	for _, resource := range resources {
		aliAccountVO := resourceToAliyunAccountVO(resource)
		accounts = append(accounts, *aliAccountVO)
	}

	return accounts
}

func (s *SyncService) SyncResource(modelUid, operator string) (bool, error) {
	ec := utils.EtcdClient{}
	uidMap := map[string]string{
		"built_domain":         "built_domain",
		"built_domain_parsing": "built_domain",
		"ali_domain":           "ali_domain",
		"ali_domain_log":       "ali_domain",
		"ali_parsing_records":  "ali_domain",
	}

	if uidMap[modelUid] != "" && !ec.ApplyLock(uidMap[modelUid]) {
		return false, fmt.Errorf("目前正在执行同步过程，请稍后再试")
	}

	switch {
	case uidMap[modelUid] == "built_domain":
		go func() {
			s.SyncSelfDomain(operator)
			ec.DeleteLease(uidMap[modelUid])
		}()
	case uidMap[modelUid] == "ali_domain":
		go func() {
			s.SyncAliDomain(operator)
			ec.DeleteLease(uidMap[modelUid])
		}()
	}

	return false, nil
}

func (s *SyncService) SyncButton(modelUid, operator string) bool {
	if modelUid == "" {
		return false
	}
	config := Config{}
	mongodb := utils.GetDefaultMongoDB()
	mongodb.FindOne("config", bson.M{"name": "sync_model"}, &config)

	if config.Value != nil {
		uids := make([]interface{}, 0)
		utils.SimpleConvert(&uids, config.Value)
		for _, uid := range uids {
			if uid == modelUid {
				return true
			}
		}
	}

	return false
}
