package command

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/shyiko/jabba/cfg"
	"github.com/shyiko/jabba/semver"
)

type byOS map[string]byArch
type byArch map[string]byDistribution
type byDistribution map[string]map[string]string

func LsRemote(os, arch string) (map[*semver.Version]string, error) {
	//读取远程配置
	cnt, err := fetch(cfg.Index())
	//开发时,读取本地配置(读取的是当前文件夹的index.json)
	// cnt, err := ioutil.ReadFile("../index.json")
	//打包时,读取本地配置(读取的是当前文件夹的index.json)
	// cnt, err := ioutil.ReadFile("index.json")
	if err != nil {
		return nil, err
	}
	var index byOS
	err = json.Unmarshal(cnt, &index)
	if err != nil {
		return nil, err
	}
	releaseMap := make(map[*semver.Version]string)
	for key, value := range index[os][arch] {
		var prefix string
		if key != "jdk" {
			if !strings.Contains(key, "@") {
				continue
			}
			prefix = key[strings.Index(key, "@")+1:] + "@"
		}
		for ver, url := range value {
			v, err := semver.ParseVersion(prefix + ver)
			if err != nil {
				return nil, err
			}
			releaseMap[v] = url
		}
	}
	return releaseMap, nil
}

func fetch(url string) (content []byte, err error) {
	client := http.Client{Transport: RedirectTracer{}}
	res, err := client.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, errors.New("GET " + url + " returned " + strconv.Itoa(res.StatusCode))
	}
	content, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}
