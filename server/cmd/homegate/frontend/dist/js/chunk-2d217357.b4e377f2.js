(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-2d217357"],{c66d:function(e,n,r){"use strict";r.r(n);var s=function(){var e=this,n=e.$createElement,r=e._self._c||n;return r("div",{staticClass:"container"},[r("header",{staticClass:"jumbotron"},[r("h3",[e._v("\n      HomeGate "),r("strong",[e._v(e._s(e.currentUser.username))]),e._v(" Profile\n    ")])]),r("p",[r("strong",[e._v("Token:")]),e._v("\n    "+e._s(e.currentUser.accessToken.substring(0,20))+" ...\n    "+e._s(e.currentUser.accessToken.substr(e.currentUser.accessToken.length-20))+"\n  ")]),r("p",[r("strong",[e._v("Id:")]),e._v("\n    "+e._s(e.currentUser.id)+"\n  ")]),r("p",[r("strong",[e._v("Email:")]),e._v("\n    "+e._s(e.currentUser.email)+"\n  ")])])},t=[],u={name:"Profile",computed:{currentUser:function(){return this.$store.state.auth.user}},mounted:function(){this.currentUser||this.$router.push("/login")}},c=u,o=r("2877"),a=Object(o["a"])(c,s,t,!1,null,null,null);n["default"]=a.exports}}]);
//# sourceMappingURL=chunk-2d217357.b4e377f2.js.map