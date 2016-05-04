package i18n

/**
{
  "English": {
    "test": "test{{.prj}}",
    "main": {
      "welcom": "welcome {{.Person}} or %s %d",
      "welcom_use": "welcome use",
      "nav": {
        "home": "HOME",
        "publish": "PUBLISH"
      }
    }
  }
}

	 bobMap := map[string]interface{}{"Person": "B-ob"}
	bobStruct := struct{ Person string }{Person: "Bob"}

	i18n.Init("/Users/justinh/Development/Workspace/GoLang/untitled/langs/testdata")
	fmt.Println("all:",i18n.Languages())
	fmt.Println("test--->:",i18n.Tr("en","test"))
	fmt.Println("main.welcom--->:",i18n.Tr("en","main.welcom",bobMap,"ss",88))
	fmt.Println("main.welcom--->:",i18n.Tr("en","main.welcom",bobStruct))
	fmt.Println("main.welcom--->:",i18n.Tr("en","main.welcom","kkd",99))

	lg:=i18n.Language("zh")
	fmt.Println(lg.Code,lg.Name)
	fmt.Println("test--->:",lg.Tr("test"))

	lg.Update("main.welcom","欢迎修改")
	fmt.Println("main.welcom--->:",lg.Tr("main.welcom"))
	lg.Add("main.www","www特色")
	fmt.Println("main.www--->:",lg.Tr("main.www"))

	fmt.Println("main.xxxx--->:",lg.Tr("main.xxxx"))

	//写文件测试
	file,_:=os.OpenFile("/Users/justinh/"+lg.TranslateFilename(),os.O_RDWR|os.O_CREATE,0666)
	lg.Write(file)
	defer file.Close()

 */
import (

	"path"
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"reflect"
	"text/template"
	"bytes"
	"io"
	"errors"
)

type (
	translation struct {
		Language
		Store map[string]interface{}
	}

	translations struct {
		path            string
		defaultLangTag string
		langs           []*translation
	}



	Language struct {
		Tag string
		Name  string
		Country string
		Translation  *translation
	}

     mlanguage struct {
	Name  string
	Country string
	Translation  map[string]interface{}
      }

        TranslateFunc func(translationID string, args ...interface{}) string
)

var locales = new(translations)

func LanguageTagAll()map[string]string{
	return   map[string]string{
		"en": "English",
		"en-us": "English(America)",
		"zh":"中文",
		"zh-cn":  "中文(简体)",
		"zh-tw":  "中文(繁體)",
		"ja":  "日本语",
		"ko":  "한국어", //朝鲜语
		"fr":  "Français", //法语
		"vi":  "Tiếng Việt", //越南
		"th":  "ไทย", //泰语
		"tr":  "Türkçe", //土耳其语
		"el":  "Ελληνικά", //希腊语
		"sv":  "Svenska", //瑞典语
		"hu":  "Magyar", //匈牙利语
		"nb":  "Bokmål", //挪威语
		"de":  "Deutsch", //德语
		"ru":  "Русский", //俄语
		"es":  "Español", //西班牙语
		"it":  "Italiano", //意大利语
		"nl":  "Nederlands", //荷兰语
		"pt":  "Português", //葡萄牙语
		"cs":  "Čeština", //捷克语
		"pl":  "Polski", //波兰语
		"he":  "עברית" ,//希伯来语
		"fa":  "فارسی" ,//波斯语
		"ar":  "العربية" ,//阿拉伯语
	}
}

func LanguageName(languageTag string)string{
	return LanguageTagAll()[strings.ToLower(languageTag)]

}

//设置默认语言
func SetDefaultLang(languageTag string) {
	locales.defaultLangTag = languageTag
}

//获取所有语言名称/代码列表
func Translations() []Language {
	ts:=[]Language{}
	for _, v := range locales.langs {
		ts=append(ts,v.Language)
	}
	return ts
}

//新建语言
func NewLanguage(languageTag,name,country string)*translation{
	lang:= new(translation)
	lang.Tag=languageTag
	lang.Name=name
	lang.Country=country
	lang.Translation=lang
	lang.Store=make(map[string]interface{})
	return lang
}

//根据TAG匹配最相近语言(优先级别: 1-完全匹配 2-国家匹配)
func TranslationMatch(languageTag string, languageTags ...string) *translation {
	//生成所有相关语言tag
	tags:=ParseAccept_Language(languageTag)
	if len(languageTags)>0{
		for _,l:=range languageTags {
			for _,ll:=range  ParseAccept_Language(l) {
				tags = append(tags, ll)
			}
		}
	}
	tags=append(tags,locales.defaultLangTag)



	for _,t:=range tags {
		matchesT := matchingTags(t)
		for _,l:=range locales.langs{
			matchesL:=matchingTags(l.Tag)

			for _,mL:=range matchesL{
				for _,mT:=range matchesT{
					if mL== mT {
						return l
					}
				}
			}
		}
	}

	return nil
}

//加载语言路径中所有语言
func Init(i18nPath string,defLanguageTag string) error {
	locales.path=i18nPath
	locales.defaultLangTag = defLanguageTag
	files,err:=ioutil.ReadDir(i18nPath)
	if err!=nil{
		return err
	}

	for _,v:=range files{
		if strings.ToLower(path.Ext(v.Name()))==".json" {
			LoadTranslation(strings.TrimSuffix(v.Name(),".json"))
		}
	}
	return nil
}

