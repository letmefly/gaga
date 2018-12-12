package services

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"sync"
	"time"
)

import (
	"types"
	"utils"

	etcdclient "github.com/coreos/etcd/client"
	"google.golang.org/grpc"
)

// https://github.com/letmefly/ant/blob/master/samples/template/sdk/rpc.gen.go

type ServiceConf struct {
	ServiceType    string   `json:"serviceType"`
	ServiceId      string   `json:"serviceId"`
	ServiceAddr    string   `json:"serviceAddr"`
	ServiceUseList []string `json:"serviceUseList"`
	ProtoUseList   []string `json:protoUseList`
	IsStream       bool     `json:"isStream"`
	TTL            int      `json:"ttl"`
}

type ServiceClient struct {
	use  int
	Conn *grpc.ClientConn
	Conf *ServiceConf
}

type RecvHandler func([]byte) error
type ServiceBundle struct {
	serviceType string
	//handlers       map[string]RecvHandler    // (sessId, handler)
	serviceClients map[string]*ServiceClient // (serviceId, ClientConn)
}

func newServiceBundle(serviceType string) *ServiceBundle {
	s := &ServiceBundle{
		serviceType: serviceType,
		//handlers:       make(map[string]RecvHandler, 0),
		serviceClients: make(map[string]*ServiceClient, 0),
	}
	return s
}

/*
func (s *ServiceBundle) AddRecvHandler(owner string, handler RecvHandler) {
	s.handlers[owner] = handler
}

func (s *ServiceBundle) Call(owner string, data []byte) ([]byte, error) {
	return nil, nil
}
*/

type ServiceManager struct {
	servicesRoot     string
	onlineServices   map[string]*ServiceConf
	serviceBundleMap map[string]*ServiceBundle
	toServiceTypeMap map[uint32]string // msgId to service type
	toMsgNameMap     map[uint32]string // msgId to msgName
	toMsgIdMap       map[string]uint32 // msgName 2 msgId
	currServiceConf  *ServiceConf
	etcdcli          etcdclient.Client
	mu               *sync.RWMutex
}

func (m *ServiceManager) init(ctx context.Context, conf *ServiceConf) {
	m.currServiceConf = conf
	m.onlineServices = make(map[string]*ServiceConf, 0)
	m.serviceBundleMap = make(map[string]*ServiceBundle, 0)
	m.toServiceTypeMap = make(map[uint32]string, 0)
	m.toMsgNameMap = make(map[uint32]string, 0)
	m.toMsgIdMap = make(map[string]uint32, 0)
	m.mu = new(sync.RWMutex)
	m.servicesRoot = "/serviceList/"

	cfg := etcdclient.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               etcdclient.DefaultTransport,
		HeaderTimeoutPerRequest: 3 * time.Second,
	}

	cli, err := etcdclient.New(cfg)
	if err != nil {
		log.Println("etcdclient new fail")
		log.Fatal(err)
	}
	m.etcdcli = cli
	m.registerMyself(ctx)
	m.fetchOnlineServices(ctx)
	m.addOtherServices(ctx)
	go m.watchOnlineServices(ctx)

}

func (m *ServiceManager) registerMyself(ctx context.Context) {
	str, err := json.Marshal(m.currServiceConf)
	if err != nil {
		log.Fatal(err)
		return
	}
	kapi := etcdclient.NewKeysAPI(m.etcdcli)
	_, err1 := kapi.Set(ctx, m.servicesRoot+m.currServiceConf.ServiceId, string(str), &etcdclient.SetOptions{TTL: time.Duration(m.currServiceConf.TTL) * time.Second})
	if err1 != nil {
		log.Fatal(err1)
	}

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(int(m.currServiceConf.TTL/2)))
			_, err2 := kapi.Set(ctx, m.servicesRoot+m.currServiceConf.ServiceId, "", &etcdclient.SetOptions{Refresh: true, TTL: time.Duration(m.currServiceConf.TTL) * time.Second})
			if err2 != nil {
				log.Fatal(err2)
				return
			}
		}
	}()
}

func (m *ServiceManager) fetchOnlineServices(ctx context.Context) {
	kapi := etcdclient.NewKeysAPI(m.etcdcli)
	ret, err := kapi.Get(ctx, m.servicesRoot, nil)
	if err != nil {
		log.Fatal(err)
	}
	if !ret.Node.Dir {
		log.Fatal("No Services Dir in Etcd")
	}
	for _, node := range ret.Node.Nodes {
		serviceId := filepath.Base(node.Key)
		confStr := node.Value
		log.Println(serviceId, confStr)
		serviceConf := ServiceConf{}
		json.Unmarshal([]byte(confStr), &serviceConf)
		m.onlineServices[serviceId] = &serviceConf
	}
}

