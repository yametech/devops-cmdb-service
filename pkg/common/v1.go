package common

type Business struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Leader   string    `json:"leader"`
	Children []*Domain `json:"children"`
}

type Domain struct {
	Id       int64      `json:"id"`
	Name     string     `json:"name"`
	Children []*Service `json:"children"`
}

type Service struct {
	Id       int64             `json:"id"`
	Name     string            `json:"name"`
	Desc     string            `json:"desc"`
	Owner    string            `json:"owner"`
	Children []*ServiceCluster `json:"children"`
}

type ServiceCluster struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (b *Business) AddAttribute(uid, value string) {
	switch uid {
	case "business_name":
		b.Name = value
	case "business_master":
		b.Leader = value
	}
}

func (b *Business) AddService(domainId, id int64, uid, value string) {
	var flag bool
	for _, d := range b.Children {
		if domainId == d.Id {
			for _, s := range d.Children {
				if s.Id == id {
					s.AddAttribute(uid, value)
					flag = true
					break
				}
			}
			if !flag {
				s := &Service{Id: id}
				s.AddAttribute(uid, value)
				d.Children = append(d.Children, s)
			}
		}
	}
}

func (b *Business) AddDomain(d *Domain) {
	var flag bool
	for _, child := range b.Children {
		if child.Id == d.Id {
			flag = true
			break
		}
	}
	if !flag {
		b.Children = append(b.Children, d)
	}
}

func (s *Service) AddAttribute(uid, value string) {
	switch uid {
	case "service_id":
		s.Name = value
	case "service_master":
		s.Owner = value
	case "service_describe":
		s.Desc = value
	}
}

func (s *Service) AddCluster(id int64, uid, value string) {
	var flag bool
	for _, c := range s.Children {
		if c.Id == id {
			c.AddAttribute(uid, value)
			flag = true
			break
		}
	}
	if !flag {
		sc := &ServiceCluster{Id: id}
		sc.AddAttribute(uid, value)
		s.Children = append(s.Children, sc)
	}
}

func (s *ServiceCluster) AddAttribute(uid, value string) {
	switch uid {
	case "cluster_name":
		s.Name = value
	case "cluster_describe":
		s.Desc = value
	}
}
