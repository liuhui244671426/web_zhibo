package main
//todo 经本机测试guest页面延迟过高,不宜使用此方案
import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8888", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var channel_msg = make(chan []byte)

var fd_map = map[int]int{}


func push(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()

		if fd_map[mt] == 0 {
			fd_map[mt]=mt
		}
		//log.Println(fd_map)

		if err != nil {
			log.Println("read push ", err)
			break
		}
		channel_msg<-message

		/*go func(msg []byte){
			channel_msg<-msg
		}(message)*/
		//log.Printf("recv push %s", mt)
		//err = c.WriteMessage(mt, message)
		//if err != nil {
		//	log.Println("write:", err)
		//	break
		//}
	}

	 func(){

	}()

}

func pull(w http.ResponseWriter, r *http.Request){
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for{
		msg := <-channel_msg
		//log.Println(fd_map)

		for mt, v := range fd_map{
			//w, err := c.NextWriter(mt)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write pull ", err, v)
				break
			}
			//log.Println(msg)
			//w.Write(msg)
		}


		/*go func(mt int, msg []byte){
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write pull", err)
				return
			}
		}(mt, msg)*/
	}

}

func mc(w http.ResponseWriter, r *http.Request) {
	mcTemplate.Execute(w, "ws://"+r.Host+"/push")
}
func guest(w http.ResponseWriter, r *http.Request) {
	guestTemplate.Execute(w, "ws://"+r.Host+"/pull")
}

func main() {


	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/mc", mc)
	http.HandleFunc("/push", push)
	http.HandleFunc("/pull", pull)
	http.HandleFunc("/guest", guest)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
var guestTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>直播游客界面</title>
</head>
<body>
<img id="receiver" style='width:640px;height:480px'/>
<script type="text/javascript" charset="utf-8">
var ws = new WebSocket("{{.}}");
var image = document.getElementById('receiver');
ws.onopen = function(){};
ws.onmessage = function(data)
{
	console.log(data);
    image.src=data.data;
}
</script>
</body>
</html>
`))
var mcTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>主播录制界面</title>
</head>
<body>
<video id="video" autoplay="" style='width:640px;height:480px'></video>
<canvas id="output" style="display:none"></canvas>
<script type="text/javascript" charset="utf-8">
    var ws = new WebSocket("{{.}}");
    var back = document.getElementById('output');
    var backcontext = back.getContext('2d');
    var video = document.getElementById("video");
    var success = function(stream){
        video.src = window.URL.createObjectURL(stream);
    };
    ws.onopen = function(){
        draw();
    };
    var draw = function(){
        try{
            backcontext.drawImage(video,0,0, back.width, back.height);
        }catch(e){
            if (e.name == "NS_ERROR_NOT_AVAILABLE") {
                return setTimeout(draw, 100);
            } else {
                throw e;
            }
        }
        if(video.src){
            ws.send(back.toDataURL("image/jpeg", 0.6));//传递90%清晰度, https://developer.mozilla.org/zh-CN/docs/Web/API/HTMLCanvasElement/toDataURL
        }
        setTimeout(draw, 100);
    };
    navigator.getUserMedia = navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia || navigator.msGetUserMedia;
    navigator.getUserMedia({video:true, audio:false}, success, console.log);
</script>
</body>
</html>
`))
