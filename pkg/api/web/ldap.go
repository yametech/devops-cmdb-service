package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"strings"
)

type LdapApi struct {
}

func (r *LdapApi) router(e *gin.Engine) {
	groupRoute := e.Group(common.WEB_API_GROUP)
	groupRoute.GET("/ldap/user-list", r.getLdapUserList)
}

func (r *LdapApi) getLdapUserList(ctx *gin.Context) {
	pageList := common.PageResultVO{}
	list, err := LdapUserList()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	pageList.TotalCount = int64(len(list))
	userList := make([]interface{}, 0)
	for _, entry := range list {
		user := common.LdapUserVO{}
		user.Uid = r.getUid(entry.DN)
		user.Name = entry.GetAttributeValue("cn")
		userList = append(userList, user)
	}
	pageList.List = userList
	Success(ctx, pageList)
}

func (r *LdapApi) getUid(dn string) string {
	for _, item := range strings.Split(dn, ",") {
		if strings.Split(item, "=")[0] == "uid" {
			return strings.Split(item, "=")[1]
		}
	}
	return ""
}

func LdapUserList() ([]*ldap.Entry, error) {
	l, err := ldap.DialURL(common.LdapServer)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	// ldap login
	err = l.Bind(common.LdapAuthDN, common.LdapAuthPassword)
	if err != nil {
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
		return nil, err
	}

	return sr.Entries, nil
}
