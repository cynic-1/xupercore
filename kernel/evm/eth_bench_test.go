package evm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/burrow/crypto"
	_ "github.com/xuperchain/xupercore/bcs/contract/evm"
	_ "github.com/xuperchain/xupercore/bcs/contract/native"
	_ "github.com/xuperchain/xupercore/bcs/contract/xvm"
	"github.com/xuperchain/xupercore/kernel/contract"
	_ "github.com/xuperchain/xupercore/kernel/contract"
	_ "github.com/xuperchain/xupercore/kernel/contract/kernel"
	_ "github.com/xuperchain/xupercore/kernel/contract/manager"
	"github.com/xuperchain/xupercore/kernel/contract/mock"
	"io/ioutil"
	"testing"
)

func BenchmarkEVM(b *testing.B) {
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	bin, err := ioutil.ReadFile("testdata/counter.bin")
	if err != nil {
		b.Error(err)
		return
	}
	abi, err := ioutil.ReadFile("testdata/counter.abi")
	if err != nil {
		b.Error(err)
		return
	}
	args := map[string][]byte{
		"contract_abi": abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}
	data, err := hex.DecodeString(string((bin)))
	if err != nil {
		b.Fatal(err)
	}
	resp, err := th.Deploy("evm", "counter", "counter", data, args)
	if err != nil {
		b.Fatal(err)
	}
	b.Run("Benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := th.Invoke("evm", "counter", "increase", map[string][]byte{
				"input":       []byte(`{"key":"xchain"}`),
				"jsonEncoded": []byte("true"),
			})
			if err != nil {
				b.Error(err)
				return
			}
		}
	})
	_ = resp

}

