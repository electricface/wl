package wl

import (
	"sync"
)

type DisplayErrorEvent struct {
	ObjectId Proxy
	Code uint32
	Message string
}

func (p *Display) AddErrorHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.errorHandlers = append(p.errorHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Display) RemoveErrorHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.errorHandlers {
		if e == h {
			p.errorHandlers = append(p.errorHandlers[:i] , p.errorHandlers[i+1:]...)
			break
		}
	}
}

type DisplayDeleteIdEvent struct {
	Id uint32
}

func (p *Display) AddDeleteIdHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.deleteIdHandlers = append(p.deleteIdHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Display) RemoveDeleteIdHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.deleteIdHandlers {
		if e == h {
			p.deleteIdHandlers = append(p.deleteIdHandlers[:i] , p.deleteIdHandlers[i+1:]...)
			break
		}
	}
}

func (p *Display) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.errorHandlers) > 0 {
			ev := DisplayErrorEvent{}
			ev.ObjectId = event.Proxy(p.Context())
			ev.Code = event.Uint32()
			ev.Message = event.String()
			p.mu.RLock()
			for _, h := range p.errorHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.deleteIdHandlers) > 0 {
			ev := DisplayDeleteIdEvent{}
			ev.Id = event.Uint32()
			p.mu.RLock()
			for _, h := range p.deleteIdHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Display struct {
	BaseProxy
	mu sync.RWMutex
	errorHandlers []Handler
	deleteIdHandlers []Handler
}

func NewDisplay(ctx *Context) *Display {
	ret := new(Display)
	ctx.register(ret)
	return ret
}

func (p *Display) Sync() (*Callback , error) {
	ret := NewCallback(p.Context())
	return ret , p.Context().sendRequest(p,0,Proxy(ret))
}

func (p *Display) GetRegistry() (*Registry , error) {
	ret := NewRegistry(p.Context())
	return ret , p.Context().sendRequest(p,1,Proxy(ret))
}

const (
	DisplayErrorInvalidObject = 0
	DisplayErrorInvalidMethod = 1
	DisplayErrorNoMemory = 2
)

type RegistryGlobalEvent struct {
	Name uint32
	Interface string
	Version uint32
}

func (p *Registry) AddGlobalHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.globalHandlers = append(p.globalHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Registry) RemoveGlobalHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.globalHandlers {
		if e == h {
			p.globalHandlers = append(p.globalHandlers[:i] , p.globalHandlers[i+1:]...)
			break
		}
	}
}

type RegistryGlobalRemoveEvent struct {
	Name uint32
}

func (p *Registry) AddGlobalRemoveHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.globalRemoveHandlers = append(p.globalRemoveHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Registry) RemoveGlobalRemoveHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.globalRemoveHandlers {
		if e == h {
			p.globalRemoveHandlers = append(p.globalRemoveHandlers[:i] , p.globalRemoveHandlers[i+1:]...)
			break
		}
	}
}

func (p *Registry) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.globalHandlers) > 0 {
			ev := RegistryGlobalEvent{}
			ev.Name = event.Uint32()
			ev.Interface = event.String()
			ev.Version = event.Uint32()
			p.mu.RLock()
			for _, h := range p.globalHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.globalRemoveHandlers) > 0 {
			ev := RegistryGlobalRemoveEvent{}
			ev.Name = event.Uint32()
			p.mu.RLock()
			for _, h := range p.globalRemoveHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Registry struct {
	BaseProxy
	mu sync.RWMutex
	globalHandlers []Handler
	globalRemoveHandlers []Handler
}

func NewRegistry(ctx *Context) *Registry {
	ret := new(Registry)
	ctx.register(ret)
	return ret
}

func (p *Registry) Bind(name uint32,iface string,version uint32,id Proxy) error {
	return p.Context().sendRequest(p,0,name,iface,version,id)
}

type CallbackDoneEvent struct {
	CallbackData uint32
}

func (p *Callback) AddDoneHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.doneHandlers = append(p.doneHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Callback) RemoveDoneHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.doneHandlers {
		if e == h {
			p.doneHandlers = append(p.doneHandlers[:i] , p.doneHandlers[i+1:]...)
			break
		}
	}
}

func (p *Callback) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.doneHandlers) > 0 {
			ev := CallbackDoneEvent{}
			ev.CallbackData = event.Uint32()
			p.mu.RLock()
			for _, h := range p.doneHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Callback struct {
	BaseProxy
	mu sync.RWMutex
	doneHandlers []Handler
}

func NewCallback(ctx *Context) *Callback {
	ret := new(Callback)
	ctx.register(ret)
	return ret
}

type Compositor struct {
	BaseProxy
}

func NewCompositor(ctx *Context) *Compositor {
	ret := new(Compositor)
	ctx.register(ret)
	return ret
}

func (p *Compositor) CreateSurface() (*Surface , error) {
	ret := NewSurface(p.Context())
	return ret , p.Context().sendRequest(p,0,Proxy(ret))
}

func (p *Compositor) CreateRegion() (*Region , error) {
	ret := NewRegion(p.Context())
	return ret , p.Context().sendRequest(p,1,Proxy(ret))
}

type ShmPool struct {
	BaseProxy
}

func NewShmPool(ctx *Context) *ShmPool {
	ret := new(ShmPool)
	ctx.register(ret)
	return ret
}

func (p *ShmPool) CreateBuffer(offset int32,width int32,height int32,stride int32,format uint32) (*Buffer , error) {
	ret := NewBuffer(p.Context())
	return ret , p.Context().sendRequest(p,0,Proxy(ret),offset,width,height,stride,format)
}

func (p *ShmPool) Destroy() error {
	return p.Context().sendRequest(p,1)
}

func (p *ShmPool) Resize(size int32) error {
	return p.Context().sendRequest(p,2,size)
}

type ShmFormatEvent struct {
	Format uint32
}

func (p *Shm) AddFormatHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.formatHandlers = append(p.formatHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Shm) RemoveFormatHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.formatHandlers {
		if e == h {
			p.formatHandlers = append(p.formatHandlers[:i] , p.formatHandlers[i+1:]...)
			break
		}
	}
}

