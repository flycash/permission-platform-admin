package domain

type BusinessConfig struct {
	ID        int64  // 业务ID
	OwnerID   int64  // 业务方ID
	OwnerType string // 业务方类型
	Name      string // 业务名称
	RateLimit int    // 每秒最大请求数
	Token     string // 业务方Token，内部包含bizID也就是上方的ID，需要先插入一个空的Token获取ID，再根据ID生成token再更新
}
