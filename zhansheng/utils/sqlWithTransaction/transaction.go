package sqlWithTransaction

import (
	"danfwing.com/m/zhansheng/config/global"
	"danfwing.com/m/zhansheng/utils/gls"
	"database/sql"
	"github.com/jinzhu/gorm"
)

func ExecSqlWithTransaction(db *sql.DB, handle func(tx *sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = handle(tx); err != nil {
		return err
	}
	return tx.Commit()
}

type D_CB func(DB **gorm.DB) (err error)

func D_Pool(cb D_CB) (err error) {
	// 获取当前 goroutine 信息
	//获取当前数据库连接
	var db *gorm.DB
	var db_pramar **gorm.DB
	db_interface := gls.Get("DB")
	if db_interface == nil {
		// 不存在则新建连接
		db = global.GDB
		// 置入 goroutine
		// 问题一、以地址传入，否则传入的将是一个 db 的克隆
		// 问题二、传入的是这个初始化的 db,而不是执行过 begin 的db,是不是需要重定义 begin？
		// 方案一：内存中存储 db 地址，无论其他地方如何使用，都改动的是地址指向的值
		// 方案二：重定义 db 的 begin 、commit 等方法，在完成操作后重新将 db 的克隆写入内存
		// 此处采用方案一
		db_pramar = &db
		gls.Put("DB", db_pramar)

	} else {
		// 获取对应 goroutine DB
		// 拿出来的 是一个 db 地址
		db_pramar = db_interface.(**gorm.DB)
	}
	return cb(db_pramar)
}

func D_Transaction(cb D_CB) (err error) {
	var db *gorm.DB
	var db_pramar **gorm.DB
	db_interface := gls.Get("DB")

	if db_interface == nil {
		db = global.GDB
		db_pramar = &db
		gls.Put("DB", db_pramar)

	} else {
		db_pramar = db_interface.(**gorm.DB)
	}

	// b = a.Begin()
	// a,b 的类型已经发生改变， a 是 sqlDb (拥有 begin 方法) ; b 是 sqlTx (拥有 commit,rollback 方法)
	// 所以即是通过 commit 提交完数据，也不能用 b.Begin 再次开启事务
	(*db_pramar) = (*db_pramar).Begin()

	//DO
	err = cb(db_pramar)

	if err != nil {
		(*db_pramar).Rollback()
		return
	}

	(*db_pramar).Commit()
	return
}
