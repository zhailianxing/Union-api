package brand

import (
	"os"
	"bufio"
	"strings"
	"fmt"
)

type BrandDict struct {
	BrandToModel map[string][]string
	BrandToBrand map[string][]string
}

func Load() {
	brandToModel := make(map[string][]string)
	brandToBrand := make(map[string][]string)
	brandDict := &BrandDict{
		brandToModel,
		brandToBrand,
	}
	datas := make([][]string, 0)

	f, err := os.Open("liuliangchi.csv")
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			panic(err)

		}
		strs := strings.Split(line, ",")
		ss := make([]string, 2)
		for i := 0; i < 2; i++ {
			s := strs[i]
			s1 := strings.ToLower(s)
			s2 := strings.TrimSpace(s1)
			ss[i] = s2
		}
		datas = append(datas, ss)
	}

	fmt.Println(datas[0:10])

	loss := make([][]string, 0)
	for _, data := range datas {
		title := CheckBrand(data[0], data[1])
		if title == "" {
			loss = append(loss, data)
		} else {
			models := brandDict.BrandToModel[title]
			if !Equal(data[0], models...) {
				models = append(models, data[0])
				brandDict.BrandToModel[title] = models
			}

			brands := brandDict.BrandToBrand[title]
			if !Equal(data[1], brands...) {
				brands = append(brands, data[1])
				brandDict.BrandToBrand[title] = brands
			}
		}
	}

	for i, k := 0, 0; i < len(loss); i++ {
		if loss[i][1] != "" {
			fmt.Println(loss[i])
			k++
		}
		if k == 10 {
			break
		}
	}

	//write file
	f2, err := os.Create("brandDict.txt")
	for brandTitle, brands := range brandDict.BrandToBrand {
		f2.WriteString("[")
		f2.WriteString(brandTitle)
		f2.WriteString("]")
		f2.WriteString("\n")
		for _, brand := range brands {
			f2.WriteString(brand)
			f2.WriteString("\n")
		}
	}

	f3, err := os.Create("modelDict.txt")
	for brandTitle, models := range brandDict.BrandToModel {
		f3.WriteString("[")
		f3.WriteString(brandTitle)
		f3.WriteString("]")
		f3.WriteString("\n")
		for _, model := range models {
			f3.WriteString(model)
			f3.WriteString("\n")
		}
	}
}

func CheckBrand(model, brand string) (brandTitle string) {
	if model == "" && brand == "" {
		return
	}
	if Contain(brand, "vivo") {
		return "vivo"
	} else if Contain(brand, "huawei", "-bac-al00-") {
		return "华为"
	} else if Contain(brand, "xiaomi", "redmi", "-mi-") {
		return "小米"
	} else if Contain(brand, "oppo") {
		return "OPPO"
	} else if Contain(brand, "iphone", "apple") {
		return "苹果"
	} else if Contain(brand, "samsung", "sm") {
		return "三星"
	} else if Contain(brand, "ivvi") {
		return "ivvi"
	} else if Contain(brand, "letv", "lemobile") {
		return "乐视"
	} else if Contain(brand, "lenovo") {
		return "联想"
	} else if Contain(brand, "meizu", "-m3-note-", "-m5s-") {
		return "魅族"
	} else if Contain(brand, "honor", "-plk-al10-", "-stf-al10-", "-che-tl00h-") {
		return "华为荣耀"
	} else if Contain(brand, "nubia") {
		return "努比亚"
	} else if Contain(brand, "gionee") {
		return "金立"
	} else if Contain(brand, "kopo") {
		return "酷珀"
	} else if Contain(brand, "htc") {
		return "HTC"
	} else if Contain(brand, "sharp") {
		return "夏普"
	} else if Contain(brand, "coolpad", "yulong") {
		return "酷派"
	} else if Contain(brand, "zte") {
		return "中兴"
	} else if Contain(brand, "alps") {
		return "阿尔卑斯"
	} else if Contain(brand, "sony") {
		return "索尼"
	} else if Contain(brand, "bird") {
		return "波导"
	} else if Contain(brand, "lephone") {
		return "乐丰"
	} else if Contain(brand, "cmdc") {
		return "中国移动"
	} else if Contain(brand, "shown") {
		return "首云"
	} else if Contain(brand, "asus") {
		return "华硕"
	} else if Contain(brand, "xiaolajiao") {
		return "小辣椒"
	} else if Contain(brand, "shijitianyuan") {
		return "世纪天元"
	} else if Contain(brand, "yusun") {
		return "语信"
	} else if Contain(brand, "qiku") {
		return "奇酷"
	} else if Contain(brand, "sanchi") {
		return "三驰"
	} else if Contain(brand, "k-touch") {
		return "天语"
	} else if Contain(brand, "oysin") {
		return "欧亚信"
	}

	return
}

func Contain(src string, subStrs ...string) bool {
	for _, subStr := range subStrs {
		if strings.Contains(src, subStr) {
			return true
		}
	}

	return false
}

func Equal(src string, subStrs ...string) bool {
	for _, subStr := range subStrs {
		if src == subStr {
			return true
		}
	}

	return false
}
