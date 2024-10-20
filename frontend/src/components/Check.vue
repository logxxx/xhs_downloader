<template>
  <div>
    <div id="main_page">
      <div class="uper_info" style="position: fixed; left: 10px; top: 0px;background-color: rgba(0, 0, 0, 0.5);">

        <div style="font-size: 20px;font-weight: bold;margin-top: 10px;margin-bottom: 5px;">{{getCurrUperInfo().name}}</div>
        <div>{{getCurrUperInfo().desc}}</div>

      </div>

      <div style="height:100px"></div>

      <div class="show_notes">
        <div class="per_note_poster" v-for="(note, idx) in this.currUperWorks" :key="note.note_id">
            <van-image radius="5" fit="cover" style="width: 160px;height: 240px;" :src="getFileURL(note.poster)"/>
         <div>{{note.title}}</div>
        </div>
      </div>

      <div style="height:400px"></div>

      <div v-show="showAct" style="position:fixed;bottom:50px;display:flex;">
        <van-space direction="vertical" fill style="margin:0 10px">
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[0].value)'>{{ this.tagTmpl[0].show }}</van-button>
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[3].value)'>{{ this.tagTmpl[3].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[6].value)'>{{ this.tagTmpl[6].show }}</van-button>
          <van-button class="btn" type="success" @click="addTag(this.tagTmpl[9].value)">{{ this.tagTmpl[9].show }}</van-button>
          <van-button class="btn" type="primary" @click="addTag(this.tagTmpl[12].value)">{{ this.tagTmpl[12].show }}</van-button>
        </van-space>
        <van-space direction="vertical" fill style="">
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[1].value)'>{{ this.tagTmpl[1].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[4].value)'>{{ this.tagTmpl[4].show }}</van-button>
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[7].value)'>{{ this.tagTmpl[7].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[10].value)'>{{ this.tagTmpl[10].show }}</van-button>
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[13].value)'>{{ this.tagTmpl[13].show }}</van-button>
        </van-space>

        <van-space direction="vertical" fill style="margin-left:10px;">
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[2].value)'>{{ this.tagTmpl[2].show }}</van-button>
          <van-button class="btn" type="warning" @click='addTag(this.tagTmpl[5].value)'>{{ this.tagTmpl[5].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[8].value)'>{{ this.tagTmpl[8].show }}</van-button>
          <van-button class="btn" type="success" @click='addTag(this.tagTmpl[11].value)'>{{ this.tagTmpl[11].show }}</van-button>
          <van-button class="btn" type="primary" @click='addTag(this.tagTmpl[14].value)'>{{ this.tagTmpl[14].show }}</van-button>
        </van-space>

      </div>

    </div>
  </div>
</template>
<script>



import axios from "axios";
import {showToast} from "vant";
import {showFailToast} from "vant/lib/toast/function-call";

export default {
  name: 'UperInfo',
  props: {
    msg: String
  },
  data() {
    return {
      getUpersNextToken: '',
      currUper: {},
      currUperWorks: [],
      showAct: true,

      tagTmpl:[

        {show:'颜值',value:'yanzhi'}, {show:'我看',value:'wokan'}, {show:'难评',value:'nanping'},
        {show:'正常',value:'zhengchang'},{show:'营销',value:'yingxiao'}, {show:'超高',value:'high'},
        {show:'jio',value:'jio'}, {show:'jio商',value:'jioshang'},{show:'幼态',value:'youtai'},
        {show:'可1',value:'ke1'},{show:'可2',value:'ke2'}, {show:'可3', value:'ke3'},
        {show:'一般', value:'yiban'}, {show:'商家',value:'shangjia'},{show:'没用',value:'meiyong'},
      ],
    }
  },
  mounted(){

  },
  created(){
    this.apiGetUpers()
  },
  methods: {

    getHost: function() {
      //return "http://testnas.com:9887/"
      //return "http://localhost:6080/"
      return ""
    },

    getCurrUperInfo() {
      return this.currUper ? this.currUper : {}
    },

    addTag: function(tag){
      let reqURL = this.getHost()+"add_tag?tag="+tag+"&uid="+this.currUper.uid

      console.log("["+tag+"]"+reqURL)
      axios.get(reqURL).then(resp=>{
        if(resp.data.err_msg){
          showFailToast(resp.data.err_msg)
        }
      })

      this.apiGetUpers()
    },

    apiGetUpers: function() {
      let reqURL = this.getHost()+"upers?with=withNoTag&limit=1&token="+this.getUpersNextToken
      axios.get(reqURL).then(resp=>{
        if(!resp.data) {
          console.log("get upers failed.")
          return
        }
        if(!resp.data.data){
          showToast("empty uper")
          return
        }

        this.currUper = resp.data.data[0]

        console.log("apiGetUpers currUper:", this.currUper)

        if(resp.data.token){
          this.getUpersNextToken = resp.data.token
        }

        this.apiGetCurrUperNotes()

      })
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
  overflow: hidden;
}

.show_images {
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
  color: whitesmoke;
}

.btn{
  width:100px;
}

</style>