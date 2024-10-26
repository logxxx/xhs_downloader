<template>
  <div>
    <div id="main_page">

      <van-overlay style="" :show="isShowUserNotesDialog" @click="isShowUserNotesDialog=false">

        <div style="color:whitesmoke;font-size: 30px;font-weight: bold;margin-top: 10px;margin-bottom: 5px;">{{getCurrUperInfo().name}}</div>
        <div style="color:whitesmoke;">{{getCurrUperInfo().desc}}</div>
        <div style="color:whitesmoke;">{{getCurrUperInfo().tags}}</div>
        <div style="color:whitesmoke;">{{getCurrUperInfo().myTags}}</div>

        <div style="display: flex; flex-wrap: wrap;height: 100vh;width:100vw;overflow-y: scroll">



        <div class="per_note_poster" v-for="(note, idx) in this.currUperNotes" :key="note.note_id">
          <van-image radius="5" fit="cover" style="max-width: 160px;max-height: 240px;" :src="getFileURL(note.poster)">
            <template slot="error">加载失败</template>
          </van-image>
          <div style="color: whitesmoke">{{note.title}}</div>
        </div>
        </div>
      </van-overlay>

      <div v-if="showAct" class="uper_info" style="width: 100%;position: fixed; left: 10px; top: 0px;background-color: rgba(0, 0, 0, 0.5);">

        <div style="font-size: 20px;font-weight: bold;margin-top: 10px;margin-bottom: 5px;">{{getCurrUperInfo().name}}</div>
        <!--<div style="width:90%;margin:5px;word-wrap:break-word;color:whitesmoke">{{this.parseB64(this.currNote.video)}}</div>-->
        <div style="color:whitesmoke"></div>
        <div style="color:whitesmoke">{{this.nextToken}}&nbsp;{{this.currNote.title}}</div>
        <div style="color:whitesmoke">{{this.currNote.show_size}}</div>

      </div>

      <video
          id="only_video"
          v-if="this.currNote.video && this.currNote.video.length>0"
          :src="getFileURL(this.currNote.video)"
          @click="onClickVideo"
          class="videoSource"

          loop="loop" autoplay="autoplay"
          webkit-playsinline="true" x-webkit-airplay="true" playsInline={true} x5-playsinline="true" x5-video-orientation="portraint"
      >
      </video>

      <div v-if="this.currNote.images && this.currNote.images.length>0" class="show_images">
        <div style="height:100px"></div>
        <div class="per_note_poster" v-for="(img, idx) in this.currNote.images" :key="img">
            <van-image radius="5" fit="cover" style="width: 160px;height: 240px;" :src="getFileURL(img)"/>
         <div>{{this.currNote.title}}</div>
        </div>
      </div>

      <!--
      <div v-if="showAct" style=" background-color: rgba(0, 0, 0, 0.5); padding: 10px; margin: 10px;position:fixed;bottom:450px;display:flex;">

      </div>
      -->


      <div v-if="showAct" style="position:fixed;bottom:100px;margin-left: 10px;">
        <van-checkbox style="margin-bottom:20px;background-color: rgba(0, 0, 0, 0.5);"  v-model="isAddTagForUperAndAllNote"><div style="color:whitesmoke;">应用到ta的所有笔记</div></van-checkbox>

      <div style="display:flex;">

        <van-space direction="vertical" fill style="">
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[0].value)'>{{ this.tagTmpl[0].show }}</van-button>
          <van-button class="btn" type="warning" @click='addTag(this.tagTmpl[3].value)'>{{ this.tagTmpl[3].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[6].value)'>{{ this.tagTmpl[6].show }}</van-button>
          <van-button class="btn" type="success" @click="addTag(this.tagTmpl[9].value)">{{ this.tagTmpl[9].show }}</van-button>
          <van-button class="btn" type="primary" @click="addTag(this.tagTmpl[12].value)">{{ this.tagTmpl[12].show }}</van-button>
          <van-button class="btn" type="warning" @click="addTag(this.tagTmpl[15].value)">{{ this.tagTmpl[15].show }}</van-button>
        </van-space>
        <van-space direction="vertical" fill style="margin-left: 20px;">
          <van-button class="btn" type="warning" @click='addTag(this.tagTmpl[1].value)'>{{ this.tagTmpl[1].show }}</van-button>
          <van-button class="btn" type="danger" @click='addTag(this.tagTmpl[4].value)'>{{ this.tagTmpl[4].show }}</van-button>
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[7].value)'>{{ this.tagTmpl[7].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[10].value)'>{{ this.tagTmpl[10].show }}</van-button>
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[13].value)'>{{ this.tagTmpl[13].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[16].value)'>{{ this.tagTmpl[16].show }}</van-button>
        </van-space>

        <van-space direction="vertical" fill style="margin-left: 20px;">
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[2].value)'>{{ this.tagTmpl[2].show }}</van-button>
          <van-button class="btn" type="warning" @click='addTag(this.tagTmpl[5].value)'>{{ this.tagTmpl[5].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[8].value)'>{{ this.tagTmpl[8].show }}</van-button>
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[11].value)'>{{ this.tagTmpl[11].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[14].value)'>{{ this.tagTmpl[14].show }}</van-button>
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[17].value)'>{{ this.tagTmpl[17].show }}</van-button>
        </van-space>
      </div>
        </div>

      <van-floating-bubble v-if="showPlayBubble" @click="onPlayVideo" icon="arrow"/>

      <van-floating-bubble v-model:offset="showUserNotesDialogOffset" @click="showUserNotesDialog" icon="user"/>

    </div>
  </div>
