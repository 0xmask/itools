package iorm

type PartitionType int //分区类型

const (
	PartitionNORMAL PartitionType = iota //无分区
	PartitionKEY                         //key分区
	PartitionRANGE                       //range分区
	PartitionLIST                        //list分区
)

type InformationSchema struct {
	TABLE_NAME     string `json:"TABLE_NAME"`
	PARTITION_NAME string `json:"PARTITION_NAME"`
}

type IndexModel struct {
	IndexName string   //索引名
	Unique    bool     //是否唯一索引
	Columns   []string //索引列
}

type TableIndex struct {
	TableName []string     //表名
	IndexArr  []IndexModel //索引
}

type TableModel struct {
	TableName      []string      //表名
	Model          any           //模型
	PartitionType  PartitionType //分区类型
	PartitionModel any           //分区条件
}

// PartitionModelKey key分区模型
type PartitionModelKey struct {
	KeyField string //分区字段
	KeyTotal int    //分区数量
}

// PartitionModelRange range按天分区模型
type PartitionModelRange struct {
	KeyField string //分区字段
	StartDay string //开始时间，格式20250601
	AddDays  int    //构建次数
}

// PartitionModelList list分区模型
type PartitionModelList struct {
	KeyField string   //分区字段
	ListArr  []string //list分区数据
}
