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
      upers: [],
      currIdx: 0,
      currUperWorks: [],
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
      return "http://localhost:8080/"
    },

    getCurrUperInfo() {
      var resp = this.upers.length > this.currIdx ? this.upers[this.currIdx] : {}
      console.log("getCurrUperInfo resp:", resp)
      return resp
    },

    apiGetCurrUperNotes() {
      var uperInfo = this.getCurrUperInfo()
      if(!uperInfo.uid) {
        console.log("apiGetCurrUperNotes failed: no uid")
        return
      }
      let reqURL = this.getHost()+"uper_notes?uid=" + uperInfo.uid
      axios.get(reqURL).then(resp=>{
        if(!resp.data) {
          console.log("get uper notes failed.")
          return
        }
        if(!resp.data.data){
          showToast("empty notes")
          return
        }

        resp.data.data.forEach(v=>{
          this.currUperWorks.push(v)
        })

        console.log("apiGetCurrUperNotes result:", this.currUperWorks)

      })
    },

    apiGetUpers: function() {
      let reqURL = this.getHost()+"upers?limit=10&token="+this.getUpersNextToken
      axios.get(reqURL).then(resp=>{
        if(!resp.data) {
          console.log("get upers failed.")
          return
        }
        if(!resp.data.data){
          showToast("empty uper")
          return
        }

        resp.data.data.forEach(v=>{
          this.upers.push(v)
        })

        console.log("apiGetUpers upers:", this.upers)

        if(resp.data.next_token){
          this.getUpersNextToken = resp.data.next_token
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
  color: whitesmoke;
}

</style>