func (p *Shm) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.formatHandlers) > 0 {
			ev := ShmFormatEvent{}
			ev.Format = event.Uint32()
			p.mu.RLock()
			for _, h := range p.formatHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Shm struct {
	BaseProxy
	mu sync.RWMutex
	formatHandlers []Handler
}

func NewShm(ctx *Context) *Shm {
	ret := new(Shm)
	ctx.register(ret)
	return ret
}

func (p *Shm) CreatePool(fd uintptr,size int32) (*ShmPool , error) {
	ret := NewShmPool(p.Context())
	return ret , p.Context().sendRequest(p,0,Proxy(ret),fd,size)
}

const (
	ShmErrorInvalidFormat = 0
	ShmErrorInvalidStride = 1
	ShmErrorInvalidFd = 2
)

const (
	ShmFormatArgb8888 = 0
	ShmFormatXrgb8888 = 1
	ShmFormatC8 = 0x20203843
	ShmFormatRgb332 = 0x38424752
	ShmFormatBgr233 = 0x38524742
	ShmFormatXrgb4444 = 0x32315258
	ShmFormatXbgr4444 = 0x32314258
	ShmFormatRgbx4444 = 0x32315852
	ShmFormatBgrx4444 = 0x32315842
	ShmFormatArgb4444 = 0x32315241
	ShmFormatAbgr4444 = 0x32314241
	ShmFormatRgba4444 = 0x32314152
	ShmFormatBgra4444 = 0x32314142
	ShmFormatXrgb1555 = 0x35315258
	ShmFormatXbgr1555 = 0x35314258
	ShmFormatRgbx5551 = 0x35315852
	ShmFormatBgrx5551 = 0x35315842
	ShmFormatArgb1555 = 0x35315241
	ShmFormatAbgr1555 = 0x35314241
	ShmFormatRgba5551 = 0x35314152
	ShmFormatBgra5551 = 0x35314142
	ShmFormatRgb565 = 0x36314752
	ShmFormatBgr565 = 0x36314742
	ShmFormatRgb888 = 0x34324752
	ShmFormatBgr888 = 0x34324742
	ShmFormatXbgr8888 = 0x34324258
	ShmFormatRgbx8888 = 0x34325852
	ShmFormatBgrx8888 = 0x34325842
	ShmFormatAbgr8888 = 0x34324241
	ShmFormatRgba8888 = 0x34324152
	ShmFormatBgra8888 = 0x34324142
	ShmFormatXrgb2101010 = 0x30335258
	ShmFormatXbgr2101010 = 0x30334258
	ShmFormatRgbx1010102 = 0x30335852
	ShmFormatBgrx1010102 = 0x30335842
	ShmFormatArgb2101010 = 0x30335241
	ShmFormatAbgr2101010 = 0x30334241
	ShmFormatRgba1010102 = 0x30334152
	ShmFormatBgra1010102 = 0x30334142
	ShmFormatYuyv = 0x56595559
	ShmFormatYvyu = 0x55595659
	ShmFormatUyvy = 0x59565955
	ShmFormatVyuy = 0x59555956
	ShmFormatAyuv = 0x56555941
	ShmFormatNv12 = 0x3231564e
	ShmFormatNv21 = 0x3132564e
	ShmFormatNv16 = 0x3631564e
	ShmFormatNv61 = 0x3136564e
	ShmFormatYuv410 = 0x39565559
	ShmFormatYvu410 = 0x39555659
	ShmFormatYuv411 = 0x31315559
	ShmFormatYvu411 = 0x31315659
	ShmFormatYuv420 = 0x32315559
	ShmFormatYvu420 = 0x32315659
	ShmFormatYuv422 = 0x36315559
	ShmFormatYvu422 = 0x36315659
	ShmFormatYuv444 = 0x34325559
	ShmFormatYvu444 = 0x34325659
)

type BufferReleaseEvent struct {
}

func (p *Buffer) AddReleaseHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.releaseHandlers = append(p.releaseHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Buffer) RemoveReleaseHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.releaseHandlers {
		if e == h {
			p.releaseHandlers = append(p.releaseHandlers[:i] , p.releaseHandlers[i+1:]...)
			break
		}
	}
}

func (p *Buffer) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.releaseHandlers) > 0 {
			ev := BufferReleaseEvent{}
			p.mu.RLock()
			for _, h := range p.releaseHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Buffer struct {
	BaseProxy
	mu sync.RWMutex
	releaseHandlers []Handler
}

func NewBuffer(ctx *Context) *Buffer {
	ret := new(Buffer)
	ctx.register(ret)
	return ret
}

func (p *Buffer) Destroy() error {
	return p.Context().sendRequest(p,0)
}

type DataOfferOfferEvent struct {
	MimeType string
}

func (p *DataOffer) AddOfferHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.offerHandlers = append(p.offerHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataOffer) RemoveOfferHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.offerHandlers {
		if e == h {
			p.offerHandlers = append(p.offerHandlers[:i] , p.offerHandlers[i+1:]...)
			break
		}
	}
}

type DataOfferSourceActionsEvent struct {
	SourceActions uint32
}

func (p *DataOffer) AddSourceActionsHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.sourceActionsHandlers = append(p.sourceActionsHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataOffer) RemoveSourceActionsHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.sourceActionsHandlers {
		if e == h {
			p.sourceActionsHandlers = append(p.sourceActionsHandlers[:i] , p.sourceActionsHandlers[i+1:]...)
			break
		}
	}
}

type DataOfferActionEvent struct {
	DndAction uint32
}

func (p *DataOffer) AddActionHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.actionHandlers = append(p.actionHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataOffer) RemoveActionHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.actionHandlers {
		if e == h {
			p.actionHandlers = append(p.actionHandlers[:i] , p.actionHandlers[i+1:]...)
			break
		}
	}
}

func (p *DataOffer) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.offerHandlers) > 0 {
			ev := DataOfferOfferEvent{}
			ev.MimeType = event.String()
			p.mu.RLock()
			for _, h := range p.offerHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.sourceActionsHandlers) > 0 {
			ev := DataOfferSourceActionsEvent{}
			ev.SourceActions = event.Uint32()
			p.mu.RLock()
			for _, h := range p.sourceActionsHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 2:
		if len(p.actionHandlers) > 0 {
			ev := DataOfferActionEvent{}
			ev.DndAction = event.Uint32()
			p.mu.RLock()
			for _, h := range p.actionHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type DataOffer struct {
	BaseProxy
	mu sync.RWMutex
	offerHandlers []Handler
	sourceActionsHandlers []Handler
	actionHandlers []Handler
}

func NewDataOffer(ctx *Context) *DataOffer {
	ret := new(DataOffer)
	ctx.register(ret)
	return ret
}

func (p *DataOffer) Accept(serial uint32,mime_type string) error {
	return p.Context().sendRequest(p,0,serial,mime_type)
}

func (p *DataOffer) Receive(mime_type string,fd uintptr) error {
	return p.Context().sendRequest(p,1,mime_type,fd)
}

func (p *DataOffer) Destroy() error {
	return p.Context().sendRequest(p,2)
}

func (p *DataOffer) Finish() error {
	return p.Context().sendRequest(p,3)
}

func (p *DataOffer) SetActions(dnd_actions uint32,preferred_action uint32) error {
	return p.Context().sendRequest(p,4,dnd_actions,preferred_action)
}

const (
	DataOfferErrorInvalidFinish = 0
	DataOfferErrorInvalidActionMask = 1
	DataOfferErrorInvalidAction = 2
	DataOfferErrorInvalidOffer = 3
)

type DataSourceTargetEvent struct {
	MimeType string
}

func (p *DataSource) AddTargetHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.targetHandlers = append(p.targetHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataSource) RemoveTargetHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.targetHandlers {
		if e == h {
			p.targetHandlers = append(p.targetHandlers[:i] , p.targetHandlers[i+1:]...)
			break
		}
	}
}

type DataSourceSendEvent struct {
	MimeType string
	Fd uintptr
}

func (p *DataSource) AddSendHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.sendHandlers = append(p.sendHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataSource) RemoveSendHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.sendHandlers {
		if e == h {
			p.sendHandlers = append(p.sendHandlers[:i] , p.sendHandlers[i+1:]...)
			break
		}
	}
}

