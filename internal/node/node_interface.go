package node

type Interface interface {
	// SetTask 给节点设置任务
	SetTask(spec string, f func() error)
}
