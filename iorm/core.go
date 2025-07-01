package iorm

import (
	"fmt"
	"github.com/dromara/carbon/v2"
	"gorm.io/gorm"
	"strings"
)

type IOrm struct {
	db *gorm.DB
}

func NewIOrm(db *gorm.DB) *IOrm {
	return &IOrm{db: db}
}

// BuildTableIndex 构建索引
func (t *IOrm) BuildTableIndex(data []TableIndex) error {
	for _, arr := range data {
		for _, tableName := range arr.TableName {
			for _, index := range arr.IndexArr {
				//检查表是否存在
				if !t.db.Migrator().HasTable(tableName) {
					return fmt.Errorf("table %s not found", tableName)
				}

				var count int64
				sql := fmt.Sprintf(`SELECT COUNT(*) AS total FROM information_schema.statistics WHERE table_name = '%s' AND index_name = '%s'`, tableName, index.IndexName)
				if err := t.db.Raw(sql).Scan(&count).Error; err != nil {
					return err
				}
				if count > 0 {
					continue
				}

				//构建索引
				if index.Unique {
					sql = fmt.Sprintf("ALTER TABLE %s ADD UNIQUE INDEX %s (%s)", tableName, index.IndexName, strings.Join(index.Columns, ","))
				} else {
					sql = fmt.Sprintf("ALTER TABLE %s ADD INDEX %s (%s)", tableName, index.IndexName, strings.Join(index.Columns, ","))
				}
				if err := t.db.Exec(sql).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// BuildTable 构建表
func (t *IOrm) BuildTable(data []TableModel) error {
	for _, v := range data {
		for _, tableName := range v.TableName {
			if v.PartitionType == PartitionRANGE {
				opt := v.PartitionModel.(PartitionModelRange)
				for i := 0; i <= opt.AddDays; i++ {
					buildTime := carbon.ParseByFormat(opt.StartDay, "Ymd").AddDays(i).Format("Ymd")
					if err := t.BuildRangePartitionByDay(tableName, opt.KeyField, buildTime); err != nil {
						return err
					}
				}
			} else if v.PartitionType == PartitionKEY {
				opt := v.PartitionModel.(PartitionModelKey)
				if err := t.BuildKeyPartition(tableName, opt.KeyField, opt.KeyTotal); err != nil {
					return err
				}
			} else if v.PartitionType == PartitionLIST {
				opt := v.PartitionModel.(PartitionModelList)
				if err := t.BuildListPartition(tableName, opt.KeyField, opt.ListArr); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// BuildKeyPartition 构建key分区
// @table 分区表
// @keyField 分区字段，必须是主键
// @keyTotal 分区数量
func (t *IOrm) BuildKeyPartition(table, keyField string, keyTotal int) error {
	partitionInfo := make([]InformationSchema, 0)
	err := t.db.Table("information_schema.partitions").Where("table_name = ? and partition_name != ''", table).Find(&partitionInfo).Error
	if err != nil {
		return err
	}

	if len(partitionInfo) > 0 {
		return nil
	}

	sql := fmt.Sprintf("ALTER TABLE %s PARTITION BY KEY(%s) PARTITIONS %d", table, keyField, keyTotal)
	if err = t.db.Exec(sql).Error; err != nil {
		return err
	}

	return nil
}

// BuildRangePartitionByDay 按天构建range分区
// @table 分区表
// @keyField 分区字段，必须是主键
// @buildTime 20250601,允许为空
func (t *IOrm) BuildRangePartitionByDay(table string, keyField string, buildTime string) error {
	todayPartition := "p" + carbon.Now().Format("Ymd")
	todayLess := carbon.Tomorrow().Format("Ymd")

	if buildTime != "" {
		todayPartition = "p" + buildTime
		todayLess = carbon.ParseByFormat(buildTime, "Ymd").AddDays(1).Format("Ymd")
	}

	if !t.db.Migrator().HasTable(table) {
		return fmt.Errorf("table %s not found", table)
	}

	partitionInfo := make([]InformationSchema, 0)
	err := t.db.Table("information_schema.partitions").Where("table_name = ? and partition_name != ''", table).Find(&partitionInfo).Error
	if err != nil {
		return err
	}

	//没有分区则初始化分区
	if len(partitionInfo) == 0 {
		sql := fmt.Sprintf("ALTER TABLE %s PARTITION BY RANGE (TO_DAYS(%s)) (PARTITION %s VALUES LESS THAN (TO_DAYS('%s')))", table, keyField, todayPartition, todayLess)
		if err = t.db.Exec(sql).Error; err != nil {
			return err
		}
	}

	//如果存在分区，则检查是否包含今日分区，没有则添加
	var exits bool
	for _, v := range partitionInfo {
		if v.PARTITION_NAME == todayPartition {
			exits = true
			break
		}
	}
	if !exits {
		sql := fmt.Sprintf("ALTER TABLE %s ADD PARTITION (PARTITION %s VALUES LESS THAN (TO_DAYS('%s')))", table, todayPartition, todayLess)
		if err = t.db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

// DropRangePartition 删除range分区
// @table 分区表
// @buildTime 待删除分区时间 20250601
func (t *IOrm) DropRangePartition(table string, buildTime string) error {
	if buildTime == "" {
		return fmt.Errorf("missing build time")
	}

	partitionName := "p" + buildTime

	var partitionInfo InformationSchema
	err := t.db.Table("information_schema.partitions").Where("table_name = ? and partition_name = ?", table, partitionName).Find(&partitionInfo).Error
	if err != nil {
		return err
	}

	if partitionInfo.PARTITION_NAME != "" {
		sql := fmt.Sprintf("ALTER TABLE %s DROP PARTITION %s", table, partitionName)
		if err = t.db.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// BuildListPartition 构建list分区
// @table 分区表
// @keyField 分区字段，必须是主键
// @listArr 分区内容
// 如果新增分区内容，可以动态拆分othen并增加list分区
// ALTER TABLE 表名 REORGANIZE PARTITION `p_other` INTO (PARTITION p_base VALUES IN ('base'),PARTITION p_other VALUES IN (DEFAULT))
func (t *IOrm) BuildListPartition(table, keyField string, listArr []string) error {
	partitionInfo := make([]InformationSchema, 0)
	err := t.db.Table("information_schema.partitions").Where("table_name = ? and partition_name != ''", table).Find(&partitionInfo).Error
	if err != nil {
		return err
	}

	//无分区，则初始化
	//有分区，则判断是否有新增差异
	if len(partitionInfo) == 0 {
		sqlItem := make([]string, 0)
		for _, item := range listArr {
			sqlItem = append(sqlItem, fmt.Sprintf("PARTITION p_%s VALUES IN ('%s')", item, item))
		}

		sql := fmt.Sprintf("ALTER TABLE %s PARTITION BY LIST COLUMNS(%s) (%s,PARTITION p_other VALUES IN (DEFAULT))", table, keyField, strings.Join(sqlItem, ","))

		if err = t.db.Exec(sql).Error; err != nil {
			return err
		}
	} else {
		tmpArr := make(map[string]int)
		for _, v := range partitionInfo {
			tmpArr[v.PARTITION_NAME] = 1
		}

		newListArr := make([]string, 0)
		for _, v := range listArr {
			_, ok := tmpArr["p_"+v]
			if !ok {
				newListArr = append(newListArr, v)
			}
		}

		if len(newListArr) == 0 {
			return fmt.Errorf("no new list")
		}

		sqlItem := make([]string, 0)
		for _, item := range newListArr {
			sqlItem = append(sqlItem, fmt.Sprintf("PARTITION p_%s VALUES IN ('%s')", item, item))
		}
		sql := fmt.Sprintf("ALTER TABLE %s REORGANIZE PARTITION `p_other` INTO (%s,PARTITION p_other VALUES IN (DEFAULT))", table, strings.Join(sqlItem, ","))
		if err = t.db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}