func (m *ServiceManager) watchOnlineServices(ctx context.Context) {
	kapi := etcdclient.NewKeysAPI(m.etcdcli)
	watcher := kapi.Watcher(m.servicesRoot, &etcdclient.WatcherOptions{Recursive: true})
	for {
		r, err := watcher.Next(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		switch r.Action {
		case "set", "create", "update", "compareAndSwap":
			serviceId := filepath.Base(r.Node.Key)
			confStr := r.Node.Value
			conf := &ServiceConf{}
			json.Unmarshal([]byte(confStr), conf)
			m.addService(serviceId, conf)
			log.Println("update:", serviceId)
		case "delete", "expire":
			serviceId := filepath.Base(r.Node.Key)
			log.Println("delete:", serviceId)
			m.removeService(serviceId)
		}
	}
}

func (m *ServiceManager) isUse(serviceType string) bool {
	for _, v := range m.currServiceConf.ServiceUseList {
		if v == serviceType || v == "*" {
			return true
		}
	}
	return false
}

// only connect services that diffient from current service type
func (m *ServiceManager) addOtherServices(ctx context.Context) {
	for serviceId, conf := range m.onlineServices {
		m.addService(serviceId, conf)
	}
}

func (m *ServiceManager) addService(serviceId string, conf *ServiceConf) {
	m.mu.Lock()
	defer m.mu.Unlock()
	serviceAddr := conf.ServiceAddr
	serviceType := conf.ServiceType

	m.createTLB(conf)
	if conf.ServiceType == m.currServiceConf.ServiceType {
		return
	}

	if !m.isUse(conf.ServiceType) {
		return
	}
	m.onlineServices[serviceId] = conf
	serviceBundle, ok1 := m.serviceBundleMap[serviceType]
	if !ok1 {
		m.serviceBundleMap[serviceType] = newServiceBundle(serviceType)
		serviceBundle = m.serviceBundleMap[serviceType]
	}

	// this service has been connected
	if _, ok2 := serviceBundle.serviceClients[serviceId]; ok2 {
		return
	}

	log.Println("[services] add service ", serviceId, "->", serviceAddr, "trying")
	conn, err := grpc.Dial(serviceAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return
	}
	serviceBundle.serviceClients[serviceId] = &ServiceClient{Conn: conn, use: 0, Conf: conf}

	log.Println("[services] add service ", serviceId, "->", serviceAddr, "success")
}

func (m *ServiceManager) removeService(serviceId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	conf, ok1 := m.onlineServices[serviceId]
	if !ok1 {
		return
	}
	serviceType := conf.ServiceType
	serviceBundle, ok2 := m.serviceBundleMap[serviceType]
	if !ok2 {
		return
	}
	serviceClient, ok3 := serviceBundle.serviceClients[serviceId]
	if !ok3 {
		return
	}
	serviceClient.Conn.Close()
	delete(serviceBundle.serviceClients, serviceId)
	delete(m.onlineServices, serviceId)
	m.deleteTLB(conf)
	log.Println("[services] remvoe service ", serviceId, "success")
}

func (m *ServiceManager) createTLB(conf *ServiceConf) {
	for _, v := range conf.ProtoUseList {
		msgName := v
		msgId := utils.HashCode(conf.ServiceType + "." + msgName)
		m.toMsgIdMap[msgName] = msgId
		m.toMsgNameMap[msgId] = msgName
		m.toServiceTypeMap[msgId] = conf.ServiceType
	}
}

func (m *ServiceManager) deleteTLB(conf *ServiceConf) {
	for _, v := range conf.ProtoUseList {
		msgName := v
		msgId := utils.HashCode(conf.ServiceType + "." + msgName)
		delete(m.toMsgIdMap, msgName)
		delete(m.toMsgNameMap, msgId)
		delete(m.toServiceTypeMap, msgId)
	}
}

func (m *ServiceManager) assignServiceId(serviceType string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	serviceBundle, ok := m.serviceBundleMap[serviceType]
	if !ok {
		return "", types.ERR_NO_SERVICE
	}
	tmpId, tmpUse := "", int(^uint(0)>>1)
	for k, v := range serviceBundle.serviceClients {
		item := v
		if item.use < tmpUse {
			tmpUse = item.use
			tmpId = k
		}
	}
	serviceBundle.serviceClients[tmpId].use += 1
	return tmpId, nil
}

func (m *ServiceManager) unassignServiceId(serviceId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	serviceType := utils.GetServiceType(serviceId)
	if serviceType == "" {
		return
	}
	serviceBundle, ok1 := m.serviceBundleMap[serviceType]
	if !ok1 {
		return
	}
	serviceClient, ok2 := serviceBundle.serviceClients[serviceId]
	if !ok2 {
		return
	}
	serviceClient.use -= 1
}

func (m *ServiceManager) getServiceClient(serviceId string) (*ServiceClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	serviceType := utils.GetServiceType(serviceId)
	if serviceType == "" {
		return nil, types.ERR_INVALID_SERVICE_ID
	}
	serviceBundle, ok1 := m.serviceBundleMap[serviceType]
	if !ok1 {
		return nil, types.ERR_NO_SERVICE
	}
	serviceClient, ok2 := serviceBundle.serviceClients[serviceId]
	if !ok2 {
		return nil, types.ERR_NO_SERVICE_CLIENT
	}
	return serviceClient, nil
}

var (
	_service_manager ServiceManager
	_once            sync.Once
)

// services module init
func Register(ctx context.Context, conf *ServiceConf) {
	_once.Do(func() { _service_manager.init(ctx, conf) })
}

func GetServiceUseList() []string {
	conf := _service_manager.currServiceConf
	return conf.ServiceUseList
}

func AssignServiceId(serviceType string) (string, error) {
	return _service_manager.assignServiceId(serviceType)
}

func UnassignServiceId(serviceId string) {
	_service_manager.unassignServiceId(serviceId)
}

func GetServiceClient(serviceId string) (*ServiceClient, error) {
	return _service_manager.getServiceClient(serviceId)
}

func ToServiceType(msgId uint32) string {
	_service_manager.mu.RLock()
	defer _service_manager.mu.RUnlock()
	return _service_manager.toServiceTypeMap[msgId]
}

func ToMsgName(msgId uint32) string {
	_service_manager.mu.RLock()
	defer _service_manager.mu.RUnlock()
	return _service_manager.toMsgNameMap[msgId]
}

func ToMsgId(msgName string) uint32 {
	_service_manager.mu.RLock()
	defer _service_manager.mu.RUnlock()
	return _service_manager.toMsgIdMap[msgName]
}
