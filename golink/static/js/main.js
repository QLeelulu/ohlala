(function(){

/* hack seajs, 避免与页面上的seajs冲突 */
var _seajsHost = {location: window.location};
(function(){
    var seajs=
/*
 SeaJS - A Module Loader for the Web
 v1.2.0 | seajs.org | MIT Licensed
*/
this.seajs={_seajs:this.seajs};seajs.version="1.2.0";seajs._util={};seajs._config={debug:"",preload:[]};
(function(a){var c=Object.prototype.toString,d=Array.prototype;a.isString=function(a){return"[object String]"===c.call(a)};a.isFunction=function(a){return"[object Function]"===c.call(a)};a.isRegExp=function(a){return"[object RegExp]"===c.call(a)};a.isObject=function(a){return a===Object(a)};a.isArray=Array.isArray||function(a){return"[object Array]"===c.call(a)};a.indexOf=d.indexOf?function(a,c){return a.indexOf(c)}:function(a,c){for(var b=0;b<a.length;b++)if(a[b]===c)return b;return-1};var b=a.forEach=
d.forEach?function(a,c){a.forEach(c)}:function(a,c){for(var b=0;b<a.length;b++)c(a[b],b,a)};a.map=d.map?function(a,c){return a.map(c)}:function(a,c){var d=[];b(a,function(a,b,f){d.push(c(a,b,f))});return d};a.filter=d.filter?function(a,c){return a.filter(c)}:function(a,c){var d=[];b(a,function(a,b,f){c(a,b,f)&&d.push(a)});return d};a.unique=function(a){var c=[],d={};b(a,function(a){d[a]=1});if(Object.keys)c=Object.keys(d);else for(var i in d)d.hasOwnProperty(i)&&c.push(i);return c};a.keys=Object.keys;
a.keys||(a.keys=function(a){var c=[],b;for(b in a)a.hasOwnProperty(b)&&c.push(b);return c});a.now=Date.now||function(){return(new Date).getTime()}})(seajs._util);(function(a,c){var d=Array.prototype;a.log=function(){if("undefined"!==typeof console){var a=d.slice.call(arguments),f="log";console[a[a.length-1]]&&(f=a.pop());if("log"!==f||c.debug)a="dir"===f?a[0]:d.join.call(a," "),console[f](a)}}})(seajs._util,seajs._config);
(function(a,c,d){function b(a){a=a.match(s);return(a?a[0]:".")+"/"}function f(a){k.lastIndex=0;k.test(a)&&(a=a.replace(k,"$1/"));if(-1===a.indexOf("."))return a;for(var c=a.split("/"),b=[],d,e=0;e<c.length;e++)if(d=c[e],".."===d){if(0===b.length)throw Error("The path is invalid: "+a);b.pop()}else"."!==d&&b.push(d);return b.join("/")}function o(a){var a=f(a),c=a.charAt(a.length-1);if("/"===c)return a;"#"===c?a=a.slice(0,-1):-1===a.indexOf("?")&&!e.test(a)&&(a+=".js");0<a.indexOf(":80/")&&(a=a.replace(":80/",
"/"));return a}function l(a){if("#"===a.charAt(0))return a.substring(1);var b=c.alias;if(b&&j(a)){var d=a.split("/"),e=d[0];b.hasOwnProperty(e)&&(d[0]=b[e],a=d.join("/"))}return a}function i(a){return 0<a.indexOf("://")||0===a.indexOf("//")}function j(a){var c=a.charAt(0);return-1===a.indexOf("://")&&"."!==c&&"/"!==c}var s=/.*(?=\/.*$)/,k=/([^:\/])\/\/+/g,e=/\.(?:css|js)$/,n=/^(.*?\w)(?:\/|$)/,g={},d=d.location,h=d.protocol+"//"+d.host+function(a){"/"!==a.charAt(0)&&(a="/"+a);return a}(d.pathname);
0<h.indexOf("\\")&&(h=h.replace(/\\/g,"/"));a.dirname=b;a.realpath=f;a.normalize=o;a.parseAlias=l;a.parseMap=function(b){var d=c.map||[];if(!d.length)return b;for(var e=b,h=0;h<d.length;h++){var f=d[h];if(a.isArray(f)&&2===f.length){var i=f[0];if(a.isString(i)&&-1<e.indexOf(i)||a.isRegExp(i)&&i.test(e))e=e.replace(i,f[1])}else a.isFunction(f)&&(e=f(e))}e!==b&&(g[e]=b);return e};a.unParseMap=function(a){return g[a]||a};a.id2Uri=function(a,d){if(!a)return"";a=l(a);d||(d=h);var e;i(a)?e=a:0===a.indexOf("./")||
0===a.indexOf("../")?(0===a.indexOf("./")&&(a=a.substring(2)),e=b(d)+a):e="/"===a.charAt(0)&&"/"!==a.charAt(1)?d.match(n)[1]+a:c.base+"/"+a;return o(e)};a.isAbsolute=i;a.isTopLevel=j;a.pageUri=h})(seajs._util,seajs._config,this);
(function(a,c){function d(a,b){a.onload=a.onerror=a.onreadystatechange=function(){k.test(a.readyState)&&(a.onload=a.onerror=a.onreadystatechange=null,a.parentNode&&!c.debug&&i.removeChild(a),a=void 0,b())}}function b(c,b){h||p?(a.log("Start poll to fetch css"),setTimeout(function(){f(c,b)},1)):c.onload=c.onerror=function(){c.onload=c.onerror=null;c=void 0;b()}}function f(a,c){var b;if(h)a.sheet&&(b=!0);else if(a.sheet)try{a.sheet.cssRules&&(b=!0)}catch(d){"NS_ERROR_DOM_SECURITY_ERR"===d.name&&(b=
!0)}setTimeout(function(){b?c():f(a,c)},1)}function o(){}var l=document,i=l.head||l.getElementsByTagName("head")[0]||l.documentElement,j=i.getElementsByTagName("base")[0],s=/\.css(?:\?|$)/i,k=/loaded|complete|undefined/,e,n;a.fetch=function(c,f,h){var k=s.test(c),g=document.createElement(k?"link":"script");h&&(h=a.isFunction(h)?h(c):h)&&(g.charset=h);f=f||o;"SCRIPT"===g.nodeName?d(g,f):b(g,f);k?(g.rel="stylesheet",g.href=c):(g.async="async",g.src=c);e=g;j?i.insertBefore(g,j):i.appendChild(g);e=null};
a.getCurrentScript=function(){if(e)return e;if(n&&"interactive"===n.readyState)return n;for(var a=i.getElementsByTagName("script"),c=0;c<a.length;c++){var b=a[c];if("interactive"===b.readyState)return n=b}};a.getScriptAbsoluteSrc=function(a){return a.hasAttribute?a.src:a.getAttribute("src",4)};a.importStyle=function(a,c){if(!c||!l.getElementById(c)){var b=l.createElement("style");c&&(b.id=c);i.appendChild(b);b.styleSheet?b.styleSheet.cssText=a:b.appendChild(l.createTextNode(a))}};var g=navigator.userAgent,
h=536>Number(g.replace(/.*AppleWebKit\/(\d+)\..*/,"$1")),p=0<g.indexOf("Firefox")&&!("onload"in document.createElement("link"))})(seajs._util,seajs._config,this);(function(a){var c=/(?:^|[^.$])\brequire\s*\(\s*(["'])([^"'\s\)]+)\1\s*\)/g;a.parseDependencies=function(d){var b=[],f,d=d.replace(/^\s*\/\*[\s\S]*?\*\/\s*$/mg,"").replace(/^\s*\/\/.*$/mg,"");for(c.lastIndex=0;f=c.exec(d);)f[2]&&b.push(f[2]);return a.unique(b)}})(seajs._util);
(function(a,c,d){function b(a,c){this.uri=a;this.status=c||0}function f(a,d){return c.isString(a)?b._resolve(a,d):c.map(a,function(a){return f(a,d)})}function o(a,r){var m=c.parseMap(a);x[m]?(e[a]=e[m],r()):p[m]?t[m].push(r):(p[m]=!0,t[m]=[r],b._fetch(m,function(){x[m]=!0;var b=e[a];b.status===h.FETCHING&&(b.status=h.FETCHED);u&&(l(a,u),u=null);q&&b.status===h.FETCHED&&(e[a]=q,q.packageUri=a);q=null;p[m]&&delete p[m];t[m]&&(c.forEach(t[m],function(a){a()}),delete t[m])},d.charset))}function l(a,d){var m=
e[a]||(e[a]=new b(a));m.status<h.SAVED&&(m.id=d.id||a,m.dependencies=f(c.filter(d.dependencies||[],function(a){return!!a}),a),m.factory=d.factory,m.status=h.SAVED);return m}function i(a,c){var b=a(c.require,c.exports,c);void 0!==b&&(c.exports=b)}function j(a){var b=a.uri,d=n[b];d&&(c.forEach(d,function(c){i(c,a)}),delete n[b])}function s(a){var b=a.uri;return c.filter(a.dependencies,function(a){v=[b];if(a=k(e[a],b))v.push(b),c.log("Found circular dependencies:",v.join(" --\> "),void 0);return!a})}
function k(a,b){if(!a||a.status!==h.SAVED)return!1;v.push(a.uri);var d=a.dependencies;if(d.length){if(-1<c.indexOf(d,b))return!0;for(var f=0;f<d.length;f++)if(k(e[d[f]],b))return!0}return!1}var e={},n={},g=[],h={FETCHING:1,FETCHED:2,SAVED:3,READY:4,COMPILING:5,COMPILED:6};b.prototype._use=function(a,b){c.isString(a)&&(a=[a]);var d=f(a,this.uri);this._load(d,function(){var a=c.map(d,function(a){return a?e[a]._compile():null});b&&b.apply(null,a)})};b.prototype._load=function(a,d){function f(a){a&&(a.status=
h.READY);0===--g&&d()}var y=c.filter(a,function(a){return a&&(!e[a]||e[a].status<h.READY)}),i=y.length;if(0===i)d();else for(var g=i,j=0;j<i;j++)(function(a){function c(){d=e[a];if(d.status>=h.SAVED){var r=s(d);r.length?b.prototype._load(r,function(){f(d)}):f(d)}else f()}var d=e[a]||(e[a]=new b(a,h.FETCHING));d.status>=h.FETCHED?c():o(a,c)})(y[j])};b.prototype._compile=function(){function a(c){c=f(c,b.uri);c=e[c];if(!c)return null;if(c.status===h.COMPILING)return c.exports;c.parent=b;return c._compile()}
var b=this;if(b.status===h.COMPILED)return b.exports;if(b.status<h.READY)return null;b.status=h.COMPILING;a.async=function(a,c){b._use(a,c)};a.resolve=function(a){return f(a,b.uri)};a.cache=e;b.require=a;b.exports={};var d=b.factory;c.isFunction(d)?(g.push(b),i(d,b),g.pop()):void 0!==d&&(b.exports=d);b.status=h.COMPILED;j(b);return b.exports};b._define=function(a,b,d){var i=arguments.length;1===i?(d=a,a=void 0):2===i&&(d=b,b=void 0,c.isArray(a)&&(b=a,a=void 0));!c.isArray(b)&&c.isFunction(d)&&(b=
c.parseDependencies(d.toString()));var i={id:a,dependencies:b,factory:d},g;if(document.attachEvent){var j=c.getCurrentScript();j&&(g=c.unParseMap(c.getScriptAbsoluteSrc(j)));g||c.log("Failed to derive URI from interactive script for:",d.toString(),"warn")}if(j=a?f(a):g){if(j===g){var k=e[g];k&&(k.packageUri&&k.status===h.SAVED)&&(e[g]=null)}i=l(j,i);if(g){if((e[g]||{}).status===h.FETCHING)e[g]=i,i.packageUri=g}else q||(q=i)}else u=i};b._getCompilingModule=function(){return g[g.length-1]};b._find=
function(a){var b=[];c.forEach(c.keys(e),function(d){if(c.isString(a)&&-1<d.indexOf(a)||c.isRegExp(a)&&a.test(d))d=e[d],d.exports&&b.push(d.exports)});var d=b.length;1===d?b=b[0]:0===d&&(b=null);return b};b._modify=function(b,c){var d=f(b),g=e[d];g&&g.status===h.COMPILED?i(c,g):(n[d]||(n[d]=[]),n[d].push(c));return a};b.STATUS=h;b._resolve=c.id2Uri;b._fetch=c.fetch;b.cache=e;var p={},x={},t={},u=null,q=null,v=[],w=new b(c.pageUri,h.COMPILED);a.use=function(c,b){var e=d.preload;e.length?w._use(e,function(){d.preload=
[];w._use(c,b)}):w._use(c,b);return a};a.define=b._define;a.cache=b.cache;a.find=b._find;a.modify=b._modify;a.pluginSDK={Module:b,util:c,config:d}})(seajs,seajs._util,seajs._config);
(function(a,c,d){var b="seajs-ts="+c.now(),f=document.getElementById("seajsnode");f||(f=document.getElementsByTagName("script"),f=f[f.length-1]);var o=c.getScriptAbsoluteSrc(f)||c.pageUri,o=c.dirname(function(a){if(a.indexOf("??")===-1)return a;var b=a.split("??"),a=b[0],b=c.filter(b[1].split(","),function(a){return a.indexOf("sea.js")!==-1});return a+b[0]}(o));c.loaderDir=o;var l=o.match(/^(.+\/)seajs\/[\d\.]+\/$/);l&&(o=l[1]);d.base=o;if(f=f.getAttribute("data-main"))d.main=f;d.charset="utf-8";
a.config=function(f){for(var j in f)if(f.hasOwnProperty(j)){var l=d[j],k=f[j];if(l&&j==="alias")for(var e in k){if(k.hasOwnProperty(e)){var n=l[e],g=k[e];/^\d+\.\d+\.\d+$/.test(g)&&(g=e+"/"+g+"/"+e);n&&n!==g&&c.log("The alias config is conflicted:","key =",'"'+e+'"',"previous =",'"'+n+'"',"current =",'"'+g+'"',"warn");l[e]=g}}else if(l&&(j==="map"||j==="preload")){c.isString(k)&&(k=[k]);c.forEach(k,function(a){a&&l.push(a)})}else d[j]=k}if((f=d.base)&&!c.isAbsolute(f))d.base=c.id2Uri("./"+f+"/");
if(d.debug===2){d.debug=1;a.config({map:[[/^.*$/,function(a){a.indexOf("seajs-ts=")===-1&&(a=a+((a.indexOf("?")===-1?"?":"&")+b));return a}]]})}if(d.debug)a.debug=!!d.debug;return this};d.debug&&(a.debug=!!d.debug)})(seajs,seajs._util,seajs._config);
(function(a,c,d){a.log=c.log;a.importStyle=c.importStyle;a.config({alias:{seajs:c.loaderDir}});if(-1<d.location.search.indexOf("seajs-debug")||-1<document.cookie.indexOf("seajs=1"))a.config({debug:2}).use("seajs/plugin-debug"),a._use=a.use,a._useArgs=[],a.use=function(){a._useArgs.push(arguments);return a}})(seajs,seajs._util,this);
(function(a,c,d){var b=a._seajs;if(b&&!b.args)d.seajs=a._seajs;else{d.define=a.define;c.main&&a.use(c.main);if(c=(b||0).args){d={"0":"config",1:"use",2:"define"};for(b=0;b<c.length;b+=2)a[d[c[b]]].apply(a,c[b+1])}delete a.define;delete a._util;delete a._config;delete a._seajs}})(seajs,seajs._config,this);

}).apply(_seajsHost);

/**
 * 框架基础函数
 */
function Base () {
}

Base.prototype = {
    toLogin: function () {
        var ru = encodeURIComponent(window.location.pathname);
        window.location.href = 'http://'+window.location.host + '/user/login?returnurl=' + ru;
    },
    /**
     * 格式化字符串 from tbra
     * eg:
     *  formatText('{{0}}天有{{1}}个小时', [1, 24]) 
     *  or
     *  formatText('{{day}}天有{{hour}}个小时', {day:1, hour:24}}
     * @param {Object} msg
     * @param {Object} values
     */
    tpFormat: function(msg, values, filter) {
        var pattern = /\{\{([\w\s\.\(\)"',-\[\]]+)?\}\}/g;
        return msg.replace(pattern, function(match, key) {
            var value = values[key] || eval('(values.' +key+')');
            return Object.prototype.toString.call(filter) === "[object Function]" ? filter(value, key) : value;
        }); 
    },
    queryString: function () {
        var querystring = {};
        window.location.href.replace(
            new RegExp("([^?=&]+)(=([^&]*))?", "g"),
            function($0, $1, $2, $3) { querystring[$1] = $3; }
        );
        return querystring;
    },
    Msg: {
        info: function (msg) {
            alert(msg);
        },
        error: function (msg) {
            alert(msg);
        },
        success: function (msg) {
            alert(msg);
        }
    }
};

/*直接使用seajs的发放释放命名空间与difine*/
var oh = _seajsHost.seajs;
window.define = _seajsHost.define;
delete _seajsHost;

var _base = new Base();
for(var n in _base){
    oh[n] = _base[n];
}

oh.config({
  alias: {
    // 'es5-safe': 'es5-safe/0.9.2/es5-safe',
    // 'json': 'json/1.0.1/json',
    'jquery': 'seajs-lib/jquery-1.7.2',
    'jquery.ui': 'seajs-lib/jquery.ui',
    'jquery.ui.widget': 'seajs-lib/jquery.ui.widget',
    'jquery.fileupload': 'seajs-lib/jquery.fileupload',
    'jquery.poshytip': 'seajs-lib/jquery.poshytip.min',
    'jquery.tagsinput': 'seajs-lib/jquery.tagsinput',
    'jquery.pagination': 'seajs-lib/jquery.pagination',
    'jquery.editable': 'seajs-lib/jquery.editable',
    'scrolltotop': 'seajs-lib/scrolltotop',
    'bootstrap': 'seajs-lib/bootstrap.min'
  },
  // preload: [
  //   Function.prototype.bind ? '' : 'es5-safe',
  //   this.JSON ? '' : 'json'
  // ],
  debug: false,
  // map: [
  //   ['http://example.com/js/app/', 'http://localhost/js/app/']
  // ],
  base: 'http://'+ window.location.host +'/assets/js/',
  charset: 'utf-8'
});

window.oh = oh;

})();



(function(){
    oh.use(['jquery', 'jquery.poshytip', 'bootstrap', 'scrolltotop'], function ($) {
        /**
         * 提示信息
         */
        var hideAt = 3000, 
            modalTpml = '<div class="modal msgmodal {{itype}}">\
  <div class="modal-header">\
    <button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>\
    <h3>{{title}}</h3>\
  </div>\
  <div class="modal-body">\
    <p>{{body}}</p>\
  </div>\
</div>';
        function showMsgModal (itype, title, body, timeout) {
            var md = $(oh.tpFormat(modalTpml, {title:title, body:body, itype:itype}));
            md.appendTo('body');
            var toutId = 0,
                rm = function () {
                    clearTimeout(toutId);
                    md.fadeOut('fast', function() {
                        md.remove();
                    });
                };
            md.find('.close').click(rm);
            toutId = setTimeout(rm, timeout);
        }
        oh.Msg.info = function (msg) {
            showMsgModal('info', '=_=', msg, hideAt);
        };
        oh.Msg.error = function (msg) {
            showMsgModal('error', '&gt;_&lt;!!!', msg, hideAt);
        };
        oh.Msg.success = function (msg) {
            showMsgModal('success', '^_^', msg, hideAt);
        };

        $(document).tooltip({
          selector: "a[rel=tooltip]"
        });

        /**
         * 详细信息浮动提示框
         */
        var popinfoCache = {};
        $('.a-pop-info').poshytip({
            liveEvents: true,
            showTimeout: 500, // hover 1秒才会触发显示
            className: 'tip-yellowsimple',
            alignTo: 'target',
            alignY: 'top',
            alignX: 'center',
            offsetX: 0,
            offsetY: 5,
            fade: false,
            slide: false,
            allowTipHover: true,
            content: function(updateCallback) {
                var url = $(this).attr('data-infourl');
                if (!url) { return '' }

                var t = $(this),
                // d = t.data('a-pop-data');
                d = popinfoCache[url];
                if (d) {
                    // return oh.tpFormat(userPopinfoTmpl, d);
                    return d;
                }

                $.ajax({
                    url: url,
                    type: "get",
                    cache: false,
                    // dataType: "json",
                    success: function (r) {
                        if (r && r.indexOf('<div') === 0){
                            popinfoCache[url] = r;
                            updateCallback(r);
                            // t.poshytip('show');
                        }
                    }
                });
                return '加载中...';
            }
        });
        
        /**
         * link投票
         */
        $(document.body).on('click', '.ulitem .vote a', function () {
            var t = $(this), vt = 0;
            if (t.hasClass('up')) {
                vt = 1;
            } else if (t.hasClass('down')) {
                vt = 2;
            } else {
                return;
            }
            var lid = t.closest('.ulitem').attr('data-id');
            if (!lid) {return}
            $.ajax({
                url: '/vote/link/' + lid + '/' + vt,
                type: "post",
                dataType: "json",
                beforeSend: function(xhr){
                    t.attr('disabled', true);
                },
                success: function(data, textStatus){
                    if (data && data.Success === true) {
                        var p = t.closest('.ulitem');
                        p.find('.vote a').removeClass('on');
                        t.addClass('on');
                        p.find('.vote .num').text(data.VoteNum);
                    } else if (data) {
                        if (data.needLogin) {
                            oh.toLogin();
                        } else {
                            oh.Msg.error( data.Errors ? data.Errors : '请求出错，请稍后重试');
                        }
                    } else {
                        oh.Msg.error('请求出错，请稍后重试');
                    }
                },
                complete: function(xhr, status){
                    t.removeAttr('disabled');
                },
                error: function(){
                    oh.Msg.error('请求出错，请稍后重试');
                }
            });
        });

        /**
         * 关注(包括话题和用户)
         */
        function doFollow (btn) {
            var ftype = btn.attr('data-ftype'), url = '';
            if (ftype === 'user') {
                url = "/user/" + btn.attr("data-id") + "/follow";
            } else if (ftype === 'topic') {
                url = "/topic/" + btn.attr("data-id") + "/follow";
            } else {
                return;
            }
            btn.attr('disabled', true).text('关注中...')
            $.ajax({
                url: url,
                type: "post",
                dataType: "json",
                success: function (r) {
                    if (r && r.success){
                        btn.text("已关注")
                            .attr('data-atype', 'unfollow')
                            .removeClass('dofollow')
                            .addClass('dounfollow');
                    }else {
                        btn.text("关注");
                        oh.Msg.error(r.errors)
                    }
                },
                complete: function(xhr, status){
                    btn.removeAttr('disabled');
                }
            });
        }
        function doUnFollow (btn) {
            var ftype = btn.attr('data-ftype'), url = '';
            if (ftype === 'user') {
                url = "/user/" + btn.attr("data-id") + "/unfollow";
            } else if (ftype === 'topic') {
                url = "/topic/" + btn.attr("data-id") + "/unfollow";
            } else {
                return;
            }
            btn.attr('disabled', true).text('取消关注中...')
            if (!btn.data('data-otext')) {
                btn.data('data-otext', btn.html())
            }
            $.ajax({
                url: url,
                type: "post",
                dataType: "json",
                success: function (r) {
                    if (r && r.success){
                        btn.html('<i class="icon-plus icon-white"></i> 关注')
                        .attr('data-atype', 'follow')
                        .removeClass('dounfollow').removeClass('btn-danger').removeClass('btn-info')
                        .addClass('dofollow').addClass('btn-primary');
                    }else {
                        btn.html(btn.data('data-otext'));
                        oh.Msg.error(r.errors)
                    }
                },
                complete: function(xhr, status){
                    btn.removeAttr('disabled');
                }
            });
        }
        $(document.body).on('click', '.dofollow, .dounfollow', function () {
            var btn = $(this), atype = btn.attr('data-atype');
            if (atype === 'follow') {
                doFollow(btn);
            } else if (atype === 'unfollow') {
                doUnFollow(btn);
            }
        })
        .on('mouseenter', '.dofollow, .dounfollow', function () {
            var t = $(this);
            if (t.hasClass('dounfollow')) {
                var ot = t.data('data-otext');
                if (!ot) {
                    t.data('data-otext', t.html());
                }
                t.removeClass('btn-info').addClass('btn-danger')
                    .text('取消关注');
            }
        })
        .on('mouseleave', '.dofollow, .dounfollow', function () {
            var t = $(this), ot = t.data('data-otext');
            if (t.hasClass('dounfollow')) {
                t.removeClass('btn-danger').addClass('btn-info')
                    .text('取消关注');
                if (ot) {
                    t.html(ot);
                    t.data('data-otext', '');
                }
            }
        });

        /**
         * 加载更多link
         */
        $("#loadmorelink").click(function (e) {
            var t = $(this), querystring = oh.queryString();
            if (!window.linkLoadedPage) {
                window.linkLoadedPage = 1;
            }
            window.linkLoadedPage++;
            querystring['page'] = window.linkLoadedPage;
            var loadTip = '(正在加载...)';
            $.ajax({
                url: $(this).attr('data-url'), // '/home/loadmorelink',
                type: "get",
                data: querystring,
                dataType: "json",
                beforeSend: function(xhr){
                    t.attr('disabled', true);
                    t.text(t.text() + loadTip);
                },
                success: function(data, textStatus){
                    if (data && data.success === true) {
                        if (data.html) {
                            t.before(data.html);
                        } else {
                            oh.Msg.info('没有更多链接了');
                        }
                        if (!data.hasmore) {
                            t.hide();
                        }
                    } else if (data) {
                        if (data.needLogin) {
                            oh.toLogin();
                        } else {
                            oh.Msg.error( data.errors ? data.errors : '请求出错，请稍后重试');
                        }
                    } else {
                        oh.Msg.error('请求出错，请稍后重试');
                    }
                },
                complete: function(xhr, status){
                    t.removeAttr('disabled');
                    t.text(t.text().replace(loadTip, ''));
                },
                error: function(){
                    oh.Msg.error('请求出错，请稍后重试');
                }
            });
        });

        /**
         * 删除
         * @t: 触发事件的元素，为jquery对象
         * @parent: 删除成功时要移除的元素，为selecter
         */
        function ajaxDelItem (t, parent) {
            var url = t.attr('data-url');
            $.ajax({
                url: url,
                type: "post",
                dataType: "json",
                beforeSend: function(xhr){
                    t.attr('disabled', true);
                },
                success: function(data, textStatus){
                    if (data && data.success === true) {
                        oh.Msg.success("已成功删除");
                        t.closest(parent).fadeOut('slow', function () {
                            $(this).remove();
                        });
                    } else if (data) {
                        if (data.needLogin) {
                            oh.toLogin();
                        } else {
                            oh.Msg.error( data.errors ? data.errors : '请求出错，请稍后重试');
                        }
                    } else {
                        oh.Msg.error('请求出错，请稍后重试');
                    }
                },
                complete: function(xhr, status){
                    t.removeAttr('disabled');
                },
                error: function(){
                    oh.Msg.error('请求出错，请稍后重试');
                }
            });
        };
        // 删除链接
        $(document.body).on('click', '.link-del', function () {
            if (!window.confirm('你确定要删除该链接吗？')) { return; }
            var btn = $(this);
            ajaxDelItem(btn, '.ulitem');
        });

        /**
         * 检查新提醒
         */
        var checkRemindId = 0;
        var reminder = '<div id="reminder" class="alert alert-block hide">\
  <button type="button" class="close" data-dismiss="alert">&times;</button>\
  <div class="c"></div>\
</div>';
        function checkRemind () {
            clearTimeout(checkRemindId);
            $.ajax({
                url: '/user/remind',
                dataType: 'json',
                success: function (r) {
                  if (r && r.success && r.remind) {
                      var rele = $('#reminder');
                      if (!rele.length) {
                          rele = $(reminder).appendTo('body');
                          rele.find('.close').click(function (e) {
                              rele.hide();
                          });
                      }
                      var c = rele.find('.c');
                      c.html('');
                      if (r.remind.Comments) {
                          c.append('<a href="/comment/inbox">' + r.remind.Comments + ' 条新评论</a>');
                      }
                      if (r.remind.Fans) {
                          c.append('<a href="/comment/inbox">' + r.remind.Fans + ' 位新粉丝</a>');
                      }
                      if (c.html()) { rele.show() } else { rele.hide() }
                  } else {
                      r = r || {};
                      window.console && console.error('获取提醒信息出错：'+r.errors);
                  }
                  checkRemindId = setTimeout(checkRemind, 120*1000);
                }
            });
        };
        checkRemind();

        /**
         * 分享按钮工具
         */
        $('#tool-bookmark-share-btn').hover(function () {
            $(this).find('.drag-arrow').show();
        }, function () {
            $(this).find('.drag-arrow').hide();
        });

        /**
         * 扩散
         */
        $(document.body).on('click', '.spread', function () {
            oh.Msg.info('功能努力开发中...');
        });

    });

})();
