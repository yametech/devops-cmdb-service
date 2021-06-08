package service

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"log"
	"strings"
	"time"
)

type LdapService struct {
	Service
}

type LdapUserList struct {
	List       []common.LdapUserVO
	CreateTime int64
}

var cache *LdapUserList

func (l *LdapService) GetLdapUserList() (*common.PageResultVO, error) {
	result, err := requestLdapUserList()
	if err != nil {
		return nil, err
	}
	return &common.PageResultVO{TotalCount: int64(len(result)), List: result}, nil
}

func (l *LdapService) GetLdapUserMap() map[string]common.LdapUserVO {
	list, err := requestLdapUserList()
	if err != nil {
		return map[string]common.LdapUserVO{}
	}
	result := map[string]common.LdapUserVO{}
	for _, vo := range list {
		result[vo.Uid] = vo
	}
	return result
}

func (l *LdapService) GetUserNameByUId(uid string) string {
	list, err := requestLdapUserList()
	if err != nil {
		return ""
	}

	for _, vo := range list {
		if vo.Uid == uid {
			return vo.Name
		}
	}
	return ""
}

func requestLdapUserList() ([]common.LdapUserVO, error) {
	if cache != nil && time.Now().Unix()-cache.CreateTime < 600 {
		return cache.List, nil
	}

	log.Println("requestLdapUserList")
	l, err := ldap.DialURL(common.LdapServer)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer l.Close()

	// ldap login
	err = l.Bind(common.LdapAuthDN, common.LdapAuthPassword)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		common.LdapSearchBaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(uid=*)",            // The filter to apply
		[]string{"dn", "cn"}, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil || len(sr.Entries) <= 0 {
		log.Println(err)
		return nil, err
	}

	userList := make([]common.LdapUserVO, 0)
	for _, entry := range sr.Entries {
		user := common.LdapUserVO{}
		user.Uid = getUid(entry.DN)
		user.Name = entry.GetAttributeValue("cn")
		userList = append(userList, user)
	}
	cache = &LdapUserList{List: userList, CreateTime: time.Now().Unix()}
	return userList, nil
}

func getUid(dn string) string {
	for _, item := range strings.Split(dn, ",") {
		if strings.Split(item, "=")[0] == "uid" {
			return strings.Split(item, "=")[1]
		}
	}
	return ""
}