type DataSourceCancelledEvent struct {
}

func (p *DataSource) AddCancelledHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.cancelledHandlers = append(p.cancelledHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataSource) RemoveCancelledHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.cancelledHandlers {
		if e == h {
			p.cancelledHandlers = append(p.cancelledHandlers[:i] , p.cancelledHandlers[i+1:]...)
			break
		}
	}
}

type DataSourceDndDropPerformedEvent struct {
}

func (p *DataSource) AddDndDropPerformedHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.dndDropPerformedHandlers = append(p.dndDropPerformedHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataSource) RemoveDndDropPerformedHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.dndDropPerformedHandlers {
		if e == h {
			p.dndDropPerformedHandlers = append(p.dndDropPerformedHandlers[:i] , p.dndDropPerformedHandlers[i+1:]...)
			break
		}
	}
}

type DataSourceDndFinishedEvent struct {
}

func (p *DataSource) AddDndFinishedHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.dndFinishedHandlers = append(p.dndFinishedHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataSource) RemoveDndFinishedHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.dndFinishedHandlers {
		if e == h {
			p.dndFinishedHandlers = append(p.dndFinishedHandlers[:i] , p.dndFinishedHandlers[i+1:]...)
			break
		}
	}
}

type DataSourceActionEvent struct {
	DndAction uint32
}

func (p *DataSource) AddActionHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.actionHandlers = append(p.actionHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataSource) RemoveActionHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.actionHandlers {
		if e == h {
			p.actionHandlers = append(p.actionHandlers[:i] , p.actionHandlers[i+1:]...)
			break
		}
	}
}

func (p *DataSource) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.targetHandlers) > 0 {
			ev := DataSourceTargetEvent{}
			ev.MimeType = event.String()
			p.mu.RLock()
			for _, h := range p.targetHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.sendHandlers) > 0 {
			ev := DataSourceSendEvent{}
			ev.MimeType = event.String()
			ev.Fd = event.FD()
			p.mu.RLock()
			for _, h := range p.sendHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 2:
		if len(p.cancelledHandlers) > 0 {
			ev := DataSourceCancelledEvent{}
			p.mu.RLock()
			for _, h := range p.cancelledHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 3:
		if len(p.dndDropPerformedHandlers) > 0 {
			ev := DataSourceDndDropPerformedEvent{}
			p.mu.RLock()
			for _, h := range p.dndDropPerformedHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 4:
		if len(p.dndFinishedHandlers) > 0 {
			ev := DataSourceDndFinishedEvent{}
			p.mu.RLock()
			for _, h := range p.dndFinishedHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 5:
		if len(p.actionHandlers) > 0 {
			ev := DataSourceActionEvent{}
			ev.DndAction = event.Uint32()
			p.mu.RLock()
			for _, h := range p.actionHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type DataSource struct {
	BaseProxy
	mu sync.RWMutex
	targetHandlers []Handler
	sendHandlers []Handler
	cancelledHandlers []Handler
	dndDropPerformedHandlers []Handler
	dndFinishedHandlers []Handler
	actionHandlers []Handler
}

func NewDataSource(ctx *Context) *DataSource {
	ret := new(DataSource)
	ctx.register(ret)
	return ret
}

func (p *DataSource) Offer(mime_type string) error {
	return p.Context().sendRequest(p,0,mime_type)
}

func (p *DataSource) Destroy() error {
	return p.Context().sendRequest(p,1)
}

func (p *DataSource) SetActions(dnd_actions uint32) error {
	return p.Context().sendRequest(p,2,dnd_actions)
}

const (
	DataSourceErrorInvalidActionMask = 0
	DataSourceErrorInvalidSource = 1
)

type DataDeviceDataOfferEvent struct {
	Id *DataOffer
}

func (p *DataDevice) AddDataOfferHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.dataOfferHandlers = append(p.dataOfferHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataDevice) RemoveDataOfferHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.dataOfferHandlers {
		if e == h {
			p.dataOfferHandlers = append(p.dataOfferHandlers[:i] , p.dataOfferHandlers[i+1:]...)
			break
		}
	}
}

type DataDeviceEnterEvent struct {
	Serial uint32
	Surface *Surface
	X float32
	Y float32
	Id *DataOffer
}

func (p *DataDevice) AddEnterHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.enterHandlers = append(p.enterHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataDevice) RemoveEnterHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.enterHandlers {
		if e == h {
			p.enterHandlers = append(p.enterHandlers[:i] , p.enterHandlers[i+1:]...)
			break
		}
	}
}

type DataDeviceLeaveEvent struct {
}

func (p *DataDevice) AddLeaveHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.leaveHandlers = append(p.leaveHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataDevice) RemoveLeaveHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.leaveHandlers {
		if e == h {
			p.leaveHandlers = append(p.leaveHandlers[:i] , p.leaveHandlers[i+1:]...)
			break
		}
	}
}

type DataDeviceMotionEvent struct {
	Time uint32
	X float32
	Y float32
}

func (p *DataDevice) AddMotionHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.motionHandlers = append(p.motionHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataDevice) RemoveMotionHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.motionHandlers {
		if e == h {
			p.motionHandlers = append(p.motionHandlers[:i] , p.motionHandlers[i+1:]...)
			break
		}
	}
}

type DataDeviceDropEvent struct {
}

func (p *DataDevice) AddDropHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.dropHandlers = append(p.dropHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataDevice) RemoveDropHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.dropHandlers {
		if e == h {
			p.dropHandlers = append(p.dropHandlers[:i] , p.dropHandlers[i+1:]...)
			break
		}
	}
}

type DataDeviceSelectionEvent struct {
	Id *DataOffer
}

func (p *DataDevice) AddSelectionHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.selectionHandlers = append(p.selectionHandlers , h)
		p.mu.Unlock()
	}
}

func (p *DataDevice) RemoveSelectionHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.selectionHandlers {
		if e == h {
			p.selectionHandlers = append(p.selectionHandlers[:i] , p.selectionHandlers[i+1:]...)
			break
		}
	}
}

