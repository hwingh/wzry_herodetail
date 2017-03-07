package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	iconv "github.com/djimenez/iconv-go"
)

type Hero struct {
	Title     string
	Name      string
	Attribute Attribute
	Skills    []Skill
}

func (hero Hero) String() string {
	result := ""
	result += "称号: " + hero.Title + "\n"
	result += "名称: " + hero.Name + "\n"
	result += "属性: " + "\n"
	result += "\t" + "生存能力: " + hero.Attribute.Surival + "\n"
	result += "\t" + "攻击伤害: " + hero.Attribute.Damage + "\n"
	result += "\t" + "技能效果: " + hero.Attribute.Effect + "\n"
	result += "\t" + "上手难度: " + hero.Attribute.Diffculty + "\n"
	result += "技能: " + "\n"
	for _, skill := range hero.Skills {
		result += "\t" + "名称: " + skill.Name + "\n"
		result += "\t" + "冷却值: " + strconv.Itoa(skill.CoolValue) + "\n"
		result += "\t" + "消耗: " + strconv.Itoa(skill.Consume) + "\n"
		result += "\t" + "描述: " + skill.Desc + "\n"
	}
	return result
}

type Attribute struct {
	Surival   string
	Damage    string
	Effect    string
	Diffculty string
}

type Skill struct {
	Name      string
	CoolValue int
	Consume   int
	Desc      string
}

func HeroDetail(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	content_s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	content := make([]byte, len(content_s)*2)
	_, _, err = iconv.Convert(content_s, content, "gbk", "utf-8")

	if err != nil {
		log.Fatalln("Cant't convert:", err)
	}

	hero := Hero{}

	//<h3 class="cover-title">沧海之曜</h3>
	titleReg := regexp.MustCompile(`<h3 class="cover-title">(.*)</h3>`)
	title := titleReg.FindSubmatch(content)

	hero.Title = string(title[1])

	//<h2 class="cover-name">大乔</h2>
	nameReg := regexp.MustCompile(`<h2 class="cover-name">(.*)</h2>`)
	name := nameReg.FindSubmatch(content)

	hero.Name = string(name[1])

	//<em class="cover-list-text fl">生存能力</em>

	// attrReg := regexp.MustCompile(`<em class="cover-list-text fl">(.*)</em>`)
	// attr := attrReg.FindAllSubmatch(content, -1)

	// for _, a := range attr {
	// 	fmt.Println(string(a[1]))
	// }

	//<i class="ibar" style="width:40%"></i>
	attrvalReg := regexp.MustCompile(`<i class="(ibar|ibar ibar1)" style="width:(.*)"></i>`)
	attrval := attrvalReg.FindAllSubmatch(content, -1)

	// for _, av := range attrval {
	// 	fmt.Println(string(av[2]))
	// }
	hero.Attribute.Surival = string(attrval[0][2])
	hero.Attribute.Damage = string(attrval[1][2])
	hero.Attribute.Effect = string(attrval[2][2])
	hero.Attribute.Diffculty = string(attrval[3][2])

	//<a href="javascript:void(0)" class="skill-btn" title="制裁仪式">川流不息</a>
	skillReg := regexp.MustCompile(`<a href="javascript:void\(0\)" class="skill-btn" title="(.*)">(.*)</a>`)
	skills := skillReg.FindAllSubmatch(content, -1)

	// for _, skill := range skills {
	// 	fmt.Println(string(skill[2]))
	// }

	//<p class="skill-p1">冷却值：0</p>

	skillP1Reg := regexp.MustCompile(`<p class="skill-p1">(.*)</p>`)
	sp1 := skillP1Reg.FindAllSubmatch(content, -1)

	//for _, sp := range sp1 {
	//	  fmt.Println(string(sp[1]))
	//}

	//<p class="skill-p2">消耗：0</p>

	skillP2Reg := regexp.MustCompile(`<p class="skill-p2">(.*)</P>`)
	sp2 := skillP2Reg.FindAllSubmatch(content, -1)

	// for _, sp := range sp2 {
	// 	fmt.Println(string(sp[1]))
	// }

	//<p class="skill-p3">被动：大乔提升自身与附近友军英雄高额移动速度</p>

	skillP3Reg := regexp.MustCompile(`<p class="skill-p3">(.*)</p>`)
	sp3 := skillP3Reg.FindAllSubmatch(content, -1)

	// for _, sp := range sp3 {
	// 	fmt.Println(string(sp[1]))
	// }

	for i := 0; i < len(sp3); i++ {
		skill := Skill{}
		skill.Name = string(skills[i][2])

		i1 := strings.IndexFunc(string(sp1[i][1]), func(c rune) bool {
			if c >= '0' && c <= '9' {
				return true
			} else {
				return false
			}
		})

		coolvalue, _ := strconv.Atoi(string(sp1[i][1])[i1:])

		i2 := strings.IndexFunc(string(sp2[i][1]), func(c rune) bool {
			if c >= '0' && c <= '9' {
				return true
			} else {
				return false
			}
		})

		consume, _ := strconv.Atoi(string(sp2[i][1])[i2:])

		skill.CoolValue = coolvalue
		skill.Consume = consume
		skill.Desc = string(sp3[i][1])
		hero.Skills = append(hero.Skills, skill)
	}

	fmt.Println(hero)

}

func main() {
	var url_format = "http://pvp.qq.com/web201605/herodetail/%d.shtml"
	num := 100
	for num <= 250 {
		url := fmt.Sprintf(url_format, num)
		num++

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		if _, ok := resp.Header["Expires"]; !ok {
			fmt.Println("*****")
			continue
		}
		fmt.Println("URL: ", url)
		HeroDetail(url)
		fmt.Println("=================================================")
	}
}
