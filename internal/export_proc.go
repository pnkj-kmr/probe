package probe

import (
	"context"
	"fmt"
	"log"
	"probe/exporter"
	M "probe/model"
	"strconv"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/safemap"
)

type exportProcess struct {
	ctx      context.Context
	id       int
	name     string
	engine   *actor.Engine
	pid      *actor.PID
	exporter *exporter.Exporter // M.Sender - need to make

	bucketSize int
	partition  int
	bucket     *safemap.SafeMap[int, []any]
}

func newExportProcess(ctx context.Context, id int, name string, e *actor.Engine, bucketSize int, partition int) (*exportProcess, error) {
	return &exportProcess{
		ctx: ctx,
		id:  id, name: name, engine: e,
		bucketSize: bucketSize,
		partition:  partition,
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
		log.Println("[EXPORT] process initialized...", p.name)
	case actor.Started:
		p.initBucket()
		log.Println("[EXPORT] process started", p.name)
	case actor.Stopped:
		if p.exporter != nil {
			p.exporter.Close()
		}
		log.Println("[EXPORT] process stopped", p.name)
	case M.ICMPRes:
		// log.Println("========> exportProcess message received========>")
		// TODO - need to group x message gourp
		// send into bulk
		// msg based control
		// p.exporter.Send() <- d
		paritionId := p.getParititionId(msg.Cid)
		p.addToBucket(paritionId, msg)
	case M.Do:
		fmt.Println("polling started here by poller ----------------> DO ")
		if p.exporter == nil {
			p.setExporter()
		}
		p.initBucket()
	case M.Done:
		p.flushBucket()
		p.initBucket()
		fmt.Println("polling completed here by poller --------------> DONE")
	default:
		log.Println("[EXPORT] default message")
		_ = msg
	}
}

func (p *exportProcess) initBucket() {
	log.Println("bucket init =====>")
	for i := 0; i < p.partition; i++ {
		p.bucket.Set(i, []any{})
	}
}

func (p *exportProcess) addToBucket(paritionId int, data any) {
	bucket, _ := p.bucket.Get(paritionId)
	bucket = append(bucket, data)
	if len(bucket) == p.bucketSize {
		fmt.Println("++++++++++++ sending bucket to next stage ---->", len(bucket))
		if p.exporter != nil {
			log.Println("====================sendin to exporter", p.getExportMsg(paritionId, bucket))
			p.exporter.Send() <- p.getExportMsg(paritionId, bucket)
		}
		log.Println(paritionId, "reset --->")
		p.bucket.Set(paritionId, []any{})
	} else {
		log.Println(paritionId, "set -->")
		p.bucket.Set(paritionId, bucket)
	}
}

func (p *exportProcess) flushBucket() {
	log.Println("bucket flush =====>")
	p.bucket.ForEach(func(i int, bucket []any) {
		fmt.Println("FLUSHING ++++++++++++ sending bucket to next stage ---->", i, "bucket", len(bucket))
		if len(bucket) > 0 {
			if p.exporter != nil {
				p.exporter.Send() <- p.getExportMsg(i, bucket)
			}
		}
	})
}

func (p *exportProcess) getParititionId(d string) int {
	i, err := strconv.Atoi(d)
	if err != nil {
		return 0
	}
	return i % p.partition
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
	switch p.id {
	case KAFKA:
		exporter, err := exporter.New(p.ctx, p.id, p.name)
		if err != nil {
			p.exporter = nil
			// TODO - need to generate the event on kafka failure
			log.Println("[EXPORT] kafka init error found ===> ", err)
		} else {
			p.exporter = exporter
			p.exporter.Spin()
		}
	}
}

func (p *exportProcess) Start() {
	// log.Println("[EXPORT] process starting...")
	p.pid = p.engine.SpawnFunc(p.Receive, p.name, actor.WithID(strconv.Itoa(p.id)))

}

func (p *exportProcess) Stop() error {
	ctx := p.engine.Poison(p.pid)
	<-ctx.Done()
	// log.Println("[EXPORT] process stopped")
	return nil
}