</template>
<script>



import axios from "axios";
import {showToast} from "vant";
import {showFailToast} from "vant/lib/toast/function-call";

export default {
  name: 'EachNote',
  props: {
    msg: String
  },
  data() {
    return {
      nextToken: '',
      currNote: {},
      currUper: {},
      currUperNotes: [],
      showAct: true,
      showPlayBubble: true,
      isShowUserNotesDialog:false,
      showUserNotesDialogOffset: {x: 20, y: 80},
      isAddTagForUperAndAllNote: false,

      tagTmpl:[
        {show:'颜值',value:'yanzhi'}, {show:'我看',value:'wokan'}, {show:'内内',value:'neinei'},
        {show:'自拍',value:'zipai'},{show:'没用',value:'meiyong'}, {show:'幼态',value:'youtai'},
        {show:'jio',value:'jio'}, {show:'ss',value:'ss'},{show:'超高',value:'high'},
        {show:'还行',value:'haixing'},{show:'跳舞',value:'tiaowu'}, {show:'摆拍', value:'baipai'},
        {show:'高级',value:'gaoji'},{show:'挺好',value:'tinghao'}, {show:'待定', value:'daiding'},
        {show:'营销',value:'yingxiao'},{show:'剧情', value:'juqing'}, {show:'日常',value:'richang'},
      ],
    }
  },
  mounted(){

  },
  created(){
    this.apiGetOneNote()
  },
  methods: {

    parseB64: function(input) {
      let resp = atob(input)
      console.log(resp)
      return resp
    },

    onClickVideo: function() {
      console.log("onClickVideo showAct:", this.showAct)
      this.showAct = !this.showAct

    },

    getHost: function() {
      //return "http://testnas.com:9887/"
      //return "http://localhost:6080/"
      //return ""
      return ""
    },

    getCurrUperInfo() {

      return this.currUper ? this.currUper : {}
    },

    addTag: function(tag){

      showToast({
        message: ""+tag+"+1",
        duration: 500,
      })

      let reqURL = ""
      if(tag == "meiyong") {
        reqURL = this.getHost()+"note/delete?note_id="+this.currNote.note_id
        if(this.isAddTagForUperAndAllNote) {
          reqURL = this.getHost()+"uper/delete?uid="+this.currNote.uper_uid
        }
      } else {
        reqURL = this.getHost()+"note/add_tag?tag="+tag+"&note_id="+this.currNote.note_id
        if(this.isAddTagForUperAndAllNote) {
          reqURL = this.getHost()+"uper/add_tag?tag="+tag+"&uid="+this.currNote.uper_uid
        }
      }

      this.isAddTagForUperAndAllNote = false

      console.log("["+tag+"]"+reqURL)
      axios.get(reqURL).then(resp=>{
        if(resp.data.err_msg){
          showFailToast(resp.data.err_msg)
        }

        this.apiGetOneNote(()=>{
          this.showAct = false
          setTimeout(()=>{
            this.showAct = true
          }, 2000)
        })
      }).catch(err=>{
        this.apiGetOneNote()
      })

    },

    showUserNotesDialog:function(){

      if(this.isShowUserNotesDialog == false) {
        this.apiGetUperInfo(this.currNote.uper_uid)
        this.apiGetUperNotes(this.currNote.uper_uid)

        this.isShowUserNotesDialog = true
      } else {
        this.isShowUserNotesDialog = false
      }

    },

    apiGetUperNotes: function(uid){
      let reqURL = this.getHost()+"uper_notes?uid=" + uid
      axios.get(reqURL).then(resp=> {
        if (!resp.data) {
          console.log("apiGetUperNotes failed.")
          return
        }
        if (!resp.data.data) {
          showToast("apiGetUperNotes empty")
          return
        }

        this.currUperNotes = resp.data.data
        this.currUper.myTags = resp.data.tags

        console.log("set currUperNotes:", this.currUperNotes)
      })
    },

    apiGetUperInfo:function(uid) {
      let reqURL = this.getHost()+"uper?with=withoutNotes&uid=" + uid
      axios.get(reqURL).then(resp=> {
        if (!resp.data) {
          console.log("get uper failed.")
          return
        }
        if (!resp.data) {
          showToast("empty uper")
          return
        }

        this.currUper = resp.data
      })
    },

    apiGetOneNote(afterFn) {
      //let reqURL = this.getHost()+"one_note?type=video&token="+this.nextToken
      let reqURL = this.getHost()+"one_video?type=video&token="+this.nextToken
      axios.get(reqURL).then(resp=>{
        if(!resp.data) {
          console.log("get one_note failed.")
          return
        }
        if(!resp.data.data){
          showToast("empty one_note")
          return
        }

        this.currNote = resp.data.data

        console.log("apiGetOneNote currNote:", this.currNote)

        this.parseB64(this.currNote.video)

        if(resp.data.token){
          this.nextToken = resp.data.token
        }

        let videoDom = document.getElementById("only_video")
        if(videoDom){
          console.log("play video start")
          videoDom.play()
        } else {
          console.log("play video failed: not found")
        }

        if(afterFn) {
          afterFn()
        }


      })
    },

    onPlayVideo: function() {

      this.showPlayBubble = false

      let videoDom = document.getElementById("only_video")
      if(videoDom){
        console.log("play video start")
        videoDom.play()
      } else {
        console.log("play video failed: not found")
      }

    },

    getFileURL: function(id) {
      if(!id) {
        return
      }
      var resp= this.getHost()+"file?id=" + id
      //console.log("getFileURL resp:", resp)
      return resp
    },

  }
}
</script>

<style>

#main_page {
  background-color: black;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  display:flex;
}

.show_notes {
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  padding: 0 5px;
}


.per_note_poster {
  width: calc(45% - 10px); /* 假设图片之间有20px间隔，每张图片的宽度为50%减去每侧的间隔 */
  margin: 10px; /* 设置图片间隔 */
  background-color: rgba(0, 0, 0, 0.5);
}

.img_container {
  width: 160px;
  height: 240px;
}

.img_container img {

}

div {

}

.btn{
  width:100px;
  margin-bottom: 10px;
}

.videoSource{
  width: 100vw;
}


</style>