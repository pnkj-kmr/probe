package probe

import (
	"context"
	"log/slog"
	"probe/exporter"
	"probe/exporter/dbdump"
	M "probe/model"
	"probe/repository"
	"strconv"
	"strings"

	"probe/safemap"

	"github.com/anthdm/hollywood/actor"
)

// exporter *exporter.Exporter // M.Sender - need to make

type exportProcess struct {
	ctx      context.Context
	pType    M.ProbeType
	name     string
	engine   *actor.Engine
	pid      *actor.PID
	exporter M.Exporter[chan any, any]
	db       *repository.Repository

	bucketSize int
	partitions int
	bucket     *safemap.SafeMap[int, []any]
}

func newExportProcess(ctx *Context, name string, db *repository.Repository, opts Opts) (*exportProcess, error) {
	var _pType M.ProbeType
	if strings.Contains(name, string(M.KAFKA)) {
		_pType = M.KAFKA
	} else if strings.Contains(name, string(M.DBDUMP)) {
		_pType = M.DBDUMP
	}

	return &exportProcess{
		ctx: ctx.Context(), pType: _pType, name: name, engine: ctx.Engine(), db: db,
		bucketSize: opts.maxBucketSize,
		partitions: opts.totalPartitions,
		bucket:     safemap.New[int, []any](),
	}, nil
}

func (p *exportProcess) Send(data any) error {
	// TODO - any channel need to be handle
	p.engine.Send(p.pid, data)
	return nil
}

func (p *exportProcess) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		p.setExporter()
		slog.Info("[EXPORT] process initialized...", "name", p.name)
	case actor.Started:
		p.initBucket()
		slog.Info("[EXPORT] process started", "name", p.name)
	case actor.Stopped:
		if p.exporter != nil {
			p.exporter.Close()
		}
		slog.Info("[EXPORT] process stopped", "name", p.name)
	case M.ICMPRes:
		slog.Info("[EXPORT] data ----====>", "name", p.name, "exporter", p.exporter)
		paritionId := p.getParititionId(msg.Cid)
		p.addToBucket(paritionId, msg)
	case M.Do:
		slog.Info("[EXPORT] polling --> DO", "name", p.name, "exporter", p.exporter)
		if p.exporter == nil {
			p.setExporter()
		}
		p.initBucket()
	case M.Done:
		p.flushBucket()
		p.initBucket()
		slog.Info("[EXPORT] polling completed --> DONE")
	default:
		slog.Info("[EXPORT] default message")
		_ = msg
	}
}

func (p *exportProcess) initBucket() {
	// slog.Info("bucket init...")
	for i := 0; i < p.partitions; i++ {
		p.bucket.Set(i, []any{})
	}
}

func (p *exportProcess) addToBucket(paritionId int, data any) {
	bucket, _ := p.bucket.Get(paritionId)
	bucket = append(bucket, data)
	if len(bucket) == p.bucketSize {
		slog.Info("[EXPORT] bucket full, next interation invoked", "bucket", len(bucket))
		if p.exporter != nil {
			slog.Info("[EXPORT] exporting now ===>", "bucket", len(bucket), "parition_id", paritionId)
			if p.pType == M.KAFKA {
				p.exporter.Push() <- p.getExportMsg(paritionId, bucket)
			} else if p.pType == M.DBDUMP {
				p.exporter.Export(bucket)
			}
		}
		slog.Info("[EXPORT] reset bucket", "parition_id", paritionId)
		p.bucket.Set(paritionId, []any{})
	} else {
		slog.Info("[EXPORT] set bucket", "parition_id", paritionId)
		p.bucket.Set(paritionId, bucket)
	}
}

func (p *exportProcess) flushBucket() {
	p.bucket.ForEach(func(i int, bucket []any) {
		slog.Info("[EXPORT] +++++ FLUSHING +++++", "bucket", len(bucket), "i", i)
		if len(bucket) > 0 {
			if p.exporter != nil {
				if p.pType == M.KAFKA {
					p.exporter.Push() <- p.getExportMsg(i, bucket)
				} else if p.pType == M.DBDUMP {
					p.exporter.Export(bucket)
				}
			}
		}
	})
}

func (p *exportProcess) getParititionId(d string) int {
	i, err := strconv.Atoi(d)
	if err != nil {
		return 0
	}
	return i % p.partitions
}

func (p *exportProcess) getExportMsg(partitionId int, data any) M.ExportMsg {
	// var newData []M.ICMPRes
	// switch d := data.(type) {
	// case M.ICMPRes:
	// 	newData = append(newData, d)
	// }
	return M.ExportMsg{
		PartitionId: int32(partitionId),
		Data:        data,
		// TODO - Agent and Timezone
		Agent:    "NOID-TODO",
		Timezone: -19800,
	}
}

func (p *exportProcess) setExporter() {
	switch p.pType {
	case M.KAFKA:
		exporter, err := exporter.New(p.ctx, M.ProbeType(p.name))
		if err != nil {
			p.exporter = nil
			// TODO - need to generate the event on kafka failure
			slog.Warn("[EXPORT] kafka init error found", "err", err)
		} else {
			p.exporter = exporter
			p.exporter.Spin()
			slog.Info("[EXPORT] initialing export --", "name", p.name, "exporter", p.exporter)
		}
	case M.DBDUMP:
		exporter := dbdump.NewDBDumpExporter(p.name, p.db)
		p.exporter = exporter
		slog.Info("[EXPORT] initialing export --", "name", p.name, "exporter", p.exporter)
	}
}

func (p *exportProcess) Start() {
	p.pid = p.engine.SpawnFunc(p.Receive, p.name)
}

func (p *exportProcess) Stop() error {
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	return nil
}
