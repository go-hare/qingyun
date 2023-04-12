package models

import "qingyun/common/store/mysql"

type Bill struct { //记账
	mysql.Model    `xorm:"extends"`
	UserId         int64  `xorm:"not null default 0 INT(11) comment('用户id')" json:"user_id"`
	CategoryId     int64  `xorm:"not null default 0 INT(11) comment('分类id')" json:"category_id"`
	Remark         string `xorm:"not null default '' VARCHAR(256) comment('备注')" json:"remark"`
	Price          int64  `xorm:"not null default 0 INT(11) comment('费用')" json:"price"`
	BillType       int64  `xorm:"not null default 0 INT(11) comment('类型 0 支出，1 收入')" json:"bill_type"`
	Week           int64  `xorm:"not null default 0 INT(11) comment('周 1 往后')" json:"week"`
	Year           int64  `xorm:"not null default 0 INT(11) comment('年')" json:"year"`
	Month          int64  `xorm:"not null default 0 INT(11) comment('月')" json:"month"`
	Day            int64  `xorm:"not null default 0 INT(11) comment('日')" json:"day"`
	AccountFundsId int64  `xorm:"not null default 0 INT(11) comment('账户id')" json:"account_funds_id"`
	CreateTime     int64  `xorm:"not null default 0 INT(11) comment('创建时间')" json:"create_time"`
}

func ListBills(filter mysql.OrmFilter) (list []*Bill, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	err = session.Find(&list)
	return
}

func GetBill(filter mysql.OrmFilter) (bill *Bill, err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	category := new(Bill)
	if has, err := session.Get(category); err != nil {
		return nil, err
	} else {
		if !has {
			return nil, nil
		}
		return category, nil
	}
}

func CreateBill(bill *Bill) (err error) {
	_, err = mysql.GetDB().Insert(bill)
	return
}

func UpdateBill(filter mysql.OrmFilter, message *Bill) (err error) {
	session := mysql.GetDB().NewSession()
	defer session.Close()
	if filter != nil {
		session = filter(session)
	}
	_, err = session.Update(message)
	return
}