//加载一个文件到内存
func LoadTranslation(languageTag string) (*translation,error) {
	filename:=path.Join(locales.path,languageTag+".json")
	bytes,err:=ioutil.ReadFile(filename)
	if err!=nil{
		return nil,err
	}

	lang:=locales.indexof(languageTag)
	if lang==nil{
		lang=new(translation)
		lang.Tag=languageTag
		lang.Translation=lang
		locales.langs=append(locales.langs,lang)
	}



        tt:=mlanguage{}
	err=json.Unmarshal(bytes,&tt)
	if err!=nil{

		return nil, err
	}
	lang.Name=tt.Name
	lang.Country=tt.Country
	lang.Store=tt.Translation


	//fmt.Println("LangCode:",lang.LangCode)
	//fmt.Println("LangDesc:",lang.desc)
	//fmt.Println("store:",lang.store)
	// fmt.Println("--->",lang.store["test"])
	//fmt.Println("--->",lang.store["main"].(map[string]interface{})["welcom"])
      return lang,nil

}

//语言合并补全
func MergeTranslation(){

}

//定位语言(tag完全匹配)
func (langs *translations)indexof(languageTag string) *translation{
	for _, v := range locales.langs {
		if strings.ToLower(v.Tag)==strings.ToLower(languageTag) {
			return v
		}
	}
	return nil
}

//func (langs *translations)indexof(languageTag string,otherLanguageTag ...string) *translation{
//
//	tags:=  make([]string, 0, len(otherLanguageTag))
//	tags=append(tags,languageTag)
//	tags=append(tags,otherLanguageTag...)
//	tags=append(tags,langs.defaultLangTag)
//	for _,v:=range tags{
//		t:=langs._indexof(v)
//		if t!=nil{
//		   return t
//		}
//	}
//
//	return nil
//}

func   matchingTags(tag string) []string {
	parts := strings.Split(tag, "-")
	var prefix, matches []string
	for _, part := range parts {
		prefix = append(prefix, part)
		match := strings.Join(prefix, "-")
		matches = append(matches, match)
	}
	return matches
}

func (lang *translation)TranslateFilename()string{
	return lang.Tag+".json"
}

func (lang *translation)Template(tag string) string{
     tags:=strings.Split(tag,".")
	store:=lang.Store

	for _,v:=range tags{
                s:=store[v]
		//不存在
		if s==nil{
		   return tag
		}
		if str,ok:=s.(string);ok{
			return str
		}
		store=s.(map[string]interface{})
	}
	return ""
}

//翻译
func Tr(langCode, translationID string, args ...interface{}) string {
	lang:=locales.indexof(langCode)
	if lang==nil{
		return ""
	}
	return  lang.Tr(translationID,args...)
}


//翻译
func (lang *translation)Tr(translationID string, args ...interface{}) (string) {
      formatstr:= lang.Template(translationID)
	if len(args)==0 {
		return formatstr
	}
	if !strings.Contains(formatstr, "{{") {
		return formatstr
	}
	var structParam interface{}
	params := make([]interface{}, 0, len(args))
	for _,v:=range args{
		val :=reflect.TypeOf(v)
		if  val.Kind() == reflect.Struct || val.Kind()==reflect.Map{
			structParam=v
		}else{
			params=append(params,v)
		}
	}
	if len(params)>0  {
		formatstr=fmt.Sprintf(formatstr,params...)
	}
	if  structParam!=nil && strings.Contains(formatstr, "{{") {
		tpl, err := template.New("data").Parse(formatstr)

		if err!=nil{
			return formatstr
		}
		var buf bytes.Buffer
		err=tpl.Execute(&buf,structParam)
		formatstr=buf.String()
	}

     return formatstr
}

//增加语言条目
func (lang *translation)Add(tag, templatestr string) error{
	tags:=strings.Split(tag,".")

	store:=lang.Store

	for i,v:=range tags{
		if i<len(tags)-1 {
			if store[v]==nil {
				store[v]= make(map[string]interface{})
			}
			if _,ok:=store[v].(string);ok{
			   return errors.New("Can't  Override the section:<"+strings.Join(tags[0:i+1],".")+">")
			}
			store=store[v].(map[string]interface{})
		}else{
			if stv,ok:=store[v].(map[string]interface{});ok && len(stv)>0{
				return errors.New("Can't  Override the section:<"+strings.Join(tags[0:i+1],".")+">")
			}
			store[v]=templatestr
			return nil
		}
	}
	return nil
}

//删除语言条目
func (lang *translation)Del(tag string){
	tags:=strings.Split(tag,".")

	store:=lang.Store

	for i,v:=range tags{
		if store[v]==nil {
			return
		}

		if i<len(tags)-1 {
			store=store[v].(map[string]interface{})
		}else{
			delete(store,v)
			return
		}
	}
}

//修改语言条目
func (lang *translation)Update(tag,templatestr string ) error{
	tags:=strings.Split(tag,".")

	store:=lang.Store
	for _,v:=range tags{
		s:=store[v]
		if s==nil{
			return errors.New("Don't exist")
		}
		if _,ok:=s.(string);ok{
			store[v]=templatestr
			return nil
		}
		store=s.(map[string]interface{})
	}
	return nil
}



//保存语言
func (lang *translation)Write(w io.Writer) error{
	tt:=mlanguage{}
	tt.Country=lang.Country
	tt.Name=lang.Name
	tt.Translation =lang.Store
	 bytes,err:=json.Marshal(&tt)
	if err!=nil{
		return err
	}
	_,err=w.Write(bytes)
	return err
}

//cookie Accept-Language 头解析
func ParseAccept_Language(src string) []string {
	var langs []string
	start := 0
	for end, chr := range src {
		switch chr {
		case ',', ';', '.':
			tag := strings.TrimSpace(src[start:end])
			if len(tag)>1 && strings.Index(tag,"=")<0 {
				langs = append(langs, tag)
			}
			start = end + 1
		}
	}

	if start > 0 {
		tag := strings.TrimSpace(src[start:])
			langs = append(langs, tag)
		return langs
	}

	langs=append(langs,src)

	return langs
}