func (p *DataDevice) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.dataOfferHandlers) > 0 {
			ev := DataDeviceDataOfferEvent{}
			ev.Id = event.Proxy(p.Context()).(*DataOffer)
			p.mu.RLock()
			for _, h := range p.dataOfferHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.enterHandlers) > 0 {
			ev := DataDeviceEnterEvent{}
			ev.Serial = event.Uint32()
			ev.Surface = event.Proxy(p.Context()).(*Surface)
			ev.X = event.Float32()
			ev.Y = event.Float32()
			ev.Id = event.Proxy(p.Context()).(*DataOffer)
			p.mu.RLock()
			for _, h := range p.enterHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 2:
		if len(p.leaveHandlers) > 0 {
			ev := DataDeviceLeaveEvent{}
			p.mu.RLock()
			for _, h := range p.leaveHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 3:
		if len(p.motionHandlers) > 0 {
			ev := DataDeviceMotionEvent{}
			ev.Time = event.Uint32()
			ev.X = event.Float32()
			ev.Y = event.Float32()
			p.mu.RLock()
			for _, h := range p.motionHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 4:
		if len(p.dropHandlers) > 0 {
			ev := DataDeviceDropEvent{}
			p.mu.RLock()
			for _, h := range p.dropHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 5:
		if len(p.selectionHandlers) > 0 {
			ev := DataDeviceSelectionEvent{}
			ev.Id = event.Proxy(p.Context()).(*DataOffer)
			p.mu.RLock()
			for _, h := range p.selectionHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type DataDevice struct {
	BaseProxy
	mu sync.RWMutex
	dataOfferHandlers []Handler
	enterHandlers []Handler
	leaveHandlers []Handler
	motionHandlers []Handler
	dropHandlers []Handler
	selectionHandlers []Handler
}

func NewDataDevice(ctx *Context) *DataDevice {
	ret := new(DataDevice)
	ctx.register(ret)
	return ret
}

func (p *DataDevice) StartDrag(source *DataSource,origin *Surface,icon *Surface,serial uint32) error {
	return p.Context().sendRequest(p,0,source,origin,icon,serial)
}

func (p *DataDevice) SetSelection(source *DataSource,serial uint32) error {
	return p.Context().sendRequest(p,1,source,serial)
}

func (p *DataDevice) Release() error {
	return p.Context().sendRequest(p,2)
}

const (
	DataDeviceErrorRole = 0
)

type DataDeviceManager struct {
	BaseProxy
}

func NewDataDeviceManager(ctx *Context) *DataDeviceManager {
	ret := new(DataDeviceManager)
	ctx.register(ret)
	return ret
}

func (p *DataDeviceManager) CreateDataSource() (*DataSource , error) {
	ret := NewDataSource(p.Context())
	return ret , p.Context().sendRequest(p,0,Proxy(ret))
}

func (p *DataDeviceManager) GetDataDevice(seat *Seat) (*DataDevice , error) {
	ret := NewDataDevice(p.Context())
	return ret , p.Context().sendRequest(p,1,Proxy(ret),seat)
}

const (
	DataDeviceManagerDndActionNone = 0
	DataDeviceManagerDndActionCopy = 1
	DataDeviceManagerDndActionMove = 2
	DataDeviceManagerDndActionAsk = 4
)

type Shell struct {
	BaseProxy
}

func NewShell(ctx *Context) *Shell {
	ret := new(Shell)
	ctx.register(ret)
	return ret
}

func (p *Shell) GetShellSurface(surface *Surface) (*ShellSurface , error) {
	ret := NewShellSurface(p.Context())
	return ret , p.Context().sendRequest(p,0,Proxy(ret),surface)
}

const (
	ShellErrorRole = 0
)

type ShellSurfacePingEvent struct {
	Serial uint32
}

func (p *ShellSurface) AddPingHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.pingHandlers = append(p.pingHandlers , h)
		p.mu.Unlock()
	}
}

func (p *ShellSurface) RemovePingHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.pingHandlers {
		if e == h {
			p.pingHandlers = append(p.pingHandlers[:i] , p.pingHandlers[i+1:]...)
			break
		}
	}
}

type ShellSurfaceConfigureEvent struct {
	Edges uint32
	Width int32
	Height int32
}

func (p *ShellSurface) AddConfigureHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.configureHandlers = append(p.configureHandlers , h)
		p.mu.Unlock()
	}
}

func (p *ShellSurface) RemoveConfigureHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.configureHandlers {
		if e == h {
			p.configureHandlers = append(p.configureHandlers[:i] , p.configureHandlers[i+1:]...)
			break
		}
	}
}

type ShellSurfacePopupDoneEvent struct {
}

func (p *ShellSurface) AddPopupDoneHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.popupDoneHandlers = append(p.popupDoneHandlers , h)
		p.mu.Unlock()
	}
}

func (p *ShellSurface) RemovePopupDoneHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.popupDoneHandlers {
		if e == h {
			p.popupDoneHandlers = append(p.popupDoneHandlers[:i] , p.popupDoneHandlers[i+1:]...)
			break
		}
	}
}

func (p *ShellSurface) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.pingHandlers) > 0 {
			ev := ShellSurfacePingEvent{}
			ev.Serial = event.Uint32()
			p.mu.RLock()
			for _, h := range p.pingHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.configureHandlers) > 0 {
			ev := ShellSurfaceConfigureEvent{}
			ev.Edges = event.Uint32()
			ev.Width = event.Int32()
			ev.Height = event.Int32()
			p.mu.RLock()
			for _, h := range p.configureHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 2:
		if len(p.popupDoneHandlers) > 0 {
			ev := ShellSurfacePopupDoneEvent{}
			p.mu.RLock()
			for _, h := range p.popupDoneHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type ShellSurface struct {
	BaseProxy
	mu sync.RWMutex
	pingHandlers []Handler
	configureHandlers []Handler
	popupDoneHandlers []Handler
}

func NewShellSurface(ctx *Context) *ShellSurface {
	ret := new(ShellSurface)
	ctx.register(ret)
	return ret
}

func (p *ShellSurface) Pong(serial uint32) error {
	return p.Context().sendRequest(p,0,serial)
}

func (p *ShellSurface) Move(seat *Seat,serial uint32) error {
	return p.Context().sendRequest(p,1,seat,serial)
}

func (p *ShellSurface) Resize(seat *Seat,serial uint32,edges uint32) error {
	return p.Context().sendRequest(p,2,seat,serial,edges)
}

func (p *ShellSurface) SetToplevel() error {
	return p.Context().sendRequest(p,3)
}

func (p *ShellSurface) SetTransient(parent *Surface,x int32,y int32,flags uint32) error {
	return p.Context().sendRequest(p,4,parent,x,y,flags)
}

func (p *ShellSurface) SetFullscreen(method uint32,framerate uint32,output *Output) error {
	return p.Context().sendRequest(p,5,method,framerate,output)
}

func (p *ShellSurface) SetPopup(seat *Seat,serial uint32,parent *Surface,x int32,y int32,flags uint32) error {
	return p.Context().sendRequest(p,6,seat,serial,parent,x,y,flags)
}

func (p *ShellSurface) SetMaximized(output *Output) error {
	return p.Context().sendRequest(p,7,output)
}

func (p *ShellSurface) SetTitle(title string) error {
	return p.Context().sendRequest(p,8,title)
}

func (p *ShellSurface) SetClass(class_ string) error {
	return p.Context().sendRequest(p,9,class_)
}

const (
	ShellSurfaceResizeNone = 0
	ShellSurfaceResizeTop = 1
	ShellSurfaceResizeBottom = 2
	ShellSurfaceResizeLeft = 4
	ShellSurfaceResizeTopLeft = 5
	ShellSurfaceResizeBottomLeft = 6
	ShellSurfaceResizeRight = 8
	ShellSurfaceResizeTopRight = 9
	ShellSurfaceResizeBottomRight = 10
)

const (
	ShellSurfaceTransientInactive = 0x1
)

const (
	ShellSurfaceFullscreenMethodDefault = 0
	ShellSurfaceFullscreenMethodScale = 1
	ShellSurfaceFullscreenMethodDriver = 2
	ShellSurfaceFullscreenMethodFill = 3
)

type SurfaceEnterEvent struct {
	Output *Output
}

func (p *Surface) AddEnterHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.enterHandlers = append(p.enterHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Surface) RemoveEnterHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.enterHandlers {
		if e == h {
			p.enterHandlers = append(p.enterHandlers[:i] , p.enterHandlers[i+1:]...)
			break
		}
	}
}

type SurfaceLeaveEvent struct {
	Output *Output
}

func (p *Surface) AddLeaveHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.leaveHandlers = append(p.leaveHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Surface) RemoveLeaveHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.leaveHandlers {
		if e == h {
			p.leaveHandlers = append(p.leaveHandlers[:i] , p.leaveHandlers[i+1:]...)
			break
		}
	}
}

