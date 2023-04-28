package generator

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

/**
建表：
	CREATE TABLE `id_generator_tab` (
	 `business` varchar(32) NOT NULL,
	 `current_id` bigint NOT NULL,
	 `step` bigint NOT NULL,
	 `desc` varchar(1024) DEFAULT '' NOT NULL,
	 `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	 PRIMARY KEY (`biz_tag`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

type Mysql struct {
	db *sql.DB
}

var GMysql *Mysql

const (
	NewSegmentRetryTimes = 5
)

func InitMysql() error {
	db, err := sql.Open("mysql", GConf.DSN)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0)
	GMysql = &Mysql{db: db}
	return nil
}

func (mysql *Mysql) NextSegment(business string) (*Segment, error) {
	var (
		currentId    int64
		step         int64
		result       sql.Result
		rowsAffected int64
		err          error
	)

	// 总耗时小于2秒
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(20000)*time.Millisecond)
	defer cancelFunc()

	for counter := 0; counter < NewSegmentRetryTimes; counter++ {
		query := "SELECT current_id,step FROM " + GConf.Table + " WHERE business=?"
		if err := mysql.db.QueryRowContext(ctx, query, business).Scan(&currentId, &step); err != nil {
			return nil, err
		}

		update := "UPDATE " + GConf.Table + " SET current_id=current_id+step WHERE business=? and current_id =?"
		if result, err = mysql.db.ExecContext(ctx, update, business, currentId); err != nil {
			return nil, err
		}

		if rowsAffected, err = result.RowsAffected(); err != nil { // 失败
			return nil, err
		} else if rowsAffected == 0 {
			// 记录不存在，继续执行
			continue
		}

		// 执行成功
		segment := &Segment{
			CurrentId: currentId + step,
			Offset:    0,
			Step:      step,
		}

		return segment, nil
	}

	return nil, errors.New("new segment error")
}
