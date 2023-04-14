package mock

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/xuperchain/xupercore/kernel/contract/bridge"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/xuperchain/xupercore/kernel/contract"
	"github.com/xuperchain/xupercore/kernel/contract/sandbox"
	"github.com/xuperchain/xupercore/kernel/ledger"
	"github.com/xuperchain/xupercore/kernel/permission/acl/utils"
	"github.com/xuperchain/xupercore/protos"
)

const (
	ContractAccount      = "XC1111111111111111@xuper"
	ContractAccount2     = "XC2222222222222222@xuper"
	FeaturesContractName = "features"
	// <AccountPrefix><AccountNumber>[@<BlockChainName>]
)

type TestHelper struct {
	basedir    string              // 目录
	utxo       *contract.UTXORWSet // UTXO读写集
	utxoReader sandbox.UtxoReader  // UTXO 读取器
	state      *sandbox.MemXModel  // 沙盒内存状态 使用红黑树保存合约执行过程中的参数 中间变量等
	manager    contract.Manager    // 合约管理器
}

func NewTestHelper(cfg *contract.ContractConfig) *TestHelper {
	// 创建临时目录
	basedir, err := ioutil.TempDir("", "contract-test")
	if err != nil {
		panic(err)
	}

	// 生成一个沙盒状态
	state := sandbox.NewMemXModel()
	// 用于测试的 区块链核心结构 需要重新实现一个
	core := new(fakeChainCore)
	// 创建合约管理器
	m, err := contract.CreateManager("default", &contract.ManagerConfig{
		Basedir:  basedir, // 临时目录
		BCName:   "xuper", // 区块链名字
		Core:     core,    // 区块链核
		XMReader: state,   // 读写状态
		Config:   cfg,     // 合约配置
	})
	if err != nil {
		panic(err)
	}

	// 使用合约管理器 和 沙盒状态 生成一个测试Helper
	th := &TestHelper{
		basedir: basedir,
		manager: m,
		state:   state,
	}

	// 初始化账户 使用的是UTXO类型账户 MIS暂不支持 需要后续实现
	th.initAccount()
	return th
}

func (t *TestHelper) Manager() contract.Manager {
	return t.manager
}

func (t *TestHelper) Basedir() string {
	return t.basedir
}

func (t *TestHelper) State() *sandbox.MemXModel {
	return t.state
}
func (t *TestHelper) UTXOState() *contract.UTXORWSet {
	return t.utxo
}

func (t *TestHelper) initAccount() {
	// 账户的相关信息都存储在这个子树中 key = XCAccount + ContractAccount value = VersionedData
	t.state.Put(utils.GetAccountBucket(), []byte(ContractAccount), &ledger.VersionedData{
		RefTxid:  []byte("txid"),
		PureData: nil,
	})

	utxoReader := sandbox.NewUTXOReaderFromInput([]*protos.TxInput{
		{
			RefTxid:      nil,                          // 事务ID
			RefOffset:    0,                            //
			FromAddr:     []byte(FeaturesContractName), // 合约名
			Amount:       big.NewInt(99999999).Bytes(), // 账户余额
			FrozenHeight: 0,                            // 冻结区块个数
		},
	})

	t.utxoReader = utxoReader
}

func (t *TestHelper) Deploy(module, lang, contractName string, bin []byte, args map[string][]byte) (*contract.Response, error) {
	m := t.Manager()
	state, err := m.NewStateSandbox(&contract.SandboxConfig{ // 沙盒配置
		XMReader:   t.State(),    // 沙盒内存状态
		UTXOReader: t.utxoReader, // UTXO读集
	})
	if err != nil {
		return nil, err
	}

	//
	ctx, err := m.NewContext(&contract.ContextConfig{
		Module:         "xkernel",
		ContractName:   "$contract",
		State:          state,
		ResourceLimits: contract.MaxLimits,
		Initiator:      ContractAccount,
	})
	if err != nil {
		return nil, err
	}

	desc := &protos.WasmCodeDesc{
		Runtime:      lang,
		ContractType: module,
	}
	descbuf, _ := proto.Marshal(desc)

	argsBuf, _ := json.Marshal(args)

	invokeArgs := map[string][]byte{
		"account_name":  []byte(ContractAccount),
		"contract_name": []byte(contractName),
		"contract_code": bin,
		"contract_desc": descbuf,
		"init_args":     argsBuf,
	}
	if bridge.ContractType(module) == bridge.TypeEvm {
		invokeArgs["contract_abi"] = args["contract_abi"]
	}
	// 调用部署合约方法
	resp, err := ctx.Invoke("deployContract", invokeArgs)

	if err != nil {
		return nil, err
	}

	// 释放资源
	ctx.Release()
	t.Commit(state) // 提交
	return resp, nil
}