func (p *Surface) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.enterHandlers) > 0 {
			ev := SurfaceEnterEvent{}
			ev.Output = event.Proxy(p.Context()).(*Output)
			p.mu.RLock()
			for _, h := range p.enterHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.leaveHandlers) > 0 {
			ev := SurfaceLeaveEvent{}
			ev.Output = event.Proxy(p.Context()).(*Output)
			p.mu.RLock()
			for _, h := range p.leaveHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Surface struct {
	BaseProxy
	mu sync.RWMutex
	enterHandlers []Handler
	leaveHandlers []Handler
}

func NewSurface(ctx *Context) *Surface {
	ret := new(Surface)
	ctx.register(ret)
	return ret
}

func (p *Surface) Destroy() error {
	return p.Context().sendRequest(p,0)
}

func (p *Surface) Attach(buffer *Buffer,x int32,y int32) error {
	return p.Context().sendRequest(p,1,buffer,x,y)
}

func (p *Surface) Damage(x int32,y int32,width int32,height int32) error {
	return p.Context().sendRequest(p,2,x,y,width,height)
}

func (p *Surface) Frame() (*Callback , error) {
	ret := NewCallback(p.Context())
	return ret , p.Context().sendRequest(p,3,Proxy(ret))
}

func (p *Surface) SetOpaqueRegion(region *Region) error {
	return p.Context().sendRequest(p,4,region)
}

func (p *Surface) SetInputRegion(region *Region) error {
	return p.Context().sendRequest(p,5,region)
}

func (p *Surface) Commit() error {
	return p.Context().sendRequest(p,6)
}

func (p *Surface) SetBufferTransform(transform int32) error {
	return p.Context().sendRequest(p,7,transform)
}

func (p *Surface) SetBufferScale(scale int32) error {
	return p.Context().sendRequest(p,8,scale)
}

func (p *Surface) DamageBuffer(x int32,y int32,width int32,height int32) error {
	return p.Context().sendRequest(p,9,x,y,width,height)
}

const (
	SurfaceErrorInvalidScale = 0
	SurfaceErrorInvalidTransform = 1
)

type SeatCapabilitiesEvent struct {
	Capabilities uint32
}

func (p *Seat) AddCapabilitiesHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.capabilitiesHandlers = append(p.capabilitiesHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Seat) RemoveCapabilitiesHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.capabilitiesHandlers {
		if e == h {
			p.capabilitiesHandlers = append(p.capabilitiesHandlers[:i] , p.capabilitiesHandlers[i+1:]...)
			break
		}
	}
}

type SeatNameEvent struct {
	Name string
}

func (p *Seat) AddNameHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.nameHandlers = append(p.nameHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Seat) RemoveNameHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.nameHandlers {
		if e == h {
			p.nameHandlers = append(p.nameHandlers[:i] , p.nameHandlers[i+1:]...)
			break
		}
	}
}

