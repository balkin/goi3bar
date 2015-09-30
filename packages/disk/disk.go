package disk

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"github.com/pivotal-golang/bytefmt"

	"fmt"
	"syscall"
)

const (
	freeFormat = "%v: %v free"
)

type DiskUsageItem struct {
	Name string
	Path string
}

type DiskUsageGenerator struct {
	CritThreshold int
	WarnThreshold int

	Items []DiskUsageItem
}

type diskUsageInfo struct {
	Free     uint64
	Total    uint64
	UsedPerc int
}

func getUsage(path string) (info diskUsageInfo, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return
	}

	info.Free = (uint64(stat.Bavail) * uint64(stat.Bsize))
	info.Total = (uint64(stat.Blocks) * uint64(stat.Bsize))

	info.UsedPerc = int(info.Total-info.Free) / 100

	return
}

func (g DiskUsageGenerator) Generate() ([]i3.Output, error) {
	items := make([]i3.Output, len(g.Items))

	for i, item := range g.Items {
		usage, err := getUsage(item.Path)
		if err != nil {
			return nil, err
		}

		free := bytefmt.ByteSize(usage.Free * bytefmt.BYTE)

		items[i].FullText = fmt.Sprintf(freeFormat, item.Name, free)

		freePercent := int(100 - usage.UsedPerc)
		var color string
		switch {
		case freePercent < g.CritThreshold:
			color = "#FF0000"
		case freePercent < g.WarnThreshold:
			color = "#FFA500"
		default:
			color = "#00FF00"
		}

		items[i].Color = color
	}

	items[len(items)-1].Separator = true

	return items, nil
}