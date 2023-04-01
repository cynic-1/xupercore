package native

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xuperchain/xupercore/kernel/contract"
	_ "github.com/xuperchain/xupercore/kernel/contract/kernel"
	_ "github.com/xuperchain/xupercore/kernel/contract/manager"
	"github.com/xuperchain/xupercore/kernel/contract/mock"
)

const (
	RUNTIME_DOCKER = "docker"
	RUNTIME_HOST   = "host"
	IMAGE_NAME     = "alpine"
)

func compile(th *mock.TestHelper, runtime string) ([]byte, error) {
	target := filepath.Join(th.Basedir(), "counter.bin")
	cmd := exec.Command("go", "build", "-o", target)
	if runtime == RUNTIME_DOCKER {
		cmd.Env = append(os.Environ(), []string{"GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0"}...)
	}
	cmd.Dir = "testdata"
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s:%s", err, out)
	}
	bin, err := ioutil.ReadFile(target)
	if err != nil {
		return nil, err
	}
	return bin, nil
}

func TestNative(t *testing.T) {
	// 生成合约配置
	var contractConfig = &contract.ContractConfig{
		EnableUpgrade: true, //允许更新合约
		Xkernel: contract.XkernelConfig{
			Enable: true,
			Driver: "default",
		},
		Native: contract.NativeConfig{ // 原生合约配置
			Enable: true,
			Driver: "native",
			Docker: contract.NativeDockerConfig{ // Docker配置
				Enable:    true,
				ImageName: IMAGE_NAME,
			},
		},
		LogDriver: mock.NewMockLogger(),
	}

	// 一个是本地跑 一个是另起一个容器跑
	runtimes := []string{RUNTIME_HOST, RUNTIME_DOCKER}

	for _, runtime := range runtimes {
		if runtime == RUNTIME_DOCKER {
			// 判断docker是否能够正常使用
			_, err := exec.Command("docker", "info").CombinedOutput()
			if err != nil {
				t.Skip("docker not available")
			}

			t.Log("pulling image......")
			// 拉取docker镜像
			pullResp, errPull := exec.Command("docker", "pull", IMAGE_NAME).CombinedOutput()
			if errPull != nil {
				t.Error(err)
				continue
				t.Log(string(pullResp))
			}
			contractConfig.Native.Docker.Enable = true
		} else {
			contractConfig.Native.Docker.Enable = false
		}
		// 部署原生合约
		t.Run("TestNativeDeploy_"+runtime, func(t *testing.T) {
			th := mock.NewTestHelper(contractConfig)
			defer th.Close()

			// 编译出二进制文件
			bin, err := compile(th, runtime)
			if err != nil {
				t.Fatal(err)
			}

			// 部署合约
			resp, err := th.Deploy("native", "go", "counter", bin, map[string][]byte{
				"creator": []byte("icexin"),
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%#v", resp)
		})
		t.Run("TestNativeInvoke_"+runtime, func(t *testing.T) {
			th := mock.NewTestHelper(contractConfig)
			defer th.Close()

			bin, err := compile(th, runtime)
			if err != nil {
				t.Fatal(err)
			}

			_, err = th.Deploy("native", "go", "counter", bin, map[string][]byte{
				"creator": []byte("icexin"),
			})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := th.Invoke("native", "counter", "increase", map[string][]byte{
				"key": []byte("k1"),
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("body:%s", resp.Body)
		})

		t.Run("TestNativeUpgrade_"+runtime, func(t *testing.T) {
			th := mock.NewTestHelper(contractConfig)
			defer th.Close()

			bin, err := compile(th, runtime)
			if err != nil {
				t.Fatal(err)
			}

			_, err = th.Deploy("native", "go", "counter", bin, map[string][]byte{
				"creator": []byte("icexin"),
			})
			if err != nil {
				t.Fatal(err)
			}

			err = th.Upgrade("counter", bin)
			if err != nil {
				t.Fatal(err)
			}
		})

	}
}