func (p *Seat) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.capabilitiesHandlers) > 0 {
			ev := SeatCapabilitiesEvent{}
			ev.Capabilities = event.Uint32()
			p.mu.RLock()
			for _, h := range p.capabilitiesHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.nameHandlers) > 0 {
			ev := SeatNameEvent{}
			ev.Name = event.String()
			p.mu.RLock()
			for _, h := range p.nameHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Seat struct {
	BaseProxy
	mu sync.RWMutex
	capabilitiesHandlers []Handler
	nameHandlers []Handler
}

func NewSeat(ctx *Context) *Seat {
	ret := new(Seat)
	ctx.register(ret)
	return ret
}

func (p *Seat) GetPointer() (*Pointer , error) {
	ret := NewPointer(p.Context())
	return ret , p.Context().sendRequest(p,0,Proxy(ret))
}

func (p *Seat) GetKeyboard() (*Keyboard , error) {
	ret := NewKeyboard(p.Context())
	return ret , p.Context().sendRequest(p,1,Proxy(ret))
}

func (p *Seat) GetTouch() (*Touch , error) {
	ret := NewTouch(p.Context())
	return ret , p.Context().sendRequest(p,2,Proxy(ret))
}

func (p *Seat) Release() error {
	return p.Context().sendRequest(p,3)
}

const (
	SeatCapabilityPointer = 1
	SeatCapabilityKeyboard = 2
	SeatCapabilityTouch = 4
)

type PointerEnterEvent struct {
	Serial uint32
	Surface *Surface
	SurfaceX float32
	SurfaceY float32
}

func (p *Pointer) AddEnterHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.enterHandlers = append(p.enterHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveEnterHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.enterHandlers {
		if e == h {
			p.enterHandlers = append(p.enterHandlers[:i] , p.enterHandlers[i+1:]...)
			break
		}
	}
}

type PointerLeaveEvent struct {
	Serial uint32
	Surface *Surface
}

func (p *Pointer) AddLeaveHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.leaveHandlers = append(p.leaveHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveLeaveHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.leaveHandlers {
		if e == h {
			p.leaveHandlers = append(p.leaveHandlers[:i] , p.leaveHandlers[i+1:]...)
			break
		}
	}
}

type PointerMotionEvent struct {
	Time uint32
	SurfaceX float32
	SurfaceY float32
}

func (p *Pointer) AddMotionHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.motionHandlers = append(p.motionHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveMotionHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.motionHandlers {
		if e == h {
			p.motionHandlers = append(p.motionHandlers[:i] , p.motionHandlers[i+1:]...)
			break
		}
	}
}

type PointerButtonEvent struct {
	Serial uint32
	Time uint32
	Button uint32
	State uint32
}

func (p *Pointer) AddButtonHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.buttonHandlers = append(p.buttonHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveButtonHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.buttonHandlers {
		if e == h {
			p.buttonHandlers = append(p.buttonHandlers[:i] , p.buttonHandlers[i+1:]...)
			break
		}
	}
}

type PointerAxisEvent struct {
	Time uint32
	Axis uint32
	Value float32
}

func (p *Pointer) AddAxisHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.axisHandlers = append(p.axisHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveAxisHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.axisHandlers {
		if e == h {
			p.axisHandlers = append(p.axisHandlers[:i] , p.axisHandlers[i+1:]...)
			break
		}
	}
}

type PointerFrameEvent struct {
}

func (p *Pointer) AddFrameHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.frameHandlers = append(p.frameHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveFrameHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.frameHandlers {
		if e == h {
			p.frameHandlers = append(p.frameHandlers[:i] , p.frameHandlers[i+1:]...)
			break
		}
	}
}

type PointerAxisSourceEvent struct {
	AxisSource uint32
}

func (p *Pointer) AddAxisSourceHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.axisSourceHandlers = append(p.axisSourceHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveAxisSourceHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.axisSourceHandlers {
		if e == h {
			p.axisSourceHandlers = append(p.axisSourceHandlers[:i] , p.axisSourceHandlers[i+1:]...)
			break
		}
	}
}

type PointerAxisStopEvent struct {
	Time uint32
	Axis uint32
}

func (p *Pointer) AddAxisStopHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.axisStopHandlers = append(p.axisStopHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveAxisStopHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.axisStopHandlers {
		if e == h {
			p.axisStopHandlers = append(p.axisStopHandlers[:i] , p.axisStopHandlers[i+1:]...)
			break
		}
	}
}

type PointerAxisDiscreteEvent struct {
	Axis uint32
	Discrete int32
}

func (p *Pointer) AddAxisDiscreteHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.axisDiscreteHandlers = append(p.axisDiscreteHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Pointer) RemoveAxisDiscreteHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.axisDiscreteHandlers {
		if e == h {
			p.axisDiscreteHandlers = append(p.axisDiscreteHandlers[:i] , p.axisDiscreteHandlers[i+1:]...)
			break
		}
	}
}