func TestExecEVMMapContract(t *testing.T) {
	// 创建合约配置，设置EVM驱动
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}

	// 根据配置创建一个执行合约的Helper
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	// 读取合约的ABI文件和BIN文件
	bin, err := ioutil.ReadFile("testdata/map.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err := ioutil.ReadFile("testdata/map.abi")
	if err != nil {
		t.Error(err)
		return
	}
	data, err := hex.DecodeString(string((bin)))
	if err != nil {
		t.Fatal(err)
	}

	// 部署合约
	args := map[string][]byte{
		"contract_abi": abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}
	resp, err := th.Deploy("evm", "MapContract", "MapContract", data, args)
	if err != nil {
		t.Fatal(err)
	}

	// 调用合约中的函数PUT 存入数据
	resp, err = th.Invoke("evm", "MapContract", "put", map[string][]byte{
		"input":       []byte(`{"_key":"chaincode","_value":"111"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}

	// 调用合约中的函数GET 获取数据
	resp, err = th.Invoke("evm", "MapContract", "get", map[string][]byte{
		"input":       []byte(`{"_key":"chaincode"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp.Body))

}

func TestExecEVMIdentityIdentifierContract(t *testing.T) {
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	bin, err := ioutil.ReadFile("testdata/IdentityIdentifier.bin")
	if err != nil {
		t.Error(err)
		return
	}
	_abi, err := ioutil.ReadFile("testdata/IdentityIdentifier.abi")
	if err != nil {
		t.Error(err)
		return
	}
	args := map[string][]byte{
		"contract_abi": _abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}
	data, err := hex.DecodeString(string((bin)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := th.Deploy("evm", "IdentityIdentifier", "test", data, args)
	if err != nil {
		t.Fatal(err)
	}

	test_args := make(map[string]interface{})
	test_args["username"] = "yzy"
	test_args["uOwner"] = crypto.Address{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	test_args["expiryTime"] = "22123121"

	input, err := json.Marshal(test_args)
	fmt.Println(string(input))

	for i := 0; i < 1; i++ {
		resp, err = th.Invoke("evm", "test", "registerIdentity", map[string][]byte{
			"input":       input,
			"jsonEncoded": []byte("true"),
		})
		if err != nil {
			t.Error(err)
			return
		}
	}

	t.Log(string(resp.Body))

}

func TestTest1Contract(t *testing.T) {
	// 创建合约配置，设置EVM驱动
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}

	// 根据配置创建一个执行合约的Helper
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	// 读取合约的ABI文件和BIN文件
	bin, err := ioutil.ReadFile("testdata/test1.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err := ioutil.ReadFile("testdata/test1.abi")
	if err != nil {
		t.Error(err)
		return
	}
	data, err := hex.DecodeString(string((bin)))
	if err != nil {
		t.Fatal(err)
	}

	// 部署合约
	args := map[string][]byte{
		"contract_abi": abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}
	resp, err := th.Deploy("evm", "Test1Contract", "Test1Contract", data, args)
	if err != nil {
		t.Fatal(err)
	}

	// 调用合约中的函数PUT 存入数据
	resp, err = th.Invoke("evm", "Test1Contract", "mint", map[string][]byte{
		"input":       []byte(`{"receiver":"5B38Da6a701c568545dCfcB03FcB875f56beddC4","amount":"1000"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}

	// 调用合约中的函数GET 获取数据
	resp, err = th.Invoke("evm", "Test1Contract", "balances", map[string][]byte{
		"input":       []byte(`{"":"5B38Da6a701c568545dCfcB03FcB875f56beddC4"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp.Body))
}

func TestTest2Contract(t *testing.T) {
	// 创建合约配置，设置EVM驱动
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}

	// 根据配置创建一个执行合约的Helper
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	// 读取合约的ABI文件和BIN文件
	bin, err := ioutil.ReadFile("testdata/test1.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err := ioutil.ReadFile("testdata/test1.abi")
	if err != nil {
		t.Error(err)
		return
	}
	data, err := hex.DecodeString(string((bin)))
	if err != nil {
		t.Fatal(err)
	}

	// 部署合约
	args := map[string][]byte{
		"contract_abi": abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}
	resp, err := th.Deploy("evm", "Test1Contract", "Test1Contract", data, args)
	if err != nil {
		t.Fatal(err)
	}

	// 调用合约中的函数PUT 存入数据
	resp, err = th.Invoke("evm", "Test1Contract", "mint", map[string][]byte{
		"input":       []byte(`{"receiver":"5B38Da6a701c568545dCfcB03FcB875f56beddC4","amount":"1000"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}

	// 调用合约中的函数GET 获取数据
	resp, err = th.Invoke("evm", "Test1Contract", "balances", map[string][]byte{
		"input":       []byte(`{"":"5B38Da6a701c568545dCfcB03FcB875f56beddC4"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp.Body))
}

func TestContractACallB(t *testing.T) {
	// 创建合约配置，设置EVM驱动
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}

	// 根据配置创建一个执行合约的Helper
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	// 读取合约的ABI文件和BIN文件
	bin, err := ioutil.ReadFile("testdata/B.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err := ioutil.ReadFile("testdata/B.abi")
	if err != nil {
		t.Error(err)
		return
	}
	data, err := hex.DecodeString(string((bin)))
	if err != nil {
		t.Fatal(err)
	}

	// 部署合约
	args := map[string][]byte{
		"contract_abi": abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}
	resp, err := th.Deploy("evm", "TestB", "TestB", data, args)
	if err != nil {
		t.Fatal(err)
	}

	// 读取合约的ABI文件和BIN文件
	bin, err = ioutil.ReadFile("testdata/A.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err = ioutil.ReadFile("testdata/A.abi")
	if err != nil {
		t.Error(err)
		return
	}
	data, err = hex.DecodeString(string((bin)))
	if err != nil {
		t.Fatal(err)
	}

	// 部署合约
	args = map[string][]byte{
		"contract_abi": abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}

	resp, err = th.Deploy("evm", "TestA", "TestA", data, args)
	if err != nil {
		t.Fatal(err)
	}

	// 调用合约中的函数GET 获取数据
	resp, err = th.Invoke("evm", "TestA", "setVars", map[string][]byte{
		"input":       []byte(`{"_contract":"313131312D2D2D2D2D2D2D2D2D2D2D5465737442","_num":"100"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp.Body))
}

func TestEMISCallIdentity(t *testing.T) {
	// 创建合约配置，设置EVM驱动
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}

	// 根据配置创建一个执行合约的Helper
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	// 读取合约的ABI文件和BIN文件
	bin, err := ioutil.ReadFile("testdata/IdentityIdentifier.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err := ioutil.ReadFile("testdata/IdentityIdentifier.abi")
	if err != nil {
		t.Error(err)
		return
	}
	data, err := hex.DecodeString(string((bin)))
	if err != nil {
		t.Fatal(err)
	}

	// 部署合约
	args := map[string][]byte{
		"contract_abi": abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}
	resp, err := th.Deploy("evm", "IdentityIdentifier", "test", data, args)
	if err != nil {
		t.Fatal(err)
	}

	// 读取合约的ABI文件和BIN文件
	bin, err = ioutil.ReadFile("testdata/EMIS.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err = ioutil.ReadFile("testdata/EMIS.abi")
	if err != nil {
		t.Error(err)
		return
	}
	data, err = hex.DecodeString(string((bin)))
	if err != nil {
		t.Fatal(err)
	}

	// 部署合约
	args = map[string][]byte{
		"contract_abi": abi,
		"input":        bin,
		"jsonEncoded":  []byte("false"),
	}

	resp, err = th.Deploy("evm", "EMIS", "EMIS", data, args)
	if err != nil {
		t.Fatal(err)
	}

	// 调用合约中的函数register 注册用户名
	resp, err = th.Invoke("evm", "EMIS", "register", map[string][]byte{
		"input":       []byte(`{"username":"zf","owner":"5B38Da6a701c568545dCfcB03FcB875f56beddC4"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}

	// 调用合约中的函数GET 获取数据
	resp, err = th.Invoke("evm", "EMIS", "setIdentity", map[string][]byte{
		"input":       []byte(`{"_identityIdentifier":"313131312D2D2D2D2D2D2D2D2D2D2D2D74657374"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}

	resp, err = th.Invoke("evm", "EMIS", "registerIdentityIdentifier", map[string][]byte{
		"input":       []byte(`{"username":"zf","transferAmount":"100"}`),
		"jsonEncoded": []byte("true"),
	})

	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp.Body))
}

func TestBallotContract(t *testing.T) {
	// 创建合约配置，设置EVM驱动
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}

	// 根据配置创建一个执行合约的Helper
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	// 读取合约的ABI文件和BIN文件
	bin, err := ioutil.ReadFile("testdata/Ballot.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err := ioutil.ReadFile("testdata/Ballot.abi")
	if err != nil {
		t.Error(err)
		return
	}
	// 带构造不能够进行编码！
	//data, err := hex.DecodeString(string((bin)))
	//if err != nil {
	//	t.Fatal(err)
	//}

	//var argsStr = "{\"proposalNames\": [\"0x6f7074696f6e4100000000000000000000000000000000000000000000000000\",\"0x6f7074696f6e4200000000000000000000000000000000000000000000000000\",\"0x6f7074696f6e4300000000000000000000000000000000000000000000000000\"]}"
	var argsStr = `{"proposalNames": ["0x6f7074696f6e4100000000000000000000000000000000000000000000000000","0x6f7074696f6e4200000000000000000000000000000000000000000000000000","0x6f7074696f6e4300000000000000000000000000000000000000000000000000","0x6f7074696f6e4400000000000000000000000000000000000000000000000000"]}`
	// 部署合约
	args := map[string]interface{}{
		//"proposalNames": []string{"0x6f7074696f6e4100000000000000000000000000000000000000000000000000", "0x6f7074696f6e4200000000000000000000000000000000000000000000000000", "0x6f7074696f6e4300000000000000000000000000000000000000000000000000", "0x6f7074696f6e4400000000000000000000000000000000000000000000000000"},
	}
	err = json.Unmarshal([]byte(argsStr), &args)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(args)

	var x3args map[string][]byte
	if x3args, err = convertToXuper3EvmArgs(args); err != nil {
		t.Error(err)
		return
	}

	resp, err := th.DeployWithABI("evm", "Ballot", "Ballot", bin, abi, x3args)
	if err != nil {
		t.Fatal(err)
	}

	// 调用合约中的函数PUT 存入数据
	resp, err = th.Invoke("evm", "Ballot", "giveRightToVote", map[string][]byte{
		"input":       []byte(`{"voter":"3131313231313131313131313131313131313132"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(resp.Body))
}

func convertToXuper3EvmArgs(args map[string]interface{}) (map[string][]byte, error) {
	input, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	// 此处与 server 端结构相同，如果 jsonEncoded 字段修改，server 端也要修改（core/contract/evm/creator.go）。
	ret := map[string][]byte{
		"input":       input,
		"jsonEncoded": []byte("true"),
	}
	return ret, nil
}

func TestExecEVMEMISContract(t *testing.T) {
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true,
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{
			Enable: true,
			Driver: "native",
		},
		EVM: contract.EVMConfig{
			Enable: true,
			Driver: "evm",
		},
		LogDriver: mock.NewMockLogger(),
	}
	th := mock.NewTestHelper(contractConfig)
	defer th.Close()

	bin, err := ioutil.ReadFile("testdata/IdentityIdentifier.bin")
	if err != nil {
		t.Error(err)
		return
	}
	_abi, err := ioutil.ReadFile("testdata/IdentityIdentifier.abi")
	if err != nil {
		t.Error(err)
		return
	}

	args := map[string]interface{}{}

	var x3args map[string][]byte
	if x3args, err = convertToXuper3EvmArgs(args); err != nil {
		t.Error(err)
		return
	}

	resp, err := th.DeployWithABI("evm", "IdentityIdentifier", "test", bin, _abi, x3args)
	if err != nil {
		t.Fatal(err)
	}

	bin, err = ioutil.ReadFile("testdata/EMIS_test.bin")
	if err != nil {
		t.Error(err)
		return
	}
	abi, err := ioutil.ReadFile("testdata/EMIS_test.abi")
	if err != nil {
		t.Error(err)
		return
	}

	args = map[string]interface{}{
		"_identityIdentifier": "313131312D2D2D2D2D2D2D2D2D2D2D2D74657374",
		"_identityPrice":      "100000000",
	}

	if x3args, err = convertToXuper3EvmArgs(args); err != nil {
		t.Error(err)
		return
	}

	resp, err = th.DeployWithABI("evm", "EMIS_test", "EMIS_test", bin, abi, x3args)
	if err != nil {
		t.Fatal(err)
	}

	// 调用合约中的函数register 注册用户名
	resp, err = th.Invoke("evm", "EMIS_test", "register", map[string][]byte{
		"input":       []byte(`{"username":"zf","owner":"5B38Da6a701c568545dCfcB03FcB875f56beddC4"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = th.Invoke("evm", "EMIS_test", "updateAboutMe", map[string][]byte{
		"input":       []byte(`{"username":"zf","aboutMe":"123456"}`),
		"jsonEncoded": []byte("true"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp.Body))
}