func (t *TestHelper) DeployWithABI(module, lang, contractName string, bin []byte, abi []byte, args map[string][]byte) (*contract.Response, error) {
	m := t.Manager()
	state, err := m.NewStateSandbox(&contract.SandboxConfig{ // 沙盒配置
		XMReader:   t.State(),    // 沙盒内存状态
		UTXOReader: t.utxoReader, // UTXO读集
	})
	if err != nil {
		return nil, err
	}

	//
	ctx, err := m.NewContext(&contract.ContextConfig{
		Module:         "xkernel",
		ContractName:   "$contract",
		State:          state,
		ResourceLimits: contract.MaxLimits,
		Initiator:      ContractAccount,
	})
	if err != nil {
		return nil, err
	}

	desc := &protos.WasmCodeDesc{
		Runtime:      lang,
		ContractType: module,
	}
	descbuf, _ := proto.Marshal(desc)

	argsBuf, _ := json.Marshal(args)
	fmt.Println(string(argsBuf))
	invokeArgs := map[string][]byte{
		"account_name":  []byte(ContractAccount),
		"contract_name": []byte(contractName),
		"contract_code": bin,
		"contract_desc": descbuf,
		"init_args":     argsBuf,
	}
	if bridge.ContractType(module) == bridge.TypeEvm {
		invokeArgs["contract_abi"] = abi
	}
	// 调用部署合约方法
	resp, err := ctx.Invoke("deployContract", invokeArgs)

	if err != nil {
		return nil, err
	}

	// 释放资源
	ctx.Release()
	t.Commit(state) // 提交
	return resp, nil
}

func (t *TestHelper) Upgrade(contractName string, bin []byte) error {
	m := t.Manager()
	state, err := m.NewStateSandbox(&contract.SandboxConfig{
		XMReader:   t.State(),
		UTXOReader: t.utxoReader,
	})
	if err != nil {
		return err
	}

	ctx, err := m.NewContext(&contract.ContextConfig{
		Module:         "xkernel",
		ContractName:   "$contract",
		State:          state,
		ResourceLimits: contract.MaxLimits,
	})
	if err != nil {
		return err
	}

	_, err = ctx.Invoke("upgradeContract", map[string][]byte{
		"contract_name": []byte(contractName),
		"contract_code": bin,
	})
	ctx.Release()
	t.Commit(state)
	return err
}

func (t *TestHelper) Invoke(module, contractName, method string, args map[string][]byte) (*contract.Response, error) {
	m := t.Manager()
	state, err := m.NewStateSandbox(&contract.SandboxConfig{
		XMReader:   t.State(),
		UTXOReader: t.utxoReader,
	})
	if err != nil {
		return nil, err
	}

	ctx, err := m.NewContext(&contract.ContextConfig{
		Module:         module,
		ContractName:   contractName,
		State:          state,
		ResourceLimits: contract.MaxLimits,
		Initiator:      ContractAccount,
	})
	if err != nil {
		return nil, err
	}
	defer ctx.Release()

	resp, err := ctx.Invoke(method, args)
	if err != nil {
		return nil, err
	}
	state.Flush()
	t.utxo = state.UTXORWSet()
	t.Commit(state)
	return resp, nil
}

func (t *TestHelper) Commit(state contract.StateSandbox) {
	rwset := state.RWSet()    // 获取读写集
	txbuf := make([]byte, 32) //
	rand.Read(txbuf)
	for i, w := range rwset.WSet { // 遍历读写集
		t.state.Put(w.Bucket, w.Key, &ledger.VersionedData{ //
			RefTxid:   txbuf,    // 一个随机数
			RefOffset: int32(i), //
			PureData: &ledger.PureData{
				Bucket: w.Bucket, // 哪个账户
				Key:    w.Key,    // 状态key
				Value:  w.Value,  // 状态value
			},
		})
	}
}

func (t *TestHelper) Close() {
	os.RemoveAll(t.basedir)
}