func (p *Pointer) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.enterHandlers) > 0 {
			ev := PointerEnterEvent{}
			ev.Serial = event.Uint32()
			ev.Surface = event.Proxy(p.Context()).(*Surface)
			ev.SurfaceX = event.Float32()
			ev.SurfaceY = event.Float32()
			p.mu.RLock()
			for _, h := range p.enterHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.leaveHandlers) > 0 {
			ev := PointerLeaveEvent{}
			ev.Serial = event.Uint32()
			ev.Surface = event.Proxy(p.Context()).(*Surface)
			p.mu.RLock()
			for _, h := range p.leaveHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 2:
		if len(p.motionHandlers) > 0 {
			ev := PointerMotionEvent{}
			ev.Time = event.Uint32()
			ev.SurfaceX = event.Float32()
			ev.SurfaceY = event.Float32()
			p.mu.RLock()
			for _, h := range p.motionHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 3:
		if len(p.buttonHandlers) > 0 {
			ev := PointerButtonEvent{}
			ev.Serial = event.Uint32()
			ev.Time = event.Uint32()
			ev.Button = event.Uint32()
			ev.State = event.Uint32()
			p.mu.RLock()
			for _, h := range p.buttonHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 4:
		if len(p.axisHandlers) > 0 {
			ev := PointerAxisEvent{}
			ev.Time = event.Uint32()
			ev.Axis = event.Uint32()
			ev.Value = event.Float32()
			p.mu.RLock()
			for _, h := range p.axisHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 5:
		if len(p.frameHandlers) > 0 {
			ev := PointerFrameEvent{}
			p.mu.RLock()
			for _, h := range p.frameHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 6:
		if len(p.axisSourceHandlers) > 0 {
			ev := PointerAxisSourceEvent{}
			ev.AxisSource = event.Uint32()
			p.mu.RLock()
			for _, h := range p.axisSourceHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 7:
		if len(p.axisStopHandlers) > 0 {
			ev := PointerAxisStopEvent{}
			ev.Time = event.Uint32()
			ev.Axis = event.Uint32()
			p.mu.RLock()
			for _, h := range p.axisStopHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 8:
		if len(p.axisDiscreteHandlers) > 0 {
			ev := PointerAxisDiscreteEvent{}
			ev.Axis = event.Uint32()
			ev.Discrete = event.Int32()
			p.mu.RLock()
			for _, h := range p.axisDiscreteHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Pointer struct {
	BaseProxy
	mu sync.RWMutex
	enterHandlers []Handler
	leaveHandlers []Handler
	motionHandlers []Handler
	buttonHandlers []Handler
	axisHandlers []Handler
	frameHandlers []Handler
	axisSourceHandlers []Handler
	axisStopHandlers []Handler
	axisDiscreteHandlers []Handler
}

func NewPointer(ctx *Context) *Pointer {
	ret := new(Pointer)
	ctx.register(ret)
	return ret
}

func (p *Pointer) SetCursor(serial uint32,surface *Surface,hotspot_x int32,hotspot_y int32) error {
	return p.Context().sendRequest(p,0,serial,surface,hotspot_x,hotspot_y)
}

func (p *Pointer) Release() error {
	return p.Context().sendRequest(p,1)
}

const (
	PointerErrorRole = 0
)

const (
	PointerButtonStateReleased = 0
	PointerButtonStatePressed = 1
)

const (
	PointerAxisVerticalScroll = 0
	PointerAxisHorizontalScroll = 1
)

const (
	PointerAxisSourceWheel = 0
	PointerAxisSourceFinger = 1
	PointerAxisSourceContinuous = 2
)

type KeyboardKeymapEvent struct {
	Format uint32
	Fd uintptr
	Size uint32
}

func (p *Keyboard) AddKeymapHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.keymapHandlers = append(p.keymapHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Keyboard) RemoveKeymapHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.keymapHandlers {
		if e == h {
			p.keymapHandlers = append(p.keymapHandlers[:i] , p.keymapHandlers[i+1:]...)
			break
		}
	}
}

type KeyboardEnterEvent struct {
	Serial uint32
	Surface *Surface
	Keys []int32
}

func (p *Keyboard) AddEnterHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.enterHandlers = append(p.enterHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Keyboard) RemoveEnterHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.enterHandlers {
		if e == h {
			p.enterHandlers = append(p.enterHandlers[:i] , p.enterHandlers[i+1:]...)
			break
		}
	}
}

type KeyboardLeaveEvent struct {
	Serial uint32
	Surface *Surface
}

func (p *Keyboard) AddLeaveHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.leaveHandlers = append(p.leaveHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Keyboard) RemoveLeaveHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.leaveHandlers {
		if e == h {
			p.leaveHandlers = append(p.leaveHandlers[:i] , p.leaveHandlers[i+1:]...)
			break
		}
	}
}

type KeyboardKeyEvent struct {
	Serial uint32
	Time uint32
	Key uint32
	State uint32
}

func (p *Keyboard) AddKeyHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.keyHandlers = append(p.keyHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Keyboard) RemoveKeyHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.keyHandlers {
		if e == h {
			p.keyHandlers = append(p.keyHandlers[:i] , p.keyHandlers[i+1:]...)
			break
		}
	}
}

type KeyboardModifiersEvent struct {
	Serial uint32
	ModsDepressed uint32
	ModsLatched uint32
	ModsLocked uint32
	Group uint32
}

func (p *Keyboard) AddModifiersHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.modifiersHandlers = append(p.modifiersHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Keyboard) RemoveModifiersHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.modifiersHandlers {
		if e == h {
			p.modifiersHandlers = append(p.modifiersHandlers[:i] , p.modifiersHandlers[i+1:]...)
			break
		}
	}
}

type KeyboardRepeatInfoEvent struct {
	Rate int32
	Delay int32
}

func (p *Keyboard) AddRepeatInfoHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.repeatInfoHandlers = append(p.repeatInfoHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Keyboard) RemoveRepeatInfoHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.repeatInfoHandlers {
		if e == h {
			p.repeatInfoHandlers = append(p.repeatInfoHandlers[:i] , p.repeatInfoHandlers[i+1:]...)
			break
		}
	}
}

func (p *Keyboard) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.keymapHandlers) > 0 {
			ev := KeyboardKeymapEvent{}
			ev.Format = event.Uint32()
			ev.Fd = event.FD()
			ev.Size = event.Uint32()
			p.mu.RLock()
			for _, h := range p.keymapHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.enterHandlers) > 0 {
			ev := KeyboardEnterEvent{}
			ev.Serial = event.Uint32()
			ev.Surface = event.Proxy(p.Context()).(*Surface)
			ev.Keys = event.Array()
			p.mu.RLock()
			for _, h := range p.enterHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 2:
		if len(p.leaveHandlers) > 0 {
			ev := KeyboardLeaveEvent{}
			ev.Serial = event.Uint32()
			ev.Surface = event.Proxy(p.Context()).(*Surface)
			p.mu.RLock()
			for _, h := range p.leaveHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 3:
		if len(p.keyHandlers) > 0 {
			ev := KeyboardKeyEvent{}
			ev.Serial = event.Uint32()
			ev.Time = event.Uint32()
			ev.Key = event.Uint32()
			ev.State = event.Uint32()
			p.mu.RLock()
			for _, h := range p.keyHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 4:
		if len(p.modifiersHandlers) > 0 {
			ev := KeyboardModifiersEvent{}
			ev.Serial = event.Uint32()
			ev.ModsDepressed = event.Uint32()
			ev.ModsLatched = event.Uint32()
			ev.ModsLocked = event.Uint32()
			ev.Group = event.Uint32()
			p.mu.RLock()
			for _, h := range p.modifiersHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 5:
		if len(p.repeatInfoHandlers) > 0 {
			ev := KeyboardRepeatInfoEvent{}
			ev.Rate = event.Int32()
			ev.Delay = event.Int32()
			p.mu.RLock()
			for _, h := range p.repeatInfoHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Keyboard struct {
	BaseProxy
	mu sync.RWMutex
	keymapHandlers []Handler
	enterHandlers []Handler
	leaveHandlers []Handler
	keyHandlers []Handler
	modifiersHandlers []Handler
	repeatInfoHandlers []Handler
}

func NewKeyboard(ctx *Context) *Keyboard {
	ret := new(Keyboard)
	ctx.register(ret)
	return ret
}

func (p *Keyboard) Release() error {
	return p.Context().sendRequest(p,0)
}

const (
	KeyboardKeymapFormatNoKeymap = 0
	KeyboardKeymapFormatXkbV1 = 1
)

const (
	KeyboardKeyStateReleased = 0
	KeyboardKeyStatePressed = 1
)

type TouchDownEvent struct {
	Serial uint32
	Time uint32
	Surface *Surface
	Id int32
	X float32
	Y float32
}

func (p *Touch) AddDownHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.downHandlers = append(p.downHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Touch) RemoveDownHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.downHandlers {
		if e == h {
			p.downHandlers = append(p.downHandlers[:i] , p.downHandlers[i+1:]...)
			break
		}
	}
}

type TouchUpEvent struct {
	Serial uint32
	Time uint32
	Id int32
}

func (p *Touch) AddUpHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.upHandlers = append(p.upHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Touch) RemoveUpHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.upHandlers {
		if e == h {
			p.upHandlers = append(p.upHandlers[:i] , p.upHandlers[i+1:]...)
			break
		}
	}
}

type TouchMotionEvent struct {
	Time uint32
	Id int32
	X float32
	Y float32
}

func (p *Touch) AddMotionHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.motionHandlers = append(p.motionHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Touch) RemoveMotionHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.motionHandlers {
		if e == h {
			p.motionHandlers = append(p.motionHandlers[:i] , p.motionHandlers[i+1:]...)
			break
		}
	}
}

type TouchFrameEvent struct {
}

func (p *Touch) AddFrameHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.frameHandlers = append(p.frameHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Touch) RemoveFrameHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.frameHandlers {
		if e == h {
			p.frameHandlers = append(p.frameHandlers[:i] , p.frameHandlers[i+1:]...)
			break
		}
	}
}

type TouchCancelEvent struct {
}

func (p *Touch) AddCancelHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.cancelHandlers = append(p.cancelHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Touch) RemoveCancelHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.cancelHandlers {
		if e == h {
			p.cancelHandlers = append(p.cancelHandlers[:i] , p.cancelHandlers[i+1:]...)
			break
		}
	}
}

func (p *Touch) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.downHandlers) > 0 {
			ev := TouchDownEvent{}
			ev.Serial = event.Uint32()
			ev.Time = event.Uint32()
			ev.Surface = event.Proxy(p.Context()).(*Surface)
			ev.Id = event.Int32()
			ev.X = event.Float32()
			ev.Y = event.Float32()
			p.mu.RLock()
			for _, h := range p.downHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.upHandlers) > 0 {
			ev := TouchUpEvent{}
			ev.Serial = event.Uint32()
			ev.Time = event.Uint32()
			ev.Id = event.Int32()
			p.mu.RLock()
			for _, h := range p.upHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 2:
		if len(p.motionHandlers) > 0 {
			ev := TouchMotionEvent{}
			ev.Time = event.Uint32()
			ev.Id = event.Int32()
			ev.X = event.Float32()
			ev.Y = event.Float32()
			p.mu.RLock()
			for _, h := range p.motionHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 3:
		if len(p.frameHandlers) > 0 {
			ev := TouchFrameEvent{}
			p.mu.RLock()
			for _, h := range p.frameHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 4:
		if len(p.cancelHandlers) > 0 {
			ev := TouchCancelEvent{}
			p.mu.RLock()
			for _, h := range p.cancelHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Touch struct {
	BaseProxy
	mu sync.RWMutex
	downHandlers []Handler
	upHandlers []Handler
	motionHandlers []Handler
	frameHandlers []Handler
	cancelHandlers []Handler
}

func NewTouch(ctx *Context) *Touch {
	ret := new(Touch)
	ctx.register(ret)
	return ret
}

func (p *Touch) Release() error {
	return p.Context().sendRequest(p,0)
}

type OutputGeometryEvent struct {
	X int32
	Y int32
	PhysicalWidth int32
	PhysicalHeight int32
	Subpixel int32
	Make string
	Model string
	Transform int32
}

func (p *Output) AddGeometryHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.geometryHandlers = append(p.geometryHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Output) RemoveGeometryHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.geometryHandlers {
		if e == h {
			p.geometryHandlers = append(p.geometryHandlers[:i] , p.geometryHandlers[i+1:]...)
			break
		}
	}
}

type OutputModeEvent struct {
	Flags uint32
	Width int32
	Height int32
	Refresh int32
}

func (p *Output) AddModeHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.modeHandlers = append(p.modeHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Output) RemoveModeHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.modeHandlers {
		if e == h {
			p.modeHandlers = append(p.modeHandlers[:i] , p.modeHandlers[i+1:]...)
			break
		}
	}
}

type OutputDoneEvent struct {
}

func (p *Output) AddDoneHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.doneHandlers = append(p.doneHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Output) RemoveDoneHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.doneHandlers {
		if e == h {
			p.doneHandlers = append(p.doneHandlers[:i] , p.doneHandlers[i+1:]...)
			break
		}
	}
}

type OutputScaleEvent struct {
	Factor int32
}

func (p *Output) AddScaleHandler(h Handler) {
	if h != nil {
		p.mu.Lock()
		p.scaleHandlers = append(p.scaleHandlers , h)
		p.mu.Unlock()
	}
}

func (p *Output) RemoveScaleHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i , e := range p.scaleHandlers {
		if e == h {
			p.scaleHandlers = append(p.scaleHandlers[:i] , p.scaleHandlers[i+1:]...)
			break
		}
	}
}

