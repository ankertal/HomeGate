(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-2d0efcba"],{"9a40":function(t,e,s){"use strict";s.r(e);var a=function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("div",{staticClass:"container"},[s("header",{staticClass:"jumbotron"},[s("h3",[t._v("\n      HomeGate's "),s("strong",[t._v(t._s(t.currentUser.username))]),t._v(" Settings:\n    ")])]),t.successful?t._e():s("div",[t.content?s("div",{staticClass:"alert",class:t.successful?"alert-success":"alert-danger"},[t._v("\n      "+t._s(t.content)+"\n    ")]):t._e()]),t.successful?s("div",[s("strong",[t._v("Gates:")]),s("ul",t._l(t.content.gates,(function(e,a){return s("li",{key:a},[s("nav",{staticClass:"navbar navbar-expand-lg navbar-light bg-light"},[s("a",{staticClass:"navbar-brand",attrs:{href:"#"}},[t.isMyGate(a)?s("span",{staticStyle:{color:"red"}},[t._v(t._s(e))]):t._e(),t.isMyGate(a)?t._e():s("span",{staticStyle:{color:"blue"}},[t._v(t._s(e))])]),t._m(0,!0),s("div",{staticClass:"collapse navbar-collapse",attrs:{id:"navbarNav"}},[s("ul",{staticClass:"navbar-nav"},[s("li",{staticClass:"nav-item"},[s("b-button",{attrs:{pill:"",variant:"outline-success",size:"sm"},on:{click:function(e){return t.gateCommand(a,"is_open")}}},[t._v("Open")])],1),t._v("\n                \n              "),s("li",{staticClass:"nav-item"},[s("b-button",{attrs:{pill:"",variant:"outline-warning",size:"sm"},on:{click:function(e){return t.gateCommand(a,"is_close")}}},[t._v("Close")])],1),t._v("\n                \n              "),s("li",{staticClass:"nav-item active"},[s("b-button",{attrs:{pill:"",variant:"outline-danger",size:"sm",disabled:t.isMyGate(a)},on:{click:function(e){return t.deleteGate(a)}}},[t._v("\n                  Delete\n                ")])],1),t._v("\n                  \n              "),t.isMyGate(a)?s("li",{staticClass:"nav-item active"},[s("b-button",{attrs:{pill:"",variant:"outline-dark",size:"sm"},on:{click:function(e){return t.addUser(a)}}},[t._v("Add User")]),t._v(" \n                "),s("input",{directives:[{name:"model",rawName:"v-model",value:t.email,expression:"email"}],staticClass:"text-left",attrs:{size:"sm",type:"text",required:""},domProps:{value:t.email},on:{input:function(e){e.target.composing||(t.email=e.target.value)}}}),t.msg.email?s("span",[t._v(t._s(t.msg.email))]):t._e()],1):t._e()])])])])})),0),s("div",{attrs:{id:"gate-list"}},[s("form",{on:{submit:function(e){return e.preventDefault(),t.addGate.apply(null,arguments)}}},[s("label",{attrs:{for:"new-gate"}},[t._v("Add a gate to "+t._s(t.currentUser.username)+": ")]),t._v("\n          \n        "),s("input",{directives:[{name:"model",rawName:"v-model",value:t.newGateText,expression:"newGateText"}],attrs:{id:"new-gate",placeholder:"e.g. gate-XXXXXX"},domProps:{value:t.newGateText},on:{input:function(e){e.target.composing||(t.newGateText=e.target.value)}}}),s("button",[t._v("Add")])])])]):t._e(),s("br"),s("br"),s("br"),t.successful?s("div",[s("strong",[t._v("My ("+t._s(t.currentUser.my_gate)+") friends: ")]),s("ul",{attrs:{id:"uses-list"}},t._l(this.content.users,(function(e,a){return s("li",[t._v("\n        "+t._s(e)+"\n      ")])})),0)]):t._e(),s("br"),s("br"),s("br"),this.showCommandStatus?s("div",[s("b-alert",{attrs:{show:t.dismissCountDown,variant:t.alertVariant,dismissible:""},on:{"dismiss-count-down":t.countDownChanged}},[t._v("\n      "+t._s(this.message))])],1):t._e()])},n=[function(){var t=this,e=t.$createElement,s=t._self._c||e;return s("button",{staticClass:"navbar-toggler",attrs:{type:"button","data-toggle":"collapse","data-target":"#navbarNav","aria-controls":"navbarNav","aria-expanded":"false","aria-label":"Toggle navigation"}},[s("span",{staticClass:"navbar-toggler-icon"})])}],i=(s("6b54"),s("1f57")),r={name:"User",computed:{currentUser:function(){return this.$store.state.auth.user}},data:function(){return{content:"",message:"",cmdError:!1,showCommandStatus:!1,successful:!1,newGateText:"",newUserText:"",msg:[],email:"",ismissCountDown:null,showDismissibleAlert:!1,alertVariant:"info"}},watch:{email:function(t){this.email=t,this.validateEmail(t)}},mounted:function(){var t=this;i["a"].getUserBoard().then((function(e){t.content=e.data,t.successful=!0,t.items=t.content.gates}),(function(e){t.successful=!1,t.content=e.response&&e.response.data&&e.response.data.message||e.message||e.toString()}))},methods:{addGate:function(){""!=this.newGateText&&(this.content.gates.push(this.newGateText),this.newGateText="")},deleteGate:function(t){var e=this.content.gates[t],s=this.content.user_gate;e!=s&&this.content.gates.splice(t,1)},isMyGate:function(t){var e=this.content.gates[t],s=this.content.my_gate;return e===s},gateCommand:function(t,e){var s=this;this.message="",this.cmdError=!1;var a=this.content.gates[t];i["a"].triggerCommand(this.currentUser,a,e).then((function(t){t.isAxiosError||t.is_error?(s.message=t.message,s.showCommandStatus=!0,s.showAlert("danger")):(s.cmdError=t.is_error,s.message=t.message,s.showCommandStatus=!s.cmdError,s.showAlert("info"))}),(function(t){s.cmdError=!0,s.message=t.response&&t.response.data&&t.response.data.message||t.message||t.toString(),s.showAlert("danger")}))},addUser:function(){/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(this.email)&&(this.content.users.push(this.email),this.newUserText="")},validateEmail:function(t){/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(t)?this.msg["email"]="":this.msg["email"]="Invalid Email Address"},countDownChanged:function(t){this.dismissCountDown=t},showAlert:function(t){this.dismissCountDown=2,this.alertVariant=t}}},o=r,l=s("2877"),c=Object(l["a"])(o,a,n,!1,null,null,null);e["default"]=c.exports}}]);
//# sourceMappingURL=chunk-2d0efcba.c6441f82.js.map