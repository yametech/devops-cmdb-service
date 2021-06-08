package service

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"log"
	"strings"
	"testing"
)

func TestCreateFullModel(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	var host, username, password = "localhost", "neo4j", "123456"
	store.Neo4jInit(host, username, password)

	body := "{\"id\":134591,\"uuid\":\"7563b700-58b5-42ed-8746-04fc1b89f55b\",\"uid\":\"built_domain_parsing\",\"name\":\"自建域名解析记录\",\"iconUrl\":\"\",\"attributeGroups\":[{\"id\":453,\"uuid\":\"b1c7c78b-f84c-4b19-a928-7776d8cb278f\",\"uid\":\"built_domain_parsing_info\",\"name\":\"基本信息\",\"modelUid\":\"built_domain_parsing\",\"attributes\":[{\"id\":107020,\"uuid\":\"b9dfb20e-32b2-447e-a693-1cb3072d10be\",\"uid\":\"parsing_records\",\"name\":\"解析记录\",\"valueType\":\"短字符\",\"editable\":false,\"required\":false,\"unique\":false,\"defaultValue\":null,\"unit\":\"\",\"maximum\":\"\",\"minimum\":\"\",\"enums\":null,\"listValues\":null,\"tips\":\"\",\"regular\":\"\",\"comment\":\"\",\"visible\":false,\"modelUid\":\"built_domain_parsing\",\"creator\":\"linxieting\",\"editor\":\"linxieting\",\"createTime\":1621333711726,\"updateTime\":1621333711726},{\"id\":107021,\"uuid\":\"c6151e26-901e-45e7-9ebe-9930aea367b2\",\"uid\":\"record_type\",\"name\":\"记录类型\",\"valueType\":\"枚举\",\"editable\":false,\"required\":false,\"unique\":false,\"defaultValue\":null,\"unit\":\"\",\"maximum\":\"\",\"minimum\":\"\",\"enums\":\"[{\\\"id\\\":\\\"0\\\",\\\"value\\\":\\\"A\\\"},{\\\"id\\\":\\\"1\\\",\\\"value\\\":\\\"CNAME\\\"},{\\\"id\\\":\\\"2\\\",\\\"value\\\":\\\"AAAA\\\"},{\\\"id\\\":\\\"3\\\",\\\"value\\\":\\\"NS\\\"},{\\\"id\\\":\\\"4\\\",\\\"value\\\":\\\"MX\\\"},{\\\"id\\\":\\\"5\\\",\\\"value\\\":\\\"SRV\\\"},{\\\"id\\\":\\\"6\\\",\\\"value\\\":\\\"TXT\\\"},{\\\"id\\\":\\\"7\\\",\\\"value\\\":\\\"CAA\\\"},{\\\"id\\\":\\\"8\\\",\\\"value\\\":\\\"显性URL\\\"},{\\\"id\\\":\\\"9\\\",\\\"value\\\":\\\"隐性URL\\\"}]\",\"listValues\":null,\"tips\":\"\",\"regular\":\"\",\"comment\":\"A:将域名指向一个IPV4地址\\n\\nCNAME:将域名指向另外一个域名\\n\\nAAAA:将域名指向一个IPV6地址\\n\\nNS:将子域名指定其他DNS服务器解析\\n\\nMX:将域名指向邮件服务器地址\\n\\nSRV:记录提供特定的服务的服务器\\n\\nTXT:文本长度限制512，通常做SPF记录（反垃圾邮件）\\n\\nCAA:CA证书颁发机构授权校验\\n显性URL:将域名重定向到另外一个地址\\n\\n隐性URL:与显性URL类似，但是会隐藏真实目标地址\",\"visible\":false,\"modelUid\":\"built_domain_parsing\",\"creator\":\"linxieting\",\"editor\":\"linxieting\",\"createTime\":1621333712023,\"updateTime\":1621333712023},{\"id\":107022,\"uuid\":\"3b0cf319-1833-4f36-8c49-f733212eac7d\",\"uid\":\"isp\",\"name\":\"解析线路（ISP）\",\"valueType\":\"枚举\",\"editable\":true,\"required\":false,\"unique\":false,\"defaultValue\":null,\"unit\":\"\",\"maximum\":\"\",\"minimum\":\"\",\"enums\":\"[{\\\"id\\\":\\\"0\\\",\\\"value\\\":\\\"默认\\\"},{\\\"id\\\":\\\"1\\\",\\\"value\\\":\\\"境外\\\"},{\\\"id\\\":\\\"2\\\",\\\"value\\\":\\\"搜索引擎\\\"}]\",\"listValues\":null,\"tips\":\"\",\"regular\":\"\",\"comment\":\"\",\"visible\":false,\"modelUid\":\"built_domain_parsing\",\"creator\":\"linxieting\",\"editor\":\"linxieting\",\"createTime\":1621333712332,\"updateTime\":1621333712332},{\"id\":107023,\"uuid\":\"3b25b32d-1de3-4b8c-8c6a-949ecb829df9\",\"uid\":\"record_value\",\"name\":\"记录值\",\"valueType\":\"短字符\",\"editable\":false,\"required\":false,\"unique\":false,\"defaultValue\":null,\"unit\":\"\",\"maximum\":\"\",\"minimum\":\"\",\"enums\":null,\"listValues\":null,\"tips\":\"\",\"regular\":\"\",\"comment\":\"\",\"visible\":false,\"modelUid\":\"built_domain_parsing\",\"creator\":\"linxieting\",\"editor\":\"linxieting\",\"createTime\":1621333712638,\"updateTime\":1621333712638},{\"id\":107024,\"uuid\":\"50674cdb-3cad-4f0f-9775-f8befe463ba9\",\"uid\":\"ttl\",\"name\":\"TTL\",\"valueType\":\"枚举\",\"editable\":true,\"required\":false,\"unique\":false,\"defaultValue\":null,\"unit\":\"\",\"maximum\":\"\",\"minimum\":\"\",\"enums\":\"[{\\\"id\\\":\\\"0\\\",\\\"value\\\":\\\"秒\\\"},{\\\"id\\\":\\\"1\\\",\\\"value\\\":\\\"分钟\\\"},{\\\"id\\\":\\\"2\\\",\\\"value\\\":\\\"小时\\\"},{\\\"id\\\":\\\"3\\\",\\\"value\\\":\\\"天\\\"}]\",\"listValues\":null,\"tips\":\"\",\"regular\":\"\",\"comment\":\"\",\"visible\":false,\"modelUid\":\"built_domain_parsing\",\"creator\":\"linxieting\",\"editor\":\"linxieting\",\"createTime\":1621333712946,\"updateTime\":1621333712946},{\"id\":107025,\"uuid\":\"7fd90431-3c62-4921-83dd-cd91e2ae8f32\",\"uid\":\"state\",\"name\":\"状态\",\"valueType\":\"枚举\",\"editable\":true,\"required\":false,\"unique\":false,\"defaultValue\":null,\"unit\":\"\",\"maximum\":\"\",\"minimum\":\"\",\"enums\":\"[{\\\"id\\\":\\\"0\\\",\\\"value\\\":\\\"正常\\\"},{\\\"id\\\":\\\"1\\\",\\\"value\\\":\\\"暂停\\\"}]\",\"listValues\":null,\"tips\":\"\",\"regular\":\"\",\"comment\":\"\",\"visible\":false,\"modelUid\":\"built_domain_parsing\",\"creator\":\"linxieting\",\"editor\":\"linxieting\",\"createTime\":1621333713323,\"updateTime\":1621333713323}],\"creator\":\"linxieting\",\"editor\":\"linxieting\",\"createTime\":1619602813,\"updateTime\":1619602813}],\"resources\":null,\"creator\":\"linxieting\",\"editor\":\"linxieting\",\"createTime\":1619345500,\"updateTime\":1619345500}"
	modelVO := &store.Model{}
	err := json.Unmarshal([]byte(body), modelVO)
	if err != nil {
		panic(err)
	}

	vo := &common.AddModelVO{ModelGroupUUID: "19b539de-402a-4766-b145-05b187e3764a", Uid: modelVO.Uid, Name: modelVO.Name}
	service := ModelService{}
	model, err := service.CreateModel(vo, "test")
	if err != nil {
		panic(err)
	}

	for _, attributeGroupIn := range modelVO.AttributeGroups {
		attributeGroupVO := &common.AddAttributeGroupVO{ModelUUID: model.UUID, Uid: attributeGroupIn.Uid, Name: attributeGroupIn.Name}
		attributeService := AttributeService{}
		group, err := attributeService.CreateAttributeGroup(attributeGroupVO, "test")
		if err != nil {
			panic(err)
		}

		for _, ins := range attributeGroupIn.Attributes {
			vo := &common.CreateAttributeVO{}
			utils.SimpleConvert(vo, ins)
			vo.AttributeGroupUUID = group.UUID
			_, err = attributeService.CreateAttribute(vo, "test")
			if err != nil {
				panic(err)
			}
		}

	}

}