func (p *Output) Dispatch(event *Event) {
	switch event.opcode {
	case 0:
		if len(p.geometryHandlers) > 0 {
			ev := OutputGeometryEvent{}
			ev.X = event.Int32()
			ev.Y = event.Int32()
			ev.PhysicalWidth = event.Int32()
			ev.PhysicalHeight = event.Int32()
			ev.Subpixel = event.Int32()
			ev.Make = event.String()
			ev.Model = event.String()
			ev.Transform = event.Int32()
			p.mu.RLock()
			for _, h := range p.geometryHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 1:
		if len(p.modeHandlers) > 0 {
			ev := OutputModeEvent{}
			ev.Flags = event.Uint32()
			ev.Width = event.Int32()
			ev.Height = event.Int32()
			ev.Refresh = event.Int32()
			p.mu.RLock()
			for _, h := range p.modeHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 2:
		if len(p.doneHandlers) > 0 {
			ev := OutputDoneEvent{}
			p.mu.RLock()
			for _, h := range p.doneHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	case 3:
		if len(p.scaleHandlers) > 0 {
			ev := OutputScaleEvent{}
			ev.Factor = event.Int32()
			p.mu.RLock()
			for _, h := range p.scaleHandlers {
				h.Handle(ev)
			}
			p.mu.RUnlock()
		}
	}
}

type Output struct {
	BaseProxy
	mu sync.RWMutex
	geometryHandlers []Handler
	modeHandlers []Handler
	doneHandlers []Handler
	scaleHandlers []Handler
}

func NewOutput(ctx *Context) *Output {
	ret := new(Output)
	ctx.register(ret)
	return ret
}

func (p *Output) Release() error {
	return p.Context().sendRequest(p,0)
}

const (
	OutputSubpixelUnknown = 0
	OutputSubpixelNone = 1
	OutputSubpixelHorizontalRgb = 2
	OutputSubpixelHorizontalBgr = 3
	OutputSubpixelVerticalRgb = 4
	OutputSubpixelVerticalBgr = 5
)

const (
	OutputTransformNormal = 0
	OutputTransform90 = 1
	OutputTransform180 = 2
	OutputTransform270 = 3
	OutputTransformFlipped = 4
	OutputTransformFlipped90 = 5
	OutputTransformFlipped180 = 6
	OutputTransformFlipped270 = 7
)

const (
	OutputModeCurrent = 0x1
	OutputModePreferred = 0x2
)

type Region struct {
	BaseProxy
}

func NewRegion(ctx *Context) *Region {
	ret := new(Region)
	ctx.register(ret)
	return ret
}

func (p *Region) Destroy() error {
	return p.Context().sendRequest(p,0)
}

func (p *Region) Add(x int32,y int32,width int32,height int32) error {
	return p.Context().sendRequest(p,1,x,y,width,height)
}

func (p *Region) Subtract(x int32,y int32,width int32,height int32) error {
	return p.Context().sendRequest(p,2,x,y,width,height)
}

type Subcompositor struct {
	BaseProxy
}

func NewSubcompositor(ctx *Context) *Subcompositor {
	ret := new(Subcompositor)
	ctx.register(ret)
	return ret
}

func (p *Subcompositor) Destroy() error {
	return p.Context().sendRequest(p,0)
}

func (p *Subcompositor) GetSubsurface(surface *Surface,parent *Surface) (*Subsurface , error) {
	ret := NewSubsurface(p.Context())
	return ret , p.Context().sendRequest(p,1,Proxy(ret),surface,parent)
}

const (
	SubcompositorErrorBadSurface = 0
)

type Subsurface struct {
	BaseProxy
}

func NewSubsurface(ctx *Context) *Subsurface {
	ret := new(Subsurface)
	ctx.register(ret)
	return ret
}

func (p *Subsurface) Destroy() error {
	return p.Context().sendRequest(p,0)
}

func (p *Subsurface) SetPosition(x int32,y int32) error {
	return p.Context().sendRequest(p,1,x,y)
}

func (p *Subsurface) PlaceAbove(sibling *Surface) error {
	return p.Context().sendRequest(p,2,sibling)
}

func (p *Subsurface) PlaceBelow(sibling *Surface) error {
	return p.Context().sendRequest(p,3,sibling)
}

func (p *Subsurface) SetSync() error {
	return p.Context().sendRequest(p,4)
}

func (p *Subsurface) SetDesync() error {
	return p.Context().sendRequest(p,5)
}

const (
	SubsurfaceErrorBadSurface = 0
)