func TestName(t *testing.T) {
	queryVO := common.ResourceListPageParamVO{QueryTags: map[string][]string{"parsing_records": {"@", "iauto360.cn", "1431.xyz", "carrieym.com"}, "record_value": {"10.10.11.82", "2"}}}
	marshal, _ := json.Marshal(queryVO)
	fmt.Println(string(marshal))
	where := ""
	if len(queryVO.QueryTags) > 0 {
		for k, v := range queryVO.QueryTags {
			where += "(b.uid='" + k + "' AND b.attributeInsValue in ["
			for _, value := range v {
				where += "'" + value + "',"
			}
			where = strings.TrimSuffix(where, ",") + "]) OR "
		}
		where = strings.TrimSuffix(strings.TrimSpace(where), "OR")
	}
	fmt.Println(where)
}

func TestUnmarshal(t *testing.T) {
	body := "{\"modelUid\":\"business\",\"modelName\":\"业务\",\"attributeGroupIns\":[{\"uid\":\"business_info\",\"attributeIns\":" +
		"[{\"uid\":\"business_name\",\"attributeInsValue\":\"中台-营销中心\"},{\"uid\":\"business_describe\",\"attributeInsValue\":\"为前端营销活动赋能\"}," +
		"{\"uid\":\"business_id\",\"attributeInsValue\":\"\"},{\"uid\":\"business_maintenance\",\"attributeInsValue\":\"\"}," +
		"{\"uid\":\"business_master\",\"attributeInsValue\":\"\"},{\"uid\":\"business_product\",\"attributeInsValue\":\"贾乐文\"}," +
		"{\"uid\":\"business_architect\",\"attributeInsValue\":\"\"},{\"uid\":\"business_test\",\"attributeInsValue\":\"\"}," +
		"{\"uid\":\"affiliated_center\",\"attributeInsValue\":\"\"}]}],\"uuid\":\"7ba08dd9-9832-4ae0-bc53-3d9a17fb69aa\"}\n"
	bodyObj := &store.Resource{}
	json.Unmarshal([]byte(body), bodyObj)
	//marshal, _ := json.Marshal(bodyObj)
	//fmt.Println(string(marshal))

	attributeInsValueMap := map[string]string{}
	attributeInsValueMap["business_info-business_product"] = "贾乐文"

	fmt.Println(attributeInsValueMap["business_info-business_product"])
}
